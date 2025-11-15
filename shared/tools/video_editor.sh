#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TOOLS_DIR="$REPO_ROOT/backend/typescript/tools"

if [[ ! -d "$TOOLS_DIR" ]]; then
  echo "[video_editor] missing backend/typescript/tools directory at $TOOLS_DIR" >&2
  exit 1
fi

if ! command -v pnpm >/dev/null 2>&1; then
  echo "[video_editor] pnpm is required. Install from https://pnpm.io/installation" >&2
  exit 1
fi

(
  cd "$TOOLS_DIR"
  PNPM_SCRIPT_ARGS=("$@")
  pnpm video:render -- "${PNPM_SCRIPT_ARGS[@]}"
)

