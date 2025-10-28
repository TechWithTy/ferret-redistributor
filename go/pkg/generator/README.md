# Generator (Planning from Python Artifacts)

Purpose: Consume ML/analytics artifacts produced by Python to plan and schedule posts in Postgres.

- Types: `types.go` (Trend, Variant, PlanInput)
- Loaders: `trends.go`, `variants.go`
- Planner: `planner.go` (PlanAndSchedule)
- Captioning: `captioner.go` (CaptionFor, MakeTags)

## Artifacts
- `_data/trends.json` (array)
```
[
  {"topic":"ai marketing","relevance":0.8,"volume":900,"competition":0.3,"trend_score":82.5}
]
```
- `_data/variants.json` (map topic -> variants)
```
{
  "ai marketing": [
    {"id":"v1","title":"The Ultimate Guide to AI Marketing","cta":"Subscribe for weekly insights","is_control":true},
    {"id":"v2","title":"AI Marketing Hacks You Need","cta":"Try the template (free)","is_control":false}
  ]
}
```

## Planning
```
in := generator.PlanInput{
  Trends: trends,
  Variants: variants,
  Platforms: []string{"linkedin","twitter"},
  StartAt: time.Now().UTC().Add(15*time.Minute),
  Spacing: 2*time.Hour,
  PerDayLimit: 10,
}
repo := calendarrepo.Repository{DB: db}
_ = generator.PlanAndSchedule(ctx, repo, in)
```
- Spreads posts across platforms using StartAt + Spacing.
- Picks control + one non-control variant per topic by default.
- Builds platform-specific captions/hashtags via `CaptionFor`.

## CLI
See `go/cmd/planner` for a ready-made CLI that loads artifacts and writes to `scheduled_posts`.

