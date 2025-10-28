package workflows

import (
    "time"
)

// ScheduleInstagramTrendAnalysis runs the provided fetch function every 24 hours in a background goroutine.
// It immediately invokes the function once, then on each tick.
func ScheduleInstagramTrendAnalysis(fetch func()) func() {
    stop := make(chan struct{})
    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()
        fetch()
        for {
            select {
            case <-ticker.C:
                fetch()
            case <-stop:
                return
            }
        }
    }()
    return func() { close(stop) }
}

