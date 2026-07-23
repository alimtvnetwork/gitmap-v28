#!/usr/bin/env bash
# full-suite-lint.sh — strict golangci-lint pass for the full-suite-guard job.
# Captures stdout to /tmp/full-suite/lint-output.txt for the PR summary
# artifact, counts findings, and propagates the linter's exit code as the
# job gate. golangci-lint must already be on PATH (install via the
# install-golangci-lint composite action).
#
# Outputs (written to GITHUB_OUTPUT when set):
#   exit_code    — raw golangci-lint exit code
#   issue_count  — number of finding lines (path:line:col: msg) in output
#
# Run from the gitmap/ working directory.
set -uo pipefail

mkdir -p /tmp/full-suite

set +e
golangci-lint run ./... \
  --timeout=5m \
  --max-issues-per-linter=0 \
  --max-same-issues=0 2>&1 | tee /tmp/full-suite/lint-output.txt
rc=${PIPESTATUS[0]}
set -e

issues=$(grep -cE '^[^[:space:]].+:[0-9]+:[0-9]+:' /tmp/full-suite/lint-output.txt || true)

if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
  echo "exit_code=$rc"     >> "$GITHUB_OUTPUT"
  echo "issue_count=$issues" >> "$GITHUB_OUTPUT"
fi

exit "$rc"
