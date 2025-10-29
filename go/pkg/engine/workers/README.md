# Workers (Go)

Reusable units .

- `poster_worker.go`: posts a single `ScheduledPostRow` via factory poster and updates DB status.

## Usage
```
w := workers.PosterWorker{DB: yourDBAdapter}
_ = w.Post(ctx, row, cfg)
```
- Calls `PosterWithID.PostWithID` when available (captures external ID), else `Poster.Post`.
- Updates `scheduled_posts` to `published` with timestamps.

