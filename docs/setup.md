# Setup and Development

## Requirements

- Go `1.25.x` (`go.mod` sets `go 1.25`).
- Docker and Docker Compose for containerized local workflows. The committed
  compose file provides API, Redis, and local observability services.
- Local Redis at `localhost:6379` for runtime health checks and several tests.

## Environment Variables

`core.NewConfigFromEnv` reads the following keys:

| Category | Variables | Defaults |
|---|---|---|
| Service | `ENVIRONMENT`, `PORT`, `SKIP_AUTH` | `development`, `3000`, `false` |
| OTel | `OTEL_DISABLE`, `OTEL_OTLP_EXPORTER_ENDPOINT`, `OTEL_OTLP_EXPORTER_INSECURE` | `true`, `localhost:4317`, `false` |
| Redis | `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`, `REDIS_USE_TLS`, `REDIS_INSECURE_SKIP_VERIFY` | `localhost:6379`, empty, `0`, `true`, `false` |
| NSC | `NSC_SUBMIT_URL`, `NSC_TOKEN_URL`, `NSC_CLIENT_SECRET`, `NSC_CLIENT_ID`, `NSC_ACCOUNT_ID` | empty |
| VA | `VA_BASE_URL`, `VA_TOKEN_URL`, `VA_CLIENT_ID`, `VA_AUD`, `VA_PRIVATE_KEY_PATH`, `VA_TIMEOUT_SECONDS` | empty, empty, empty, empty, empty, `5` |

- The table above reflects code defaults from `pkg/core/config.go`.
- `.env.example` overrides the local example port to `PORT=8000`, sets
  `SKIP_AUTH=true`, and includes placeholders for VA veteran-verification
  credentials.
- VA authentication uses a signed JWT client assertion, so the configured
  private key path must point to a readable RSA PEM file on disk.
- Populate the VA values before exercising
  `POST /api/v0/veteran-disability-ratings`.

## Local Run

### 1. Configure env

Create `.env.local` and/or `.env` from `.env.example`. Adjust variables to your
preferred values.

For local Redis started via `docker compose` or `redis-server`, set
`REDIS_USE_TLS=false`. The code default is `true`, which is appropriate for
TLS-enabled deployments but will cause local startup to hang or fail against
the plain `redis:7` container in this repo's compose stack.

### 2. Run service directly

```bash
go run .
```

### 3. Run with live reload (Air)

Air is a development watcher that rebuilds and restarts the app when Go files
change, so you can iterate without re-running `go run .` manually.

Install Air (Go toolchain install):

```bash
go install github.com/air-verse/air@latest
```

If `air` is not found after install, add your Go bin directory to `PATH`
(commonly `$(go env GOPATH)/bin`).

Run:

```bash
air
```

`air` is optional. This repo includes `.air.toml` with build command:

```bash
go build -o ./tmp/main -ldflags "-X github.com/cmsgov/emmy-api/pkg/core.ServiceVersion=local" .
```

## Docker Workflows

### App + Observability stack

```bash
docker compose up --build
```

Services:

- API (`:8000` from the compose/example env)
- Redis (`:6379`)
- OTel Collector (`:4317`, `:4318`, metrics endpoints)
- Jaeger UI (`:16686`)
- Prometheus (`:9090`)

The API container is configured with `REDIS_ADDR=redis:6379`, so the compose
stack includes the Redis dependency needed for local startup and
circuit-breaker/health behavior. It also sets `REDIS_USE_TLS=false` because the
local Redis container does not serve TLS.

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
- `pkg/core/TestLoadEnv` currently expects a non-nil error even when
  `LoadEnv()` may return `nil`; behavior appears logically inconsistent with the
  assertion message.

## Telemetry Notes

- OTel service is enabled unless `OTEL_DISABLE=true`.
- OTel collector config: `otel-collector-config.yml`.
- Prometheus scrape config: `prometheus.yml`.
- Logger fanout can include OTel log bridge via `core.NewLoggerWithOtel`.

## Assumptions

- **High confidence:** Local Redis is mandatory for meaningful integration
  testing in this repo's current state.
- **Medium confidence:** CI test behavior may differ if CI provisions Redis
  implicitly.
