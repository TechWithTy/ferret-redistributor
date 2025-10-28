import os
import time
import math
import json
from dataclasses import dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple

import requests
import aiohttp


@dataclass
class InstagramPost:
    id: str
    caption: str
    media_url: str
    permalink: str
    timestamp: str
    like_count: int
    comments_count: int
    media_type: str
    media_product_type: str
    engagement: int = 0
    comments: List[Dict[str, Any]] = field(default_factory=list)
    insights: Dict[str, Any] = field(default_factory=dict)
    saved_assets: Dict[str, str] = field(default_factory=dict)


class InstagramClient:
    """
    Lightweight Instagram Graph API client focused on hashtag discovery and virality analysis.

    Requirements:
    - Instagram Business or Creator account connected to a Facebook Page
    - Facebook app with instagram_basic, instagram_manage_insights, pages_show_list, pages_read_engagement
    - Access token with required scopes and IG User (Business Account) ID
    """

    def __init__(
        self,
        access_token: str,
        ig_user_id: Optional[str] = None,
        *,
        version: str = None,
        base_url: Optional[str] = None,
        request_timeout: int = 30,
        max_retries: int = 3,
        retry_backoff: float = 0.8,
    ) -> None:
        self.access_token = access_token
        self.ig_user_id = ig_user_id or os.getenv("IG_USER_ID") or os.getenv("INSTAGRAM_BUSINESS_ACCOUNT_ID")
        if not self.ig_user_id:
            raise ValueError("IG User ID is required (set IG_USER_ID or pass ig_user_id)")
        if not version:
            version = os.getenv("IG_GRAPH_VERSION") or os.getenv("INSTAGRAM_API_VERSION") or "v19.0"
        self.version = version
        self.BASE_URL = base_url or f"https://graph.facebook.com/{self.version}"
        self.session = requests.Session()
        self.timeout = request_timeout
        self.max_retries = max_retries
        self.retry_backoff = retry_backoff

    # ----- low-level request with basic retries -----
    def _make_request(self, endpoint: str, params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        url = f"{self.BASE_URL}/{endpoint.lstrip('/')}"
        params = params.copy() if params else {}
        params["access_token"] = self.access_token
        backoff = self.retry_backoff
        for attempt in range(self.max_retries):
            resp = self.session.get(url, params=params, timeout=self.timeout)
            if 200 <= resp.status_code < 300:
                return resp.json()
            # Retry on 429/5xx
            if resp.status_code in (429, 500, 502, 503, 504) and attempt < self.max_retries - 1:
                time.sleep(backoff)
                backoff = min(backoff * 2, 8.0)
                continue
            # Raise for other errors
            try:
                payload = resp.json()
            except Exception:
                payload = {"error": {"message": resp.text}}
            raise requests.HTTPError(f"HTTP {resp.status_code}: {payload}")
        return {}

    # ----- hashtags & media -----
    def search_hashtag(self, hashtag: str) -> Optional[str]:
        """Search a hashtag and return its ID (or None if not found)."""
        endpoint = "ig_hashtag_search"
        params = {"user_id": self.ig_user_id, "q": hashtag.strip("#")}
        data = self._make_request(endpoint, params)
        items = data.get("data", [])
        return (items[0].get("id") if items else None)

    def get_hashtag_recent_media(
        self,
        hashtag_id: str,
        *,
        limit: int = 50,
        fields: Optional[str] = None,
        min_likes: int = 0,
        min_comments: int = 0,
    ) -> List[Dict[str, Any]]:
        """Fetch recent media for a hashtag with optional engagement filters (paginated)."""
        if not fields:
            fields = (
                "id,caption,media_type,media_url,permalink,comments_count,like_count,timestamp,media_product_type"
            )
        endpoint = f"{hashtag_id}/recent_media"
        params: Dict[str, Any] = {
            "user_id": self.ig_user_id,
            "fields": fields,
            "limit": 100,
        }
        results: List[Dict[str, Any]] = []
        after: Optional[str] = None
        remaining = max(1, limit)
        while remaining > 0:
            if after:
                params["after"] = after
            resp = self._make_request(endpoint, params)
            data = resp.get("data", [])
            for post in data:
                if post.get("like_count", 0) >= min_likes and post.get("comments_count", 0) >= min_comments:
                    results.append(post)
                    remaining -= 1
                    if remaining <= 0:
                        break
            paging = resp.get("paging", {}).get("cursors", {})
            after = paging.get("after")
            if not after:
                break
        return results

    # ----- comments & insights -----
    def get_post_comments(self, post_id: str, *, max_total: int = 1000) -> List[Dict[str, Any]]:
        endpoint = f"{post_id}/comments"
        params = {"fields": "id,text,username,timestamp,like_count", "limit": 100}
        out: List[Dict[str, Any]] = []
        after: Optional[str] = None
        remaining = max_total
        while remaining > 0:
            if after:
                params["after"] = after
            data = self._make_request(endpoint, params)
            items = data.get("data", [])
            out.extend(items)
            remaining -= len(items)
            after = data.get("paging", {}).get("cursors", {}).get("after")
            if not after or not items:
                break
        return out

    def get_post_insights(self, post_id: str, *, metrics: Optional[List[str]] = None) -> Dict[str, Any]:
        # Valid metrics depend on media type; these are common ones.
        if metrics is None:
            metrics = ["engagement", "impressions", "reach", "saved", "video_views"]
        endpoint = f"{post_id}/insights"
        params = {"metric": ",".join(metrics), "period": "lifetime"}
        return self._make_request(endpoint, params)

    # ----- assets -----
    def download_media(self, url: str, save_dir: str, *, filename: Optional[str] = None) -> str:
        """Download media from URL and save to save_dir. Returns saved path."""
        Path(save_dir).mkdir(parents=True, exist_ok=True)
        name = filename or url.split("?")[0].split("/")[-1] or f"asset_{int(time.time())}"
        path = Path(save_dir) / name
        with self.session.get(url, stream=True, timeout=self.timeout) as r:
            r.raise_for_status()
            with open(path, "wb") as f:
                for chunk in r.iter_content(chunk_size=8192):
                    if chunk:
                        f.write(chunk)
        return str(path)

    async def download_media_async(self, url: str, save_path: str) -> Optional[str]:
        """Async media download using aiohttp; returns saved path or None."""
        os.makedirs(os.path.dirname(save_path) or ".", exist_ok=True)
        timeout = aiohttp.ClientTimeout(total=self.timeout)
        async with aiohttp.ClientSession(timeout=timeout) as session:
            async with session.get(url) as resp:
                if resp.status == 200:
                    content = await resp.read()
                    with open(save_path, "wb") as f:
                        f.write(content)
                    return save_path
        return None

    async def save_post_assets(self, post: Dict[str, Any], base_dir: str = "downloads") -> Dict[str, str]:
        """Download all assets and metadata for a post (images/videos + comments/insights)."""
        post_id = post.get("id", "")
        ts = datetime.now().strftime("%Y%m%d_%H%M%S")
        post_dir = os.path.join(base_dir, f"post_{post_id}_{ts}")
        os.makedirs(post_dir, exist_ok=True)

        assets: Dict[str, str] = {}
        # Save raw metadata
        with open(os.path.join(post_dir, "post_metadata.json"), "w", encoding="utf-8") as f:
            json.dump(post, f, indent=2, ensure_ascii=False)

        # Download main media
        media_url = post.get("media_url")
        if media_url:
            ext = "mp4" if post.get("media_type") == "VIDEO" else "jpg"
            media_path = os.path.join(post_dir, f"media.{ext}")
            saved = await self.download_media_async(media_url, media_path)
            if saved:
                assets["media"] = saved

        # Fetch and save comments (sync API call)
        try:
            comments = self.get_post_comments(post_id)
        except Exception:
            comments = []
        with open(os.path.join(post_dir, "comments.json"), "w", encoding="utf-8") as f:
            json.dump(comments, f, indent=2, ensure_ascii=False)

        # Fetch and save insights (sync API call)
        try:
            insights = self.get_post_insights(post_id)
        except Exception:
            insights = {}
        with open(os.path.join(post_dir, "insights.json"), "w", encoding="utf-8") as f:
            json.dump(insights, f, indent=2, ensure_ascii=False)

        return assets

    # ----- helpers for virality analysis -----
    @staticmethod
    def _parse_timestamp(ts: str) -> datetime:
        try:
            return datetime.fromisoformat(ts.replace("Z", "+00:00")).astimezone(timezone.utc)
        except Exception:
            return datetime.now(timezone.utc)

    @staticmethod
    def virality_score(likes: int, comments: int, *, hours_since: float) -> float:
        # Simple engagement-velocity score (tune as needed)
        base = likes + 2 * comments
        hours = max(1.0, hours_since)
        return base / (hours ** 0.7)

    def analyze_hashtag(
        self,
        hashtag: str,
        *,
        limit: int = 50,
        min_likes: int = 500,
        min_comments: int = 5,
        with_comments: bool = False,
        with_insights: bool = False,
        download_assets_dir: Optional[str] = None,
    ) -> List[InstagramPost]:
        """End-to-end: resolve hashtag, fetch recent media, compute virality scores, optional enrichments."""
        hashtag_id = self.search_hashtag(hashtag)
        if not hashtag_id:
            return []
        media = self.get_hashtag_recent_media(
            hashtag_id,
            limit=limit,
            min_likes=min_likes,
            min_comments=min_comments,
        )
        results: List[InstagramPost] = []
        for m in media:
            ts = self._parse_timestamp(m.get("timestamp", ""))
            hours_since = max(1.0, (datetime.now(timezone.utc) - ts).total_seconds() / 3600.0)
            score = self.virality_score(m.get("like_count", 0), m.get("comments_count", 0), hours_since=hours_since)
            post = InstagramPost(
                id=m.get("id", ""),
                caption=m.get("caption", "") or "",
                media_url=m.get("media_url", "") or "",
                permalink=m.get("permalink", "") or "",
                timestamp=m.get("timestamp", "") or ts.isoformat(),
                like_count=int(m.get("like_count", 0)),
                comments_count=int(m.get("comments_count", 0)),
                media_type=m.get("media_type", "") or "",
                media_product_type=m.get("media_product_type", "") or "",
                engagement=int(score),
            )
            # optional enrichments
            if with_comments:
                try:
                    post.comments = self.get_post_comments(post.id)
                except Exception:
                    post.comments = []
            if with_insights:
                try:
                    post.insights = self.get_post_insights(post.id)
                except Exception:
                    post.insights = {}
            if download_assets_dir and post.media_url:
                try:
                    saved = self.download_media(post.media_url, download_assets_dir)
                    post.saved_assets = {"media": saved}
                except Exception:
                    post.saved_assets = {}
            results.append(post)
        # sort by computed engagement score desc
        results.sort(key=lambda p: p.engagement, reverse=True)
        return results

    # ----- follow/unfollow decision stubs (no API actions) -----
    def should_follow(self, profile: Dict[str, Any], settings: Dict[str, Any]) -> bool:
        """Decision helper: whether a profile is a good follow candidate.

        Note: IG Graph API does not allow programmatic follow. This is only decisioning.
        """
        min_posts = int(settings.get('min_posts', 3))
        min_ratio = float(settings.get('min_followers_ratio', 0.5))
        max_following = int(settings.get('max_following', 7500))
        # Extract signals (caller supplies profile fields)
        posts = int(profile.get('media_count', 0))
        followers = int(profile.get('followers_count', 0))
        following = int(profile.get('follows_count', 1) or 1)
        if following <= 0:
            following = 1
        ratio = followers / float(following)
        if posts < min_posts:
            return False
        if ratio < min_ratio:
            return False
        if following >= max_following:
            return False
        return True

    def unfollow_due(self, followed_at_unix: int, settings: Dict[str, Any], now_unix: Optional[int] = None) -> bool:
        """Decision helper: whether an account followed at time should be considered for unfollow.
        Only decisioning; no API unfollow supported.
        """
        now = now_unix or int(time.time())
        min_age = int(settings.get('min_follow_duration', 24 * 3600))
        return (now - int(followed_at_unix)) >= min_age

    # ----- unsupported actions (follow/unfollow, follower lists) -----
    def follow_user(self, user_id: str) -> bool:
        """Follow a user (NOT SUPPORTED by Instagram Graph API).

        This method is a stub for decision/testing. Always returns False.
        """
        print("follow_user: Instagram Graph API does not support following users programmatically.")
        return False

    def unfollow_user(self, user_id: str) -> bool:
        """Unfollow a user (NOT SUPPORTED by Instagram Graph API).

        This method is a stub for decision/testing. Always returns False.
        """
        print("unfollow_user: Instagram Graph API does not support unfollowing users programmatically.")
        return False

    def get_user_info(self, username: str) -> Dict[str, Any]:
        """Get limited public info for a business/creator via Business Discovery (if enabled).

        Note: Business Discovery requires additional permissions and only works for business/creator accounts.
        Pass a username (not user_id). Returns empty dict on failure.
        """
        try:
            endpoint = f"{self.ig_user_id}"
            fields = (
                f"business_discovery.username({username})"
                "{followers_count,media_count,follows_count,username,website,name,profile_picture_url}"
            )
            data = self._make_request(endpoint, {"fields": fields})
            bd = (data or {}).get("business_discovery") or {}
            return bd
        except Exception as e:
            print(f"Error getting user info via business discovery: {e}")
            return {}

    def get_following(self, user_id: Optional[str] = None) -> List[Dict[str, Any]]:
        """List of following (NOT SUPPORTED by Instagram Graph API). Stub returns []."""
        print("get_following: Instagram Graph API does not expose following lists.")
        return []

    def get_followers(self, user_id: Optional[str] = None) -> List[Dict[str, Any]]:
        """List of followers (NOT SUPPORTED by Instagram Graph API). Stub returns []."""
        print("get_followers: Instagram Graph API does not expose follower lists.")
        return []


if __name__ == "__main__":
    # Quick demo (reads env IG_ACCESS_TOKEN, IG_USER_ID)
    token = os.getenv("IG_ACCESS_TOKEN") or os.getenv("INSTAGRAM_ACCESS_TOKEN")
    if not token:
        raise SystemExit("Set IG_ACCESS_TOKEN in env")
    client = InstagramClient(token)
    tag = os.getenv("IG_TREND_TAG", "startup")
    posts = client.analyze_hashtag(tag, limit=20, min_likes=100, min_comments=5, with_comments=False, with_insights=False)
    out = [p.__dict__ for p in posts]
    Path("_data").mkdir(exist_ok=True)
    Path("_data/trending_instagram.json").write_text(json.dumps(out, indent=2), encoding="utf-8")
    print(f"Saved {_data_path := '_data/trending_instagram.json'} with {len(posts)} posts for #{tag}")
