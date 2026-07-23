#!/usr/bin/env bash
# check-changelog-version-sync.sh — CI gate that fails when the
# version pinned in gitmap/constants/constants.go is not documented
# in CHANGELOG.md.
#
# Motivation: version bumps have historically drifted from the
# changelog (see the v4.28 → v4.30 cycle documented in project
# memory). Every constants.Version = "X.Y.Z" MUST have a matching
# `## vX.Y.Z` heading in CHANGELOG.md. This script enforces that
# invariant in CI so the drift cannot land on main.
#
# Exit codes:
#   0 — version present in CHANGELOG.md
#   1 — version missing or constants.go unreadable
#   2 — bad invocation

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CONSTANTS_FILE="$REPO_ROOT/gitmap/constants/constants.go"
CHANGELOG_FILE="$REPO_ROOT/CHANGELOG.md"

if [[ ! -f "$CONSTANTS_FILE" ]]; then
  echo "✗ constants.go not found at $CONSTANTS_FILE" >&2
  exit 1
fi
if [[ ! -f "$CHANGELOG_FILE" ]]; then
  echo "✗ CHANGELOG.md not found at $CHANGELOG_FILE" >&2
  exit 1
fi

# Extract the pinned version literal. Matches:
#   const Version = "6.74.0"
VERSION=$(grep -E '^const Version = "[0-9]+\.[0-9]+\.[0-9]+"' "$CONSTANTS_FILE" \
  | head -1 \
  | sed -E 's/^const Version = "([^"]+)"/\1/')

if [[ -z "$VERSION" ]]; then
  echo "✗ Could not parse Version from $CONSTANTS_FILE" >&2
  echo "  Expected line: const Version = \"X.Y.Z\"" >&2
  exit 1
fi

echo "→ constants.Version = $VERSION"

# Look for a matching heading. Accept `## vX.Y.Z` with an optional
# suffix (date, tag, etc): "## v6.74.0 — 2026-07-01" is fine.
HEADING_RE="^##[[:space:]]+v${VERSION//./\\.}([[:space:]]|$)"

if ! grep -Eq "$HEADING_RE" "$CHANGELOG_FILE"; then
  echo "" >&2
  echo "✗ CHANGELOG drift: constants.Version is $VERSION but no matching" >&2
  echo "  '## v$VERSION' heading exists in CHANGELOG.md." >&2
  echo "" >&2
  echo "  Fix: add a '## v$VERSION' section to $CHANGELOG_FILE describing" >&2
  echo "  the release, then re-run this check." >&2
  exit 1
fi

echo "✓ CHANGELOG.md has entry for v$VERSION"
