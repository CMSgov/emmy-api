# Feature: Cognito Auth

## Feature Overview

Captures the intended Cognito access-token validation model for this repository.
The current runtime branch does not register a Cognito verifier in `api.New`;
it only provides `SkipAuthMiddleware` for local identity injection when
`SKIP_AUTH=true`.

## Business Logic

- No Cognito validation path is currently wired into the request pipeline.
- When `SKIP_AUTH=true`, the app can inject local `sub`, `username`, `scope`,
  and `groups` values from `x-skip-auth-*` headers for downstream handlers.
- Bearer-token and JWKS-backed validation remain deferred work relative to the
  checked-in public contract.

## Package Location

- `api/middleware/middleware.go`
- `api/app.go`

## Key Structs and Interfaces

- `SkipAuthMiddleware`

## Real Code Excerpt

```go
if cfg.Core.SkipAuth {
    app.Use(middleware.SkipAuthMiddleware())
}
```

## Edge Cases Handled Today

- The local skip-auth path supplies stable default identity values when callers
  do not provide override headers.
- Group overrides are parsed from a comma-separated header and empty values are
  ignored.

## Performance and Operational Considerations

- Runtime behavior is currently limited to local identity injection for
  development and testing scenarios.
- Public contract docs still describe bearer-token authentication even though
  the verifier implementation is not wired into the current branch.

## Future Improvements

- Add explicit middleware unit/integration tests.
- Implement and wire a real bearer-token/Cognito verifier.
- Support configurable token header name for proxy variations.
- Improve unauthorized response detail for operator troubleshooting while preserving security posture.

## Assumptions

- **High confidence:** Cognito enforcement is currently documented intent, not
  live runtime behavior, on this branch.
