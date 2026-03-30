---
name: docs-technical-research
description: Create and update technical research papers in docs/research using the repository's standard structure: Problem Statement, Alternatives Considered, Trade-offs, Why Current Approach Was Selected, Benchmarks / Status, References, and Assumptions. Use when asked to write or revise architecture/implementation research notes, compare design options, or document evidence-backed rationale in this repo.
---

# Technical Research Documentation

## Workflow

1. Identify scope and target file.
- Default location: `docs/research/`.
- Convert the topic into a kebab-case filename (for example, `status-endpoint-dependency-injection.md`).
- Reuse an existing file when the user asks for an update; create a new file when the topic is new.

2. Align with repository format.
- Follow this section order exactly:
  - `## Problem Statement`
  - `## Alternatives Considered`
  - `## Trade-offs`
  - `## Why Current Approach Was Selected (Observed/Inferred)` or `(Inferred)` depending on evidence level
  - `## Benchmarks / Status`
  - `## References`
  - `## Assumptions`
- Mirror the concise style used in `docs/research/*.md`.

3. Ground claims in repository evidence.
- Read relevant code/docs before writing conclusions.
- Use concrete file references in `## References` (for example, `main.go`, `api/app.go`, `docs/api.md`).
- Mark uncertain reasoning as inferred; do not present assumptions as facts.

4. Write actionable, decision-oriented analysis.
- In `Alternatives Considered`, list plausible options, including current approach.
- In `Trade-offs`, provide balanced pros/cons for each option.
- In `Benchmarks / Status`, explicitly state when data is unavailable.
- In `Assumptions`, label confidence (`High`, `Medium`, or `Low`) for each key assumption.

5. Validate final output quality.
- Ensure section names and order match the standard.
- Ensure file is in `docs/research/` unless user explicitly requests a different location.
- Ensure references are repository-local paths rather than external links unless the task explicitly requires external sources.
- Verify every path listed under `## References` exists in the repository before finalizing; correct or remove any invalid path.

## Template

Use this starter and customize by topic:

```markdown
# Research: <Topic>

## Problem Statement
<What decision or gap this research resolves>

## Alternatives Considered
- <Option A>
- <Option B>
- <Option C (current or recommended)>

## Trade-offs
- <Option A>:
  - Pros: <...>
  - Cons: <...>
- <Option B>:
  - Pros: <...>
  - Cons: <...>

## Why Current Approach Was Selected (Observed/Inferred)
<Evidence-backed rationale; mark inferred reasoning clearly>

## Benchmarks / Status
- <Existing data or "Not available">
- <Current implementation/documentation status>

## References
- `<repo/path-1>`
- `<repo/path-2>`

## Assumptions
- **High confidence:** <...>
- **Medium confidence:** <...>
```

For a copyable template file, use `references/research-paper-template.md`.
