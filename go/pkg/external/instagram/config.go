package instagram

import (
    "os"
    "time"
)

type Config struct {
    BaseURL     string
    Version     string
    IGUserID    string
    AccessToken string
    HTTPTimeout time.Duration
}

func NewFromEnv() Config {
    version := os.Getenv("IG_GRAPH_VERSION")
    if version == "" {
        version = "v19.0"
    }
    baseURL := "https://graph.facebook.com/" + version
    timeout := 30 * time.Second
    return Config{
        BaseURL:     baseURL,
        Version:     version,
        IGUserID:    os.Getenv("IG_USER_ID"),
        AccessToken: os.Getenv("IG_ACCESS_TOKEN"),
        HTTPTimeout: timeout,
    }
}
