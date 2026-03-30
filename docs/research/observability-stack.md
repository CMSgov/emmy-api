# Research: Observability Stack

## Problem Statement
The service requires trace and metric visibility for integration operations and runtime diagnostics.

## Alternatives Considered
- OpenTelemetry OTLP collector + Jaeger/Prometheus (current).
- Direct vendor-specific SDK integration.
- Minimal logging-only observability.

## Trade-offs
- OTel collector pipeline:
  - Pros: vendor-neutral, flexible exporters, consistent instrumentation surface.
  - Cons: operational overhead for collector deployment/config.
- Vendor SDKs:
  - Pros: deep vendor features.
  - Cons: lock-in and harder migration paths.
- Logging-only:
  - Pros: simplest.
  - Cons: weak distributed tracing and metric correlation.

## Why Current Approach Was Selected (Observed/Inferred)
Code initializes OTLP trace/meter exporters and compose files include collector, Jaeger, and Prometheus, indicating an explicit standards-based observability direction.

## Benchmarks / Status
- Not available.
- No repository benchmark data for exporter overhead.

## References
- `pkg/core/otel.go`
- `pkg/core/logger.go`
- `pkg/redis/redis.go`
- `api/app.go`
- `otel-collector-config.yml`
- `prometheus.yml`
- `docker-compose.yml`

## Assumptions
- **High confidence:** Observability is a first-class operational requirement, not an optional local-only feature.
