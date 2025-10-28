import os
import json
from pathlib import Path
from typing import Dict, List, Any


def load_trending_posts(path: str) -> List[Dict[str, Any]]:
    if not os.path.exists(path):
        return []
    with open(path, "r", encoding="utf-8") as f:
        try:
            data = json.load(f)
            if isinstance(data, list):
                return data
            return []
        except Exception:
            return []


def compute_trends(posts: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    by_tag: Dict[str, List[Dict[str, Any]]] = {}
    for p in posts:
        tag = (p.get("hashtag") or "").strip().lower()
        if not tag:
            continue
        by_tag.setdefault(tag, []).append(p)

    trends = []
    for tag, items in by_tag.items():
        if not items:
            continue
        # compute average engagement (likes + comments)
        scores = [(p.get("like_count", 0) + p.get("comments_count", 0)) for p in items]
        avg = sum(scores) / max(1, len(scores))
        trend_score = min(100.0, max(0.0, avg / 50.0))  # heuristic scale
        trends.append({
            "topic": tag,
            "relevance": 0.8,
            "volume": len(items),
            "competition": 0.5,
            "trend_score": trend_score,
        })
    return trends


def build_variants(posts: List[Dict[str, Any]]) -> Dict[str, List[Dict[str, Any]]]:
    by_tag: Dict[str, List[Dict[str, Any]]] = {}
    for p in posts:
        tag = (p.get("hashtag") or "").strip().lower()
        if not tag:
            continue
        by_tag.setdefault(tag, []).append(p)

    variants: Dict[str, List[Dict[str, Any]]] = {}
    for tag, items in by_tag.items():
        # sort by engagement
        items.sort(key=lambda x: (x.get("like_count", 0) + x.get("comments_count", 0)), reverse=True)
        top_caption = (items[0].get("caption") or "").strip() if items else ""
        vlist = []
        # control variant: top caption
        vlist.append({"id": f"{tag}_v1", "title": top_caption[:120], "cta": "Follow for weekly insights", "is_control": True})
        # templated variant
        vlist.append({"id": f"{tag}_v2", "title": f"Top {tag} tips this week", "cta": "Try the template (free)", "is_control": False})
        variants[tag] = vlist
    return variants


def main() -> int:
    posts_path = "_data/ig_trending_posts.json"
    trends_path = "_data/trends.json"
    variants_path = "_data/variants.json"
    Path("_data").mkdir(parents=True, exist_ok=True)

    posts = load_trending_posts(posts_path)
    trends = compute_trends(posts)
    variants = build_variants(posts)

    with open(trends_path, "w", encoding="utf-8") as f:
        json.dump(trends, f, indent=2, ensure_ascii=False)
    with open(variants_path, "w", encoding="utf-8") as f:
        json.dump(variants, f, indent=2, ensure_ascii=False)
    print(f"Wrote planner artifacts: {trends_path}, {variants_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

