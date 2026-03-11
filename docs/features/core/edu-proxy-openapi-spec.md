# Education Proxy OpenAPI 3.1 Contract

## Intent

This feature introduces a design-first OpenAPI 3.1 contract for a public
education verification proxy endpoint, `POST /v1/edu`.

The contract is intentionally minimal and user friendly while preserving
compliance and consent attestation requirements.

## What Is New

- Canonical public endpoint: `POST /v1/edu`
- Canonical request shape:
  - Required: `applicant.firstName`, `applicant.lastName`,
    `applicant.dateOfBirth`
  - Optional: `applicant.ssnLast4`
  - Required structured consent:
    `consent.attested` (must be `true`), `consent.attestedAt`, and
    `consent.purpose`
  - Optional `clientReferenceId` for caller correlation
- Canonical response shape:
  - `requestId`
  - `status` (`code`, `severity`, `message`)
  - `result` (`verified`, `nscHit`, `hasEnrollmentRecords`,
    `matchedSchoolCount`)
  - `transaction` metadata and charges
  - `schools` with enrollment snapshots
- Error envelope for non-2xx responses:
  - RFC 7807 `application/problem+json`

## Mapping Rationale (Canonical -> Provider)

The proxy maps user-friendly canonical fields to provider payload fields and
keeps provider-owned fields internal.

Selected mapping examples:

- `applicant.firstName` -> NSC `firstName`
- `applicant.lastName` -> NSC `lastName`
- `applicant.dateOfBirth` -> NSC `dateOfBirth`
- `applicant.ssnLast4` -> NSC `ssn` (last-4 handling policy is implemented by proxy)
- `consent.attested=true` -> NSC `terms="y"` (internal mapping)
- `clientReferenceId` -> NSC `caseReferenceId` (when provided)

Provider account and routing controls such as `accountId`, `endClient`, and
other service options are not public contract fields in v1.

## Compliance Fields

Structured consent fields are mandatory for request acceptance:

- `consent.attested` must be `true`
- `consent.attestedAt` must be a valid RFC 3339 date-time
- `consent.purpose` must be a non-empty business/legal purpose statement

These controls preserve regulatory intent while keeping the external API
minimal.

## Sample Flow

1. Client sends `POST /v1/edu` with canonical identity and consent payload.
1. Proxy validates schema, auth, and consent fields.
1. Proxy maps canonical payload to provider request format.
1. Proxy returns normalized education verification result and enrollment summary.
1. If request/auth/upstream errors occur, proxy returns RFC 7807 problem details.

## Spec and Governance Files

- `api-spec/openapi.yaml`
- `api-spec/paths/edu.yaml`
- `api-spec/components/`
- `api-spec/dist/openapi.bundled.yaml`
- `api-spec/dist/openapi.bundled.json`
- `.spectral.yaml`

