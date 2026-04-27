package circuitbreaker

import (
	"context"
	"errors"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

const (
	defaultFailureThreshold = 5
	defaultFailWindow       = 10
	defaultOpenCooldown     = 30
	defaultHalfOpenLease    = 5
	defaultFailOpen         = true
	defaultPrefix           = "cb:"
)

type Breaker interface {
	Allow(ctx context.Context) error
	OnSuccess(ctx context.Context)
	OnFailure(ctx context.Context)
}

type Options struct {
	Prefix           string
	FailureThreshold int
	FailWindow       time.Duration
	OpenCoolDown     time.Duration
	HalfOpenLease    time.Duration
	FailOpen         bool
}

func DefaultOptions() Options {
	return Options{
		FailureThreshold: defaultFailureThreshold,
		FailWindow:       defaultFailWindow * time.Second,
		OpenCoolDown:     defaultOpenCooldown * time.Second,
		HalfOpenLease:    defaultHalfOpenLease * time.Second,
		FailOpen:         defaultFailOpen,
		Prefix:           defaultPrefix,
	}
}
