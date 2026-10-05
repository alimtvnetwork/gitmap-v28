#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTHON_SCRIPT="${SCRIPT_DIR}/heal-u1-pull-errors.py"

echo "=== GitMap Remote Node Diagnostic & Healing Wrapper ==="

if ! command -v python3 >/dev/null 2>&1; then
    echo "[ERROR] python3 is required but not installed." >&2
    exit 1
fi

python3 "${PYTHON_SCRIPT}" "$@"
