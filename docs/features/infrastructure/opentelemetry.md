# Feature: OpenTelemetry

## Feature Overview
Provides application-level telemetry for traces and metrics with optional log fanout integration.

## Business Logic
- Initialize OTel service from config via `core.NewOtelService`.
- Use noop providers when `OTEL_DISABLE=true`.
- Build OTLP trace and metric exporters over gRPC when enabled.
- Attach Fiber tracing middleware for request spans.
- Optionally fan out structured logs through OTel logger provider.

## Package Location
- `pkg/core/otel.go`
- `pkg/core/logger.go`
- `api/app.go`
- `otel-collector-config.yml`
- `docker-compose.yml`

## Key Structs and Interfaces
- `OtelService`
- `otelService` (concrete)
- `NewOtelService(ctx, cfg)`
- `NewLoggerWithOtel(cfg, otel)`

## Real Code Excerpt
```go
func NewOtelService(ctx context.Context, cfg *Config) (OtelService, error) {
    if cfg.Otel.Disable {
        return newOtelServiceNoop(), nil
    }

    return newOtelServiceGRPC(ctx, cfg)
}

app.Use(otelfiber.Middleware())
```

## Edge Cases Handled Today
- Disabled mode uses noop providers and avoids exporter setup.
- Shutdown errors are logged and do not panic process.
- Exporter/provider creation failures return errors to startup path.

## Performance and Operational Considerations
- Trace provider uses batch span processing.
- Meter provider uses periodic reader export.
- Collector endpoint mismatch caveat: app config reads `OTEL_OTLP_EXPORTER_ENDPOINT`, while compose currently sets `OTEL_EXPORTER_OTLP_ENDPOINT`.
- Telemetry pipeline depends on OTel collector availability for sink delivery.

## Future Improvements
- Align OTLP env key usage between app config loader and compose env.
- Add explicit telemetry health indicators for operator diagnostics.
- Introduce sampling controls for high-throughput environments.

## Assumptions
- **High confidence:** OTel is expected in local/dev workflows through the collector stack.
- **High confidence:** `OTEL_DISABLE=true` is the supported fallback when collector/export path is unavailable.
