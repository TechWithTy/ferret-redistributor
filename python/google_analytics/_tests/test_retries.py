from __future__ import annotations

import types
from unittest.mock import patch, Mock

from ..client import GA4Client


def test_retries_on_429_then_success():
    client = GA4Client(access_token="test")

    class Resp:
        def __init__(self, status_code, headers=None, payload=None):
            self.status_code = status_code
            self.headers = headers or {}
            self._payload = payload or {}
            self.content = b"{}"
        def json(self):
            return self._payload
        @property
        def text(self):
            return json.dumps(self._payload)

    import json

    calls = []

    def side_effect(method, url, json=None, timeout=None):  # type: ignore[override]
        calls.append(url)
        if len(calls) == 1:
            return Resp(429, headers={"Retry-After": "0.01"}, payload={"error": {"message": "rate"}})
        return Resp(200, payload={})

    with patch.object(client._session, "request", side_effect=side_effect):
        client.get_metadata("123")
        assert len(calls) == 2
