package behiiv

import (
    "os"
    "time"
)

type Config struct {
    BaseURL     string
    Version     string
    Token       string
    HTTPTimeout time.Duration
}

func NewFromEnv() Config {
    base := os.Getenv("BEEHIIV_BASE_URL")
    if base == "" { base = "https://api.beehiiv.com" }
    return Config{
        BaseURL:     base,
        Version:     "v2",
        Token:       os.Getenv("BEEHIIV_TOKEN"),
        HTTPTimeout: 30 * time.Second,
    }
}
