# CI/CD Fix: Go 1.25 Upgrade & Strict Linter Failures

## Error Summary
The CI/CD pipeline failed with a massive wave of linter errors after bumping to `go1.25.0`. Initially, it threw `package requires newer Go version go1.25 (application built with go1.24)`. After upgrading `golangci-lint` to `v2.13.2`, the linter successfully parsed the Go 1.25 codebase but surfaced 40+ strict linting violations (misspell, nolintlint, unparam, wastedassign, errorlint, go vet unreachable code).

## Root Cause Analysis
1. **Go Version Mismatch:** `golangci-lint v1.64.8` was compiled against Go 1.24 and lacked the updated `x/tools` required to parse Go 1.25 modules, causing an outright parsing failure.
2. **Stricter Linter Rules in v2.x:** Upon upgrading `golangci-lint` to `v2.13.2`, newer, stricter analyzer versions flagged previously ignored code:
   - `misspell`: Caught British spellings (e.g., `behaviour`, `recognised`, `cancelled`).
   - `nolintlint`: Flagged `//nolint:gosec` directives that are no longer necessary.
   - `errorlint`: Flagged direct `err == ...` comparisons instead of `errors.Is`.
   - `govet`: Flagged unreachable code that was likely introduced during a recent refactor.

## Solution Applied
- Upgraded `golangci-lint` to `v2.13.2` in both `ci.yml` and the local environment.
- Ran `02-cicd-local-runner.py` to identify all strict violations natively.
- Fixed spelling errors (British to US English).
- Removed unused `//nolint:gosec` directives.
- Replaced direct error comparisons with `errors.Is()`.
- Removed unreachable code and wasted assignments.

## What NOT to Repeat
- NEVER bump `go.mod` to a new minor Go release (like 1.25) without also confirming that pinned linters (like golangci-lint) support that Go release. (Already in `strictly-avoid.md`).
- NEVER use British English spelling (e.g., `behaviour`, `recognise`) in the codebase, as the strict `misspell` linter enforces US English.
