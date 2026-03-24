# Feature: Jaeger

## Feature Overview
Provides local trace sink and UI for distributed tracing data emitted through the OpenTelemetry collector.

## Business Logic
- Application emits trace data through OTel pipeline.
- OTel collector receives OTLP traces and exports them to Jaeger.
- Docker compose runs Jaeger all-in-one and exposes UI for trace inspection.

## Package Location
- `otel-collector-config.yml`
- `docker-compose.yml`
- `pkg/core/otel.go`

## Key Structs and Interfaces
- Collector trace pipeline: `service.pipelines.traces`
- Collector Jaeger exporter: `exporters.jaeger`
- Compose service: `jaeger-all-in-one`

## Real Code Excerpt
```yaml
# otel-collector-config.yml
exporters:
  jaeger:
    endpoint: jaeger-all-in-one:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, jaeger]

# docker-compose.yml
jaeger-all-in-one:
  image: jaegertracing/all-in-one:latest
  ports:
    - "16686:16686"
    - "14268"
    - "14250"
```

## Edge Cases Handled Today
- Collector retains logging exporter in trace pipeline for local debugging visibility.
- Jaeger unavailability affects trace sink visibility but is not designed as a request-serving gate.
- Local TLS mode is explicitly insecure for internal compose-network transport.

## Performance and Operational Considerations
- Jaeger all-in-one mode is suitable for local/dev workflows, not production scale.
- Trace visibility depends on both collector and Jaeger service health.
- UI is exposed on `http://localhost:16686` in the compose stack.

## Future Improvements
- Move to production-grade distributed Jaeger/OpenTelemetry backend strategy.
- Add retention/storage policy documentation for trace data.
- Add trace quality checks (sampling strategy and span coverage expectations).

## Assumptions
- **High confidence:** Jaeger is currently used as a local debugging/inspection sink for traces.
- **Medium confidence:** Long-term production tracing backend may differ from local all-in-one topology.
