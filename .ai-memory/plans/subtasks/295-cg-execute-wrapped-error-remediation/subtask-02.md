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
   - Import `github.com/alimtvnetwork/gitmap-v28/cli/result`.
   - Update `(db *DB) ListReposByTag(tag string) result.Result[[]model.ScanRecord]`.
   - Update `(db *DB) ListReposUnpushed() result.Result[[]model.ScanRecord]`.
2. In `cli/store/repo.go`:
   - Import `github.com/alimtvnetwork/gitmap-v28/cli/result`.
   - Update `(db *DB) ListReposByBranch(branch string) result.Result[[]model.ScanRecord]`.
   - Update `(db *DB) ListReposDirtyStatus() result.Result[[]model.ScanRecord]`.
3. Preserve existing legacy methods to protect micro-batch boundaries.
