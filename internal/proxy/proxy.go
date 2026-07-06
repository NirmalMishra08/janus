package proxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Service struct {
	ServiceName string
	Instances  []string
	RetryCount int
	Timeout    time.Duration
}

type Handler struct {
	balancer *RoundRobinBalancer
	proxy    *httputil.ReverseProxy
	service  Service
}

func NewHandler(service Service) (http.Handler, error) {
	if len(service.Instances) == 0 {
		return nil, errors.New("no instances are there in service")
	}

	balancer := NewRoundRobinBalancer(service.Instances, service.ServiceName)

	director := func(req *http.Request) {
		targetUrl := balancer.Next()
		if targetUrl == "" {
			return
		}
		target, err := url.Parse(targetUrl)
		if err != nil {
			return 
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
	}

	proxy := &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadGateway)
		},
	}

	return &Handler{
		balancer: balancer,
		proxy:    proxy,
		service:  service,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}
