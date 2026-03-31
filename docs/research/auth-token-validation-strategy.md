# Research: Auth Token Validation Strategy

## Problem Statement

The checked-in v0 contract expects authenticated requests, but the current Go
runtime does not install request-auth middleware when `SKIP_AUTH=false`.
This note captures possible future directions for closing that gap.

## Alternatives Considered

- Offline JWT validation using a JWKS-backed verifier.
- Token introspection against upstream auth server.
- API gateway-only auth with no in-app verification.
- Keep the current runtime behavior and rely on trusted upstream infrastructure.

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

## Why A JWKS-Backed Approach Still Looks Plausible (Inferred)

The repository still vendors JWT/JWX dependencies for other upstream auth
flows, so a local JWKS-backed verifier would fit the existing stack if service
level auth is restored.

## Benchmarks / Status

- Not available.
- No runtime auth latency or JWKS refresh metrics are currently instrumented in
  repo.

## References

- `api/middleware/middleware.go`
- `api/app.go`
- Dependencies already present in `go.mod`, including
  `github.com/lestrrat-go/jwx/v2`

## Assumptions

- **High confidence:** The current branch has an auth-behavior gap between the
  published v0 contract and the live runtime wiring.
