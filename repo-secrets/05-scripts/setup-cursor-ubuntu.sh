#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_NODE="${1:-u1}"

echo "● Launching Cursor Ubuntu Fleet Setup for node: ${TARGET_NODE}"
if command -v python3 >/dev/null 2>&1; then
    python3 "${SCRIPT_DIR}/setup-cursor-ubuntu.py" --node "${TARGET_NODE}" "$@"
else
    echo "⚠ python3 not found, attempting apt bootstrap..."
    if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update -y && sudo apt-get install -y python3 curl
        python3 "${SCRIPT_DIR}/setup-cursor-ubuntu.py" --node "${TARGET_NODE}" "$@"
    else
        echo "❌ Cannot bootstrap python3 without apt-get" >&2
        exit 1
    fi
fi
