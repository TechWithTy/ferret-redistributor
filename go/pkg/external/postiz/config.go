package postiz

import (
	"os"
	"time"
)

const defaultBaseURL = "https://api.postiz.com/public/v1"

// Config controls how the Postiz client is created.
type Config struct {
	APIKey      string
	BaseURL     string
	HTTPTimeout time.Duration
}

// NewConfigFromEnv reads configuration from environment variables.
// POSTIZ_API_KEY is required, POSTIZ_BASE_URL is optional.
func NewConfigFromEnv() Config {
	return Config{
		APIKey:      os.Getenv("POSTIZ_API_KEY"),
		BaseURL:     os.Getenv("POSTIZ_BASE_URL"),
		HTTPTimeout: 30 * time.Second,
	}
}

func (c Config) normalizedBaseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return defaultBaseURL
}

func (c Config) normalizedTimeout() time.Duration {
	if c.HTTPTimeout <= 0 {
		return 30 * time.Second
	}
	return c.HTTPTimeout
}
