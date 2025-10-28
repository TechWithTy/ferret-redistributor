from __future__ import annotations

from app.core.landing_page.google_trends.api.deps import get_client
from app.core.landing_page.google_trends import GoogleTrendsClient


def test_get_client_returns_instance():
    client = get_client(authorization=None)
    assert isinstance(client, GoogleTrendsClient)
