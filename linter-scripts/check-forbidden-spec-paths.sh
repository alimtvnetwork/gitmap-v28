#!/usr/bin/env bash
# ============================================================
# Forbidden Spec Paths Guard
# ============================================================
# Fails CI on any of these violations:
#
#   1. Re-appearance of the deprecated, pre-consolidation update
#      folders under spec/:
#        - spec/14-generic-update/
#        - spec/15-self-update-app-update/
#      (Both were merged into spec/14-update/ on 2026-04-17.)
#
#   2. Any MERGE-PROPOSAL.md (any case variant) under spec/.
#      It was a transient planning doc and must not be committed.
#
#   3. Any uppercase-letter .md filename anywhere under spec/
#      or release-artifacts/. The lowercase-markdown standard
#      requires every Markdown filename to be lowercase
#      (e.g. readme.md, not readme.md or ReadMe.md).
#
# Usage:
#   bash linter-scripts/check-forbidden-spec-paths.sh
#
# Exit codes:
#   0  no violations
#   1  one or more violations present
# ============================================================

set -uo pipefail

SPEC_ROOT="spec"
RELEASE_ROOT="release-artifacts"
EXIT_CODE=0

echo "🔍 Checking for forbidden spec paths and uppercase .md filenames..."

# ── 1. Forbidden folders (re-split guards) ──────────────────
FORBIDDEN_DIRS=(
  "$SPEC_ROOT/14-generic-update"
  "$SPEC_ROOT/15-self-update-app-update"
)

for DIR in "${FORBIDDEN_DIRS[@]}"; do
  if [[ -e "$DIR" ]]; then
    echo "::error file=$DIR::Forbidden folder present: $DIR (merged into spec/14-update/, must not re-appear)"
    EXIT_CODE=1
  fi
done

# ── 2. Forbidden files (any-case MERGE-PROPOSAL.md) ─────────
if [[ -d "$SPEC_ROOT" ]]; then
  MERGE_PROPOSAL_HITS=$(find "$SPEC_ROOT" -type f -iname 'merge-proposal.md' 2>/dev/null || true)
  if [[ -n "$MERGE_PROPOSAL_HITS" ]]; then
    while IFS= read -r HIT; do
      echo "::error file=$HIT::Forbidden file: MERGE-PROPOSAL.md must not be committed under spec/"
    done <<< "$MERGE_PROPOSAL_HITS"
    EXIT_CODE=1
  fi
fi

# ── 3. Uppercase .md filenames under spec/ + release-artifacts/ ──
check_uppercase_md() {
  local root="$1"
  [[ -d "$root" ]] || return 0

  # Match any .md whose basename contains at least one uppercase letter.
  local hits
  hits=$(find "$root" -type f -name '*.md' 2>/dev/null \
    | awk -F/ '{ if ($NF ~ /[A-Z]/) print }' || true)

  if [[ -z "$hits" ]]; then
    return 0
  fi

  while IFS= read -r HIT; do
    [[ -z "$HIT" ]] && continue
    echo "::error file=$HIT::Uppercase letters in .md filename — rename to lowercase: $(basename "$HIT")"
    EXIT_CODE=1
  done <<< "$hits"
}

check_uppercase_md "$SPEC_ROOT"
check_uppercase_md "$RELEASE_ROOT"

# ── Summary ─────────────────────────────────────────────────
echo ""
if [[ "$EXIT_CODE" -eq 0 ]]; then
  echo "✅ No forbidden paths or uppercase .md filenames detected."
else
  echo "❌ Violations detected. See errors above."
  echo "   - Consolidated update home: spec/14-update/"
  echo "   - Markdown filenames must be all lowercase (e.g. readme.md)."
fi

exit "$EXIT_CODE"
