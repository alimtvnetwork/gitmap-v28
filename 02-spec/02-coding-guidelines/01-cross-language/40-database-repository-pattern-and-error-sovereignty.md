# Database Repository Pattern & Error Sovereignty — Coding Guidelines

## 1. Overview & Core Mission

Direct SQL execution scattered across user interfaces, CLI command handlers, and business logic coordinates creates architectural chaos:
1. SQL injection vulnerabilities and syntax discrepancies proliferate across files.
2. Raw SQL string changes require global codebase refactoring rather than isolated repository modifications.
3. Database errors are swallowed via anti-patterns like `_ = db.Close()` or `_ = conn.Exec(...)`.
4. Command layers call `os.Exit` or exit handlers directly inside database query helpers, breaking database audit log capture and terminating runtime processes unexpectedly.

This specification establishes the **Database Repository Pattern & Error Sovereignty Standard**. All database interactions MUST be encapsulated behind strongly-typed repository structures in the storage layer (`cli/store/` or equivalent). Leaf repository methods MUST return strongly-typed `*apperror.AppError` values directly. Swallowing errors (`_ = db.Close()`) is strictly banned.

---

## 2. The Split-DB Repository Architecture

Under the Split-DB pattern (such as GitMap's Tier 1 `gitmap.db`, Tier 2 `tasks.db`, and Tier 3 `pipeline.db`):
1. **Repository Abstraction:** Each database domain (e.g. `TasksSplitDB`, `PipelineSplitDB`, `HeatmapSplitDB`) provides a dedicated repository struct that encapsulates the underlying `*sql.DB` connection.
2. **Metadata Enums & Table Constants:** Table names, column names, and query constants MUST be defined as PascalCase constants or enums with `Type` suffix. Raw inline table string literals are forbidden in high-level business logic.
3. **Dedicated CRUD Methods:** Instead of executing raw `UPDATE` or `SELECT` queries in command handlers, handlers invoke clean repository methods (e.g. `db.UpdateTaskQueueStatus(queueId, status)`).

---

## 3. Anti-Patterns & Mandated Solutions

### 3.1 Raw SQL in High-Level Command Files

#### ❌ Anti-Pattern: Direct SQL in CLI Handlers
```go
// ❌ FORBIDDEN: Direct SQL string and unencapsulated query in CLI command file
func updateQueue(db *sql.DB, id, status string) {
	query := "UPDATE TaskQueue SET Status = ? WHERE QueueId = ?"
	_, _ = db.Exec(query, status, id) // raw SQL + swallowed error
}
```

#### ✅ MANDATORY: Storage Repository Encapsulation
```go
// ✅ MANDATORY: Repository method in cli/store/tasks_split_db.go
func (db *TasksSplitDB) UpdateTaskQueueStatus(queueId, status string) *apperror.AppError {
	if db == nil || db.Conn() == nil {
		return apperror.NewValidationError("tasks database connection is nil")
	}

	now := timeutil.NowUTCISO()
	query := "UPDATE TaskQueue SET Status = ?, UpdatedAt = ? WHERE QueueId = ?"
	res := ExecWrapper(db.Conn(), query, status, now, queueId)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "UpdateTaskQueueStatus")
	}

	return nil
}
```

### 3.2 Swallowed Connection Close Errors (`_ = db.Close()`)

#### ❌ Anti-Pattern: Swallowing Close Errors
```go
// ❌ FORBIDDEN: Blindly discarding database close errors
func closeDb(db *store.TasksSplitDB) {
	if db != nil {
		_ = db.Close() // Swallowed close error! Corrupted WAL locks go unrecorded!
	}
}
```

#### ✅ MANDATORY: Structured Audit Logging on Close Errors
```go
// ✅ MANDATORY: Capture and log close errors into internal audit database
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

### 3.3 Leaf Error Sovereignty vs. Outer Caller Exits

#### ❌ Anti-Pattern: Exit Handler in Leaf Database Method
```go
// ❌ FORBIDDEN: Leaf function terminating the process
func FindTask(db *sql.DB, id string) *Task {
	var task Task
	err := db.QueryRow("SELECT ...", id).Scan(&task.Id)
	if err != nil {
		cliexit.HandleError(err) // CRITICAL VIOLATION: Terminates process before DB audit capture!
		return nil
	}
	return &task
}
```

#### ✅ MANDATORY: Monadic Result Return Sovereignty with DRY `types.go` Reuse
```go
// Step 1: In pkg/task/types.go (in the package where it is defined)
type (
	// TaskWrap is the canonical single reusable result container for Task.
	TaskWrap = result.Wrap[Task]

	// TaskSliceWrap is the canonical single reusable collection container for Task.
	TaskSliceWrap = result.WrapSlice[Task]
)

// Step 2: In pkg/task/repository.go
// ✅ MANDATORY: Return strongly-typed TaskWrap directly to caller — zero leaked generic brackets
func FindTask(db *sql.DB, id string) TaskWrap {
	var task Task
	err := db.QueryRow("SELECT ...", id).Scan(&task.Id)
	if err != nil {
		return result.FailWrap[Task](apperror.WrapSimple(err, "FindTask"))
	}

	return result.OkWrap(task)
}

// ✅ MANDATORY: Return strongly-typed TaskSliceWrap for list queries
func ListTasks(db *sql.DB) TaskSliceWrap {
	rows, err := db.Query("SELECT ...")
	if err != nil {
		return result.FailWrapSlice[Task](apperror.WrapSimple(err, "ListTasks"))
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if scanErr := rows.Scan(&task.Id); scanErr != nil {
			return result.FailWrapSlice[Task](apperror.WrapSimple(scanErr, "ListTasks.Scan"))
		}
		tasks = append(tasks, task)
	}

	return result.OkWrapSlice(tasks)
}
```
*The outer command orchestrator receives the single Result object, inspects `res.HasError()`, passes structured errors to `finishCommandAudit`, and decides whether to display an error or exit.*

---

## 4. Database Metadata Enums Standard

All database file paths, table names, and schema constants MUST be declared in dedicated constants:

```go
package store

type TasksTableType string

const (
	TableParentTask TasksTableType = "ParentTask"
	TableSubtask    TasksTableType = "Subtask"
	TableTaskQueue  TasksTableType = "TaskQueue"
	TableTaskLog    TasksTableType = "TaskLog"
)

type DatabaseMetadataType string

const (
	DbMasterAgent DatabaseMetadataType = ".ai-memory/temp-agents/master-agent.db"
	DbTasksRoot   DatabaseMetadataType = ".ai-memory/temp-agents/tasks-root.db"
)
```

---

## 5. Summary Checklist

- [ ] Zero raw SQL queries in UI, CLI command, or coordinator files.
- [ ] All database queries reside in `cli/store/` within dedicated repository structs.
- [ ] Swallowed errors like `_ = db.Close()` or `_ = conn.Exec()` are 100% eliminated.
- [ ] Leaf database methods return `*apperror.AppError` directly; zero calls to `os.Exit` or `cliexit.HandleError` in storage functions.
- [ ] Database errors bubble up to root dispatchers so they are recorded by command audit loggers.
