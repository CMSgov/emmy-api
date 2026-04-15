# Research: Auth Token Validation Strategy

## Problem Statement

This note captures an earlier Cognito/JWKS direction that is not currently
implemented on this branch. Keep it as historical research rather than current
runtime documentation.

## Alternatives Considered

- Offline JWT validation using Cognito JWKS (previously explored).
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

Earlier middleware design appears to have preferred low-latency local
validation with explicit issuer/client claim checks, but that implementation is
not present in the current branch.

## Benchmarks / Status

- Not available.
- No auth latency or JWKS refresh metrics are currently instrumented in repo.

## References

- `api/middleware/middleware.go`
- `api/app.go`
- Review git history if you need the prior Cognito-specific implementation details.

## Assumptions

- **Medium confidence:** A production auth layer still exists outside the
  currently committed skip-auth middleware, but this repository snapshot does
  not document it in runnable code.
