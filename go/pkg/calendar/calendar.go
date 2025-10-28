package calendar

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "time"
)

// Platform mirrors the Python enum for platform names.
type Platform string

const (
    PlatformInstagram Platform = "instagram"
    PlatformLinkedIn  Platform = "linkedin"
    PlatformTwitter   Platform = "twitter"
    PlatformFacebook  Platform = "facebook"
    PlatformYouTube   Platform = "youtube"
    PlatformBeehiiv   Platform = "behiiv"
)

// ScheduledStatus mirrors the Python enum for scheduled post status.
type ScheduledStatus string

const (
    StatusScheduled ScheduledStatus = "scheduled"
    StatusPublished ScheduledStatus = "published"
    StatusFailed    ScheduledStatus = "failed"
    StatusCanceled  ScheduledStatus = "canceled"
)

// ScheduledPostRow represents a scheduled post with minimal joined context.
type ScheduledPostRow struct {
    ID          string
    CampaignID  string
    CampaignName string
    ContentID   sql.NullString
    ContentTitle sql.NullString
    ContentURL  sql.NullString
    Platform    Platform
    Caption     sql.NullString
    Hashtags    sql.NullString
    ScheduledAt time.Time
    Status      ScheduledStatus
    ExternalID  sql.NullString
    PublishedAt sql.NullTime
    Metadata    json.RawMessage
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// FetchScheduledPostsWithin returns posts in [start, end) with status=scheduled.
// The caller must pass a *sql.DB connected to Postgres (import a driver in main).
func FetchScheduledPostsWithin(ctx context.Context, db *sql.DB, start, end time.Time) ([]ScheduledPostRow, error) {
    const q = `
SELECT sp.id, sp.campaign_id, c.name AS campaign_name,
       sp.content_id, ci.title AS content_title, ci.canonical_url AS content_url,
       sp.platform, sp.caption, sp.hashtags,
       sp.scheduled_at, sp.status, sp.external_id, sp.published_at, sp.metadata,
       sp.created_at, sp.updated_at
FROM scheduled_posts sp
LEFT JOIN campaigns c ON sp.campaign_id = c.id
LEFT JOIN content_items ci ON sp.content_id = ci.id
WHERE sp.status = 'scheduled' AND sp.scheduled_at >= $1 AND sp.scheduled_at < $2
ORDER BY sp.scheduled_at ASC`

    rows, err := db.QueryContext(ctx, q, start, end)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    out := make([]ScheduledPostRow, 0, 32)
    for rows.Next() {
        var r ScheduledPostRow
        var platform string
        var status string
        var metaBytes sql.NullString
        if err := rows.Scan(
            &r.ID, &r.CampaignID, &r.CampaignName,
            &r.ContentID, &r.ContentTitle, &r.ContentURL,
            &platform, &r.Caption, &r.Hashtags,
            &r.ScheduledAt, &status, &r.ExternalID, &r.PublishedAt, &metaBytes,
            &r.CreatedAt, &r.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        r.Platform = Platform(platform)
        r.Status = ScheduledStatus(status)
        if metaBytes.Valid && len(metaBytes.String) > 0 {
            r.Metadata = json.RawMessage(metaBytes.String)
        }
        out = append(out, r)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return out, nil
}

// FetchDuePosts returns posts scheduled from now up to now+within.
func FetchDuePosts(ctx context.Context, db *sql.DB, within time.Duration) ([]ScheduledPostRow, error) {
    if within <= 0 {
        return nil, errors.New("within must be > 0")
    }
    now := time.Now().UTC()
    end := now.Add(within)
    return FetchScheduledPostsWithin(ctx, db, now, end)
}

// FetchAndClaimDuePosts atomically moves due posts to 'processing' and returns them.
// This prevents multiple runners from posting the same items. Uses SKIP LOCKED.
func FetchAndClaimDuePosts(ctx context.Context, db *sql.DB, within time.Duration, limit int) ([]ScheduledPostRow, error) {
    if within <= 0 { return nil, errors.New("within must be > 0") }
    if limit <= 0 { limit = 50 }
    now := time.Now().UTC()
    end := now.Add(within)

    const q = `
WITH cte AS (
    SELECT id FROM scheduled_posts
    WHERE status = 'scheduled' AND scheduled_at >= $1 AND scheduled_at < $2
    ORDER BY scheduled_at ASC
    FOR UPDATE SKIP LOCKED
    LIMIT $3
)
UPDATE scheduled_posts sp
SET status = 'processing', updated_at = NOW()
FROM cte
WHERE sp.id = cte.id
RETURNING sp.id, sp.campaign_id, c.name AS campaign_name,
          sp.content_id, ci.title AS content_title, ci.canonical_url AS content_url,
          sp.platform, sp.caption, sp.hashtags,
          sp.scheduled_at, sp.status, sp.external_id, sp.published_at, sp.metadata,
          sp.created_at, sp.updated_at
`

    rows, err := db.QueryContext(ctx, q, now, end, limit)
    if err != nil { return nil, err }
    defer rows.Close()

    out := make([]ScheduledPostRow, 0, limit)
    for rows.Next() {
        var r ScheduledPostRow
        var platform string
        var status string
        var metaBytes sql.NullString
        if err := rows.Scan(
            &r.ID, &r.CampaignID, &r.CampaignName,
            &r.ContentID, &r.ContentTitle, &r.ContentURL,
            &platform, &r.Caption, &r.Hashtags,
            &r.ScheduledAt, &status, &r.ExternalID, &r.PublishedAt, &metaBytes,
            &r.CreatedAt, &r.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        r.Platform = Platform(platform)
        r.Status = ScheduledStatus(status)
        if metaBytes.Valid && len(metaBytes.String) > 0 {
            r.Metadata = json.RawMessage(metaBytes.String)
        }
        out = append(out, r)
    }
    if err := rows.Err(); err != nil { return nil, err }
    return out, nil
}

// UpdatePostStatus updates status and optional external metadata after publish/fail.
// Pass publishedAt non-zero when marking as published.
func UpdatePostStatus(ctx context.Context, db *sql.DB, id string, status ScheduledStatus, externalID *string, publishedAt *time.Time, metadata json.RawMessage) error {
    // Build dynamic update based on provided fields.
    const base = `UPDATE scheduled_posts SET status = $1, updated_at = NOW()`
    args := []any{string(status)}
    set := ""
    idx := 2
    if externalID != nil {
        set += ", external_id = $" + itoa(idx)
        args = append(args, *externalID)
        idx++
    }
    if publishedAt != nil {
        set += ", published_at = $" + itoa(idx)
        args = append(args, *publishedAt)
        idx++
    }
    if metadata != nil {
        set += ", metadata = $" + itoa(idx)
        args = append(args, string(metadata))
        idx++
    }
    where := ", updated_at = NOW() WHERE id = $" + itoa(idx)
    args = append(args, id)
    q := base + set + where
    _, err := db.ExecContext(ctx, q, args...)
    return err
}

// itoa is a tiny helper avoiding strconv to keep deps minimal.
func itoa(n int) string {
    if n == 0 {
        return "0"
    }
    var buf [16]byte
    i := len(buf)
    for n > 0 {
        i--
        buf[i] = byte('0' + n%10)
        n /= 10
    }
    return string(buf[i:])
}
