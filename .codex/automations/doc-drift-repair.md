# Doc Drift Repair Automation

This file captures the repository-owned definition for the recurring doc drift
repair flow.

## Intent

- Audit documentation against the current branch using repository-local skills.
- Fix safe documentation drift grounded in code and config.
- Perform the work on a dedicated branch named `codex/doc-drift-YYYY-MM-DD`.
- Stage, commit, and prepare a draft PR when the run produces a reviewable doc
  change set.

## Skills

- [`docs-audit`](../skills/docs-audit/SKILL.md)
- [`docs-feature`](../skills/docs-feature/SKILL.md)
- [`github:yeet`](/Users/ssnl/.codex/plugins/cache/openai-curated/github/f78e3ad49297672a905eb7afb6aa0cef34edc79e/skills/yeet/SKILL.md)

## Suggested Automation

```text
::automation-update{mode="suggested create" name="Doc Drift Repair" prompt="Audit documentation drift in this repository by comparing README.md, docs/, docs/features/README.md, routes, handlers, pkg/, main.go, api-spec sources, and runtime/config files. Treat code, config, and CI as source of truth when prose disagrees. Use [$docs-audit](/Users/ssnl/cmsgov/emmy-api/.codex/skills/docs-audit/SKILL.md) to generate a dated audit report in docs/audit and capture evidence-backed findings. Start each run on a dedicated branch that can become a pull request: create or switch to a branch named codex/doc-drift-YYYY-MM-DD for that run and keep all changes isolated there. Use [$docs-feature](/Users/ssnl/cmsgov/emmy-api/.codex/skills/docs-feature/SKILL.md) when updating docs/features pages or reconciling docs/features/README.md. When the audit produces safe, repo-grounded documentation fixes, stage them, create a focused commit, and prepare a draft pull request using [$github:yeet](/Users/ssnl/.codex/plugins/cache/openai-curated/github/f78e3ad49297672a905eb7afb6aa0cef34edc79e/skills/yeet/SKILL.md). Make only documentation changes in docs/, docs/features/, docs/audit/, and README.md when safe, and keep any example data synthetic. Do not change workflows, security policy files, deploy scripts, public.jwk, bundled OpenAPI outputs, or anything that would alter auth behavior or public API semantics without explicit approval. If drift is found in one of those restricted areas, leave files unchanged there and summarize the needed follow-up in the inbox item instead of forcing a PR. In the inbox item, report what drift was found, which files were updated, what branch was used, whether a draft PR was prepared, what was skipped, and any remaining risks or assumptions." rrule="FREQ=WEEKLY;BYDAY=MO;BYHOUR=9;BYMINUTE=0" cwds="/Users/ssnl/cmsgov/emmy-api" status="ACTIVE"}
```

## Notes

- The live scheduled automation is stored outside the repository in Codex's
  user-level automation store.
- Keep this file updated when the prompt, schedule, branch policy, or skill
  usage changes.
