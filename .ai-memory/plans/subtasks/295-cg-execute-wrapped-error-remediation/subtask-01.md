# Subtask 01: Refactor RepoRepository to Return Monadic Result Containers

## Status: Completed
- **Parent Plan:** [295-cg-execute-wrapped-error-remediation.md](../../pending/295-cg-execute-wrapped-error-remediation.md)
- **Target File:** `cli/store/repo_repository.go`
- **Objective:** Eliminate two-value return tuples in favor of canonical `ScanRecordSliceWrap` (`result.WrapSlice[model.ScanRecord]`), `ScanRecordWrap` (`result.Wrap[model.ScanRecord]`), and `RepoPathWrap` (`result.Wrap[string]`) declared in `cli/store/repository_types.go`.

---

## Action Items
1. In `cli/store/repository_types.go`, declare `ScanRecordWrap`, `ScanRecordSliceWrap`, `ScanRecordMapWrap`, and `RepoPathWrap`.
2. Update `ListByBranch`, `ListDirty`, `ListByTag`, `ListUnpushed`, `ListAll` return types to `ScanRecordSliceWrap`.
3. Add `FindBySlug(slug string)` returning `ScanRecordWrap`.
4. Update `FindPathBySlug` return type to `RepoPathWrap`.
5. Update `queryScanRecords` helper to return `ScanRecordSliceWrap`, utilizing `result.OkSlice` and `result.FailSlice` (eliminating redundant `FailWrapSlice` double-wrap calls).
6. Maintain $\le 15$ lines per function and file cap $\le 300$ lines.
