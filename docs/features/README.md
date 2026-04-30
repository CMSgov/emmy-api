# Features

This directory organizes feature documentation by domain so readers can quickly
navigate core application behavior, supporting infrastructure, security
controls, and resilience patterns without scanning a single flat list.

## Core

| Component | Purpose | Functionality |
|---|---|---|
| [Verification API v0 Contract](core/edu-openapi-spec.md) | Document the versioned public API contract defined in OpenAPI and JSON Schema. | Covers the v0 operations, shared request schema, response shapes, and contract-governance expectations. |
| [NSC Education](core/nsc-education.md) | Describe NSC-based education verification flow. | Covers request/response behavior, service boundaries, and operational caveats for the education path. |

## Infrastructure

| Component | Purpose | Functionality |
|---|---|---|
| [Redis](infrastructure/redis.md) | Document Redis runtime dependency and client behavior. | Covers config defaults, connection/pool settings, instrumentation hooks, ping usage, and operational failure modes. |
| [Orchestrion](infrastructure/orchestrion.md) | Document Datadog Orchestrion integration. | Covers build-time instrumentation, environment configuration, and agent connectivity. |

## Security

| Component | Purpose | Functionality |
|---|---|---|
| [Request Identity Handling](security/cognito-auth.md) | Document the current request identity middleware and local skip-auth behavior. | Covers subject extraction, local identity injection, and auth-related caveats on the current branch. |

## Resilience

| Component | Purpose | Functionality |
|---|---|---|
| [Circuit Breaker](resilience/circuit-breaker.md) | Document request admission and breaker behavior. | Covers Redis-backed `Allow` checks, fail-open behavior, and current transition-hook limitations. |
