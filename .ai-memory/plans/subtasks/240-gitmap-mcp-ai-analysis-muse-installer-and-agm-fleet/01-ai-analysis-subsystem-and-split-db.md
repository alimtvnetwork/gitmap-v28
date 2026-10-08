# Subtask 01: AI Analysis Subsystem & Split-DB Schema

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/240-gitmap-mcp-ai-analysis-muse-installer-and-agm-fleet.md`  
> **Target Subsystems:** `cli/store/ai_analysis_split_db.go`, `cli/cmdai/`, `cli/store/split_db_path.go`  
> **Owned Files:**  
> - `cli/store/ai_analysis_split_db.go`  
> - `cli/store/split_db_path.go`  
> - `cli/cmdai/ai_analysis_types.go`  
> - `cli/cmdai/ai_analysis_db.go`  
> - `cli/cmdai/ai_analysis_cmd.go`  
> - `cli/cmdai/ai_analysis_test.go`  

---

## 1. Concrete Objectives

1. **Split-DB Path & Partition Resolution:**
   - In `cli/store/split_db_path.go`, declare constant `SectionAiAnalysis = "ai-analysis"`.
   - Implement `ResolveAiAnalysisSplitDbPath(slug, repoRoot string) string` resolving to `.gitmap/data/ai-analysis/<slug>/sql.db`.
   - Ensure target directories are created with `0755` permissions via `os.MkdirAll`.

2. **Schema Definition & Migration Engine:**
   - In `cli/store/ai_analysis_split_db.go`, implement schema initialization with WAL mode (`PRAGMA journal_mode=WAL`), foreign key enforcement (`PRAGMA foreign_keys=ON`), and busy timeout (`PRAGMA busy_timeout=5000`).
   - Create master table `AiAnalysisTask` tracking task metadata, executing model, status, timestamps, file counters, and reasoning summary.
   - Create child table `AiAnalysisLine` with foreign key cascade to `AiAnalysisTask`, capturing file path, operation type (`read`, `edit`, `delete`), line range, content snippet, and explicit reasoning.
   - Create optimal composite indexes on `(status, created_at)` and `(task_id, file_path)`.

3. **Database Access Layer (DAL):**
   - In `cli/cmdai/ai_analysis_db.go`, implement:
     - `OpenAiAnalysisSplitDB(repoRoot string) (*sql.DB, error)`
     - `CreateAiAnalysisTask(db *sql.DB, task AiAnalysisTask) error`
     - `RecordAiAnalysisLine(db *sql.DB, line AiAnalysisLine) error`
     - `BatchRecordAiAnalysisLines(db *sql.DB, lines []AiAnalysisLine) error`
     - `GetAiAnalysisTask(db *sql.DB, taskID string) (*AiAnalysisTask, error)`
     - `ListAiAnalysisTasks(db *sql.DB, limit, offset int) ([]AiAnalysisTask, error)`
     - `GetAiAnalysisLinesForTask(db *sql.DB, taskID string) ([]AiAnalysisLine, error)`
     - `UpdateAiAnalysisTaskStatus(db *sql.DB, taskID string, status AiTaskStatusType, summary string) error`

4. **CLI Command Tree Integration:**
   - In `cli/cmdai/ai_analysis_cmd.go`, register `analysisCmd` under root `AiCmd` (`gitmap ai analysis`).
   - Provide subcommands:
     - `gitmap ai analysis start --task <id> --desc <description> [--model <name>]`
     - `gitmap ai analysis record --task <id> --file <path> --lines <start-end> --op <read|edit|delete> --reason <text>`
     - `gitmap ai analysis finish --task <id> [--status <completed|failed>] [--summary <text>]`
     - `gitmap ai analysis list [--limit <n>]`
     - `gitmap ai analysis inspect <task_id>`
   - Implement ANSI table rendering for list and inspect outputs.

5. **Error Management & Validation:**
   - Enforce structured errors with `apperror.NewWithDetails("cmd.ai.analysis", "E3011", ...)` for missing tasks, invalid line ranges, or malformed operations.
   - Validate that referenced file paths exist on disk or normalize them relative to repository root.

---

## 2. Core Domain Types & Structs (`cli/cmdai/ai_analysis_types.go`)

```go
package cmdai

import "time"

// AiTaskStatusType defines valid task statuses.
type AiTaskStatusType string

const (
	AiTaskStatusInProgress AiTaskStatusType = "in_progress"
	AiTaskStatusCompleted  AiTaskStatusType = "completed"
	AiTaskStatusFailed     AiTaskStatusType = "failed"
	AiTaskStatusCancelled  AiTaskStatusType = "cancelled"
)

// AiOperationType defines the file interaction performed.
type AiOperationType string

const (
	AiOperationRead   AiOperationType = "read"
	AiOperationEdit   AiOperationType = "edit"
	AiOperationDelete AiOperationType = "delete"
)

// AiAnalysisTask represents a high-level AI analysis session.
type AiAnalysisTask struct {
	TaskId             string           `json:"taskId" db:"task_id"`
	RepoPath           string           `json:"repoPath" db:"repo_path"`
	TaskDescription    string           `json:"taskDescription" db:"task_description"`
	ModelName          string           `json:"modelName" db:"model_name"`
	Status             AiTaskStatusType `json:"status" db:"status"`
	StartedAt          time.Time        `json:"startedAt" db:"started_at"`
	CompletedAt        *time.Time       `json:"completedAt,omitempty" db:"completed_at"`
	TotalFilesRead     int              `json:"totalFilesRead" db:"total_files_read"`
	TotalFilesModified int              `json:"totalFilesModified" db:"total_files_modified"`
	TotalFilesDeleted  int              `json:"totalFilesDeleted" db:"total_files_deleted"`
	ReasoningSummary   string           `json:"reasoningSummary" db:"reasoning_summary"`
	CreatedAt          time.Time        `json:"createdAt" db:"created_at"`
}

// AiAnalysisLine represents a granular file observation or mutation record.
type AiAnalysisLine struct {
	Id             int64           `json:"id" db:"id"`
	TaskId         string          `json:"taskId" db:"task_id"`
	FilePath       string          `json:"filePath" db:"file_path"`
	OperationType  AiOperationType `json:"operationType" db:"operation_type"`
	StartLine      int             `json:"startLine" db:"start_line"`
	EndLine        int             `json:"endLine" db:"end_line"`
	LineCount      int             `json:"lineCount" db:"line_count"`
	ContentSnippet string          `json:"contentSnippet" db:"content_snippet"`
	Reasoning      string          `json:"reasoning" db:"reasoning"`
	CreatedAt      time.Time       `json:"createdAt" db:"created_at"`
}

// StartTaskOptions configures task initiation.
type StartTaskOptions struct {
	TaskId      string
	Description string
	ModelName   string
	RepoRoot    string
}

// RecordLineOptions configures recording an individual line observation.
type RecordLineOptions struct {
	TaskId    string
	FilePath  string
	Lines     string
	Operation string
	Reasoning string
	Snippet   string
	RepoRoot  string
}
```

---

## 3. Step-by-Step Implementation Plan

### Step 1: Split-DB Section & Path Resolver

- Open `cli/store/split_db_path.go`.
- Add `SectionAiAnalysis = "ai-analysis"` to the section constants.
- Implement `ResolveAiAnalysisSplitDbPath(slug, repoRoot string) string`:
  - Sanitize slug and compute partition directory via `ResolveSplitDbDir(SectionAiAnalysis, slug, repoRoot)`.
  - Return canonical path to `sql.db`.

### Step 2: DDL Migration Engine

- Create `cli/store/ai_analysis_split_db.go`.
- Define migration string containing `AiAnalysisTask` and `AiAnalysisLine` table creation and index definitions.
- Implement `MigrateAiAnalysisSplitDb(db *sql.DB) error`.
- Verify foreign key cascades and WAL configuration.

### Step 3: Database Access Layer

- Create `cli/cmdai/ai_analysis_db.go`.
- Implement `OpenAiAnalysisSplitDB(repoRoot string) (*sql.DB, error)`.
- Implement parameterized queries for task creation, line recording, and status updates.
- Ensure all queries use prepared statements and defensive error handling.

### Step 4: Line Range Parser & Validation Helper

- In `cli/cmdai/ai_analysis_types.go`, implement `ParseLineRange(raw string) (int, int, error)`.
- Support single line (`"42"` -> `42, 42`) and line range (`"10-50"` -> `10, 50`).
- Validate that `startLine > 0` and `endLine >= startLine`.

### Step 5: Cobra Command Integration

- Create `cli/cmdai/ai_analysis_cmd.go`.
- Bind subcommands `start`, `record`, `finish`, `list`, and `inspect` under `analysisCmd`.
- Add flags:
  - `--task`, `-t` (task identifier)
  - `--desc`, `-d` (task description)
  - `--model`, `-m` (AI model name)
  - `--file`, `-f` (relative file path)
  - `--lines`, `-l` (line numbers or range)
  - `--op`, `-o` (operation: read, edit, delete)
  - `--reason`, `-r` (reasoning text)
  - `--summary`, `-s` (completion summary)
- Register `analysisCmd` into `AiCmd` in `cli/cmdai/ai_cmd.go`.

### Step 6: Automated Testing & Verification

- Create `cli/cmdai/ai_analysis_test.go`.
- Write unit tests for:
  - Table migration in memory / temporary directory.
  - Creating a task and retrieving it.
  - Recording single and batch lines with reasoning.
  - Querying lines for a task.
  - Updating task status and completion counters.
  - Cascading deletion when a task is removed.

---

## 4. Acceptance Criteria

- [ ] `ResolveAiAnalysisSplitDbPath` returns a clean relative or configured path under `.gitmap/data/ai-analysis/<slug>/sql.db`.
- [ ] Database initialization applies WAL mode and enforces foreign key constraints.
- [ ] `AiAnalysisTask` correctly tracks `total_files_read`, `total_files_modified`, and `total_files_deleted`.
- [ ] `AiAnalysisLine` stores `operation_type` in `('read', 'edit', 'delete')`, line ranges, and non-empty reasoning.
- [ ] `gitmap ai analysis start` creates a valid task and exits with code 0.
- [ ] `gitmap ai analysis record` validates line ranges and file paths, inserting line telemetry.
- [ ] `gitmap ai analysis finish` updates status to `completed` and records end timestamp.
- [ ] `gitmap ai analysis list` renders recent tasks in an ANSI-formatted table.
- [ ] All unit tests in `cli/cmdai/ai_analysis_test.go` pass with 100% success rate.
- [ ] Zero absolute paths or `file:///` URIs exist in any owned file.

---

## 5. Verification Commands

```bash
  # 1. Run unit tests for AI analysis database and commands
  go test -v ./cli/cmdai -run "TestAiAnalysis.*"

  # 2. Test CLI task initiation
  gitmap ai analysis start --task "test-subtask-01" --desc "Verify AI analysis subsystem" --model "gemini-2.5-pro"

  # 3. Test recording file reads and reasoning
  gitmap ai analysis record --task "test-subtask-01" --file "cli/store/ai_analysis_split_db.go" --lines "1-45" --op "read" --reason "Inspecting DDL migration definitions"

  # 4. Test recording file edits and reasoning
  gitmap ai analysis record --task "test-subtask-01" --file "cli/store/split_db_path.go" --lines "25-30" --op "edit" --reason "Added SectionAiAnalysis constant"

  # 5. Test inspection and listing
  gitmap ai analysis list --limit 5
  gitmap ai analysis inspect "test-subtask-01"

  # 6. Test task completion
  gitmap ai analysis finish --task "test-subtask-01" --status "completed" --summary "Subtask 01 verified successfully"

  # 7. Run repository relative paths linter
  python linter-scripts/check-relative-paths.py
```
