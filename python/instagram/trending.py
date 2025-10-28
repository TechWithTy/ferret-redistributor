import os
import json
from typing import List, Dict, Any, Optional

from .client import InstagramClient


class InstagramTrending:
    def __init__(self, access_token: str, page_id: Optional[str] = None) -> None:
        self.client = InstagramClient(access_token, page_id)

    async def find_trending_posts(
        self,
        hashtags: List[str],
        *,
        min_likes: int = 1000,
        min_comments: int = 10,
        limit_per_hashtag: int = 10,
    ) -> List[Dict[str, Any]]:
        """Return a list of trending posts (dicts) across hashtags, sorted by engagement."""
        all_posts: List[Dict[str, Any]] = []
        for raw_tag in hashtags:
            hashtag = raw_tag.strip().lstrip("#")
            if not hashtag:
                continue
            try:
                hashtag_id = self.client.search_hashtag(hashtag)
                if not hashtag_id:
                    print(f"No results for #{hashtag}")
                    continue
                print(f"Found hashtag ID for #{hashtag}: {hashtag_id}")
                posts = self.client.get_hashtag_recent_media(
                    hashtag_id,
                    limit=limit_per_hashtag,
                    min_likes=min_likes,
                    min_comments=min_comments,
                )
                print(f"Found {len(posts)} trending posts for #{hashtag}")
                # Annotate with source hashtag for downstream conversion
                for p in posts:
                    if isinstance(p, dict):
                        p.setdefault("hashtag", hashtag)
                        all_posts.append(p)
            except Exception as e:
                print(f"Error processing #{hashtag}: {e}")
                continue

        all_posts.sort(
            key=lambda x: (x.get("like_count", 0) + x.get("comments_count", 0)),
            reverse=True,
        )
        return all_posts

    async def find_trending_content(
        self,
        hashtags: List[str],
        *,
        min_engagement: int = 1000,
        max_downloads: int = None,
        download_base: str = "downloads",
    ) -> List[Dict[str, Any]]:
        """Find high-performing posts across hashtags and optionally download assets.

        - min_engagement is likes + comments threshold
        - max_downloads caps the number of asset downloads (from env MAX_DAILY_DOWNLOADS if unset)
        """
        posts = await self.find_trending_posts(hashtags, min_likes=0, min_comments=0, limit_per_hashtag=25)
        # Filter by total engagement
        filtered = [p for p in posts if (p.get("like_count", 0) + p.get("comments_count", 0)) >= min_engagement]

        # Optional asset downloads
        if max_downloads is None:
            try:
                max_downloads = int(os.getenv("MAX_DAILY_DOWNLOADS", "50"))
            except Exception:
                max_downloads = 50
        downloaded = 0
        for i, post in enumerate(filtered):
            if downloaded >= max_downloads:
                break
            try:
                assets = await self.client.save_post_assets(post, base_dir=os.path.join(download_base, f"post_{i+1}"))
                post["saved_assets"] = assets
                downloaded += 1
            except Exception as e:
                print(f"Download error for {post.get('id')}: {e}")
                continue
        return filtered


if __name__ == "__main__":
    token = os.getenv("IG_ACCESS_TOKEN") or os.getenv("INSTAGRAM_ACCESS_TOKEN")
    if not token:
        raise SystemExit("Set IG_ACCESS_TOKEN")
    tags = os.getenv("IG_TREND_TAGS", "startup,marketing,founders").split(",")
    import asyncio

    trending = InstagramTrending(token)
    posts = asyncio.run(trending.find_trending_posts([t.strip() for t in tags if t.strip()], min_likes=200, min_comments=5, limit_per_hashtag=10))
    os.makedirs("_data", exist_ok=True)
    with open("_data/ig_trending_posts.json", "w", encoding="utf-8") as f:
        json.dump(posts, f, indent=2, ensure_ascii=False)
    print(f"Saved {_ := len(posts)} posts to _data/ig_trending_posts.json")
