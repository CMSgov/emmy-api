# Static Swagger Site

`swagger-site/` contains the static Swagger UI payload uploaded to S3 website
hosting.

Tracked source files:

- `index.html`
- `swagger-initializer.js`

Generated files are produced by `./scripts/build-swagger-site`:

- `openapi.json` (copied from `api-spec/v0/dist/openapi.bundled.json`)
- Swagger UI runtime assets copied from `swagger-ui/`
