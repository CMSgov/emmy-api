---
name: docs-agent-md
description: Create, refresh, or reconcile the repository root `AGENTS.md` using current branch truth. Use when Codex needs to update `AGENTS.md` after code, config, CI, audit, or documentation changes, or when the repo's agent guidance must be checked against current files without widening scope to the rest of the docs set.
---

# Docs Agent MD

Use this skill to maintain the repository root `AGENTS.md` as a minimal,
repo-grounded document.

## Required Inputs

- Repository root containing `AGENTS.md`
- Current branch files used as source of truth

If `AGENTS.md` is missing, create it at the repository root and still keep the
scope limited to that file.

## Audit Check

Before updating `AGENTS.md`, determine the current local date in `YYYY-MM-DD`
format and inspect `docs/audit/`.

- If any audit filename starts with today's date, treat that as the current
  audit input.
- If no audit filename starts with today's date, ask the user whether to run
  `$docs-audit` before continuing.
- If the user wants the audit, use `$docs-audit` first and then continue with
  the refreshed report.

Accept both date-only and timestamped report names, for example:

- `docs/audit/2026-03-24.md`
- `docs/audit/2026-03-24_09:15:00.md`

## Workflow

1. Read the current `AGENTS.md`.
2. Check for a same-day audit in `docs/audit/`.
3. Inspect current repository truth before editing:
   - `git status --short`
   - recent commits such as `git log --oneline -n 10`
   - runtime and config files referenced by `AGENTS.md`
   - CI and policy files referenced by `AGENTS.md`
4. Use `references/agents-md-template.md` as the required output structure.
5. Update only `AGENTS.md`.
6. Re-read the final file for structure, wording, and claim accuracy.

## Content Rules

- Keep the document grounded in repository-observable facts.
- Prefer code, config, and CI over prose docs when they conflict.
- Preserve the current 5-section layout unless the user explicitly asks for a
  different structure.
- Keep the document minimal; avoid expanding it into a policy manual.
- State when a rule is observed guidance rather than hard enforcement.
- Do not invent ownership, deployment guarantees, or runtime behavior.

## Section Contract

Follow the exact section order in `references/agents-md-template.md`:

1. `## 1. Purpose`
2. `## 2. Repository Overview`
3. `## 3. Agent Roles`
4. `## 4. Coding Standards`
5. `## 5. Safety & Constraints`

## Boundaries

Allowed work:

- Read docs, audits, code, config, CI, and git metadata needed to verify
  `AGENTS.md`
- Update `AGENTS.md`

Default write scope:

- `AGENTS.md` only

Do not do this by default:

- Rewrite broader repository docs
- Edit `.github/workflows/*`, policy files, or API contract artifacts as part of
  an `AGENTS.md` refresh
- Expand the task into a full docs audit when the user only asked for
  `AGENTS.md`

## Redirects

- Use `$docs-audit` for full documentation accuracy audits and dated audit
  reports in `docs/audit/`.
- Use `$docs-feature` for `docs/features/*` authoring and index reconciliation.
- Use other repository docs skills when the target is not the root `AGENTS.md`.

## Bundled Resources

- `references/agents-md-template.md`: Required output template and section
  intent for the root `AGENTS.md`.
