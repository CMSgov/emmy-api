# Open Questions

## Q1
- **Question:** Should `/api/edu` remain a GET endpoint, or become POST with caller-supplied payload?
- **Context:** Current implementation issues an NSC submit call using hardcoded request fields.
- **Why it matters:** HTTP semantics, security posture, and client contract design depend on this decision.
- **Suggested investigation direction:** Align with product/API consumer expectations and define versioned contract.
- **Owner:** TDB
- **Target decision date:** TBD

## Q2
- **Question:** What is the expected error contract shape for external clients?
- **Context:** Current errors may be plain text via Fiber error handler.
- **Why it matters:** Client interoperability and observability improve with stable machine-readable errors.
- **Suggested investigation direction:** Define and adopt JSON error schema with code/message/details fields.
- **Owner:** TBD
- **Target decision date:** TBD

## Q3
- **Question:** Which telemetry signals are required for production SLOs?
- **Context:** OTel traces/metrics exist, but business-level metrics are limited.
- **Why it matters:** Scaling, alerting, and incident response depend on explicit SLO instrumentation.
- **Suggested investigation direction:** Establish route/provider latency, error rate, and breaker-state metrics.
- **Owner:** TBD
- **Target decision date:** TBD
