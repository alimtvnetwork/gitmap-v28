#!/bin/bash
# ============================================================
# Axios Version Safeguard
# ============================================================
# Validates that Axios is pinned to an approved safe version
# and not using any range symbols (^, ~, >=, *).
#
# Blocked versions: 1.14.1, 0.30.4
# Approved versions: 1.14.0, 0.30.3
#
# Usage: bash linter-scripts/check-axios-version.sh
# Add to CI pipeline for automated enforcement.
# ============================================================

set -euo pipefail

BLOCKED_VERSIONS=("1.14.1" "0.30.4")
APPROVED_VERSIONS=("1.14.0" "0.30.3")

# Extract Axios version from package.json
CURRENT=$(node -e "
  const pkg = require('./package.json');
  const v = (pkg.dependencies && pkg.dependencies.axios) || (pkg.devDependencies && pkg.devDependencies.axios) || 'NOT_FOUND';
  console.log(v);
")

echo "📦 Axios version in package.json: $CURRENT"

# Check: not found
if [[ "$CURRENT" == "NOT_FOUND" ]]; then
  echo "⚠️  Axios is not declared in package.json"
  exit 0
fi

# Check: range symbols
if [[ "$CURRENT" == ^* ]] || [[ "$CURRENT" == ~* ]] || [[ "$CURRENT" == ">="* ]] || [[ "$CURRENT" == "*" ]] || [[ "$CURRENT" == "latest" ]]; then
  echo "❌ FAIL: Axios version uses a range symbol or tag: $CURRENT"
  echo "   Fix: Use an exact version like \"axios\": \"1.14.0\""
  exit 1
fi

# Check: blocked versions
for blocked in "${BLOCKED_VERSIONS[@]}"; do
  if [[ "$CURRENT" == "$blocked" ]]; then
    echo "❌ FAIL: Axios version $CURRENT is BLOCKED (known security vulnerability)"
    echo "   Approved versions: ${APPROVED_VERSIONS[*]}"
    exit 1
  fi
done

# Check: approved versions
IS_APPROVED=false
for approved in "${APPROVED_VERSIONS[@]}"; do
  if [[ "$CURRENT" == "$approved" ]]; then
    IS_APPROVED=true
    break
  fi
done

if [[ "$IS_APPROVED" == true ]]; then
  echo "✅ PASS: Axios $CURRENT is an approved safe version"
else
  echo "⚠️  WARNING: Axios $CURRENT is not in the approved list (${APPROVED_VERSIONS[*]})"
  echo "   This version has not been verified. Please review spec/01-app/axios-version-control/"
  exit 1
fi
