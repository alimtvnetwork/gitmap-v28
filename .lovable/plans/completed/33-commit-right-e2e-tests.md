# Fix Commit-Right Path Resolution and Add E2E Tests

## Root Cause

The `commit-right` command only processed 8 commits because the user provided the path `..\prompt-architect-v2\` from `D:\work\commit-fix`, which resolved to `D:\work\prompt-architect-v2` (a directory that explicitly only contains those 8 commits). The actual 238-commit repository was located at `.\prompt-architect-v2\` (which was cloned subsequently). 

However, per user instructions, we must bolster the system's end-to-end testing and ensure `apperror` management is perfectly applied across `committransfer`.

## Architectural Plan

1. **Audit `committransfer` Error Management**: Review `gitmap/cmd/committransfer.go`, `gitmap/cmd/dispatchcommittransfer.go`, and `gitmap/committransfer/` to ensure all errors are properly wrapped using `apperror.Wrap` or `*apperror.AppError`, avoiding swallowed errors or bare `os.Exit(1)`.
2. **End-to-End Tests**: Create comprehensive E2E tests for `commit-right` in `gitmap/tests/committransfer_test/` (or similar) that simulate transferring commits from a local source to a local target, verifying exact commit counts and metadata integrity.

## Code Review Guides

- Follow `spec/02-coding-guidelines/`.
- Ensure all booleans start with `is`, `has`, `can`, or `should`.
- Do not use generic variables like `data`, `res`, `temp`.
- All temporary scripts go to `.lovable/temp-scripts/`.

## Subtasks

- **01-audit-errors**: Scan `committransfer` for bare `fmt.Errorf` or missing `apperror.Wrap` and fix them.
- **02-write-e2e**: Write the E2E test suite for `commit-right`.
