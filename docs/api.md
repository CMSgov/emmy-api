# Runtime API Notes (Current Go Implementation)

## Overview

The server port is configured through `PORT` and defaults to `3000` in code.
The local example environment sets `PORT=8000`. The intended public contract
for this branch is defined in `api-spec/v0/openapi.yaml`; this page documents
the currently wired Go runtime behavior.

## Request Identity Behavior

`api.New` always installs `SubjectMiddleware`, which populates `c.Locals("sub")`
using this precedence:

1. `X-Sub`
2. `Authorization: Bearer <JWT>` parsed without signature verification, using
   only the token `sub` claim
3. fallback value `unknown-subject`

When `SKIP_AUTH=true`, `SkipAuthMiddleware` is also installed and injects a
deterministic local identity before `SubjectMiddleware` runs. Optional override
headers are:

- `x-skip-auth-sub`
- `x-skip-auth-username`
- `x-skip-auth-scope`
- `x-skip-auth-groups`

The current branch does not wire an active Cognito verifier into the Fiber app.
Treat the checked-in v0 auth scheme as contract documentation rather than a
runtime-enforced guarantee.

## Circuit Breaker Behavior

`GET /health`, `POST /api/v0/education-enrollments`, and
`POST /api/v0/veteran-disability-ratings` are wrapped by the Redis-backed
circuit-breaker middleware.

- Open breaker or breaker error response: `503 Service Unavailable`
- Redis read failures follow the breaker package's `FailOpen=true` default and
  allow the request through

## Runtime Endpoints

| Method | Path | Description | Success | Notes |
|---|---|---|---|---|
| `GET` | `/` | Liveness string | `200` text | Returns `Backend running!` |
| `GET` | `/health` | Redis health check | `200` empty | Uses a 2-second Redis ping timeout and is circuit-breaker wrapped |
| `GET` | `/api-spec/v1/verify` | Bundled OpenAPI JSON artifact | `200` JSON | Serves `api-spec/v0/dist/openapi.bundled.json` |
| `POST` | `/api/v0/education-enrollments` | NSC education verification | `200` JSON | Parses caller JSON and requires `firstName`, `lastName`, and `dateOfBirth` |
| `POST` | `/api/v0/veteran-disability-ratings` | Veteran disability lookup | `200` JSON | Requires `firstName`, `lastName`, `dateOfBirth`, and either `ssn` or a full address |

## Example: `/health`

```bash
curl -i http://localhost:8000/health
```

## Example: `/api-spec/v1/verify`

```bash
curl -i http://localhost:8000/api-spec/v1/verify
```

Returns the checked-in bundled OpenAPI JSON artifact with
`Content-Type: application/json`.

## Example: `/api/v0/education-enrollments`

```bash
curl -i --request POST http://localhost:8000/api/v0/education-enrollments \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer <ACCESS_TOKEN>' \
  --data '{
    "firstName": "Lynette",
    "lastName": "Oyola",
    "dateOfBirth": "1988-10-24"
  }'
```

## Example: `/api/v0/veteran-disability-ratings`

```bash
curl -i --request POST http://localhost:8000/api/v0/veteran-disability-ratings \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer <ACCESS_TOKEN>' \
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

- The v0 OpenAPI contract documents OAuth 2.0 client credentials, but the
  runtime currently uses `SubjectMiddleware` and optional skip-auth locals
  instead of a verifying auth middleware.
- Error responses come from Fiber error handling and are plain text unless a
  handler returns JSON explicitly.
- The education and veteran handlers add metadata such as `apiVersion`,
  `environment`, timestamps, and datasource duration to successful responses.

## Assumptions

- **High confidence:** This page is a runtime reference, not the public API
  contract reference.
- **High confidence:** `POST /api/v0/education-enrollments` and
  `POST /api/v0/veteran-disability-ratings` are the only current runtime POST
  endpoints under `/api/v0/`.
