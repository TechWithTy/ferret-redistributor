# Deploy & Run

Prereqs
- Go 1.21+
- Postgres 13+
- (Optional) Valkey/Redis for dedupe and caching
- (Optional) Pulsar for events
- (Optional) OTel collector for telemetry

Steps
1. Configure env in `.env` or shell (template: `backend/env.example`)
2. Apply migrations: `DATABASE_URL=postgres://... make db-migrate`
3. Build poster/scheduler binaries or run via `go run`
4. API: install Gin (network required): `go get github.com/gin-gonic/gin`
   - Run API: `make run-api`
   - For password hashing: `go get golang.org/x/crypto/bcrypt` and build with `-tags=secure`

Security
- Set `JWT_SECRET` to enable JWT auth.
- Use strong passwords and do not store raw tokens.

