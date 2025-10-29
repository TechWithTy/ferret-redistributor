# Metrics (Go)

Tiny in-process metrics for counters, gauges, histograms.

- `metrics.go`: Counter, Gauge, Histogram, Snapshot.

## Example
```
lat := metrics.NewHistogram("poster_post_seconds")
lat.Time(func(){
  // do work
})
metrics.NewCounter("poster_posted_total").Inc(1)

c,g,h := metrics.Snapshot()
log.Printf("counters=%v gauges=%v histo_keys=%d", c, g, len(h))
```

