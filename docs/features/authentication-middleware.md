# Feature: Authentication Middleware

## Feature Overview
Validates AWS Cognito access tokens for incoming requests when auth is enabled.

## Business Logic
- Read token from `x-amzn-oidc-accesstoken`.
- Load JWKS from Cognito issuer URL.
- Parse and validate JWT claims/signature.
- Enforce `client_id` match with configured app client.
- Add selected claims (`sub`, `username`, `scope`, `groups`) to Fiber locals.

## Package Location
- `api/middleware/middleware.go`
- `api/app.go`

## Key Structs and Interfaces
- `CognitoConfig`
- `CognitoVerifier`
- `NewCognitoVerifier`
- `FiberMiddleware`

## Real Code Excerpt
```go
tok, err := jwt.Parse(
    []byte(raw),
    jwt.WithKeySet(keyset),
    jwt.WithValidate(true),
    jwt.WithIssuer(v.issuer),
    jwt.WithClaimValue("token_use", "access"),
)
if err != nil {
    return fiber.ErrUnauthorized
}
```

## Edge Cases Handled Today
- Missing token header returns `401`.
- JWKS retrieval failures return unauthorized error.
- Invalid or mismatched `client_id` returns `401`.
- Config validation blocks startup if required cognito settings are missing.

## Performance and Operational Considerations
- JWKS uses `jwk.Cache` to avoid repeated key fetches.
- Request-time auth check includes a 5-second context timeout.
- Middleware is globally applied unless `SKIP_AUTH=true`.

## Future Improvements
- Add explicit middleware unit/integration tests.
- Support configurable token header name for proxy variations.
- Improve unauthorized response detail for operator troubleshooting while preserving security posture.

## Assumptions
- **High confidence:** Current claim checks are intentionally minimal and focused on access-token validity plus client binding.
