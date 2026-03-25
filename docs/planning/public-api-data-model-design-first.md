---
date: 2026-03-09
status: Proposed
scope: Public API contract definitions only
out_of_scope: Internal implementation objects and transient internal models
---

# Public API Data Model: Design-First Format and Tooling Decision

## Decision Summary

We will adopt a design-first API process with **OpenAPI 3.1 (YAML) as the canonical public API contract**.

Selected tooling stack:

- Authoring and preview: VS Code + SpecLynx
- Bundling for multi-file specs: Redocly CLI
- Validation: swagger-cli
- Linting and governance style rules: Spectral + repo-owned `.spectral.yaml`
- Breaking-change detection: oasdiff
- Documentation publishing: ReDoc for public docs (next iteration), Swagger UI for internal/dev debugging (initial phase)
- Mocking: Prism
- Go code generation: oapi-codegen (Fiber target)
- Contract/property-based testing: Schemathesis

This gives one source of truth that supports:

- Human-readable docs (ReDoc / Swagger UI)
- Programmatic consumers (OpenAPI YAML/JSON artifacts + code generation)
- Automatic contract enforcement (lint/validate/breaking checks + Schemathesis)

## Why This Decision

### Requirements Fit

| Requirement | Decision fit |
| --- | --- |
| Human-readable docs | ReDoc and Swagger UI render directly from OpenAPI |
| Programmatic processing | OpenAPI supports code generation, schema reuse, mocks, and diff tooling |
| Contract testing | Schemathesis runs generated tests from the OpenAPI contract |
| Design-first governance | Spec reviewed and approved before implementation changes |
| Go implementation alignment | oapi-codegen generates types and Fiber server interfaces |

### OpenAPI Version Choice

- OpenAPI 2.0: rejected (legacy)
- OpenAPI 3.0.x: acceptable, but not preferred for new work
- **OpenAPI 3.1.x: selected** (best balance of modern schema support and tooling maturity)
- OpenAPI 3.2.x: deferred until ecosystem support is stable for our required toolchain

Validation nuance:

- OpenAPI uses an OpenAPI Schema dialect, not pure JSON Schema.
- Use OpenAPI-aware validators and linters (`swagger-cli`, Spectral).
- Do not rely only on generic JSON Schema validation for contract correctness.

## Design-First Workflow (Required)

1. PRD defines resources, operations, consumers, constraints, and non-functional requirements.
1. (Optional) LLM digestible rules as a code in machine readable spec.
1. LLM creates an initial OpenAPI draft from the PRD.
1. Engineers refine the spec in editor with preview + linting.
1. CI validates and lints the spec, then runs breaking-change checks.
1. Spec PR requires peer approval before merge.
1. Merge publishes versioned artifacts (docs, mock, bundled specs, changelog, generated code/SDK updates).

Rule: Public API behavior is implemented to match the merged spec. The implementation does not define the contract.

## Go Governance Model

### Ownership and Review

- API spec changes require at least one reviewer approval from API owners.
- Spec PRs are gated by validation, lint, and breaking-change checks.
- Breaking changes require explicit approval label (for example `breaking-change-approved`) and corresponding version bump plan.

### Versioning and Deprecation

- Versioning strategy: URL path versioning (`/v1`, `/v2`) for public endpoints.
- Breaking change policy: never break a published major version in place.
- Deprecation policy: minimum 6-month notice with `Deprecation` and `Sunset` headers (RFC 8594).
- Changelog policy: update `CHANGELOG.md` from `oasdiff` output on each merged spec PR.
- SemVer tagging policy:
  - Patch: documentation/example clarifications with no API shape change
  - Minor: backward-compatible endpoint/field additions
  - Major: breaking changes

### CI Policy for Spec PRs

Required checks:

1. Bundle spec (if multi-file).
1. Validate syntax/structure.
1. Lint style and policy rules.
1. Fail on unapproved breaking changes.

Recommended optional checks:

1. Spin up Prism mock and run smoke tests.
1. Run Schemathesis against staging/ephemeral environment.

### Drift Detection and Runtime Conformance

- Use Schemathesis against deployed API for request/response schema conformance.
- Add runtime request/response validation middleware where practical in Go using OpenAPI-aware validators.
- Regularly compare deployed behavior against spec to catch undocumented endpoints or response drift.

## Repository Conventions

Recommended structure:

```text
api-spec/
  openapi.yaml
  components/
  paths/
  dist/
    openapi.bundled.yaml
    openapi.bundled.json
.spectral.yaml
docs/planning/public-api-data-model-design-first.md
```

Contract boundary:

- Include only public API request/response models and public endpoint behavior.
- Do not include internal DTOs, persistence entities, or internal service payloads unless externally visible.

Current EDU contract conventions in this repo:

- Public path: `POST /v1/edu`
- OpenAPI version: `3.1.0`
- Security model: HTTP bearer auth plus required `X-EMMY-Consent-Token` header on the EDU operation
- Request model style: narrow public request wrapper with reusable component schemas and `additionalProperties: false`
- Success response model style: normalized business outcome model instead of downstream vendor payload passthrough
- Error model style: reusable RFC 7807-style `ProblemDetails` schema with shared response components and named examples for each documented failure mode
- Current documented EDU error responses: `400`, `401`, `403`, `429`, `502`, and `503`

This repo’s recent spec updates reinforce the design-first rule that public docs
must track the contract as written in reusable OpenAPI components, not internal
service behavior or superseded response envelopes.

## Local Usage Instructions

### 1. Author or Update the Spec

- Create or update `api-spec/openapi.yaml` (and referenced files).
- Keep schema components reusable and explicitly named.
- Prefer explicit `operationId`, `summary`, response schemas, and examples.
- Define shared error responses as reusable components when the same problem
  schema is returned across operations.
- Keep public success payloads normalized and implementation-agnostic.

### 2. Bundle, Validate, and Lint

```bash
# 1. Bundle the checked-in YAML and JSON artifacts (required for multi-file refs)
./scripts/bundle-api-spec

# 2. Validate OpenAPI structure
./scripts/validate-api-spec

# 3. Lint style and governance rules
./scripts/lint-api-spec
```

Using `mise`, the same workflow is available as:

```bash
mise install
mise run bundle-api-spec
mise run validate-api-spec
mise run lint-api-spec
mise run check-api-spec
```

Starter `.spectral.yaml` (commit and version in repo):

```yaml
extends:
  - spectral:oas
rules:
  operation-summary-required:
    description: Every operation must define a summary.
    given: "$.paths[*][*]"
    severity: warn
    then:
      field: summary
      function: truthy
```

### 3. Preview Documentation

```bash
# Internal/dev preview
npx swagger-ui-watcher api-spec/openapi.yaml --no-open

# Public docs build (ReDoc CLI example)
npx redoc-cli build api-spec/dist/openapi.bundled.yaml -o api-spec/dist/index.html
```

### 4. Run Mock Server

```bash
# Mock server from contract
npx @stoplight/prism-cli mock api-spec/dist/openapi.bundled.yaml

# Proxy validation against a running environment
npx @stoplight/prism-cli proxy api-spec/dist/openapi.bundled.yaml https://api.example.gov
```

### 5. Generate Go Types and Server Interfaces

Add `oapi-codegen.yaml`:

```yaml
package: api
generate:
  - types
  - fiber-server
output: internal/api/api.gen.go
```

Generate code:

```bash
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -config oapi-codegen.yaml api-spec/dist/openapi.bundled.yaml
```

Team rule: choose one strategy and document it in `README.md`:

- Commit generated code to git, or
- Generate in CI only

Do not mix strategies across services.

### 6. Run Contract and Property-Based Tests

```bash
# Example against local server
schemathesis run --url http://localhost:8080 api-spec/dist/openapi.bundled.yaml
```

Use Schemathesis for contract conformance and edge-case generation. Keep separate hand-written integration tests for business logic and authorization behavior.

## Artifact Publication Requirements

After a spec PR merges, publish all of the following from the merged contract:

1. Documentation site (ReDoc or Swagger UI) on a stable URL.
1. Mock endpoint (Prism) for consumer integration.
1. Bundled machine-readable artifacts (`openapi.bundled.yaml`, `openapi.bundled.json`).
1. Generated Go server types/interfaces (and client SDKs if applicable).
1. Changelog entry from `oasdiff` output.

## CI Reference (GitHub Actions)

```yaml
- name: Bundle OpenAPI
  run: npx @redocly/cli bundle api-spec/openapi.yaml -o api-spec/dist/openapi.bundled.yaml

- name: Validate OpenAPI
  run: npx @apidevtools/swagger-cli validate api-spec/dist/openapi.bundled.yaml

- name: Lint OpenAPI
  run: npx @stoplight/spectral-cli lint api-spec/dist/openapi.bundled.yaml --ruleset .spectral.yaml

- name: Prepare base spec
  run: |
    git fetch --no-tags --depth=1 origin ${{ github.base_ref }}
    git show origin/${{ github.base_ref }}:api-spec/openapi.yaml > /tmp/openapi.base.yaml
    npx @redocly/cli bundle /tmp/openapi.base.yaml -o /tmp/openapi.base.bundled.yaml

- name: Detect breaking changes
  if: ${{ !contains(github.event.pull_request.labels.*.name, 'breaking-change-approved') }}
  run: oasdiff breaking /tmp/openapi.base.bundled.yaml api-spec/dist/openapi.bundled.yaml --fail-on ERR
```

Implementation note: use a stable "base spec" artifact from `main` in CI (download artifact or check out base commit) before running `oasdiff`.

## Selected Governance Defaults

- Spec format: OpenAPI 3.1 YAML
- Canonical artifact: `api-spec/openapi.yaml`
- Machine artifact: bundled YAML + JSON in `api-spec/dist/`
- Public docs renderer: ReDoc
- Internal debugging renderer: Swagger UI
- Linter: Spectral + repo-owned ruleset
- Breaking change gate: oasdiff (block by default)
- Mocking: Prism (shared sandbox)
- Go codegen: oapi-codegen (`fiber-server` + `types`)
- Contract testing: Schemathesis

## Risks and Mitigations

- Toolchain sprawl: mitigate with pinned tool versions in CI and documented local scripts.
- Spec drift: mitigate with runtime conformance checks + Schemathesis + breaking-change gate.
- Generated-code churn: mitigate by choosing one generation strategy (commit or CI-only) and enforcing it consistently.

## References

- [Stoplight: The Right Way to API (Design-First rationale)](https://blog.stoplight.io/the-right-way-to-api)
- [OpenAPI Specification 3.1.0](https://spec.openapis.org/oas/v3.1.0.html)
- [Spectral (Stoplight) repository and docs](https://github.com/stoplightio/spectral)
- [Redocly CLI `bundle` command](https://redocly.com/docs/cli/commands/bundle)
- [Prism (Stoplight) repository and CLI usage](https://github.com/stoplightio/prism)
- [oasdiff (breaking change detection)](https://www.oasdiff.com/)
- [Schemathesis (property-based API contract testing)](https://schemathesis.io/)
- [RFC 8594: The Sunset HTTP Header Field](https://www.rfc-editor.org/rfc/rfc8594)
