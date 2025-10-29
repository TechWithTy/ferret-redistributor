package workers

import "fmt"

type DataWorker struct{}

func (w DataWorker) Transform(s string) string {
    fmt.Println("DataWorker: Transforming data")
    return s + " + transformed"
}

