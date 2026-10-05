#!/usr/bin/env bash
# ==============================================================================
# migrate-cursor-memories-conversations.sh
# Shell executable wrapper for Cursor memories, conversations & projects migration.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Discover Python 3 executable across Linux and Windows Git Bash environments
PYTHON_BIN=""
if command -v python3 >/dev/null 2>&1; then
    PYTHON_BIN="python3"
elif command -v python >/dev/null 2>&1; then
    PYTHON_BIN="python"
fi

# Fallback: Attempt apt-get bootstrap on Debian/Ubuntu systems if python3 is missing
if [ -z "${PYTHON_BIN}" ]; then
    echo "⚠ python3 not found in PATH, attempting apt bootstrap..."
    if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update -y && sudo apt-get install -y python3 tar curl
        PYTHON_BIN="python3"
    else
        echo "❌ Cannot bootstrap python3: apt-get not available" >&2
        exit 1
    fi
fi

echo "● Launching Cursor Memories, Conversations & Projects Migration Engine"
exec "${PYTHON_BIN}" "${SCRIPT_DIR}/migrate-cursor-memories-conversations.py" "$@"
