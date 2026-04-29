# .github AGENTS

## Scope

These instructions apply to work inside `.github/`, especially
`.github/workflows/`.

## Workflow Runner Requirement

- Any workflow job in `.github/workflows/*.yml` that declares `runs-on` must
  use the exact label `codebuild-emmy-github-runner-emmy-api-${{
  github.run_id }}-${{ github.run_attempt }}`
- Do not replace that label with GitHub-hosted runners such as `ubuntu-latest`
  or any other runner name
- Reusable workflow-call jobs that use top-level `uses:` do not accept
  `runs-on`; leave those as valid reusable-workflow invocations unless the task
  explicitly requires restructuring them
- If you add a new direct-execution job or convert a reusable-workflow call
  into one, set:

```yaml
runs-on: "codebuild-emmy-github-runner-emmy-api-${{ github.run_id }}-${{ github.run_attempt }}"
```

## Safety

- Treat workflow runner selection as repository policy, not a convenience
- Escalate before changing workflow semantics beyond the requested task
