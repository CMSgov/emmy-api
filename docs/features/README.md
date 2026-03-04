# Features

This directory organizes feature documentation by domain so readers can quickly navigate core application behavior, supporting infrastructure, security controls, and resilience patterns without scanning a single flat list.

## Core

| Component | Purpose | Functionality |
|---|---|---|
| [NSC Education](core/nsc-education.md) | Describe NSC-based education verification flow. | Covers request/response behavior, service boundaries, and operational caveats for the education path. |
| [Fiber Status Endpoint](core/fiber-status-endpoint.md) | Document status endpoint behavior in Fiber. | Covers `/status` logic, auth gating, Redis ping path, and known wiring caveats. |

## Infrastructure

| Component | Purpose | Functionality |
|---|---|---|
| [Redis](infrastructure/redis.md) | Document Redis runtime dependency and client behavior. | Covers config defaults, connection/pool settings, instrumentation hooks, ping usage, and operational failure modes. |
| [OpenTelemetry Logging](infrastructure/opentelemetry-logging.md) | Document observability and structured logging stack. | Covers OTel setup, logger fanout, instrumentation behavior, and telemetry caveats. |

## Security

| Component | Purpose | Functionality |
|---|---|---|
| [Cognito Auth](security/cognito-auth.md) | Document Cognito access-token validation middleware. | Covers token header/claims checks, local context propagation, and auth-related edge cases. |

## Resilience

| Component | Purpose | Functionality |
|---|---|---|
| [Circuit Breaker](resilience/circuit-breaker.md) | Document request admission and breaker behavior. | Covers Redis-backed `Allow` checks, fail-open behavior, and current transition-hook limitations. |
