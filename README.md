# Verification Service API

An API that provides a unified, standard interface to the data sources needed to verify community engagement for the purpose of evaluating Medicaid eligibility.

## About the Project

This project evolved out of the [IVaaS](https://github.com/DSACMS/iv-cbv-payroll "IVaaS repository") tool for consent based verification. As the need for more complex forms of validation developed, it became clear that providing a way for agencies to integrate directly with an API was becoming increasingly useful, particularly for ex parte renewals. The ultimate goals of this and related projects is to remove as much friction as possible between the applicant and receiving their benefits by reducing the burden placed on them to manually provide evidence of eligibility.

Current repository-focused documentation starts with:

- [docs/overview.md](docs/overview.md) for runtime shape and system context
- [docs/setup.md](docs/setup.md) for local setup, testing, and tooling
- [api-spec/README.md](api-spec/README.md) for the contract-first API workflow

## Core Team

Repository contact channels and observable contributor information are listed in
[COMMUNITY.md](COMMUNITY.md).

## Local Development

This project uses [pre-commit](https://pre-commit.com/ "pre-commit Docs") to register git hooks that run our various linters. You can [install](https://pre-commit.com/#install "Installation Instructions") pre-commit and then run:

```sh
pre-commit install
```

For OpenAPI spec maintenance, use the commands that are currently wired on this
branch:

```sh
# 1. Rebuild the bundled JSON artifact from the design-time source spec
pnpm run bundle:api-spec

# 2. Validate the bundled JSON artifact
pnpm exec swagger-cli validate api-spec/v0/dist/openapi.bundled.json

# 3. Lint the bundled JSON artifact with the repo Spectral ruleset
pnpm exec spectral lint api-spec/v0/dist/openapi.bundled.json --ruleset .spectral.yaml

# 4. Check formatting for contract YAML and JSON files
pnpm exec prettier --check --no-error-on-unmatched-pattern "api-spec/**/*.yaml" "schema/**/*.json"

# 5. Compile standalone JSON Schemas
pnpm run check:schemas
```

If you use [mise](https://mise.jdx.dev/), install the pinned runtimes from
`mise.toml`, then run the individual tasks that are backed by checked-in
commands:

```sh
# Install the pinned Go and Node runtimes for this repo
mise install

# Run the available local contract/spec tasks
mise run bundle-api-spec
pnpm exec prettier --check --no-error-on-unmatched-pattern "api-spec/**/*.yaml" "schema/**/*.json"
pnpm run check:schemas
```

The aggregate `mise run check-contract-files` and `mise run lint` tasks are not
fully usable on this branch yet because `scripts/lint-yaml-files` and
`scripts/check-openapi-breaking` are not checked in.

## Policies

### Open Source Policy

We adhere to the [CMS Open Source
Policy](https://github.com/CMSGov/cms-open-source-policy). If you have any
questions, just [shoot us an email](mailto:opensource@cms.hhs.gov).

### Security and Responsible Disclosure Policy

_Submit a vulnerability:_ Vulnerability reports can be submitted through [Bugcrowd](https://bugcrowd.com/cms-vdp). Reports may be submitted anonymously. If you share contact information, we will acknowledge receipt of your report within 3 business days.

For more information about our Security, Vulnerability, and Responsible Disclosure Policies, see [SECURITY.md](SECURITY.md).

### Software Bill of Materials (SBOM)

A Software Bill of Materials (SBOM) is a formal record containing the details and supply chain relationships of various components used in building software.

In the spirit of [Executive Order 14028 - Improving the Nation’s Cyber Security](https://www.gsa.gov/technology/it-contract-vehicles-and-purchasing-programs/information-technology-category/it-security/executive-order-14028), the current dependency graph for this repository is available at:
[https://github.com/CMSgov/emmy-api/network/dependencies](https://github.com/CMSgov/emmy-api/network/dependencies)

For more information and resources about SBOMs, visit: [https://www.cisa.gov/sbom](https://www.cisa.gov/sbom).

## Public domain

This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/) as indicated in [LICENSE](LICENSE).

All contributions to this project will be released under the CC0 dedication. By submitting a pull request or issue, you are agreeing to comply with this waiver of copyright interest.
