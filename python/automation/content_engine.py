from __future__ import annotations

import asyncio
import json
import random
import time
from dataclasses import dataclass
from typing import List, Dict, Optional, Any, Protocol


class LLMClient(Protocol):
    async def generate(self, prompt: str, **kwargs: Any) -> str: ...


class ContentDB(Protocol):
    async def save_generated(self, content: "GeneratedContent") -> None: ...


@dataclass
class GeneratedContent:
    id: str
    title: str
    body: str
    metadata: Dict[str, Any]
    score: float


class ContentEngine:
    def __init__(self, llm_client: LLMClient, db: ContentDB) -> None:
        self.llm = llm_client
        self.db = db

    async def generate_content(
        self,
        *,
        topic: str,
        trend_data: Dict[str, Any],
        style: str = "professional",
        max_retries: int = 3,
    ) -> GeneratedContent:
        if not topic:
            raise ValueError("topic required")
        prompt = self._build_prompt(topic, trend_data, style)
        body = await self._retry(lambda: self.llm.generate(prompt, temperature=0.7))
        title = self._extract_title(body) or self._fallback_title(topic)
        score = float(trend_data.get("trend_score", 50.0))
        gen = GeneratedContent(
            id=self._id(),
            title=title,
            body=body,
            metadata={"topic": topic, "style": style, "trend": trend_data},
            score=score,
        )
        await self.db.save_generated(gen)
        return gen

    async def generate_best_content(
        self,
        *,
        topic: str,
        trend_data: Dict[str, Any],
        related_topics: Optional[List[str]] = None,
        style: str = "professional",
        variants: int = 3,
    ) -> GeneratedContent:
        """Generate multiple variants and select best by composite score.

        - Uses readability, simple SEO proxy, and engagement heuristics.
        - Falls back gracefully if LLM calls fail.
        """
        related_topics = related_topics or []
        tasks = [self._generate_variant(topic, related_topics, style) for _ in range(max(1, variants))]
        gens = await asyncio.gather(*tasks, return_exceptions=True)
        candidates: List[Dict[str, Any]] = [g for g in gens if isinstance(g, dict)]
        if not candidates:
            # Fallback: single basic generation
            single = await self.generate_content(topic=topic, trend_data=trend_data, style=style)
            return single
        scored = await self._score_variants(candidates, trend_data)
        best = max(scored, key=lambda x: x.score)
        await self.db.save_generated(best)
        return best

    async def _generate_variant(self, topic: str, related: List[str], style: str) -> Dict[str, Any]:
        prompt = self._build_prompt(topic, {"related": related}, style)
        text = await self._retry(lambda: self.llm.generate(prompt, temperature=0.8))
        title = self._extract_title(text) or self._fallback_title(topic)
        return {
            "id": self._id(),
            "title": title,
            "body": text,
            "metadata": {"topic": topic, "style": style, "related": related},
        }

    async def _score_variants(self, variants: List[Dict[str, Any]], trend_data: Dict[str, Any]) -> List[GeneratedContent]:
        out: List[GeneratedContent] = []
        for v in variants:
            body = str(v.get("body", ""))
            title = str(v.get("title", ""))
            readability = self._calculate_readability(body)
            seo_score = await self._calculate_seo_score(v, trend_data)
            engagement = self._predict_engagement(title, body)
            total = 0.3 * readability + 0.4 * seo_score + 0.3 * engagement
            out.append(GeneratedContent(id=v.get("id", self._id()), title=title, body=body, metadata=v.get("metadata", {}), score=float(total)))
        return out

    def _calculate_readability(self, text: str) -> float:
        # Basic Flesch-like proxy: shorter sentences/words score higher.
        if not text:
            return 0.5
        words = max(1, len(text.split()))
        sentences = max(1, text.count(".") + text.count("! ") + text.count("? "))
        avg_words_per_sentence = words / sentences
        score = 1.0 / (1.0 + (avg_words_per_sentence - 12) ** 2 / 144.0)
        return float(max(0.0, min(1.0, score)))

    async def _calculate_seo_score(self, variant: Dict[str, Any], trend: Dict[str, Any]) -> float:
        # Simple proxy: include topic keywords, reasonable length title/body
        title = str(variant.get("title", "")).lower()
        body = str(variant.get("body", "")).lower()
        topic = str((variant.get("metadata", {}) or {}).get("topic", "")).lower()
        hits = 0
        if topic and topic in title:
            hits += 1
        if topic and topic in body:
            hits += 1
        # Prefer bodies between 250-1200 words
        wc = len(body.split())
        length_bonus = 1.0 if 250 <= wc <= 1200 else 0.6
        volume = float(trend.get("volume", 0))
        vol_norm = min(1.0, volume / 5000.0)
        return float(min(1.0, 0.4 * hits + 0.3 * length_bonus + 0.3 * vol_norm))

    def _predict_engagement(self, title: str, body: str) -> float:
        # Heuristic: title punch + body structure
        emojis = sum(1 for ch in title if ch in "🔥🚀✅⭐⚡✨")
        questions = title.count("?")
        headers = body.count("\n#") + body.count("\n##")
        score = min(1.0, 0.2 * emojis + 0.2 * questions + 0.6 * min(1.0, headers / 10.0))
        return float(score)

    async def _retry(self, fn, *, attempts: int = 3, base: float = 0.5) -> str:
        for i in range(attempts):
            try:
                return await asyncio.wait_for(fn(), timeout=30.0)
            except Exception:
                if i == attempts - 1:
                    raise
                await asyncio.sleep(base * (2 ** i) + random.random() * 0.2)
        raise RuntimeError("unreachable")

    def _build_prompt(self, topic: str, trend: Dict[str, Any], style: str) -> str:
        return (
            f"Write a {style} article about '{topic}'.\n"
            f"Incorporate current interest: trend_score={trend.get('trend_score', 'n/a')}, "
            f"volume={trend.get('volume', 'n/a')}, competition={trend.get('competition', 'n/a')}.\n"
            f"Use clear headings, actionable steps, and end with a concise CTA."
        )

    def _extract_title(self, text: str) -> Optional[str]:
        lines = [l.strip() for l in (text or "").splitlines() if l.strip()]
        return lines[0] if lines else None

    def _fallback_title(self, topic: str) -> str:
        return f"{topic.title()}: A Practical Guide"

    def _id(self) -> str:
        import uuid
        return f"gen_{uuid.uuid4()}"


class InstagramContentGenerator:
    def __init__(self, ig_client, trend_analyzer):
        self.ig = ig_client
        self.analyzer = trend_analyzer

    async def generate_from_hashtags(self, hashtags, *, min_likes=200, min_comments=5, limit_per_hashtag=5):
        """Fetch trending posts by hashtags and return a simplified list for planning."""
        from python.instagram.trending import InstagramTrending
        trending = InstagramTrending(self.ig.access_token, getattr(self.ig, "ig_user_id", None))
        posts = await trending.find_trending_posts(hashtags, min_likes=min_likes, min_comments=min_comments, limit_per_hashtag=limit_per_hashtag)
        # Return top captions with permalinks
        out = []
        for p in posts:
            out.append({
                "id": p.get("id"),
                "caption": (p.get("caption") or "").strip(),
                "permalink": p.get("permalink"),
                "like_count": p.get("like_count", 0),
                "comments_count": p.get("comments_count", 0),
            })
        return out
