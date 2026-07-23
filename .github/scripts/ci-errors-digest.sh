#!/usr/bin/env bash
# ci-errors-digest.sh — Consolidate every CI failure signal into one
# contiguous, copy-pasteable block printed to stdout.
#
# Usage: bash .github/scripts/ci-errors-digest.sh <artifacts-root>
#
# Reads (read-only) any of these subdirectories that exist under
# <artifacts-root>:
#   test-results-*/test-output.txt        (per-suite go test logs)
#   e2e-output-*/...                      (per-OS e2e logs)
#   full-suite-outputs/test-output.txt    (full-suite go test log)
#   full-suite-outputs/lint-output.txt    (golangci-lint strict log)
#   golangci-lint-report/*.json           (lint baseline-diff JSON)
#   lint-suggestions/*.txt                (actionable suggestions)
#
# Writes NOTHING to disk. The full digest is printed between two
# sentinel lines so the user can copy a single contiguous block and
# paste it back to the assistant:
#
#   ===BEGIN CI ERROR DIGEST===
#   ...
#   ===END CI ERROR DIGEST===
#
# Exit code is always 0 — this job is informational. The actual gate
# remains the upstream test/lint jobs.

set -uo pipefail

readonly ROOT="${1:-results}"

# Caps to keep the digest pasteable. Per-section line caps; the whole
# digest is also clipped to MAX_TOTAL_LINES at the end.
readonly MAX_PER_SECTION=80
readonly MAX_TOTAL_LINES=2000

# Buffer the entire digest in memory so the sentinels and the
# truncation tail can be emitted as one contiguous block (no
# interleaved `::group::` collapsibles inside the sentinels).
buf=""
emit() { buf+="$1"$'\n'; }

emit "===BEGIN CI ERROR DIGEST==="
emit "Repo:   ${GITHUB_REPOSITORY:-unknown}"
emit "SHA:    ${GITHUB_SHA:-unknown}"
emit "Run:    ${GITHUB_SERVER_URL:-https://github.com}/${GITHUB_REPOSITORY:-x/x}/actions/runs/${GITHUB_RUN_ID:-0}"
emit "Branch: ${GITHUB_REF_NAME:-unknown}"
emit ""

totalErrors=0

# ── Helper: emit a section if matching lines exist ────────────────
section() {
  local title="$1" file="$2" pattern="$3"
  [ -f "$file" ] || return 0
  local hits
  hits=$(grep -E "$pattern" "$file" 2>/dev/null | head -n "$MAX_PER_SECTION" || true)
  [ -n "$hits" ] || return 0
  emit "── ${title} ──"
  emit "(source: ${file#"$ROOT/"})"
  while IFS= read -r line; do emit "  $line"; done <<< "$hits"
  emit ""
  local count
  count=$(grep -cE "$pattern" "$file" 2>/dev/null || echo 0)
  totalErrors=$((totalErrors + count))
}

# ── 1. Per-suite go test failures with assertion context ──────────
for dir in "$ROOT"/test-results-* "$ROOT"/e2e-output-* "$ROOT"/full-suite-outputs; do
  [ -d "$dir" ] || continue
  for log in "$dir"/test-output.txt "$dir"/*.log; do
    [ -f "$log" ] || continue
    fails=$(grep -cE '^--- FAIL:' "$log" 2>/dev/null || echo 0)
    [ "$fails" -gt 0 ] || continue
    suite=$(basename "$dir" | sed 's/^test-results-//; s/^e2e-output-//')
    emit "── TEST FAILURES — ${suite} (${fails} failed) ──"
    # For each failing test, dump its name + the body lines that
    # carry .go:NN:CC or expected/got/Error markers. Capture awk
    # output to a variable first — piping into `while read | emit`
    # would run emit in a subshell and silently drop the appends.
    body=$(awk '
      /^=== RUN[[:space:]]/   { name=$3; body=""; capturing=1; next }
      /^--- FAIL:/            {
                                 if (capturing) {
                                   print "  --- FAIL: " name
                                   if (body != "") printf "%s", body
                                 }
                                 capturing=0; body=""
                                 next
                               }
      /^--- PASS:/            { capturing=0; body=""; next }
      capturing && /\.go:[0-9]+:/                                { body=body "    " $0 "\n" }
      capturing && /(expected|got:|want:|Error:|panic:|undefined|mismatch)/ { body=body "    " $0 "\n" }
    ' "$log" | head -n "$MAX_PER_SECTION")
    while IFS= read -r line; do emit "$line"; done <<< "$body"
    emit ""
    totalErrors=$((totalErrors + fails))
  done
done

# ── 2. golangci-lint strict findings (full-suite) ─────────────────
section "GOLANGCI-LINT (strict)" \
  "$ROOT/full-suite-outputs/lint-output.txt" \
  '^[^[:space:]].+:[0-9]+:[0-9]+:'

# ── 3. golangci-lint baseline-diff JSON (NEW findings only) ───────
for f in "$ROOT"/golangci-lint-report/*.json; do
  [ -f "$f" ] || continue
  if command -v jq >/dev/null 2>&1; then
    new=$(jq -r '.Issues[]? | "  \(.Pos.Filename):\(.Pos.Line):\(.Pos.Column): [\(.FromLinter)] \(.Text)"' "$f" 2>/dev/null | head -n "$MAX_PER_SECTION")
    if [ -n "$new" ]; then
      emit "── GOLANGCI-LINT (baseline-diff, NEW findings) ──"
      while IFS= read -r line; do emit "$line"; done <<< "$new"
      emit ""
      totalErrors=$((totalErrors + $(printf '%s\n' "$new" | wc -l)))
    fi
  fi
done

# ── 4. Actionable lint suggestions (if produced) ──────────────────
for f in "$ROOT"/lint-suggestions/*.txt "$ROOT"/lint-suggestions/*.md; do
  [ -f "$f" ] || continue
  body=$(head -n "$MAX_PER_SECTION" "$f")
  [ -n "$body" ] || continue
  emit "── LINT SUGGESTIONS (${f##*/}) ──"
  while IFS= read -r line; do emit "  $line"; done <<< "$body"
  emit ""
done

# ── 5. Build / vet / compile errors ───────────────────────────────
for log in "$ROOT"/full-suite-outputs/test-output.txt "$ROOT"/test-results-*/test-output.txt; do
  [ -f "$log" ] || continue
  build=$(grep -E '^(FAIL|# )[[:space:]]|cannot find|undefined:|cannot use' "$log" 2>/dev/null | grep -v '^FAIL[[:space:]]\+[a-z]' | head -n 40 || true)
  [ -n "$build" ] || continue
  emit "── BUILD / VET / COMPILE — ${log#"$ROOT/"} ──"
  while IFS= read -r line; do emit "  $line"; done <<< "$build"
  emit ""
done

if [ "$totalErrors" -eq 0 ]; then
  emit "(no errors detected across collected artifacts)"
fi

emit "===END CI ERROR DIGEST==="

# Clip to MAX_TOTAL_LINES so the digest stays pasteable even on
# catastrophic failures. Truncation marker is inserted BEFORE the
# end sentinel so the closing fence is always present.
lines=$(printf '%s' "$buf" | wc -l)
if [ "$lines" -gt "$MAX_TOTAL_LINES" ]; then
  head=$(printf '%s' "$buf" | head -n $((MAX_TOTAL_LINES - 2)))
  printf '%s\n' "$head"
  printf '... [digest truncated at %d lines — download artifacts for full output] ...\n' "$MAX_TOTAL_LINES"
  printf '===END CI ERROR DIGEST===\n'
else
  printf '%s' "$buf"
fi
