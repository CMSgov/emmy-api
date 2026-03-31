# Architecture

## Package Structure

- `main.go`
- `api/app.go`
- `api/routes/*.go`
- `api/handlers/*.go`
- `api/middleware/middleware.go`
- `pkg/core/*.go`
- `pkg/education/*.go`
- `pkg/circuitbreaker/*.go`
- `pkg/redis/*.go`
- `pkg/choice/choice.go`

This page describes the current Go implementation. The intended public API
contract for the branch is defined separately in `api-spec/v0/openapi.yaml` and
`schema/v0/`.

## Dependency Relationships

```mermaid
flowchart LR
    main --> core
    main --> api
    main --> redispkg[pkg/redis]
    main --> routes

    api --> middleware
    api --> routes
    api --> core

    routes --> handlers
    routes --> middleware
    routes --> education
    routes --> circuitbreaker

    handlers --> education
    handlers --> redispkg

    middleware --> circuitbreaker

    education --> core
    circuitbreaker --> redisclient[go-redis]
    redispkg --> redisclient
```

## Interfaces and Abstractions

- `pkg/education/service.go`
  - `type EducationService interface { Submit(ctx context.Context, req Request) (Response, error) }`
  - `type HTTPTransport interface { Do(req *http.Request) (*http.Response, error) }`
- `pkg/core/otel.go`
  - `type OtelService interface { SpanFromContext; LoggerProvider; Shutdown }`
- `pkg/circuitbreaker/circuitbreaker.go`
  - `type Breaker interface { Allow; OnSuccess; OnFailure }`

These abstractions support unit testing and integration boundary replacement
without route-layer rewrites.

## Concurrency Model

- Server lifecycle:
  - `runServer` starts `app.Listen` in a goroutine and selects on server error or signal context cancellation.
  - graceful shutdown uses `app.ShutdownWithTimeout(5 * time.Second)`.
- Request lifecycle:
  - handlers create per-request contexts with timeout (`/health`: 2s, `POST /api/v0/education-enrollments`: 5s, `POST /api/v0/veteran-disability-ratings`: 5s).
- Circuit-breaker middleware:
  - breaker registry map guarded with `sync.RWMutex`.
  - lazy breaker initialization via double-check lock pattern.

## Error Handling Strategy

- Global Fiber error handler (`api/errorHandler`) converts errors into HTTP status and message.
- `fiber.NewError` used for explicit gateway semantics in education handler.
- Wrapped errors with context (`fmt.Errorf("...: %w", err)`) in service layers.
- Circuit breaker denies with `503 Service Unavailable` when `Allow` returns `ErrCircuitOpen`.
- With default `FailOpen=true`, Redis state-read/parse errors in breaker checks allow request pass-through instead of denying.
- Panic recovery middleware logs stack traces (`recover.Config{EnableStackTrace:true}`).

## Middleware Stack

Ordered middleware in `api.New`:

1. Recover
2. CORS (`*` origin/headers/methods)
3. OpenTelemetry Fiber middleware
4. Structured request logging (trace/span/request IDs)
5. Conditional `SkipAuthMiddleware` when `SKIP_AUTH=true`

## Dependency Injection Pattern

Observed constructor and options-based DI:

- `education.New(cfg, education.Options{HTTPClient, Logger, Timeout})`
- Current `main` call: `api.New(&api.Config{Core, Logger, Otel, Redis})`
- Circuit breaker injection via higher-order middleware factory:
  - `WithCircuitBreaker(func(name string) *RedisBreaker { ... })`

Note: `api.Config` includes a `Redis` field and the current `main` path now
injects it.

## Technical Caveats (Current State)

- `POST /api/v0/education-enrollments` builds a hardcoded request payload
  instead of binding user input from the caller body.
- `/health` is registered before the auth middleware, so it remains a runtime-only unauthenticated health route.
- Current runtime routes are `GET /`, `GET /health`, `GET /api-spec/v1/verify`,
  `POST /api/v0/education-enrollments`, and
  `POST /api/v0/veteran-disability-ratings`.
- `POST /api/v0/education-enrollments` uses a contract-aligned path but remains
  runtime scaffolding, while `POST /api/v0/veteran-disability-ratings` matches
  the checked-in v0 route and binds caller input.
- When `SKIP_AUTH=false`, no request-auth middleware is currently installed in
  `api.New`.
- Some tests require local Redis and fail when unavailable.

## Assumptions

- **High confidence:** Current layering is intentionally thin and
  integration-oriented rather than strict clean-architecture separation.
- **Medium confidence:** Additional provider adapters are expected to follow
  existing interface patterns in `pkg/education` and middleware wrappers.
