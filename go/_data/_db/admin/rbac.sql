-- RBAC schema: roles, permissions, bindings (org-scoped)

BEGIN;

CREATE TABLE IF NOT EXISTS roles (
  id           TEXT PRIMARY KEY,
  name         TEXT NOT NULL UNIQUE,
  description  TEXT,
  scope        TEXT NOT NULL DEFAULT 'org', -- org, global
  is_system    BOOLEAN NOT NULL DEFAULT TRUE,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permissions (
  name         TEXT PRIMARY KEY,
  description  TEXT
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id      TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission   TEXT NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
  PRIMARY KEY(role_id, permission)
);

-- Org-scoped user role assignments
CREATE TABLE IF NOT EXISTS user_roles (
  org_id       TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id      TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(org_id, user_id, role_id)
);

-- Seed permissions (MVP)
INSERT INTO permissions(name, description) VALUES
  ('org.admin',          'Organization administration (all privileges)'),
  ('posts.read',         'Read posts and schedules'),
  ('posts.write',        'Create/update/delete posts and schedules'),
  ('schedule.claim',     'Claim scheduled posts for processing'),
  ('schedule.manage',    'Manage scheduling windows and limits'),
  ('billing.read',       'Read billing and plans'),
  ('billing.manage',     'Manage billing and subscriptions'),
  ('analytics.read',     'Read analytics dashboards')
ON CONFLICT DO NOTHING;

-- Seed roles
INSERT INTO roles(id, name, description) VALUES
  ('role_admin',  'admin',  'Full org administration'),
  ('role_editor', 'editor', 'Manage content and scheduling'),
  ('role_viewer', 'viewer', 'Read-only access')
ON CONFLICT DO NOTHING;

-- Map role -> permissions
-- admin: all
INSERT INTO role_permissions(role_id, permission)
SELECT 'role_admin', name FROM permissions
ON CONFLICT DO NOTHING;

-- editor
INSERT INTO role_permissions(role_id, permission) VALUES
  ('role_editor','posts.read'),
  ('role_editor','posts.write'),
  ('role_editor','schedule.claim'),
  ('role_editor','schedule.manage'),
  ('role_editor','analytics.read')
ON CONFLICT DO NOTHING;

-- viewer
INSERT INTO role_permissions(role_id, permission) VALUES
  ('role_viewer','posts.read'),
  ('role_viewer','analytics.read')
ON CONFLICT DO NOTHING;

COMMIT;

