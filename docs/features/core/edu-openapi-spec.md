# Education Verification OpenAPI 3.1 Contract

## Intent

This feature introduces a design-first OpenAPI 3.1 contract for the public
education verification service endpoint, `POST /v1/edu`.

The contract is intentionally minimal and user friendly while preserving
compliance and consent-token requirements.

## What Is New

- Canonical public endpoint: `POST /v1/edu`
- Canonical request shape:
  - Required header: `X-EMMY-Consent-Token` (non-empty string)
  - Required: `applicant.firstName`, `applicant.lastName`,
    `applicant.dateOfBirth`
  - Optional: `applicant.ssnLast4`
  - Optional `clientReferenceId` for caller correlation
- Canonical response shape:
  - `requestId`
  - `status` (`code`, `severity`, `message`)
  - `result` (`verified`, `matchFound`, `hasEnrollmentRecords`,
    `matchedSchoolCount`)
  - `transaction` metadata and charges
  - `schools` with enrollment snapshots
- Error envelope for non-2xx responses:
  - RFC 7807 `application/problem+json`

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

## Sample Flow

1. Client sends `POST /v1/edu` with consent token header and canonical identity payload.
1. EDU service validates schema, auth, and required consent token header.
1. EDU service executes verification through internal dependency orchestration.
1. EDU service returns normalized education verification results and enrollment summary.
1. If validation/auth/dependency errors occur, EDU service returns RFC 7807 problem details.

## Spec and Governance Files

- `api-spec/openapi.yaml`
- `api-spec/paths/edu.yaml`
- `api-spec/components/`
- `api-spec/dist/openapi.bundled.yaml`
- `api-spec/dist/openapi.bundled.json`
- `.spectral.yaml`
