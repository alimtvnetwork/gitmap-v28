#!/usr/bin/env bash
# fix-repo-gofmt-audit.sh — Targeted pre-gate that fires BEFORE the
# global `gofmt -l .` check. Detects PRs that look like a fix-repo bump
# (any .go file in the diff contains a `*-v<digits>` token, which is
# the canonical fix-repo rewrite signature) and runs gofmt on ONLY
# those files.
#
# Narrows the failure message from "some file isn't gofmt-clean" to
# "fix-repo left these files gofmt-dirty — re-run with the v4.8.0+
# binary or invoke fix-repo.ps1 v4.9.0+", which is actionable in a
# single read. Pure read-only — never writes back, never auto-commits;
# CI auto-fixes would mask the upstream bug. See
# .lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md.
#
# Required env:
#   GH_PR_BASE_SHA  — github.event.pull_request.base.sha (or empty on push;
#                     script falls back to HEAD~1)
#
# Run from the gitmap/ working directory.

set -uo pipefail

base="${GH_PR_BASE_SHA:-HEAD~1}"
[ -z "$base" ] && base="HEAD~1"

changed=$(git -C .. diff --name-only --diff-filter=AM "$base" -- 'gitmap/**/*.go' 2>/dev/null | sed 's|^gitmap/||' || true)
if [ -z "$changed" ]; then
  echo "fix-repo audit: no .go files changed vs $base — skipping."
  exit 0
fi

# A fix-repo bump leaves `<base>-v<N>` literals scattered across the
# changed files. We grep for that signature and only audit the matching
# subset so a plain refactor PR doesn't trip this targeted gate.
affected=""
for f in $changed; do
  [ -f "$f" ] || continue
  if grep -Eq '[A-Za-z0-9._-]+-v[0-9]+' "$f"; then
    affected="$affected $f"
  fi
done
affected=$(echo "$affected" | xargs -n1 echo | sort -u | xargs)
if [ -z "$affected" ]; then
  echo "fix-repo audit: no fix-repo-shaped tokens in changed .go files — skipping."
  exit 0
fi

echo "fix-repo audit: scanning $(echo "$affected" | wc -w) affected .go file(s)…"
# shellcheck disable=SC2086
dirty=$(gofmt -l $affected)
if [ -n "$dirty" ]; then
  echo "::error title=fix-repo gofmt regression::These .go files contain fix-repo-style version tokens AND are not gofmt-clean — the v4.8.0+ post-rewrite gofmt step did not run on them."
  echo "$dirty" | sed 's/^/  /'
  echo ""
  echo "----- gofmt -d (diff) -----"
  # shellcheck disable=SC2086
  gofmt -d $dirty
  echo "---------------------------"
  echo ""
  echo "Root cause: the fix-repo invocation that produced these changes did NOT auto-run gofmt."
  echo "Most likely you ran a pre-v4.8.0 'gitmap fix-repo' or pre-v4.9.0 'fix-repo.ps1'."
  echo ""
  echo "Fix locally (one of):"
  echo "  • cd gitmap && gofmt -w . && git add -u && git commit --amend --no-edit"
  echo "  • Re-run with the current binary: gitmap fix-repo --all (auto-formats since v4.8.0)"
  echo "  • Re-run the .ps1 helper: pwsh -File fix-repo.ps1 -All (auto-formats since v4.9.0)"
  exit 1
fi
echo "✅ fix-repo audit: all affected .go files are gofmt-clean."
