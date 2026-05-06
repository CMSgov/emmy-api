# Publish Swagger Docs to S3

This repository currently uses a contract-first OpenAPI workflow.

- Source contract: `api-spec/v0/openapi.yaml`
- Runtime-served bundled artifact: `api-spec/v0/dist/openapi.bundled.json`
- App route that serves this artifact: `GET /api-spec/v1/verify`

Swaggo (`swag init`) is not currently part of the active generation flow in this
branch; the generated static site builds from the checked-in OpenAPI source and
bundle command already used in this repository.

## Generate Static Swagger Site Locally

From the repository root:

```bash
pnpm install --frozen-lockfile
./scripts/build-swagger-site
```

This produces a deployable static payload under `swagger-site/`:

- `swagger-site/index.html`
- `swagger-site/swagger-initializer.js`
- `swagger-site/openapi.json`
- Swagger UI assets copied from `swagger-ui/`

`swagger-site/swagger-initializer.js` is configured with:

```js
url: "./openapi.json"
```

So Swagger UI loads the OpenAPI artifact from the same S3 website bucket.

## Manual Publish to S3

Set environment variables and run:

```bash
AWS_REGION=us-east-1 \
SWAGGER_S3_BUCKET=my-swagger-site-bucket \
./scripts/publish-swagger-site
```

Equivalent direct sync command:

```bash
aws s3 sync swagger-site "s3://my-swagger-site-bucket" --delete
```

## CI/CD Workflow

The workflow at `.github/workflows/publish-swagger-site.yml` runs on:

- Push to `main`
- Manual `workflow_dispatch`

It performs:

1. Setup Go and download module dependencies
2. Setup Node/pnpm and install dependencies
3. Run `./scripts/build-swagger-site`
4. Assume AWS credentials via OIDC (`aws-actions/configure-aws-credentials`)
5. Run `aws s3 sync swagger-site s3://<bucket> --delete`

## Required GitHub Actions Configuration

Set these in GitHub repository **Variables** (preferred) or **Secrets**:

- `SWAGGER_S3_BUCKET`: destination S3 bucket name
- `AWS_REGION`: AWS region for bucket/CLI operations (defaults to `us-east-1` if
  not set)

For AWS auth (same style as existing deploy workflow), provide either:

- `AWS_ROLE_TO_ASSUME`

or both:

- `AWS_ACCOUNT_ID`
- `AWS_ROLE_NAME`

## Required S3 Website Settings

For direct public S3 website hosting:

1. Enable **Static website hosting**
2. Set **Index document** to `index.html`
3. If the site must be public, allow public policy access (disable blocking
   public access settings that prevent policy-based reads)
4. Apply a bucket policy allowing object reads

Example policy (replace `YOUR_BUCKET_NAME`):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowPublicReadForSwaggerSite",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::YOUR_BUCKET_NAME/*"
    }
  ]
}
```

## Expected Website Endpoint Format

S3 website endpoint format:

```text
http://<bucket-name>.s3-website-<region>.amazonaws.com
```

Example (`us-east-1`):

```text
http://my-swagger-site-bucket.s3-website-us-east-1.amazonaws.com
```
