from __future__ import annotations

from typing import List, Optional
from pydantic import BaseModel
from ._requests import Interval


class TimePoint(BaseModel):
    date: str  # YYYY-MM-DD
    value: float
    scaled: bool = True


class Series(BaseModel):
    term: str
    region: Optional[str] = None
    interval: Interval
    points: List[TimePoint]


class TrendsQueryResponse(BaseModel):
    series: Series
    coverage: Optional[str] = None
    notes: Optional[list[str]] = None


class TrendsCompareResponse(BaseModel):
    series: List[Series]
    mergeHints: Optional[list[str]] = None


class RegionInfo(BaseModel):
    code: str
    name: str
    level: str  # country, subregion, etc.


class RegionsResponse(BaseModel):
    regions: List[RegionInfo]


class TrendsMetadata(BaseModel):
    supportedIntervals: List[Interval]
    maxLookbackYears: int
    regionsAvailable: List[str]
    categoriesAvailable: Optional[List[str]] = None
