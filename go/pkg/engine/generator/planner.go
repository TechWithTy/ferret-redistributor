package generator

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/bitesinbyte/ferret/pkg/adapters/calendarrepo"
)

// PlanAndSchedule writes scheduled_posts based on trends + variants with spacing per platform.
func PlanAndSchedule(ctx context.Context, repo calendarrepo.Repository, in PlanInput) error {
    // Sort trends by score desc
    trends := append([]Trend(nil), in.Trends...)
    sort.Slice(trends, func(i, j int) bool { return trends[i].TrendScore > trends[j].TrendScore })
    now := in.StartAt
    spacing := in.Spacing
    if spacing <= 0 { spacing = 2 * time.Hour }
    perDay := in.PerDayLimit
    if perDay <= 0 { perDay = 10 }

    // Per-platform next-available times
    nextAt := map[string]time.Time{}
    for _, p := range in.Platforms { nextAt[p] = now }

    var batch []calendarrepo.ScheduleInput
    for _, t := range trends {
        vlist := in.Variants[t.Topic]
        if len(vlist) == 0 { continue }
        // Use control + top variant (if present) as example
        toSchedule := pickTopVariants(vlist, 2)
        for _, v := range toSchedule {
            for _, p := range in.Platforms {
                when := nextAt[p]
                // advance spacing
                nextAt[p] = when.Add(spacing)
                caption, tags := CaptionFor(p, t.Topic, v)
                meta := toJSON(map[string]any{"topic": t.Topic, "variant_id": v.ID, "title": v.Title, "cta": v.CTA, "hashtags": tags})
                cap := ptr(caption)
                var hashPtr *string
                if strings.TrimSpace(tags) != "" { hashPtr = &tags }
                item := calendarrepo.ScheduleInput{
                    ID: newID(),
                    CampaignID: nil,
                    ContentID:  nil,
                    Platform:   p,
                    Caption:    cap,
                    Hashtags:   hashPtr,
                    ScheduledAt: when,
                    MetadataJSON: &meta,
                }
                batch = append(batch, item)
            }
        }
    }
    return repo.BulkSchedule(ctx, batch)
}

func pickTopVariants(vs []Variant, n int) []Variant {
    if n >= len(vs) { return vs }
    // Prefer control + first non-control
    out := make([]Variant, 0, n)
    for _, v := range vs { if v.IsControl { out = append(out, v); break } }
    for _, v := range vs { if !v.IsControl { out = append(out, v); break } }
    for len(out) < n && len(out) < len(vs) { out = append(out, vs[len(out)]) }
    return out
}

func newID() string { var b [8]byte; _, _ = rand.Read(b[:]); return "sp_" + hex.EncodeToString(b[:]) }
func toJSON(m map[string]any) string { b, _ := json.Marshal(m); return string(b) }
func ptr(s string) *string { return &s }
