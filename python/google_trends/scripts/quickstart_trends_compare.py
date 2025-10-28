from __future__ import annotations

"""Quickstart (alpha scaffold): compare multiple terms.
This demonstrates building compare requests; client methods are stubs.
"""

import os
from app.core.landing_page.google_trends import GoogleTrendsClient
from app.core.landing_page.google_trends.api.utils import build_query, build_compare
from app.core.landing_page.google_trends.api._requests import Interval


def main():
    terms = os.getenv("TRENDS_TERMS", "ai,startups,fastapi").split(",")
    region = os.getenv("TRENDS_REGION", "US")

    client = GoogleTrendsClient()
    q = [build_query(terms=[t], interval=Interval.daily, last_n_days=30, region=region) for t in terms]
    compare_req = build_compare(queries=q)

    try:
        # Placeholder; implement when backend is wired
        raise NotImplementedError("compare endpoint not implemented")
    except NotImplementedError as e:
        print("[Alpha] compare not implemented yet:", e)
        print("Built compare request:", compare_req.model_dump())


if __name__ == "__main__":
    main()
