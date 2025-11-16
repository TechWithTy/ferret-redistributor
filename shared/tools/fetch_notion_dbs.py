#!/usr/bin/env python3
"""Synchronize Notion databases that back the Experiments stack.

The script queries the Notion API for each configured database, then writes the
raw response payloads to JSON snapshot files (one per database) so ingestion
jobs can consume data offline.

Usage:
    shared/tools/fetch_notion_dbs.py \
        --output shared/generated/notion_snapshots
"""

from __future__ import annotations

import argparse
import json
import os
from pathlib import Path
from typing import Dict, Iterable, Optional

import requests

NOTION_API_BASE = os.getenv("NOTION_API_BASE", "https://api.notion.com/v1")
NOTION_API_VERSION = "2022-06-28"
NOTION_TOKEN_ENV = "NOTION_API_TOKEN"

# Mapping of logical names to env vars (with sane defaults pulled from Notion URLs)
DB_SPECS: Dict[str, Dict[str, str]] = {
    "experiments": {
        "env": "NOTION_DB_EXPERIMENTS_ID",
        "default": "6be08800e8d54cde93f14c26061f7b79",
    },
    "creative_assets": {
        "env": "NOTION_DB_CREATIVE_ASSETS_ID",
        "default": "299e9c25ecb081808164f32b83e31b45",
    },
    "iterations_actions": {
        "env": "NOTION_DB_ITERATIONS_ID",
        "default": "295e9c25ecb0816ebd87f774cfb3312e",
    },
    "kpi_definitions": {
        "env": "NOTION_DB_KPI_DEFS_ID",
        "default": "299e9c25ecb08111b2c4d96bddc6fc54",
    },
    "kpi_progress": {
        "env": "NOTION_DB_KPI_PROGRESS_ID",
        "default": "292e9c25ecb080a382d3f7b8aee64f7b",
    },
    "scripts_variants": {
        "env": "NOTION_DB_SCRIPTS_ID",
        "default": "299e9c25ecb0819ab02efca74143cf19",
    },
    "copy_calendar": {
        "env": "NOTION_DB_COPY_CALENDAR_ID",
        "default": "299e9c25ecb081cd8512eaf0bca8e2c7",
    },
    "channels": {
        "env": "NOTION_DB_CHANNELS_ID",
        "default": "299e9c25ecb0813b9dc5c78fb53d8866",
    },
    "platforms": {
        "env": "NOTION_DB_PLATFORMS_ID",
        "default": "299e9c25ecb081b59e41c9e258976fba",
    },
    "media_edit_tracker": {
        "env": "NOTION_DB_MEDIA_EDIT_TRACKER_ID",
        "default": "2ade9c25ecb080c495bbe7c5dbecf0e2",
    },
}


class NotionClient:
    def __init__(self, token: str):
        self._session = requests.Session()
        self._session.headers.update(
            {
                "Authorization": f"Bearer {token}",
                "Notion-Version": NOTION_API_VERSION,
                "Content-Type": "application/json",
            }
        )

    def fetch_database(self, database_id: str) -> dict:
        """Fetch all pages for a database (handles pagination)."""
        url = f"{NOTION_API_BASE}/databases/{database_id}/query"
        payload: Dict[str, Optional[str]] = {"page_size": 100}
        all_results = []

        while True:
            resp = self._session.post(url, json=payload)
            resp.raise_for_status()
            data = resp.json()
            all_results.extend(data.get("results", []))

            if not data.get("has_more"):
                break
            payload["start_cursor"] = data.get("next_cursor")

        return {"database_id": database_id, "results": all_results}


def resolve_database_id(name: str) -> str:
    spec = DB_SPECS[name]
    return os.getenv(spec["env"], spec["default"])


def ensure_output_dir(path: Path) -> None:
    path.mkdir(parents=True, exist_ok=True)


def write_snapshot(output_dir: Path, name: str, payload: dict) -> Path:
    target = output_dir / f"{name}.json"
    target.write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return target


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Fetch Notion databases for Experiments context.")
    parser.add_argument(
        "--output",
        type=Path,
        default=Path("shared/generated/notion_snapshots"),
        help="Directory to write JSON snapshots (default: %(default)s)",
    )
    parser.add_argument(
        "--only",
        nargs="+",
        choices=sorted(DB_SPECS.keys()),
        help="Limit fetches to the provided database names.",
    )
    return parser.parse_args()


def determine_targets(selection: Optional[Iterable[str]]) -> Iterable[str]:
    return selection if selection else DB_SPECS.keys()


def main() -> None:
    args = parse_args()
    ensure_output_dir(args.output)

    token = os.getenv(NOTION_TOKEN_ENV)
    if not token:
        raise SystemExit(f"Missing {NOTION_TOKEN_ENV}. Export your Notion integration token first.")

    client = NotionClient(token=token)
    fetched = []

    for name in determine_targets(args.only):
        database_id = resolve_database_id(name)
        print(f"[notion] Fetching {name} ({database_id})...")
        payload = client.fetch_database(database_id)
        snapshot_path = write_snapshot(args.output, name, payload)
        fetched.append(snapshot_path)

    print(f"[notion] Wrote {len(fetched)} snapshots → {args.output.resolve()}")


if __name__ == "__main__":
    main()




