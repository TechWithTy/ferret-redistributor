from __future__ import annotations

from datetime import date, timedelta
from typing import Iterable

from ._requests import TrendsQueryRequest, TrendsCompareRequest, DateRange, Interval


def build_date_range(*, last_n_days: int | None = None, start_date: str | None = None, end_date: str | None = None) -> list[DateRange] | None:
    if last_n_days is not None:
        end = date.today()
        start = end - timedelta(days=last_n_days)
        return [DateRange(startDate=start.isoformat(), endDate=end.isoformat())]
    if start_date and end_date:
        return [DateRange(startDate=start_date, endDate=end_date)]
    return None


def build_query(
    *,
    terms: Iterable[str],
    interval: Interval = Interval.daily,
    last_n_days: int | None = None,
    start_date: str | None = None,
    end_date: str | None = None,
    region: str | None = None,
    category: str | None = None,
) -> TrendsQueryRequest:
    ranges = build_date_range(last_n_days=last_n_days, start_date=start_date, end_date=end_date)
    return TrendsQueryRequest(
        terms=list(terms),
        interval=interval,
        dateRange=ranges[0] if ranges else None,
        lastNDays=last_n_days,
        region=region,
        category=category,
    )


def build_compare(*, queries: list[TrendsQueryRequest], interval: Interval | None = None) -> TrendsCompareRequest:
    # Optionally enforce same interval across queries for consistency
    if interval is not None:
        queries = [q.copy(update={"interval": interval}) for q in queries]
    return TrendsCompareRequest(queries=queries)
