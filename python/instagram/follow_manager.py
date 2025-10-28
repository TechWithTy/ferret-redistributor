from __future__ import annotations

import json
import os
import random
from dataclasses import dataclass
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Tuple

from .config import Config
from .client import InstagramClient


@dataclass
class FollowManager:
    """Local follow/unfollow decision manager (no API follow/unfollow calls).

    Tracks intended follows in a JSON file and enforces daily limits and
    time windows for unfollow consideration. All follow/unfollow actions are
    recorded locally; actual API calls are not performed because the IG Graph
    API does not provide endpoints for follow/unfollow.
    """

    config: Config
    client: InstagramClient
    follow_data_file: Path = Path("data/follow_data.json")

    def __post_init__(self) -> None:
        self.follow_data = self._load_follow_data()

    def _load_follow_data(self) -> Dict:
        if self.follow_data_file.exists():
            try:
                with open(self.follow_data_file, "r", encoding="utf-8") as f:
                    data = json.load(f)
                    # basic shape validation
                    data.setdefault("following", {})
                    data.setdefault("daily_counts", {
                        "follows": 0,
                        "unfollows": 0,
                        "last_reset": datetime.now().isoformat(),
                    })
                    return data
            except Exception as e:
                print(f"Error loading follow data: {e}")
        return {
            "following": {},  # user_id(str) -> {timestamp, followed_back}
            "daily_counts": {
                "follows": 0,
                "unfollows": 0,
                "last_reset": datetime.now().isoformat(),
            },
        }

    def _save_follow_data(self) -> None:
        try:
            self.follow_data_file.parent.mkdir(parents=True, exist_ok=True)
            with open(self.follow_data_file, "w", encoding="utf-8") as f:
                json.dump(self.follow_data, f, indent=2, ensure_ascii=False)
        except Exception as e:
            print(f"Error saving follow data: {e}")

    def reset_daily_counts(self) -> None:
        dc = self.follow_data.get("daily_counts", {})
        last_reset_str = dc.get("last_reset")
        try:
            last_reset = datetime.fromisoformat(last_reset_str) if last_reset_str else datetime.now()
        except Exception:
            last_reset = datetime.now()
        if datetime.now().date() > last_reset.date():
            self.follow_data["daily_counts"] = {
                "follows": 0,
                "unfollows": 0,
                "last_reset": datetime.now().isoformat(),
            }
            self._save_follow_data()

    def can_follow_more(self) -> bool:
        self.reset_daily_counts()
        return (self.follow_data["daily_counts"]["follows"] < self.config.FOLLOW_SETTINGS["daily_follows"])

    def can_unfollow_more(self) -> bool:
        self.reset_daily_counts()
        return (self.follow_data["daily_counts"]["unfollows"] < self.config.FOLLOW_SETTINGS["daily_unfollows"])

    def is_good_candidate(self, user_info: Dict) -> bool:
        if not user_info:
            return False
        # Delegate to client decision helper if available
        return self.client.should_follow(user_info, self.config.FOLLOW_SETTINGS)

    def record_follow(self, user_id: str) -> None:
        self.follow_data["following"][str(user_id)] = {
            "timestamp": datetime.now().isoformat(),
            "followed_back": False,
        }
        self.follow_data["daily_counts"]["follows"] += 1
        self._save_follow_data()

    def record_unfollow(self, user_id: str) -> None:
        if str(user_id) in self.follow_data["following"]:
            del self.follow_data["following"][str(user_id)]
        self.follow_data["daily_counts"]["unfollows"] += 1
        self._save_follow_data()

    def follow_user(self, user_id: str) -> bool:
        if not self.can_follow_more():
            return False
        # Enforce max following threshold
        if len(self.follow_data["following"]) >= self.config.FOLLOW_SETTINGS["max_following"]:
            print("Max following limit reached (local threshold)")
            return False
        # No API call: record intention only
        self.record_follow(user_id)
        return True

    def get_users_to_unfollow(self) -> List[Tuple[str, Dict]]:
        now = datetime.now()
        users: List[Tuple[str, Dict]] = []
        min_dur = self.config.FOLLOW_SETTINGS["min_follow_duration"]
        max_dur = self.config.FOLLOW_SETTINGS["max_follow_duration"]
        for uid, data in list(self.follow_data["following"].items()):
            try:
                ts = datetime.fromisoformat(data.get("timestamp", ""))
            except Exception:
                ts = now - timedelta(days=2)
            age = (now - ts).total_seconds()
            if age > max_dur:
                users.append((uid, data))
            elif age > min_dur and random.random() < 0.3:
                users.append((uid, data))
        return users

    def unfollow_user(self, user_id: str) -> bool:
        if not self.can_unfollow_more():
            return False
        # No API call: record intention only
        self.record_unfollow(user_id)
        return True

    def check_follow_back(self, user_id: str) -> bool:
        # No direct follower list available via Graph API; caller may supply signals via user_info.
        # Here we conservatively return the stored flag if present.
        data = self.follow_data["following"].get(str(user_id))
        return bool(data and data.get("followed_back"))


def main() -> int:
    """Housekeeping pass: reset daily counters if needed and list unfollow candidates.

    Usage:
        python -m python.instagram.follow_manager
    Outputs summary JSON to _data/follow_housekeeping.json
    """
    # Load config and client (token/env)
    cfg = Config()
    token = os.getenv("IG_ACCESS_TOKEN") or os.getenv("INSTAGRAM_ACCESS_TOKEN")
    if not token:
        print("Set IG_ACCESS_TOKEN in env")
        return 1
    client = InstagramClient(token)
    fm = FollowManager(cfg, client)

    # Reset daily counters if a new day
    fm.reset_daily_counts()
    candidates = fm.get_users_to_unfollow()
    summary = {
        "daily_counts": fm.follow_data.get("daily_counts", {}),
        "following_total": len(fm.follow_data.get("following", {})),
        "unfollow_candidates": [uid for uid, _ in candidates],
    }
    os.makedirs("_data", exist_ok=True)
    with open("_data/follow_housekeeping.json", "w", encoding="utf-8") as f:
        json.dump(summary, f, indent=2, ensure_ascii=False)
    print(json.dumps(summary, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
