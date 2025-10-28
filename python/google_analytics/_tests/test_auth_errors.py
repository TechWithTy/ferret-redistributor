from __future__ import annotations

from unittest.mock import patch, Mock

from ..client import GA4Client
from ..api._exceptions import AuthError


def test_401_raises_auth_error():
    client = GA4Client(access_token="test")

    fake_resp = Mock(status_code=401, json=lambda: {"error": {"message": "bad"}}, content=b"{}", text="bad")

    with patch.object(client._session, "request", return_value=fake_resp):
        try:
            client.get_metadata("123")
        except AuthError as e:
            assert "Unauthorized" in str(e) or "Forbidden" in str(e)
        else:
            assert False, "Expected AuthError"
