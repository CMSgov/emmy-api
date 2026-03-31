# Feature: Skip Auth Middleware

## Feature Overview

Injects a stable local identity into Fiber request locals when `SKIP_AUTH=true`.
This is the only request-auth-related middleware currently wired in the runtime.

## Business Logic

- Read optional local override headers:
  - `x-skip-auth-sub`
  - `x-skip-auth-username`
  - `x-skip-auth-scope`
  - `x-skip-auth-groups`
- Fall back to stable defaults when those headers are absent.
- Add selected values (`sub`, `username`, `scope`, `groups`) to Fiber locals.
- Skip any token parsing or upstream auth validation entirely.

## Package Location

- `api/middleware/middleware.go`
- `api/app.go`

## Key Structs and Interfaces

- `SkipAuthMiddleware`
- `parseGroups`

## Real Code Excerpt

```go
groupsHeader := c.Get(skipAuthHeaderGroups)
groups := []string{defaultSkipAuthGroup}
if groupsHeader != "" {
    parsed := parseGroups(groupsHeader)
    if len(parsed) > 0 {
        groups = parsed
    }
}

c.Locals("sub", sub)
c.Locals("username", username)
c.Locals("scope", scope)
c.Locals("groups", groups)
```

## Edge Cases Handled Today

- Missing override headers fall back to stable local defaults.
- Empty or whitespace-only group entries are ignored.
- `x-skip-auth-groups` supports comma-separated values.
- `/health` is registered before this middleware and remains unauthenticated.

## Performance and Operational Considerations

- The middleware is only installed when `SKIP_AUTH=true`.
- Header parsing is in-process and does not trigger network calls.
- When `SKIP_AUTH=false`, the current branch does not install any request-auth
  middleware, so runtime auth enforcement remains a known gap.

## Future Improvements

- Restore a real request-auth path for the non-skip-auth case so runtime
  behavior matches the published contract.
- Add route-level tests that verify when middleware is and is not applied.
- Document how local skip-auth should coexist with future OAuth2 enforcement.

## Assumptions

- **High confidence:** This middleware exists for local/development identity
  injection, not for production-grade request authentication.
