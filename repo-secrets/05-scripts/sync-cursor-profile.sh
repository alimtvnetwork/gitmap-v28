#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

ensure_python3() {
    if ! command -v python3 >/dev/null 2>&1; then
        echo "⚠ python3 not found, attempting apt bootstrap..."
        if command -v apt-get >/dev/null 2>&1; then
            sudo apt-get update -y && sudo apt-get install -y python3 tar curl
        else
            echo "❌ Cannot bootstrap python3 without apt-get" >&2
            exit 1
        fi
    fi
}

ensure_python3

echo "● Launching Cursor Profile & AI Memory Synchronization"
exec python3 "${SCRIPT_DIR}/sync-cursor-profile.py" "$@"
