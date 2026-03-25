# Postman Collections

**NOTE: EXAMPLE AI GENERATED README - NEEDS TO CHANGE** This directory contains Postman collections and environment configurations for
exploring, testing, and integrating with the **My API** (v1).

These collections complement the [Getting Started Guide](../../docs/guides/getting-started.md)
and the [OpenAPI specification](../../openapi.yaml) at the root of this repository.

---

## 📁 Directory Structure

```bash
/postman
  collections/
    my-api.postman_collection.json       ← Main API collection
    my-api-auth-examples.postman_collection.json  ← Auth flow examples
  environments/
    local.postman_environment.json       ← Local development environment
    staging.postman_environment.json     ← Staging environment
  README.md                              ← You are here
```

---

## 🚀 Getting Started

### Prerequisites

- [Postman](https://www.postman.com/downloads/) (Desktop or Web) — v10 or later recommended
- A running instance of the API, or access to the staging environment
- Valid API credentials (see [Authentication Guide](../../docs/guides/authentication.md))

---

## 📥 Importing the Collection & Environment

1. Open **Postman**
2. Click **Import** (top-left)
3. Select or drag-and-drop the desired files from this directory:
   - `collections/my-api.postman_collection.json`
   - `environments/local.postman_environment.json` *(or `staging` as appropriate)*
4. In the top-right environment dropdown, select **My API – Local** (or **Staging**)
5. You're ready to make requests! 🎉

---

## 🔧 Environment Variables

Each environment file uses the following variables. **Never commit real secrets** —
populate sensitive values in your local Postman environment after importing.

| Variable | Description | Example Value |
|---|---|---|
| `base_url` | Base URL of the API | `http://localhost:8080` |
| `api_version` | API version prefix | `v1` |
| `api_key` | Your API key *(set locally — do not commit)* | `your-api-key-here` |
| `auth_token` | Bearer token after authentication | *(populated by auth request)* |
| `user_id` | Sample user ID for testing | `12345` |

> ⚠️ **Security Note:** The environment files in this repository use **placeholder
> values only**. Never replace placeholders with real credentials and commit them.
> Use Postman's built-in secret masking or a `.env`-style local override for
> sensitive values.

---

## 📬 Collection Overview

### `my-api.postman_collection.json`

The primary collection covering all major API endpoints, organized into folders:

| Folder | Description |
|---|---|
| **Auth** | Obtain and refresh authentication tokens |
| **Users** | Create, read, update, and delete user resources |
| **Items** | Core resource CRUD operations |
| **Search** | Query and filter endpoints |
| **Webhooks** | Register and manage webhook subscriptions |
| **Health** | Health check and status endpoints |

### `my-api-auth-examples.postman_collection.json`

Focused collection demonstrating various authentication flows:

- API Key authentication
- OAuth 2.0 Client Credentials
- Bearer token usage

---

## ▶️ Running the Collection (Automated / CI)

You can run the full collection headlessly using the
[Postman CLI (`newman`)](https://github.com/postmanlabs/newman):

```bash
# Install newman
npm install -g newman

# Run the collection against the local environment
newman run collections/my-api.postman_collection.json \
  --environment environments/local.postman_environment.json \
  --reporters cli,json \
  --reporter-json-export results/newman-report.json
```

This is also integrated into the CI/CD pipeline — see
[`.github/workflows/integration-tests.yml`](../../.github/workflows/integration-tests.yml)
for details.

---

## 🧪 Testing Notes

- Requests marked with ⚡ include **automated test scripts** (visible in the
  Postman **Tests** tab) that validate response status codes, schema, and key fields.
- Run the full collection via the **Collection Runner** in Postman to see a
  pass/fail summary across all requests.
- For load or performance testing, consider exporting to
  [k6](https://k6.io/) using the
  [postman-to-k6](https://github.com/apideck-libraries/postman-to-k6) converter.

---

## 🔄 Keeping Collections Up to Date

These collections are maintained alongside the codebase. When API changes are made:

1. Update the relevant requests/folders in the collection
2. Re-export the collection JSON from Postman (**⋯ → Export → Collection v2.1**)
3. Replace the existing `.json` file and commit the changes with your PR

If you notice a collection is out of sync with the current API, please
[open an issue](../../issues) or submit a PR.

---

## 📚 Related Resources

- [Getting Started Guide](../../docs/guides/getting-started.md)
- [Authentication Guide](../../docs/guides/authentication.md)
- [OpenAPI Specification](../../openapi.yaml)
- [Swagger UI](https://api.example.gov/docs) *(live reference)*
- [Contributing Guidelines](../../CONTRIBUTING.md)

---

## 🙋 Questions or Issues?

If you run into problems with the collections or environments, please:

- Check the [Troubleshooting Guide](../../docs/guides/error-handling.md)
- Open an [issue](../../issues/new) in this repository
- Reach out via [support@example.gov](mailto:support@example.gov)
