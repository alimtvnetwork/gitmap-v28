# Subtask Plan 22: Fast Repository Counter (`cli/store/repo.go`)

> **Subtask ID:** `subtask-22`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/repo.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `CountRepos() (int, error)` on `*DB` in `cli/store/repo.go`.
2. Provide an efficient, O(1) aggregate count query using SQLite `SELECT COUNT(*) FROM Repo`.
3. Eliminate the need for callers to load the entire repository dataset via `ListRepos()` simply to obtain the count.
4. Maintain strict adherence to Error Return Sovereignty (Rule 1) and Function Length Limits (<= 8 lines target, Rule 2).

---

## 2. Targeted Problem & Architectural Rationale
Multiple status, telemetry, and health check commands need to report the total count of tracked repositories in GitMap. Currently, these commands call `db.ListRepos()` and measure `len(records)`.

This causes:
1. Significant I/O overhead on large repository installations (hundreds or thousands of repositories), reading all columns from disk.
2. In-memory struct allocations for every repository row.
3. Wasteful CPU cycles scanning string fields that are immediately discarded.

---

## 3. Concrete Implementation Details

### 3.1 Method Addition in `cli/store/repo.go`
Add `CountRepos` to `*DB`:

```go
// CountRepos returns the total number of repositories tracked in the database.
func (db *DB) CountRepos() (int, error) {
	var total int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM Repo").Scan(&total)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrDBQuery, err)
	}

	return total, nil
}
```

### 3.2 Architectural Benefits
- Direct aggregate query executed by SQLite engine using internal B-Tree metadata or table row count.
- Zero slice or model struct allocations.
- Fully backwards compatible; existing `ListRepos` remains intact.

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Direct error propagation using `fmt.Errorf(constants.ErrDBQuery, err)`. Zero exit handler calls.
- **Rule 2 (Function Limit):** 8 lines total (matches the optimal <= 8 line preferred target).
- **Rule 3 (Zero Nested Ifs):** Flat single guard conditional.
- **Rule 6 (Vertical Whitespace):** Single blank line after `}` and before final `return total, nil`.

---

## 5. Verification Criteria
- [x] Method `CountRepos() (int, error)` exists on `*DB` in `cli/store/repo.go`.
- [x] Returns the exact count of rows in the `Repo` table.
- [x] Propagates database errors cleanly without panic or exit.
- [x] Function is exactly 8 lines with zero nested ifs.
- [x] All file paths are strictly relative lowercase Git paths.
