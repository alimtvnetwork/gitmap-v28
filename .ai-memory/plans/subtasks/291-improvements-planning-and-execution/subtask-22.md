# Subtask Plan 22: Repository Dirty Status Listing (`cli/store/repo.go`)

> **Subtask ID:** `subtask-22`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/repo.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `ListReposDirtyStatus() ([]model.ScanRecord, error)` on `*DB` in `cli/store/repo.go`.
2. Allow status, push, and commit commands to query repositories flagged dirty or needing sync directly from SQLite without scanning every clean repository disk path.
3. Adhere strictly to Error Return Sovereignty (Rule 1), Function Length Limits (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
Currently, commands such as `gitmap see`, `cpar`, and `pull` frequently iterate across all 40+ repositories in the workspace executing expensive `git status` subprocesses for clean repositories that have no pending changes.

By providing a database-level query `ListReposDirtyStatus()`, GitMap can:
1. Immediately prioritize repositories with uncommitted or flagged changes.
2. Reduce redundant OS filesystem and subprocess overhead for bulk sync commands.
3. Integrate seamlessly with the `HeatmapLocalTracking` table state.

---

## 3. Concrete Implementation Details

### 3.1 Method Addition in `cli/store/repo.go`
Add `ListReposDirtyStatus` after `ListReposByBranch`:

```go
const sqlSelectDirtyRepos = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Notes LIKE '%dirty%' OR RepoName IN (SELECT RepoSlug FROM HeatmapLocalTracking WHERE IsDirty = 1) ORDER BY Slug ASC"

// ListReposDirtyStatus returns repositories flagged as dirty or requiring synchronization.
func (db *DB) ListReposDirtyStatus() ([]model.ScanRecord, error) {
	rows, err := QueryWrapper(db.conn, sqlSelectDirtyRepos).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}

	defer rows.Close()

	return scanRows(rows)
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Database query errors bubble directly to the caller via `return nil, fmt.Errorf(...)`.
- **Rule 2 (Function Limit):** 11 lines of code, well within the <= 15 line rule.
- **Rule 3 (Zero Nested Ifs):** Linear execution with guard error return.
- **Rule 4 (Positive Booleans):** Positive naming conventions only.
- **Rule 6 (Vertical Whitespace):** Proper blank lines before returns and after `}`.

---

## 5. Verification Criteria
- [x] Method `ListReposDirtyStatus() ([]model.ScanRecord, error)` is declared on `*DB`.
- [x] Returns records flagged with dirty notes or registered with `IsDirty = 1` in `HeatmapLocalTracking`.
- [x] Function is <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
