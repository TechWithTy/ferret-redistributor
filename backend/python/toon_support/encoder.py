"""Lightweight TOON encoder tailored for Social Scale workflows.

The goal is not to be a perfect implementation of the reference spec but to
provide deterministic, token-efficient payloads for automation workflows. The
encoder supports:

* Arbitrary mappings (dict, dataclass, pydantic BaseModel via `model_dump`)
* Lists of primitives (rendered inline) and nested complex structures
* Tabular emission for homogeneous lists of mappings
* Strings that only quote when the TOON syntax requires it
"""

from __future__ import annotations

from dataclasses import asdict, is_dataclass
from typing import Any, Iterable, Mapping, Sequence

DEFAULT_INDENT = 2


def encode(value: Any, *, indent: int = DEFAULT_INDENT) -> str:
    """Serialize ``value`` into a TOON string."""

    builder: list[str] = []
    _emit_value(_normalize(value), builder, level=0, indent=indent, prefix="")
    return "\n".join(line.rstrip() for line in builder).rstrip()


def encode_to_toon(name: str, payload: Any, *, indent: int = DEFAULT_INDENT) -> str:
    """Encode a named workflow payload (e.g., ``encode_to_toon("post", {...})``)."""

    return encode({name: payload}, indent=indent)


def _normalize(value: Any) -> Any:
    if value is None:
        return None

    if is_dataclass(value):
        return {k: _normalize(v) for k, v in asdict(value).items()}

    if hasattr(value, "model_dump"):
        return _normalize(value.model_dump())

    if isinstance(value, Mapping):
        return {str(k): _normalize(v) for k, v in value.items()}

    if isinstance(value, Sequence) and not isinstance(value, (str, bytes, bytearray)):
        return [_normalize(v) for v in value]

    return value


def _emit_value(value: Any, builder: list[str], level: int, indent: int, prefix: str) -> None:
    if isinstance(value, Mapping):
        _emit_mapping(value, builder, level, indent, prefix)
    elif isinstance(value, Sequence) and not isinstance(value, (str, bytes, bytearray)):
        _emit_sequence(value, builder, level, indent, prefix)
    else:
        builder.append(f"{prefix}{_format_scalar(value)}")


def _emit_mapping(value: Mapping[str, Any], builder: list[str], level: int, indent: int, prefix: str) -> None:
    keys = sorted(value.keys())
    for key in keys:
        val = value[key]
        line_prefix = f"{prefix}{_indent(level, indent)}{key}:"
        if _is_leaf(val):
            builder.append(f"{line_prefix} {_format_scalar(val)}")
        else:
            builder.append(line_prefix)
            _emit_value(val, builder, level + 1, indent, prefix)


def _emit_sequence(seq: Sequence[Any], builder: list[str], level: int, indent: int, prefix: str) -> None:
    if not seq:
        builder.append(f"{prefix}{_indent(level, indent)}[]")
        return

    if _is_tabular(seq):
        headers = sorted(seq[0].keys())
        header_line = ",".join(headers)
        builder.append(
            f"{prefix}{_indent(level, indent)}[{len(seq)}]{{{header_line}}}:"
        )
        for row in seq:
            values = ",".join(_format_scalar(row[h]) for h in headers)
            builder.append(f"{prefix}{_indent(level + 1, indent)}{values}")
        return

    if all(_is_leaf(item) for item in seq):
        inline_values = ",".join(_format_scalar(item) for item in seq)
        builder.append(f"{prefix}{_indent(level, indent)}[{len(seq)}]: {inline_values}")
        return

    builder.append(f"{prefix}{_indent(level, indent)}[{len(seq)}]:")
    for item in seq:
        item_prefix = f"{prefix}{_indent(level + 1, indent)}- "
        if _is_leaf(item):
            builder.append(f"{item_prefix}{_format_scalar(item)}")
        else:
            builder.append(item_prefix.strip())
            _emit_value(item, builder, level + 2, indent, prefix)


def _is_leaf(value: Any) -> bool:
    return not isinstance(value, (Mapping, Sequence)) or isinstance(value, (str, bytes, bytearray))


def _is_tabular(seq: Sequence[Any]) -> bool:
    if not seq or not all(isinstance(item, Mapping) for item in seq):
        return False

    first_keys = set(seq[0].keys())
    if not first_keys:
        return False

    return all(set(item.keys()) == first_keys and all(_is_leaf(v) for v in item.values()) for item in seq)


def _indent(level: int, indent: int) -> str:
    return " " * (level * indent)


def _format_scalar(value: Any) -> str:
    if value is None:
        return "null"
    if isinstance(value, bool):
        return "true" if value else "false"
    if isinstance(value, (int, float)):
        return repr(value)

    text = str(value)
    if not text:
        return '""'

    if any(char in text for char in [",", ":", "\n", '"']) or text.strip() != text:
        return f'"{text}"'

    lower = text.lower()
    if lower in {"null", "true", "false"}:
        return f'"{text}"'

    if _looks_numeric(text):
        return f'"{text}"'

    return text


def _looks_numeric(text: str) -> bool:
    try:
        float(text)
        return True
    except ValueError:
        return False

