package circuitbreaker

import (
	"context"
	"io"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	redisClientAddr   = "localhost:6379"
	redisPassword     = ""
	redisDB           = 0
	redisDialTimeout  = 2 * time.Second
	redisReadTimeout  = 2 * time.Second
	redisWriteTimeout = 2 * time.Second
	redisPoolTimeout  = 2 * time.Second
	redisPoolSize     = 20
	redisMinIdleConns = 2

	testFailureThreshold = 5
	testFailWindow       = 10 * time.Second
	testOpenCoolDown     = 30 * time.Second
	testHalfOpenLease    = 5 * time.Second
	testFailOpen         = true
	testPrefix           = "cb:"
)

func newTestRedisClient(t *testing.T) *redis.Client {
	t.Helper()

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisClientAddr,
		Password:     redisPassword,
		DB:           redisDB,
		DialTimeout:  redisDialTimeout,
		ReadTimeout:  redisReadTimeout,
		WriteTimeout: redisWriteTimeout,
		PoolTimeout:  redisPoolTimeout,
		PoolSize:     redisPoolSize,
		MinIdleConns: redisMinIdleConns,
	})

	// Fail fast with a clear error if Redis isn't running
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	require.NoError(t, rdb.Ping(ctx).Err(), "Redis must be running at %s for these tests", redisClientAddr)

	return rdb
}

func newTestBreakerOptions(t *testing.T) Options {
	t.Helper()

	return Options{
		FailureThreshold: testFailureThreshold,
		FailWindow:       testFailWindow,
		OpenCoolDown:     testOpenCoolDown,
		HalfOpenLease:    testHalfOpenLease,
		FailOpen:         testFailOpen,
		Prefix:           testPrefix,
	}
}

func TestNewRedisBreaker(t *testing.T) {
	rdb := newTestRedisClient(t)
	testBreakerOpts := newTestBreakerOptions(t)

	name := "redisBreaker:" + t.Name()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	result := NewRedisBreaker(rdb, name, testBreakerOpts, logger)

	require.NotNil(t, result, "NewRedisBreaker should not return nil")
	assert.Same(t, rdb, result.rdb, "Expected breaker to keep the passed-in redis client instance")
	assert.Equal(t, name, result.name)
	assert.Equal(t, testBreakerOpts, result.opts)
}

func TestNewRedisBreaker_keys(t *testing.T) {
	rdb := newTestRedisClient(t)
	testBreakerOpts := newTestBreakerOptions(t)

	name := "redisBreaker:" + t.Name()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	breaker := NewRedisBreaker(rdb, name, testBreakerOpts, logger)

	resultStateKey, resultFailsKey := breaker.keys()

	expectedStateKey := testPrefix + name
	expectedFailsKey := testPrefix + name + ":fails"

	assert.Equalf(t, expectedStateKey, resultStateKey, "Got: %q; Expected: %q", resultStateKey, expectedStateKey)
	assert.Equalf(t, expectedFailsKey, resultFailsKey, "Got: %q; Expected: %q", resultFailsKey, expectedFailsKey)
}

func TestNewRedisBreaker_Allow(t *testing.T) {
	rdb := newTestRedisClient(t)
	testBreakerOpts := newTestBreakerOptions(t)

	name := "redisBreaker:" + t.Name()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	breaker := NewRedisBreaker(rdb, name, testBreakerOpts, logger)

	ctx := context.Background()

	stateKey, failsKey := breaker.keys()
	t.Cleanup(func() { _ = rdb.Del(ctx, stateKey, failsKey).Err() })
	_ = rdb.Del(ctx, stateKey, failsKey).Err()

	err := breaker.Allow(ctx)
	require.NoErrorf(t, err, "The Allow method returned an error: %v", err)
}

func TestRedisBreaker_OnFailure_TransitionsToOpen(t *testing.T) {
	rdb := newTestRedisClient(t)

	opts := newTestBreakerOptions(t)
	opts.FailureThreshold = 2

	ctx := context.Background()

	name := "redisBreaker:" + t.Name()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	breaker := NewRedisBreaker(rdb, name, opts, logger)

	stateKey, failsKey := breaker.keys()
	t.Cleanup(func() { _ = rdb.Del(ctx, stateKey, failsKey).Err() })
	_ = rdb.Del(ctx, stateKey, failsKey).Err()

	breaker.OnFailure(ctx)

	fails, err := rdb.Get(ctx, failsKey).Int64()
	require.NoError(t, err)
	require.Equal(t, int64(1), fails)

	_, err = rdb.Get(ctx, stateKey).Result()
	require.ErrorIs(t, err, redis.Nil)

	breaker.OnFailure(ctx)

	val, err := rdb.Get(ctx, stateKey).Result()
	require.NoError(t, err)

	timeToHalfOpenMs, err := strconv.ParseInt(val, 10, 64)
	require.NoError(t, err)
	require.Greater(t, timeToHalfOpenMs, time.Now().UnixMilli())

	exists, err := rdb.Exists(ctx, failsKey).Result()
	require.NoError(t, err)
	require.Equal(t, int64(0), exists)

	err = breaker.Allow(ctx)
	require.ErrorIs(t, err, ErrCircuitOpen)
}
