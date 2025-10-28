from __future__ import annotations

import argparse
import os

from .db import create_all


def main() -> int:
    ap = argparse.ArgumentParser(description="Calendar DB utilities")
    ap.add_argument("command", choices=["init-db"], help="command")
    ap.add_argument("--url", dest="url", help="DATABASE_URL override")
    ns = ap.parse_args()

    url = ns.url or os.getenv("DATABASE_URL")
    if ns.command == "init-db":
        create_all(url)
        print("Initialized DB schema")
        return 0
    return 1


if __name__ == "__main__":
    raise SystemExit(main())

