//go:build otel
// +build otel


import (
    "context"
    "log"
    "os"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/metric"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
    meter  metric.Meter
    ctrs   = map[string]metric.Int64Counter{}
    histos = map[string]metric.Float64Histogram{}
)

// InitFromEnv bootstraps OTel SDK if OTEL_ENABLED is true-ish.
func InitFromEnv(ctx context.Context) context.Context {
    if !boolEnv("OTEL_ENABLED") { return ctx }
    svc := os.Getenv("OTEL_SERVICE_NAME")
    if svc == "" { svc = "ferret" }
    endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
    if endpoint == "" { endpoint = "http://localhost:4318" }

    // Traces
    texp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
    if err != nil { log.Printf("otel: trace exporter: %v", err); return ctx }
    res, _ := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String(svc),
        ),
    )
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(texp),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)
    // Metrics
    mexp, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL(endpoint))
    if err != nil { log.Printf("otel: metric exporter: %v", err); return ctx }
    mp := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(sdkmetric.NewPeriodicReader(mexp, sdkmetric.WithInterval(10*time.Second))),
        sdkmetric.WithResource(res),
    )
    otel.SetMeterProvider(mp)
    meter = mp.Meter(svc)
    return ctx
}

// StartSpan starts a real OTel span.
func StartSpan(ctx context.Context, name string, attrs map[string]string) (context.Context, func()) {
    tr := otel.Tracer(otel.GetMeterProvider().Meter("noop").InstrumentationLibrary().Name)
    // fallback tracer name: use service name
    if tr == nil { tr = otel.Tracer("ferret") }
    var kvs []attribute.KeyValue
    for k, v := range attrs { kvs = append(kvs, attribute.String(k, v)) }
    ctx, span := tr.Start(ctx, name)
    if len(kvs) > 0 { span.SetAttributes(kvs...) }
    return ctx, span.End
}

func RecordCounter(ctx context.Context, name string, delta float64, attrs map[string]string) {
    if meter == nil { return }
    c, ok := ctrs[name]
    if !ok {
        cc, err := meter.Int64Counter(name)
        if err != nil { log.Printf("otel: counter %s: %v", name, err); return }
        c = cc
        ctrs[name] = c
    }
    var opts []metric.AddOption
    if len(attrs) > 0 {
        var kvs []attribute.KeyValue
        for k, v := range attrs { kvs = append(kvs, attribute.String(k, v)) }
        opts = append(opts, metric.WithAttributes(kvs...))
    }
    c.Add(ctx, int64(delta), opts...)
}

func RecordHistogram(ctx context.Context, name string, value float64, attrs map[string]string) {
    if meter == nil { return }
    h, ok := histos[name]
    if !ok {
        hh, err := meter.Float64Histogram(name)
        if err != nil { log.Printf("otel: histogram %s: %v", name, err); return }
        h = hh
        histos[name] = h
    }
    var opts []metric.RecordOption
    if len(attrs) > 0 {
        var kvs []attribute.KeyValue
        for k, v := range attrs { kvs = append(kvs, attribute.String(k, v)) }
        opts = append(opts, metric.WithAttributes(kvs...))
    }
    h.Record(ctx, value, opts...)
}

