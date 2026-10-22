# Component Specification: Plan 294 — Code Refactoring, Coding Guideline Modernization & Subtask JSON Envelopes

## Status: APPROVED
- **Spec ID:** 294-02
- **Task Slug:** `294-code-refactoring-and-coding-guideline-update`
- **Architecture Spec:** [01-architecture-spec.md](01-architecture-spec.md)
- **Master Plan:** [.ai-memory/plans/pending/294-code-refactoring-and-coding-guideline-update.md](../../../.ai-memory/plans/pending/294-code-refactoring-and-coding-guideline-update.md)
- **Author:** Spec Author 02
- **Execution Mandate:** Zero-Build, Zero-Test, Strictly Relative Lowercase Git Paths

---

## 1. Executive Summary & Component Scope

This component specification details the surgical refactoring blueprints, data structures, and command enhancements across four core subsystems:

1. **Authentication Fix Subsystem (`cli/cmdfixauth/fixauth.go`):**
   - Enforce leaf error return sovereignty by updating `checkGitRepo` to return `*apperror.AppError`.
   - Decouple variable initialization from `if` statements (statement extraction).
   - Enforce multiline parameter formatting across all functions with $>2$ arguments.
2. **Centralized Time Utility (`cli/timeutil/timeutil.go`) & CPAR Subsystem (`cli/cmdcpar/cpar.go`):**
   - Architect standalone `cli/timeutil/` package providing standardized ISO 8601 UTC timestamp generation (`timeutil.NowUTCISO()`).
   - Eliminate direct inline SQL from `cpar.go` by routing `updateCPARTaskQueue` through a dedicated `store.TasksSplitDB` helper.
   - Eliminate swallowed error anti-pattern (`_ = db.Close()`) in `closeCPARTaskDb` via structured internal error audit logging.
3. **Agent Subtask Decomposition & File Size Cap Compliance (`cli/cmdagent/`):**
   - Decompose `cli/cmdagent/agent_subtask.go` and `cli/cmdagent/agent_subtask_resolve.go` into `agent_subtask_types.go` and `agent_subtask_cli.go`.
   - Maintain all subtask CLI files strictly $\le 300$ lines.
4. **Structured JSON Response Envelopes:**
   - Standardize `gitmap agent subtask ls --json` with a top-level metadata envelope (database path, parent task ID, status rollups, and subtasks list).
   - Standardize `gitmap agent subtask add` output to group commonalities (database path and parent ID printed once, added count, and subtasks list).

---

## 2. Component Blueprint: `cli/cmdfixauth/fixauth.go`

### 2.1 Leaf Error Sovereignty: `checkGitRepo`
- **Target File:** `cli/cmdfixauth/fixauth.go`
- **Current Signature:** `func checkGitRepo() error`
- **Refactored Signature:** `func checkGitRepo() *apperror.AppError`
- **Rationale:** Under repository error architecture guidelines, leaf verification functions returning domain errors must return strongly-typed `*apperror.AppError` rather than generic `error`.

```go
func checkGitRepo() *apperror.AppError {
	isGitRepo := cmdpull.IsGitRepoCWD()
	if !isGitRepo {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run inside the repo you want to fix)")

		return apperror.NewSimple("not a git repository (run inside the repo you want to fix)", "E_NOT_GIT_REPO")
	}

	return nil
}
```

Caller updates in `RunFixAuth`:
```go
func RunFixAuth(args []string) error {
	user, email, assumeYes, force := parseFixAuthFlags(args)

	userErr := checkUserFlag(user)
	if userErr != nil {
		return userErr
	}

	repoErr := checkGitRepo()
	if repoErr != nil {
		return repoErr
	}

	return executeFixAuthPipeline(user, email, assumeYes, force)
}
```

### 2.2 Control Flow Modernization: Extraction from `if` Statement
- **Location:** `executeFixAuthPipeline` in `cli/cmdfixauth/fixauth.go`
- **Violation:** Inlining call inside conditional: `if err := setupSSHDirAndKey(...); err != nil`
- **Guideline Rule:** Statements must be extracted before `if` conditions with vertical blank line rhythm.

```go
func executeFixAuthPipeline(
	user string,
	email string,
	assumeYes bool,
	force bool,
) error {
	keyPath := fixAuthKeyPath(user)

	err := setupSSHDirAndKey(keyPath, email, assumeYes, force)
	if err != nil {
		return err
	}

	return bindAndAnnounceAuth(user, keyPath)
}
```

### 2.3 Multiline Parameter Formatting for Functions with $>2$ Arguments
All function signatures with more than two parameters must be formatted across multiple lines with trailing commas on opening and closing parameters.

#### Functions to Format in `cli/cmdfixauth/fixauth.go`:

1. `executeFixAuthPipeline`:
```go
func executeFixAuthPipeline(
	user string,
	email string,
	assumeYes bool,
	force bool,
) error
```

2. `setupSSHDirAndKey`:
```go
func setupSSHDirAndKey(
	keyPath string,
	email string,
	assumeYes bool,
	force bool,
) error
```

3. `fixAuthGenerate`:
```go
func fixAuthGenerate(
	keyPath string,
	email string,
	assumeYes bool,
	force bool,
) error
```

4. `validateKeyGeneration`:
```go
func validateKeyGeneration(
	hasKey bool,
	keyPath string,
	assumeYes bool,
	force bool,
) bool
```

---

## 3. Component Blueprint: `cli/timeutil/` and `cli/cmdcpar/`

### 3.1 New Package: `cli/timeutil/timeutil.go`
- **Path:** `cli/timeutil/timeutil.go`
- **Package Name:** `timeutil`
- **Purpose:** Centralize ISO 8601 UTC timestamp formatting, preventing inconsistent RFC3339 formats and duplicative standard library calls across packages.

```go
// Package timeutil provides centralized RFC3339 UTC time formatting helpers.
package timeutil

import (
	"time"
)

// NowUTCISO returns the current time formatted as UTC RFC3339.
func NowUTCISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// FormatUTCISO formats the specified time as UTC RFC3339.
func FormatUTCISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// NowTimestampCompact returns a compact UTC timestamp (HHMMSS).
func NowTimestampCompact() string {
	return time.Now().UTC().Format("150405")
}
```

### 3.2 Split-DB Store Helper: `cli/store/tasks_split_db.go`
- **Method Addition:** `(db *TasksSplitDB) UpdateTaskQueueStatus(queueId, status string) error`
- **Rationale:** Repositories must encapsulate all SQL operations in `cli/store/`. CLI command files must never execute raw `UPDATE TaskQueue` SQL queries directly.

```go
// UpdateTaskQueueStatus updates the status and UpdatedAt timestamp for an enqueued task.
func (db *TasksSplitDB) UpdateTaskQueueStatus(queueId, status string) error {
	if db == nil || db.Conn() == nil {
		return apperror.NewValidationError("tasks database connection is nil")
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

### 3.3 Refactoring `cli/cmdcpar/cpar.go`
1. **Centralize Timestamps:** Replace inline `time.Now().UTC().Format(time.RFC3339)` with `timeutil.NowUTCISO()`.
2. **Encapsulate Task Queue Updates:** Update `updateCPARTaskQueue` to invoke `db.UpdateTaskQueueStatus`.
3. **Eliminate Swallowed Close Error (`_ = db.Close()`):** Log errors with `store.LogInternalErrorRecord` upon connection close failures.

```go
func enqueueCPARTaskQueue(action, target string) (string, *store.TasksSplitDB) {
	queueId := fmt.Sprintf("%s-%d", action, time.Now().UnixNano())
	tasksDb, err := store.OpenTasksRootSplitDB()
	hasErr := err != nil
	if hasErr {
		return "", nil
	}

	now := timeutil.NowUTCISO()
	store.ExecWrapper(
		tasksDb.Conn(),
		"INSERT INTO TaskQueue (QueueId, Section, Action, Target, Status, CreatedAt, UpdatedAt) VALUES (?, 'cpar', ?, ?, 'pending', ?, ?)",
		queueId,
		action,
		target,
		now,
		now,
	)

	return queueId, tasksDb
}

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

## 4. Component Blueprint: `cli/cmdagent/` Decomposition

### 4.1 Decomposition Strategy
Both `cli/cmdagent/agent_subtask.go` (297 lines) and `cli/cmdagent/agent_subtask_resolve.go` (309 lines) are dangerously close to or exceeding the 300-line repository ceiling.

The code is restructured into four discrete files:

```
cli/cmdagent/
├── agent_subtask_types.go    (~75 lines)  - Options, response envelopes, summary models
├── agent_subtask_cli.go      (~150 lines) - Cobra command declarations, flag wiring, table formatting
├── agent_subtask.go          (~190 lines) - Command runners (Add, Claim, Start, Complete, Fail, Ls)
└── agent_subtask_resolve.go  (~220 lines) - Resolution, slug parsing, JSON input normalizers
```

### 4.2 Type Definitions: `cli/cmdagent/agent_subtask_types.go`

```go
package cmdagent

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

type SubtaskAddOptions struct {
	ParentId   string
	JsonInput  string
	Title      string
	TaskCode   string
	Slug       string
	OwnedFiles string
	IsJson     bool
}

type SubtaskClaimOptions struct {
	AgentRole string
	TaskId    string
}

type SubtaskStartOptions struct {
	SubtaskId string
	TaskCode  string
	Slug      string
	AgentRole string
	TaskId    string
}

type SubtaskCompleteOptions struct {
	SubtaskId string
	TaskCode  string
	Slug      string
	AgentRole string
	Evidence  string
	TaskId    string
}

type SubtaskFailOptions struct {
	SubtaskId string
	TaskCode  string
	Slug      string
	AgentRole string
	Reason    string
	TaskId    string
}

type SubtaskLsOptions struct {
	TaskId string
	IsJson bool
}

// SubtaskStatusSummary contains aggregate counts for subtasks in a task database.
type SubtaskStatusSummary struct {
	TotalCount      int `json:"totalCount"`
	PendingCount    int `json:"pendingCount"`
	InProgressCount int `json:"inProgressCount"`
	CompletedCount  int `json:"completedCount"`
	FailedCount     int `json:"failedCount"`
}

// SubtaskLsEnvelope is the standard JSON response envelope for `gitmap agent subtask ls --json`.
type SubtaskLsEnvelope struct {
	Status       string               `json:"status"`
	DatabasePath string               `json:"databasePath"`
	ParentTaskId string               `json:"parentTaskId"`
	Summary      SubtaskStatusSummary `json:"summary"`
	Subtasks     []types.Subtask      `json:"subtasks"`
}

// SubtaskAddEnvelope is the standard JSON response envelope for `gitmap agent subtask add`.
type SubtaskAddEnvelope struct {
	Status       string          `json:"status"`
	DatabasePath string          `json:"databasePath"`
	ParentTaskId string          `json:"parentTaskId"`
	AddedCount   int             `json:"addedCount"`
	Subtasks     []types.Subtask `json:"subtasks"`
}
```

### 4.3 CLI Wiring: `cli/cmdagent/agent_subtask_cli.go`
Responsible for:
- Cobra command definitions (`subtaskAddCmd`, `subtaskClaimCmd`, `subtaskStartCmd`, `subtaskCompleteCmd`, `subtaskFailCmd`, `subtaskLsCmd`).
- Flag registration (`initSubtaskCommands`, `registerTargetFlags`).
- Terminal table renderer (`renderSubtasksTable`).

---

## 5. JSON Response Envelope Specifications

### 5.1 `gitmap agent subtask ls --json` Response Schema

```json
{
  "status": "success",
  "databasePath": ".gitmap/data/tasks/task-294.db",
  "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
  "summary": {
    "totalCount": 5,
    "pendingCount": 1,
    "inProgressCount": 2,
    "completedCount": 2,
    "failedCount": 0
  },
  "subtasks": [
    {
      "subtaskId": "sub-120000-subtask-01",
      "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
      "taskCode": "Subtask-01",
      "taskSlug": "294-code-refactoring-and-coding-guideline-update-subtask-01",
      "title": "Architecture Specification & Subtask Ledger",
      "status": "COMPLETED",
      "assignedAgentRole": "Spec Author 01",
      "ownedFilesJson": "[\"02-spec/21-app/294-code-refactoring-and-coding-guideline-update/01-architecture-spec.md\"]",
      "files": [
        {
          "path": "02-spec/21-app/294-code-refactoring-and-coding-guideline-update/01-architecture-spec.md",
          "exists": true,
          "fileType": "spec"
        }
      ]
    }
  ]
}
```

Implementation in `renderSubtasksJsonEnvelope`:
```go
func renderSubtasksJsonEnvelope(subtasks []types.Subtask, dbPath, parentId string) *appfault.AppError {
	cwd, _ := os.Getwd()
	relDbPath := toRelativePath(cwd, dbPath)

	summary := computeSubtaskSummary(subtasks)
	for i := range subtasks {
		subtasks[i].Files = fsutil.ClassifyFileTargets(subtasks[i].OwnedFilesJson, "", cwd)
	}

	envelope := SubtaskLsEnvelope{
		Status:       "success",
		DatabasePath: relDbPath,
		ParentTaskId: parentId,
		Summary:      summary,
		Subtasks:     subtasks,
	}

	b, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return appfault.WrapSimple(err, "renderSubtasksJsonEnvelope")
	}

	fmt.Println(string(b))

	return nil
}

func computeSubtaskSummary(subtasks []types.Subtask) SubtaskStatusSummary {
	summary := SubtaskStatusSummary{TotalCount: len(subtasks)}
	for _, s := range subtasks {
		switch strings.ToUpper(s.Status) {
		case "PENDING":
			summary.PendingCount++
		case "IN_PROGRESS", "RUNNING":
			summary.InProgressCount++
		case "COMPLETED", "DONE":
			summary.CompletedCount++
		case "FAILED":
			summary.FailedCount++
		}
	}

	return summary
}
```

### 5.2 `gitmap agent subtask add` Response Schema
Eliminates duplicated metadata across individual subtask items by printing the database path and parent task ID once at the top level.

```json
{
  "status": "success",
  "databasePath": ".gitmap/data/tasks/task-294.db",
  "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
  "addedCount": 3,
  "subtasks": [
    {
      "subtaskId": "sub-120000-subtask-03",
      "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
      "taskCode": "Subtask-03",
      "taskSlug": "294-code-refactoring-and-coding-guideline-update-subtask-03",
      "title": "Refactor cmdfixauth Error Handling and Formatting",
      "status": "PENDING"
    }
  ]
}
```

---

## 6. Coding Guidelines & Invariant Adherence Matrix

| Guideline Rule | Component | Enforcement Strategy |
| :--- | :--- | :--- |
| **Error Sovereignty** | `cli/cmdfixauth/fixauth.go` | `checkGitRepo` returns `*apperror.AppError` directly; callers propagate error upwards. |
| **No Swallowed Errors** | `cli/cmdcpar/cpar.go` | `closeCPARTaskDb` handles `db.Close()` error and writes to `store.InternalErrorRecord`. |
| **No Raw SQL in Commands** | `cli/cmdcpar/cpar.go` | Task queue SQL updates encapsulated within `TasksSplitDB.UpdateTaskQueueStatus`. |
| **Clean Control Flow** | `cli/cmdfixauth/fixauth.go` | Variables assigned before `if` statements; no inline assignments inside `if`. |
| **Function Line Cap** | All files | All functions partitioned $\le 15$ lines ($\le 8$ lines preferred). |
| **File Line Cap** | `cli/cmdagent/` | Files partitioned into 4 discrete modules, ensuring each file is $\le 300$ lines. |
| **Multiline Parameters** | `cli/cmdfixauth/fixauth.go` | All functions with $>2$ arguments formatted on separate lines with trailing commas. |
| **Path Hygiene** | All components | All file paths relative and strictly lowercase; no drive prefixes. |
