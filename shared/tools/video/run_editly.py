#!/usr/bin/env python3
"""Invoke the backend TypeScript editly toolchain from shared tooling."""

from __future__ import annotations

import argparse
import os
import shutil
import subprocess
from pathlib import Path
from typing import List

REPO_ROOT = Path(__file__).resolve().parents[2]
TOOLS_DIR = REPO_ROOT / "backend" / "typescript" / "tools"


def main() -> None:
    parser = argparse.ArgumentParser(description="Run editly renders via backend/typescript/tools.")
    parser.add_argument("--config", "-c", required=True, help="Path to an editly JSON/JSON5 spec.")
    parser.add_argument("--output", "-o", help="Override outPath.")
    parser.add_argument("--width", type=int, help="Override width.")
    parser.add_argument("--height", type=int, help="Override height.")
    parser.add_argument("--fast", action="store_true", help="Enable editly fast mode.")
    parser.add_argument("--pnpm", help="Optional path to pnpm executable.")
    parser.add_argument("--cwd", default=str(TOOLS_DIR), help="Override working directory.")
    parser.add_argument("--dry-run", action="store_true", help="Print command without executing.")
    args = parser.parse_args()

    cwd = Path(args.cwd).resolve()
    if not cwd.exists():
        raise SystemExit(f"Workspace missing: {cwd}")

    pnpm_cmd = resolve_pnpm(args.pnpm)
    command = build_command(pnpm_cmd, args)

    print("[video_pipeline]", " ".join(command))
    if args.dry_run:
        return

    subprocess.run(command, check=True, cwd=cwd)
    print("[video_pipeline] render finished")


def resolve_pnpm(explicit: str | None) -> str:
    candidates = [explicit, os.getenv("PNPM_BIN"), "pnpm"]
    for candidate in candidates:
        if not candidate:
            continue
        resolved = shutil.which(candidate)
        if resolved:
            return resolved
    raise SystemExit("pnpm not found. Install pnpm or set PNPM_BIN/--pnpm.")


def build_command(pnpm_cmd: str, args: argparse.Namespace) -> List[str]:
    cmd: List[str] = [pnpm_cmd, "video:render", "--", "--config", args.config]
    if args.output:
        cmd += ["--output", args.output]
    if args.width:
        cmd += ["--width", str(args.width)]
    if args.height:
        cmd += ["--height", str(args.height)]
    if args.fast:
        cmd.append("--fast")
    return cmd


if __name__ == "__main__":
    main()

