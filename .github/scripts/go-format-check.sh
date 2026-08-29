#!/usr/bin/env bash
exec python3 "$(dirname "$0")/go-format-check.py" "$@"
