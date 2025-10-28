from __future__ import annotations

from typing import List, Optional
from pydantic import BaseModel

class DimensionHeader(BaseModel):
    name: str

class MetricHeader(BaseModel):
    name: str
    type: Optional[str] = None

class Row(BaseModel):
    dimensionValues: List[dict]
    metricValues: List[dict]

class Metadata(BaseModel):
    dimensions: Optional[list] = None
    metrics: Optional[list] = None

class RunReportResponse(BaseModel):
    dimensionHeaders: List[DimensionHeader]
    metricHeaders: List[MetricHeader]
    rows: List[Row] = []
    totals: Optional[List[Row]] = None
    metadata: Optional[Metadata] = None
    rowCount: Optional[int] = None

class RealtimeReportResponse(BaseModel):
    dimensionHeaders: List[DimensionHeader]
    metricHeaders: List[MetricHeader]
    rows: List[Row] = []
    totals: Optional[List[Row]] = None
