-- Support tickets and agent messaging to assist users

BEGIN;

CREATE TABLE IF NOT EXISTS support_queues (
  id           TEXT PRIMARY KEY,
  name         TEXT NOT NULL UNIQUE,
  description  TEXT,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS support_tickets (
  id             TEXT PRIMARY KEY,
  org_id         TEXT REFERENCES organizations(id) ON DELETE SET NULL,
  user_id        TEXT REFERENCES users(id) ON DELETE SET NULL,
  queue_id       TEXT REFERENCES support_queues(id) ON DELETE SET NULL,
  assigned_agent TEXT REFERENCES admin_profiles(id) ON DELETE SET NULL,
  subject        TEXT NOT NULL,
  status         TEXT NOT NULL DEFAULT 'open', -- open, pending, resolved, closed
  priority       TEXT NOT NULL DEFAULT 'normal', -- low, normal, high, urgent
  tags           JSONB NOT NULL DEFAULT '[]'::jsonb,
  metadata       JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  closed_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS support_messages (
  id            TEXT PRIMARY KEY,
  ticket_id     TEXT NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
  author_user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
  author_admin_id TEXT REFERENCES admin_profiles(id) ON DELETE SET NULL,
  author_role   TEXT NOT NULL, -- user | agent | system
  body          TEXT NOT NULL,
  attachments   JSONB NOT NULL DEFAULT '[]'::jsonb,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_support_tickets_org_status ON support_tickets(org_id, status);
CREATE INDEX IF NOT EXISTS idx_support_messages_ticket ON support_messages(ticket_id);

COMMIT;

