from __future__ import annotations

import json
import os
from datetime import datetime
from typing import Any, Dict, List, Optional

from python.calendar.db import session_scope
from python.calendar.models import ScheduledPost, ScheduledStatus, Platform


class FileExperimentDB:
    """Simple file-backed experiment store (JSONL) until a DB table is added."""

    def __init__(self, path: str = "_data/experiments.jsonl") -> None:
        self.path = path
        os.makedirs(os.path.dirname(self.path), exist_ok=True)

    async def save_experiment(self, experiment: Any) -> None:  # matches ExperimentDB Protocol
        with open(self.path, "a", encoding="utf-8") as f:
            f.write(json.dumps(getattr(experiment, "__dict__", dict(experiment))) + "\n")


class SQLAlchemyCalendarDB:
    """DB adapter implementing schedule and basic fetches for Python scheduler."""

    def __init__(self, db_url: Optional[str] = None) -> None:
        self.db_url = db_url

    async def schedule_post(self, payload: Dict[str, Any]) -> None:
        # payload: content_id, variant_id, publish_time (ISO or datetime), platforms (list), experiment_id
        from uuid import uuid4
        publish_time = payload.get("publish_time")
        if isinstance(publish_time, str):
            from datetime import datetime
            try:
                when = datetime.fromisoformat(publish_time)
            except Exception:
                when = datetime.utcnow()
        elif isinstance(publish_time, datetime):
            when = publish_time
        else:
            when = datetime.utcnow()

        content_id = payload.get("content_id")
        caption = payload.get("caption", "")
        hashtags = payload.get("hashtags", "")
        platforms = payload.get("platforms") or ["linkedin"]
        meta = {k: v for k, v in payload.items() if k not in {"content_id", "publish_time", "platforms", "caption", "hashtags"}}

        with session_scope(self.db_url) as s:
            for p in platforms:
                sp = ScheduledPost(
                    id=str(uuid4()),
                    campaign_id=meta.get("campaign_id"),
                    content_id=content_id,
                    platform=Platform(p),
                    caption=caption,
                    hashtags=hashtags,
                    scheduled_at=when,
                    status=ScheduledStatus.scheduled,
                    metadata=meta,
                )
                s.add(sp)

    async def get_scheduled_posts(self, start_time: datetime, end_time: datetime) -> List[Dict[str, Any]]:
        from python.calendar.models import ScheduledPost as SP
        with session_scope(self.db_url) as s:
            q = (
                s.query(SP)
                .filter(SP.status == ScheduledStatus.scheduled)
                .filter(SP.scheduled_at >= start_time)
                .filter(SP.scheduled_at < end_time)
                .order_by(SP.scheduled_at.asc())
            )
            rows = []
            for sp in q.all():
                rows.append({
                    "id": sp.id,
                    "campaign_id": sp.campaign_id,
                    "content_id": sp.content_id,
                    "platform": sp.platform.value if sp.platform else None,
                    "caption": sp.caption,
                    "hashtags": sp.hashtags,
                    "scheduled_at": sp.scheduled_at,
                    "status": sp.status.value if sp.status else None,
                    "metadata": sp.metadata or {},
                })
            return rows

    async def get_upcoming_posts(self, limit: int = 50) -> List[Dict[str, Any]]:
        from python.calendar.models import ScheduledPost as SP
        with session_scope(self.db_url) as s:
            q = (
                s.query(SP)
                .filter(SP.status == ScheduledStatus.scheduled)
                .order_by(SP.scheduled_at.asc())
                .limit(limit)
            )
            rows = []
            for sp in q.all():
                rows.append({
                    "id": sp.id,
                    "scheduled_at": sp.scheduled_at,
                    "platform": sp.platform.value if sp.platform else None,
                    "metadata": sp.metadata or {},
                })
            return rows

    async def reschedule_post(self, post_id: str, new_time: datetime) -> None:
        from python.calendar.models import ScheduledPost as SP
        with session_scope(self.db_url) as s:
            sp = s.query(SP).get(post_id)
            if sp:
                sp.scheduled_at = new_time

