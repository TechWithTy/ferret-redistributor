from __future__ import annotations

from typing import Optional
from fastapi import Header, HTTPException, status

from ..client import GoogleTrendsClient
from ..config import HttpSettings


def get_client(authorization: Optional[str] = Header(None)) -> GoogleTrendsClient:
    """Return a GoogleTrendsClient.

    Alpha note: real API auth model is unknown. This dependency accepts an
    optional Bearer token header to keep parity with other SDKs. The token is
    not used yet but reserved for future integration.
    """
    # If needed later, parse and store token in client state.
    # For now, just construct the client with default http settings.
    http = HttpSettings()
    return GoogleTrendsClient(user_agent=http.user_agent, timeout=http.timeout)


class ApiKeyProvider:
    """Factory for API-key based clients (if the alpha requires an API key)."""

    def __init__(self, api_key: Optional[str] = None) -> None:
        self.api_key = api_key

    def __call__(self) -> GoogleTrendsClient:
        # Placeholder: wire api_key when real HTTP calls are implemented
        return GoogleTrendsClient()
