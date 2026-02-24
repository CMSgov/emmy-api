---
id: architecture-1
category: architecture
title: "Should breaker state transitions be fully managed in middleware (`Allow` + `OnSuccess`/`OnFailure`) or delegated to handlers/services?"
status: open
owner: TBD
target_decision_date: TBD
priority: medium
---

## Q: Should breaker state transitions be fully managed in middleware (`Allow` + `OnSuccess`/`OnFailure`) or delegated to handlers/services?

### Context
Current middleware only performs admission check.

### Why It Matters
Placement affects consistency of resilience behavior and implementation complexity.

### Suggested Investigation Direction
Prototype middleware-wrapped result classification with unit tests.

## A:
