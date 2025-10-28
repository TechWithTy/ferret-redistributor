# Planner CLI

Reads Python artifacts and schedules posts into Postgres `scheduled_posts`.

## Usage
```
DATABASE_URL=postgres://... go run ./cmd/planner \
  --trends _data/trends.json \
  --variants _data/variants.json \
  --spacing 2h \
  --start-offset 15m
```

- Platforms default: `linkedin`, `twitter` (adjust in code or extend flags).
- Spacing applies per platform to avoid audience fatigue.

## Artifacts
- See `go/pkg/generator/README.md` for JSON shapes.

