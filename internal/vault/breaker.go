package vault

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

// BreakerConfig holds circuit breaker configuration
type BreakerConfig struct {
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
}

// WithCircuitBreaker wraps the client with a circuit breaker
func (c *Client) WithCircuitBreaker(config BreakerConfig, onStateChange func(string, string)) {
	settings := gobreaker.Settings{
		Name:        "vault-client",
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if onStateChange != nil {
				onStateChange(from.String(), to.String())
			}
		},
	}

	c.breaker = gobreaker.NewCircuitBreaker(settings)
}

// executeWithBreaker executes a function with circuit breaker protection
func (c *Client) executeWithBreaker(fn func() (interface{}, error)) (interface{}, error) {
	if c.breaker == nil {
		return fn()
	}

	result, err := c.breaker.Execute(fn)
	if err != nil {
		return nil, fmt.Errorf("circuit breaker: %w", err)
	}
	return result, nil
}
