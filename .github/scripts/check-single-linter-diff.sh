#!/usr/bin/env bash
# Hard-gate CI step that runs ONE golangci-lint linter and fails when a
# NEW finding (with a guaranteed full file path) appears vs the baseline.
#
# Generalization of check-misspell-diff.sh — same contract, parameterized
# by linter name. Used to wire dedicated diff jobs for `misspell`,
# `gocritic`, and `exhaustive` so a regression in any of those classes
# can never be masked by stale path-less errors from other analyzers
# (e.g. typecheck `r.Mode undefined` annotations with no Pos.Filename).
#
# Why one focused linter per invocation?
#   - --disable-all + --enable=<one> guarantees the JSON report contains
#     only that analyzer's issues, so the path-required filter can't
#     accidentally drop or admit a finding from another linter.
#   - A failing class is reported in isolation: the PR annotation says
#     "[gocritic] appendAssign ..." with the exact file:line:col,
#     instead of being buried in a 400-line mixed report.
#   - The repo's .golangci.yml is bypassed (--no-config) so excludes
#     can't silence the floor we're enforcing here.
#
# Usage:
#   LINTER=misspell BASELINE=/tmp/x/baseline.json CURRENT_OUT=/tmp/x/cur.json \
#     bash .github/scripts/check-single-linter-diff.sh gitmap
#
# Inputs (env or args):
#   $1            — directory to lint (default "gitmap")
#   $LINTER       — REQUIRED; the single golangci-lint analyzer to enable
#   $BASELINE     — path to baseline JSON report (optional; missing/empty
#                   means "seeding mode": warn, don't fail)
#   $CURRENT_OUT  — where to write the current JSON report (REQUIRED)
#   $TEXT_FILTER  — OPTIONAL regex; when set, only findings whose .Text
#                   matches are kept (used to scope `gosec` to a single
#                   rule like G115). Applied uniformly to current AND
#                   baseline so the diff stays apples-to-apples.
#   $LABEL        — OPTIONAL display label for log/annotation banners
#                   (defaults to $LINTER; set to e.g. "gosec-G115" when
#                   TEXT_FILTER scopes a single rule).
#
# Exit codes:
#   0 — no new findings (or seeding mode)
#   1 — at least one new finding from $LINTER with a full file path
#   2 — toolchain missing / unrecoverable error / missing required env

set -euo pipefail

LINT_DIR="${1:-gitmap}"
LINTER="${LINTER:-}"
CURRENT_OUT="${CURRENT_OUT:-}"
BASELINE="${BASELINE:-}"
TEXT_FILTER="${TEXT_FILTER:-}"
LABEL="${LABEL:-$LINTER}"

if [ -z "$LINTER" ]; then
  echo "ERROR: LINTER env var is required (e.g. LINTER=gocritic)" >&2
  exit 2
fi
if [ -z "$CURRENT_OUT" ]; then
  echo "ERROR: CURRENT_OUT env var is required" >&2
  exit 2
fi
if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "ERROR: golangci-lint not on PATH" >&2
  exit 2
fi
if ! command -v jq >/dev/null 2>&1; then
  echo "ERROR: jq not on PATH" >&2
  exit 2
fi

mkdir -p "$(dirname "$CURRENT_OUT")"

# --no-config: ignore repo .golangci.yml so its excludes can't mask a
#   regression. We're enforcing a focused floor here.
# --disable-all + --enable=$LINTER: only the target analyzer runs, so
#   unrelated linters cannot leak path-less errors into this report.
# --issues-exit-code=0: never let the linter's exit drive this script —
#   we parse JSON ourselves and enforce the "must have file path"
#   contract before emitting any annotation.
(
  cd "$LINT_DIR"
  golangci-lint run \
    --no-config \
    --disable-all \
    --enable="$LINTER" \
    --timeout=5m \
    --issues-exit-code=0 \
    --out-format=json \
    ./... > "$CURRENT_OUT"
)

# Normalize a JSON report to a stable set of "file|line|text" keys,
# refusing entries with empty file paths. When TEXT_FILTER is set, only
# findings whose .Text matches the regex are kept (applied to BOTH
# current and baseline so the diff stays apples-to-apples). Emits one
# key per line.
normalize() {
  local path="$1"
  if [ -z "$path" ] || [ ! -s "$path" ]; then
    return 0
  fi
  jq -r --arg linter "$LINTER" --arg filter "$TEXT_FILTER" '
    .Issues // []
    | map(select(.FromLinter == $linter))
    | map(select((.Pos.Filename // "") | length > 0))
    | map(select($filter == "" or (.Text | test($filter))))
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
echo "  ${LABEL^^} DIFF (baseline-diff, full-path only)"
echo "========================================================================"
echo "  current  : $CURRENT_OUT"
echo "  baseline : ${BASELINE:-<none — seeding mode>}"
echo "  + NEW    : $NEW_COUNT"
echo "========================================================================"

if [ "$NEW_COUNT" = "0" ]; then
  echo "OK: no new $LABEL findings."
  exit 0
fi

# Re-emit each new finding with a guaranteed full path. The JSON
# lookup re-reads the report so the annotation carries the column
# too — needed for the GitHub Actions PR-files view to underline
# the exact location.
while IFS='|' read -r FILE LINE TEXT; do
  [ -z "$FILE" ] && continue
  COL=$(jq -r --arg f "$FILE" --argjson l "$LINE" --arg t "$TEXT" --arg linter "$LINTER" '
    .Issues // []
    | map(select(.FromLinter == $linter
        and .Pos.Filename == $f
        and .Pos.Line == $l
        and .Text == $t))
    | .[0].Pos.Column // 1
  ' "$CURRENT_OUT")
  if [ "$SEEDING" = "true" ]; then
    echo "::warning file=${FILE},line=${LINE},col=${COL}::[${LABEL}] ${TEXT} (seeding baseline)"
  else
    echo "::error file=${FILE},line=${LINE},col=${COL}::[${LABEL}] ${TEXT} (NEW vs baseline)"
  fi
done <<EOF
$NEW_KEYS
EOF

if [ "$SEEDING" = "true" ]; then
  echo "Seeding mode — not failing the build." >&2
  exit 0
fi

echo "" >&2
echo "FAIL: $NEW_COUNT new $LABEL finding(s). Fix the issues above." >&2
exit 1
