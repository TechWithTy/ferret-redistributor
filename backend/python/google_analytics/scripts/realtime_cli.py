"""Realtime CLI for GA4.
Fetch realtime metrics/dimensions dynamically.

Usage:
  python -m google_analytics.scripts.realtime_cli \
    --property $GA_PROPERTY_ID \
    --metrics activeUsers,eventCount \
    --dimensions eventName,streamId \
    --limit 10

Auth:
  - Service Account (recommended): set GA_SA_KEY and GA_SCOPE env; script uses SA by default
  - Installed App: set GA_CLIENT_SECRETS and GA_TOKEN_PATH; pass --installed-app
"""
from __future__ import annotations

import argparse
import json
import os
import sys
from typing import List

from ..client import GA4Client
from ..config import SCOPE_ANALYTICS_READONLY
from ..api.utils import build_realtime_request


def parse_args(argv: List[str]) -> argparse.Namespace:
    p = argparse.ArgumentParser(description="GA4 Realtime dynamic fetch")
    p.add_argument("--property", dest="property_id", required=True, help="GA4 property id (numeric)")
    p.add_argument("--metrics", required=True, help="Comma-separated metrics, e.g. activeUsers,eventCount")
    p.add_argument("--dimensions", default="", help="Comma-separated dimensions, e.g. eventName,streamId")
    p.add_argument("--limit", type=int, default=10)
    p.add_argument("--installed-app", action="store_true", help="Use installed-app OAuth instead of service account")
    return p.parse_args(argv)


def main(argv: List[str] | None = None) -> int:
    ns = parse_args(argv or sys.argv[1:])

    if ns.installed_app:
        client = GA4Client.from_installed_app(scopes=[SCOPE_ANALYTICS_READONLY])
    else:
        client = GA4Client.from_service_account(scopes=[SCOPE_ANALYTICS_READONLY])

    metrics = [m.strip() for m in ns.metrics.split(",") if m.strip()]
    dimensions = [d.strip() for d in ns.dimensions.split(",") if d.strip()]

    req = build_realtime_request(dimensions=dimensions, metrics=metrics, limit=ns.limit)
    resp = client.run_realtime_report(ns.property_id, req)
    print(json.dumps(resp.model_dump(), indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

