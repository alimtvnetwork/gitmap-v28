#!/usr/bin/env bash
# Smoke test wrapper: forwards to Python 3 cross-platform runner.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTHON_BIN="python3"
if ! command -v python3 >/dev/null 2>&1; then
  PYTHON_BIN="python"
fi

exec "$PYTHON_BIN" "$SCRIPT_DIR/smoke-installer.py" "$@"
