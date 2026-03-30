# Feature: Verification API v0 Contract

## Feature Overview
Defines the current public verification API contract in OpenAPI 3.1. The
versioned spec under `api-spec/v0/` and the reusable schemas under `schema/v0/`
are the source of truth for intended API behavior in this repository.

## Business Logic
- Expose two public operations:
  - `POST /v0/education-enrollments`
  - `POST /v0/veteran-disability-ratings`
- Apply the global `OAuth2ClientCredentials` security scheme to the contract.
- Reuse the shared `Identity` schema for both request bodies.
- Return operation-specific response shapes:
  - `enrollmentStatus` for education verification
  - `combinedDisabilityRating` for veteran verification

## Package Location
- `api-spec/v0/openapi.yaml`
- `schema/v0/identity.schema.json`
- `schema/v0/school_enrollment_status.schema.json`
- `schema/v0/combined_disability_rating.schema.json`
- `api-spec/README.md`

## Key Structs and Interfaces
- `OAuth2ClientCredentials`
- `GetEducationEnrollmentRequest`
- `EducationEnrollmentResponse`
- `GetVeteranDisabilityStatusRequest`
- `VeteranDisabilityStatusResponse`

## Real Code Excerpt
```yaml
paths:
  /v0/education-enrollments:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/GetEducationEnrollmentRequest"

components:
  schemas:
    GetEducationEnrollmentRequest:
      $ref: ../../schema/v0/identity.schema.json
```

## Edge Cases Handled Today
- Both operations reuse the same required identity fields: `firstName`,
  `lastName`, and `dateOfBirth`.
- `schema/v0/identity.schema.json` currently allows additional properties, so
  the request model is intentionally permissive beyond the required fields.
- The checked-in v0 OpenAPI currently documents `200` success responses only; it
  does not yet define a shared non-2xx error envelope.

## Performance and Operational Considerations
- Contract files are versioned so future breaking changes can land in a new API
  version without rewriting the v0 source files in place.
- The OpenAPI contract is implementation-agnostic: downstream provider payloads
  and current runtime scaffolding are intentionally outside the public boundary.
- Bundled artifacts are generated from the versioned source files and should be
  treated as build outputs, not authoring inputs.

## Future Improvements
- Add explicit non-2xx response documentation when the public error model is
  finalized.
- Publish consumer-facing examples and artifact-serving docs that point at the
  versioned bundle locations used in this repository.

## Assumptions
- **High confidence:** `api-spec/v0/openapi.yaml` and `schema/v0/*.json` define
  intended public API behavior for this branch.
