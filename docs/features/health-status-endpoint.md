# Feature: Health Status Endpoint

## Feature Overview
Provides a lightweight Redis-backed health endpoint at `GET /status`.

## Business Logic
- Wrap endpoint with circuit breaker middleware.
- Create request-scoped context with 2-second timeout.
- Execute Redis `PING` through shared redis helper.
- Return `200 OK` on success; propagate error through Fiber handler chain on failure.

## Package Location
- `api/routes/status_router.go`
- `api/handlers/status_handler.go`
- `pkg/redis/redis.go`

## Key Structs and Interfaces
- `GetRDBStatus(rdb *redis.Client) fiber.Handler`
- `redis.Ping(ctx, rdb)`
- `WithCircuitBreaker(...)`

## Real Code Excerpt
```go
ctx, cancel := context.WithTimeout(c.Context(), contextTimeout*time.Second)
defer cancel()

err := redisLocal.Ping(ctx, rdb)
if err != nil {
    return err
}

return c.SendStatus(fiber.StatusOK)
```

## Edge Cases Handled Today
- Timeout-bound ping prevents indefinite blocking.
- Redis failures bubble to Fiber error handling.
- Circuit breaker may return `503` before ping is attempted.

## Performance and Operational Considerations
- Endpoint is cheap but still network-bound on Redis.
- Suitable for dependency-aware health checks, not pure process liveness.
- Circuit breaker wrapping can mask direct Redis error semantics with uniform `503`.

## Future Improvements
- Split `/live` (process-only) and `/ready` (dependency-aware) semantics.
- Return structured JSON with dependency states.
- Add optional shallow health mode bypassing breaker for diagnostics.

## Assumptions
- **High confidence:** `/status` is intended as readiness/dependency check rather than a simple heartbeat.
