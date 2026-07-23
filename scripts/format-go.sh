#!/usr/bin/env bash
# scripts/format-go.sh — Auto-format Go files with gofmt.
#
# Modes:
#   ./scripts/format-go.sh                # format every .go file under gitmap/
#   ./scripts/format-go.sh --staged       # format only staged .go files (hook mode)
#   ./scripts/format-go.sh path/a.go ...  # format an explicit file list
#
# In --staged mode, files that were *fully* staged are re-added with `git add`
# after rewrite so the formatted version is what actually gets committed.
# Files with unstaged worktree changes (partial staging) are formatted but
# NOT re-added — re-adding would silently swallow the unstaged hunks. The
# script exits 1 in that case so the contributor can resolve manually.
#
# Exit codes:
#   0 — nothing to do, or all files formatted + safely re-staged
#   1 — partial-staging conflict (see message), or gofmt missing in --staged mode
#   2 — usage error

set -uo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"

if ! command -v gofmt &>/dev/null; then
  echo "⚠ gofmt not found — install Go toolchain (https://go.dev/dl/)" >&2
  # Fail hard in --staged (hook) mode so commits don't slip through unformatted.
  if [ "${1:-}" = "--staged" ]; then exit 1; fi
  exit 0
fi

mode="${1:-all}"

case "$mode" in
  --staged)
    FILES=$(git diff --cached --name-only --diff-filter=ACM -- '*.go' || true)
    ;;
  all)
    FILES=$(find "$REPO_ROOT/gitmap" -name '*.go' -type f 2>/dev/null || true)
    ;;
  --help|-h)
    sed -n '2,12p' "$0"; exit 0
    ;;
  -*)
    echo "✗ unknown flag: $mode" >&2; exit 2
    ;;
  *)
    FILES="$*"
    mode="explicit"
    ;;
esac

if [ -z "$FILES" ]; then
  echo "  (no .go files to format)"
  exit 0
fi

# Find which files actually need rewriting — keeps `git add` noise to a minimum.
DIRTY=$(echo "$FILES" | xargs gofmt -l 2>/dev/null || true)

if [ -z "$DIRTY" ]; then
  echo "  ✓ all Go files already gofmt-clean"
  exit 0
fi

echo "Formatting $(echo "$DIRTY" | wc -l | tr -d ' ') file(s):"
echo "$DIRTY" | sed 's/^/    /'
echo "$DIRTY" | xargs gofmt -w

# In --staged mode, re-add fully-staged files; warn about partial-staging.
if [ "$mode" = "--staged" ]; then
  CONFLICTS=""
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    # `git diff --name-only -- <f>` lists the file iff worktree differs from index.
    # Before our gofmt -w, that meant unstaged user edits; after, it's just our
    # formatting delta. To distinguish, check if the file appeared in the
    # unstaged diff *before* this script ran — we do that by checking if there
    # were any non-staged lines for it pre-rewrite. Simpler heuristic: compare
    # the staged blob to the worktree-content's gofmt'd form. If equal, the
    # only worktree drift is our own rewrite → safe to re-add.
    staged_fmt=$(git show ":$f" 2>/dev/null | gofmt 2>/dev/null || true)
    worktree=$(cat "$f" 2>/dev/null || true)
    if [ "$staged_fmt" = "$worktree" ]; then
      git add -- "$f"
    else
      CONFLICTS="$CONFLICTS\n    $f"
    fi
  done <<< "$DIRTY"

  if [ -n "$CONFLICTS" ]; then
    echo ""
    echo "✗ Partial-staging detected — these files have unstaged edits and" >&2
    printf "  were formatted but NOT re-staged:%b\n" "$CONFLICTS" >&2
    echo ""
    echo "  Resolve with:  git add <file>   (commits the formatted version)" >&2
    echo "             or: git stash && git add <file> && git stash pop" >&2
    exit 1
  fi
  echo "  ✓ formatted files re-staged"
fi
