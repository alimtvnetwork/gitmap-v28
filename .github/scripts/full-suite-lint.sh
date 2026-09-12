#!/usr/bin/env bash
# Fast native shell runner for full suite golangci-lint
set -uo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
GITMAP_DIR="${1:-$REPO_ROOT/cli}"

if ! command -v golangci-lint >/dev/null 2>&1; then
    if command -v python3 >/dev/null 2>&1; then
        exec python3 "$(dirname "$0")/full-suite-lint.py" "$@"
    fi
    echo "ERROR: golangci-lint not on PATH" >&2
    exit 2
fi

OUTPUT_FILE="/tmp/full-suite-lint-$$.log"
mkdir -p /tmp

cd "$GITMAP_DIR" || exit 1
set +e
golangci-lint run ./... --timeout=5m --max-issues-per-linter=0 --max-same-issues=0 2>&1 | tee "$OUTPUT_FILE"
EXIT_CODE=${PIPESTATUS[0]}
set -e

ISSUE_COUNT=$(grep -cE '^\S+:[0-9]+:[0-9]+:' "$OUTPUT_FILE" 2>/dev/null || true)
rm -f "$OUTPUT_FILE"

if [[ -n "${GITHUB_OUTPUT:-}" && -f "$GITHUB_OUTPUT" ]]; then
    echo "exit_code=$EXIT_CODE" >> "$GITHUB_OUTPUT"
    echo "issue_count=$ISSUE_COUNT" >> "$GITHUB_OUTPUT"
fi

exit "$EXIT_CODE"
