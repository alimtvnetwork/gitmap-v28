# Subtask 02: Heatmap Tracker Error Sovereignty Refactoring

> **Subtask ID:** `SUBTASK-296-02`  
> **Parent Task:** `296-result-wrapper-and-apperror-returns`  
> **Status:** Completed  
> **Target Files:**  
> - `cli/heatmap/heatmap_tracker.go` (modify)  
> - `cli/heatmap/heatmap_tracker_agent.go` (modify)  
> **Related Architecture Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/01-architecture-spec.md`  

---

## 1. Objective

Refactor telemetry and auditing functions in `cli/heatmap/heatmap_tracker.go` and `cli/heatmap/heatmap_tracker_agent.go` to return strongly-typed `*apperror.AppError` instead of standard library `error`. Enforce error return sovereignty, eliminate bare error interfaces, and ensure no double-wrapping occurs.

---

## 2. Action Checklist

- [ ] 1. Refactor `cli/heatmap/heatmap_tracker.go`:
  - [ ] Convert `TrackFileTouch(filePath string, touchType store.FileTouchType) *apperror.AppError`:
    - Ensure clean path validation returns `nil` for empty or root paths.
    - Wrap `os.Getwd()` failure with `apperror.WrapSimple(err, "heatmap.TrackFileTouch.cwd")`.
    - Return `recordFileTouchAtRoot(cwd, cleanPath, touchType)`.
  - [ ] Convert `recordFileTouchAtRoot(cwd, cleanPath string, touchType store.FileTouchType) *apperror.AppError`:
    - Wrap `store.OpenHeatmapRootDB(cwd)` failure with `apperror.WrapSimple(err, "heatmap.recordFileTouchAtRoot.openDB")`.
    - Wrap `hdb.LogAction` failure with `apperror.WrapSimple(logErr, "heatmap.recordFileTouchAtRoot.logAction")`.
    - Return `nil` on success.
  - [ ] Convert `TrackBatchFileTouches(cwd string, touches []store.FileTouchCount) *apperror.AppError`:
    - Early return `nil` if `len(touches) == 0`.
    - Wrap `store.OpenHeatmapRootDB(cwd)` error.
    - Return `insertBatchTouchesTx(hdb.Conn(), touches)`.
  - [ ] Convert `insertBatchTouchesTx(conn *sql.DB, touches []store.FileTouchCount) *apperror.AppError`:
    - If `conn == nil`, return `apperror.NewSimple("database connection is nil", "E_NIL_CONN")`.
    - Wrap transaction begin, prepare, statement execution, and commit errors using `apperror.WrapSimple`.
    - Return `nil` on successful transaction commit.

- [ ] 2. Refactor `cli/heatmap/heatmap_tracker_agent.go`:
  - [ ] Convert `TrackAgentAction(taskSlug, actionType, summary string) *apperror.AppError`:
    - Early return `nil` if `taskSlug` or `actionType` is empty.
    - Wrap `os.Getwd()` failure with `apperror.WrapSimple(err, "heatmap.TrackAgentAction.cwd")`.
    - Return `recordAgentActionAtRoot(cwd, cleanSlug, cleanAction, summary)`.
  - [ ] Convert `recordAgentActionAtRoot(cwd, taskSlug, actionType, summary string) *apperror.AppError`:
    - Wrap `store.OpenHeatmapRootDB(cwd)` error.
    - Wrap `hdb.LogAction` error with `apperror.WrapSimple(logErr, "heatmap.recordAgentActionAtRoot.logAction")`.
    - Return `nil` on success.

- [ ] 3. Anti-Pattern & Invariant Checks:
  - [ ] Zero double-wrapping: Never call `apperror.WrapSimple` on an error that is already an `*apperror.AppError`.
  - [ ] Zero bare `error` return signatures in modified functions.
  - [ ] Functions stay within 15-line limit with zero nested if statements.
  - [ ] Strictly relative Git paths only.
  - [ ] Zero compiler builds and zero tests (`go build` / `go test` strictly banned).
  - [ ] Zero git commands.

---

## 3. Concrete Code Specifications

### Refactoring in `cli/heatmap/heatmap_tracker.go`

```go
// TrackFileTouch records a file access operation into the root heatmap database.
func TrackFileTouch(filePath string, touchType store.FileTouchType) *apperror.AppError {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimSpace(filePath)))
	isValidPath := len(cleanPath) > 0 && cleanPath != "."
	if !isValidPath {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.TrackFileTouch.cwd")
	}

	return recordFileTouchAtRoot(cwd, cleanPath, touchType)
}

func recordFileTouchAtRoot(cwd, cleanPath string, touchType store.FileTouchType) *apperror.AppError {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.recordFileTouchAtRoot.openDB")
	}

	defer hdb.Close()

	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
		ActionType: string(touchType),
		FilePath:   cleanPath,
	})
	if logErr != nil {
		return apperror.WrapSimple(logErr, "heatmap.recordFileTouchAtRoot.logAction")
	}

	return nil
}

// TrackBatchFileTouches records multiple file touch records in a single database transaction.
func TrackBatchFileTouches(cwd string, touches []store.FileTouchCount) *apperror.AppError {
	if len(touches) == 0 {
		return nil
	}

	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.TrackBatchFileTouches.openDB")
	}

	defer hdb.Close()

	return insertBatchTouchesTx(hdb.Conn(), touches)
}

func insertBatchTouchesTx(conn *sql.DB, touches []store.FileTouchCount) *apperror.AppError {
	if conn == nil {
		return apperror.NewSimple("database connection is nil", "E_NIL_CONN")
	}

	tx, err := conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.insertBatchTouchesTx.begin")
	}

	stmt, err := tx.Prepare(sqlInsertFileTouch)
	if err != nil {
		_ = tx.Rollback()

		return apperror.WrapSimple(err, "heatmap.insertBatchTouchesTx.prepare")
	}
	defer stmt.Close()

	for _, touch := range touches {
		if _, execErr := stmt.Exec(touch.FilePath, touch.Count, touch.ReadCount, touch.EditCount, touch.CreatedCount); execErr != nil {
			_ = tx.Rollback()

			return apperror.WrapSimple(execErr, "heatmap.insertBatchTouchesTx.exec")
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return apperror.WrapSimple(commitErr, "heatmap.insertBatchTouchesTx.commit")
	}

	return nil
}
```

### Refactoring in `cli/heatmap/heatmap_tracker_agent.go`

```go
// TrackAgentAction records a discrete agent lifecycle event into the root heatmap database.
func TrackAgentAction(taskSlug, actionType, summary string) *apperror.AppError {
	cleanSlug := strings.TrimSpace(taskSlug)
	cleanAction := strings.TrimSpace(actionType)
	hasMissingInput := cleanSlug == "" || cleanAction == ""
	if hasMissingInput {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.TrackAgentAction.cwd")
	}

	return recordAgentActionAtRoot(cwd, cleanSlug, cleanAction, summary)
}

func recordAgentActionAtRoot(cwd, taskSlug, actionType, summary string) *apperror.AppError {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.recordAgentActionAtRoot.openDB")
	}

	defer hdb.Close()

	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
		TaskSlug:       taskSlug,
		ActionType:     actionType,
		ReasoningNotes: summary,
		ReasoningState: "EXPLORATION",
	})
	if logErr != nil {
		return apperror.WrapSimple(logErr, "heatmap.recordAgentActionAtRoot.logAction")
	}

	return nil
}
```

---

## 4. Verification & Done Criteria

- [ ] `TrackFileTouch` and `TrackBatchFileTouches` signatures in `cli/heatmap/heatmap_tracker.go` return `*apperror.AppError`.
- [ ] `TrackAgentAction` signature in `cli/heatmap/heatmap_tracker_agent.go` returns `*apperror.AppError`.
- [ ] All internal helpers (`recordFileTouchAtRoot`, `insertBatchTouchesTx`, `recordAgentActionAtRoot`) return `*apperror.AppError`.
- [ ] No double-wrapping or unhandled errors.
- [ ] No builds or tests executed.
