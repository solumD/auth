package circuitbreaker

import (
	"fmt"
	"time"

	"github.com/solumD/auth/internal/logger"

	"github.com/sony/gobreaker"
)

const (
	serviceName = "auth_service"
	maxRequests = 3
	timeout     = 5 * time.Second
)

// New инициализирует новый объект circuit breaker
func New() *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        serviceName,
		MaxRequests: maxRequests,
		Timeout:     timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Info(fmt.Sprintf("Circuit Breaker: %s, changed from %v, to %v\n", name, from, to))
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}
