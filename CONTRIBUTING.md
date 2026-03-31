# How to Contribute

## Getting Started

Start with [docs/setup.md](docs/setup.md), which is the repository's detailed
setup guide. The current branch expects:

- Go `1.25.x`
- Node `20` for contract and docs tooling
- Local Redis on `localhost:6379` for several tests and health-check behavior
- Docker for the local compose workflow

If you use `mise`, install the pinned toolchain with:

```sh
mise install
```

To enable the repo's pre-commit hooks:

```sh
pre-commit install
```

## Local Development Workflow

Run the service directly:

```sh
go run .
```

Run the local container stack:

```sh
docker compose up --build
```

Run the Go test suite:

```sh
go test ./...
```

Run the routine lint and scan tasks with `mise`:

```sh
mise run lint
```

For contract-focused changes, run the narrower contract workflow:

```sh
mise run check-contract-files
```

That task bundles, validates, lints, and formatting-checks the OpenAPI and JSON
Schema artifacts that live under `api-spec/` and `schema/`.

## Pull Requests and Issues

- Keep pull requests scoped to one logical change when possible.
- If you change source spec files under `api-spec/v0/`, include the generated
  `api-spec/v0/dist/` and `api-spec/dist/v0/` updates in the same change.
- Run the relevant local checks before opening a pull request.
- Use GitHub issues and pull requests in this repository for bug reports,
  feature discussions, and proposed changes.
- Use synthetic data only in docs, tests, fixtures, and examples.

## Policies

### Open Source Policy

We adhere to the [CMS Open Source
Policy](https://github.com/CMSGov/cms-open-source-policy). If you have any
questions, just [shoot us an email](mailto:opensource@cms.hhs.gov).

### Security and Responsible Disclosure Policy

_Submit a vulnerability:_ Vulnerability reports can be submitted through [Bugcrowd](https://bugcrowd.com/cms-vdp). Reports may be submitted anonymously. If you share contact information, we will acknowledge receipt of your report within 3 business days.

For more information about our Security, Vulnerability, and Responsible Disclosure Policies, see [SECURITY.md](SECURITY.md).

## Public domain

This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/) as indicated in [LICENSE](LICENSE).

All contributions to this project will be released under the CC0 dedication. By submitting a pull request or issue, you are agreeing to comply with this waiver of copyright interest.
