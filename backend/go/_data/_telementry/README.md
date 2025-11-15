Telemetry & Tracing

OpenTelemetry (optional)
- Enable with build tag `otel` and env:
  - OTEL_ENABLED=1
- OTEL_SERVICE_NAME=social-scale
  - OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318

Spans
- scheduler.claim — wraps DB claim of due posts
- poster.publish — wraps publish call per platform

Metrics
- See ../_metrics/metrics.yaml for list and attributes

Fallback
- app_metrics table (_data/_db/006_telemetry.sql) for local storage

