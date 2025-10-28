# RBAC

Tables
- `roles`, `permissions`, `role_permissions`, `user_roles`, `team_roles`.

Seed roles
- admin: full org privileges (includes `org.admin`).
- editor: content and schedule management.
- viewer: read-only.

Helpers
- See `pkg/auth/rbac.go` for constants and in-memory expansions.
- See `pkg/auth/repo.go` for Postgres-backed permission checks.

