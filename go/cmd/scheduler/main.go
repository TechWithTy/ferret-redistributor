package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/lib/pq"

    "github.com/bitesinbyte/ferret/pkg/calendar"
    "github.com/bitesinbyte/ferret/pkg/metrics"
    "github.com/bitesinbyte/ferret/pkg/telemetry"
    "github.com/bitesinbyte/ferret/pkg/queue"
)

func main() {
    within := flag.Duration("within", 15*time.Minute, "time window to fetch due posts (e.g., 15m, 1h)")
    dsn := flag.String("database", os.Getenv("DATABASE_URL"), "Postgres DSN (or set DATABASE_URL)")
    jsonArray := flag.Bool("json-array", false, "output as a single JSON array instead of JSONL")
    flag.Parse()

    if *dsn == "" {
        log.Fatal("missing database DSN (set --database or DATABASE_URL)")
    }
    db, err := sql.Open("postgres", *dsn)
    if err != nil { log.Fatal(err) }
    defer db.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    ctx = telemetry.InitFromEnv(ctx)
    defer cancel()

    claimTimer := metrics.NewHistogram("scheduler_claim_seconds")
    // Optional Pulsar events emitter (schedule events)
    var prod *queue.PulsarProducer
    if pc, cfg, err := queue.NewPulsarClientFromEnv(); err == nil && cfg.ServiceURL != "" {
        topic := os.Getenv("PULSAR_TOPIC_SCHEDULE_EVENTS")
        if topic == "" { topic = "persistent://public/default/schedule-events" }
        if p, err := pc.NewProducer(topic); err == nil { prod = p; defer prod.Close() }
    }
    var rows []calendar.ScheduledPostRow
    claimTimer.Time(func() {
        var end func()
        ctx, end = telemetry.StartSpan(ctx, "scheduler.claim", map[string]string{"unit":"post"})
        rows, err = calendar.FetchAndClaimDuePosts(ctx, db, *within, 50)
        end()
    })
    if err != nil { log.Fatal(err) }
    metrics.NewCounter("scheduler_claimed_total").Inc(float64(len(rows)))

    enc := json.NewEncoder(os.Stdout)
    if *jsonArray {
        if err := enc.Encode(rows); err != nil { log.Fatal(err) }
        return
    }
    for _, r := range rows {
        if err := enc.Encode(r); err != nil { log.Fatal(err) }
    }

    fmt.Fprintf(os.Stderr, "claimed %d due posts\n", len(rows))
    telemetry.RecordCounter(ctx, "scheduler_claimed_total", float64(len(rows)), nil)
    // Emit schedule.claimed events
    if prod != nil {
        for _, r := range rows {
            _ = prod.SendJSON(map[string]any{
                "type": "schedule.claimed",
                "id": r.ID,
                "platform": string(r.Platform),
                "scheduled_at": r.ScheduledAt.Format(time.RFC3339),
                "campaign_id": r.CampaignID,
                "content_id": r.ContentID.String,
            })
        }
    }
    // Metrics snapshot for CI visibility
    if c,g,h := metrics.Snapshot(); true {
        log.Printf("metrics counters=%v gauges=%v histo_keys=%d", c, g, len(h))
    }
}
