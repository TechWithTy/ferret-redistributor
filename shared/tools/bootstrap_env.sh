#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
VENV_DIR="$ROOT_DIR/.venv"

echo "[shared/tools] Root directory: $ROOT_DIR"
export SOCIAL_SCALE_ROOT="$ROOT_DIR"

if [[ ! -d "$VENV_DIR" ]]; then
  echo "[shared/tools] Creating Python virtualenv in .venv"
  python -m venv "$VENV_DIR"
fi

ACTIVATE_SCRIPT="$VENV_DIR/bin/activate"
if [[ "$OS" == "Windows_NT" ]]; then
  ACTIVATE_SCRIPT="$VENV_DIR/Scripts/activate"
fi

echo "[shared/tools] Activating virtualenv: $ACTIVATE_SCRIPT"
# shellcheck disable=SC1090
source "$ACTIVATE_SCRIPT"

if [[ -f "$ROOT_DIR/requirements.txt" ]]; then
  echo "[shared/tools] Installing Python dependencies"
  pip install --upgrade pip
  pip install -r "$ROOT_DIR/requirements.txt"
fi

cat <<'EOF'
[shared/tools] Environment bootstrapped.
- SOCIAL_SCALE_ROOT exported
- Python virtualenv ready
- Add more setup steps here as shared tooling grows
EOF

