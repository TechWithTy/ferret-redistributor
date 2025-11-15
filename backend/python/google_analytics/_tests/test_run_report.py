from __future__ import annotations

from unittest.mock import patch, Mock

from ..client import GA4Client
from ..api.utils import build_run_report_request


def test_run_report_success():
    client = GA4Client(access_token="test")

    fake_json = {
        "dimensionHeaders": [{"name": "date"}],
        "metricHeaders": [{"name": "activeUsers", "type": "TYPE_INTEGER"}],
        "rows": [
            {
                "dimensionValues": [{"value": "20250101"}],
                "metricValues": [{"value": "123"}],
            }
        ],
        "rowCount": 1,
    }

    with patch.object(client._session, "request", return_value=Mock(status_code=200, json=lambda: fake_json, content=b"{}")):
        req = build_run_report_request(dimensions=["date"], metrics=["activeUsers"], last_n_days=7, limit=10)
        resp = client.run_report("123", req)
        assert resp.rowCount == 1
        assert len(resp.rows) == 1
