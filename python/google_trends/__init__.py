"""Google Trends SDK (Alpha scaffold) package init."""

from .client import GoogleTrendsClient
from .api.utils import build_query, build_compare
from .api.deps import get_client as get_trends_client_dependency, ApiKeyProvider

__all__ = [
    "GoogleTrendsClient",
    "build_query",
    "build_compare",
    "get_trends_client_dependency",
    "ApiKeyProvider",
]
