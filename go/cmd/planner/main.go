package main

import (
    "context"
    "database/sql"
    "flag"
    "log"
    "os"
    "time"

    _ "github.com/lib/pq"

    "github.com/bitesinbyte/ferret/pkg/adapters/calendarrepo"
    "github.com/bitesinbyte/ferret/pkg/generator"
)

func main() {
    var (
        trendsFile   = flag.String("trends", "_data/trends.json", "path to trends JSON")
        variantsFile = flag.String("variants", "_data/variants.json", "path to variants JSON")
        dsn          = flag.String("database", os.Getenv("DATABASE_URL"), "Postgres DSN")
        spacing      = flag.Duration("spacing", 2*time.Hour, "min spacing per platform")
        start        = flag.Duration("start-offset", 0, "offset from now for first slot")
    )
    flag.Parse()
    if *dsn == "" { log.Fatal("missing database DSN") }

    trends, err := generator.LoadTrends(*trendsFile)
    if err != nil { log.Fatal(err) }
    variants, err := generator.LoadVariants(*variantsFile)
    if err != nil { log.Fatal(err) }

    db, err := sql.Open("postgres", *dsn)
    if err != nil { log.Fatal(err) }
    defer db.Close()
    repo := calendarrepo.Repository{DB: db}

    in := generator.PlanInput{
        Trends:    trends,
        Variants:  variants,
        Platforms: []string{"linkedin", "twitter"},
        StartAt:   time.Now().UTC().Add(*start),
        Spacing:   *spacing,
        PerDayLimit: 10,
    }

    if err := generator.PlanAndSchedule(context.Background(), repo, in); err != nil { log.Fatal(err) }
}
