# Architecture Overview

- Go monorepo with multiple commands under `cmd/` and packages under `pkg/`.
- External integrations in `pkg/external/*`.
- Database migrations under `_data/_db/**` (recursive).
- Cache layer supports file cache and Valkey.
- Telemetry via OpenTelemetry behind a build tag (`otel`).
- Optional Pulsar integration behind a build tag (`pulsar`).

Key services
- Scheduler: claims due posts and hands off to poster.
- Poster: publishes posts and records outcomes.
- API: FastAPI‑style Gin server for auth, profiles, ICP, etc.

