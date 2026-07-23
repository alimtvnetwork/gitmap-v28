#!/usr/bin/env bash
# Coverage-floor enforcer (#13).
#
# Inputs:
#   $1 = path to coverage.out (from `go test -coverprofile=...`)
#
# Floor logic:
#   - Default floor for any package is FLOOR_DEFAULT (70%).
#   - Packages listed in .github/coverage.floor override the default
#     with a per-package floor:  <import-path> <floor-percent>
#   - A package below its floor fails the build with a line per package.
#
# Existing packages that don't meet 70% should pin themselves at their
# current value (rounded down) and ratchet upward over time.
set -euo pipefail

COVER="${1:?usage: coverage-floor.sh coverage.out}"
FLOORS_FILE=".github/coverage.floor"
FLOOR_DEFAULT=70

if [[ ! -s "$COVER" ]]; then
  echo "coverage-floor: empty coverage profile at $COVER — skipping"
  exit 0
fi

declare -A FLOORS=()
if [[ -f "$FLOORS_FILE" ]]; then
  while read -r pkg pct; do
    [[ -z "$pkg" || "$pkg" =~ ^# ]] && continue
    FLOORS["$pkg"]="$pct"
  done < "$FLOORS_FILE"
fi

# Per-package coverage % via `go tool cover -func`.
TMP="$(mktemp)"
go tool cover -func="$COVER" \
  | awk '/\.go:/ { print }' \
  | awk -F'[/ \t:]+' '{
      # collapse to import-path by stripping the file basename + line range.
      n=split($0, a, "\t");
      file=a[1]; pct=a[n]; gsub("%", "", pct);
      sub("/[^/]+\\.go:[0-9]+$", "", file);
      pkgs[file] += pct; counts[file]++;
    }
    END {
      for (p in pkgs) printf "%s %.1f\n", p, pkgs[p]/counts[p];
    }' > "$TMP"

fail=0
while read -r pkg avg; do
  floor="${FLOORS[$pkg]:-$FLOOR_DEFAULT}"
  awk -v a="$avg" -v f="$floor" 'BEGIN { if (a+0 < f+0) exit 1 }' \
    || { echo "coverage-floor: $pkg below floor (avg=$avg%, floor=$floor%)"; fail=1; }
done < "$TMP"

rm -f "$TMP"
exit "$fail"
