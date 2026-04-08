# Runtime API Notes (Current Go Implementation)

## Overview

The server port is configured through `PORT` and defaults to `3000` in code.
The local example environment sets `PORT=8000`. The intended public API
contract for this branch is defined in `api-spec/v0/openapi.yaml`; this page
documents the endpoints and middleware that are actually wired by the current
Go service.

## Authentication Behavior

- `api.New` only adds `SkipAuthMiddleware` when `SKIP_AUTH=true`.
- That middleware injects a stable local identity into Fiber locals:
  `sub`, `username`, `scope`, and `groups`.
- When `SKIP_AUTH=false`, the current branch does not add an alternate
  in-app bearer-token validator.
- `GET /health` is registered before the optional skip-auth middleware and
  remains unauthenticated in the current branch.

The checked-in OpenAPI contract still models OAuth 2.0 client credentials for
public integrations. Treat that contract and the current runtime behavior as
separate concerns until runtime auth enforcement lands.

## Circuit Breaker Behavior

`/health`, `POST /api/v0/education-enrollments`, and
`POST /api/v0/veteran-disability-ratings` are wrapped by Redis-backed circuit
breaker middleware.

- On breaker deny/open state: `503 Service Unavailable`.
- On Redis state read failures with fail-open (default): request is allowed.

## Runtime Endpoints

| Method | Path | Description | Success | Notes |
|---|---|---|---|---|
| `GET` | `/` | Liveness string | `200` text | Returns `Backend running!`. |
| `GET` | `/health` | Redis health check | `200` empty | Pings Redis with a 2-second timeout. |
| `GET` | `/api-spec/v1/verify` | Bundled OpenAPI JSON artifact | `200` JSON | Serves `api-spec/v0/dist/openapi.bundled.json`. |
| `POST` | `/api/v0/education-enrollments` | Education enrollment lookup | `200` JSON | Requires `firstName`, `lastName`, and `dateOfBirth`; returns `enrollmentStatus`. |
| `POST` | `/api/v0/veteran-disability-ratings` | Veteran disability lookup | `200` JSON | Requires `firstName`, `lastName`, `dateOfBirth`, plus either `ssn` or a full `address`; returns `combinedDisabilityRating`. |

## Error Semantics

- Invalid JSON or missing required identity fields return `400`.
- Education and veteran lookups return `404` when the upstream service reports
  the subject was not found.
- Veteran disability requests also return `404` when neither `ssn` nor a full
  address is supplied.
- Upstream lookup failures return `502`.
- Circuit-breaker denies return `503`.
- Fiber-generated error bodies are plain text in the current branch.

## Example: `/health`

```bash
curl -i http://localhost:8000/health
```

## Example: `/api-spec/v1/verify`

```bash
curl -i http://localhost:8000/api-spec/v1/verify
```

## Example: `POST /api/v0/education-enrollments`

```bash
curl -i --request POST http://localhost:8000/api/v0/education-enrollments \
  --header 'Content-Type: application/json' \
  --data '{
    "firstName": "Lynette",
    "middleName": "Marie",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24",
    "ssn": "123-45-6789"
  }'
```

Example success body:

```json
{
  "enrollmentStatus": "FULL_TIME"
}
```

## Example: `POST /api/v0/veteran-disability-ratings`

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

Example success body:

```json
{
  "combinedDisabilityRating": 70
}
```

## Current-State Caveats

- The runtime and the checked-in OpenAPI contract both expose the same two POST
  verification routes, but runtime auth enforcement does not yet match the
  contract's bearer-token description.
- `main` injects Redis into `api.New`, so the health route and breaker
  middleware share the same Redis dependency.
- The bundled spec route is an implementation convenience and should not be
  treated as the authoring source of truth.

## Assumptions

- **High confidence:** This page is a runtime reference, not the public API
  contract reference.
- **Medium confidence:** Real bearer-token enforcement will be added later to
  align the running service with `api-spec/v0/openapi.yaml`.
