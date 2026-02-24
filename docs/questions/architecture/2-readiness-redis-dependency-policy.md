---
id: architecture-2
category: architecture
title: "Should Redis remain a hard dependency for readiness, or should readiness degrade gracefully when Redis is unavailable?"
status: open
owner: TBD
target_decision_date: TBD
priority: high
---

## Q: Should Redis remain a hard dependency for readiness, or should readiness degrade gracefully when Redis is unavailable?

### Context
`/status` and many tests assume Redis connectivity.

### Why It Matters
Deployment resilience and startup behavior differ by dependency policy.

### Suggested Investigation Direction
Define liveness/readiness split and platform expectations.

## A:
