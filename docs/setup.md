# Setup and Development

## Requirements
- Go `1.25.x` (`go.mod` sets `go 1.25`).
- Docker and Docker Compose (for telemetry and optional Redis stack).
- Local Redis at `localhost:6379` for runtime health checks and several tests.

## Environment Variables
`core.NewConfigFromEnv` reads the following keys:

| Category | Variables | Defaults |
|---|---|---|
| Service | `ENVIRONMENT`, `PORT`, `SKIP_AUTH` | `development`, `8000`, `false` |
| OTel | `OTEL_DISABLE`, `OTEL_OTLP_EXPORTER_ENDPOINT`, `OTEL_OTLP_EXPORTER_INSECURE` | `false`, `localhost:4317`, `false` |
| Cognito | `COGNITO_REGION`, `COGNITO_USER_POOL_ID`, `COGNITO_APP_CLIENT_ID` | `us-east-1`, `UNSET`, `UNSET` |
| Redis | `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB` | `localhost:6379`, empty, `0` |
| NSC | `NSC_SUBMIT_URL`, `NSC_TOKEN_URL`, `NSC_CLIENT_SECRET`, `NSC_CLIENT_ID`, `NSC_ACCOUNT_ID` | empty |

### `.env.example` Caveat
`.env.example` currently uses `Port` and `Environment` (mixed case), while code expects `PORT` and `ENVIRONMENT`.

## Local Run
### 1) Configure env
Use `.env.local` and/or `.env` with uppercase keys expected by the app.

### 2) Run service directly
```bash
go run .
```

### 3) Run with live reload (Air)
```bash
air
```

Air config is in `.air.toml` and builds with:
```bash
go build -o ./tmp/main -ldflags "-X verification-service-api/pkg/core.ServiceVersion=local" .
```

## Docker Workflows
### App + Observability stack
```bash
docker compose up --build
```
Services:
- API (`:8000`)
- OTel Collector (`:4317`, `:4318`, metrics endpoints)
- Jaeger UI (`:16686`)
- Prometheus (`:9090`)

### App + Observability + Redis
```bash
docker compose -f docker-compose-with-redis.yml up --build
```
Adds Redis (`:6379`) for local circuit-breaker/status behavior.

## Build
```bash
go build ./...
```

Container build:
```bash
docker build .
```

## Test
```bash
go test ./...
```

### Test Prerequisites
- Redis must be running on `localhost:6379` for:
  - `api/routes/status_router_test.go`
  - `pkg/circuitbreaker/*_test.go`
  - `pkg/redis/redis_test.go`

### Known Test Behavior (Observed)
- Without Redis, Redis-dependent tests fail with connection refused/timeouts.
- `pkg/core/TestLoadEnv` currently expects a non-nil error even when `LoadEnv()` may return `nil`; behavior appears logically inconsistent with the assertion message.

## Telemetry Notes
- OTel service is enabled unless `OTEL_DISABLE=true`.
- OTel collector config: `otel-collector-config.yml`.
- Prometheus scrape config: `prometheus.yml`.
- Logger fanout can include OTEL log bridge via `core.NewLoggerWithOtel`.

## Assumptions
- **High confidence:** Local Redis is mandatory for meaningful integration testing in this repo's current state.
- **Medium confidence:** CI test behavior may differ if CI provisions Redis implicitly.
