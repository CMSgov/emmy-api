# Research: Circuit Breaker State Management

## Problem Statement
The API needs consistent failure-gating behavior across instances while tolerating transient Redis issues.

## Alternatives Considered
- In-memory breaker per process only.
- Redis-backed breaker state (current).
- Service-mesh or gateway-level circuit breaking.

## Trade-offs
- Redis-backed state:
  - Pros: shared state across instances, explicit TTL control, simple operational model.
  - Cons: Redis dependency and additional request-time network calls.
- In-memory only:
  - Pros: low latency, no external dependency.
  - Cons: inconsistent behavior across pods.
- Gateway/service mesh:
  - Pros: centralized resilience policy.
  - Cons: extra platform complexity and reduced app-layer control.

## Why Current Approach Was Selected (Inferred)
Implementation indicates preference for application-owned resilience semantics with distributed state and fail-open behavior when Redis state cannot be determined.

## Benchmarks / Status
- Not available.
- Existing tests validate transitions and key structure but not throughput/latency under load.

## References
- `pkg/circuitbreaker/circuitbreaker.go`
- `pkg/circuitbreaker/redis_breaker.go`
- `api/middleware/middleware.go`
- `pkg/circuitbreaker/redis_breaker_test.go`

## Assumptions
- **High confidence:** Distributed breaker state is required for multi-instance consistency.
