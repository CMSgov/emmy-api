# Research: Auth Token Validation Strategy

## Problem Statement

If the runtime restores token enforcement, the API will need request
authentication that can validate Cognito-issued access tokens efficiently and
securely. The current branch does not implement that verifier.

## Alternatives Considered

- Offline JWT validation using Cognito JWKS.
- Token introspection against upstream auth server.
- API gateway-only auth with no in-app verification.
- Keep the current subject-extraction-only middleware and rely on upstream
  controls.

## Trade-offs

- Offline JWKS validation:
  - Pros: no per-request introspection call, low latency, direct control over claims checks.
  - Cons: key cache management and claim-policy maintenance in service code.
- Introspection:
  - Pros: centralized revocation semantics.
  - Cons: network dependency per request.
- Gateway-only:
  - Pros: less app code.
  - Cons: reduced defense-in-depth and local context extraction flexibility.
- Subject extraction only:
  - Pros: simple local behavior and low overhead.
  - Cons: no token verification or claim enforcement in service code.

## Current Branch Status

The checked-in middleware only extracts a `sub` value from headers or an
unverified bearer token. The stronger Cognito/JWKS validation path described in
older docs is not present in the current implementation.

## Benchmarks / Status

- Not available.
- No auth latency or JWKS refresh metrics are currently instrumented in repo.

## References

- `api/middleware/middleware.go`
- `api/app.go`
- Dependency used by current subject extraction: `github.com/lestrrat-go/jwx/v2`

## Assumptions

- **Medium confidence:** A future verified-auth middleware will likely build on
  the same request-local identity fields (`sub`, `username`, `scope`,
  `groups`) used today.
