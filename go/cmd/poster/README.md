# Poster CLI

Posts claimed items with retries and per-platform rate limiting; updates DB status.

## Usage
```
DATABASE_URL=postgres://... go run ./cmd/poster --input go/due_posts.json --config config.json
```

## Rate Limiting
Set ops/sec per platform via env (defaults to 2/sec shared):
```
POSTER_RATE_LINKEDIN=1
POSTER_RATE_INSTAGRAM=1
POSTER_RATE_TWITTER=3
POSTER_RATE_FACEBOOK=1
POSTER_RATE_THREAD=1
POSTER_RATE_MASTODON=3
```

## Behavior
- Retries transient errors (429/timeout/temporary) with exponential backoff.
- Captures external IDs when supported and marks `published` with timestamps.
- On error, marks `failed` with error metadata.
- Logs metrics snapshot at the end.

