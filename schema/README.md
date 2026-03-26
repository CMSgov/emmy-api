# JSON Schema

This directory contains standalone JSON Schema documents that define shared data
models independent of any single API endpoint.

The intent is for these schemas to become reusable, version-controlled data
standards in their own right. OpenAPI documents in this repository should
reference these schemas instead of redefining the same data structures inline
when that separation is appropriate.

## Why This Directory Exists

- To separate reusable data-model design from endpoint design
- To allow schema-focused pull request review
- To support versioning and compatibility discussions at the data-standard layer
- To make shared models portable beyond this service implementation

## What Should Live Here

- Versioned JSON Schema files such as `*.schema.json`
- Shared domain models that may be reused across operations or services
- Schema examples or companion documentation when they improve reviewability

Avoid placing endpoint-specific concerns here, such as:

- HTTP methods
- status codes
- auth requirements
- request routing details
- API-specific error response wrappers

Those belong in `api-spec/`.

## Suggested Layout

As this directory grows, prefer organizing by domain and version, for example:

```text
schema/
  education/
    v1/
      applicant.schema.json
      enrollment.schema.json
```

Keep file names stable and descriptive so diffs are easy to review.

## Versioning Guidance

- Use explicit version directories when compatibility matters
- Treat removing fields, changing required fields, or narrowing allowed values
  as breaking changes
- Treat additive optional fields as backward-compatible changes
- Document intended migration expectations in the PR when introducing a new
  schema version

## Working Agreement

- Hand-author and review schema files in git
- Keep schemas focused on data shape and validation rules
- Prefer explicit constraints, examples, and descriptions where helpful
- Avoid mixing multiple conceptual resources into a single schema file
- Update OpenAPI references after schema changes are approved

## Pull Request Expectations

- Explain whether the change is breaking, additive, or editorial
- Include example instances when changing complex data structures
- Note which OpenAPI documents or services are expected to consume the change
- Keep schema-only changes reviewable without requiring implementation context

## Relationship To `api-spec/`

`schema/` defines reusable data structures.

`api-spec/` defines how those structures are used in the public API contract.

The preferred workflow is:

1. Design or revise the standalone schema here
1. Review and approve the schema change
1. Reference the schema from OpenAPI in `api-spec/`
1. Implement application code against the approved contract
