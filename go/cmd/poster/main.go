package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "strings"
    "time"

    _ "github.com/lib/pq"

    "github.com/bitesinbyte/ferret/pkg/calendar"
    "github.com/bitesinbyte/ferret/pkg/config"
    "github.com/bitesinbyte/ferret/pkg/external"
    "github.com/bitesinbyte/ferret/pkg/metrics"
    "strconv"
)

func main() {
    input := flag.String("input", "go/due_posts.json", "path to JSON array of due posts")
    dsn := flag.String("database", os.Getenv("DATABASE_URL"), "Postgres DSN (or set DATABASE_URL)")
    cfgPath := flag.String("config", "config.json", "path to config.json")
    flag.Parse()

    if *dsn == "" { log.Fatal("missing database DSN (set --database or DATABASE_URL)") }

    // Read input file
    f, err := os.Open(*input)
    if err != nil { log.Fatal(err) }
    defer f.Close()
    body, err := io.ReadAll(f)
    if err != nil { log.Fatal(err) }
    var rows []calendar.ScheduledPostRow
    if err := json.Unmarshal(body, &rows); err != nil { log.Fatalf("invalid JSON: %v", err) }
    if len(rows) == 0 {
        fmt.Fprintln(os.Stderr, "no items to post")
        return
    }

    cfg := config.LoadConfig(*cfgPath)

    db, err := sql.Open("postgres", *dsn)
    if err != nil { log.Fatal(err) }
    defer db.Close()

    ctx := context.Background()

    postedCounter := metrics.NewCounter("poster_posted_total")
    failedCounter := metrics.NewCounter("poster_failed_total")
    postLatency := metrics.NewHistogram("poster_post_seconds")
    limiter := newPlatformLimiters(loadRatesFromEnv(), time.Second)

    for _, r := range rows {
        platform := strings.ToLower(string(r.Platform))
        poster, supported := createPoster(platform)
        if !supported {
            markFailed(ctx, db, r.ID, fmt.Sprintf("unsupported platform: %s", platform))
            failedCounter.Inc(1)
            continue
        }
        // Build the external.Post
        title := firstNonEmpty(r.ContentTitle.String, r.CampaignName)
        hashtags := r.Hashtags.String
        link := r.ContentURL.String
        post := external.Post{Title: title, Link: link, HashTags: hashtags, Description: ""}

        publishedAt := time.Now().UTC()
        // If poster supports ID, capture and update
        var perr error
        postLatency.Time(func() {
            // pace
            limiter.Take(platform)
            perr = doWithRetry(func() error {
                if pwid, ok := poster.(external.PosterWithID); ok {
                    id, err := pwid.PostWithID(cfg, post)
                    if err != nil { return err }
                    return calendar.UpdatePostStatus(ctx, db, r.ID, calendar.StatusPublished, &id, &publishedAt, nil)
                }
                if err := poster.Post(cfg, post); err != nil { return err }
                return calendar.UpdatePostStatus(ctx, db, r.ID, calendar.StatusPublished, nil, &publishedAt, nil)
            })
        })
        if perr != nil {
            markFailed(ctx, db, r.ID, perr.Error())
            failedCounter.Inc(1)
        } else {
            postedCounter.Inc(1)
        }
    }
    // Metrics snapshot for CI visibility
    if c,g,h := metrics.Snapshot(); true {
        // compute simple stats for poster_post_seconds if present
        if samples, ok := h["poster_post_seconds"]; ok && len(samples) > 0 {
            minv, maxv, sum := samples[0], samples[0], 0.0
            for _, v := range samples {
                if v < minv { minv = v }
                if v > maxv { maxv = v }
                sum += v
            }
            avg := sum / float64(len(samples))
            log.Printf("latency poster_post_seconds: count=%d min=%.3fs avg=%.3fs max=%.3fs", len(samples), minv, avg, maxv)
        }
        log.Printf("metrics counters=%v gauges=%v histo_keys=%d", c, g, len(h))
    }
}

func createPoster(platform string) (external.Poster, bool) {
    switch platform {
    case "linkedin":
        return external.Linkedin{}, true
    case "mastodon":
        return external.Mastodon{}, true
    case "twitter":
        return external.Twitter{}, true
    case "facebook":
        return external.Facebook{}, true
    case "thread":
        return external.Thread{}, true
    case "instagram":
        return external.Instagram{}, true
    default:
        return nil, false
    }
}

func markFailed(ctx context.Context, db *sql.DB, id, msg string) {
    meta := map[string]any{"error": msg}
    b, _ := json.Marshal(meta)
    _ = calendar.UpdatePostStatus(ctx, db, id, calendar.StatusFailed, nil, nil, b)
}

func firstNonEmpty(val ...string) string {
    for _, v := range val {
        if strings.TrimSpace(v) != "" { return v }
    }
    return ""
}

// Retry helper with exponential backoff and jitter.
func doWithRetry(fn func() error) error {
    backoff := 500 * time.Millisecond
    for attempt := 0; attempt < 5; attempt++ {
        if err := fn(); err != nil {
            if !isRetryable(err) { return err }
            time.Sleep(backoff + time.Duration(attempt*100)*time.Millisecond)
            if backoff < 5*time.Second { backoff *= 2 }
            continue
        }
        return nil
    }
    return errors.New("exhausted retries")
}

func isRetryable(err error) bool {
    // naive mapping; could inspect error strings or types from platform clients
    s := strings.ToLower(err.Error())
    return strings.Contains(s, "429") || strings.Contains(s, "rate") || strings.Contains(s, "timeout") || strings.Contains(s, "temporar")
}

// simple token bucket: allow n tokens per window.
type rateLimiter struct {
    ch chan struct{}
}

func newRateLimiter(n int, window time.Duration) *rateLimiter {
    if n <= 0 { n = 1 }
    rl := &rateLimiter{ch: make(chan struct{}, n)}
    // fill
    for i := 0; i < n; i++ { rl.ch <- struct{}{} }
    go func() {
        ticker := time.NewTicker(window)
        defer ticker.Stop()
        for range ticker.C {
            for len(rl.ch) < n { rl.ch <- struct{}{} }
        }
    }()
    return rl
}

func (r *rateLimiter) Take() { <-r.ch }

// per-platform limiter
type platformLimiters struct{
    m map[string]*rateLimiter
    def *rateLimiter
}

func newPlatformLimiters(per map[string]int, window time.Duration) *platformLimiters {
    m := make(map[string]*rateLimiter, len(per))
    for k, n := range per {
        if n <= 0 { n = 1 }
        m[strings.ToLower(k)] = newRateLimiter(n, window)
    }
    return &platformLimiters{m: m, def: newRateLimiter(2, window)}
}
func (p *platformLimiters) Take(platform string) {
    if rl, ok := p.m[strings.ToLower(platform)]; ok { rl.Take(); return }
    p.def.Take()
}

// read poster rates from env, like POSTER_RATE_LINKEDIN=2, POSTER_RATE_INSTAGRAM=1
func loadRatesFromEnv() map[string]int {
    mapping := map[string]string{
        "linkedin":  os.Getenv("POSTER_RATE_LINKEDIN"),
        "twitter":   os.Getenv("POSTER_RATE_TWITTER"),
        "facebook":  os.Getenv("POSTER_RATE_FACEBOOK"),
        "thread":    os.Getenv("POSTER_RATE_THREAD"),
        "instagram": os.Getenv("POSTER_RATE_INSTAGRAM"),
        "mastodon":  os.Getenv("POSTER_RATE_MASTODON"),
    }
    out := map[string]int{}
    for k, v := range mapping {
        if v == "" { continue }
        if n, err := strconv.Atoi(v); err == nil && n > 0 { out[k] = n }
    }
    // defaults if unset
    if _, ok := out["linkedin"]; !ok { out["linkedin"] = 1 }
    if _, ok := out["instagram"]; !ok { out["instagram"] = 1 }
    return out
}
