#!/usr/bin/env python3
"""Copy TOON schemas into shared/tools/schemas for MCP discovery."""

from __future__ import annotations

import shutil
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[1]
MODELS_DIR = REPO_ROOT / "shared" / "models"
TARGET_DIR = REPO_ROOT / "shared" / "tools" / "schemas"


def main() -> None:
    TARGET_DIR.mkdir(parents=True, exist_ok=True)

    # Clean old copies
    for toon_file in TARGET_DIR.glob("*.toon"):
        toon_file.unlink()

    copied = 0
    for source in sorted(MODELS_DIR.glob("*.toon")):
        shutil.copy2(source, TARGET_DIR / source.name)
        copied += 1

    print(f"[sync_schemas] copied {copied} schema(s) into {TARGET_DIR.relative_to(REPO_ROOT)}")


if __name__ == "__main__":
    main()


