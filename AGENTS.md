# AGENTS

## 1. Purpose

Agents in this repository exist to make scoped, safe changes to the verification
service and its supporting artifacts.

This repo currently contains:

- Go application code for a Fiber HTTP service
- Redis-backed status and circuit-breaker behavior
- NSC education integration code
- Datadog/Orchestrion local observability config
- ECS deployment and image-publish helper scripts
- OpenAPI contract files and repository documentation

When repo prose and current implementation disagree, treat code, config, and CI
as the source of truth.

## 2. Repository Overview

Core structure:

- `main.go`: process bootstrap, config load, Redis init, app startup
- `api/`: HTTP app construction, routes, handlers, middleware
- `pkg/`: core config/logging plus integration packages
- `api-spec/`: OpenAPI source files and bundled artifacts
- `docs/`: setup, architecture, API, feature, research, and audit docs
- `.github/workflows/`: CI checks for tests, linting, markdown, spelling, and secrets
- `scripts/`: helper scripts for image build/push and ECS deployment
- `Dockerfile`, `docker-compose.yml`: local container and observability setup

API boundaries:

- `api/` owns HTTP routing, middleware, and handler wiring
- `pkg/` owns service logic, config, Redis, circuit breaker, and education integration
- `api-spec/` owns public contract artifacts

Observed entry points:

- `main.go`
- `api.New`
- `routes.RegisterRoutes`
- `api-spec/openapi.yaml`

Observed deployment helpers:

- `scripts/push-image`: builds and pushes images to the configured private registry
- `scripts/deploy-ecs`: updates ECS task definitions and services for named environments
- `scripts/emmy-common.sh`: defines the supported deploy environments: `dev`, `test`, `demo`, `uat`, `sandbox`, `prod`

Do not modify these without an explicit task and approval:

- `.github/workflows/*`
- `.github/.gitleaks.toml`
- `SECURITY.md`
- `LICENSE`
- `public.jwk`
- `api-spec/dist/*` unless the source spec changed too

## 3. Agent Roles

### Runtime Agent

- Allowed: change Go code under `api/`, `pkg/`, and `main.go`
- Forbidden: silently change auth behavior, secrets handling, workflow files, or policy files
- Escalate when: a change crosses into public API contract, security behavior, CI, or repo policy
- Decision authority: within runtime code boundaries only

### Contract & Docs Agent

- Allowed: change `docs/`, `README.md`, and `api-spec/`
- Forbidden: describe behavior that is not observable in the current branch
- Escalate when: changing the public contract, auth semantics, or bundled spec outputs
- Decision authority: within docs and contract boundaries only

### Security/Workflow Agent

- Allowed: review security, auth, secret-scanning, and CI/workflow files; edit them only when explicitly asked
- Forbidden: make unapproved changes to auth, secrets policy, workflow policy, or sensitive root files
- Escalate when: any requested change touches security posture, scanning rules, or workflow enforcement
- When editing `.github/workflows/publish-image.yml`, use the GitHub-hosted `ubuntu-latest` runner for its direct jobs
- For every other workflow job that declares `runs-on`, retain the exact runner label `runs-on: "codebuild-emmy-github-runner-emmy-api-${{ github.run_id }}-${{ github.run_attempt }}"` unless a separate runner-migration task explicitly changes that workflow
- Reusable workflow-call jobs that use top-level `uses:` cannot declare `runs-on`; if you convert one into a direct job or add a new direct job, add the required runner label
- Decision authority: review by default; edit only with explicit approval

Agents may decide within their boundary, but must escalate before crossing into
security, workflow, or public API contract changes.

## 4. Coding Standards

Formatting and linting are defined by repo config:

- `.editorconfig` controls indentation, newline, and whitespace rules
- `.golangci.yml` defines Go linters and formatters
- `.markdownlint.yml` defines markdown rules
- `.pre-commit-config.yaml` wires the main local hooks

Observed commit guidance:

- No hard-enforced commit convention is declared in repo config
- Recent history commonly uses short prefixes such as `docs:`, `feat:`, and `lint:`
- Follow that lightweight style instead of inventing a new convention

Testing expectations by touched area:

- Go changes: `go test ./...`
- Go runtime changes: `golangci-lint`
- Markdown or docs changes: markdownlint or `pre-commit`
- Container or runtime config changes: `docker build .`

Known caveat:

- Some tests require local Redis on `localhost:6379`

## 5. Safety & Constraints

Secrets handling:

- Never commit real secrets, access tokens, or private keys
- Respect `.github/.gitleaks.toml`
- Use environment variables for sensitive values instead of committing them to files

Data privacy:

- Request models in `pkg/education` include names, date of birth, and SSN fields
- Use synthetic example data only in docs, tests, and fixtures

Deployment restrictions:

- Do not claim full cloud topology, account layout, or approval policy unless it is observable in the current branch
- The repo does include deployment automation artifacts: ECS deployment helpers under `rails-app/script/`, a required ECR repository input in the image helpers, and CI workflow definitions under `.github/workflows/`
- Treat environment names and deploy script behavior as source-of-truth only for what the scripts actually encode

Approval required before changing:

- Auth behavior
- Secret-scanning rules
- Workflow files
- Deployment scripts or registry/deploy environment conventions
- Policy files
- Bundled OpenAPI outputs
- `public.jwk`
