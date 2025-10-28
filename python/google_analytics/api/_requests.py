from __future__ import annotations

from typing import List, Optional
from pydantic import BaseModel, Field

class DateRange(BaseModel):
    startDate: str
    endDate: str

class Dimension(BaseModel):
    name: str

class Metric(BaseModel):
    name: str

class OrderBy(BaseModel):
    fieldName: str
    desc: bool = False

class StringFilter(BaseModel):
    matchType: str = Field(default="EXACT")  # GA4 supports BEGINS_WITH, EXACT, etc.
    value: str

class Filter(BaseModel):
    fieldName: str
    stringFilter: Optional[StringFilter] = None

class FilterExpression(BaseModel):
    filter: Optional[Filter] = None

class RunReportRequest(BaseModel):
    dimensions: List[Dimension]
    metrics: List[Metric]
    dateRanges: Optional[List[DateRange]] = None
    limit: Optional[int] = None
    offset: Optional[int] = None
    orderBys: Optional[List[OrderBy]] = None
    dimensionFilter: Optional[FilterExpression] = None
    metricFilter: Optional[FilterExpression] = None

class RealtimeReportRequest(BaseModel):
    dimensions: List[Dimension]
    metrics: List[Metric]
    limit: Optional[int] = None
    orderBys: Optional[List[OrderBy]] = None
