package resilience

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCircuitBreaker_StateTransitions(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerOptions{
		FailureThreshold:  2,
		OpenTimeout:       20 * time.Millisecond,
		HalfOpenMaxProbes: 1,
	})

	require.Equal(t, StateClosed, breaker.State())

	require.NoError(t, breaker.Allow())
	breaker.OnFailure()
	require.Equal(t, StateClosed, breaker.State())

	breaker.OnFailure()
	require.Equal(t, StateOpen, breaker.State())
	require.ErrorIs(t, breaker.Allow(), ErrCircuitOpen)

	time.Sleep(25 * time.Millisecond)

	require.NoError(t, breaker.Allow())
	require.Equal(t, StateHalfOpen, breaker.State())

	breaker.OnSuccess()
	require.Equal(t, StateClosed, breaker.State())
	require.NoError(t, breaker.Allow())
}

func TestCircuitBreaker_HalfOpenProbeLimit(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerOptions{
		FailureThreshold:  1,
		OpenTimeout:       20 * time.Millisecond,
		HalfOpenMaxProbes: 1,
	})

	breaker.OnFailure()
	require.Equal(t, StateOpen, breaker.State())

	time.Sleep(25 * time.Millisecond)

	require.NoError(t, breaker.Allow())
	require.Equal(t, StateHalfOpen, breaker.State())
	require.ErrorIs(t, breaker.Allow(), ErrCircuitOpen)
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerOptions{
		FailureThreshold:  1,
		OpenTimeout:       20 * time.Millisecond,
		HalfOpenMaxProbes: 1,
	})

	breaker.OnFailure()
	require.Equal(t, StateOpen, breaker.State())

	time.Sleep(25 * time.Millisecond)

	require.NoError(t, breaker.Allow())
	breaker.OnFailure()

	require.Equal(t, StateOpen, breaker.State())
	require.True(t, errors.Is(breaker.Allow(), ErrCircuitOpen))
}

func TestDefaultCircuitBreakerOptions_AreUsable(t *testing.T) {
	opts := DefaultCircuitBreakerOptions()

	assert.Greater(t, opts.FailureThreshold, 0)
	assert.Greater(t, opts.OpenTimeout, time.Duration(0))
	assert.Greater(t, opts.HalfOpenMaxProbes, 0)
}
