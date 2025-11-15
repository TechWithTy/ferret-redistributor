package linkedin

import (
    "os"
    "time"
)

type Config struct {
    AccessToken string
    HTTPTimeout time.Duration
}

func NewFromEnv() Config {
    return Config{
        AccessToken: os.Getenv("LINKEDIN_ACCESS_TOKEN"),
        HTTPTimeout: 30 * time.Second,
    }
}

