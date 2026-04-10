# Runtime API Notes (Current Go Implementation)

## Overview

The server port is configured through `PORT` and defaults to `3000` in code;
the local example environment sets `PORT=8000`. The intended public API
contract for this branch is defined in `api-spec/v0/openapi.yaml`; this page
documents the currently wired Go runtime endpoints and their operational
caveats.

## Authentication Behavior

- The checked-in public contract in `api-spec/v0/openapi.yaml` requires bearer
  tokens via `OAuth2ClientCredentials`.
- The current runtime wiring in `api.New` does not register a Cognito or bearer
  token enforcement middleware.
- When `SKIP_AUTH=true`, `SkipAuthMiddleware` injects local identity values from
  optional `x-skip-auth-*` headers into Fiber locals for routes registered
  after that middleware.
- `/health` is registered before `SkipAuthMiddleware`, so it remains an
  unauthenticated health route in the current branch.

## Circuit Breaker Behavior

`/health`, `POST /api/v0/education-enrollments`, and
`POST /api/v0/veteran-disability-ratings` are wrapped by Redis-backed circuit
breaker middleware.

- On breaker deny/open state: `503 Service Unavailable`.
- On Redis state read failures with fail-open (default): request is allowed.

## Runtime Endpoints

| Method | Path | Description | Success | Notes |
| ------ | ---- | ----------- | ------- | ----- |
| `GET` | `/` | Liveness string | `200` text | Returns `Backend running!` |
| `GET` | `/health` | Redis health check | `200` empty | Uses a 2-second Redis ping timeout and is wrapped by the circuit breaker |
| `GET` | `/api-spec/v1/verify` | Bundled OpenAPI JSON artifact | `200` JSON | Returns `api-spec/v0/dist/openapi.bundled.json` |
| `POST` | `/api/v0/education-enrollments` | Education enrollment lookup | `200` JSON | Parses caller-provided JSON into `education.Request` |
| `POST` | `/api/v0/veteran-disability-ratings` | Veteran disability status from v0 spec | `200` JSON | Parses caller-provided JSON into `veteran.Request` |

### Education Submit Request model (`pkg/education/models_request.go`)

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

### Education Submit Response model (`pkg/education/models_response.go`)

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
    "dateOfBirth": "1988-10-24",
    "ssn": "123-45-6789"
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

- The current runtime does not enforce the bearer-token auth described in the
  public OpenAPI contract.
- `main` injects Redis into `api.New`, so the health route and breaker
  middleware share the same Redis dependency.
- The intended public contract for this branch is versioned under
  `api-spec/v0/`, and both versioned POST routes are registered in the current
  runtime.
- Error response bodies come from Fiber error handling and may be plain text.

## Assumptions

- **High confidence:** This page is a runtime reference, not the public API
  contract reference.
- **High confidence:** `POST /api/v0/education-enrollments` is the current
  education verification route in the runtime.
