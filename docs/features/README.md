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
| [OpenTelemetry](infrastructure/opentelemetry.md) | Document app and collector telemetry integration. | Covers OTel initialization, middleware tracing, logger fanout, and collector integration caveats. |
| [Prometheus](infrastructure/prometheus.md) | Document metrics scraping and exposure path. | Covers collector metrics export, scrape configuration, and operational metrics caveats. |
| [Jaeger](infrastructure/jaeger.md) | Document trace sink and local trace UI integration. | Covers trace pipeline export, Jaeger service wiring, and local operational caveats. |

## Security

| Component | Purpose | Functionality |
|---|---|---|
| [Skip Auth Middleware](security/skip-auth-middleware.md) | Document the local identity injection used when auth is skipped. | Covers the `x-skip-auth-*` headers, default locals, and the current runtime auth gap. |

## Resilience

| Component | Purpose | Functionality |
|---|---|---|
| [Circuit Breaker](resilience/circuit-breaker.md) | Document request admission and breaker behavior. | Covers Redis-backed `Allow` checks, fail-open behavior, and current transition-hook limitations. |
