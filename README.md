# Eligibility Made Easy (Emmy) API Verification Service

The Emmy API is a backend data service that connects to federal and commercial data sources to facilitate eligibility determination for state agencies.

## Quick Links

- **[About Emmy Software](https://cms.gov/eligibility-made-easy)**
  - [Emmy API Overview](https://cms.gov/eligibility-made-easy)
  - [Emmy Application Github](https://github.com/DSACMS/iv-cbv-payroll/blob/main/README.md)

- **[Technical Guide to Getting Started](docs/guides/01-getting-started.md)**
  - [Usage Examples](docs/guides/03-usage-examples.md)
  - [Emmy API Specifications (Swagger)](https://cmsgov.github.io/emmy-api/swagger-ui/)

- [Developer and Repo Information](#local-development)

## About the Project

The Emmy API is a part of a suite of open-source tools developed by the Centers for Medicare & Medicaid Services (CMS). Emmy supports states in implementing the new Medicaid eligibility and renewal community engagement requirements under the Working Families Tax Cut legislation.

Emmy's suite of tools does not replace a state's existing eligibility system — rather, it offers streamlined, efficient ways to help states implement the new community engagement requirements.

## Core Team

Repository contact channels and observable contributor information are listed in
[COMMUNITY.md](COMMUNITY.md).

## Local Development

This project uses [pre-commit](https://pre-commit.com/ "pre-commit Docs") to register git hooks that run our various linters. You can [install](https://pre-commit.com/#install "Installation Instructions") pre-commit and then run:

```sh
pre-commit install
```

For OpenAPI spec maintenance, use the supported local commands from the project
root:

```sh
# 1. Rebuild the bundled artifacts from the design-time source spec
pnpm run bundle:api-spec

# 2. Validate the bundled artifact
pnpm run validate:api-spec

# 3. Lint the bundled artifact with the repo Spectral ruleset
pnpm run lint:api-spec

# 4. Check formatting for contract YAML and JSON files
./scripts/check-format-contract-files

# 5. Compile standalone JSON Schemas
./scripts/check-json-schemas
```

If you use [mise](https://mise.jdx.dev/), install the pinned runtimes from
`mise.toml` and run:

```sh
# Install the pinned Go and Node runtimes for this repo
mise install
```

You can then run the wired `mise` tasks individually with
`mise run bundle-api-spec`, `mise run validate-api-spec`,
`mise run lint-api-spec`, `mise run check-format-contract-files`, and
`mise run check-json-schemas`.

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

In the spirit of [Executive Order 14028 - Improving the Nation's Cyber Security](https://www.gsa.gov/technology/it-contract-vehicles-and-purchasing-programs/information-technology-category/it-security/executive-order-14028), the current dependency graph for this repository is available at:
[https://github.com/CMSgov/emmy-api/network/dependencies](https://github.com/CMSgov/emmy-api/network/dependencies)

For more information and resources about SBOMs, visit: [https://www.cisa.gov/sbom](https://www.cisa.gov/sbom).

## Public domain

This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/) as indicated in [LICENSE](LICENSE).

All contributions to this project will be released under the CC0 dedication. By submitting a pull request or issue, you are agreeing to comply with this waiver of copyright interest.
