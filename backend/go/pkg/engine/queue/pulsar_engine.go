package queue

import "fmt"

type PulsarEngine struct{}

func (p PulsarEngine) Send(topic string, payload any) {
    fmt.Printf("Queue: Sending to Pulsar topic=%s (mock)\n", topic)
}

