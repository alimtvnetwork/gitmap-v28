#!/usr/bin/env bash
exec python3 "$(dirname "$0")/check-single-linter-diff.py" "$@"
