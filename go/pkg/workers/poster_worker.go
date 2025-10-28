package workers

import (
    "context"
    "encoding/json"
    "time"

    "github.com/bitesinbyte/ferret/pkg/calendar"
    "github.com/bitesinbyte/ferret/pkg/config"
    "github.com/bitesinbyte/ferret/pkg/external"
    "github.com/bitesinbyte/ferret/pkg/factory"
)

type PosterWorker struct {
    DB        DB
    Now       func() time.Time
}

type DB interface {
    UpdateStatus(ctx context.Context, id string, status calendar.ScheduledStatus, externalID *string, publishedAt *time.Time, metadata json.RawMessage) error
}

func (w *PosterWorker) Post(ctx context.Context, row calendar.ScheduledPostRow, cfg config.Config) error {
    poster := factory.CreateSocialPoster(string(row.Platform))
    title := firstNonEmpty(row.ContentTitle.String, row.CampaignName)
    post := external.Post{Title: title, Link: row.ContentURL.String, Description: "", HashTags: row.Hashtags.String}
    publishedAt := w.clockNow()
    if pwid, ok := poster.(external.PosterWithID); ok {
        id, err := pwid.PostWithID(cfg, post)
        if err != nil { return err }
        return w.DB.UpdateStatus(ctx, row.ID, calendar.StatusPublished, &id, &publishedAt, nil)
    }
    if err := poster.Post(cfg, post); err != nil { return err }
    return w.DB.UpdateStatus(ctx, row.ID, calendar.StatusPublished, nil, &publishedAt, nil)
}

func (w *PosterWorker) clockNow() time.Time {
    if w.Now != nil { return w.Now() }
    return time.Now().UTC()
}

func firstNonEmpty(vals ...string) string {
    for _, v := range vals {
        if len(v) > 0 { return v }
    }
    return ""
}

