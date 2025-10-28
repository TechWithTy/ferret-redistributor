from __future__ import annotations

from dataclasses import dataclass, field
import os


@dataclass
class Config:
    """Instagram automation configuration (decisioning only, no follow/unfollow API calls).

    Note: Instagram Graph API does not support programmatic follow/unfollow actions.
    These settings are intended for decision logic and scheduling only.
    """

    # Follow/Unfollow Settings
    FOLLOW_SETTINGS: dict = field(default_factory=lambda: {
        'daily_follows': int(os.getenv('IG_DAILY_FOLLOWS', '50')),
        'daily_unfollows': int(os.getenv('IG_DAILY_UNFOLLOWS', '50')),
        'min_follow_duration': int(os.getenv('IG_MIN_FOLLOW_DURATION', str(24 * 3600))),  # 24h
        'max_follow_duration': int(os.getenv('IG_MAX_FOLLOW_DURATION', str(3 * 24 * 3600))),  # 72h
        'follow_chance': float(os.getenv('IG_FOLLOW_CHANCE', '0.7')),
        'unfollow_older_than': int(os.getenv('IG_UNFOLLOW_OLDER_THAN', str(24 * 3600))),
        'max_following': int(os.getenv('IG_MAX_FOLLOWING', '7500')),
        'min_followers_ratio': float(os.getenv('IG_MIN_FOLLOWERS_RATIO', '0.5')),
        'min_posts': int(os.getenv('IG_MIN_POSTS', '3')),
    })

    # Engagement daily targets
    DAILY_TARGETS: dict = field(default_factory=lambda: {
        'likes': int(os.getenv('IG_DAILY_LIKES', '100')),
        'comments': int(os.getenv('IG_DAILY_COMMENTS', '30')),
        'stories': int(os.getenv('IG_DAILY_STORIES', '50')),
        'follows': int(os.getenv('IG_DAILY_FOLLOWS', '50')),
        'unfollows': int(os.getenv('IG_DAILY_UNFOLLOWS', '50')),
    })

