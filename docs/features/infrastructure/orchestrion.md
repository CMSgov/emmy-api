# Orchestrion

The `emmy-api` service uses [Datadog Orchestrion](https://github.com/DataDog/orchestrion) for automatic instrumentation. This replaces manual OpenTelemetry instrumentation with build-time code injection.

## Overview

Orchestrion automatically instruments common Go libraries (like Fiber and Go-Redis) to provide distributed tracing and metrics without requiring manual code changes in handlers or clients.

## Configuration

The instrumentation is configured via environment variables, typically reporting to a Datadog agent or compatible collector on port `8126`.

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `DD_AGENT_HOST` | The hostname of the Datadog agent. | `localhost` |
| `DD_TRACE_AGENT_PORT` | The port the agent is listening on for traces. | `8126` |
| `DD_SERVICE` | The name of the service. | `emmy-api` |

## Build Integration

The `Dockerfile` is updated to use the `orchestrion` tool during the build process:

```dockerfile
RUN go install github.com/datadog/orchestrion@latest
RUN orchestrion build -o apiserver .
```

## Local Development

When running locally via `docker-compose`, a Datadog agent container is provided to receive and process traces.

```yaml
  datadog-agent:
    image: gcr.io/datadoghq/agent:latest
    ports:
      - "8126:8126"
```
