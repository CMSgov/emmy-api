# Runtime API Notes (Current Go Implementation)

## Overview

The server port is configured through `PORT` and defaults to `3000` in code;
the local example environment sets `PORT=8000`. The intended public API
contract for this branch is defined in `api-spec/v0/openapi.yaml`; this page
documents the endpoints and operational behavior that are actually wired in the
current Go runtime.

## Authentication Behavior

- The current branch does not register a Cognito or bearer-token validation
  middleware.
- When `SKIP_AUTH=true`, `SkipAuthMiddleware` injects a stable local identity
  into Fiber locals using the local values `sub`, `username`, `scope`, and
  `groups`.
- When `SKIP_AUTH=false`, requests still reach the route handlers without an
  additional auth gate in the current code path.

The checked-in v0 contract still documents OAuth 2.0 client credentials. Treat
that as the intended public contract, not proof of current runtime enforcement.

## Circuit Breaker Behavior

`/health`, `POST /api/v0/education-enrollments`, and
`POST /api/v0/veteran-disability-ratings` are wrapped by Redis-backed circuit
breaker middleware.

- On breaker deny/open state: `503 Service Unavailable`.
- On Redis state read failures with fail-open (default): request is allowed.

## Runtime Endpoints

| Method | Path | Description | Success | Notes |
|---|---|---|---|---|
| `GET` | `/` | Liveness string | `200` text | Returns `Backend running!` |
| `GET` | `/health` | Redis health check | `200` empty | Pings Redis with a 2-second timeout and is wrapped by the breaker middleware |
| `GET` | `/api-spec/v1/verify` | Bundled OpenAPI JSON artifact | `200` JSON | Returns `api-spec/v0/dist/openapi.bundled.json` |
| `POST` | `/api/v0/education-enrollments` | Education enrollment lookup | `200` JSON | Binds request JSON, validates required identity fields, then calls NSC service |
| `POST` | `/api/v0/veteran-disability-ratings` | Veteran disability lookup | `200` JSON | Binds request JSON and requires either SSN or a complete address block |

## Request Validation and Error Semantics

- Education requests require `firstName`, `lastName`, and `dateOfBirth`.
- Veteran requests require `firstName`, `lastName`, and `dateOfBirth`, plus
  either `ssn` or a complete address (`street1`, `city`, `state`,
  `postalCode`, `country`).
- Fiber error handling sends plain-text response bodies for `400`, `502`, and
  `503` cases produced with `fiber.NewError(...)`.
- Both verification handlers return bare `404` responses for not-found cases.

## Examples

### `/health`

```bash
curl -i http://localhost:8000/health
```

### `/api-spec/v1/verify`

```bash
curl -i http://localhost:8000/api-spec/v1/verify
```

### `/api/v0/education-enrollments`

```bash
curl -i --request POST http://localhost:8000/api/v0/education-enrollments \
  --header 'Content-Type: application/json' \
  --data '{
    "firstName": "Lynette",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24",
    "ssn": "123-45-6789"
  }'
```

### `/api/v0/veteran-disability-ratings`

```bash
curl -i --request POST http://localhost:8000/api/v0/veteran-disability-ratings \
  --header 'Content-Type: application/json' \
  --data '{
    "firstName": "Lynette",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24",
    "address": {
      "street1": "17020 Tortoise St",
      "city": "Round Rock",
      "state": "TX",
      "postalCode": "78664",
      "country": "USA"
    }
  }'
```

## Current-State Caveats

- The route surface now matches the checked-in v0 contract for both public
  verification operations, but the runtime still serves plain-text error bodies
  instead of a versioned public error envelope.
- `GET /api-spec/v1/verify` is a runtime convenience route; the design-time
  source of truth remains `api-spec/v0/openapi.yaml`.
- The current branch does not enforce the contract's documented OAuth 2.0
  security scheme in Fiber middleware.

## Assumptions

- **High confidence:** This page is a runtime reference, not the public API
  contract reference.
- **High confidence:** `POST /api/v0/education-enrollments` and
  `POST /api/v0/veteran-disability-ratings` are the current verification
  endpoints exposed by the Go service.
