package circuitbreaker

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

//nolint:govet // Field order prioritizes readability of runtime dependencies.
type RedisBreaker struct {
	// Logger for observability.
	logger *slog.Logger
	// Redis client used to read and update the circuit state.
	rdb *redis.Client
	// Defines the behavior and timing characteristics of the breaker.
	opts Options
	// Name of the redis circuit breaker used when constructing Redis keys.
	name string
}

var _ Breaker = (*RedisBreaker)(nil)

func NewRedisBreaker(rdb *redis.Client, name string, opts Options, logger *slog.Logger) *RedisBreaker {
	if rdb == nil {
		panic("redis client cannot be nil for RedisBreaker")
	}
	if opts.FailureThreshold <= 0 {
		opts = DefaultOptions()
	}
	if opts.Prefix == "" {
		opts.Prefix = "cb:"
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &RedisBreaker{rdb: rdb, name: name, opts: opts, logger: logger}
}

func (b *RedisBreaker) keys() (string, string) {
	base := b.opts.Prefix + b.name
	return base, base + ":fails"
}

func (b *RedisBreaker) Allow(ctx context.Context) error {
	stateKey, _ := b.keys()

	val, err := b.rdb.Get(ctx, stateKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		if b.opts.FailOpen {
			b.logger.WarnContext(ctx, "Redis GET failed; defaulting to allow(assume closed).", "key", stateKey, "err", err)
			return nil
		}
		return ErrCircuitOpen
	}

	timeToHalfOpenMs, convErr := strconv.ParseInt(val, 10, 64)
	if convErr != nil {
		if b.opts.FailOpen {
			b.logger.WarnContext(ctx, "Invalid redis value; defaulting to allow (assume closed).", "key", stateKey, "value", val, "err", convErr)
			return nil
		}
		return ErrCircuitOpen
	}

	nowMs := time.Now().UnixMilli()

	if nowMs >= timeToHalfOpenMs {
		return nil
	}

	return ErrCircuitOpen
}

func (b *RedisBreaker) OnSuccess(ctx context.Context) {
	stateKey, failsKey := b.keys()
	if err := b.rdb.Del(ctx, stateKey, failsKey).Err(); err != nil {
		b.logger.DebugContext(ctx, "redis DEL failed", "state_key", stateKey, "fails_key", failsKey, "err", err)
	}
}

func (b *RedisBreaker) OnFailure(ctx context.Context) {
	stateKey, failsKey := b.keys()

	fails, err := b.rdb.Incr(ctx, failsKey).Result()
	if err != nil {
		b.logger.DebugContext(ctx, "redis INCR failed", "key", failsKey, "err", err)
		return
	}

	ttl, err := b.rdb.PTTL(ctx, failsKey).Result()
	if err == nil && ttl < 0 {
		if pexpErr := b.rdb.PExpire(ctx, failsKey, b.opts.FailWindow).Err(); pexpErr != nil {
			b.logger.DebugContext(ctx, "redis PEXPIRE failed", "key", failsKey, "err", pexpErr)
		}
	}

	if int(fails) >= b.opts.FailureThreshold {
		timeToHalfOpenMs := time.Now().Add(b.opts.OpenCoolDown).UnixMilli()

		stateTTL := b.opts.OpenCoolDown + b.opts.HalfOpenLease
		if stateTTL <= 0 {
			stateTTL = b.opts.OpenCoolDown
		}

		if setErr := b.rdb.Set(ctx, stateKey, strconv.FormatInt(timeToHalfOpenMs, 10), stateTTL).Err(); setErr != nil {
			b.logger.DebugContext(ctx, "redis SET failed", "key", stateKey, "err", setErr)
		}

		if delErr := b.rdb.Del(ctx, failsKey).Err(); delErr != nil {
			b.logger.DebugContext(ctx, "redis DEL failed", "key", failsKey, "err", delErr)
		}
	}
}
