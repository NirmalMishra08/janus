package proxy

import "sync/atomic"

type RoundRobinBalancer struct {
	instances   []string
	index       uint64
	serviceName string
}

func NewRoundRobinBalancer(instances []string, serviceName string) *RoundRobinBalancer {
	return &RoundRobinBalancer{
		serviceName: serviceName,
		instances: instances,
	}
}

func (lb *RoundRobinBalancer) Next() string {
	if len(lb.instances) == 0 {
		return ""
	}

	for i := 0; i < len(lb.instances); i++ {
		idx := atomic.AddUint64(&lb.index, 1) % uint64(len(lb.instances))
		instance := lb.instances[idx]

		if healthChecker.IsHealthy(lb.serviceName, instance) {
			return instance
		}
	}
	return lb.instances[0]
}
