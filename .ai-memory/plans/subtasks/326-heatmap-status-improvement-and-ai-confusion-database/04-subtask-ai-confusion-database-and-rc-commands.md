# Subtask 04: AI Confusion Database, Store & Repo-Cache CLI Commands

## Status: PENDING
- **Parent Plan**: `.ai-memory/plans/pending/326-heatmap-status-improvement-and-ai-confusion-database.md`
- **Specification Reference**: `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/02-component-and-cli-spec.md`
- **Owner**: Worker 02 (Database & CLI Subsystem Engineer)
- **Target Subsystem**: `cli/store`, `cli/cmd`
- **Owned Files**:
  - `cli/store/ai_confusion_types.go` (new file)
  - `cli/store/ai_confusion_schema.go` (new file)
  - `cli/store/ai_confusion_store.go` (new file)
  - `cli/store/ai_confusion_sync.go` (new file)
  - `cli/cmd/rc_ai_confusion.go` (new file)
  - `cli/cmd/special_repos_cli.go`

---

## 1. Context & Objectives

1. **`repo-cache` AI Confusion Database (`cli/store/`)**:
   - Establish a dedicated SQLite database partitioned by repository inside `repo-cache` (`repo-cache/<repo-slug>/ai-confusion.db`).
   - Store comprehensive context for AI agent confusion episodes: agent identity, parent task, subtask, spec path, affected file path, line ranges, snippet, intention, reasoning/thinking traces, and unresolved question.
   - Enforce GitMap Split-DB standards: PascalCase table name (`AiConfusion`), PascalCase column names, positive booleans (`IsCompleted`), WAL journal mode, and 5000ms busy timeout.

2. **Command Suite: `gitmap rc ai-confusion` (`cli/cmd/rc_ai_confusion.go`)**:
   - Provide subcommands:
     - `add`: Record a new confusion episode via terminal flags or `--file <instruction.json>`.
     - `view`: Inspect records filtered by repo, relative file path, completion state, or ID; pretty-print confusion card.
     - `complete`: Mark an episode as resolved with optional resolution notes (`IsCompleted = 1`).
     - `sync`: Execute the 6-stage non-destructive pull/export/replay synchronization.
   - Centralize CLI routing in `cli/cmd/special_repos_cli.go` under `gitmap rc` (`rc ai-confusion`, `rc aic`, `rc confusion`).

3. **6-Stage Safe Synchronization Protocol (`cli/store/ai_confusion_sync.go`)**:
   - Stage 1: Export uncommitted/pending records to a temporary snapshot file.
   - Stage 2: Discard uncommitted binary modifications on `ai-confusion.db`.
   - Stage 3: Pull upstream commits using `git pull --rebase`.
   - Stage 4: Replay and merge snapshot records via `INSERT ... ON CONFLICT(Id) DO UPDATE`.
   - Stage 5: Create atomic chore commit (`gitmap cpc`).
   - Stage 6: Push to upstream remote.

4. **Code Quality & Coding Guideline Discipline**:
   - Maximum 15 lines per function (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses and early returns).
   - Positive boolean naming (`isCompleted`, `hasSnippet`, `isAutoCommit`).
   - AppError wrapping on all error return paths.
   - Zero test / zero build runs during specification authoring.

---

## 2. Step-by-Step Implementation Plan

### Step 4.1: Implement Data Models in `cli/store/ai_confusion_types.go`
1. Define `AiConfusionRecord`:
   ```go
   type AiConfusionRecord struct {
       Id                string     `json:"id"`
       RepoSlug          string     `json:"repoSlug"`
       SequenceNumber    int        `json:"sequenceNumber"`
       AgentSlug         string     `json:"agentSlug"`
       SpecPath          string     `json:"specPath"`
       TaskSlug          string     `json:"taskSlug"`
       SubtaskSlug       string     `json:"subtaskSlug"`
       FilePath          string     `json:"filePath"`
       LineStart         int        `json:"lineStart"`
       LineEnd           int        `json:"lineEnd"`
       CodeSnippet       string     `json:"codeSnippet"`
       AiIntention       string     `json:"aiIntention"`
       AiThinking        string     `json:"aiThinking"`
       ConfusionQuestion string     `json:"confusionQuestion"`
       IsCompleted       bool       `json:"isCompleted"`
       ResolutionNotes   string     `json:"resolutionNotes"`
       CreatedAt         time.Time  `json:"createdAt"`
       UpdatedAt         time.Time  `json:"updatedAt"`
       CompletedAt       *time.Time `json:"completedAt,omitempty"`
   }
   ```
2. Define `AiConfusionFilter`:
   ```go
   type AiConfusionFilter struct {
       RepoSlug    string `json:"repoSlug,omitempty"`
       FilePath    string `json:"filePath,omitempty"`
       IsCompleted *bool  `json:"isCompleted,omitempty"`
       TaskSlug    string `json:"taskSlug,omitempty"`
       Limit       int    `json:"limit,omitempty"`
   }
   ```
3. Define ingestion payload struct for JSON files (`AiConfusionInput`).

### Step 4.2: Implement Schema & Connection Management in `cli/store/ai_confusion_schema.go`
1. Define table schema DDL constant:
   - `CREATE TABLE IF NOT EXISTS AiConfusion (...)` with PascalCase columns.
   - Create indexes on `RepoSlug`, `IsCompleted`, `FilePath`, and `(TaskSlug, SubtaskSlug)`.
2. Implement `ResolveAiConfusionDbPath(repoSlug string) (string, error)`:
   - Look up `repo-cache` root via `ResolveRepoCachePath()`.
   - Ensure directory `repo-cache/<repo-slug>/` exists.
   - Return `filepath.Join(repoCacheDir, repoSlug, "ai-confusion.db")`.
3. Implement `OpenAiConfusionDB(repoSlug string) (*sql.DB, error)`:
   - Open SQLite connection with PRAGMAs: WAL mode, busy_timeout=5000, foreign_keys=ON.
   - Execute table creation DDL.

### Step 4.3: Implement SQLite CRUD Store in `cli/store/ai_confusion_store.go`
1. Define `AiConfusionStore` struct wrapping `*sql.DB` and `dbPath string`.
2. Implement `InsertConfusion(rec AiConfusionRecord) error`:
   - Generate UUID for `Id` if blank.
   - Set `CreatedAt` and `UpdatedAt` to current UTC time.
   - Execute parameterized `INSERT INTO AiConfusion (...) VALUES (...)`.
3. Implement `GetConfusionByID(id string) (*AiConfusionRecord, error)`:
   - Execute parameterized query by `Id`.
   - Scan PascalCase columns into `AiConfusionRecord`.
4. Implement `ListConfusions(filter AiConfusionFilter) ([]AiConfusionRecord, error)`:
   - Dynamically build `WHERE` clause using discrete guard clauses (no string concatenation vulnerabilities).
   - Order by `CreatedAt DESC`.
5. Implement `MarkCompleted(id string, resolutionNotes string) error`:
   - Execute `UPDATE AiConfusion SET IsCompleted = 1, ResolutionNotes = ?, CompletedAt = ?, UpdatedAt = ? WHERE Id = ?`.

### Step 4.4: Implement Safe 6-Stage Synchronization in `cli/store/ai_confusion_sync.go`
1. Implement `SyncAiConfusion(repoSlug string) (*SyncResult, error)`:
   - **Stage 1 (Export)**: Query local records where `IsCompleted = 0` or updated within 24h. Write to temporary JSON snapshot.
   - **Stage 2 (Checkout)**: Run `git -C <repo-cache> checkout -- <dbPath>` to discard binary modifications.
   - **Stage 3 (Rebase Pull)**: Run `git -C <repo-cache> pull --rebase origin main`.
   - **Stage 4 (Replay & Merge)**: Read snapshot JSON. Execute `INSERT INTO AiConfusion ... ON CONFLICT(Id) DO UPDATE SET ...` for each record.
   - **Stage 5 (Commit)**: Stage database file and commit using semantic commit verb:
     `git -C <repo-cache> commit -m "AI_CONFUSION - sync and replay local agent episodes"`
   - **Stage 6 (Push)**: Run `git -C <repo-cache> push origin main`.
2. Return structured `SyncResult` detailing pulled commits and replayed record counts.

### Step 4.5: Implement CLI Commands in `cli/cmd/rc_ai_confusion.go`
1. Define command dispatcher:
   ```go
   func RunRCAiConfusion(args []string) error
   ```
2. Implement subcommands:
   - `handleAiConfusionAdd(args []string) error`: parse `--file` JSON or flags (`--repo`, `--agent`, `--file-path`, `--lines`, `--intention`, `--thinking`, `--question`). Render preview card and persist.
   - `handleAiConfusionView(args []string) error`: parse filter flags (`--repo`, `--file`, `--pending`, `--id`, `--json`). Render formatted terminal cards or compact JSON.
   - `handleAiConfusionComplete(args []string) error`: parse ID and `--notes`. Mark completed in database and print confirmation.
   - `handleAiConfusionSync(args []string) error`: invoke `store.SyncAiConfusion` and print stage-by-stage status.
3. Implement help renderer:
   ```go
   func RenderAiConfusionHelp() string
   ```

### Step 4.6: Wire Routing in `cli/cmd/special_repos_cli.go`
1. In `dispatchSpecialRepoSubcmd(shortKey string, opts specialCLIOptions)`:
   - Add check: if `shortKey == "rc"` and `isAiConfusionSubcmd(opts.subcmd)`:
     ```go
     if shortKey == "rc" && isAiConfusionSubcmd(opts.subcmd) {
         return RunRCAiConfusion(opts.extraArgs)
     }
     ```
2. Define `isAiConfusionSubcmd(sub string) bool`:
   - Match `ai-confusion`, `aic`, `confusion`.
3. Update `RenderSpecialRepoHelp("rc")` to include `gitmap rc ai-confusion [add|view|complete|sync]`.

---

## 3. Acceptance Criteria & Quality Gates

- [ ] SQLite database is created at `repo-cache/<repo-slug>/ai-confusion.db` with PascalCase `AiConfusion` table and columns.
- [ ] Positive booleans used exclusively (`isCompleted`, `hasSnippet`).
- [ ] `gitmap rc ai-confusion add --file <path.json>` ingests structured JSON and renders a preview card.
- [ ] `gitmap rc ai-confusion view` displays syntax-highlighted / boxed confusion details.
- [ ] `gitmap rc ai-confusion complete <id>` sets `IsCompleted = 1` and records `CompletedAt`.
- [ ] `gitmap rc ai-confusion sync` executes the 6-stage non-destructive export/pull/replay cycle.
- [ ] All Go functions adhere to the 15-line limit (<= 8 lines preferred).
- [ ] Zero build or test commands executed during this spec authoring phase.
