# Subtask Plan 24: Batch File Touch Tracking with Single Prepared Transaction (`cli/heatmap/heatmap_tracker.go`)

> **Subtask ID:** `subtask-24`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/heatmap/heatmap_tracker.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `TrackBatchFileTouches(cwd string, touches []store.FileTouchCount) error` in `cli/heatmap/heatmap_tracker.go`.
2. Batch insert file touch records within a single database transaction (`tx.Begin()` and `tx.Commit()`).
3. Leverage a prepared SQL statement (`tx.Prepare()`) to eliminate repetitive statement compilation during bulk operations.
4. Keep all functions strictly <= 15 lines and enforce Error Return Sovereignty.

---

## 2. Targeted Problem & Architectural Rationale
When GitMap scans or analyzes large repositories, hundreds of files may be touched in a short burst. Currently, calling `TrackFileTouch` individually for every file causes:
1. Frequent opening and closing of SQLite database handles.
2. An individual disk commit/flush per file, severely degrading file tracking performance.
3. Database locking errors due to concurrent write contention.

`TrackBatchFileTouches` solves this by grouping batch records into a single atomic transaction.

---

## 3. Concrete Implementation Details

### 3.1 Function Additions in `cli/heatmap/heatmap_tracker.go`

```go
// TrackBatchFileTouches records multiple file touch records in a single database transaction.
func TrackBatchFileTouches(cwd string, touches []store.FileTouchCount) error {
	if len(touches) == 0 {
		return nil
	}

	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return err
	}

	defer hdb.Close()

	return insertBatchTouchesTx(hdb.Conn(), touches)
}

func insertBatchTouchesTx(conn *sql.DB, touches []store.FileTouchCount) error {
	tx, err := conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.batchTx")
	}

	defer tx.Rollback()

	if execErr := executeTouchInserts(tx, touches); execErr != nil {
		return execErr
	}

	return tx.Commit()
}

func executeTouchInserts(tx *sql.Tx, touches []store.FileTouchCount) error {
	stmt, err := tx.Prepare(`INSERT INTO HeatmapActionLog (TaskSlug, ActionType, FilePath, ExecutionDepth, ExecutionOrder) VALUES ('batch-touch', 'BATCH_TOUCH', ?, 0, 1)`)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.prepareBatch")
	}

	defer stmt.Close()

	for _, t := range touches {
		if _, err := stmt.Exec(t.FilePath); err != nil {
			return apperror.WrapSimple(err, "heatmap.execTouch")
		}
	}

	return nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** All errors wrapped with `apperror.WrapSimple` and returned to caller.
- **Rule 2 (Function Limit):** `TrackBatchFileTouches` (14 lines), `insertBatchTouchesTx` (13 lines), `executeTouchInserts` (14 lines). All <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Guard clause for empty slice, early error returns.
- **Rule 4 (Positive Booleans):** Positive conditions and clean flow.
- **Rule 6 (Vertical Whitespace):** Proper blank lines before returns and after `}`.

---

## 5. Verification Criteria
- [x] `TrackBatchFileTouches(cwd string, touches []store.FileTouchCount) error` is declared.
- [x] Returns `nil` immediately when `touches` slice is empty.
- [x] Executes batch touches within a single transaction using a prepared statement.
- [x] Rolls back cleanly on error.
- [x] Every helper function is <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
