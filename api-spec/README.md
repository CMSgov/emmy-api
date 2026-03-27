# API Specification

This directory contains the hand-authored OpenAPI source for the Verification
Service API.

The team should treat the files here as the design-time source of truth for the
public API contract. We do not generate this specification from application
code. Instead, we design and review the contract first, then implement the
service to match the approved contract.

## What Lives Here

- `openapi.yaml`: the root OpenAPI document
- `paths/`: operation and path definitions
- `components/`: reusable OpenAPI components such as requests, responses,
  examples, and security schemes
- `dist/`: generated bundled artifacts for tooling and runtime use

## Working Agreement

- Edit source files in `api-spec/`; do not hand-edit files in `api-spec/dist/`
- Prefer small, reviewable files with `$ref` links instead of large inline
  definitions
- Keep API behavior, examples, and error contracts explicit in the spec
- Treat schema and contract changes as product and integration changes, not just
  implementation details

## Recommended Workflow

1. Design or update the contract in the source files under `api-spec/`
1. If the change affects shared data structures, update the standalone JSON
   Schemas under `schema/` first and reference them from OpenAPI
1. Review the contract diff in pull request before implementation work begins
1. Rebuild the bundled artifacts
1. Validate and lint the bundled spec
1. Implement or update service code to match the approved contract

## Bundling And Checks

Use the repo scripts from the project root:

```sh
./scripts/bundle-api-spec
./scripts/validate-api-spec
./scripts/lint-api-spec
```

Or, with `mise`:

```sh
mise install
mise run check-api-spec
```

Bundling produces checked-in artifacts in `api-spec/dist/`. Those artifacts are
for machine consumption and runtime serving; they should always be regenerated
from the source files in this directory.

## Pull Request Expectations

- Make contract intent clear in the PR description
- Include example request and response updates when behavior changes
- Call out breaking changes explicitly
- Keep generated `dist/` updates in the same PR as the source changes
- Prefer reviewing contract changes separately from implementation when possible

## Design Notes

This directory defines API surface area: operations, authentication
requirements, headers, status codes, examples, and response shapes.

The goal is for `api-spec/` to describe how clients integrate with the service,
while `schema/` describes reusable data structures that can evolve into a
broader data standard outside this specific API.
