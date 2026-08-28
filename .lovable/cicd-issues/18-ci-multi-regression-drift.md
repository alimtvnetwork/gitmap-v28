# CI/CD Fix: Multi-Regression Drift

## Error Summary
Multiple pipelines failed:
1. Generate drift (`gitmap/completion/allcommands_generated.go` out of sync).
2. Missing `setup-go-cached` action because `actions/checkout` was inside the composite instead of the caller.
3. Python syntax error in `query_wrapper.py` due to injected literal `` `n ``.
4. CHANGELOG drift for version 6.135.0.
5. Installer smoke test regex failing to parse `var Version = "6.135.0"` in `constants.go`.
6. GolangCI-Lint failing with Go 1.25.0 target vs Go 1.24 binary.

## Root Cause Analysis (RCA)
- The `commit-push` suite addition lacked a `go generate ./...` run.
- `setup-go-cached` failed because GitHub Actions requires the repo to be checked out *before* invoking local composite actions.
- `query_wrapper.py` was corrupted with PowerShell newline characters during a previous replacement.
- Version 6.135.0 was released without adding its section to `changelog.md`.
- `constants.go` `Version` was changed from `const` to `var` to support ldflags, breaking the installer script's naive regex.
- `golangci-lint v1.64.8` does not support the newly bumped `go 1.25.0` natively.

## Solution Applied
- Ran `go generate ./...` and committed `allcommands_generated.go`.
- Moved `actions/checkout@v6` into `ci.yml` callers (spell-check, lint) and removed it from `setup-go-cached/action.yml`.
- Fixed the Python syntax error in `query_wrapper.py`.
- Added the `## [v6.135.0]` heading to `changelog.md`.
- Updated `smoke-installer.ps1` regex to `^(const|var) Version`.
- Added `go: '1.24'` to `.golangci.yml` to explicitly instruct the older linter how to parse the AST.

## What NOT to Repeat
- NEVER invoke a local GitHub Action composite (`uses: ./.github/actions/...`) without running `actions/checkout` FIRST in the job.
- NEVER use PowerShell literal newlines (`` `n ``) as text replacement values in non-PowerShell files (like Python scripts).
- ALWAYS verify regex patterns in CI scripts if core constants change format (e.g., `const` to `var`).
