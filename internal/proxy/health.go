package proxy

import (
	"net/http"
	"sync"
	"time"
)

type HealthChecker struct {
	services map[string]*ServiceHealth
	mu       sync.RWMutex
}

type ServiceHealth struct {
	Instances map[string]*InstanceHealth
}

type InstanceHealth struct {
	Healthy      bool
	LastChecked  time.Time
	FailureCount int
}

var healthChecker = &HealthChecker{
	services: make(map[string]*ServiceHealth),
}

func StartHealthCheck(services map[string]Service) {
	go func() {
		tick := time.NewTicker(10 * time.Second)
		defer tick.Stop()

		for range tick.C {
			healthChecker.checkAll(services)
		}
	}()
}

func (hc *HealthChecker) checkAll(services map[string]Service) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	for name, service := range services {
		if _, exists := hc.services[name]; !exists {
			hc.services[name] = &ServiceHealth{
				Instances: make(map[string]*InstanceHealth),
			}
		}

		for _, instance := range service.Instances {
			hc.CheckInstance(name, instance)
		}

	}
}

func (hc *HealthChecker) CheckInstance(serviceName string, instance string) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(instance + "/health")
	if resp.StatusCode != http.StatusOK || err != nil {
		if _, ok := hc.services[serviceName].Instances[instance]; !ok {
			hc.services[serviceName].Instances[instance] = &InstanceHealth{}
		}
		hc.services[serviceName].Instances[instance].Healthy = false
		hc.services[serviceName].Instances[instance].FailureCount++
	} else {
		hc.services[serviceName].Instances[instance].Healthy = true
		hc.services[serviceName].Instances[instance].FailureCount = 0

	}

	hc.services[serviceName].Instances[instance].LastChecked = time.Now()
}

func (hc *HealthChecker) IsHealthy(serviceName string, instanceName string) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	if sh, ok := hc.services[serviceName]; ok {
		if ih, ok := sh.Instances[instanceName]; ok {
			return ih.Healthy
		}
	}
	return true // default to healthy if not checked yet
}
