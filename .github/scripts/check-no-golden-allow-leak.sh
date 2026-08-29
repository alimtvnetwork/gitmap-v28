#!/usr/bin/env bash
exec python3 "$(dirname "$0")/check-no-golden-allow-leak.py" "$@"
