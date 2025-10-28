from __future__ import annotations

from app.core.landing_page.google_trends.api.utils import build_query, build_compare
from app.core.landing_page.google_trends.api._requests import Interval


def test_build_compare_multiple_terms():
    q1 = build_query(terms=["ai"], interval=Interval.daily, last_n_days=7, region="US")
    q2 = build_query(terms=["fastapi"], interval=Interval.daily, last_n_days=7, region="US")
    cmp_req = build_compare(queries=[q1, q2])
    assert len(cmp_req.queries) == 2
    assert cmp_req.queries[0].terms == ["ai"]
    assert cmp_req.queries[1].terms == ["fastapi"]
