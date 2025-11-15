package workers

import "fmt"

type CRMWorker struct{}

func (c CRMWorker) Push(record string) {
    fmt.Printf("CRMWorker: Pushing record %q to CRM (mock)\n", record)
}

