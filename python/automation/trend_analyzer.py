from __future__ import annotations

import asyncio
import math
import time
from dataclasses import dataclass
from typing import Any, Dict, List, Optional, Protocol

from python.cache.file_cache import FileCache


class SocialAPI(Protocol):
    async def get_topic_metrics(self, topic: str) -> Dict[str, Any]:
        ...


class SEOClient(Protocol):
    async def get_keyword_data(self, topic: str) -> Dict[str, Any]:
        ...


@dataclass
class TrendScore:
    topic: str
    relevance: float  # 0-1
    volume: int
    competition: float  # 0-1
    trend_score: float  # 0-100


class TrendAnalyzer:
    """Concurrent, cached trend analyzer with robust fallbacks.

    - Deduplicates/normalizes topics
    - Bounded concurrency to respect upstream rate limits
    - Per-call timeouts and safe defaults
    - Short‑lived cache to reduce cost/latency
    """

    def __init__(
        self,
        social_api: SocialAPI,
        seo_client: SEOClient,
        *,
        cache_dir: str = "_data_cache",
        concurrency: int = 6,
        ttl_seconds: int = 3600,
        call_timeout: float = 8.0,
    ) -> None:
        self.social_api = social_api
        self.seo = seo_client
        self.cache = FileCache(cache_dir)
        self.sem = asyncio.Semaphore(max(1, concurrency))
        self.ttl = int(max(60, ttl_seconds))
        self.call_timeout = max(1.0, float(call_timeout))

    async def analyze_trends(self, topics: List[str]) -> List[TrendScore]:
        topics = [t.strip().lower() for t in topics if t and t.strip()]
        if not topics:
            return []
        # Deduplicate while preserving order
        seen = set()
        uniq: List[str] = []
        for t in topics:
            if t not in seen:
                seen.add(t)
                uniq.append(t)

        tasks = [self._score_topic(t) for t in uniq]
        results = await asyncio.gather(*tasks, return_exceptions=True)
        out: List[TrendScore] = []
        for r in results:
            if isinstance(r, TrendScore):
                out.append(r)
        # Sort by score desc
        out.sort(key=lambda s: s.trend_score, reverse=True)
        return out

    async def _score_topic(self, topic: str) -> TrendScore:
        # Cache key
        ck = f"trend:{topic}"
        cached = self.cache.get(ck)
        if cached:
            try:
                return TrendScore(**cached)
            except Exception:
                pass

        async with self.sem:
            social_data, seo_data = await asyncio.gather(
                self._timeout(self.social_api.get_topic_metrics(topic)),
                self._timeout(self.seo.get_keyword_data(topic)),
            )

        # Defaults on failure
        social_data = social_data or {}
        seo_data = seo_data or {}

        relevance = self._calculate_relevance(topic, social_data)
        competition = float(self._clamp(seo_data.get("competition", 0.5), 0.0, 1.0))
        volume = int(max(0, seo_data.get("search_volume", 0)))

        trend_score = (
            (relevance * 0.4)
            + (self._normalize_volume(volume) * 0.3)
            + ((1.0 - competition) * 0.3)
        ) * 100.0
        trend_score = float(self._clamp(trend_score, 0.0, 100.0))

        ts = TrendScore(
            topic=topic,
            relevance=relevance,
            volume=volume,
            competition=competition,
            trend_score=trend_score,
        )
        self.cache.set(ck, ts.__dict__, ttl_seconds=self.ttl)
        return ts

    async def get_related_topics(self, topic: str) -> List[str]:
        """Fetch related topics/entities with caching and fallbacks.

        Default implementation returns an empty list; plug in your own logic by
        extending TrendAnalyzer or injecting a richer SocialAPI/SEOClient.
        """
        ck = f"related:{topic}"
        cached = self.cache.get(ck)
        if cached and isinstance(cached, list):
            return [str(x) for x in cached]
        try:
            async with self.sem:
                social = await self._timeout(self.social_api.get_topic_metrics(topic)) or {}
                seo = await self._timeout(self.seo.get_keyword_data(topic)) or {}
            rel = list({*map(str, social.get("related", []) or []), *map(str, seo.get("related", []) or [])})
            rel = [r.strip() for r in rel if r and isinstance(r, str)]
            self.cache.set(ck, rel, ttl_seconds=self.ttl)
            return rel
        except Exception:
            return []

    async def get_optimal_posting_time(self, *, content_type: str, target_audience: str) -> str:
        """Return an ISO timestamp (UTC) for a reasonable posting time.

        This is a heuristic placeholder that can later use historical analytics.
        For now, schedule within the next hour at the next half-hour mark.
        """
        now = int(time.time())
        # next half-hour boundary
        mins = ((now // 60) % 60)
        add = 30 - (mins % 30)
        ts = now + add * 60
        from datetime import datetime, timezone
        return datetime.fromtimestamp(ts, tz=timezone.utc).isoformat()

    async def _timeout(self, coro: asyncio.Future) -> Optional[Dict[str, Any]]:
        try:
            return await asyncio.wait_for(coro, timeout=self.call_timeout)
        except Exception:
            return None

    def _calculate_relevance(self, topic: str, social_data: Dict[str, Any]) -> float:
        # Heuristic: engagement weighted by recency
        eng = float(social_data.get("engagement", 0.0))  # arbitrary scale
        recency_hours = float(social_data.get("recency_hours", 48.0))
        # Normalize engagement via tanh; recency decay
        eng_norm = math.tanh(eng / 1000.0)
        decay = math.exp(-max(0.0, recency_hours) / 72.0)
        rel = eng_norm * decay
        return float(self._clamp(rel, 0.0, 1.0))

    def _normalize_volume(self, volume: int) -> float:
        # Smoothly saturating normalization (0..1)
        return float(self._clamp(math.tanh(volume / 2000.0), 0.0, 1.0))

    @staticmethod
    def _clamp(v: float, lo: float, hi: float) -> float:
        return min(hi, max(lo, v))
