#!/usr/bin/env bash
# setup.sh — One-time project setup after cloning.
# Usage: ./setup.sh

set -euo pipefail

HOOK_SRC="hooks/pre-commit"
HOOK_DST=".git/hooks/pre-commit"

echo "Setting up gitmap development environment..."

# Install pre-commit hook
if [ -f "$HOOK_SRC" ]; then
  cp "$HOOK_SRC" "$HOOK_DST"
  chmod +x "$HOOK_DST"
  echo "  ✓ Pre-commit hook installed"
else
  echo "  ✗ Hook source not found: $HOOK_SRC"
fi

# Verify Go toolchain
if command -v go &>/dev/null; then
  echo "  ✓ Go $(go version | awk '{print $3}')"
else
  echo "  ⚠ Go not found — install from https://go.dev/dl/"
fi

# Install golangci-lint if missing
if command -v golangci-lint &>/dev/null; then
  echo "  ✓ golangci-lint available"
else
  echo "  → Installing golangci-lint..."
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 && echo "  ✓ golangci-lint installed" || echo "  ⚠ Failed to install golangci-lint"
fi

# Download Go dependencies
if [ -f "gitmap/go.mod" ]; then
  echo "  → Downloading Go dependencies..."
  (cd gitmap && go mod download) && echo "  ✓ Dependencies ready" || echo "  ⚠ go mod download failed"
fi

echo ""
echo "Done! Run 'cd gitmap && go test ./...' to verify."
