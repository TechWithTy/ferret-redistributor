package queue

import "fmt"

type WorkerQueue struct{}

func (w WorkerQueue) Enqueue(name string, payload any) {
    fmt.Printf("Queue: Enqueued job %s\n", name)
}

