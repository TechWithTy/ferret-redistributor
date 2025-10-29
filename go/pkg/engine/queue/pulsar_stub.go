package queue

import (
    "errors"
    "os"
)

type PulsarClient struct{}
type PulsarProducer struct{}

type PulsarConfig struct {
    ServiceURL string
    Token      string
    Topic      string
}

func NewPulsarClient(cfg PulsarConfig) (*PulsarClient, error) {
    // Default build has no real client
    return nil, errors.New("pulsar: disabled (build with -tags=pulsar)")
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
    return nil, errors.New("pulsar: disabled")
}

func (p *PulsarProducer) SendJSON(v any) error { return errors.New("pulsar: disabled") }
func (p *PulsarProducer) Close() error         { return nil }

