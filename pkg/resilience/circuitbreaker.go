package resilience

import (
	"errors"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type State string

const (
	StateClosed   State = "CLOSED"
	StateOpen     State = "OPEN"
	StateHalfOpen State = "HALF_OPEN"

	defaultFailureThreshold  = 3
	defaultOpenTimeout       = 30 * time.Second
	defaultHalfOpenMaxProbes = 1
)

type CircuitBreakerOptions struct {
	FailureThreshold  int
	OpenTimeout       time.Duration
	HalfOpenMaxProbes int
}

func DefaultCircuitBreakerOptions() CircuitBreakerOptions {
	return CircuitBreakerOptions{
		FailureThreshold:  defaultFailureThreshold,
		OpenTimeout:       defaultOpenTimeout,
		HalfOpenMaxProbes: defaultHalfOpenMaxProbes,
	}
}

type CircuitBreaker struct {
	mu sync.Mutex

	state State

	consecutiveFailures int
	halfOpenInFlight    int
	openedAt            time.Time

	opts CircuitBreakerOptions
}

func NewCircuitBreaker(opts CircuitBreakerOptions) *CircuitBreaker {
	if opts.FailureThreshold <= 0 {
		opts.FailureThreshold = defaultFailureThreshold
	}
	if opts.OpenTimeout <= 0 {
		opts.OpenTimeout = defaultOpenTimeout
	}
	if opts.HalfOpenMaxProbes <= 0 {
		opts.HalfOpenMaxProbes = defaultHalfOpenMaxProbes
	}

	return &CircuitBreaker{
		state: StateClosed,
		opts:  opts,
	}
}

func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.refreshStateLocked(time.Now())

	switch cb.state {
	case StateClosed:
		return nil
	case StateOpen:
		return ErrCircuitOpen
	case StateHalfOpen:
		if cb.halfOpenInFlight >= cb.opts.HalfOpenMaxProbes {
			return ErrCircuitOpen
		}

		cb.halfOpenInFlight++
		return nil
	default:
		return ErrCircuitOpen
	}
}

func (cb *CircuitBreaker) OnSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		cb.consecutiveFailures = 0
	case StateHalfOpen:
		if cb.halfOpenInFlight > 0 {
			cb.halfOpenInFlight--
		}
		if cb.halfOpenInFlight == 0 {
			cb.state = StateClosed
			cb.consecutiveFailures = 0
			cb.openedAt = time.Time{}
		}
	}
}

func (cb *CircuitBreaker) OnFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case StateClosed:
		cb.consecutiveFailures++
		if cb.consecutiveFailures >= cb.opts.FailureThreshold {
			cb.openLocked(now)
		}
	case StateHalfOpen:
		if cb.halfOpenInFlight > 0 {
			cb.halfOpenInFlight--
		}
		cb.openLocked(now)
	case StateOpen:
		// Extend the cooldown window on repeated failures while open.
		cb.openedAt = now
	}
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.refreshStateLocked(time.Now())
	return cb.state
}

func (cb *CircuitBreaker) openLocked(now time.Time) {
	cb.state = StateOpen
	cb.openedAt = now
	cb.halfOpenInFlight = 0
}

func (cb *CircuitBreaker) refreshStateLocked(now time.Time) {
	if cb.state != StateOpen {
		return
	}

	if now.Sub(cb.openedAt) < cb.opts.OpenTimeout {
		return
	}

	cb.state = StateHalfOpen
	cb.halfOpenInFlight = 0
}
