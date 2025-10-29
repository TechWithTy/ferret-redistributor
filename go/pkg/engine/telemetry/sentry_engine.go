package telemetry

import "fmt"

type Sentry struct{}

func (s Sentry) TrackEvent(event string) {
    fmt.Printf("Telemetry: Tracking event -> %s\n", event)
}

