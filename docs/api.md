# Runtime API Notes (Current Go Implementation)

## Overview

The server port is configured through `PORT` and defaults to `3000` in code;
the local example environment sets `PORT=8000`. The intended public API
contract for this branch is defined in `api-spec/v0/openapi.yaml`; this page
documents the currently wired Go runtime endpoints and their operational
caveats.

## Authentication Behavior

- When `SKIP_AUTH=true`, `api.New` installs `SkipAuthMiddleware`, which injects
  a stable local identity into Fiber locals.
- Optional override headers for local development:
  - `x-skip-auth-sub`
  - `x-skip-auth-username`
  - `x-skip-auth-scope`
  - `x-skip-auth-groups`
- Default injected values are:
  - `sub=skip-auth-user`
  - `username=skip-auth-user`
  - `scope=local`
  - `groups=["local-dev"]`
- When `SKIP_AUTH=false`, the current branch does not install any request-auth
  middleware.

`/health` is registered before `SkipAuthMiddleware` and remains unauthenticated
in all cases on this branch.

## Circuit Breaker Behavior

`/health`, `POST /api/v0/education-enrollments`, and
`POST /api/v0/veteran-disability-ratings` are wrapped by Redis-backed circuit
breaker middleware.

- On breaker deny/open state: `503 Service Unavailable`.
- On Redis state read failures with fail-open (default): request is allowed.
- `GET /` and `GET /api-spec/v1/verify` are not wrapped by the breaker.

## Runtime Endpoints

| Method | Path | Description | Success | Notes |
|---|---|---|---|---|
| `GET` | `/` | Liveness string | `200` text | Returns `Backend running!`. |
| `GET` | `/health` | Redis health check | `200` empty | Pings Redis with a 2-second timeout and is wrapped by the circuit breaker. |
| `GET` | `/api-spec/v1/verify` | Bundled OpenAPI JSON artifact | `200` JSON | Returns the checked-in `api-spec/v0/dist/openapi.bundled.json` artifact. |
| `POST` | `/api/v0/education-enrollments` | NSC education verification route | `200` JSON | Path matches the v0 contract, but the handler still ignores the caller body and submits a hardcoded `education.Request`. |
| `POST` | `/api/v0/veteran-disability-ratings` | Veteran disability status from the v0 route | `200` JSON | Accepts caller-provided identity payload and validates required fields before upstream lookup. |

## Example: `/health`

```bash
curl -i http://localhost:8000/health
```

### `/api-spec/v1/verify`

```bash
curl -i http://localhost:8000/api-spec/v1/verify
```

Returns the checked-in bundled OpenAPI JSON artifact with `Content-Type: application/json`.

## Example: `POST /api/v0/education-enrollments`

```bash
curl -i --request POST http://localhost:8000/api/v0/education-enrollments \
  --header 'Content-Type: application/json' \
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

- `POST /api/v0/education-enrollments` currently ignores caller-provided JSON
  and submits a hardcoded upstream NSC request built in the handler.
- `main` injects Redis into `api.New`, so the health route and breaker-backed
  routes use the same Redis dependency that is checked at startup.
- The intended public contract for this branch is versioned under
  `api-spec/v0/`. The route paths match that contract, but the education
  runtime path is still scaffolded and current runtime auth does not yet match
  the contract's OAuth2 security model.
- Error response bodies come from Fiber error handling and may be plain text.
- When `SKIP_AUTH=false`, routes remain unauthenticated on the current branch.

## Assumptions

- **High confidence:** This page is a runtime reference, not the public API
  contract reference.
- **High confidence:** `POST /api/v0/education-enrollments` is intended to
  converge on the published v0 contract but does not yet bind the public
  request body.
