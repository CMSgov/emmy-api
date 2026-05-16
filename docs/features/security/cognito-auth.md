# Feature: Request Subject Context

## Feature Overview

Captures a request subject for logging, reporting, and handler context.

## Business Logic

- If `SkipAuthMiddleware` already ran, preserve the injected local identity.
- Otherwise read `X-Sub` when present.
- Otherwise parse the `Authorization: Bearer <token>` value without signature
  verification and copy the JWT `sub` claim.
- Fall back to `unknown-subject` when neither source is available.
- When `SKIP_AUTH=true`, allow local override headers such as
  `x-skip-auth-sub` and `x-skip-auth-groups`.

## Package Location

- `api/middleware/middleware.go`
- `api/app.go`

## Key Structs and Interfaces

- `SubjectMiddleware`
- `SkipAuthMiddleware`
- `parseGroups`

## Real Code Excerpt

```go
auth := c.Get(fiber.HeaderAuthorization)
if strings.HasPrefix(auth, "Bearer ") {
    tokenString := strings.TrimPrefix(auth, "Bearer ")
    token, err := jwt.ParseString(tokenString, jwt.WithVerify(false), jwt.WithValidate(false))
    if err != nil {
        return fiber.ErrUnauthorized
    }
    sub = token.Subject()
}
```

## Edge Cases Handled Today

- Missing headers fall back to `unknown-subject`.
- Malformed bearer tokens return `401`.
- `SkipAuthMiddleware` injects stable local defaults when auth is disabled.
- Comma-separated `x-skip-auth-groups` values are trimmed and normalized.

## Performance and Operational Considerations

- Subject extraction is cheap and local; no JWKS fetch or external auth call is
  performed in the current branch.
- Middleware is globally applied so reporting always has a `client_id`-like
  value to log, even if it falls back to `unknown-subject`.
- This behavior does not enforce the OAuth 2.0 client-credentials contract
  described in `api-spec/v0/openapi.yaml`.

## Future Improvements

- Reintroduce a verifying auth middleware if the runtime is expected to enforce
  the checked-in OAuth contract.
- Add explicit docs once token verification, header semantics, and local-dev
  behavior are finalized.
- Keep security feature docs aligned with the runtime and spec separately
  instead of merging the two concerns.

## Assumptions

- **High confidence:** The current middleware stack propagates request identity
  but does not validate external tokens beyond parsing a bearer token for `sub`.
