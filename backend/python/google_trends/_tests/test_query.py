from __future__ import annotations

from app.core.landing_page.google_trends import GoogleTrendsClient
from app.core.landing_page.google_trends.api.utils import build_query
from app.core.landing_page.google_trends.api._requests import Interval


def test_build_query_last_30_days():
    req = build_query(terms=["dealscale"], interval=Interval.daily, last_n_days=30, region="US")
    assert req.terms == ["dealscale"]
    assert req.interval == Interval.daily
    assert req.lastNDays == 30
    assert req.region == "US"


def test_interest_over_time_not_implemented():
    client = GoogleTrendsClient()
    try:
        client.interest_over_time("dealscale", geo="US", timeframe="today 1-m")
        assert False, "Expected NotImplementedError"
    except NotImplementedError:
        assert True
