# Feature: Health Status Endpoint

## Feature Overview
Intended behavior is a lightweight Redis-backed health endpoint at `GET /status`.
Current `main` wiring has a known caveat: `/status` is registered via `api.New` using `api.Config.Redis`, but `main` does not inject that field.

## Business Logic
- Intended route behavior:
  - Wrap endpoint with circuit breaker middleware.
  - Create request-scoped context with 2-second timeout.
  - Execute Redis `PING` through shared redis helper.
  - Return `200 OK` on success; propagate error through Fiber handler chain on failure.
- Current caveat:
  - Route wiring path in `main` currently passes nil Redis into status route setup.

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
- Redis failures from handler path bubble to Fiber error handling (`500` by default).
- Circuit breaker may return `503` before ping is attempted when breaker denies.
- `/status` is auth-gated when `SKIP_AUTH=false` (global Cognito middleware).

## Performance and Operational Considerations
- Endpoint is cheap but still network-bound on Redis.
- Suitable for dependency-aware health checks, not pure process liveness.
- Failure semantics differ by path: breaker deny returns `503`, while handler Redis ping errors resolve through Fiber global error handler (`500` for non-`fiber.Error`).

## Future Improvements
- Split `/live` (process-only) and `/ready` (dependency-aware) semantics.
- Return structured JSON with dependency states.
- Add optional shallow health mode bypassing breaker for diagnostics.

## Assumptions
- **High confidence:** `/status` is intended as readiness/dependency check rather than a simple heartbeat.
