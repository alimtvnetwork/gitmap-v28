# Subtask 03: Split SQLite Isolation & AI Analysis Data Layer (`ai-analysis.db`)

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `03-ai-analysis-split-db`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md` (Section 2.1 & Section 4)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Design and implement the dedicated **Split SQLite Storage Engine** for GitMap AI analysis tasks, file tracking, line-by-line AST/refactoring diffs, and architectural justifications.

### Architectural Invariant Alignment
- **Split SQLite Isolation:** Persistence is strictly segregated into `.gitmap/data/ai-analysis/ai-analysis.db`. Never mix or inject tables into `gitmap.db`, `installation.db`, or `repodb/pipeline.db`.
- **Database Conventions:** PascalCase table names (`AiTask`, `AiTaskFile`, `AiTaskLine`), primary key `{Table}Id`, positive affirmative booleans (`IsActive`, `HasFailed`, `IsModified`, `IsRemoved`, `IsApplied`), and standard documentation columns (`Description`, `Notes`, `Comments`).
- **Relative Path Hygiene:** The `RelPath` column in `AiTaskFile` must enforce clean forward-slash relative Git repository paths (`02-spec/...`, `cmd/...`, `cli/...`).
- **Concurrency & WAL:** Configured with `PRAGMA journal_mode=WAL;`, `PRAGMA synchronous=NORMAL;`, and `SetMaxOpenConns(1)` to eliminate SQLite locking contention.

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/store/ai_analysis_models.go` | Create | Go struct models for `AiTask`, `AiTaskFile`, and `AiTaskLine` with JSON tags and helper constructors. |
| `cli/store/ai_analysis_split_db.go` | Create | Database connection provider, directory resolver, schema DDL migrator, and WAL configuration. |
| `cli/store/ai_analysis_ops.go` | Create | High-performance CRUD routines (Insert task, upsert file records, batch insert diff lines, query active tasks). |
| `cli/store/ai_analysis_split_db_test.go` | Create | Comprehensive unit tests validating schema creation, foreign key cascading, WAL mode, and CRUD atomicity. |

---

## 3. Database Schema & Data Models

### 3.1 Tables Overview
1. **`AiTask`**: Root record for an autonomous AI task or session.
   - `AiTaskId` (PK), `TaskUuid` (UNIQUE), `Goal`, `Status`, `Category`, `Reasoning`, `TotalFiles`, `TotalLines`, `IsActive`, `HasFailed`, `Description`, `Notes`, `Comments`, `CreatedAt`, `CompletedAt`.
2. **`AiTaskFile`**: Catalogs every file analyzed, modified, staged, or removed.
   - `AiTaskFileId` (PK), `AiTaskId` (FK), `RelPath`, `AbsPath`, `Action`, `BeforeSha256`, `AfterSha256`, `BackupPath`, `LineCount`, `Reasoning`, `IsModified`, `IsRemoved`, `Description`, `Notes`, `Comments`, `CreatedAt`, `UpdatedAt`.
3. **`AiTaskLine`**: Detailed diff ledger recording line mutations and rule citations.
   - `AiTaskLineId` (PK), `AiTaskFileId` (FK), `LineNumber`, `OriginalContent`, `ProposedContent`, `DiffKind`, `RuleViolation`, `Reasoning`, `IsApplied`, `Description`, `Notes`, `Comments`, `CreatedAt`.

### 3.2 CRUD Methods in `cli/store/ai_analysis_ops.go`
- `GetAiAnalysisDb(repoRoot string) (*sql.DB, error)`
- `InsertAiTask(ctx context.Context, task *AiTask) (int64, error)`
- `GetAiTaskByUuid(ctx context.Context, uuid string) (*AiTask, error)`
- `GetAiTaskById(ctx context.Context, taskId int64) (*AiTask, error)`
- `UpdateAiTaskStatus(ctx context.Context, taskId int64, status string, isFailed bool) error`
- `InsertAiTaskFileBatch(ctx context.Context, files []AiTaskFile) error`
- `UpdateAiTaskFileRemoval(ctx context.Context, fileId int64, isRemoved bool, action string) error`
- `InsertAiTaskLineBatch(ctx context.Context, lines []AiTaskLine) error`
- `ListAiTasks(ctx context.Context, isActiveOnly bool, limit int) ([]AiTask, error)`
- `GetAiTaskFiles(ctx context.Context, taskId int64) ([]AiTaskFile, error)`
- `GetAiTaskLines(ctx context.Context, fileId int64) ([]AiTaskLine, error)`

---

## 4. Acceptance Criteria

- [ ] `ai-analysis.db` is initialized under `.gitmap/data/ai-analysis/ai-analysis.db` without polluting root or other Split-DB instances.
- [ ] Schema migration executes automatically on first access with WAL mode enabled.
- [ ] All table names strictly PascalCase (`AiTask`, `AiTaskFile`, `AiTaskLine`) with `{Table}Id` PK.
- [ ] Boolean fields strictly use affirmative naming (`IsActive`, `HasFailed`, `IsModified`, `IsRemoved`, `IsApplied`).
- [ ] Cascading deletion verified: deleting an `AiTask` purges all related `AiTaskFile` and `AiTaskLine` records.
- [ ] Unit tests in `cli/store/ai_analysis_split_db_test.go` pass with 100% coverage of core CRUD workflows.

---

## 5. Verification Commands

```powershell
# 1. Run unit tests for AI Analysis Split-DB
go test -v ./cli/store -run "TestAiAnalysis"

# 2. Verify schema integrity with temporary SQLite instance
go test -v ./cli/store -run "TestAiAnalysisSchemaIntegrity"
```
