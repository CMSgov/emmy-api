# Research: Auth Token Validation Strategy

## Problem Statement
The API requires request authentication that can validate Cognito-issued access tokens efficiently and securely.

## Alternatives Considered
- Offline JWT validation using Cognito JWKS (current).
- Token introspection against upstream auth server.
- API gateway-only auth with no in-app verification.

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

## Why Current Approach Was Selected (Inferred)
The middleware design and `jwk.Cache` usage imply preference for low-latency local validation with explicit issuer/client claim checks.

## Benchmarks / Status
- Not available.
- No auth latency or JWKS refresh metrics are currently instrumented in repo.

## References
- `api/middleware/middleware.go`
- `api/app.go`
- Dependencies: `github.com/lestrrat-go/jwx/v2`

## Assumptions
- **Medium confidence:** Upstream architecture expects service-level auth enforcement even when requests may already pass through trusted infrastructure.
