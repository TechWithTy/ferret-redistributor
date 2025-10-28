package metrics

import (
    "sync"
    "time"
)

type Counter struct{ name string }
type Gauge struct{ name string }
type Histogram struct{ name string }

var (
    mu       sync.Mutex
    counters = map[string]float64{}
    gauges   = map[string]float64{}
    histos   = map[string][]float64{}
)

func NewCounter(name string) Counter { return Counter{name: name} }
func NewGauge(name string) Gauge     { return Gauge{name: name} }
func NewHistogram(name string) Histogram { return Histogram{name: name} }

func (c Counter) Inc(delta float64) {
    mu.Lock(); defer mu.Unlock()
    counters[c.name] += delta
}
func (g Gauge) Set(v float64) {
    mu.Lock(); defer mu.Unlock()
    gauges[g.name] = v
}
func (h Histogram) Observe(v float64) {
    mu.Lock(); defer mu.Unlock()
    histos[h.name] = append(histos[h.name], v)
}

// Helper to time a function and observe its duration in seconds.
func (h Histogram) Time(fn func()) {
    start := time.Now()
    fn()
    h.Observe(time.Since(start).Seconds())
}

// Snapshot returns a shallow copy of current metric values.
func Snapshot() (map[string]float64, map[string]float64, map[string][]float64) {
    mu.Lock(); defer mu.Unlock()
    c := make(map[string]float64, len(counters))
    for k, v := range counters { c[k] = v }
    g := make(map[string]float64, len(gauges))
    for k, v := range gauges { g[k] = v }
    h := make(map[string][]float64, len(histos))
    for k, v := range histos { cp := make([]float64, len(v)); copy(cp, v); h[k] = cp }
    return c, g, h
}

