# Feature: OpenTelemetry Logging

## Feature Overview
Implements trace/metric export and structured logging with optional OpenTelemetry log fanout.

## Business Logic
- Initialize OTel service from config (`disabled` or `grpc exporter` path).
- Register trace and meter providers with OTLP exporters.
- Attach Fiber OTel middleware for request tracing.
- Create structured slog logger and optionally fan out to OTel log handler.
- Instrument Redis client for tracing/metrics.

## Package Location
- `pkg/core/otel.go`
- `pkg/core/logger.go`
- `api/app.go`
- `pkg/redis/redis.go`
- `otel-collector-config.yml`
- `prometheus.yml`

## Key Structs and Interfaces
- `OtelService`
- `otelService` (concrete)
- `NewLogger`
- `NewLoggerWithOtel`

## Real Code Excerpt
```go
app.Use(otelfiber.Middleware())

app.Use(slogfiber.NewWithConfig(
    cfg.Logger,
    slogfiber.Config{
        WithRequestID: true,
        WithSpanID:    true,
        WithTraceID:   true,
    },
))
```

## Edge Cases Handled Today
- OTel disabled mode uses noop providers.
- OTel shutdown errors are logged but do not panic process.
- Redis OTel instrumentation failures degrade to warning logs.

## Performance and Operational Considerations
- Trace provider uses batch span processor.
- Meter provider uses periodic reader export.
- Logging output format changes by environment (`JSON` in production, text otherwise).

## Future Improvements
- Add log sampling strategy for high-throughput scenarios.
- Add custom business metrics around NSC calls and breaker state.
- Align environment key naming for OTLP endpoint consistency across compose and app config (compose currently uses `OTEL_EXPORTER_OTLP_ENDPOINT`, while app config loader reads `OTEL_OTLP_EXPORTER_ENDPOINT`).

## Assumptions
- **High confidence:** Current stack is sufficient for local observability validation but needs production hardening around metrics and SLO-driven instrumentation.
