# Issue: Go 1.25 gofmt and goimports Drift Across Repository Packages

- date: 2026-08-29
- status: resolved
- scope: style/ci

## 1. Why it happened

Upgrading the project to Go 1.25 introduced new comment and docstring formatting rules in `gofmt`. Furthermore, a previous set of changes modified import statements and comments without running recursive formatting and import normalization across all subdirectories.

## 2. How it happened

When CI ran `bash ../.github/scripts/go-format-check.sh`, the check script executed `gofmt -l .` and detected 36 files whose doc comments, inline comments, or blank line spacing did not match Go 1.25 `gofmt` output, causing the build gate to fail with exit code 1.

## 3. Root Cause

- Discrepancy between older Go 1.24 formatting and Go 1.25 `gofmt` expectations for doc comments and directives (e.g. blank comment lines before `//nolint:unused` and comment alignment).
- Subdirectory files had unorganized import blocks not grouped according to local-prefix rules (`github.com/alimtvnetwork/gitmap-v28/gitmap`).

## 4. Code Fix

- Ran recursive `gofmt -w` across all `.go` files in the repository.
- Ran recursive `goimports -w -local "github.com/alimtvnetwork/gitmap-v28/gitmap"` across all `.go` files in the repository.
- Verified that both `gofmt -l .` and `goimports -l -local ...` return 0 unformatted files.
- Verified that `.lovable/ai-fix-scripts/02-cicd-local-runner.py` passes all checks.
