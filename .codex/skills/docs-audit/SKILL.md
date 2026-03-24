---
name: docs-audit
description: Perform repository documentation accuracy audits that compare docs claims to current-branch runtime code truth and produce a dated report in docs/audit with mandatory Sections A-E, severity-ranked findings, score, and PASS/FAIL gate. Use when asked to audit root/docs markdown accuracy, generate recurring audit reports, validate setup/API contract documentation, or produce evidence-backed documentation backlog items.
---

# Documentation Accuracy Audit

## Audit Inputs

Set these values before auditing:

- `report_date=<YYYY-MM-DD_hh:mm:ss>`
- `scope=root docs + docs/**/*.md` by default (expand only when explicitly requested)

Use current checked-out branch files as source of truth.

## Required Workflow

1. Audit docs claims against runtime truth.
- Verify runtime contract and wiring first: API behavior, auth, health/status, initialization flow, dependency injection, and error semantics.
- Evaluate setup and contributor docs next: setup steps, env vars, tooling instructions, placeholders, and broken references.
- Treat runtime code as canonical when docs conflict.

2. Classify findings and score.
- Assign severity:
  - `P1`: Contract or operationally misleading statements in primary docs.
  - `P2`: Stale setup/runtime guidance, broken references, or placeholders that affect contributors.
  - `P3`: Lower-risk consistency/editorial issues.
- Compute one overall accuracy score out of `100`.
- Set completion gate:
  - `PASS` only if no meaningful unresolved `P1` or `P2` findings remain.
  - Otherwise `FAIL`.

3. Capture evidence and reproducibility.
- Cite every material finding with repository file paths and line references.
- Keep caveats inline at the claim location; do not hide primary caveats only in late sections.
- Include a reproducible command log in Section D.

4. Run hygiene checks.
- Run a sensitive-file scan for private credential artifacts in tracked files.
- Run markdown link/reference/path validation for reviewed docs.

## Output Contract

Produce exactly one Markdown report at:

- `docs/audit/<report_date>.md`

Use exactly these sections and order:

1. `## Section A: Executive Summary`
2. `## Section B: Severity-Ranked Findings`
3. `## Section C: Update Backlog Checklist by Doc File`
4. `## Section D: Hygiene Appendix`
5. `## Section E: Deferred Watchlist (Non-blocking)`

Mandatory content requirements:

- Include score (`NN/100`), risk statement, and gate (`PASS` or `FAIL`) in Section A.
- Use the required Section B table columns and sort findings by severity then impact.
- Map Section C checklist items directly to findings.
- Include commands, sensitive-file scan, markdown validation, and hygiene verdict in Section D.
- Keep Section E as non-blocking triggers mapped to impacted docs.

## Template

Use this starter and populate with current-branch evidence:

```markdown
# Documentation Accuracy Review (as of <Month DD, YYYY>)

## Section A: Executive Summary
<...>

## Section B: Severity-Ranked Findings
| Severity | Doc location | Observed mismatch | Source-of-truth evidence on current branch | Recommended fix |
|---|---|---|---|---|

## Section C: Update Backlog Checklist by Doc File
<...>

## Section D: Hygiene Appendix
<...>

## Section E: Deferred Watchlist (Non-blocking)
<...>
```

For a copyable template file, use `references/documentation-audit-template.md`.

## Quality Bar

- Be precise, evidence-backed, and neutral.
- Prioritize behavior accuracy over style commentary.
- Do not invent findings without code/doc evidence.
- Explicitly state when a severity band has no findings.
- Prefer repository-local evidence; avoid external links unless requested.
