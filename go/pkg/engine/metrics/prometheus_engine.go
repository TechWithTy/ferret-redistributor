package metrics

import "fmt"

type Prometheus struct{}

func (p Prometheus) IncCounter(name string) {
    fmt.Printf("Metrics: increment %s\n", name)
}

