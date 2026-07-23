#!/usr/bin/env bash
# check-cmd-naming.sh
#
# Guards the `cmd/` package against collision-prone helper names.
#
# Background: every file in gitmap/cmd shares one Go namespace, so two files
# both declaring `func invokeRelease(...)` produces a "redeclared in this
# block" build break (see history of releasescan.go vs releasealias.go).
#
# Rule: helpers using a generic verb prefix MUST be qualified with a domain
# word AFTER the prefix that narrows the scope (e.g. `invokeAliasRelease`,
# `runOnePullJob`, `runOneScanRelease`, `persistReleaseToDB`). Bare or
# generically-suffixed forms (`invokeRelease`, `runOne`, `persistAll`) are
# forbidden.
#
# Note: `execute<Cmd>` and `handle<Cmd>` are the canonical CLI router pattern
# (one per top-level command, e.g. `executeRelease`, `handleChangelogOpen`)
# and are NOT flagged here. Only `invoke`, `persist`, and `runOne` — which
# are *helper* verbs that frequently get reinvented across files — trigger.
#
# Forbidden patterns (regex, anchored at func declaration):
#   1. ^func +(invoke|persist) *\(            -> bare verb, no noun
#   2. ^func +runOne *\(                       -> bare runOne
#   3. ^func +(invoke|persist)(Release|Task|Job|Item|All|One|Cmd) *\(
#      -> verb + over-generic noun (the noun must be project-specific,
#         e.g. AliasRelease/ScanRelease/PullJob qualify the scope)
#
# Exit 1 on any violation; 0 otherwise.

set -euo pipefail

CMD_DIR="${1:-gitmap/cmd}"

if [ ! -d "$CMD_DIR" ]; then
  echo "::error::cmd directory not found: $CMD_DIR"
  exit 1
fi

# Patterns are written as extended regex (-E).
PATTERNS=(
  '^func +(invoke|persist) *\('
  '^func +runOne *\('
  '^func +(invoke|persist|runOne)(Release|Task|Job|Item|All|One|Cmd) *\('
)

violations=0
for pat in "${PATTERNS[@]}"; do
  # -H: print filename, -n: line number, -E: extended regex
  if matches=$(grep -HnE "$pat" "$CMD_DIR"/*.go 2>/dev/null); then
    if [ -n "$matches" ]; then
      echo "::error::Generic helper name detected (collision risk in single cmd/ namespace)."
      echo "::error::Pattern: $pat"
      echo "$matches" | while IFS= read -r line; do
        echo "::error::  $line"
      done
      echo ""
      violations=$((violations + 1))
    fi
  fi
done

if [ "$violations" -gt 0 ]; then
  echo ""
  echo "::error::Found $violations collision-prone naming pattern(s) in $CMD_DIR"
  echo "::error::Rename helpers with a domain prefix that narrows the scope, e.g.:"
  echo "::error::  invokeRelease       -> invokeAliasRelease"
  echo "::error::  runOne              -> runOnePullJob, runOneScanRelease"
  echo "::error::  persistAll          -> persistReleaseToDB"
  echo "::error::Note: executeXxx / handleXxx are reserved for the canonical"
  echo "::error::      command handler per top-level command and are allowed."
  exit 1
fi

echo "All cmd/ helper names are domain-qualified."
