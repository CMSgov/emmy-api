# Feature: Skip Auth Identity

## Feature Overview

Injects a stable local identity into Fiber locals when `SKIP_AUTH=true`.

## Business Logic

- Read optional override headers:
  - `x-skip-auth-sub`
  - `x-skip-auth-username`
  - `x-skip-auth-scope`
  - `x-skip-auth-groups`
- Fall back to deterministic local defaults when headers are absent.
- Add `sub`, `username`, `scope`, and `groups` values to Fiber locals.

## Package Location

- `api/middleware/middleware.go`
- `api/app.go`

## Key Structs and Interfaces

- `SkipAuthMiddleware`
- `parseGroups`

## Real Code Excerpt

```go
sub := c.Get(skipAuthHeaderSub)
if sub == "" {
    sub = defaultSkipAuthSub
}

scope := c.Get(skipAuthHeaderScope)
if scope == "" {
    scope = defaultSkipAuthScope
}

c.Locals("sub", sub)
c.Locals("scope", scope)
```

## Edge Cases Handled Today

- Missing override headers fall back to stable defaults.
- Empty values in `x-skip-auth-groups` are trimmed out.
- Empty `x-skip-auth-username` falls back to the resolved `sub`.

## Performance and Operational Considerations

- No network calls or token parsing happen in this middleware.
- Middleware is only added when `SKIP_AUTH=true`.
- This branch does not wire alternate request-auth middleware when `SKIP_AUTH=false`.

## Future Improvements

- Add the production request-auth path back to the docs once it returns to the branch.
- Document how upstream infrastructure should enforce auth when `SKIP_AUTH=false`.
- Add an explicit feature page for the non-local auth path if it lands as a separate middleware implementation.

## Assumptions

- **High confidence:** This file now documents the only auth-related middleware behavior observable in the current branch.
