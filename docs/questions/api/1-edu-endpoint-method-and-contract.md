---
id: api-1
category: api
title: "Should `/api/edu` remain a GET endpoint, or become POST with caller-supplied payload?"
status: open
owner: TDB
target_decision_date: TBD
priority: high
---

## Q: Should `/api/edu` remain a GET endpoint, or become POST with caller-supplied payload?

### Context
Current implementation issues an NSC submit call using hardcoded request fields.

### Why It Matters
HTTP semantics, security posture, and client contract design depend on this decision.

### Suggested Investigation Direction
Align with product/API consumer expectations and define versioned contract.

## A:
