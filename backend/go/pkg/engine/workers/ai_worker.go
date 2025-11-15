package workers

import "fmt"

type AIWorker struct{}

func (a AIWorker) Enrich(data string) string {
    fmt.Println("AIWorker: Enriching data with AI")
    return data + " + enriched"
}

