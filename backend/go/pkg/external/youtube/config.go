package youtube

import (
    "os"
    "time"
)

type Config struct {
    APIKey      string
    ChannelID   string
    HTTPTimeout time.Duration
}

func NewFromEnv() Config {
    return Config{
        APIKey:      os.Getenv("YOUTUBE_API_KEY"),
        ChannelID:   os.Getenv("YOUTUBE_CHANNEL_ID"),
        HTTPTimeout: 30 * time.Second,
    }
}

