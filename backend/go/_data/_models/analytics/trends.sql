-- Trend and analytics time-series for accounts, campaigns, and hashtags

BEGIN;

-- Generic metric series table (rollups by bucket)
CREATE TABLE IF NOT EXISTS trend_metrics (
  id            BIGSERIAL PRIMARY KEY,
  org_id        TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  source        TEXT NOT NULL, -- instagram,user_insights | linkedin,page_insights | youtube,video_stats | etc
  dimension     TEXT NOT NULL, -- account:{id} | campaign:{id} | hashtag:{tag}
  metric        TEXT NOT NULL, -- impressions, reach, likes, comments, shares, clicks, views, ctr, etc
  bucket_start  TIMESTAMPTZ NOT NULL,
  bucket_end    TIMESTAMPTZ NOT NULL,
  value         DOUBLE PRECISION NOT NULL,
  meta          JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id, source, dimension, metric, bucket_start, bucket_end)
);

-- Helpful index for range queries
CREATE INDEX IF NOT EXISTS idx_trend_metrics_lookup
  ON trend_metrics(org_id, source, dimension, metric, bucket_start);

COMMIT;

