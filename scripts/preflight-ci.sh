#!/usr/bin/env bash
# preflight-ci.sh — Run the same test + lint suite as CI locally.
#
# Mirrors .github/workflows/ci.yml `full-suite-guard` job:
#   1. go test ./... -count=1                (every package, no cache)
#   2. golangci-lint run --max-issues-per-linter=0 --max-same-issues=0
#
# Versions are pinned to match CI exactly (see .lovable/memory/tech/dependency-management.md):
#   - golangci-lint v1.64.8
#   - govulncheck   v1.1.4 (run separately by vulncheck.yml)
#
# Usage:
#   ./scripts/preflight-ci.sh           # run both phases
#   ./scripts/preflight-ci.sh test      # tests only
#   ./scripts/preflight-ci.sh lint      # lint only
#
# Exit 0 = clean, exit 1 = failures (matches CI gate behavior).

set -uo pipefail

readonly GOLANGCI_LINT_VERSION="v1.64.8"
readonly GITMAP_DIR="gitmap"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PHASE="${1:-all}"

hasModule=true
if [ -d "$REPO_ROOT/$GITMAP_DIR" ]; then
  hasModule=true
else
  hasModule=false
fi

if [ "$hasModule" = "false" ]; then
  echo "✗ preflight-ci: '$GITMAP_DIR/' not found at $REPO_ROOT" >&2
  exit 1
fi

cd "$REPO_ROOT/$GITMAP_DIR" || exit 1

run_tests() {
  echo "=== [1/2] go test ./... (every package, no cache) ==="
  if ! go test ./... -count=1; then
    echo "" >&2
    echo "✗ preflight-ci: go test failed" >&2
    return 1
  fi
  echo "  ✓ tests passed"
  return 0
}

ensure_golangci_lint() {
  hasLinter=true
  if command -v golangci-lint >/dev/null 2>&1; then
    hasLinter=true
  else
    hasLinter=false
  fi

  if [ "$hasLinter" = "false" ]; then
    echo "✗ preflight-ci: golangci-lint not installed" >&2
    echo "  Install pinned version with:" >&2
    echo "    go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}" >&2
    return 1
  fi

  installedVersion="$(golangci-lint version --format short 2>/dev/null || echo "unknown")"
  isPinnedVersion=true
  if [ "$installedVersion" = "${GOLANGCI_LINT_VERSION#v}" ]; then
    isPinnedVersion=true
  else
    isPinnedVersion=false
  fi

  if [ "$isPinnedVersion" = "false" ]; then
    echo "⚠ preflight-ci: golangci-lint version mismatch" >&2
    echo "  installed: $installedVersion" >&2
    echo "  expected:  ${GOLANGCI_LINT_VERSION#v}" >&2
    echo "  CI uses ${GOLANGCI_LINT_VERSION} — fix mismatch with:" >&2
    echo "    go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}" >&2
  fi
  return 0
}

run_lint() {
  echo "=== [2/2] golangci-lint (strict, full suite) ==="
  if ! ensure_golangci_lint; then
    return 1
  fi
  if ! golangci-lint run ./... \
        --timeout=5m \
        --max-issues-per-linter=0 \
        --max-same-issues=0; then
    echo "" >&2
    echo "✗ preflight-ci: golangci-lint failed" >&2
    return 1
  fi
  echo "  ✓ lint clean"
  return 0
}

testsOk=true
lintOk=true

case "$PHASE" in
  all)
    run_tests || testsOk=false
    echo ""
    run_lint  || lintOk=false
    ;;
  test|tests)
    run_tests || testsOk=false
    ;;
  lint)
    run_lint  || lintOk=false
    ;;
  *)
    echo "Usage: $0 [all|test|lint]" >&2
    exit 1
    ;;
esac

echo ""
echo "=== preflight-ci summary ==="
echo "  tests: $([ "$testsOk" = "true" ] && echo "PASS" || echo "FAIL")"
echo "  lint:  $([ "$lintOk"  = "true" ] && echo "PASS" || echo "FAIL")"

if [ "$testsOk" = "true" ] && [ "$lintOk" = "true" ]; then
  echo ""
  echo "✓ preflight-ci: ready to push"
  exit 0
fi

echo ""
echo "✗ preflight-ci: fix failures above before pushing" >&2
exit 1
