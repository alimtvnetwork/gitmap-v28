#!/usr/bin/env bash
# Hard-gate CI step for the `misspell` linter.
#
# Why a separate script (not piggy-backing on lint-baseline-diff)?
#   - The baseline-diff job is intentionally permissive: it gates on
#     ANY new finding from ANY enabled linter, and a noisy unrelated
#     class (e.g. a flood of typecheck `r.Mode undefined` errors with
#     no Pos.Filename attached) can produce path-less annotations
#     that drown out a genuine misspell regression in the PR UI.
#   - This script:
#       * runs only the misspell linter (--disable-all + --enable=misspell)
#       * diffs current vs. baseline so existing misspells in legacy
#         files don't block PRs that don't touch them
#       * REFUSES to emit any annotation that lacks a file path —
#         every reported finding is guaranteed to be `path:line:col`,
#         so stale duplicate errors with no location can never mask
#         a real failure.
#       * exits non-zero only when at least one NEW misspell with a
#         valid file path is present.
#
# Inputs (env or args):
#   $1            — directory to lint (default "gitmap")
#   $BASELINE     — path to baseline JSON report (optional; missing/empty
#                   means "seeding mode": warn, don't fail)
#   $CURRENT_OUT  — where to write the current JSON report
#                   (default /tmp/lint-misspell-current/report.json)
#
# Exit codes:
#   0 — no new misspell findings (or seeding mode)
#   1 — at least one new misspell finding with a full file path
#   2 — toolchain missing / unrecoverable error

set -euo pipefail

LINT_DIR="${1:-gitmap}"
CURRENT_OUT="${CURRENT_OUT:-/tmp/lint-misspell-current/report.json}"
BASELINE="${BASELINE:-}"

if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "ERROR: golangci-lint not on PATH" >&2
  exit 2
fi
if ! command -v jq >/dev/null 2>&1; then
  echo "ERROR: jq not on PATH" >&2
  exit 2
fi

mkdir -p "$(dirname "$CURRENT_OUT")"

# --no-config: ignore repo .golangci.yml so its excludes can't mask
#   a misspell. We're enforcing a focused floor here.
# --disable-all + --enable=misspell: only the misspell analyzer runs,
#   so unrelated linters (typecheck, gosec, etc.) cannot leak path-less
#   errors into this report.
# --issues-exit-code=0: never let the linter exit drive this script —
#   we parse JSON ourselves so we can enforce the "must have file path"
#   contract before emitting any annotation.
(
  cd "$LINT_DIR"
  golangci-lint run \
    --no-config \
    --disable-all \
    --enable=misspell \
    --timeout=5m \
    --issues-exit-code=0 \
    --out-format=json \
    ./... > "$CURRENT_OUT"
)

# Normalize a JSON report to a stable set of "file|line|text" keys,
# refusing entries with empty file paths. Emits one key per line.
normalize() {
  local path="$1"
  if [ -z "$path" ] || [ ! -s "$path" ]; then
    return 0
  fi
  jq -r '
    .Issues // []
    | map(select(.FromLinter == "misspell"))
    | map(select((.Pos.Filename // "") | length > 0))
    | .[]
    | "\(.Pos.Filename)|\(.Pos.Line)|\(.Text)"
  ' "$path"
}

CURRENT_KEYS=$(normalize "$CURRENT_OUT" | sort -u)
BASELINE_KEYS=""
SEEDING="false"
if [ -n "$BASELINE" ] && [ -s "$BASELINE" ]; then
  BASELINE_KEYS=$(normalize "$BASELINE" | sort -u)
else
  SEEDING="true"
fi

# Findings present in current but not in baseline.
NEW_KEYS=$(comm -23 \
  <(printf '%s\n' "$CURRENT_KEYS" | sed '/^$/d') \
  <(printf '%s\n' "$BASELINE_KEYS" | sed '/^$/d'))

NEW_COUNT=$(printf '%s\n' "$NEW_KEYS" | sed '/^$/d' | wc -l | tr -d ' ')

echo "========================================================================"
echo "  MISSPELL DIFF (hard-gate, full-path only)"
echo "========================================================================"
echo "  current  : $CURRENT_OUT"
echo "  baseline : ${BASELINE:-<none — seeding mode>}"
echo "  + NEW    : $NEW_COUNT"
echo "========================================================================"

if [ "$NEW_COUNT" = "0" ]; then
  echo "OK: no new misspell findings."
  exit 0
fi

# Re-emit each new finding with a guaranteed full path. The JSON
# lookup re-reads the report so the annotation carries the column
# too — needed for the GitHub Actions PR-files view to underline the
# exact word.
while IFS='|' read -r FILE LINE TEXT; do
  [ -z "$FILE" ] && continue
  COL=$(jq -r --arg f "$FILE" --argjson l "$LINE" --arg t "$TEXT" '
    .Issues // []
    | map(select(.FromLinter == "misspell"
        and .Pos.Filename == $f
        and .Pos.Line == $l
        and .Text == $t))
    | .[0].Pos.Column // 1
  ' "$CURRENT_OUT")
  if [ "$SEEDING" = "true" ]; then
    echo "::warning file=${FILE},line=${LINE},col=${COL}::[misspell] ${TEXT} (seeding baseline)"
  else
    echo "::error file=${FILE},line=${LINE},col=${COL}::[misspell] ${TEXT} (NEW vs baseline)"
  fi
done <<EOF
$NEW_KEYS
EOF

if [ "$SEEDING" = "true" ]; then
  echo "Seeding mode — not failing the build." >&2
  exit 0
fi

echo "" >&2
echo "FAIL: $NEW_COUNT new misspell finding(s). Fix the spellings above." >&2
exit 1
