"""Quickstart for GA4 service account flow.
Ensure the service account email is added as a user on the GA4 property (Viewer/Analyst).
"""
from __future__ import annotations

import os
from datetime import date, timedelta

from ..client import GA4Client
from ..config import SCOPE_ANALYTICS_READONLY
from ..api._requests import RunReportRequest, Dimension, Metric, DateRange, RealtimeReportRequest


def main() -> None:
    client = GA4Client.from_service_account(scopes=[SCOPE_ANALYTICS_READONLY])

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
    print(f"RunReport rows: {len(resp.rows)}")

    rreq = RealtimeReportRequest(
        dimensions=[Dimension(name="eventName")],
        metrics=[Metric(name="activeUsers")],
        limit=10,
    )
    rresp = client.run_realtime_report(property_id, rreq)
    print(f"Realtime rows: {len(rresp.rows)}")


if __name__ == "__main__":
    main()
