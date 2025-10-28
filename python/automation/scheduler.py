from __future__ import annotations

import asyncio
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from typing import Any, Dict, List, Protocol


class SchedulerDB(Protocol):
    async def get_scheduled_posts(self, start_time: datetime, end_time: datetime) -> List[Dict[str, Any]]: ...
    async def get_upcoming_posts(self, limit: int = 50) -> List[Dict[str, Any]]: ...
    async def reschedule_post(self, post_id: str, new_time: datetime) -> None: ...


class Publisher(Protocol):
    async def publish(self, post: Dict[str, Any]) -> None: ...


class Analytics(Protocol):
    async def get_performance_metrics(self, lookback_days: int = 30) -> Dict[str, Any]: ...


@dataclass
class ContentScheduler:
    db: SchedulerDB
    publisher: Publisher
    analytics: Analytics
    interval_seconds: int = 300
    _stop: bool = False

    async def run_schedule_cycle(self) -> None:
        backoff = 10
        while not self._stop:
            try:
                await self._process_scheduled_posts()
                await self._optimize_future_posts()
                await asyncio.sleep(self.interval_seconds)
                backoff = 10
            except Exception as e:  # pragma: no cover - operational guard
                print(f"Scheduler error: {e}")
                await asyncio.sleep(backoff)
                backoff = min(backoff * 2, 300)

    async def _process_scheduled_posts(self) -> None:
        now = datetime.now(timezone.utc)
        window = now + timedelta(minutes=5)
        upcoming = await self.db.get_scheduled_posts(start_time=now, end_time=window)
        # Fan-out publishing tasks with bounded concurrency
        sem = asyncio.Semaphore(8)
        async def _task(p: Dict[str, Any]):
            async with sem:
                try:
                    await self.publisher.publish(p)
                except Exception as e:
                    print(f"Publish error for {p.get('id')}: {e}")
        await asyncio.gather(*[_task(p) for p in upcoming])

    async def _optimize_future_posts(self) -> None:
        perf = await self.analytics.get_performance_metrics(lookback_days=30)
        best_times = self._calculate_optimal_times(perf)
        upcoming = await self.db.get_upcoming_posts(limit=50)
        for post in upcoming:
            try:
                if self._should_reschedule(post, best_times):
                    new_time = self._calculate_best_time(post, best_times)
                    await self.db.reschedule_post(post["id"], new_time)
            except Exception:
                continue

    def _calculate_optimal_times(self, perf_data: Dict[str, Any]) -> Dict[str, List[str]]:
        # Placeholder: derive per-day time slots; fall back to defaults
        return {
            "monday": ["09:00", "13:00", "19:00"],
            "tuesday": ["09:00", "13:00", "19:00"],
            "wednesday": ["09:00", "13:00", "19:00"],
            "thursday": ["09:00", "13:00", "19:00"],
            "friday": ["09:00", "13:00", "19:00"],
            "saturday": ["10:00", "14:00"],
            "sunday": ["10:00", "14:00"],
        }

    def _should_reschedule(self, post: Dict[str, Any], best_times: Dict[str, List[str]]) -> bool:
        # Simple heuristic: if not in a top slot for its weekday, consider rescheduling
        when: datetime = post.get("scheduled_at")
        if not isinstance(when, datetime):
            return False
        weekday = when.strftime("%A").lower()
        slot = when.strftime("%H:%M")
        return slot not in set(best_times.get(weekday, []))

    def _calculate_best_time(self, post: Dict[str, Any], best_times: Dict[str, List[str]]) -> datetime:
        when: datetime = post.get("scheduled_at")
        if not isinstance(when, datetime):
            return datetime.now(timezone.utc) + timedelta(hours=1)
        weekday = when.strftime("%A").lower()
        slots = best_times.get(weekday) or [when.strftime("%H:%M")]
        hh, mm = map(int, slots[0].split(":"))
        return when.replace(hour=hh, minute=mm, second=0, microsecond=0)

