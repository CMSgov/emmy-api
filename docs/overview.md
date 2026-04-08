# Verification Service API Overview

## Purpose

The Verification Service API provides a unified HTTP interface for eligibility
verification workflows, currently focused on Redis-backed health checks plus two
versioned verification routes for education enrollment and veteran disability.

This service evolved from consent-based verification work and is intended to
reduce manual burden during benefits eligibility evaluation.

The intended public API contract for this branch is defined in
`api-spec/v0/openapi.yaml` and the reusable schemas in `schema/v0/`. This page
describes the current repository and runtime shape, including the places where
contract-first docs are ahead of the running service.

## System Context

Runtime dependencies in current implementation:

- Fiber (`github.com/gofiber/fiber/v2`) for HTTP server and routing.
- Redis (`github.com/redis/go-redis/v9`) for health checks and distributed
  circuit-breaker state.
- NSC endpoints (`NSC_TOKEN_URL`, `NSC_SUBMIT_URL`) for education verification.
- VA endpoints (`VA_TOKEN_URL`, `VA_BASE_URL`) for veteran verification.
- Optional local identity injection via `SkipAuthMiddleware` when
  `SKIP_AUTH=true`.
- OpenTelemetry OTLP exporter for tracing/metrics/log fanout.

## Key Packages

- `main`: process bootstrap, env/config load, OTel startup, Redis client init,
  route registration, graceful shutdown.
- `api`: Fiber app construction and shared middleware setup.
- `api/routes`: endpoint registration (`/`, `/health`,
  `/api-spec/v1/verify`, `/api/v0/education-enrollments`,
  `/api/v0/veteran-disability-ratings`).
- `api/handlers`: HTTP handlers for Redis health, education verification,
  veteran verification, and spec serving.
- `api/middleware`: local skip-auth identity injection and circuit-breaker
  middleware.
- `pkg/core`: configuration, logger, OTel service abstractions/utilities.
- `pkg/education`: NSC service abstraction and HTTP/OAuth submit flow.
- `pkg/veteran`: VA service abstraction and JWT client-assertion flow.
- `pkg/circuitbreaker`: Redis-backed circuit-breaker implementation.
- `pkg/redis`: Redis client factory and health ping.

## Design Principles (Observed)

- Explicit startup configuration from environment via `core.NewConfigFromEnv()`.
- Interface-driven boundaries for integration points (`EducationService`,
  `HTTPTransport`, `OtelService`, `Breaker`).
- Middleware-first cross-cutting concerns (request IDs, logging, recovery,
  CORS, tracing, circuit breaking, and optional local auth bypass).
- Operational defaults favoring availability in unknown breaker state
  (`FailOpen=true` by default).

## High-Level Request Flow

```mermaid
flowchart TD
    A[Client] --> B[Fiber App]
    B --> C[Request ID + logging + recovery + CORS + OTel]
    C --> D{Route}
    D -->|"/"| E[Return liveness string]
    D -->|"/api-spec/v1/verify"| F[Serve bundled OpenAPI JSON]
    D -->|"/health"| G{Circuit Breaker Allow?}
    D -->|"/api/v0/*"| H{Circuit Breaker Allow?}

    G -->|No| I[503 Service Unavailable]
    G -->|Yes| J[Redis Ping]
    J --> K[200 OK or Fiber Error]

    H -->|No| I
    H -->|Yes| L{SKIP_AUTH == true?}
    L -->|Yes| M[Inject local identity]
    L -->|No| N[Continue without extra auth middleware]

    M --> O{Verification route}
    N --> O

    O -->|Education| P[EducationService.LookupEnrollmentStatus]
    P --> Q[OAuth2 client credentials token]
    Q --> R[NSC submit endpoint]
    R --> S[JSON response]

    O -->|Veteran| T[VeteranService.LookupDisabilityRating]
    T --> U[VA token exchange]
    U --> V[VA disability endpoint]
    V --> W[JSON response]

    B -.-> X[OpenTelemetry exporter]
    J -.-> X
    P -.-> X
    T -.-> X
```

Current wiring caveat on `main`: `api.New` receives a Redis client from `main`,
so `/health` and the circuit-breaker middleware share the same dependency.

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
- **High confidence:** `POST /api/v0/education-enrollments` and
  `POST /api/v0/veteran-disability-ratings` are the two verification routes
  wired by the current branch.
- **Medium confidence:** Additional verification domains beyond the current
  runtime routes may be introduced in future versions.
