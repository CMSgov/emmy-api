package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultDialTimeout  = 2 * time.Second
	defaultReadTimeout  = 2 * time.Second
	defaultWriteTimeout = 2 * time.Second
	defaultPoolTimeout  = 2 * time.Second

	defaultPoolSize     = 20
	defaultMinIdleConns = 2
)

type Config struct {
	// Typically "localhost:6379"
	Addr               string
	Password           string
	DB                 int
	UseTLS             bool
	InsecureSkipVerify bool
}

func NewClient(c Config, logger *slog.Logger) *redis.Client {
	if logger == nil {
		logger = slog.Default()
	}

	// Attach component metadata once
	logger = logger.With(
		slog.String("component", "redis"),
		slog.String("addr", c.Addr),
		slog.Int("db", c.DB),
	)

	opts := &redis.Options{
		Addr:         c.Addr,
		Password:     c.Password, // No password
		DB:           c.DB,
		DialTimeout:  defaultDialTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		PoolTimeout:  defaultPoolTimeout,
		PoolSize:     defaultPoolSize,
		MinIdleConns: defaultMinIdleConns,
	}
	if c.UseTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: c.InsecureSkipVerify, //nolint:gosec // Configurable for local/dev environments and controlled by config.
		}
	}

	logger.Info("initializing redis client")

	rdb := redis.NewClient(opts)

	return rdb
}

func Ping(ctx context.Context, rdb *redis.Client) error {
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	return nil
}
