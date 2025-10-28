"""Google Analytics Data API (GA4) SDK client."""
from __future__ import annotations

import json
import time
from typing import Any, Dict, Iterable, Optional

import requests

from .config import BASE_URL, HTTP_SETTINGS, OAUTH_SETTINGS, SCOPE_ANALYTICS_READONLY, FEATURE_USE_BETA_METADATA
from .api._exceptions import ApiError, AuthError, RateLimitError, RetryableError
from .api._requests import RunReportRequest, RealtimeReportRequest
from .api._responses import RunReportResponse, RealtimeReportResponse


class GA4Client:
    """Main client for GA4 Analytics Data API v1."""

    def __init__(
        self,
        access_token: str,
        *,
        timeout: Optional[int] = None,
        user_agent: Optional[str] = None,
    ) -> None:
        self._access_token = access_token
        self._timeout = timeout or HTTP_SETTINGS.timeout
        self._user_agent = user_agent or HTTP_SETTINGS.user_agent
        self._session = requests.Session()
        self._session.headers.update(
            {
                "Authorization": f"Bearer {self._access_token}",
                "Accept": "application/json",
                "Content-Type": "application/json",
                "User-Agent": self._user_agent,
            }
        )

    # ----- Auth helpers -----
    @staticmethod
    def from_installed_app(
        *,
        client_secrets_path: Optional[str] = None,
        scopes: Optional[Iterable[str]] = None,
        token_path: Optional[str] = None,
    ) -> "GA4Client":
        try:
            from google.oauth2.credentials import Credentials  # type: ignore
            from google_auth_oauthlib.flow import InstalledAppFlow  # type: ignore
            from google.auth.transport.requests import Request as GoogleAuthRequest  # type: ignore
        except Exception as e:  # pragma: no cover
            raise AuthError(
                "google-auth and google-auth-oauthlib are required for installed-app flow"
            ) from e

        scopes = list(scopes) if scopes else [SCOPE_ANALYTICS_READONLY]
        client_secrets_path = client_secrets_path or OAUTH_SETTINGS.client_secrets_path
        token_path = token_path or OAUTH_SETTINGS.token_path

        creds: Optional[Credentials] = None
        try:
            creds = Credentials.from_authorized_user_file(token_path, scopes)  # type: ignore[arg-type]
        except Exception:
            creds = None

        if not creds or not creds.valid:
            if creds and getattr(creds, "expired", False) and getattr(creds, "refresh_token", None):
                try:
                    creds.refresh(GoogleAuthRequest())  # type: ignore[arg-type]
                except Exception:
                    creds = None
            if not creds or not creds.valid:
                flow = InstalledAppFlow.from_client_secrets_file(client_secrets_path, scopes)
                creds = flow.run_local_server(port=0)
                try:
                    with open(token_path, "w", encoding="utf-8") as f:
                        f.write(creds.to_json())
                except Exception:
                    pass

        return GA4Client(access_token=creds.token)

    @staticmethod
    def from_service_account(
        *,
        key_path: Optional[str] = None,
        scopes: Optional[Iterable[str]] = None,
    ) -> "GA4Client":
        try:
            from google.oauth2 import service_account  # type: ignore
            from google.auth.transport.requests import Request as GoogleAuthRequest  # type: ignore
        except Exception as e:  # pragma: no cover
            raise AuthError("google-auth is required for service account flow") from e

        scopes = list(scopes) if scopes else [SCOPE_ANALYTICS_READONLY]
        key_path = key_path or OAUTH_SETTINGS.service_account_key_path
        creds = service_account.Credentials.from_service_account_file(key_path, scopes=scopes)  # type: ignore
        # Create fresh token
        auth_request = GoogleAuthRequest()
        try:
            creds.refresh(auth_request)  # type: ignore[arg-type]
        except Exception as e:
            raise AuthError(f"Failed to refresh service account token: {e}") from e
        return GA4Client(access_token=creds.token)

    # ----- Low-level HTTP with retries -----
    def _request(self, method: str, url: str, *, json_body: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        backoff = 1.0
        max_backoff = 16.0
        attempts = 0
        while True:
            attempts += 1
            resp = self._session.request(method, url, json=json_body, timeout=self._timeout)
            if 200 <= resp.status_code < 300:
                if resp.content:
                    return resp.json()
                return {}

            # Parse error
            try:
                err = resp.json()
                message = err.get("error", {}).get("message") or err.get("message") or resp.text
                reason = None
                if "error" in err and isinstance(err["error"], dict):
                    details = err["error"].get("details") or []
                    reason = err["error"].get("status")
                    if not reason:
                        reasons = err["error"].get("errors") or []
                        if reasons and isinstance(reasons, list) and isinstance(reasons[0], dict):
                            reason = reasons[0].get("reason")
            except Exception:
                message = resp.text
                reason = None

            if resp.status_code == 401 or resp.status_code == 403:
                raise AuthError(f"Unauthorized/Forbidden: {message}")
            if resp.status_code == 429:
                if attempts <= 5:
                    sleep_for = float(resp.headers.get("Retry-After", backoff))
                    time.sleep(sleep_for)
                    backoff = min(backoff * 2, max_backoff)
                    continue
                raise RateLimitError(resp.status_code, message, reason)
            if 500 <= resp.status_code < 600:
                if attempts <= 5:
                    time.sleep(backoff)
                    backoff = min(backoff * 2, max_backoff)
                    continue
                raise RetryableError(resp.status_code, message, reason)

            raise ApiError(resp.status_code, message, reason)

    # ----- High-level API methods -----
    def run_report(self, property_id: str | int, req: RunReportRequest) -> RunReportResponse:
        url = f"{BASE_URL}/properties/{property_id}:runReport"
        data = self._request("POST", url, json_body=json.loads(req.model_dump_json(exclude_none=True)))
        return RunReportResponse.model_validate(data)

    def run_realtime_report(self, property_id: str | int, req: RealtimeReportRequest) -> RealtimeReportResponse:
        url = f"{BASE_URL}/properties/{property_id}:runRealtimeReport"
        data = self._request("POST", url, json_body=json.loads(req.model_dump_json(exclude_none=True)))
        return RealtimeReportResponse.model_validate(data)

    def get_metadata(self, property_id: str | int) -> Dict[str, Any]:
        base = "https://analyticsdata.googleapis.com/v1beta" if FEATURE_USE_BETA_METADATA else BASE_URL
        url = f"{base}/properties/{property_id}/metadata"
        return self._request("GET", url)

    # Optional
    def batch_run_reports(self, property_id: str | int, requests_list: list[RunReportRequest]) -> Dict[str, Any]:
        url = f"{BASE_URL}/properties/{property_id}:batchRunReports"
        payload = {"requests": [json.loads(r.model_dump_json(exclude_none=True)) for r in requests_list]}
        return self._request("POST", url, json_body=payload)
