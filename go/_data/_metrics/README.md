Metrics Plan

Goals
- Uniform naming and types for counters, gauges, histograms
- Map to OpenTelemetry where possible
- Minimal local fallback via app_metrics table

Naming
- Use dot-separated names, all lowercase
- Include subsystem prefix: scheduler., poster., analytics., cache.

Core Metrics
- scheduler.claimed_total (counter) — number of posts claimed
- poster.posted_total (counter){platform}
- poster.failed_total (counter){platform}
- poster.post_seconds (histogram){platform}
- cache.valkey.ops_total (counter){op}

OTel Mapping
- Counters -> OTel counter (int)
- Histograms -> OTel histogram (float)
- Attributes map to OTel attributes (string)

Env
- OTEL_ENABLED=1
- OTEL_SERVICE_NAME=ferret
- OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318

