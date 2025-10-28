"""Quickstart for GA4 installed-app OAuth flow."""
from __future__ import annotations

import os
from datetime import date, timedelta

from ..client import GA4Client
from ..config import SCOPE_ANALYTICS_READONLY
from ..api._requests import RunReportRequest, Dimension, Metric, DateRange


def main() -> None:
    client = GA4Client.from_installed_app(scopes=[SCOPE_ANALYTICS_READONLY])

    property_id = os.getenv("GA_PROPERTY_ID")
    if not property_id:
        raise SystemExit("Set GA_PROPERTY_ID in environment")

    end = date.today()
    start = end - timedelta(days=7)

    req = RunReportRequest(
        dimensions=[Dimension(name="date"), Dimension(name="pagePath")],
        metrics=[Metric(name="screenPageViews"), Metric(name="activeUsers")],
        dateRanges=[DateRange(startDate=start.isoformat(), endDate=end.isoformat())],
        limit=10,
    )

    resp = client.run_report(property_id, req)
    print(f"Rows: {len(resp.rows)}")
    if resp.rows:
        print("Sample row:", resp.rows[0].model_dump())


if __name__ == "__main__":
    main()
