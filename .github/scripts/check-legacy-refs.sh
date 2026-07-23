#!/usr/bin/env bash
# check-legacy-refs.sh
#
# Fails the build if any forbidden legacy version strings remain anywhere
# in the repository. This is the CI counterpart to `gitmap audit-legacy`
# (gitmap/cmd/auditlegacy.go) — it runs without needing the Go toolchain
# so it can gate every CI run, including remixes that haven't built yet.
#
# Default forbidden patterns: gitmap-v27, gitmap-v27, gitmap-v27
# Override via env: LEGACY_PATTERN="gitmap-v[567]\\b|old-org-name"
#
# Exit codes:
#   0 — no matches; repo is clean
#   1 — at least one match; offending file:line printed to stderr
#   2 — internal error (missing tool, bad regex)

set -euo pipefail

PATTERN="${LEGACY_PATTERN:-gitmap-v[567]\\b}"
ROOT="${1:-.}"

# Directories we never want to scan: VCS metadata, generated artifacts,
# vendored deps, and the .gitmap state dir which contains release JSONs
# that may legitimately reference historical version names.
EXCLUDE_DIRS=(
  ".git"
  "node_modules"
  "dist"
  "build"
  "bin"
  ".next"
  ".gitmap"
  "vendor"
  "coverage"
)

# File globs that are binary or otherwise meaningless to grep.
EXCLUDE_GLOBS=(
  "*.png" "*.jpg" "*.jpeg" "*.gif" "*.webp" "*.ico"
  "*.pdf" "*.zip" "*.gz" "*.tar"
  "*.exe" "*.dll" "*.so" "*.dylib" "*.bin"
  "*.db" "*.sqlite"
  "*.woff" "*.woff2" "*.ttf"
)

if ! command -v grep >/dev/null 2>&1; then
  echo "::error::grep not found" >&2
  exit 2
fi


# -H forces grep to always prefix matches with the filename, even when
# scanning a single-file ROOT. Without it, single-file scans emit just
# `lineno:text` and the file:line:text parsing below would mis-attribute.
GREP_ARGS=(-RHInE)
for d in "${EXCLUDE_DIRS[@]}"; do
  GREP_ARGS+=(--exclude-dir="$d")
done
for g in "${EXCLUDE_GLOBS[@]}"; do
  GREP_ARGS+=(--exclude="$g")
done
# Exclude this script itself + the spec doc that legitimately documents
# the rename to avoid self-matching.
GREP_ARGS+=(--exclude="check-legacy-refs.sh")

echo "  [legacy-refs] scanning '$ROOT' for pattern: $PATTERN"

# `|| true` lets us inspect the result without -e killing the script on
# the "no match" exit-code-1 from grep. Then strip lines carrying the
# whitelist marker so docs/comments that legitimately reference the old
# names (e.g. the audit-legacy command's own help text) don't self-trip
# the guard.
RAW="$(grep "${GREP_ARGS[@]}" "$PATTERN" "$ROOT" 2>/dev/null || true)"
MATCHES="$(printf '%s\n' "$RAW" | grep -v 'gitmap-legacy-ref-allow' || true)"

if [ -z "$MATCHES" ]; then
  echo "  [legacy-refs] OK — no forbidden legacy refs found."
  exit 0
fi

# Count and pretty-print for the GitHub Actions log.
MATCH_COUNT="$(printf '%s\n' "$MATCHES" | wc -l | tr -d ' ')"
FILES="$(printf '%s\n' "$MATCHES" | cut -d: -f1 | sort -u)"
FILE_COUNT="$(printf '%s\n' "$FILES" | wc -l | tr -d ' ')"

echo "::error::Found $MATCH_COUNT legacy reference(s) across $FILE_COUNT file(s) matching /$PATTERN/"
echo ""
echo "  Offending files:"
printf '%s\n' "$FILES" | sed 's/^/    - /'
echo ""
echo "  Match details (file:line:text):"
printf '%s\n' "$MATCHES" | sed 's/^/    /'
echo ""
echo "  To audit locally:"
echo "    gitmap audit-legacy"
echo "  Or scope a check:"
echo "    LEGACY_PATTERN='$PATTERN' bash .github/scripts/check-legacy-refs.sh ."

# Per-file annotations so GitHub surfaces them inline in the PR diff.
while IFS= read -r line; do
  file="$(printf '%s' "$line" | cut -d: -f1)"
  lineno="$(printf '%s' "$line" | cut -d: -f2)"
  text="$(printf '%s' "$line" | cut -d: -f3-)"
  echo "::error file=$file,line=$lineno::Legacy ref in $file: $text"
done <<< "$MATCHES"

exit 1
