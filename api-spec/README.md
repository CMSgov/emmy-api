# API Specification

This directory contains the hand-authored OpenAPI source for the Verification
Service API.

The team should treat the versioned files here as the design-time source of
truth for the public API contract. We do not generate this specification from
application code. Instead, we design and review the contract first, then
implement the service to match the approved contract.

## What Lives Here

- `v0/openapi.yaml`: the current v0 root OpenAPI document
- `v0/dist/`: generated bundled artifacts produced from the v0 source file
- `dist/v0/`: additional generated artifacts consumed by tooling and runtime
- `README.md`: working agreement for contract-first API documentation in this
  repository

The public contract in this branch is currently defined in
`api-spec/v0/openapi.yaml` and references versioned JSON Schemas under
`schema/v0/`.

## Working Agreement

- Edit source files under `api-spec/v0/`; do not hand-edit generated artifacts
  under `api-spec/v0/dist/` or `api-spec/dist/v0/`
- Prefer reviewable, versioned contract documents over implementation-derived
  API descriptions
- Keep API behavior, examples, and error contracts explicit in the spec
- Treat schema and contract changes as product and integration changes, not just
  implementation details

## Recommended Workflow

1. Design or update the contract in the source files under `api-spec/v0/`
1. If the change affects shared data structures, update the standalone JSON
   Schemas under `schema/v0/` first and reference them from OpenAPI
1. Review the contract diff in pull request before implementation work begins
1. Rebuild the bundled artifacts
1. Validate and lint the bundled spec
1. Implement or update service code to match the approved contract

## Bundling And Checks

Use the supported commands from the project root:

```sh
pnpm run bundle:api-spec
pnpm run validate:api-spec
pnpm run lint:api-spec
```

Or, with `mise`:

```sh
mise install
mise run bundle-api-spec
mise run validate-api-spec
mise run lint-api-spec
```

Bundling produces checked-in versioned artifacts under `api-spec/v0/dist/` and
`api-spec/dist/v0/`. Those artifacts are for machine consumption and runtime
serving; they should always be regenerated from the source files in this
directory.

## Pull Request Expectations

- Make contract intent clear in the PR description
- Include example request and response updates when behavior changes
- Call out breaking changes explicitly
- Keep generated `dist/` updates in the same PR as the source changes
- Prefer reviewing contract changes separately from implementation when possible

## Design Notes

This directory defines API surface area: operations, authentication
requirements, status codes, examples, and response shapes.

The goal is for `api-spec/` to describe how clients integrate with the service,
while `schema/` describes reusable data structures that can evolve into a
broader data standard outside this specific API.
