package proxy

import (
	"time"

	"github.com/sony/gobreaker"
)

var circuitBreakers = make(map[string]*gobreaker.CircuitBreaker)

func GetCricuitBreaker(serviceName string) *gobreaker.CircuitBreaker {
	if cb, exists := circuitBreakers[serviceName]; exists {
		return cb
	}

	settings := gobreaker.Settings{
		Name:        serviceName,
		Timeout:     30 * time.Second, // How long to stay open
		MaxRequests: 5,                // How many requests allowed in Half-Open state
		Interval:    0,                // No periodic reset
		IsSuccessful: func(err error) bool {
			return err == nil
		},
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip after 5 failures out of 10 requests
			return counts.ConsecutiveFailures > 5
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)
	circuitBreakers[serviceName] = cb

	return cb

}
