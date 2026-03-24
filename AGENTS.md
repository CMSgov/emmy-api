# AGENTS

## 1. Purpose

Agents in this repository exist to make scoped, safe changes to the verification
service and its supporting artifacts.

This repo currently contains:

- Go application code for a Fiber HTTP service
- Redis-backed status and circuit-breaker behavior
- NSC education integration code
- OpenTelemetry, Prometheus, and Jaeger local observability config
- OpenAPI contract files and repository documentation

When repo prose and current implementation disagree, treat code, config, and CI
as the source of truth.

## 2. Repository Overview

Core structure:

- `main.go`: process bootstrap, config load, Redis init, app startup
- `api/`: HTTP app construction, routes, handlers, middleware
- `pkg/`: core config/logging/otel plus integration packages
- `api-spec/`: OpenAPI source files and bundled artifacts
- `docs/`: setup, architecture, API, feature, research, and audit docs
- `.github/workflows/`: CI checks for tests, linting, markdown, spelling, and secrets
- `Dockerfile`, `docker-compose.yml`, `otel-collector-config.yml`, `prometheus.yml`: local container and observability setup

API boundaries:

- `api/` owns HTTP routing, middleware, and handler wiring
- `pkg/` owns service logic, config, Redis, circuit breaker, and education integration
- `api-spec/` owns public contract artifacts

Observed entry points:

- `main.go`
- `api.New`
- `routes.RegisterRoutes`
- `api-spec/openapi.yaml`

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

- Do not claim or modify production deployment behavior based on this repo alone
- The observable deployment-related assets here are local Docker/Compose and CI build/test workflows

Approval required before changing:

- Auth behavior
- Secret-scanning rules
- Workflow files
- Policy files
- Bundled OpenAPI outputs
- `public.jwk`
