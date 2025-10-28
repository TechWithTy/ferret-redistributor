package generator

import "time"

type Trend struct {
    Topic       string  `json:"topic"`
    Relevance   float64 `json:"relevance"`
    Volume      int     `json:"volume"`
    Competition float64 `json:"competition"`
    TrendScore  float64 `json:"trend_score"`
}

type Variant struct {
    ID        string `json:"id"`
    Title     string `json:"title,omitempty"`
    CTA       string `json:"cta,omitempty"`
    IsControl bool   `json:"is_control"`
}

type PlanInput struct {
    Trends      []Trend
    Variants    map[string][]Variant // topic -> variants
    Platforms   []string
    StartAt     time.Time
    Spacing     time.Duration // min spacing per platform
    PerDayLimit int
}

