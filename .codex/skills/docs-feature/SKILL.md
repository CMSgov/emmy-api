---
name: docs-feature
description: Create or update `docs/features/*` documentation for this repository using the established feature-doc structure and keep `docs/features/README.md` synchronized. Use when Codex needs to add a new feature page, revise an existing feature page, choose the correct feature category (`core`, `infrastructure`, `resilience`, `security`), or reconcile the features index with the files that actually exist on disk.
---

# Docs Feature

Use this skill to keep feature documentation consistent and indexable inside this repository.

## Required Inputs

- Repository root containing `docs/features/`
- Feature scope grounded in real code, config, or operational behavior

If `docs/features/` is missing, stop and ask the user to confirm the target docs location before writing.

## Repository Rules

- Use only these top-level categories:
  - `core`
  - `infrastructure`
  - `resilience`
  - `security`
- Do not create a new top-level category unless the user explicitly asks.
- Treat `docs/features/README.md` as the canonical feature index.
- Keep documentation grounded in the repository. Do not invent behavior that is not supported by code, config, or existing docs.
- Treat `docs/features/core/edu-openapi-spec.md` as a useful repository document, but not as the default formatting model for ordinary feature pages.

## Workflow

1. Identify the target feature and its category.
2. Inspect the existing docs in the same category before drafting content.
3. Use `references/feature-template.md` as the default structure for standard feature pages.
4. Gather evidence from code, config, routes, handlers, middleware, and operational files before writing claims.
5. Add or update the feature page under `docs/features/<category>/`.
6. Reconcile `docs/features/README.md` so the index matches the current files on disk.
7. Re-read the changed docs for structure, category fit, and wording consistency.

Do not skip the README reconciliation step when a feature page is added, moved, renamed, or removed.

## Category Selection

- Use `core` for business capabilities, public service behavior, request/response flows, and domain-owned contracts.
- Use `infrastructure` for runtime dependencies, observability, telemetry, data stores, and supporting platform services.
- Use `security` for authentication, authorization, identity, token validation, and security controls.
- Use `resilience` for failover behavior, admission control, circuit breaking, retries, and fault-tolerance patterns.

If a feature could fit multiple categories, prefer the category that best matches its primary operational purpose and note the rationale in your working notes.

## Feature Page Standard

- Default to the section order in `references/feature-template.md`.
- Keep headings and ordering stable unless the user explicitly asks for a different format.
- Prefer concise bullets over long narrative paragraphs.
- Include file paths for the main implementation touchpoints.
- Use a short real code or config excerpt when it clarifies behavior.
- Record uncertainty in the `Assumptions` section with an explicit confidence level.

## README Maintenance

Preserve this category order in `docs/features/README.md`:

1. Core
2. Infrastructure
3. Security
4. Resilience

Maintain the existing table format under each category:

- `Component`
- `Purpose`
- `Functionality`

When reconciling the README:

- Add entries for feature files that exist on disk but are missing from the index.
- Remove or fix entries whose linked files do not exist.
- Keep links relative to `docs/features/README.md`.
- Keep component names human-readable and aligned with the feature title.
- Keep purpose/functionality summaries short and concrete.

## Bundled Resources

- `references/feature-template.md`: Default section structure for standard feature pages in this repository.
