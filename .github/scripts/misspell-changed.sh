#!/usr/bin/env bash
# misspell-changed.sh — Run pinned `misspell -locale US` on the set of
# files added or modified in the current PR / push, with include/exclude
# globs loaded from gitmap/data/config.json (misspell.include /
# misspell.exclude). Emits per-finding GitHub annotations so each typo
# is clickable in the PR "Files changed" view, and exits 1 on any
# finding (preserving the previous `-error` hard-gate behavior).
#
# Required env (all set by .github/workflows/ci.yml):
#   GH_EVENT_NAME   — "pull_request" or "push" (etc.)
#   GH_BASE_REF     — for PRs, the target branch (used to compute diff base)
#
# This is the body of the spell-check job, lifted out of ci.yml so the
# YAML stays readable. Logic unchanged from the previous inline version.
# See .lovable/plan.md (Step 3) for the extraction rationale.

set -uo pipefail

readonly CONFIG="gitmap/data/config.json"
readonly EVENT_NAME="${GH_EVENT_NAME:-push}"
readonly BASE_REF="${GH_BASE_REF:-}"

# Per Core CI JSON parsing rule: jq → python3 → safe defaults. Each
# helper prints one pattern per line for the lists.
read_patterns() {
  local key="$1"
  if command -v jq >/dev/null 2>&1 && [ -f "$CONFIG" ]; then
    jq -r --arg k "$key" '.misspell[$k] // [] | .[]' "$CONFIG" 2>/dev/null || true
  elif command -v python3 >/dev/null 2>&1 && [ -f "$CONFIG" ]; then
    python3 -c "import json; d=json.load(open('$CONFIG')).get('misspell',{}).get('$1',[]); print('\n'.join(d))" 2>/dev/null || true
  fi
}

mapfile -t excludes < <(read_patterns exclude)
mapfile -t includes < <(read_patterns include)

if [ "${#excludes[@]}" -eq 0 ]; then
  # Hard-coded fallback (matches pre-config behavior exactly).
  excludes=(
    "*.png" "*.jpg" "*.jpeg" "*.gif" "*.ico" "*.svg" "*.webp"
    "*.zip" "*.tar" "*.gz" "*.exe"
    "*/testdata/*" "*/golden/*"
    "*/.gitmap/release/*" "*/.gitmap/release-assets/*"
    "gitmap/completion/allcommands_generated.go"
    ".lovable/*"
  )
fi

echo "misspell config: ${#excludes[@]} exclude pattern(s), ${#includes[@]} include pattern(s)"

if [ "$EVENT_NAME" = "pull_request" ]; then
  base="origin/${BASE_REF}"
else
  base="HEAD~1"
fi

mapfile -t changed < <(git diff --name-only --diff-filter=AM "$base"...HEAD || true)
if [ "${#changed[@]}" -eq 0 ]; then
  echo "no changed files vs $base — nothing to scan"
  exit 0
fi

# Apply exclude (deny-list) then include (allow-list, if any).
matches_any() {
  local path="$1"; shift
  local pat
  for pat in "$@"; do
    # shellcheck disable=SC2254
    case "$path" in $pat) return 0 ;; esac
  done
  return 1
}

targets=()
for changedFile in "${changed[@]}"; do
  if matches_any "$changedFile" "${excludes[@]}"; then
    continue
  fi
  if [ "${#includes[@]}" -gt 0 ] && ! matches_any "$changedFile" "${includes[@]}"; then
    continue
  fi
  # Skip anything that no longer exists (e.g. deleted in a follow-up
  # commit on the same PR).
  [ -f "$changedFile" ] && targets+=("$changedFile")
done

if [ "${#targets[@]}" -eq 0 ]; then
  echo "no scannable changed files — nothing to scan"
  exit 0
fi

echo "Scanning ${#targets[@]} changed file(s) for US-English misspellings:"
printf '  %s\n' "${targets[@]}"

# Run misspell WITHOUT -error first so we capture every finding, then
# emit a GitHub line-level annotation per finding so each typo is
# clickable in the PR "Files changed" view. Misspell's default output
# format is:
#   path:line:col: "wrong" is a misspelling of "right"
# Exits non-zero if any finding was emitted.
report=$(mktemp)
misspell -locale US "${targets[@]}" | tee "$report" || true

findings=0
while IFS= read -r reportLine; do
  if [[ "$reportLine" =~ ^([^:]+):([0-9]+):([0-9]+):[[:space:]]*(.*)$ ]]; then
    file="${BASH_REMATCH[1]}"
    lno="${BASH_REMATCH[2]}"
    col="${BASH_REMATCH[3]}"
    msg="${BASH_REMATCH[4]}"
    # Strip CR and embedded ::/newlines that would break the workflow
    # command syntax. Annotation messages may not contain literal newlines.
    msg="${msg//$'\r'/}"
    msg="${msg//$'\n'/ }"
    echo "::error file=${file},line=${lno},col=${col}::[misspell] ${msg}"
    findings=$((findings + 1))
  fi
done < "$report"

if [ "$findings" -gt 0 ]; then
  echo "FAIL: ${findings} misspelling(s) found. See annotations above." >&2
  exit 1
fi

echo "OK: no misspellings in changed files."
