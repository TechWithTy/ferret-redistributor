from __future__ import annotations

"""Quickstart (alpha scaffold): single-term query.
This demonstrates request building; actual client method is a stub.
"""

import os
from app.core.landing_page.google_trends import GoogleTrendsClient
from app.core.landing_page.google_trends.api.utils import build_query
from app.core.landing_page.google_trends.api._requests import Interval


def main():
    term = os.getenv("TRENDS_TERM", "dealscale")
    region = os.getenv("TRENDS_REGION", "US")

    client = GoogleTrendsClient()
    req = build_query(terms=[term], interval=Interval.daily, last_n_days=30, region=region)

    # NOTE: NotImplementedError until backend is wired (pytrends or official API)
    try:
        data = client.interest_over_time(keyword=term, geo=region, timeframe="today 1-m")
        print(data)
    except NotImplementedError as e:
        print("[Alpha] interest_over_time not implemented yet:", e)
        print("Built request:", req.model_dump())


if __name__ == "__main__":
    main()
