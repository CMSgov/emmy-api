# Runtime API Notes (Current Go Implementation)

## Overview

The server port is configured through `PORT` and defaults to `3000` in code;
the local example environment sets `PORT=8000`. The intended public API
contract for this branch is defined in `api-spec/v0/openapi.yaml`; this page
documents the currently wired Go runtime endpoints and their operational
caveats.

## Authentication Behavior

- The checked-in v0 contract declares OAuth 2.0 bearer auth in
  `api-spec/v0/openapi.yaml`.
- The current Go runtime on this branch does not register Cognito or bearer-token
  verification middleware.
- When `SKIP_AUTH=true`, the app attaches `SkipAuthMiddleware`, which injects a
  stable local identity into Fiber locals using optional `x-skip-auth-*`
  override headers.
- `/health` is registered before that development middleware and remains
  unauthenticated.

## Circuit Breaker Behavior

`/health`, `POST /api/v0/education-enrollments`, and
`POST /api/v0/veteran-disability-ratings` are wrapped by Redis-backed
circuit-breaker middleware.

- On breaker deny/open state: `503 Service Unavailable`.
- On Redis state read failures with fail-open (default): request is allowed.

## Runtime Endpoints

| Method | Path                           | Description                            | Success     | Notes                                                              |
| ------ | ------------------------------ | -------------------------------------- | ----------- | ------------------------------------------------------------------ |
| `GET`  | `/`                            | Liveness string                        | `200` text  | Returns `Backend running!`                                         |
| `GET`  | `/health`                      | Redis health check                     | `200` empty | Registered before skip-auth middleware; pings Redis with 2s timeout |
| `GET`  | `/api-spec/v1/verify`          | Bundled OpenAPI JSON artifact         | `200` JSON  | Returns the checked-in spec artifact served by `OpenAPISpecHandler` |
| `POST` | `/api/v0/education-enrollments` | Education enrollment verification      | `200` JSON  | Parses caller JSON and returns `enrollmentStatus`                   |
| `POST` | `/api/v0/veteran-disability-ratings` | Veteran disability status from v0 spec | `200` JSON  | Accepts caller-provided identity payload and returns rating data    |

### Education request model (`pkg/education/models_request.go`)

```go
type Request struct {
    FirstName   string   `json:"firstName"`
    MiddleName  string   `json:"middleName,omitempty"`
    LastName    string   `json:"lastName"`
    DateOfBirth string   `json:"dateOfBirth"`
    SSN         string   `json:"ssn,omitempty"`
    Address     *Address `json:"address,omitempty"`
}
```

### Education response model (`pkg/education/models_response.go`)

```go
type Response struct {
    EnrollmentStatus EnrollmentStatus `json:"enrollmentStatus"`
}
```

## Example: `/health`

```bash
curl -i http://localhost:8000/health
```

## Example: `/api-spec/v1/verify`

```bash
curl -i http://localhost:8000/api-spec/v1/verify
```

Returns the checked-in bundled OpenAPI artifact with `Content-Type:
application/json`.

## Example: `/api/v0/education-enrollments`

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

## Example: `/api/v0/veteran-disability-ratings`

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

- `main` now injects Redis into `api.New`, so the health route has the Redis client it expects.
- The intended public contract for this branch is versioned under `api-spec/v0/`, and both versioned POST routes are registered in the current runtime.
- The current runtime docs do not imply active bearer-token enforcement because no Cognito or bearer verification middleware is wired on this branch.
- Error response bodies come from Fiber error handling and may be plain text.

## Assumptions

- **High confidence:** This page is a runtime reference, not the public API
  contract reference.
- **Medium confidence:** The checked-in contract and runtime auth behavior are
  still converging, so future auth-related docs may need another refresh.
