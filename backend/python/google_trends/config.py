from __future__ import annotations

import os
from dataclasses import dataclass


# Alpha: endpoint and auth method are not final. Keep overridable via env.
TRENDS_API_BASE_URL = os.getenv("TRENDS_API_BASE_URL", "https://trends.googleapis.com/v1/")
TRENDS_API_KEY = os.getenv("TRENDS_API_KEY")  # if API key model is used
TRENDS_SCOPE = os.getenv("TRENDS_SCOPE")  # if OAuth model is used

DEFAULT_TIMEOUT = int(os.getenv("TRENDS_HTTP_TIMEOUT", "30"))
USER_AGENT = os.getenv("TRENDS_USER_AGENT", "CyberOni-GTrends-SDK/0.1")

# Feature flags
ALPHA_MODE = os.getenv("TRENDS_ALPHA_MODE", "true").lower() == "true"
ENABLE_MERGE = os.getenv("TRENDS_ENABLE_MERGE", "true").lower() == "true"
DEFAULT_INTERVAL = os.getenv("TRENDS_DEFAULT_INTERVAL", "daily")


@dataclass
class HttpSettings:
    timeout: int = DEFAULT_TIMEOUT
    user_agent: str = USER_AGENT


@dataclass
class AuthSettings:
    api_key: str | None = TRENDS_API_KEY
    scope: str | None = TRENDS_SCOPE
