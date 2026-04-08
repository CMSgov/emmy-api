# Research: Deferred Runtime Auth Token Validation Strategy

## Problem Statement

The checked-in public contract expects bearer-token authentication, but the
current Go runtime on this branch does not yet implement in-process token
validation. This note preserves earlier design directions for future auth work
rather than describing current runtime behavior.

## Alternatives Considered

- Offline JWT validation using a JWKS-backed verifier.
- Token introspection against an upstream auth server.
- API gateway-only auth with no in-app verification.

## Trade-offs

- Offline JWKS validation:
  - Pros: no per-request introspection call, low latency, direct control over
    claims checks.
  - Cons: key cache management and claim-policy maintenance in service code.
- Introspection:
  - Pros: centralized revocation semantics.
  - Cons: network dependency per request.
- Gateway-only:
  - Pros: less app code.
  - Cons: reduced defense-in-depth and local context extraction flexibility.

## Why This Approach Was Previously Favored (Inferred)

Earlier iterations appear to have favored low-latency local validation with
explicit issuer and client claim checks.

## Benchmarks / Status

- Not available.
- No auth latency metrics are currently instrumented in repo.

## References

- `api/middleware/middleware.go`
- `api/app.go`
- Historical references in dated audit reports under `docs/audit/`

## Assumptions

- **High confidence:** This page is forward-looking research, not a description
  of the middleware currently wired in `api/app.go`.
