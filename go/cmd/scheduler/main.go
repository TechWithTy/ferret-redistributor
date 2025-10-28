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
    defer cancel()

    claimTimer := metrics.NewHistogram("scheduler_claim_seconds")
    var rows []calendar.ScheduledPostRow
    claimTimer.Time(func() {
        rows, err = calendar.FetchAndClaimDuePosts(ctx, db, *within, 50)
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
    // Metrics snapshot for CI visibility
    if c,g,h := metrics.Snapshot(); true {
        log.Printf("metrics counters=%v gauges=%v histo_keys=%d", c, g, len(h))
    }
}
