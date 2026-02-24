# Scalability Plan

## Current Limitations
- No explicit rate limiting or workload shaping.
- No documented horizontal scaling policy for Redis breaker state load.
- Limited endpoint set, but provider calls are synchronous and latency-bound.

## Proposed Improvements
- Add request-level concurrency/rate controls.
- Add caching strategy for repeat verification lookups where policy allows.
- Define Redis capacity and key TTL tuning guidelines.
- Introduce async job flow for long-running verification providers.

## Implementation Steps
1. Capture baseline latency and throughput metrics.
2. Add middleware rate-limit policy per route/client class.
3. Evaluate response caching constraints and data-sensitivity requirements.
4. Define autoscaling signals (CPU, latency, error rate, queue depth if async introduced).
5. Load-test `/status` and verification routes with realistic dependency latency.

## Risks
- Caching may conflict with freshness/compliance requirements.
- Async patterns increase operational complexity and consistency concerns.
- Aggressive throttling could impact legitimate user workflows.

## Estimated Complexity
- **L**: Requires architecture, operations, and policy coordination.

## Assumptions
- **Medium confidence:** Traffic growth and provider fanout will require stronger controls than current synchronous flow.
