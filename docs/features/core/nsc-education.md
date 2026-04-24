# Feature: NSC Education

## Feature Overview

Provides education verification through NSC by accepting a caller JSON payload
and submitting a structured upstream request with OAuth2 client-credentials
authentication.

## Business Logic

- Parse the request body into `education.Request`.
- Serialize request body to JSON.
- Create HTTP POST to `NSC_SUBMIT_URL`.
- Use OAuth2-enabled HTTP client sourced from NSC token endpoint.
- Decode NSC JSON response into `education.Response`.

## Package Location

- `pkg/education/service.go`
- `pkg/education/submit.go`
- `pkg/education/oauth.go`
- `api/handlers/education_handler.go`
- `api/routes/router.go`

## Key Structs and Interfaces

- `EducationService`
- `HTTPTransport`
- `Options`
- `Request`
- `Response`
- `service` (concrete implementation)

## Real Code Excerpt

```go
req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.SubmitURL, bytes.NewReader(body))
if err != nil {
    return Response{}, fmt.Errorf("create submit request: %w", err)
}

resp, err := s.client.Do(req)
if err != nil {
    return Response{}, fmt.Errorf("submit request: %w", err)
}
```

## Edge Cases Handled Today

- Optional timeout injection if caller context has no deadline.
- Non-2xx NSC response returns wrapped error with status code.
- Long error body logging is truncated (800 chars).
- JSON marshal/unmarshal failures are surfaced with context.
- Redirect flow preserves auth header in OAuth client transport.

## Performance and Operational Considerations

- Network latency is measured and logged per submit call.
- Service relies on upstream token + submit endpoint availability.
- Timeout control exists but is optional unless configured.
- No retry/backoff logic in submit path yet.

## Future Improvements

- Add richer request validation beyond the current required
  `firstName`/`lastName`/`dateOfBirth` checks.
- Introduce retry policy with bounded backoff for transient 5xx errors.
- Add contract tests against NSC sandbox with fixtures.

## Assumptions

- **High confidence:** The HTTP handler now accepts caller-provided payloads and
  is no longer the earlier hardcoded-request scaffold.
