#!/usr/bin/env python

from __future__ import annotations

import argparse
import sys
from dataclasses import dataclass
from pathlib import Path
from textwrap import indent
import re

ROOT = Path(__file__).resolve().parents[2]
SCHEMA = ROOT / "shared" / "models" / "social_scale_core.toon"
GENERATED_GO = ROOT / "shared" / "generated" / "go"
GENERATED_TS = ROOT / "shared" / "generated" / "ts"


@dataclass
class Field:
    name: str
    type_: str
    required: bool


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Shared TOON codegen entrypoint")
    parser.add_argument(
        "--schema",
        type=Path,
        default=SCHEMA,
        help="Path to the TOON schema (default: shared/models/social_scale_core.toon)",
    )
    return parser.parse_args()


def read_schema(path: Path) -> list[str]:
    if not path.exists():
        raise FileNotFoundError(f"Schema not found: {path}")
    lines = path.read_text(encoding="utf-8").splitlines()
    if not lines:
        raise ValueError("Schema is empty")
    return lines


def parse_enums(lines: list[str]) -> dict[str, list[str]]:
    enums: dict[str, list[str]] = {}
    in_block = False
    current: str | None = None
    for line in lines:
        if line.startswith("enums:"):
            in_block = True
            continue
        if in_block:
            if not line.startswith("  "):
                break
            stripped = line.strip()
            if stripped.endswith(":") and "[" in stripped:
                current = stripped.split("[", 1)[0]
                enums[current] = []
            elif stripped.startswith("- ") and current:
                enums[current].append(stripped[2:].strip())
    return enums


def parse_field_line(line: str) -> tuple[str, str, bool]:
    raw = line.strip()[2:]
    if "," not in raw:
        raise ValueError(f"invalid field line: {raw}")
    name, remainder = raw.split(",", 1)
    depth = 0
    type_chars: list[str] = []
    idx = 0
    while idx < len(remainder):
        ch = remainder[idx]
        if ch == "," and depth == 0:
            break
        if ch == "<":
            depth += 1
        elif ch == ">":
            depth = max(0, depth - 1)
        type_chars.append(ch)
        idx += 1
    type_str = "".join(type_chars).strip()
    rest = remainder[idx + 1 :].strip() if idx < len(remainder) else ""
    required_token = rest.split(",", 1)[0] if rest else "true"
    required = required_token.lower() not in {"false", "0", "no", ""}
    return name.strip(), type_str, required


def parse_entities(lines: list[str]) -> dict[str, list[Field]]:
    entities: dict[str, list[Field]] = {}
    in_entities = False
    current: str | None = None
    in_fields = False

    for line in lines:
        if line.startswith("entities:"):
            in_entities = True
            continue

        if not in_entities:
            continue
        if line and not line.startswith("  "):
            break

        stripped = line.strip()
        if stripped.endswith(":") and not stripped.startswith("-"):
            indent_level = len(line) - len(line.lstrip())
            if indent_level == 2:
                current = stripped.rstrip(":")
                entities[current] = []
                in_fields = False
            elif indent_level == 4 and "fields" in stripped:
                in_fields = True
            else:
                in_fields = False
            continue

        if in_fields and stripped.startswith("- ") and current:
            field_name, field_type, required = parse_field_line(stripped)
            entities[current].append(Field(field_name, field_type, required))
    return entities


ARRAY_RE = re.compile(r"array<(.+)>", re.IGNORECASE)
MAP_RE = re.compile(r"map<([^,>]+),([^>]+)>", re.IGNORECASE)


def go_type(field_type: str, required: bool) -> tuple[str, bool]:
    needs_time = False

    def base(typ: str) -> str:
        nonlocal needs_time
        typ = typ.strip()
        mapping = {
            "string": "string",
            "int": "int",
            "int64": "int64",
            "float": "float64",
            "bool": "bool",
            "datetime": "time.Time",
            "any": "any",
            "json": "map[string]any",
        }
        array_match = ARRAY_RE.match(typ)
        if array_match:
            inner = base(array_match.group(1))
            return f"[]{inner}"

        map_match = MAP_RE.match(typ)
        if map_match:
            key = base(map_match.group(1))
            value = base(map_match.group(2))
            if key != "string":
                key = "string"
            return f"map[{key}]{value}"

        target = mapping.get(typ, typ)
        if target == "time.Time":
            needs_time = True
        return target

    resolved = base(field_type)
    if not required and resolved in {"string", "int", "int64", "float64", "bool", "time.Time"}:
        resolved = f"*{resolved}"
    return resolved, needs_time


def ts_type(field_type: str, required: bool) -> str:
    def base(typ: str) -> str:
        typ = typ.strip()
        mapping = {
            "string": "string",
            "int": "number",
            "int64": "number",
            "float": "number",
            "bool": "boolean",
            "datetime": "string",
            "any": "any",
            "json": "Record<string, any>",
        }
        array_match = ARRAY_RE.match(typ)
        if array_match:
            inner = base(array_match.group(1))
            return f"{inner}[]"
        if MAP_RE.match(typ):
            return "Record<string, any>"
        return mapping.get(typ, typ)

    resolved = base(field_type)
    if not required:
        return f"{resolved} | null"
    return resolved


def to_pascal(name: str) -> str:
    return "".join(part.capitalize() for part in name.split("_"))


def generate_go(enums: dict[str, list[str]], entities: dict[str, list[Field]]) -> Path:
    GENERATED_GO.mkdir(parents=True, exist_ok=True)
    lines = [
        "// Code generated by shared/tools/run_codegen.py. DO NOT EDIT.",
        "package sharedmodels",
        "",
    ]
    needs_time_any = False

    for enum_name, values in sorted(enums.items()):
        lines.append(f"type {enum_name} string")
        lines.append("const (")
        for value in values:
            const_name = to_pascal(value)
            lines.append(f'\t{enum_name}{const_name} {enum_name} = "{value}"')
        lines.append(")")
        lines.append("")

    for entity, fields in sorted(entities.items()):
        lines.append(f"type {entity} struct {{")
        for field in fields:
            go_typ, needs_time = go_type(field.type_, field.required)
            needs_time_any = needs_time_any or needs_time
            field_name = to_pascal(field.name)
            tag = f'`json:"{field.name}'
            if not field.required and not go_typ.startswith("*"):
                tag += ",omitempty"
            tag += '"`'
            lines.append(f"\t{field_name} {go_typ} {tag}")
        lines.append("}")
        lines.append("")

    if needs_time_any:
        lines.insert(2, 'import "time"\n')

    output = GENERATED_GO / "models_gen.go"
    output.write_text("\n".join(lines).strip() + "\n", encoding="utf-8")
    return output


def generate_ts(enums: dict[str, list[str]], entities: dict[str, list[Field]]) -> Path:
    GENERATED_TS.mkdir(parents=True, exist_ok=True)
    lines = [
        "// Code generated by shared/tools/run_codegen.py. DO NOT EDIT.",
    ]

    for enum_name, values in sorted(enums.items()):
        joined = " | ".join(f'"{value}"' for value in values)
        lines.append(f"export type {enum_name} = {joined};")
    lines.append("")

    for entity, fields in sorted(entities.items()):
        lines.append(f"export interface {entity} {{")
        for field in fields:
            ts_typ = ts_type(field.type_, field.required)
            optional = "" if field.required else "?"
            lines.append(f"  {field.name}{optional}: {ts_typ};")
        lines.append("}")
        lines.append("")

    output = GENERATED_TS / "models.ts"
    output.write_text("\n".join(lines).strip() + "\n", encoding="utf-8")
    return output


def main() -> int:
    args = parse_args()
    try:
        lines = read_schema(args.schema)
        enums = parse_enums(lines)
        entities = parse_entities(lines)
        go_path = generate_go(enums, entities)
        ts_path = generate_ts(enums, entities)
    except Exception as exc:  # pragma: no cover - CLI guard
        print(f"[codegen] failed: {exc}", file=sys.stderr)
        return 1

    print("[codegen] generated:")
    print(indent(f"go -> {go_path}", "  "))
    print(indent(f"typescript -> {ts_path}", "  "))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

