# Feature: Redis Circuit Breaker

## Feature Overview
Protects endpoints from repeated downstream failures by gating requests based on Redis-stored breaker state.

## Business Logic
- Middleware computes breaker key per endpoint (`METHOD + route path`).
- A single breaker instance is lazily reused per endpoint.
- `Allow` reads Redis state key:
  - no key => closed (allow)
  - future half-open timestamp => open (deny)
- `OnFailure` increments failure counter and opens breaker when threshold reached.
- `OnSuccess` clears breaker and failure keys.

## Package Location
- `pkg/circuitbreaker/circuitbreaker.go`
- `pkg/circuitbreaker/redis_breaker.go`
- `api/middleware/middleware.go`
- `api/routes/router.go`
- `api/routes/status_router.go`

## Key Structs and Interfaces
- `Breaker`
- `RedisBreaker`
- `Options`
- `WithCircuitBreaker`

## Real Code Excerpt
```go
if int(fails) >= b.opts.FailureThreshold {
    timeToHalfOpenMs := time.Now().Add(b.opts.OpenCoolDown).UnixMilli()
    stateTTL := b.opts.OpenCoolDown + b.opts.HalfOpenLease
    _ = b.rdb.Set(ctx, stateKey, strconv.FormatInt(timeToHalfOpenMs, 10), stateTTL).Err()
    _ = b.rdb.Del(ctx, failsKey).Err()
}
```

## Edge Cases Handled Today
- Redis unavailable + `FailOpen=true` defaults to allow.
- Redis value parse errors follow fail-open/fail-closed option behavior.
- Missing state key is treated as closed circuit.
- Unknown route path falls back to raw request path for breaker naming.

## Performance and Operational Considerations
- Shared breaker map protected by RWMutex avoids race conditions.
- Redis round-trips occur per request for `Allow` checks.
- Breaker state is distributed across instances via Redis keys.

## Future Improvements
- Call `OnFailure`/`OnSuccess` from middleware wrapper automatically based on handler result.
- Introduce jitter and configurable half-open probe strategy.
- Emit breaker state transition metrics.

## Assumptions
- **Medium confidence:** Current middleware only enforces `Allow`; transition hooks are expected to be integrated later for full circuit behavior.
