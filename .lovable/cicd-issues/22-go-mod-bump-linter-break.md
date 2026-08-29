# 22-go-mod-bump-linter-break

## Error Summary
The CI/CD pipeline failed during the `full-suite-lint.sh` script with multiple typecheck errors:
`package requires newer Go version go1.25 (application built with go1.24)`

## Root Cause Analysis
During a previous task, `gitmap/go.mod` was updated from `go 1.24.13` to `go 1.25.0`. However, the CI pipeline installs and pins `golangci-lint v1.64.8` via GitHub Actions (`ci.yml`), which was built with and only fully supports Go 1.24. When the pinned linter attempted to parse the Go 1.25 packages, its internal `go/packages` parser crashed, causing a global linter failure.

## Solution Applied
Reverted `gitmap/go.mod` back to `go 1.24.13` and ran `go mod tidy`. Used `.lovable/ai-fix-scripts/02-cicd-local-runner.py` to run `golangci-lint` locally and verified that the lint and vet steps now pass with 0 errors.

## What NOT to Repeat
- NEVER bump `go.mod` to a new minor Go release (like 1.25) without also confirming that pinned linters (like `golangci-lint`) support that Go release and updating `.github/workflows/ci.yml`.
