---
name: generate-rails-openapi
description: Generate or update the OpenAPI specification from the Rails application tests (specs). Use when you need to synchronize the public API contract with the actual implementation in the Rails service.
---

# Generate Rails OpenAPI

Use this skill to automatically generate an OpenAPI specification by running the RSpec integration tests. This ensures that the documentation accurately reflects the behavior of the API.

## Required Inputs

- Rails application directory (`rails-app`)
- Integration tests (request specs) in `spec/requests/`
- `rswag` gems (`rswag-api`, `rswag-ui`, `rswag-specs`) installed in the Rails app

## Workflow

1. **Ensure dependencies are installed**:

   ```bash
   cd rails-app
   bundle install
   ```

2. **Run tests with OpenAPI generation**:
   Execute the RSpec tests using the Rswag task. This will run the specs and generate the Swagger file based on the `path`, `post`, `response`, etc. blocks.

   ```bash
   cd rails-app
   # Generate Swagger documentation
   bundle exec rake rswag:specs:swaggerize
   ```

3. **Locate the generated spec**:
   By default, the spec is generated at `swagger/v1/swagger.yaml` (as configured in `spec/swagger_helper.rb`).

4. **Review**:
   Review the generated `swagger.yaml` to ensure it correctly captures the API behavior.

## Repository Configuration

- **Gem**: `rswag` is used to capture request/response data and generate the UI.
- **Helper**: `spec/swagger_helper.rb` configures the output path and global metadata (title, version, servers).
- **Test Inclusion**: Ensure request specs use `type: :request` and follow the Rswag DSL.

  ```ruby
  # spec/requests/api/v0/my_spec.rb

  require 'swagger_helper'

  RSpec.describe 'api/v0/my_resource', type: :request do
    path '/api/v0/my_resource' do
      post 'Creates a resource' do
        tags 'MyResource'
        consumes 'application/json'
        parameter name: :resource, in: :body, schema: { ... }

        response '201', 'resource created' do
          run_test!
        end
      end
    end
  end

  ```

## Troubleshooting

- **No examples found**: Ensure you are running `rake rswag:specs:swaggerize` and that your specs use the Rswag DSL (`path`, `response`, etc.).
- **Database Errors**: Ensure a test database is available or mocked. The app may require `AWS_REGION` and `DB_IAM_AUTH=false` for local tests.
- **Schema Drift**: We use a shared example `swagger_schema_drift_detection` to ensure models and their `to_swagger_schema` definitions stay in sync. Run these via:

  ```bash
  bundle exec rspec spec/models
  ```
