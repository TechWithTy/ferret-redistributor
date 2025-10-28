package calendarrepo

import (
    "context"
    "database/sql"
    "time"
)

// Repository provides write helpers for scheduling posts.
type Repository struct{ DB *sql.DB }

type ScheduleInput struct {
    ID          string
    CampaignID  *string
    ContentID   *string
    Platform    string
    Caption     *string
    Hashtags    *string
    ScheduledAt time.Time
    MetadataJSON *string // JSON string; nullable
}

// SchedulePost inserts a single scheduled post row.
func (r Repository) SchedulePost(ctx context.Context, in ScheduleInput) error {
    const q = `INSERT INTO scheduled_posts
    (id, campaign_id, content_id, platform, caption, hashtags, scheduled_at, status, metadata, created_at, updated_at)
    VALUES ($1,$2,$3,$4,$5,$6,$7,'scheduled',COALESCE($8,'{}'::json), NOW(), NOW())`
    _, err := r.DB.ExecContext(ctx, q,
        in.ID, in.CampaignID, in.ContentID, in.Platform, in.Caption, in.Hashtags, in.ScheduledAt, in.MetadataJSON,
    )
    return err
}

// BulkSchedule inserts multiple posts in a transaction.
func (r Repository) BulkSchedule(ctx context.Context, items []ScheduleInput) error {
    if len(items) == 0 { return nil }
    tx, err := r.DB.BeginTx(ctx, nil)
    if err != nil { return err }
    const q = `INSERT INTO scheduled_posts
    (id, campaign_id, content_id, platform, caption, hashtags, scheduled_at, status, metadata, created_at, updated_at)
    VALUES ($1,$2,$3,$4,$5,$6,$7,'scheduled',COALESCE($8,'{}'::json), NOW(), NOW())`
    stmt, err := tx.PrepareContext(ctx, q)
    if err != nil { _ = tx.Rollback(); return err }
    defer stmt.Close()
    for _, in := range items {
        if _, err := stmt.ExecContext(ctx,
            in.ID, in.CampaignID, in.ContentID, in.Platform, in.Caption, in.Hashtags, in.ScheduledAt, in.MetadataJSON,
        ); err != nil { _ = tx.Rollback(); return err }
    }
    return tx.Commit()
}

