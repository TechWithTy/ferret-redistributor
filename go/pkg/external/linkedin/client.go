package linkedin

import (
    "net/http"
    "time"
)

type Client struct {
    httpClient *http.Client
    cfg        Config
}

func New(cfg Config) *Client {
    to := cfg.HTTPTimeout
    if to <= 0 { to = 30 * time.Second }
    return &Client{httpClient: &http.Client{Timeout: to}, cfg: cfg}
}

