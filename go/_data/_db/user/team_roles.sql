-- Team-scoped RBAC assignments

BEGIN;

-- Assign roles to users within a team (inherits org via teams.org_id)
CREATE TABLE IF NOT EXISTS team_roles (
  team_id      TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
  user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id      TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(team_id, user_id, role_id)
);

COMMIT;

