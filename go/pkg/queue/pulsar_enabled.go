//go:build pulsar
// +build pulsar

package queue

import (
    "context"
    "encoding/json"
    "os"

    "github.com/apache/pulsar-client-go/pulsar"
)

type PulsarClient struct {
    c pulsar.Client
}
type PulsarProducer struct {
    p pulsar.Producer
}

type PulsarConfig struct {
    ServiceURL string
    Token      string
    Topic      string
}

func NewPulsarClient(cfg PulsarConfig) (*PulsarClient, error) {
    opts := pulsar.ClientOptions{URL: cfg.ServiceURL}
    if cfg.Token != "" { opts.Authentication = pulsar.NewAuthenticationToken(cfg.Token) }
    client, err := pulsar.NewClient(opts)
    if err != nil { return nil, err }
    return &PulsarClient{c: client}, nil
}

func NewPulsarClientFromEnv() (*PulsarClient, PulsarConfig, error) {
    cfg := PulsarConfig{
        ServiceURL: os.Getenv("PULSAR_SERVICE_URL"),
        Token:      os.Getenv("PULSAR_TOKEN"),
        Topic:      os.Getenv("PULSAR_TOPIC_POST_EVENTS"),
    }
    if cfg.Topic == "" { cfg.Topic = "persistent://public/default/post-events" }
    c, err := NewPulsarClient(cfg)
    return c, cfg, err
}

func (c *PulsarClient) NewProducer(topic string) (*PulsarProducer, error) {
    p, err := c.c.CreateProducer(pulsar.ProducerOptions{Topic: topic})
    if err != nil { return nil, err }
    return &PulsarProducer{p: p}, nil
}

func (p *PulsarProducer) SendJSON(v any) error {
    b, err := json.Marshal(v)
    if err != nil { return err }
    _, err = p.p.Send(context.Background(), &pulsar.ProducerMessage{Payload: b})
    return err
}

func (p *PulsarProducer) Close() error { p.p.Close(); return nil }

