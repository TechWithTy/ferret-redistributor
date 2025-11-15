from __future__ import annotations

import json
import os
import time
from pathlib import Path
from typing import Any, Optional


class FileCache:
    def __init__(self, base_dir: str | os.PathLike = "_data_cache") -> None:
        self.base = Path(base_dir)
        self.base.mkdir(parents=True, exist_ok=True)

    def _path_for(self, key: str) -> Path:
        safe = key.replace(os.sep, "_").replace(":", "_")
        return self.base / f"{safe}.json"

    def get(self, key: str) -> Optional[Any]:
        p = self._path_for(key)
        if not p.exists():
            return None
        try:
            data = json.loads(p.read_text(encoding="utf-8"))
        except Exception:
            return None
        expires_at = data.get("expires_at", 0)
        if expires_at and time.time() > float(expires_at):
            try:
                p.unlink(missing_ok=True)
            except Exception:
                pass
            return None
        return data.get("value")

    def set(self, key: str, value: Any, ttl_seconds: int = 3600) -> None:
        p = self._path_for(key)
        tmp = p.with_suffix(".json.tmp")
        payload = {"expires_at": time.time() + ttl_seconds, "value": value}
        tmp.write_text(json.dumps(payload, ensure_ascii=False), encoding="utf-8")
        os.replace(tmp, p)

    def delete(self, key: str) -> None:
        p = self._path_for(key)
        try:
            p.unlink(missing_ok=True)
        except Exception:
            pass

