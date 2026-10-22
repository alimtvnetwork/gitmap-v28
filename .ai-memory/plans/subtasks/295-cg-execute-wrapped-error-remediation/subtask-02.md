# Subtask 02: Refactor DB Delegators to Return Monadic Result Containers

## Status: Completed
- **Parent Plan:** [295-cg-execute-wrapped-error-remediation.md](../../pending/295-cg-execute-wrapped-error-remediation.md)
- **Target Files:**
  - `cli/store/repo_query_filter.go`
  - `cli/store/repo.go`
- **Objective:** Eliminate multi-value error tuples on `ListReposByTag`, `ListReposUnpushed`, `ListReposByBranch`, and `ListReposDirtyStatus`.

---

## Action Items
1. In `cli/store/repo_query_filter.go`:
   - Update `(db *DB) ListReposByTag(tag string) ScanRecordSliceWrap`.
   - Update `(db *DB) ListReposUnpushed() ScanRecordSliceWrap`.
2. In `cli/store/repo.go`:
   - Update `(db *DB) ListReposByBranch(branch string) ScanRecordSliceWrap`.
   - Update `(db *DB) ListReposDirtyStatus() ScanRecordSliceWrap`.
3. Preserve existing legacy methods to protect micro-batch boundaries.
