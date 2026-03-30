# Verification Service API Overview

## Purpose

The Verification Service API provides a unified HTTP interface for eligibility
verification workflows, currently focused on education verification through NSC
(National Student Clearinghouse) integration and system health verification
through Redis-backed checks.

This service evolved from consent-based verification work and is intended to
reduce manual burden during benefits eligibility evaluation.

The intended public API contract for this branch is defined in
`api-spec/v0/openapi.yaml` and the reusable schemas in `schema/v0/`. This page
describes the current repository and runtime shape, which still contains
implementation scaffolding that has not yet converged on the full v0 contract.

## System Context

Runtime dependencies in current implementation:

- Fiber (`github.com/gofiber/fiber/v2`) for HTTP server and routing.
- Redis (`github.com/redis/go-redis/v9`) for health checks and distributed circuit-breaker state.
- NSC endpoints (`NSC_TOKEN_URL`, `NSC_SUBMIT_URL`) for education verification.
- AWS Cognito JWKS/JWT validation for request authentication (when `SKIP_AUTH=false`).
- OpenTelemetry OTLP exporter for tracing/metrics/log fanout.

## Key Packages

- `main`: process bootstrap, env/config load, OTel startup, Redis client init, route registration, graceful shutdown.
- `api`: Fiber app construction and shared middleware setup.
- `api/routes`: endpoint registration (`/`, `/status`, `/api/edu`).
- `api/handlers`: HTTP handlers for Redis status and education verification.
- `api/middleware`: Cognito auth and circuit-breaker middleware.
- `pkg/core`: configuration, logger, OTel service abstractions/utilities.
- `pkg/education`: NSC service abstraction and HTTP/OAuth submit flow.
- `pkg/circuitbreaker`: Redis-backed circuit-breaker implementation.
- `pkg/redis`: Redis client factory and health ping.

## Design Principles (Observed)

- Explicit startup configuration from environment via `core.NewConfigFromEnv()`.
- Interface-driven boundaries for integration points (`EducationService`, `HTTPTransport`, `OtelService`, `Breaker`).
- Middleware-first cross-cutting concerns (recovery, CORS, tracing, request logging, auth, circuit breaking).
- Operational defaults favoring availability in unknown breaker state (`FailOpen=true` by default).

## High-Level Request Flow

```mermaid
flowchart TD
    A[Client] --> B[Fiber App]
    B --> C[Recover + CORS + OTel + Slog middleware]
    C --> D{SKIP_AUTH == false?}
    D -->|Yes| E[Cognito JWT Verifier]
    D -->|No| F[Route Handler]
    E --> F

    F --> G{Circuit Breaker Allow?}
    G -->|No| H[503 Service Unavailable]

    G -->|Yes: /status| I[Redis Ping]
    I --> J[200 OK or Fiber Error]

    G -->|Yes: /api/edu| K[EducationService.Submit]
    K --> L[OAuth2 client credentials token]
    L --> M[NSC submit endpoint]
    M --> N[JSON response]

    B -.-> O[OpenTelemetry exporter]
    I -.-> O
    K -.-> O
```

Current wiring caveat on `main`: `/status` is registered in `api.New` using
`api.Config.Redis`, but `main` currently builds `api.Config` without a Redis
client. Treat the `/status` flow above as intended behavior pending that wiring
fix.

## Documentation Map

- [Architecture](architecture.md)
- [Setup](setup.md)
- [Runtime API Notes](api.md)
- [Features](features/)
- [Research](research/)
- [Planning](planning/)
- [API Specification](../api-spec/README.md)
- [JSON Schemas](../schema/README.md)

Feature docs are categorized by domain under `features/core`,
`features/infrastructure`, `features/security`, and `features/resilience`.

## Naming Note

Initial requirements referenced `/docs/planing`; this repo standardizes on
`/docs/planning`.

## Assumptions

- **High confidence:** Redis is the only persistent/shared runtime store
  currently used by this service.
- **High confidence:** `/api/edu` is presently implementation scaffolding and
  should not be treated as the public contract for this branch.
- **Medium confidence:** Additional verification domains beyond the two v0
  public operations may be introduced in future versions.
