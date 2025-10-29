package telemetry

import (
    "context"
    "log"
    "os"
    "strconv"
    "time"
)

// Lightweight, optional telemetry facade. By default, all functions are no-ops
// to keep builds dependency-free. Build with -tags=otel to enable OpenTelemetry
// via pkg/telemetry/otel_enabled.go.

// InitFromEnv initializes telemetry (noop by default). With the 'otel' build
// tag, it configures OTel using standard env vars:
//   - OTEL_SERVICE_NAME (default: ferret)
//   - OTEL_EXPORTER_OTLP_ENDPOINT (e.g., http://localhost:4318)
//   - OTEL_EXPORTER_PROTOCOL (http/protobuf)
//   - OTEL_ENABLED=1 to toggle on
func InitFromEnv(ctx context.Context) context.Context {
    if os.Getenv("OTEL_ENABLED") == "1" {
        // In non-otel builds this just logs once as a hint.
        log.Printf("telemetry: OTEL_ENABLED=1 but 'otel' build tag not enabled; running with noop")
    }
    return ctx
}

// StartSpan starts a span and returns a context and end function. No-op by default.
func StartSpan(ctx context.Context, name string, attrs map[string]string) (context.Context, func()) {
    start := time.Now()
    return ctx, func() {
        _ = start // keep for symmetry; noop in default build
    }
}

// RecordCounter emits a counter increment. No-op by default.
func RecordCounter(ctx context.Context, name string, delta float64, attrs map[string]string) {}

// RecordHistogram emits a histogram sample. No-op by default.
func RecordHistogram(ctx context.Context, name string, value float64, attrs map[string]string) {}

// Helpers
func boolEnv(key string) bool { v := os.Getenv(key); b, _ := strconv.ParseBool(v); return b }

