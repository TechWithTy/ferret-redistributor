from __future__ import annotations

from typing import Optional
from fastapi import Depends, Header, HTTPException, status

from ..client import GA4Client
from ..config import SCOPE_ANALYTICS_READONLY


def get_client(authorization: Optional[str] = Header(None)) -> GA4Client:
    """Create a GA4Client from a Bearer token in the Authorization header.

    Expect header: "Authorization: Bearer <token>"
    """
    if not authorization:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Missing Authorization header")
    parts = authorization.split()
    if len(parts) != 2 or parts[0].lower() != "bearer":
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid Authorization header format")
    token = parts[1]
    if not token:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Empty bearer token")
    return GA4Client(access_token=token)


# Optional utility: service-account client via settings/env (not used as default dependency)
class ServiceAccountProvider:
    def __init__(self, key_path: Optional[str] = None) -> None:
        self.key_path = key_path

    def __call__(self) -> GA4Client:
        return GA4Client.from_service_account(key_path=self.key_path, scopes=[SCOPE_ANALYTICS_READONLY])
