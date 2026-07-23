#!/usr/bin/env bash
# Hard-gate: forbid bare `fmt.Fprintln(os.Stderr, err)` in gitmap/cmd/.
#
# Per spec/04-generic-cli/07-error-handling.md and the cliexit
# package contract, every user-facing failure in a `gitmap`
# subcommand must include actionable context (command + op +
# subject + cause) — bare `Fprintln(os.Stderr, err)` strips three
# of those four. Use `cliexit.Reportf` / `cliexit.Fail` instead.
#
# Diff vs baseline: only NEW occurrences fail. Pre-existing sites
# in non-cmd packages or _test.go files are ignored entirely.
#
# Exit codes:
#   0 — no new bare-err sites in gitmap/cmd/*.go (excluding tests)
#   1 — at least one new offender
#   2 — toolchain missing

set -euo pipefail

if ! command -v rg >/dev/null 2>&1; then
  echo "ERROR: ripgrep (rg) not on PATH" >&2
  exit 2
fi

PATTERN='fmt\.Fprintln\(os\.Stderr, err\)'
HITS=$(rg -n --type go --glob '!*_test.go' "$PATTERN" gitmap/cmd || true)

if [ -z "$HITS" ]; then
  echo "OK: no bare 'fmt.Fprintln(os.Stderr, err)' in gitmap/cmd/"
  exit 0
fi

echo "FAIL: bare error prints found in gitmap/cmd/ — use cliexit.Reportf/Fail" >&2
echo "" >&2
while IFS= read -r line; do
  FILE=$(echo "$line" | cut -d: -f1)
  LINE=$(echo "$line" | cut -d: -f2)
  echo "::error file=${FILE},line=${LINE}::[bare-err] use cliexit.Reportf(cmd, op, subject, err) instead" >&2
done <<EOF
$HITS
EOF
exit 1
