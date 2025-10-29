package workflows

import (
    "context"
    "database/sql"
    "time"

    "github.com/bitesinbyte/ferret/pkg/calendar"
)

type Scheduler struct {
    DB *sql.DB
}

// Claim returns up to limit due posts within the window and moves them to processing.
func (s *Scheduler) Claim(ctx context.Context, within time.Duration, limit int) ([]calendar.ScheduledPostRow, error) {
    return calendar.FetchAndClaimDuePosts(ctx, s.DB, within, limit)
}

