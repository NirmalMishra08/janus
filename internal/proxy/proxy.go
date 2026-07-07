package proxy

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Service struct {
	ServiceName string
	Instances   []string
	RetryCount  int
	Timeout     time.Duration
}

type Handler struct {
	balancer    *RoundRobinBalancer
	service     Service
	serviceName string
}

func NewHandler(service Service) (http.Handler, error) {
	if len(service.Instances) == 0 {
		return nil, errors.New("no instances are there in service")
	}

	balancer := NewRoundRobinBalancer(service.Instances, service.ServiceName)

	return &Handler{
		balancer:    balancer,
		service:     service,
		serviceName: service.ServiceName,
	}, nil
}

// ServeHTTP with Retry + Health Check awareness
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cb := GetCricuitBreaker(h.serviceName)

	// Execute request through circuit breaker
	_, err := cb.Execute(func() (interface{}, error) {
		var lastErr error

		for attempt := 0; attempt <= h.service.RetryCount; attempt++ {
			targetURL := h.balancer.Next()
			if targetURL == "" {
				return nil, errors.New("no healthy instances")
			}

			target, _ := url.Parse(targetURL)
			proxy := httputil.NewSingleHostReverseProxy(target)

			errChan := make(chan error, 1)
			go func() {
				proxy.ServeHTTP(w, r)
				errChan <- nil
			}()

			select {
			case <-errChan:
				return nil, nil // Success
			case <-time.After(h.service.Timeout):
				lastErr = errors.New("timeout")
				continue
			}
		}
		return nil, lastErr
	})

	if err != nil {
		log.Printf("Circuit Breaker blocked or failed for %s: %v", h.serviceName, err)
		http.Error(w, "Service unavailable", http.StatusBadGateway)
	}
}