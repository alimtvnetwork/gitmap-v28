#!/usr/bin/env bash
# check-no-golden-allow-leak.sh
#
# Fails CI if GITMAP_ALLOW_GOLDEN_UPDATE is *set* (not merely
# referenced) anywhere outside the goldenguard test suite.
#
# Why this exists: GITMAP_ALLOW_GOLDEN_UPDATE is the second of two
# defense-in-depth gate vars guarding fixture regeneration. The
# project policy is that **CI and shipped scripts must never set
# this variable** — only humans set it locally when they
# deliberately regenerate goldens. A leak in a workflow `env:`
# block, a Make target, a `.sh`/`.ps1` setup script, or a stray
# `os.Setenv(...)` outside the goldenguard package would silently
# unlock fixture rewrites on every CI run.
#
# Detection strategy (false-positive avoidance):
#   - We do NOT grep all files for the bare token. Every Go test in
#     this repo legitimately mentions the var in error messages and
#     comments; Markdown docs cite it dozens of times; the regoldens
#     command embeds it as a child-process env entry via
#     goldenguard.AllowUpdateEnv+"=...". Naïve grep produces ~100%
#     false positives.
#   - Instead we scan ONLY files that can actually export an env var
#     at runtime, and within each file type we apply patterns that
#     match the *assignment syntax* of that language.
#
# Forbidden patterns by file type:
#   *.sh / *.bash            (^|[^A-Z_])GITMAP_ALLOW_GOLDEN_UPDATE[[:space:]]*=
#                            export[[:space:]]+GITMAP_ALLOW_GOLDEN_UPDATE
#   *.ps1 / *.psm1           \$env:GITMAP_ALLOW_GOLDEN_UPDATE
#                            Set-Item.*GITMAP_ALLOW_GOLDEN_UPDATE
#   *.yml / *.yaml           ^[[:space:]]*GITMAP_ALLOW_GOLDEN_UPDATE[[:space:]]*:
#                            plus shell-assignment patterns inside run: blocks
#   Makefile, *.mk           (^|[^A-Z_])GITMAP_ALLOW_GOLDEN_UPDATE[[:space:]]*=
#                            export[[:space:]]+GITMAP_ALLOW_GOLDEN_UPDATE
#   Dockerfile*              ^[[:space:]]*ENV[[:space:]]+GITMAP_ALLOW_GOLDEN_UPDATE
#   *.go                     (os|t|tt|child)\.Setenv\([[:space:]]*"GITMAP_ALLOW_GOLDEN_UPDATE"
#                            \.Setenv\([[:space:]]*(goldenguard\.)?AllowUpdateEnv
#
# Whitelist: gitmap/goldenguard/ legitimately sets the var inside
# its own unit tests (t.Setenv) — that is the one package whose job
# is to test the gate itself.

set -euo pipefail

ALLOW_VAR="GITMAP_ALLOW_GOLDEN_UPDATE"
WHITELIST_PREFIX="gitmap/goldenguard/"

# Track violations across all checks so we can report the full set
# in one CI run instead of bailing on the first hit.
violations=0
tmp_report="$(mktemp)"
trap 'rm -f "$tmp_report"' EXIT

# report_match prints a GitHub-Actions-style error line and bumps
# the violation counter. file:line:content is the standard format.
report_match() {
  local file="$1" line="$2" content="$3" reason="$4"
  echo "::error file=${file},line=${line}::${reason}" | tee -a "$tmp_report"
  echo "::error::  ${file}:${line}: ${content}" | tee -a "$tmp_report"
  violations=$((violations + 1))
}

# is_whitelisted returns 0 (true) when the path is inside the
# goldenguard package, which owns the var and may legitimately
# t.Setenv it from its own unit tests.
is_whitelisted() {
  case "$1" in
    "${WHITELIST_PREFIX}"*) return 0 ;;
    *) return 1 ;;
  esac
}

# scan_with_pattern_per_ext greps a file with one ERE pattern and
# reports every hit. Skips whitelisted paths. Used by every file
# type below — keeps the per-extension loops short.
scan_with_pattern_per_ext() {
  local file="$1" pattern="$2" reason="$3"
  is_whitelisted "$file" && return 0
  while IFS=: read -r lineno content; do
    [ -z "$lineno" ] && continue
    report_match "$file" "$lineno" "$(printf '%s' "$content" | head -c 200)" "$reason"
  done < <(grep -nE "$pattern" "$file" 2>/dev/null || true)
}

# list_tracked emits NUL-separated tracked paths matching one or
# more shell globs. Falls back to `find -name <glob> -o -name ...`
# when not inside a git checkout (e.g. a tarball release rebuild).
list_tracked() {
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    git ls-files -z -- "$@"
    return
  fi
  buildFindAndRun "$@"
}

# buildFindAndRun assembles a `find . -type f \( -name g1 -o -name g2 ... \)`
# expression. Splitting it out keeps list_tracked's positive branch
# obvious and lets the find-builder respect the 15-line function cap.
buildFindAndRun() {
  local args=(. -type f '(')
  local first=1
  for glob in "$@"; do
    if [ "$first" -eq 1 ]; then
      first=0
    else
      args+=(-o)
    fi
    args+=(-name "$glob")
  done
  args+=(')' -print0)
  find "${args[@]}"
}


echo "▸ Scanning for ${ALLOW_VAR} assignments outside ${WHITELIST_PREFIX}..."

# 1. Shell scripts.
while IFS= read -r -d '' f; do
  scan_with_pattern_per_ext "$f" \
    '(^|[^A-Z_])GITMAP_ALLOW_GOLDEN_UPDATE[[:space:]]*=' \
    "Shell assignment of ${ALLOW_VAR} is forbidden in shipped scripts"
  scan_with_pattern_per_ext "$f" \
    'export[[:space:]]+GITMAP_ALLOW_GOLDEN_UPDATE' \
    "export of ${ALLOW_VAR} is forbidden in shipped scripts"
done < <(list_tracked '*.sh' '*.bash')

# 2. PowerShell.
while IFS= read -r -d '' f; do
  scan_with_pattern_per_ext "$f" \
    '\$env:GITMAP_ALLOW_GOLDEN_UPDATE' \
    "PowerShell \$env: assignment of ${ALLOW_VAR} is forbidden"
  scan_with_pattern_per_ext "$f" \
    'Set-Item[^\n]*GITMAP_ALLOW_GOLDEN_UPDATE' \
    "PowerShell Set-Item of ${ALLOW_VAR} is forbidden"
done < <(list_tracked '*.ps1' '*.psm1')

# 3. YAML (workflows + composite actions). Catches both top-level
# env: keys and shell assignments embedded in run: blocks.
while IFS= read -r -d '' f; do
  scan_with_pattern_per_ext "$f" \
    '^[[:space:]]*GITMAP_ALLOW_GOLDEN_UPDATE[[:space:]]*:' \
    "YAML env key ${ALLOW_VAR}: is forbidden in workflows/actions"
  scan_with_pattern_per_ext "$f" \
    '(^|[^A-Z_])GITMAP_ALLOW_GOLDEN_UPDATE[[:space:]]*=' \
    "Shell-style assignment of ${ALLOW_VAR} inside YAML run: block is forbidden"
  scan_with_pattern_per_ext "$f" \
    'export[[:space:]]+GITMAP_ALLOW_GOLDEN_UPDATE' \
    "export of ${ALLOW_VAR} inside YAML run: block is forbidden"
done < <(list_tracked '*.yml' '*.yaml')

# 4. Makefiles (and *.mk includes).
while IFS= read -r -d '' f; do
  scan_with_pattern_per_ext "$f" \
    '(^|[^A-Z_])GITMAP_ALLOW_GOLDEN_UPDATE[[:space:]]*=' \
    "Makefile assignment of ${ALLOW_VAR} is forbidden"
  scan_with_pattern_per_ext "$f" \
    'export[[:space:]]+GITMAP_ALLOW_GOLDEN_UPDATE' \
    "Makefile export of ${ALLOW_VAR} is forbidden"
done < <(list_tracked 'Makefile' '*.mk')

# 5. Dockerfiles.
while IFS= read -r -d '' f; do
  scan_with_pattern_per_ext "$f" \
    '^[[:space:]]*ENV[[:space:]]+GITMAP_ALLOW_GOLDEN_UPDATE' \
    "Dockerfile ENV ${ALLOW_VAR} is forbidden"
done < <(list_tracked 'Dockerfile' 'Dockerfile.*' '*.dockerfile')

# 6. Go runtime sets — Setenv variants. The goldenguard package
# itself is whitelisted (handled by is_whitelisted) because its own
# unit tests legitimately exercise the gate via t.Setenv. The
# regoldens command builds a child-process env slice via string
# concatenation (goldenguard.AllowUpdateEnv+"="+value) and is NOT
# matched by these patterns — that is intentional, the slice is
# scoped to the spawned `go test` process and never exported.
while IFS= read -r -d '' f; do
  scan_with_pattern_per_ext "$f" \
    '\.Setenv\([[:space:]]*"GITMAP_ALLOW_GOLDEN_UPDATE"' \
    "Setenv(\"${ALLOW_VAR}\", ...) outside goldenguard/ is forbidden"
  scan_with_pattern_per_ext "$f" \
    '\.Setenv\([[:space:]]*(goldenguard\.)?AllowUpdateEnv' \
    "Setenv(AllowUpdateEnv, ...) outside goldenguard/ is forbidden"
done < <(list_tracked '*.go')

if [ "$violations" -gt 0 ]; then
  echo ""
  echo "::error::❌ ${violations} ${ALLOW_VAR} leak(s) detected outside ${WHITELIST_PREFIX}"
  echo "::error::Policy: only humans regenerating goldens locally may set this var."
  echo "::error::See spec/05-coding-guidelines/21-golden-fixture-regeneration.md §6 (CI Posture)."
  exit 1
fi

echo "✅ No ${ALLOW_VAR} leaks detected outside ${WHITELIST_PREFIX}"
