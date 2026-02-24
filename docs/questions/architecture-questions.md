# Architecture Questions

## Q1
- **Question:** Should breaker state transitions be fully managed in middleware (`Allow` + `OnSuccess`/`OnFailure`) or delegated to handlers/services?
- **Context:** Current middleware only performs admission check.
- **Why it matters:** Placement affects consistency of resilience behavior and implementation complexity.
- **Suggested investigation direction:** Prototype middleware-wrapped result classification with unit tests.
- **Owner:** TBD
- **Target decision date:** TBD

## Q2
- **Question:** Should Redis remain a hard dependency for readiness, or should readiness degrade gracefully when Redis is unavailable?
- **Context:** `/status` and many tests assume Redis connectivity.
- **Why it matters:** Deployment resilience and startup behavior differ by dependency policy.
- **Suggested investigation direction:** Define liveness/readiness split and platform expectations.
- **Owner:** TBD
- **Target decision date:** TBD

## Q3
- **Question:** Is the current package layering sufficient, or is a stricter ports/adapters split required now?
- **Context:** Service uses interface boundaries but still has thin, integration-centric layering.
- **Why it matters:** Future provider expansion may demand clearer boundary contracts.
- **Suggested investigation direction:** Evaluate growth scenarios and refactor cost vs benefit.
- **Owner:** TBD
- **Target decision date:** TBD
