from __future__ import annotations

import asyncio
from dataclasses import dataclass
from typing import Any, List, Protocol


class TopicsProvider(Protocol):
    async def get_trending_topics(self) -> List[str]: ...
    async def get_industry_news(self) -> List[str]: ...
    async def get_competitor_content(self) -> List[str]: ...


@dataclass
class GrowthEngine:
    db: Any
    analyzer: TopicsProvider
    content_engine: Any
    experiment_engine: Any
    is_running: bool = False

    async def start(self) -> None:
        self.is_running = True
        backoff = 60
        while self.is_running:
            try:
                await self._run_cycle()
                await asyncio.sleep(3600)
                backoff = 60
            except Exception as e:  # pragma: no cover - operational guard
                print(f"Growth engine error: {e}")
                await asyncio.sleep(backoff)
                backoff = min(backoff * 2, 900)

    async def _run_cycle(self) -> None:
        trends = await self._analyze_trends()
        if not trends:
            return
        # Generate for top topics
        top = trends[:3]
        for topic in top:
            try:
                related = []
                if hasattr(self.content_engine, "trend_analyzer") and hasattr(self.content_engine.trend_analyzer, "get_related_topics"):
                    related = await self.content_engine.trend_analyzer.get_related_topics(topic)
                content = await self.content_engine.generate_best_content(topic=topic, trend_data={"trend_score": 50}, related_topics=related)
                exp = await self.experiment_engine.create_experiment(
                    content_id=content.id,
                    experiment_type="title",
                    base_content={"title": content.title},
                    num_variants=3,
                )
                await self._schedule_content(content, exp)
            except Exception:
                continue

    async def _analyze_trends(self) -> List[str]:
        results = await asyncio.gather(
            self.analyzer.get_trending_topics(),
            self.analyzer.get_industry_news(),
            self.analyzer.get_competitor_content(),
            return_exceptions=True,
        )
        topics: List[str] = []
        for r in results:
            if isinstance(r, list):
                topics.extend([t for t in r if isinstance(t, str)])
        # Deduplicate while preserving order
        seen = set()
        uniq: List[str] = []
        for t in [x.strip() for x in topics if x and x.strip()]:
            if t not in seen:
                seen.add(t)
                uniq.append(t)
        return uniq

    async def _schedule_content(self, content: Any, experiment: Any) -> None:
        # Placeholder: write to calendar via DB layer
        pass

