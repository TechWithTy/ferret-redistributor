from __future__ import annotations

from datetime import date, timedelta
from typing import Iterable, List

from ._requests import (
    RunReportRequest,
    RealtimeReportRequest,
    Dimension,
    Metric,
    DateRange,
)


def build_run_report_request(
    *,
    dimensions: Iterable[str],
    metrics: Iterable[str],
    start_date: str | None = None,
    end_date: str | None = None,
    last_n_days: int | None = None,
    limit: int | None = None,
):
    """Convenience builder for RunReportRequest.

    - Provide explicit start_date/end_date (YYYY-MM-DD) OR last_n_days.
    - Dimensions/Metrics as simple string names.
    """
    dims = [Dimension(name=d) for d in dimensions]
    mets = [Metric(name=m) for m in metrics]

    ranges: list[DateRange] | None = None
    if last_n_days is not None:
        end = date.today()
        start = end - timedelta(days=last_n_days)
        ranges = [DateRange(startDate=start.isoformat(), endDate=end.isoformat())]
    elif start_date and end_date:
        ranges = [DateRange(startDate=start_date, endDate=end_date)]

    return RunReportRequest(
        dimensions=dims,
        metrics=mets,
        dateRanges=ranges,
        limit=limit,
    )


def build_realtime_request(
    *, dimensions: Iterable[str], metrics: Iterable[str], limit: int | None = None
):
    """Convenience builder for RealtimeReportRequest with string names."""
    dims = [Dimension(name=d) for d in dimensions]
    mets = [Metric(name=m) for m in metrics]
    return RealtimeReportRequest(dimensions=dims, metrics=mets, limit=limit)
