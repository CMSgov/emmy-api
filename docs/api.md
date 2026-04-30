# Runtime API Notes (Current Go Implementation)

## Overview

The server port is configured through `PORT` and defaults to `3000` in code;
the local example environment sets `PORT=8000`. The intended public API
contract for this branch is defined in `api-spec/v0/openapi.yaml`; this page
documents the currently wired Go runtime endpoints and their operational
caveats.

## Authentication Behavior

- The current branch does not install a token-verification middleware.
- `SubjectMiddleware` runs for all requests and sets `c.Locals("sub")` from:
  - `X-Sub`, if present
  - otherwise the `sub` claim in an `Authorization: Bearer ...` token, parsed
    without signature verification
  - otherwise the fallback value `unknown-subject`
- When `SKIP_AUTH=true`, `SkipAuthMiddleware` overwrites identity fields with a
  stable local subject and supports opt-in override headers such as
  `x-skip-auth-sub`.
- An invalid bearer token parse returns `401 Unauthorized`, but requests
  without a token are not rejected by auth middleware on this branch.

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
| `GET` | `/health` | Redis health check | `200` empty | Uses 2-second Redis ping timeout; wrapped by circuit breaker |
| `GET` | `/api-spec/v1/verify` | Bundled OpenAPI JSON artifact | `200` JSON | Returns the checked-in OpenAPI bundle |
| `POST` | `/api/v0/education-enrollments` | Education enrollment lookup | `200` JSON | Requires `firstName`, `lastName`, `dateOfBirth`; wraps NSC service |
| `POST` | `/api/v0/veteran-disability-ratings` | Veteran disability lookup | `200` JSON | Requires `firstName`, `lastName`, `dateOfBirth`, plus `ssn` or a complete address |

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

## Example: `/api-spec/v1/verify`

```bash
curl -i http://localhost:8000/api-spec/v1/verify
```

Returns the checked-in bundled OpenAPI JSON artifact with `Content-Type: application/json`.

## Example: `/api/v0/education-enrollments`

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

- The intended public contract requires OAuth 2.0 bearer auth, but this branch
  currently derives request identity without verifying bearer tokens.
- Error response bodies come from Fiber error handling and may be plain text.
- `POST /api/v0/education-enrollments` and
  `POST /api/v0/veteran-disability-ratings` both enrich successful responses
  with runtime metadata such as `apiVersion`, `environment`, and a request ID.
- Veteran lookups return `404` immediately when the request omits both `ssn`
  and a complete address.

## Assumptions

- **High confidence:** This page is a runtime reference, not the public API
  contract reference.
- **High confidence:** The education and veteran runtime paths now match the
  checked-in v0 contract paths.
- **Medium confidence:** Verified auth middleware will need to be reintroduced
  or documented separately when the runtime converges on the contract's OAuth
  requirements.
