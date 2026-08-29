#!/usr/bin/env bash
exec python3 "$(dirname "$0")/misspell-changed.py" "$@"
