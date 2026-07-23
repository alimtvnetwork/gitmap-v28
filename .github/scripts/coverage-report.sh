#!/usr/bin/env bash
# coverage-report.sh — Merge per-package coverage profiles produced by the
# matrixed `test` jobs and print a per-package breakdown plus a total.
#
# Invoked from the test-summary job in .github/workflows/ci.yml after
# `actions/download-artifact` has unpacked every `test-results-*` bundle
# into the directory passed as $1. Output is plain text intended for the
# job's stdout (which becomes the Actions log).
#
# Inputs:
#   $1  results_dir  — directory containing test-results-*/coverage-*.out files
#                      (default: "results", matching the workflow step)
#
# Side effects:
#   Writes a merged "combined-coverage.out" inside the gitmap/ Go module
#   directory (cwd at call time) so `go tool cover -func` can read it.
#
# Exits 0 even when no coverage data is present (matrix legs may opt
# out of -coverprofile in fast-path mode); the calling step never gates
# on this.

set -uo pipefail

readonly RESULTS_DIR="${1:-results}"

echo ""
echo "========================================="
echo "  COVERAGE BY PACKAGE"
echo "========================================="

# Merge all coverage profiles. The 'mode:' header must appear exactly
# once at the top — strip it from every input file before appending.
echo "mode: atomic" > combined-coverage.out
for coverFile in "../${RESULTS_DIR}"/test-results-*/coverage-*.out; do
  if [ -f "$coverFile" ]; then
    grep -v '^mode:' "$coverFile" >> combined-coverage.out 2>/dev/null || true
  fi
done

if [ "$(wc -l < combined-coverage.out)" -le 1 ]; then
  echo "No coverage data collected."
  exit 0
fi

# Per-package breakdown. `go tool cover -func` emits one line per
# function plus a final "total:" line. We collapse functions to their
# parent package (everything before the last "/...") and average the
# percentages for a package-level view.
go tool cover -func=combined-coverage.out | \
  awk '
    /^total:/ { total=$NF; next }
    {
      pkg=$1; sub(/\/[^/]+$/, "", pkg)
      if (!(pkg in seen)) { seen[pkg]=1; order[++n]=pkg }
      split($NF, a, "%"); pct[pkg] += a[1]; cnt[pkg]++
    }
    END {
      for (i=1; i<=n; i++) {
        p = order[i]
        avg = pct[p] / cnt[p]
        printf "  %-55s %5.1f%%\n", p, avg
      }
      printf "\n  %-55s %s\n", "Total", total
    }
  '

echo "========================================="
