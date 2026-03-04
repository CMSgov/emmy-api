# Feature: Redis

## Feature Overview
Provides the shared Redis integration used for runtime connectivity checks and distributed state in this service.

## Business Logic
- Build a Redis client from `redis.Config` values (address, password, database).
- Apply deterministic connection/pool settings for dial/read/write/pool timeouts and pool sizing.
- Attach OpenTelemetry tracing and metrics instrumentation hooks to the client.
- Expose a `Ping` helper used by startup and status-check paths.

## Package Location
- `pkg/redis/redis.go`
- `pkg/core/config.go`
- `api/handlers/status_handler.go`
- `main.go`

## Key Structs and Interfaces
- `redis.Config`
- `redis.NewClient(c Config, logger *slog.Logger) *redis.Client`
- `redis.Ping(ctx context.Context, rdb *redis.Client) error`
- `core.RedisConfig`

## Real Code Excerpt
```go
opts := &redis.Options{
    Addr:         c.Addr,
    Password:     c.Password,
    DB:           c.DB,
    DialTimeout:  defaultDialTimeout,
    ReadTimeout:  defaultReadTimeout,
    WriteTimeout: defaultWriteTimeout,
    PoolTimeout:  defaultPoolTimeout,
    PoolSize:     defaultPoolSize,
    MinIdleConns: defaultMinIdleConns,
}

rdb := redis.NewClient(opts)

if err := redisotel.InstrumentTracing(rdb); err != nil {
    logger.Warn("Otel Tracing Instrumentation Failed", "err", err)
}
if err := redisotel.InstrumentMetrics(rdb); err != nil {
    logger.Warn("Otel Metrics instrumentation Failed", "err", err)
}
```

## Edge Cases Handled Today
- Nil logger input is handled by falling back to `slog.Default()`.
- OTel instrumentation failures degrade to warning logs; client creation continues.
- `Ping` failures propagate to callers for explicit failure handling.

## Performance and Operational Considerations
- Connection defaults in code:
  - `DialTimeout`, `ReadTimeout`, `WriteTimeout`, `PoolTimeout` = 2s
  - `PoolSize` = 20
  - `MinIdleConns` = 2
- Config defaults from `core.DefaultConfig()`:
  - `REDIS_ADDR=localhost:6379`, empty password, `REDIS_DB=0`
- Redis connectivity is currently a hard startup dependency in `main.run` (`redis.Ping` failure aborts startup).

## Future Improvements
- Add explicit startup/readiness policy split so Redis dependency can be tuned by environment.
- Add retry/backoff strategy for startup ping in non-production local workflows.
- Add dedicated dashboards/alerts for Redis pool pressure and ping latency.

## Assumptions
- **High confidence:** Redis is an operational dependency for current startup and status-check behavior.
- **High confidence:** There is an active wiring caveat on `main` where status route setup may receive nil Redis via `api.New` config path until Redis injection is corrected there.
