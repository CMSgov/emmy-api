# Emmy API Verification Service

The Emmy API is a backend data service designed to equip states with the automated verification capabilities needed to determine if an applicant has met H.R. 1 Community Engagement (CE) standards.

## Quick Links

- **[Technical Guide to Getting Started](docs/guides/01-getting-started.md)** (✅ Implementers start here!)
  - [Usage Examples](docs/guides/03-usage-examples.md)
  - [Emmy API Specifications (Swagger)](https://cmsgov.github.io/emmy-api/swagger-ui/)
- **[About Emmy Software](https://cms.gov/eligibility-made-easy)**
  - [Emmy API Stakeholder Overview](https://cms.gov/eligibility-made-easy) (➡️ Great for decision-makers & teams)
  - [Interested in the **Emmy App** instead?](https://github.com/DSACMS/iv-cbv-payroll)
- [Developer and Repo Information](#local-development)
  - [Emmy API Repo Homepage](https://github.com/DSACMS/iv-cbv-payroll)

## About the Project

This project evolved out of the [IVaaS](https://github.com/DSACMS/iv-cbv-payroll "IVaaS repository") tool for consent based verification. As the need for more complex forms of validation developed, it became clear that providing a way for agencies to integrate directly with an API was becoming increasingly useful, particularly for ex parte renewals. The ultimate goals of this and related projects is to remove as much friction as possible between the applicant and receiving their benefits by reducing the burden placed on them to manually provide evidence of eligibility.

The Emmy API serves as a "Decision Brain" for the state. It does not replace the state's eligibility system; instead, it feeds it the definitive "Verified/Not Verified" signals needed to clear the backlog.

## Core Team

Repository contact channels and observable contributor information are listed in
[COMMUNITY.md](COMMUNITY.md).

## Local Development

This project uses [pre-commit](https://pre-commit.com/ "pre-commit Docs") to register git hooks that run our various linters. You can [install](https://pre-commit.com/#install "Installation Instructions") pre-commit and then run:

```sh
pre-commit install
```

For OpenAPI spec maintenance, the repo also includes utility scripts under
`scripts/`:

```sh
# 1. Rebuild the bundled YAML and JSON artifacts from the design-time source spec
./scripts/bundle-api-spec

# 2. Validate the bundled YAML artifact
./scripts/validate-api-spec

# 3. Lint the bundled YAML artifact with the repo Spectral ruleset
./scripts/lint-api-spec

# 4. Lint hand-authored YAML files
./scripts/lint-yaml-files

# 5. Check formatting for contract YAML and JSON files
./scripts/check-format-contract-files

# 6. Compile standalone JSON Schemas
./scripts/check-json-schemas

# 7. Check for breaking OpenAPI changes against a base ref
./scripts/check-openapi-breaking [base-ref]
```

If you use [mise](https://mise.jdx.dev/), install the pinned runtimes from
`mise.toml` and run:

```sh
# Install the pinned Go and Node runtimes for this repo
mise install

# Run the full local contract/spec workflow
mise run check-contract-files
```

You can also run each step individually with `mise run bundle-api-spec`,
`mise run validate-api-spec`, `mise run lint-api-spec`,
`mise run lint-yaml-files`, `mise run check-format-contract-files`,
`mise run check-json-schemas`, and `mise run check-openapi-breaking`.

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
