# Subtask Plan 22: Unpushed Repository Listing (`cli/store/repo_query_filter.go`)

> **Subtask ID:** `subtask-22`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/repo_query_filter.go`  
> **Status:** Pending  

---

## 1. Objectives & Scope
1. Implement `ListReposUnpushed() ([]model.ScanRecord, error)` on `*DB` in `cli/store/repo_query_filter.go`.
2. Allow CLI commands (`cmdstatus`, `cmdpush`, `cmdcpar`) to query repositories flagged with unpushed commits in metadata notes or recorded in `HeatmapLocalTracking.IsDirty`.
3. Co-locate unpushed querying with `ListReposByTag` in `cli/store/repo_query_filter.go` to maintain `cli/store/repo.go` well under 500 lines.
4. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
Identifying unpushed repositories currently requires inspecting each repository directory via git subprocess calls or iterating over the entire repository slice in memory.

By providing a database-backed query:
1. Operations can immediately fetch candidates flagged with unpushed commits without invoking dozens of subprocess commands.
2. Integrates with the heatmap tracking layer (`HeatmapLocalTracking`) where dirty repository states are recorded during development sessions.
3. Keeps query logic encapsulated within the storage repository pattern.

---

## 3. Concrete Implementation Details

### 3.1 Function Implementation in `cli/store/repo_query_filter.go`
Add `ListReposUnpushed` following `ListReposByTag`:

```go
const sqlSelectUnpushedRepos = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE (Notes LIKE '%unpushed%' OR Notes LIKE '%ahead%') OR RepoName IN (SELECT RepoSlug FROM HeatmapLocalTracking WHERE IsDirty = 1) ORDER BY Slug ASC"

// ListReposUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
func (db *DB) ListReposUnpushed() ([]model.ScanRecord, error) {
	rows, err := QueryWrapper(db.conn, sqlSelectUnpushedRepos).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}

	defer rows.Close()

	return scanRows(rows)
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Errors from `QueryWrapper` or query execution bubble directly to the caller as `return nil, fmt.Errorf(constants.ErrDBQuery, err)`.
- **Rule 2 (Function Limit):** 9 lines of code, well within the 15-line limit (and meeting the <= 8-10 line preferred standard).
- **Rule 3 (Zero Nested Ifs):** Linear error handling with early return.
- **Rule 4 (Positive Booleans):** Subquery checks `IsDirty = 1` (positive polarity).
- **Rule 6 (Vertical Whitespace):** Proper blank lines around control blocks and before returns.

---

## 5. Verification Criteria
- [ ] Method `ListReposUnpushed() ([]model.ScanRecord, error)` is declared on `*DB` in `cli/store/repo_query_filter.go`.
- [ ] Executes `sqlSelectUnpushedRepos` filtering notes for `unpushed` / `ahead` or dirty heatmap status.
- [ ] Closes `rows` via `defer rows.Close()`.
- [ ] Maps rows using package helper `scanRows(rows)`.
- [ ] Function is 9 lines with zero nested ifs.
