# Feature: Skip Auth Middleware

## Feature Overview

Provides a stable local development identity when `SKIP_AUTH=true`.

## Business Logic

- Read optional override headers:
  - `x-skip-auth-sub`
  - `x-skip-auth-username`
  - `x-skip-auth-scope`
  - `x-skip-auth-groups`
- Fall back to deterministic local defaults when override headers are absent.
- Add `sub`, `username`, `scope`, and `groups` values to Fiber locals.
- Continue request handling without external token verification.

## Package Location

- `api/middleware/middleware.go`
- `api/app.go`

## Key Structs and Interfaces

- `SkipAuthMiddleware`
- `parseGroups`

## Real Code Excerpt

```go
scope := c.Get(skipAuthHeaderScope)
if scope == "" {
    scope = defaultSkipAuthScope
}

c.Locals("sub", sub)
c.Locals("username", username)
c.Locals("scope", scope)
c.Locals("groups", groups)
```

## Edge Cases Handled Today

- Missing override headers still produce a stable local identity.
- Empty `x-skip-auth-groups` values are filtered out after splitting on commas.
- The middleware is only attached when `cfg.Core.SkipAuth` is true.

## Performance and Operational Considerations

- The middleware is fully in-process and does not call an external identity
  provider.
- Local identity defaults make test and dev workflows deterministic.
- The current branch does not implement Cognito JWT verification middleware.

## Future Improvements

- Add a separate production auth middleware once runtime auth enforcement is
  implemented.
- Document how local identity locals should map to any future bearer-token
  claims model.

## Assumptions

- **High confidence:** This page reflects the only auth-related middleware
  currently wired in the Go application.
