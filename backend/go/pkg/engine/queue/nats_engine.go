package queue

import "fmt"

type NatsEngine struct{}

func (n NatsEngine) Publish(payload any) {
    fmt.Println("Queue: Publishing to NATS (mock)")
}

