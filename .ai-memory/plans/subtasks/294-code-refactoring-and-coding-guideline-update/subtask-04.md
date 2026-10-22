# Subtask 294.4: Centralized Time Utility & CPAR Database Hygiene

- **Parent Plan:** [.ai-memory/plans/pending/294-code-refactoring-and-coding-guideline-update.md](../../pending/294-code-refactoring-and-coding-guideline-update.md)
- **Status:** Completed
- **Spec Reference:** [02-spec/21-app/294-code-refactoring-and-coding-guideline-update/02-component-spec.md](../../../02-spec/21-app/294-code-refactoring-and-coding-guideline-update/02-component-spec.md)
- **File Boundaries:**
  - `cli/timeutil/timeutil.go`
  - `cli/store/tasks_split_db.go`
  - `cli/cmdcpar/cpar.go`

---

## 1. Objectives & Scope

1. **Centralized Time Utility (`cli/timeutil/timeutil.go`):**
   - Create new `timeutil` package providing standardized RFC3339 UTC time generation helpers.
   - Eliminate scattered duplicate calls to `time.Now().UTC().Format(time.RFC3339)` across the codebase.
2. **Encapsulate Task Queue Updates in `cli/store/tasks_split_db.go`:**
   - Add `UpdateTaskQueueStatus(queueId, status string) error` method to `*TasksSplitDB`.
   - Prevent command layer files from issuing inline SQL commands directly against SQLite connections.
3. **CPAR DB Hygiene & Swallowed Error Elimination (`cli/cmdcpar/cpar.go`):**
   - Adopt `timeutil.NowUTCISO()` across CPAR task queue functions.
   - Refactor `updateCPARTaskQueue` to invoke the store helper method.
   - Eliminate Code Red violation in `closeCPARTaskDb` where `_ = db.Close()` swallowed potential database closure errors; route errors to `store.LogInternalErrorRecord`.
4. **Style & Rhythm Guidelines:**
   - Enforce positive boolean naming (`hasDb`, `hasErr`, `hasFailures`).
   - Keep all functions $\le 15$ lines.

---

## 2. Step-by-Step Execution Plan

### Step 1: Create `cli/timeutil/timeutil.go`
- Author `cli/timeutil/timeutil.go` with functions:
  - `NowUTCISO() string`: Returns `time.Now().UTC().Format(time.RFC3339)`.
  - `FormatUTCISO(t time.Time) string`: Formats any `time.Time` to RFC3339 UTC.
  - `NowTimestampCompact() string`: Compact string representation for identifiers.

### Step 2: Implement `UpdateTaskQueueStatus` in `cli/store/tasks_split_db.go`
- Add method to `*TasksSplitDB`:
  ```go
  func (db *TasksSplitDB) UpdateTaskQueueStatus(queueId, status string) error {
  	if db == nil || db.Conn() == nil {
  		return apperror.NewValidationError("tasks split db is nil")
  	}

  	now := time.Now().UTC().Format(time.RFC3339)
  	query := "UPDATE TaskQueue SET Status = ?, UpdatedAt = ? WHERE QueueId = ?"
  	res := ExecWrapper(db.Conn(), query, status, now, queueId)
  	if res.IsFailure {
  		return apperror.WrapSimple(res.Error, "UpdateTaskQueueStatus")
  	}

  	return nil
  }
  ```

### Step 3: Refactor `enqueueCPARTaskQueue` in `cli/cmdcpar/cpar.go`
- Replace inline time formatting with `timeutil.NowUTCISO()`:
  ```go
  now := timeutil.NowUTCISO()
  ```

### Step 4: Refactor `updateCPARTaskQueue` in `cli/cmdcpar/cpar.go`
- Replace inline SQL `store.ExecWrapper` with `db.UpdateTaskQueueStatus(queueId, status)`:
  ```go
  func updateCPARTaskQueue(db *store.TasksSplitDB, queueId, status string) {
  	hasDb := db != nil
  	if !hasDb {
  		return
  	}

  	err := db.UpdateTaskQueueStatus(queueId, status)
  	if err != nil {
  		rec := store.InternalErrorRecord{
  			ErrorCode:  "E_CPAR_QUEUE_UPDATE",
  			ErrorType:  "DB_ERROR",
  			Command:    "cpar",
  			Message:    fmt.Sprintf("failed to update task queue %s to %s: %v", queueId, status, err),
  			SourceFile: "cli/cmdcpar/cpar.go",
  		}
  		store.LogInternalErrorRecord(rec)
  	}
  }
  ```

### Step 5: Eliminate Code Red Swallowed Error in `closeCPARTaskDb`
- Replace `_ = db.Close()` with explicit error inspection and audit logging:
  ```go
  func closeCPARTaskDb(db *store.TasksSplitDB) {
  	hasDb := db != nil
  	if !hasDb {
  		return
  	}

  	closeErr := db.Close()
  	if closeErr != nil {
  		rec := store.InternalErrorRecord{
  			ErrorCode:  "E_CPAR_DB_CLOSE",
  			ErrorType:  "DB_CLOSE_ERROR",
  			Command:    "cpar",
  			Message:    fmt.Sprintf("failed to close tasks split db: %v", closeErr),
  			SourceFile: "cli/cmdcpar/cpar.go",
  		}
  		store.LogInternalErrorRecord(rec)
  	}
  }
  ```

---

## 3. Verification & Acceptance Criteria

- [x] `cli/timeutil/timeutil.go` exists and provides `NowUTCISO()`.
- [x] Task queue timestamps adopt `timeutil.NowUTCISO()` in `cli/cmdcpar/cpar.go`.
- [x] `closeCPARTaskDb` inspects `db.Close()` and logs failures to `store.LogInternalErrorRecord`.
- [x] All functions adhere strictly to the 15-line limit and positive boolean naming conventions.
