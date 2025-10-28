# Scheduler CLI

Atomically claims due posts (status=`scheduled` -> `processing`) within a window and prints them as JSON.

## Usage
```
DATABASE_URL=postgres://... go run ./cmd/scheduler --within 10m --json-array > go/due_posts.json
```

- Uses `FOR UPDATE SKIP LOCKED` to avoid duplicate claims.
- Metrics are logged at the end for CI visibility.

