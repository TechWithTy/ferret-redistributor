#!/usr/bin/env python3
"""Apply a MagickCore pipeline spec using the ImageMagick CLI."""

from __future__ import annotations

import argparse
import json
import os
import shutil
import subprocess
from pathlib import Path
from typing import Any, Dict, List

REPO_ROOT = Path(__file__).resolve().parents[2]


def main() -> None:
    parser = argparse.ArgumentParser(description="Run ImageMagick pipeline from a TOON spec.")
    parser.add_argument("--spec", required=True, help="Path to a JSON/TOON MagickPipeline file.")
    parser.add_argument(
        "--magick",
        help="Optional path to the `magick` executable (defaults to IMAGEMAGICK_BIN or PATH lookup).",
    )
    parser.add_argument("--dry-run", action="store_true", help="Print the command without executing.")
    args = parser.parse_args()

    spec = load_spec(Path(args.spec))
    magick_bin = resolve_magick(args.magick)

    cmd = build_command(spec, magick_bin)
    print("[magick_pipeline]", " ".join(cmd))

    if args.dry_run:
        return

    subprocess.run(cmd, check=True)
    print("[magick_pipeline] wrote", spec["output_path"])


def load_spec(path: Path) -> Dict[str, Any]:
    data = json.loads(path.read_text(encoding="utf-8"))
    required = ("source_path", "output_path", "operations")
    for key in required:
        if key not in data:
            raise SystemExit(f"Missing required field `{key}` in {path}")
    return data


def resolve_magick(explicit: str | None) -> str:
    candidates = [explicit, os.getenv("IMAGEMAGICK_BIN"), "magick"]
    for candidate in candidates:
        if not candidate:
            continue
        resolved = shutil.which(candidate)
        if resolved:
            return resolved
        candidate_path = Path(candidate)
        if candidate_path.exists():
            return str(candidate_path)
    raise SystemExit("Unable to locate `magick`. Install ImageMagick or set IMAGEMAGICK_BIN.")


def build_command(spec: Dict[str, Any], magick_bin: str) -> List[str]:
    cmd: List[str] = [magick_bin]

    # Global options
    if spec.get("colorspace"):
        cmd += ["-colorspace", spec["colorspace"].upper()]
    if spec.get("depth"):
        cmd += ["-depth", str(spec["depth"])]
    if spec.get("density"):
        cmd += ["-density", str(spec["density"])]
    if spec.get("background_color"):
        cmd += ["-background", spec["background_color"]]
    if spec.get("quality"):
        cmd += ["-quality", str(spec["quality"])]
    if spec.get("alpha") is not None:
        cmd += ["-alpha", "on" if spec["alpha"] else "off"]
    if spec.get("compression"):
        cmd += ["-compress", spec["compression"]]
    for profile in spec.get("profiles", []):
        cmd += ["-profile", profile]

    cmd.append(spec["source_path"])

    for op in spec["operations"]:
        cmd.extend(operation_to_args(op))

    cmd.append(spec["output_path"])
    return cmd


def operation_to_args(op: Dict[str, Any]) -> List[str]:
    name = op.get("op_name")
    if not name:
        raise SystemExit("Operation is missing `op_name`.")
    norm = name.strip().lower()
    geometry = op.get("geometry") or {}
    args: List[str] = []

    if norm == "resizeimage":
        dims = ensure_geometry_wh(geometry)
        args += optional_filter(op)
        args += ["-resize", dims]
    elif norm == "cropimage":
        dims = ensure_geometry_wh(geometry)
        args += ["-crop", append_offset(dims, geometry)]
    elif norm == "extentimage":
        dims = ensure_geometry_wh(geometry)
        args += ["-extent", append_offset(dims, geometry)]
    elif norm == "blurimage":
        radius = op.get("radius", 0)
        sigma = op.get("sigma", radius)
        args += ["-blur", f"{radius}x{sigma}"]
    elif norm == "sharpenimage":
        radius = op.get("radius", 0)
        sigma = op.get("sigma", radius)
        args += ["-sharpen", f"{radius}x{sigma}"]
    elif norm == "rotateimage":
        args += ["-rotate", str(op.get("angle", 0))]
    elif norm == "flipimage":
        args += ["-flip"]
    elif norm == "flopimage":
        args += ["-flop"]
    elif norm == "shearimage":
        args += ["-shear", f"{geometry.get('x_offset', 0)}x{geometry.get('y_offset', 0)}"]
    else:
        raise SystemExit(f"Unsupported op_name '{name}'. Please extend pipeline.py.")

    return args


def optional_filter(op: Dict[str, Any]) -> List[str]:
    filt = op.get("filter")
    if not filt:
        return []
    mapping = {
        "undefined": None,
        "point": "Point",
        "box": "Box",
        "triangle": "Triangle",
        "hermite": "Hermite",
        "lanczos": "Lanczos",
        "gaussian": "Gaussian",
        "catrom": "Catrom",
    }
    mapped = mapping.get(filt.lower())
    return ["-filter", mapped] if mapped else []


def ensure_geometry_wh(geometry: Dict[str, Any]) -> str:
    width = geometry.get("width")
    height = geometry.get("height")
    if width is None or height is None:
        raise SystemExit("Geometry requires width and height.")
    return f"{int(width)}x{int(height)}"


def append_offset(dims: str, geometry: Dict[str, Any]) -> str:
    x_off = geometry.get("x_offset", 0)
    y_off = geometry.get("y_offset", 0)
    return f"{dims}+{int(x_off)}+{int(y_off)}"


if __name__ == "__main__":
    main()

