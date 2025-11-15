"""
Google Analytics Data API (GA4) SDK configuration.

Defines constants, default settings, and simple env-driven configuration
for OAuth2 credentials and token storage.
"""
from __future__ import annotations

import os
from dataclasses import dataclass

# REST base
BASE_URL = "https://analyticsdata.googleapis.com/v1"

# Scopes
SCOPE_ANALYTICS_READONLY = "https://www.googleapis.com/auth/analytics.readonly"

# Defaults
DEFAULT_TIMEOUT = 30  # seconds
DEFAULT_USER_AGENT = "CyberOni-GA4-SDK/1.0"

# Feature flag to optionally use v1beta for metadata
FEATURE_USE_BETA_METADATA = os.getenv("GA_USE_BETA_METADATA", "false").lower() in {"1", "true", "yes"}


@dataclass(frozen=True)
class OAuthSettings:
    """Location of client secrets and token cache.

    Env overrides:
      - GA_CLIENT_SECRETS: path to OAuth client_secrets.json
      - GA_TOKEN_PATH: path to token storage file (JSON)
      - GA_SA_KEY: service account JSON key path
    """

    client_secrets_path: str = os.getenv("GA_CLIENT_SECRETS", "./client_secrets.json")
    token_path: str = os.getenv("GA_TOKEN_PATH", "./.ga4_token.json")
    service_account_key_path: str = os.getenv("GA_SA_KEY", "./service_account.key.json")


@dataclass(frozen=True)
class HttpSettings:
    timeout: int = int(os.getenv("GA_HTTP_TIMEOUT", str(DEFAULT_TIMEOUT)))
    user_agent: str = os.getenv("GA_USER_AGENT", DEFAULT_USER_AGENT)


OAUTH_SETTINGS = OAuthSettings()
HTTP_SETTINGS = HttpSettings()
