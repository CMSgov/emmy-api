---
id: product-1
category: product
title: "What is the expected error contract shape for external clients?"
status: open
owner: TBD
target_decision_date: TBD
priority: medium
---

## Q: What is the expected error contract shape for external clients?

### Context
Current errors may be plain text via Fiber error handler.

### Why It Matters
Client interoperability and observability improve with stable machine-readable errors.

### Suggested Investigation Direction
Define and adopt JSON error schema with code/message/details fields.

## A:
