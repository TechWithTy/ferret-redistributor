"""Convenience helpers for encoding Social Scale payloads into TOON.

The actual format is specified at https://github.com/johannschopplich/toon.
This module provides a tiny, dependency-free encoder that mirrors the
primitives we need inside Social Scale's Python automations.
"""

from .encoder import encode, encode_to_toon

__all__ = ["encode", "encode_to_toon"]

