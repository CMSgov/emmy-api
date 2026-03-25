# Education Verification OpenAPI 3.1 Contract

## Intent

This feature introduces a design-first OpenAPI 3.1 contract for the public
education verification service endpoint, `POST /v1/edu`.

The contract is intentionally minimal and user friendly while preserving
compliance, authentication, and consent-token requirements.

## What Is New

- Education: `POST /v1/edu`
- Spec artifact endpoint: `GET /api-spec/v1/verify`
- Security requirements:
  - Bearer token authentication
  - Required header: `X-EMMY-Consent-Token` (non-empty string)
- Education request shape:
  - Required top-level object: `applicant`
  - Required applicant fields: `firstName`, `lastName`, `dateOfBirth`
  - Optional applicant fields: `middleName`, `ssn`
  - `additionalProperties: false` on the request and applicant objects
- Education success response shape:
  - `currentlyEnrolled` with enum values `Y` or `N`
  - `enrollmentStatus` with enum values `F`, `Q`, `H`, or `L`
- Education error envelope for non-2xx responses:
  - RFC 7807 style `application/problem+json`
  - Shared fields: `type`, `title`, `status`, optional `detail`, optional `instance`
- Explicit error responses:
  - `400` invalid request
  - `401` authentication failed
  - `403` authenticated but not authorized
  - `429` throttled or blocked by protection controls
  - `502` downstream dependency failure
  - `503` service temporarily unavailable

## Service Ownership and Boundary

The EDU service defines and owns the public API contract. Consumers integrate to
this contract directly.

Dependency-specific request/response formats, routing controls, and integration
mechanics are internal implementation details and are not part of the public
contract.

## Consent Header

`X-EMMY-Consent-Token` is mandatory for request acceptance:

- Header name: `X-EMMY-Consent-Token`
- Location: request header
- Requirement: required, non-empty string

These controls preserve regulatory intent while keeping the external API
minimal.

## Authentication

The contract now also requires bearer authentication for `POST /v1/edu`:

- Scheme: HTTP bearer
- Bearer format: JWT
- Scope model: no OAuth scopes are declared in the contract

Authentication is defined both globally and on the EDU operation so the
security requirement is visible in the Education path definition.

The bundled OpenAPI JSON endpoint, `GET /api-spec/v1/verify`, is also served
through the application’s existing middleware stack and therefore remains
behind the current bearer-auth behavior.

## Spec Artifact Endpoint

The service also exposes the checked-in bundled OpenAPI JSON artifact at:

- `GET /api-spec/v1/verify`

Behavior:

- Returns the contents of `api-spec/dist/openapi.bundled.json`
- Responds with `Content-Type: application/json`
- Does not rebuild or transform the OpenAPI document at request time
- Uses the checked-in bundled JSON file as the runtime source of truth

## Request and Response Notes

The request body has been simplified to the minimal identity payload needed for
an education enrollment lookup:

- `applicant.firstName`
- `applicant.lastName`
- `applicant.dateOfBirth`
- Optional `applicant.middleName`
- Optional `applicant.ssn`

The success response is also intentionally narrow. Instead of a larger
transaction-and-schools envelope, the current contract returns a normalized
enrollment result object with only:

- `currentlyEnrolled`
- `enrollmentStatus`

This keeps the public contract focused on the verification outcome while leaving
dependency-specific detail outside the API boundary.

## Error Model

All defined non-2xx responses reuse the shared `ProblemDetails` schema under
`application/problem+json`.

Recent spec updates expanded the documented error surface to include rate-limit,
downstream, and temporary availability failures in addition to validation and
authorization cases. The spec also includes concrete named examples for these
problem responses under `api-spec/components/examples/responses/problems/`.

## Sample Flow

1. Client sends `POST /v1/edu` with bearer auth, the consent token header, and the Education applicant payload.
1. EDU service validates schema, authentication, and the required consent token header.
1. EDU service executes verification through internal dependency orchestration.
1. EDU service returns the normalized enrollment result with `currentlyEnrolled` and `enrollmentStatus`.
1. Clients that need the machine-readable contract can request `GET /api-spec/v1/verify` to retrieve the bundled OpenAPI JSON artifact.
1. If validation, auth, throttling, or dependency errors occur, EDU service returns RFC 7807 problem details.

## Spec and Governance Files

- `api-spec/openapi.yaml`
- `api-spec/paths/edu.yaml`
- `api-spec/components/security-schemes.yaml`
- `api-spec/components/`
- `api-spec/dist/openapi.bundled.yaml`
- `api-spec/dist/openapi.bundled.json`
- `.spectral.yaml`
