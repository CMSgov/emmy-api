# Feature: Prometheus

## Feature Overview
Provides metrics scraping and visualization backend integration for collector-exported telemetry.

## Business Logic
- OTel collector exports metrics via Prometheus exporter on `:8889`.
- Prometheus scrapes collector targets defined in `prometheus.yml`.
- Docker compose runs Prometheus service and exposes UI on `:9090`.

## Package Location
- `prometheus.yml`
- `otel-collector-config.yml`
- `docker-compose.yml`

## Key Structs and Interfaces
- Collector exporter: `exporters.prometheus.endpoint`
- Collector metrics pipeline: `service.pipelines.metrics`
- Prometheus scrape config: `scrape_configs`

## Real Code Excerpt
```yaml
# otel-collector-config.yml
exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, prometheus]

# prometheus.yml
scrape_configs:
  - job_name: "otel-collector"
    scrape_interval: 10s
    static_configs:
      - targets: ["otel-collector:8889"]
      - targets: ["otel-collector:8888"]
```

## Edge Cases Handled Today
- Collector-side logging exporter remains available even if Prometheus is not consuming metrics.
- Scrape targets are static and simple for local stack predictability.
- Prometheus failure does not directly block API request handling.

## Performance and Operational Considerations
- Current metrics are mostly framework/instrumentation level; business-specific metrics are limited.
- End-to-end metrics visibility depends on OTel collector availability.
- Scrape interval is fixed at `10s` in current config.

## Future Improvements
- Add application business metrics (NSC, breaker, auth outcomes).
- Add alerts and recording rules for latency/error SLO tracking.
- Expand scrape config for multi-service production deployment patterns.

## Assumptions
- **High confidence:** Prometheus is currently intended for local observability through collector metrics endpoints.
- **Medium confidence:** Production environments may use managed metrics backends with different scrape/discovery models.
