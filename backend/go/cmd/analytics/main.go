package main

import (
    "bufio"
    "context"
    "encoding/json"
    "errors"
    "flag"
    "log"
    "os"
    "time"

    ig "github.com/bitesinbyte/ferret/pkg/external/instagram"
    li "github.com/bitesinbyte/ferret/pkg/external/linkedin"
    yt "github.com/bitesinbyte/ferret/pkg/external/youtube"
    "github.com/bitesinbyte/ferret/pkg/analytics/storage"
)

func main() {
    platform := flag.String("platform", "instagram", "platform: instagram|linkedin|youtube")
    id := flag.String("id", "", "media/post ID or URN")
    // Instagram user insights
    period := flag.String("period", "day", "instagram user insights period: day|week|days_28|month|lifetime")
    metrics := flag.String("metrics", "", "comma-separated metrics (instagram insights)")
    // Bulk mode (reads stored published IDs)
    bulk := flag.Bool("bulk", false, "bulk mode: read published IDs from file and fetch analytics")
    file := flag.String("file", "data/published_posts.jsonl", "path to published posts JSONL file")
    since := flag.String("since", "", "RFC3339 timestamp or duration (e.g., 72h) to filter records")
    flag.Parse()

    ctx := context.Background()
    switch *platform {
    case "instagram":
        cfg := ig.NewFromEnv()
        client := ig.New(cfg)
        if *bulk {
            cutoff, err := parseSince(*since)
            if err != nil { log.Fatal(err) }
            res, err := bulkInstagram(ctx, client, *file, cutoff, *metrics)
            if err != nil { log.Fatal(err) }
            out(res)
            return
        }
        if *id != "" {
            // Fetch media basic + insights
            basic, err := client.GetMediaBasic(ctx, *id)
            if err != nil { log.Fatal(err) }
            out(basic)
            if *metrics != "" {
                ms := splitComma(*metrics)
                ins, err := client.GetMediaInsights(ctx, *id, ms)
                if err != nil { log.Fatal(err) }
                out(ins)
            }
            return
        }
        if *metrics != "" {
            ms := splitComma(*metrics)
            ins, err := client.GetUserInsights(ctx, ms, *period, "", "")
            if err != nil { log.Fatal(err) }
            out(ins)
            return
        }
        log.Fatalf("instagram: require --id or --metrics for user insights")
    case "linkedin":
        cfg := li.NewFromEnv()
        client := li.New(cfg)
        if *bulk {
            cutoff, err := parseSince(*since)
            if err != nil { log.Fatal(err) }
            res, err := bulkLinkedIn(ctx, client, *file, cutoff)
            if err != nil { log.Fatal(err) }
            out(res)
            return
        }
        if *id == "" { log.Fatalf("linkedin: require --id post URN or ID") }
        stats, err := client.GetPostStatistics(ctx, *id)
        if err != nil { log.Fatal(err) }
        out(stats)
    case "youtube":
        if *id == "" { log.Fatalf("youtube: require --id video ID") }
        cfg := yt.NewFromEnv()
        client := yt.New(cfg)
        stats, err := client.GetVideoStatistics(ctx, *id)
        if err != nil { log.Fatal(err) }
        out(stats)
    default:
        log.Fatalf("unknown platform: %s", *platform)
    }
}

func out(v any) {
    enc := json.NewEncoder(os.Stdout)
    enc.SetIndent("", "  ")
    _ = enc.Encode(v)
}

func splitComma(s string) []string {
    parts := []string{}
    cur := []rune{}
    for _, r := range s {
        if r == ',' {
            parts = append(parts, string(cur))
            cur = cur[:0]
            continue
        }
        cur = append(cur, r)
    }
    parts = append(parts, string(cur))
    out := make([]string, 0, len(parts))
    for _, p := range parts {
        p = trimSpace(p)
        if p != "" { out = append(out, p) }
    }
    return out
}

func trimSpace(s string) string {
    // simple space trim without strings import
    runes := []rune(s)
    i, j := 0, len(runes)-1
    for i <= j && (runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\n' || runes[i] == '\r') { i++ }
    for j >= i && (runes[j] == ' ' || runes[j] == '\t' || runes[j] == '\n' || runes[j] == '\r') { j-- }
    if i > j { return "" }
    return string(runes[i : j+1])
}

func parseSince(s string) (time.Time, error) {
    if s == "" {
        return time.Now().Add(-24 * time.Hour), nil
    }
    if t, err := time.Parse(time.RFC3339, s); err == nil {
        return t, nil
    }
    if d, err := time.ParseDuration(s); err == nil {
        return time.Now().Add(-d), nil
    }
    // try date only
    if t, err := time.Parse("2006-01-02", s); err == nil {
        return t, nil
    }
    return time.Time{}, errors.New("invalid --since; use RFC3339, YYYY-MM-DD, or duration like 72h")
}

func bulkInstagram(ctx context.Context, client *ig.Client, path string, cutoff time.Time, metrics string) (any, error) {
    f, err := os.Open(path)
    if err != nil { return nil, err }
    defer f.Close()
    type rec struct {
        Platform    string    `json:"platform"`
        ID          string    `json:"id"`
        Link        string    `json:"link"`
        ContentType string    `json:"content_type"`
        PublishedAt time.Time `json:"published_at"`
    }
    scanner := bufio.NewScanner(f)
    out := []map[string]any{}
    var mets []string
    if metrics != "" { mets = splitComma(metrics) }
    for scanner.Scan() {
        line := scanner.Bytes()
        var r rec
        if err := json.Unmarshal(line, &r); err != nil { continue }
        if r.Platform != "instagram" { continue }
        if r.PublishedAt.Before(cutoff) { continue }
        if r.ID == "" { continue }
        basic, err := client.GetMediaBasic(ctx, r.ID)
        if err != nil { continue }
        row := map[string]any{"id": r.ID, "link": r.Link, "published_at": r.PublishedAt, "basic": basic}
        if len(mets) > 0 {
            ins, err := client.GetMediaInsights(ctx, r.ID, mets)
            if err == nil { row["insights"] = ins }
        }
        out = append(out, row)
    }
    return out, scanner.Err()
}

func bulkLinkedIn(ctx context.Context, client *li.Client, path string, cutoff time.Time) (any, error) {
    f, err := os.Open(path)
    if err != nil { return nil, err }
    defer f.Close()
    type rec storage.PublishedPost
    scanner := bufio.NewScanner(f)
    out := []map[string]any{}
    for scanner.Scan() {
        line := scanner.Bytes()
        var r rec
        if err := json.Unmarshal(line, &r); err != nil { continue }
        if r.Platform != "linkedin" { continue }
        if r.PublishedAt.Before(cutoff) { continue }
        if r.ID == "" { continue }
        stats, err := client.GetPostStatistics(ctx, r.ID)
        if err != nil { continue }
        row := map[string]any{"id": r.ID, "link": r.Link, "published_at": r.PublishedAt, "stats": stats}
        out = append(out, row)
    }
    return out, scanner.Err()
}
