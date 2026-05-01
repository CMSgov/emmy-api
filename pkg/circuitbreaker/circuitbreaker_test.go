package circuitbreaker

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultOptions_AreUsable(t *testing.T) {
	opts := DefaultOptions()

	assert.Positive(t, opts.FailureThreshold)
	assert.Greater(t, opts.FailWindow, time.Duration(0))
	assert.Greater(t, opts.OpenCoolDown, time.Duration(0))
	assert.Greater(t, opts.HalfOpenLease, time.Duration(0))
}

func TestNewBreaker_UsesDefaults(t *testing.T) {
	rdb := newTestRedisClient(t)

	logger := slog.New(slog.DiscardHandler)

	breaker := NewRedisBreaker(rdb, "test", Options{}, logger)

	def := DefaultOptions()
	assert.Equal(t, def.FailureThreshold, breaker.opts.FailureThreshold)
	assert.Equal(t, def.FailWindow, breaker.opts.FailWindow)
	assert.Equal(t, def.OpenCoolDown, breaker.opts.OpenCoolDown)
	assert.Equal(t, def.HalfOpenLease, breaker.opts.HalfOpenLease)
	assert.Equal(t, def.FailOpen, breaker.opts.FailOpen)
	assert.Equal(t, def.Prefix, breaker.opts.Prefix)
}
