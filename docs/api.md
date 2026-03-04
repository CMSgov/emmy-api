# API Documentation (Current Implementation)

## Overview
Base server runs on port `8000` by default. Current routes are small and implementation-focused.

## Authentication Behavior
- When `SKIP_AUTH=false`, Cognito middleware is enabled globally.
- Middleware reads access token from header: `x-amzn-oidc-accesstoken`.
- Token checks include:
  - valid signature via JWKS
  - issuer match
  - `token_use=access`
  - `client_id` claim equals configured app client ID

If auth fails, response is `401 Unauthorized`.
This applies to both `/status` and `/api/edu` when auth is enabled.

## Circuit Breaker Behavior
`/status` and `/api/edu` are wrapped by Redis-backed circuit breaker middleware.
- On breaker deny/open state: `503 Service Unavailable`.
- On Redis state read failures with fail-open (default): request is allowed.

## Endpoints

| Method | Path | Description | Success | Notes |
|---|---|---|---|---|
| `GET` | `/` | Liveness string | `200` text | Returns `Backend running!` |
| `GET` | `/status` | Redis health check | `200` empty | Auth required unless `SKIP_AUTH=true`; pings Redis with 2s timeout |
| `GET` | `/api/edu` | NSC education verification passthrough | `200` JSON | Uses hardcoded request payload in handler |

## Request/Response Models
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

## Example: `/status`
```bash
curl -i http://localhost:8000/status
```
If `SKIP_AUTH=false`, include `x-amzn-oidc-accesstoken` header.

## Example: `/api/edu` (auth skipped locally)
```bash
curl -i http://localhost:8000/api/edu
```

## Current-State Caveats
- `/api/edu` currently does not accept caller-provided payload; it submits a hardcoded sample request from handler code.
- Current `main` wiring registers `/status` through `api.New` with a nil Redis client (because `main` does not inject `api.Config.Redis`), so runtime behavior can fail/panic until code wiring is corrected.
- This endpoint should be treated as implementation scaffold unless product/API contract is formalized.
- Error response bodies come from Fiber error handling and may be plain text.

## Assumptions
- **High confidence:** `/status` is the only operationally stable endpoint for external health checking today.
- **Medium confidence:** `/api/edu` contract and HTTP method may change when request binding/validation is introduced.
