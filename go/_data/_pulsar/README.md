Pulsar Integration Plan

Topics (public/default unless noted)
- post-events — events emitted after posting
- schedule-events — events for scheduled/claimed/canceled posts

Message Schema (JSON)
- post-events
  {
    "type": "post.published",
    "id": "<scheduled_post_id>",
    "platform": "instagram|linkedin|twitter|...",
    "external_id": "<platform_post_id>",
    "published_at": "RFC3339",
    "campaign_id": "...",
    "content_id": "...",
    "org_id": "..."
  }

Env
- PULSAR_SERVICE_URL=pulsar://localhost:6650 or http://localhost:8080 (web)
- PULSAR_TOKEN= (optional)
- PULSAR_TOPIC_POST_EVENTS=persistent://public/default/post-events

Build
- Default builds use a no-op stub.
- Build with `-tags=pulsar` to enable real client (requires `github.com/apache/pulsar-client-go/pulsar`).

