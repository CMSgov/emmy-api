package redis

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_Ping_Set_Get(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	rdb := NewClient(Config{Addr: addr}, logger)

	err := Ping(ctx, rdb)

	require.NoErrorf(t, err, "Ping(ctx, rdb) returned an error: %v", err)

	// test namespace
	key := "cb:test:foo"

	err = rdb.Set(ctx, key, "bar", 5*time.Second).Err()

	require.NoErrorf(t, err, `rdb.Set(ctx, key, "bar", 5*time.Second).Err() return an error: %v`, err)

	val, err := rdb.Get(ctx, key).Result()
	result := val

	expected := "bar"

	assert.Equalf(t, expected, result, "Expected: %q; Got: %q", expected, result)

	_ = rdb.Del(ctx, key).Err
}

func NewClient_Works(t *testing.T) {

	t.Setenv("REDIS_ADDR", "localhost:6379")

	addr := os.Getenv("REDIS_ADDR")

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	rdb := NewClient(Config{Addr: addr}, logger)

	assert.NotNil(t, rdb, "NewClient(Config{Addr: addr}) should not be nil")
}
