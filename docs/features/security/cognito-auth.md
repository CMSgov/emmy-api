# Feature: Request Identity Handling

## Feature Overview

Derives request identity metadata for handlers and supports a local skip-auth
mode for development workflows. The current branch does not perform Cognito
signature verification.

## Business Logic

- If `SkipAuthMiddleware` has already populated locals, preserve that identity.
- Otherwise read `X-Sub`, or parse the `sub` claim from an
  `Authorization: Bearer ...` token without verification.
- Fall back to `unknown-subject` when no subject information is available.
- When `SKIP_AUTH=true`, populate `sub`, `username`, `scope`, and `groups`
  locals from stable defaults or `x-skip-auth-*` override headers.

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

- Invalid bearer token parsing returns `401`.
- Missing auth headers do not reject the request; they fall back to
  `unknown-subject`.
- Skip-auth header overrides support local impersonation without changing code.

## Performance and Operational Considerations

- No remote JWKS fetch or token introspection occurs on this branch.
- Subject extraction is cheap, but it should not be treated as real auth
  enforcement.
- `SkipAuthMiddleware` is only added when `SKIP_AUTH=true`.

## Future Improvements

- Reintroduce a signature-verifying bearer-token middleware when contract
  requirements and runtime behavior are aligned.
- Add explicit middleware unit/integration tests for verified auth if that
  middleware returns.
- Separate subject extraction from authorization enforcement more clearly in
  docs and code.

## Assumptions

- **High confidence:** The current middleware is identity plumbing, not a full
  authentication control.
