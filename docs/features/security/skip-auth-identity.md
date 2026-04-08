# Feature: Skip Auth Identity

## Feature Overview

Injects a stable local identity into Fiber request locals when `SKIP_AUTH=true`.
This is a development/runtime convenience on the current branch, not the public
API authentication contract.

## Business Logic

- Read optional override headers:
  - `x-skip-auth-sub`
  - `x-skip-auth-username`
  - `x-skip-auth-scope`
  - `x-skip-auth-groups`
- Fall back to stable defaults when those headers are absent.
- Populate Fiber locals for downstream handlers:
  - `sub`
  - `username`
  - `scope`
  - `groups`

## Package Location

- `api/middleware/middleware.go`
- `api/app.go`
- `.env.example`

## Key Structs and Interfaces

- `SkipAuthMiddleware`
- `parseGroups`
- `core.Config.SkipAuth`

## Real Code Excerpt

```go
if cfg.Core.SkipAuth {
    app.Use(middleware.SkipAuthMiddleware())
}
```

## Edge Cases Handled Today

- Missing override headers fall back to stable defaults.
- Empty `x-skip-auth-groups` falls back to the default local-dev group.
- Comma-separated group values are trimmed before being stored.
- `GET /health` is registered before this middleware is attached.

## Performance and Operational Considerations

- The middleware is local-only request shaping; it does not perform crypto,
  network calls, or token validation.
- Enabling `SKIP_AUTH=true` makes local testing easier but does not exercise the
  bearer-token semantics described by the versioned OpenAPI contract.

## Future Improvements

- Replace or complement this middleware with a real runtime bearer-token
  validator when auth implementation work resumes.
- Add more explicit operator-facing docs once runtime auth and contract auth
  converge.

## Assumptions

- **High confidence:** This middleware exists to support local development and
  non-auth test paths on the current branch.
