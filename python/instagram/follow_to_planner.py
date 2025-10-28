import json
import os
from pathlib import Path
from typing import Any, Dict, List


def load_follow_data(path: str) -> Dict[str, Any]:
    if not os.path.exists(path):
        return {"following": {}, "daily_counts": {"follows": 0, "unfollows": 0, "last_reset": ""}}
    with open(path, "r", encoding="utf-8") as f:
        try:
            return json.load(f)
        except Exception:
            return {"following": {}, "daily_counts": {"follows": 0, "unfollows": 0, "last_reset": ""}}


def merge_trends(existing: List[Dict[str, Any]], extra: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    topics = {t["topic"] for t in existing if isinstance(t, dict) and "topic" in t}
    out = list(existing)
    for t in extra:
        if t.get("topic") not in topics:
            out.append(t)
    return out


def merge_variants(existing: Dict[str, List[Dict[str, Any]]], extra: Dict[str, List[Dict[str, Any]]]) -> Dict[str, List[Dict[str, Any]]]:
    out = dict(existing)
    for k, v in extra.items():
        if k in out:
            # append non-duplicate ids
            seen = {i.get("id") for i in out[k]}
            out[k].extend([i for i in v if i.get("id") not in seen])
        else:
            out[k] = v
    return out


def main() -> int:
    follow_path = "data/follow_data.json"
    trends_path = "_data/trends.json"
    variants_path = "_data/variants.json"
    Path("_data").mkdir(parents=True, exist_ok=True)

    # Load follow data
    fd = load_follow_data(follow_path)
    following = fd.get("following", {})
    count = len(following)

    # Build a simple trend and variants for follow-up content
    topic = "instagram-follow-ups"
    trend = {
        "topic": topic,
        "relevance": 0.7,
        "volume": count,
        "competition": 0.5,
        "trend_score": min(100.0, max(0.0, count * 2.5)),
    }
    variants = {
        topic: [
            {"id": "ifu_v1", "title": "Checking in with our new followers this week", "cta": "DM us if you need the guide", "is_control": True},
            {"id": "ifu_v2", "title": "Welcome to the community!", "cta": "Comment READY for templates", "is_control": False},
        ]
    }

    # Merge with existing artifacts if present
    trends: List[Dict[str, Any]] = []
    if os.path.exists(trends_path):
        try:
            with open(trends_path, "r", encoding="utf-8") as f:
                trends = json.load(f)
                if not isinstance(trends, list):
                    trends = []
        except Exception:
            trends = []
    trends = merge_trends(trends, [trend])
    with open(trends_path, "w", encoding="utf-8") as f:
        json.dump(trends, f, indent=2, ensure_ascii=False)

    vmap: Dict[str, List[Dict[str, Any]]] = {}
    if os.path.exists(variants_path):
        try:
            with open(variants_path, "r", encoding="utf-8") as f:
                vmap = json.load(f)
                if not isinstance(vmap, dict):
                    vmap = {}
        except Exception:
            vmap = {}
    vmap = merge_variants(vmap, variants)
    with open(variants_path, "w", encoding="utf-8") as f:
        json.dump(vmap, f, indent=2, ensure_ascii=False)

    print(f"Wrote/merged planner artifacts with topic '{topic}' (volume={count}).")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

