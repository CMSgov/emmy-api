# API Documentation (Current Implementation)

## Overview

The server listens on port `8000` by default and currently exposes a small set
of implementation-oriented routes plus a bundled OpenAPI JSON artifact route.

## Authentication Behavior

- Authentication is expected to be enforced upstream by the load balancer.
- When `SKIP_AUTH=true`, local identity values are injected through skip-auth middleware for development and testing.

## Circuit Breaker Behavior

`/status` and `/api/edu` are wrapped by Redis-backed circuit breaker middleware.

- On breaker deny/open state: `503 Service Unavailable`
- On Redis state read failures with fail-open behavior: request is allowed to continue

`/api-spec/v1/verify` is not wrapped by the circuit breaker.

## Endpoints

| Method | Path                  | Description                         | Success     | Notes |
| ------ | --------------------- | ----------------------------------- | ----------- | ----- |
| `GET`  | `/`                   | Liveness string                     | `200` text  | Returns `Backend running!` |
| `GET`  | `/status`             | Redis health check                  | `200` empty | Uses 2s Redis ping timeout; wrapped by circuit breaker |
| `GET`  | `/api-spec/v1/verify` | Bundled OpenAPI JSON artifact       | `200` JSON  | Returns `api-spec/dist/openapi.bundled.json` |
| `GET`  | `/api/edu`            | Education verification passthrough  | `200` JSON  | Uses hardcoded request payload in handler; wrapped by circuit breaker |

## Request and Response Models

### Education verify request sent to NSC (`pkg/education/models_request.go`)

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

### Education verify response returned from NSC (`pkg/education/models_response.go`)

```go
type Response struct {
    ClientData          ClientDataResponse          `json:"clientData"`
    IdentityDetails     []IdentityDetailsResponse   `json:"identityDetails"`
    Status              StatusResponse              `json:"status"`
    StudentInfoProvided StudentInfoProvidedResponse `json:"studentInfoProvided"`
    TransactionDetails  TransactionDetailsResponse  `json:"transactionDetails"`
}
```

## Examples

### `/status`

```bash
curl -i http://localhost:8000/status
```

### `/api-spec/v1/verify`

```bash
curl -i http://localhost:8000/api-spec/v1/verify
```

Returns the checked-in bundled OpenAPI JSON artifact with `Content-Type: application/json`.

### `/api/edu` (auth skipped locally)

```bash
curl -i http://localhost:8000/api/edu
```

## Current-State Caveats

- `/api/edu` currently does not accept caller-provided payload; it submits a hardcoded sample request from handler code.
- `/api-spec/v1/verify` serves the checked-in artifact from disk and does not rebuild the spec at request time.
- Error response bodies come from Fiber error handling and may be plain text.
- `main.go` currently constructs the app without passing a Redis client into `api.New`, so the `/status` route registered there may not behave correctly until that wiring is aligned.

## Assumptions

- **High confidence:** `/api-spec/v1/verify` is the most stable machine-readable contract endpoint currently exposed by the app.
- **Medium confidence:** `/api/edu` contract and HTTP method may change when request binding and public request validation are introduced.
