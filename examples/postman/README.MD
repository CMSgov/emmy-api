# Postman Collections

This directory contains Postman workspace exports used during earlier v1
integration experiments.

## Current Status

- These files are **not** the source of truth for the current public contract.
- The intended public API contract for this branch lives in
  `api-spec/v0/openapi.yaml` and `schema/v0/`.
- Some requests in this directory still use older v1 paths and should be
  treated as legacy reference material until they are refreshed.

## What Lives Here

- `v1/postman/collections/`: legacy request definitions
- `v1/postman/environments/`: example environment variables for local/manual use

## Working Guidance

- Prefer the docs under `docs/guides/` and the versioned OpenAPI contract under
  `api-spec/v0/` when building new integrations.
- If these Postman exports are updated, keep them aligned with the checked-in
  contract rather than the current runtime scaffolding.
- Never commit real client credentials, access tokens, or other secrets in
  environment exports.

## Related Resources

- `docs/guides/01-getting-started.md`
- `docs/guides/02-authentication.md`
- `docs/guides/03-usage-examples.md`
- `api-spec/v0/openapi.yaml`
