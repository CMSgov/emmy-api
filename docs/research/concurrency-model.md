# Research: Concurrency Model

## Problem Statement
The service must handle concurrent HTTP requests safely while maintaining predictable lifecycle behavior and shared middleware state.

## Alternatives Considered
- Fiber-native concurrency with shared mutable state guarded by mutexes.
- Per-request stateless components with no shared in-memory registry.
- Externalized breaker registry only in Redis (no in-process map).

## Trade-offs
- Current approach (mutex-guarded map + Redis state):
  - Pros: low overhead, stable breaker identity per route, simple implementation.
  - Cons: complexity in lock management, local cache lifecycle concerns.
- Fully stateless middleware:
  - Pros: simpler in-memory model.
  - Cons: more Redis object churn and repeated construction overhead.

## Why Current Approach Was Selected (Inferred)
The code favors a pragmatic middle ground: shared in-process breaker lookup for efficiency, with Redis as distributed state authority.

## Benchmarks / Status
- Not available in repository.
- No microbenchmarks for middleware lock contention or Redis key access path are currently present.

## References
- `main.go` (`runServer` goroutine + graceful shutdown select)
- `api/middleware/middleware.go` (`sync.RWMutex`, breaker map)
- `pkg/circuitbreaker/redis_breaker.go`

## Assumptions
- **Medium confidence:** Current traffic profile is modest enough that lock contention is not yet a bottleneck.
