from __future__ import annotations

import random
import string
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from typing import Dict, List, Optional, Protocol


class ExperimentDB(Protocol):
    async def save_experiment(self, experiment: "Experiment") -> None: ...


@dataclass
class Experiment:
    id: str
    content_id: str
    variants: List[Dict]
    metrics: List[str]
    start_date: str
    duration_days: int
    status: str = "draft"


class ExperimentEngine:
    def __init__(self, db: ExperimentDB, content_analyzer: Optional[object] = None) -> None:
        self.db = db
        self.analyzer = content_analyzer
        self.variant_templates = {
            "title": self._generate_title_variants,
            "cta": self._generate_cta_variants,
            "media": self._generate_media_variants,
        }

    async def create_experiment(
        self,
        *,
        content_id: str,
        experiment_type: str,
        base_content: Dict,
        num_variants: int = 3,
        duration_days: int = 7,
        start: Optional[datetime] = None,
    ) -> Experiment:
        if not content_id:
            raise ValueError("content_id required")
        if experiment_type not in self.variant_templates:
            raise ValueError(f"Unknown experiment type: {experiment_type}")
        if num_variants < 1:
            raise ValueError("num_variants must be >= 1")

        variants = await self._generate_variants(experiment_type, base_content, num_variants)
        start_iso = (start or datetime.now(timezone.utc)).isoformat()
        exp = Experiment(
            id=self._generate_id(),
            content_id=content_id,
            variants=variants,
            metrics=["ctr", "engagement", "conversion"],
            start_date=start_iso,
            duration_days=max(1, duration_days),
        )
        await self.db.save_experiment(exp)
        return exp

    async def _generate_variants(self, exp_type: str, base: Dict, count: int) -> List[Dict]:
        generator = self.variant_templates[exp_type]
        return await generator(base, count)

    async def _generate_title_variants(self, base: Dict, count: int) -> List[Dict]:
        templates = [
            lambda t: f"The Ultimate Guide to {t}",
            lambda t: f"{t}: What Nobody Tells You",
            lambda t: f"10 {t} Hacks You Need to Try",
            lambda t: f"{t} Explained in Simple Terms",
            lambda t: f"How to Master {t} in 2025",
        ]
        variants: List[Dict] = []
        base_title = str(base.get("title", "")).strip()
        # Include control if provided
        if base_title:
            variants.append({"id": self._generate_id(), "title": base_title, "is_control": True})
        pool = templates.copy()
        random.shuffle(pool)
        for tmpl in pool:
            if len(variants) >= count:
                break
            variants.append({"id": self._generate_id(), "title": tmpl(base_title or "Your Topic"), "is_control": False})
        if not variants:  # fallback if no base
            variants.append({"id": self._generate_id(), "title": "Your Topic: A Practical Guide", "is_control": True})
        return variants

    async def _generate_cta_variants(self, base: Dict, count: int) -> List[Dict]:
        ctas = [
            "Subscribe for weekly insights",
            "Try the template (free)",
            "DM me ‘READY’ for the guide",
            "Join 1,000+ builders learning growth",
            "Start your 7‑day free trial",
        ]
        variants: List[Dict] = []
        # control CTA if present
        if base.get("cta"):
            variants.append({"id": self._generate_id(), "cta": str(base["cta"]).strip(), "is_control": True})
        for c in ctas:
            if len(variants) >= count:
                break
            variants.append({"id": self._generate_id(), "cta": c, "is_control": False})
        if not variants:
            variants.append({"id": self._generate_id(), "cta": "Learn more", "is_control": True})
        return variants

    async def _generate_media_variants(self, base: Dict, count: int) -> List[Dict]:
        thumbs = base.get("thumbnails") or []
        variants: List[Dict] = []
        if base.get("thumbnail"):
            variants.append({"id": self._generate_id(), "thumbnail": base["thumbnail"], "is_control": True})
        for t in thumbs:
            if len(variants) >= count:
                break
            variants.append({"id": self._generate_id(), "thumbnail": t, "is_control": False})
        if not variants:
            variants.append({"id": self._generate_id(), "thumbnail": base.get("thumbnail", ""), "is_control": True})
        return variants

    def _generate_id(self, n: int = 12) -> str:
        alphabet = string.ascii_lowercase + string.digits
        return "exp_" + "".join(random.choice(alphabet) for _ in range(n))
