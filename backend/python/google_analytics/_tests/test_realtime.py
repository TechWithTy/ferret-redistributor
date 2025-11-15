from __future__ import annotations

from unittest.mock import patch, Mock

from ..client import GA4Client
from ..api.utils import build_realtime_request


def test_realtime_success():
    client = GA4Client(access_token="test")

    fake_json = {
        "dimensionHeaders": [{"name": "eventName"}],
        "metricHeaders": [{"name": "activeUsers"}],
        "rows": [
            {
                "dimensionValues": [{"value": "page_view"}],
                "metricValues": [{"value": "5"}],
            }
        ],
    }

    with patch.object(client._session, "request", return_value=Mock(status_code=200, json=lambda: fake_json, content=b"{}")):
        req = build_realtime_request(dimensions=["eventName"], metrics=["activeUsers"], limit=10)
        resp = client.run_realtime_report("123", req)
        assert len(resp.rows) == 1
