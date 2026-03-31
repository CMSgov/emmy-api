# Runtime API Notes (Current Go Implementation)

## Overview

The server port is configured through `PORT` and defaults to `3000` in code;
the local example environment sets `PORT=8000`. The intended public API
contract for this branch is defined in `api-spec/v0/openapi.yaml`; this page
documents the currently wired Go runtime endpoints and their operational
caveats.

## Authentication Behavior

- When `SKIP_AUTH=false`, Cognito middleware is enabled globally.
- Middleware reads access token from header: `x-amzn-oidc-accesstoken`.
- Token checks include:
  - valid signature via JWKS
  - issuer match
  - `token_use=access`
  - `client_id` claim equals configured app client ID

If auth fails, response is `401 Unauthorized`.
This applies to `/api/edu` and `/v0/veteran-disability-ratings` when auth is
enabled. `/health` is registered before the auth middleware and remains
unauthenticated in the current branch.

## Circuit Breaker Behavior

`/health`, `/api/edu`, and `/v0/veteran-disability-ratings` are wrapped by
Redis-backed circuit breaker middleware.

- On breaker deny/open state: `503 Service Unavailable`.
- On Redis state read failures with fail-open (default): request is allowed.

## Runtime Endpoints

| Method | Path                           | Description                            | Success     | Notes                                                              |
| ------ | ------------------------------ | -------------------------------------- | ----------- | ------------------------------------------------------------------ |
| `GET`  | `/`                            | Liveness string                        | `200` text  | Returns `Backend running!`                                         |
| `GET`  | `/health`                      | Redis health check                     | `200` empty | Registered before auth middleware; pings Redis with 2s timeout    |
| `GET`  | `/api/edu`                     | NSC education verification scaffold    | `200` JSON  | Uses a hardcoded request payload in handler; not the v0 contract   |
| `POST` | `/v0/veteran-disability-ratings` | Veteran disability status from v0 spec | `200` JSON  | Accepts caller-provided identity payload and matches the v0 route  |

| Method | Path                  | Description                         | Success     | Notes |
| ------ | --------------------- | ----------------------------------- | ----------- | ----- |
| `GET`  | `/`                   | Liveness string                     | `200` text  | Returns `Backend running!` |
| `GET`  | `/status`             | Redis health check                  | `200` empty | Uses 2s Redis ping timeout; wrapped by circuit breaker |
| `GET`  | `/api-spec/v1/verify` | Bundled OpenAPI JSON artifact       | `200` JSON  | Returns `api-spec/dist/openapi.bundled.json` |
| `GET`  | `/api/edu`            | Education verification passthrough  | `200` JSON  | Uses hardcoded request payload in handler; wrapped by circuit breaker |

### NSC Submit Request model (`pkg/education/models_request.go`)

```go
type Request struct {
    AccountID        string `json:"accountId"`
    OrganizationName string `json:"organizationName,omitempty"`
    CaseReferenceID  string `json:"caseReferenceId,omitempty"`
    ContactEmail     string `json:"contactEmail,omitempty"`
    DateOfBirth      string `json:"dateOfBirth"`
    LastName         string `json:"lastName"`
    FirstName        string `json:"firstName"`
    SSN              string `json:"ssn,omitempty"`
    IdentityDetails  []IdentityDetails `json:"identityDetails,omitempty"`
    EndClient        string `json:"endClient"`
    PreviousNames    []PreviousName `json:"previousNames,omitempty"`
    Terms            string `json:"terms"`
}
```

### NSC Submit Response model (`pkg/education/models_response.go`)

```go
type Response struct {
    ClientData          ClientDataResponse          `json:"clientData"`
    IdentityDetails     []IdentityDetailsResponse   `json:"identityDetails"`
    Status              StatusResponse              `json:"status"`
    StudentInfoProvided StudentInfoProvidedResponse `json:"studentInfoProvided"`
    TransactionDetails  TransactionDetailsResponse  `json:"transactionDetails"`
}
```

## Example: `/health`

```bash
curl -i http://localhost:8000/health
```

### `/api-spec/v1/verify`

```bash
curl -i http://localhost:8000/api-spec/v1/verify
```

Returns the checked-in bundled OpenAPI JSON artifact with `Content-Type: application/json`.

## Example: `/api/edu` (auth skipped locally)

```bash
curl -i http://localhost:8000/api/edu
```

## Example: `/v0/veteran-disability-ratings`

```bash
curl -i --request POST http://localhost:8000/v0/veteran-disability-ratings \
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

- `/api/edu` currently does not accept caller-provided payload; it submits a hardcoded sample request from handler code.
- `main` now injects Redis into `api.New`, so the health route has the Redis client it expects.
- The intended public contract for this branch is versioned under `api-spec/v0/`, and the veteran disability route matches that contract while `/api/edu` remains a runtime-only scaffold.
- Error response bodies come from Fiber error handling and may be plain text.

## Assumptions

- **High confidence:** This page is a runtime reference, not the public API
  contract reference.
- **Medium confidence:** `/api/edu` will be removed or reshaped as the runtime
  converges on the published v0 contract.
