# Subtask 01: Refactor RepoRepository to Return Monadic Result Containers

## Status: Completed
- **Parent Plan:** [295-cg-execute-wrapped-error-remediation.md](../../pending/295-cg-execute-wrapped-error-remediation.md)
- **Target File:** `cli/store/repo_repository.go`
- **Objective:** Eliminate two-value return tuples `([]model.ScanRecord, *apperror.AppError)` and `(string, *apperror.AppError)` in favor of `result.Result[[]model.ScanRecord]` and `result.Result[string]`.

---

## Action Items
1. Import `github.com/alimtvnetwork/gitmap-v28/cli/result`.
2. Update `ListByBranch`, `ListDirty`, `ListByTag`, `ListUnpushed` return types to `result.Result[[]model.ScanRecord]`.
3. Update `FindPathBySlug` return type to `result.Result[string]`.
4. Update `queryScanRecords` helper to return `result.Result[[]model.ScanRecord]`, utilizing `result.Ok` and `result.Fail`.
5. Maintain $\le 15$ lines per function and file cap $\le 300$ lines.
