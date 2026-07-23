#!/usr/bin/env bash
set -uo pipefail

# test-summary.sh — Aggregates test results from CI matrix jobs.
# Produces a detailed, copy-pasteable failure report with test names and reasons.
#
# Usage: bash .github/scripts/test-summary.sh <results-dir>

RESULTS_DIR="${1:?Usage: test-summary.sh <results-dir>}"

echo "========================================="
echo "  ALL TEST RESULTS"
echo "========================================="

overall=0
all_failures=""

for dir in "$RESULTS_DIR"/test-results-*; do
  suite=$(basename "$dir" | sed 's/test-results-//')
  file="$dir/test-output.txt"

  if [ ! -f "$file" ]; then
    continue
  fi

  pass=$(grep -c '^--- PASS:' "$file" || true)
  fail=$(grep -c '^--- FAIL:' "$file" || true)

  if [ "$fail" -gt 0 ]; then
    echo ""
    echo "❌ $suite: $fail failed, $pass passed"
    overall=1

    # Collect detailed failure info for this suite
    suite_failures=""
    while IFS= read -r fail_line; do
      test_name=$(echo "$fail_line" | sed 's/^--- FAIL: //' | sed 's/ (.*$//')

      # Extract the failure reason: capture lines between "=== RUN <test>" and "--- FAIL: <test>"
      # Look for the actual assertion/error messages
      reason=$(awk -v test="$test_name" '
        $0 ~ "=== RUN[[:space:]]+" test { capturing=1; next }
        $0 ~ "--- FAIL:[[:space:]]+" test { capturing=0 }
        capturing && /\.go:[0-9]+:/ { print "        " $0 }
        capturing && /(expected|got|Error|FAIL|panic|undefined|mismatch)/ && !/=== RUN/ { print "        " $0 }
      ' "$file" | head -10)

      suite_failures="${suite_failures}    --- FAIL: ${test_name}
"
      if [ -n "$reason" ]; then
        suite_failures="${suite_failures}${reason}
"
      fi
    done < <(grep '^--- FAIL:' "$file")

    all_failures="${all_failures}
-----------------------------------------
  Suite: ${suite} (${fail} failed)
-----------------------------------------
${suite_failures}"
  else
    echo "✅ $suite: $pass passed"
  fi
done

echo ""
echo "========================================="

if [ "$overall" -ne 0 ]; then
  echo ""
  echo "========================================="
  echo "  FAILURE REPORT (copy-paste ready)"
  echo "========================================="
  echo "$all_failures"
  echo "========================================="
  echo ""
  echo "::error::Some test suites failed — see failure report above."
  exit 1
fi

echo "All test suites passed."
