from __future__ import annotations

from enum import Enum
from typing import List, Optional
from pydantic import BaseModel, Field


class Interval(str, Enum):
    daily = "daily"
    weekly = "weekly"
    monthly = "monthly"
    yearly = "yearly"


class DateRange(BaseModel):
    startDate: str  # YYYY-MM-DD
    endDate: str    # YYYY-MM-DD


class TrendsQueryRequest(BaseModel):
    terms: List[str] = Field(min_items=1)
    interval: Interval = Interval.daily
    dateRange: Optional[DateRange] = None
    lastNDays: Optional[int] = Field(default=None, ge=1)
    region: Optional[str] = None  # e.g. "US" or subregion like "US-CA"
    category: Optional[str] = None


class TrendsCompareRequest(BaseModel):
    queries: List[TrendsQueryRequest] = Field(min_items=2)
    # Optionally enforce same interval/region across queries (SDK can check)
