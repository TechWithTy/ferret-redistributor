package generator

import (
    "context"
    "encoding/json"
    "os"
    "strings"
    "time"

    ig "github.com/bitesinbyte/ferret/pkg/external/instagram"
)

// LoadTrends reads JSON array of trends from a file (e.g., _data/trends.json).
func LoadTrends(path string) ([]Trend, error) {
    b, err := os.ReadFile(path)
    if err != nil { return nil, err }
    var items []Trend
    if err := json.Unmarshal(b, &items); err != nil { return nil, err }
    return items, nil
}

// TrendAnalyzer bridges the Instagram client with generator Trend types.
type TrendAnalyzer struct {
    IG *ig.Client
}

// AnalyzeHashtag fetches recent media for a hashtag and produces an aggregate Trend.
func (t *TrendAnalyzer) AnalyzeHashtag(ctx context.Context, hashtag string) (Trend, error) {
    tag := strings.TrimPrefix(strings.TrimSpace(hashtag), "#")
    if tag == "" {
        return Trend{}, nil
    }
    id, err := t.IG.SearchHashtag(ctx, tag)
    if err != nil {
        return Trend{Topic: tag, TrendScore: 0}, nil
    }
    media, err := t.IG.GetHashtagRecentMedia(ctx, id, 20, 0, 0)
    if err != nil {
        return Trend{Topic: tag, TrendScore: 0}, nil
    }
    var vol int
    var totalScore float64
    for _, m := range media {
        vol++
        totalScore += t.CalculateVirality(m)
    }
    avg := 0.0
    if vol > 0 {
        avg = totalScore / float64(vol)
    }
    // Map to Trend fields; simple heuristic values for relevance/competition
    return Trend{
        Topic:       tag,
        Relevance:   min1(avg/1000.0),
        Volume:      vol,
        Competition: 0.5,
        TrendScore:  min100(avg),
    }, nil
}

// CalculateVirality provides a simple engagement-velocity style score.
func (t *TrendAnalyzer) CalculateVirality(m ig.MediaItem) float64 {
    // likes + 2*comments, (optionally) adjusted by age
    base := float64(m.LikeCount + 2*m.CommentsCount)
    // parse timestamp if available
    ageH := 1.0
    if ts := parseTime(m.Timestamp); !ts.IsZero() {
        ageH = float64(time.Since(ts).Hours())
        if ageH < 1.0 { ageH = 1.0 }
    }
    return base / (ageH * 0.7)
}

func parseTime(s string) time.Time {
    if s == "" { return time.Time{} }
    // Accept RFC3339-like timestamps (Instagram returns ISO 8601 with Z)
    if ts, err := time.Parse(time.RFC3339, strings.ReplaceAll(s, "Z", "+00:00")); err == nil {
        return ts
    }
    return time.Time{}
}

func min1(v float64) float64 { if v < 0 { return 0 }; if v > 1 { return 1 }; return v }
func min100(v float64) float64 { if v < 0 { return 0 }; if v > 100 { return 100 }; return v }

