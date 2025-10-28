# Telemetry

OpenTelemetry (optional)
- Build with `-tags=otel` and set:
  - `OTEL_ENABLED=1`
  - `OTEL_SERVICE_NAME=ferret`
  - `OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318`

Metrics
- In-memory metrics under `pkg/metrics` for lightweight counters/histograms.
- OTel facade in `pkg/telemetry` records spans/counters/histograms when enabled.

Tracing
- Scheduler: `scheduler.claim` spans around DB claims.
- Poster: `poster.publish` spans around publish operations.

