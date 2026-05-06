# Datadog Rails Integration

The Rails application uses the `datadog` gem for automatic instrumentation and tracing.

## Overview

The integration is configured via `config/initializers/datadog.rb` and leverages the `datadog` gem to automatically instrument:

- **Rails**: ActionPack, ActiveRecord, ActionView, etc.
- **HTTP**: Outgoing HTTP requests via standard libraries.
- **Shoryuken**: Background job processing (if used).

## Configuration

Configuration is primarily handled through environment variables set in the `Dockerfile` or at runtime.

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `DD_AGENT_HOST` | The hostname of the Datadog agent. | `datadog-agent` |
| `DD_TRACE_ENABLED` | Whether tracing is enabled. | `true` |
| `DD_SERVICE` | The name of the service in Datadog. | `emmy-rails-app` |
| `DD_ENV` | The environment name (e.g., production, staging). | `production` |
| `DD_VERSION` | The version of the application. | `1.0.0` |
| `DD_PROFILING_ENABLED` | Whether continuous profiling is enabled. | `false` |

## Docker Integration

The `Dockerfile` includes these environment variables by default to ensure the application is ready for instrumentation when deployed.

```dockerfile
ENV DD_SERVICE="emmy-rails-app" \
    DD_TRACE_ENABLED="true" \
    DD_AGENT_HOST="datadog-agent" \
    DD_ENV="production" \
    DD_VERSION="1.0.0"
```

## Initialization

The initializer `rails-app/config/initializers/datadog.rb` sets up the instrumentation:

```ruby
require 'datadog'

Datadog.configure do |c|
  c.tracing.instrument :rails, service_name: ENV.fetch('DD_SERVICE', 'emmy-rails-app')
  c.tracing.instrument :http, service_name: ENV.fetch('DD_SERVICE', 'emmy-rails-app')
  c.tracing.instrument :shoryuken if defined?(Shoryuken)
end
```
