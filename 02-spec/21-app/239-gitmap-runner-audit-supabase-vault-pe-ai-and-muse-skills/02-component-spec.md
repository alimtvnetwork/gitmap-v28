# Component Specification: Universal File Runner, Supabase Vault, Chore Commits, PE-AI & Muse Skills

> **Specification ID:** 239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills  
> **Status:** Ratified  
> **Version:** 1.0.0  
> **Target Subsystems:**  
> - `cli/cmdrun/` (Universal file runner & task audit / errors DB)  
> - `cli/cmdsupabase/` (Multi-Supabase encrypted vault & one-liner CLI)  
> - `cli/cmd/` & `cli/constants/` (`cpc` chore commit & full forms help)  
> - `cli/cmdpipeline/` (`--ai`, `-ud`, `-t` 2-min abort, log source details, instant hash)  
> - `01-prompts/` & `.agents/skills/` (Muse master prompt & skills modernization)  

---

## 1. Subsystem: Universal File Runner (`cli/cmdrun`)

### 1.1 Overview & Responsibilities
The `cli/cmdrun` package provides an autonomous universal file runner. When invoked via `gitmap run <file> [args...]`, it:
1. Resolves the target file path relative to current working directory or repository root.
2. Identifies file extension (`.py`, `.ps1`, `.sh`, `.js`, `.ts`, `.go`).
3. Auto-discovers the verified interpreter runtime (Python, PowerShell Core/Desktop, Bash/POSIX sh, Node/Bun/Deno, Go).
4. Enqueues an execution task record in the SQLite Task Audit DB (`status = "queued"` -> `"running"`).
5. Executes the file with passed arguments, streaming stdout/stderr in real-time.
6. Upon exit code 0, marks the task `"completed"`.
7. Upon non-zero exit code or runtime panic, marks the task `"failed"`, persists execution errors into the dedicated Run Errors SQLite DB, and displays error diagnostics.
8. Exposes `gitmap run errors`, `gitmap run-errors`, and `gitmap run history` to review and diagnose past failures.

### 1.2 Data Structures & Types

```go
package cmdrun

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SupportedRunExtension enumerates supported universal runner file extensions.
type SupportedRunExtension string

const (
	ExtPython     SupportedRunExtension = ".py"
	ExtPowerShell SupportedRunExtension = ".ps1"
	ExtBash       SupportedRunExtension = ".sh"
	ExtJavaScript SupportedRunExtension = ".js"
	ExtTypeScript SupportedRunExtension = ".ts"
	ExtGolang     SupportedRunExtension = ".go"
)

// RunTaskRecord represents an enqueued file execution in the audit task database.
type RunTaskRecord struct {
	Id          string    `json:"id"`
	FilePath    string    `json:"filePath"`
	Interpreter string    `json:"interpreter"`
	Args        string    `json:"args"`
	Status      string    `json:"status"` // "queued", "running", "completed", "failed"
	ExitCode    int       `json:"exitCode"`
	ErrorId     string    `json:"errorId,omitempty"`
	EnqueuedAt  time.Time `json:"enqueuedAt"`
	StartedAt   time.Time `json:"startedAt,omitempty"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	DurationMs  int64     `json:"durationMs,omitempty"`
}

// RunErrorRecord represents a captured execution failure in the errors database.
type RunErrorRecord struct {
	Id           string    `json:"id"`
	RunTaskId    string    `json:"runTaskId"`
	FilePath     string    `json:"filePath"`
	Interpreter  string    `json:"interpreter"`
	ExitCode     int       `json:"exitCode"`
	ErrorMessage string    `json:"errorMessage"`
	StackTrace   string    `json:"stackTrace,omitempty"`
	StdErrTail   string    `json:"stdErrTail,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

// RunCommandOptions encapsulates CLI options for gitmap run.
type RunCommandOptions struct {
	TargetFile  string   `json:"targetFile"`
	Args        []string `json:"args"`
	IsJSON      bool     `json:"isJSON"`
	IsDryRun    bool     `json:"isDryRun"`
	Limit       int      `json:"limit"`
	IsClear     bool     `json:"isClear"`
	FilterExt   string   `json:"filterExt,omitempty"`
}
```

### 1.3 Database Schemas (SQLite Split-DB)

Stored in `.gitmap/data/installation/run_audit.db` and `.gitmap/data/installation/run_errors.db`:

```sql
-- Run Audit Tasks Schema (.gitmap/data/installation/run_audit.db)
CREATE TABLE IF NOT EXISTS RunAuditTasks (
    Id TEXT PRIMARY KEY,
    FilePath TEXT NOT NULL,
    Interpreter TEXT NOT NULL,
    Args TEXT,
    Status TEXT NOT NULL, -- 'queued', 'running', 'completed', 'failed'
    ExitCode INTEGER DEFAULT 0,
    ErrorId TEXT,
    EnqueuedAt DATETIME NOT NULL,
    StartedAt DATETIME,
    CompletedAt DATETIME,
    DurationMs INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_run_audit_status ON RunAuditTasks(Status);
CREATE INDEX IF NOT EXISTS idx_run_audit_filepath ON RunAuditTasks(FilePath);
CREATE INDEX IF NOT EXISTS idx_run_audit_enqueued ON RunAuditTasks(EnqueuedAt DESC);

-- Run Errors Schema (.gitmap/data/installation/run_errors.db)
CREATE TABLE IF NOT EXISTS RunErrors (
    Id TEXT PRIMARY KEY,
    RunTaskId TEXT NOT NULL,
    FilePath TEXT NOT NULL,
    Interpreter TEXT NOT NULL,
    ExitCode INTEGER NOT NULL,
    ErrorMessage TEXT NOT NULL,
    StackTrace TEXT,
    StdErrTail TEXT,
    CreatedAt DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_run_errors_filepath ON RunErrors(FilePath);
CREATE INDEX IF NOT EXISTS idx_run_errors_created ON RunErrors(CreatedAt DESC);
CREATE INDEX IF NOT EXISTS idx_run_errors_task_id ON RunErrors(RunTaskId);
```

### 1.4 Function Signatures & Contract

```go
// RunUniversalCommand dispatches universal runner subcommands and file targets.
func RunUniversalCommand(args []string) error

// ResolveInterpreter maps file extension to verified system runtime binary.
func ResolveInterpreter(filePath string) (interpreter string, prefixArgs []string, err *apperror.AppError)

// ExecuteFileWithAudit wraps execution, DB logging, and error recording.
func ExecuteFileWithAudit(opts RunCommandOptions) error

// EnqueueRunTask writes initial queued task to RunAuditTasks.
func EnqueueRunTask(dbPath string, record RunTaskRecord) (*RunTaskRecord, error)

// CompleteRunTask updates task status, exit code, and completion timestamp.
func CompleteRunTask(dbPath string, taskId string, exitCode int, errorId string) error

// RecordRunError persists execution failure into RunErrors.
func RecordRunError(dbPath string, rec RunErrorRecord) error

// ListRunErrors queries recent run errors for diagnostic display.
func ListRunErrors(limit int, isJSON bool) error

// ListRunHistory queries recent run executions.
func ListRunHistory(limit int, isJSON bool) error

// ClearRunErrors purges recorded run errors from the database.
func ClearRunErrors() error
```

### 1.5 CLI Usage & Flag Contracts
- `gitmap run <filepath> [args...]` — Detect extension and execute file with audit.
- `gitmap run errors [--limit <n>] [--json]` — Display failed executions with stack traces.
- `gitmap run history [--limit <n>] [--json]` — Display execution audit trail.
- `gitmap run errors --clear` — Purge recorded run errors.
- `gitmap run-errors` — Shortcut alias for `gitmap run errors`.

---

## 2. Subsystem: Multi-Supabase Database Vault (`cli/cmdsupabase`)

### 2.1 Overview & Responsibilities
The `cli/cmdsupabase` package provides multi-project Supabase connection management. Multiple machine nodes can store and coordinate Supabase configurations securely without exposing credentials.
- **Strict Zero Cleartext Invariant:** Anon keys, service role keys, and database connection strings (`postgres://...`) MUST NEVER be stored in plaintext. All secrets are encrypted using AES-256-GCM before writing to the local SQLite database.
- **One-Liner Registration:** Users and scripts register projects via `gitmap supabase add <alias> <url> <anon_key> <service_key> [db_url]`.
- **Vault Location:** Stored in `.gitmap/data/installation/supabase/sql.db`.

### 2.2 Data Structures & Types

```go
package cmdsupabase

import (
	"time"
)

// SupabaseProjectRecord represents an encrypted Supabase configuration in SQLite.
type SupabaseProjectRecord struct {
	Alias               string    `json:"alias"`
	Url                 string    `json:"url"`
	AnonKeyEncrypted    string    `json:"anonKeyEncrypted"`
	ServiceKeyEncrypted string    `json:"serviceKeyEncrypted"`
	DbUrlEncrypted      string    `json:"dbUrlEncrypted,omitempty"`
	Nonce               string    `json:"nonce"`
	KeyFingerprint      string    `json:"keyFingerprint"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// SupabaseProjectView represents masked display output for CLI terminal / JSON.
type SupabaseProjectView struct {
	Alias          string `json:"alias"`
	Url            string `json:"url"`
	AnonKeyMasked  string `json:"anonKeyMasked"`
	HasServiceKey  bool   `json:"hasServiceKey"`
	HasDbUrl       bool   `json:"hasDbUrl"`
	KeyFingerprint string `json:"keyFingerprint"`
	UpdatedAt      string `json:"updatedAt"`
}

// SupabaseAddOptions holds arguments parsed from CLI command.
type SupabaseAddOptions struct {
	Alias      string `json:"alias"`
	Url        string `json:"url"`
	AnonKey    string `json:"anonKey"`
	ServiceKey string `json:"serviceKey"`
	DbUrl      string `json:"dbUrl,omitempty"`
}
```

### 2.3 Database Schema (SQLite Split-DB)

Stored in `.gitmap/data/installation/supabase/sql.db`:

```sql
CREATE TABLE IF NOT EXISTS SupabaseProjects (
    Alias TEXT PRIMARY KEY,
    Url TEXT NOT NULL,
    AnonKeyEncrypted TEXT NOT NULL,
    ServiceKeyEncrypted TEXT NOT NULL,
    DbUrlEncrypted TEXT,
    Nonce TEXT NOT NULL,
    KeyFingerprint TEXT NOT NULL,
    CreatedAt DATETIME NOT NULL,
    UpdatedAt DATETIME NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_supabase_alias ON SupabaseProjects(Alias);
```

### 2.4 Vault Cryptography Contract
1. Key derivation uses machine-unique machine ID / RSA vault token combined with local salt.
2. Cipher: AES-256-GCM authenticated encryption.
3. Each row receives a distinct cryptographic 12-byte nonce (hex-encoded).
4. `EncryptVaultSecret(plaintext string) (ciphertextHex, nonceHex string, err error)`
5. `DecryptVaultSecret(ciphertextHex, nonceHex string) (plaintext string, err error)`

### 2.5 Function Signatures & Contract

```go
// RunSupabaseCommand routes gitmap supabase subcommands.
func RunSupabaseCommand(args []string) error

// AddSupabaseProject encrypts secrets and inserts/updates project in SQLite.
func AddSupabaseProject(opts SupabaseAddOptions) error

// ListSupabaseProjects outputs registered projects with masked secrets.
func ListSupabaseProjects(isJSON bool) error

// RemoveSupabaseProject deletes project entry by alias.
func RemoveSupabaseProject(alias string) error

// TestSupabaseProject verifies REST / Auth endpoint reachability.
func TestSupabaseProject(alias string) error

// GetDecryptedProject returns decrypted credentials for authenticated operations.
func GetDecryptedProject(alias string) (*SupabaseProjectRecord, error)
```

### 2.6 CLI Commands & Flags
- `gitmap supabase add <alias> <url> <anon_key> <service_key> [db_url]` — Register Supabase database with encrypted credentials.
- `gitmap supabase list [--json]` — Display registered Supabase projects with masked keys.
- `gitmap supabase remove <alias>` (or `rm`) — Delete registered Supabase configuration.
- `gitmap supabase test <alias>` — Test REST ping and authentication against Supabase URL.
- `gitmap supabase help` — Two-column help menu documenting syntax and encrypted vault safety.

---

## 3. Subsystem: Chore Commits (`cpc`) & Help Menu Full Forms (`cli/cmd`, `cli/constants`)

### 3.1 Overview & Responsibilities
GitMap provides semantic commit helpers (`cpf`, `cpb`, `cpr`, `pcp`). Chore commits are now standardized with:
- `gitmap commit-push-chore "<module> - <summary>"`
- `gitmap cpc "<module> - <summary>"` (alias)
- Commit message auto-prefixed with `Chore: `.
- All CLI help menus, help text documents, and command catalogs are expanded to display full forms of all semantic commit and transfer shortcuts (`cpf` = commit-push-feature, `cpb` = commit-push-bug, `cpc` = commit-push-chore, `cpr` = commit-push-release, `pcp` = pull-commit-push, `pas` = pull-all-ssh, `pae` = pull-all-efficient).

### 3.2 Constants (`cli/constants/constants_cli.go`)

```go
const (
	// CmdCommitPushChore stages all changes, commits with "Chore: " prefix, and pushes.
	CmdCommitPushChore      = "commit-push-chore"
	CmdCommitPushChoreAlias = "cpc"
)
```

### 3.3 Function Signatures (`cli/cmd/commit_push.go`)

```go
// runCommitPushChore stages all changes, commits with "Chore: " prefix, and pushes.
func runCommitPushChore(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushChore, []string{"--help"})
		return nil
	}
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-chore \"<what chore was done>\"", "E9000")
	}
	commitMessage := "Chore: " + strings.Join(args, " ")
	return executeCommitPush(commitMessage)
}
```

### 3.4 CLI Help Menu Updates (`cli/cmd/commit_help_menu.go`)

```go
func buildCommitAIWorkflowSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Commit & AI Workflow Shortcuts (Full Forms)",
		Entries: []termhelp.CommandEntry{
			{Command: "commit, cm <msg>", Description: "commit: flat git commit with auto-stage (git add -A)"},
			{Command: "cpf <msg>", Description: "commit-push-feature: stage all, commit feat, and push"},
			{Command: "cpb <msg>", Description: "commit-push-bug: stage all, commit bugfix, and push"},
			{Command: "cpc <msg>", Description: "commit-push-chore: stage all, commit maintenance chore, and push"},
			{Command: "cpr <msg>", Description: "commit-push-release: stage, commit release chore, and push"},
			{Command: "pcp <msg>", Description: "pull-commit-push: pull latest rebase, stage, commit, and push"},
			{Command: "pas", Description: "pull-all-ssh: pull all repositories using SSH transport"},
		},
	}
}
```

### 3.5 Help Markdown Document (`cli/helptext/commit-push-chore.md`)
Created in `cli/helptext/commit-push-chore.md` detailing purpose, shell command comparison, alias (`cpc`), and examples.

---

## 4. Subsystem: Pipeline PE-AI, Until-Done, 2-Min Poll & Source Details (`cli/cmdpipeline`)

### 4.1 Overview & Responsibilities
Enhances `gitmap pe` (pipeline error inspection) and `gitmap pe -t` (telemetry mode):
1. **`--ai` Flag (Clipboard & Memory Suppression):** When `--ai` is supplied, suppresses all clipboard copying (`writeClipboard`), suppresses `📋 Copied...` terminal notices, and emits clean stdout/stderr for automated AI agent parsing without clobbering system clipboard.
2. **`pe -t` (2-Minute Error Abort Loop):** Polls every 2 minutes. If errors/stack traces are detected in logs during active workflow execution, immediately showcases failure summary via `gitmap pe` and terminates execution instead of waiting for timeout.
3. **`-ud` / `--until-done` Flag:** Continuous watch mode running until the entire CI/CD run concludes (all jobs succeeded or failed), then surfaces final status and diagnostics.
4. **Immediate Git Hash Retrieval:** If the target or HEAD commit SHA is already stored in the local SQLite database and all runs for that commit are complete, retrieves telemetry directly from SQLite without network delays.
5. **Enhanced Log Source Details:** Displays source metadata for every failure: logger name, workflow name, run URL, job name, step name, and relative saved log file path.

### 4.2 Data Structures & Types (`cli/cmdpipeline/pipeline_flags.go`, `cli/cmdpipeline/pipeline.go`)

```go
package cmdpipeline

// PipelineErrorFlags additions
type PipelineErrorFlags struct {
	// ... existing fields ...
	HasAI        bool // Parsed from --ai
	HasUntilDone bool // Parsed from -ud, --until-done
}

// SectionFailure additions
type SectionFailure struct {
	LoggerName     string   `json:"loggerName,omitempty"` // Captured logger / runner agent name
	WorkflowName   string   `json:"workflowName"`
	RunId          uint64   `json:"runId"`
	RunUrl         string   `json:"runUrl,omitempty"`     // Web URL to GitHub Actions run
	JobName        string   `json:"jobName"`
	StepName       string   `json:"stepName"`
	FailureSummary string   `json:"failureSummary"`
	ErrorLines     []string `json:"errorLines"`
	Warnings       []string `json:"warnings,omitempty"`
	SavedLogFile   string   `json:"savedLogFile,omitempty"` // Strictly relative path
	CreatedAt      string   `json:"createdAt,omitempty"`
	StackTrace     string   `json:"stackTrace,omitempty"`
}
```

### 4.3 Immediate Git Hash Retrieval Contract (`cli/cmdpipeline/pipeline_cache_eval.go`)

```go
// CheckCompletedCommitShaInstant retrieves cached telemetry instantly if HEAD or target SHA is completed in DB.
func CheckCompletedCommitShaInstant(db *pipelinedb.PipelineSplitDb, repo string, targetSha string) (PipelineCacheDecision, bool) {
	if len(targetSha) == 0 {
		targetSha = ResolveRepoHeadCommitSha(repo)
	}
	if len(targetSha) == 0 {
		return PipelineCacheDecision{}, false
	}
	runRes := db.QueryRunsBySha(targetSha)
	if runRes.IsFailure() || len(runRes.Data) == 0 {
		return PipelineCacheDecision{}, false
	}
	if isCommitRunsCompleted(runRes.Data, targetSha) && !hasAnyActiveRun(runRes.Data, targetSha) {
		return buildCacheHitDecision(runRes.Data, targetSha, "instant_sha_completed_hit"), true
	}
	return PipelineCacheDecision{}, false
}
```

### 4.4 2-Minute Poll & `--until-done` Execution Loop (`cli/cmdpipeline/pipeline_dynamic_timeline.go`)

```go
func watchTimelineWithInterval(repo, workflowName string, isUntilDone bool, isAI bool) error {
	const pollInterval = 120 * time.Second // 2 minutes
	for {
		runs := queryWorkflowRuns(repo)
		active := findActiveWorkflowRun(runs)
		
		// If an active run has error traces detected in log stream, showcase immediately and exit
		if active != nil && hasActiveRunErrors(repo, active.DatabaseId) && !isUntilDone {
			fmt.Printf("\n  %s✖ Failure detected in active log stream [%s #%d]. Showcasing errors immediately:%s\n\n",
				constants.ColorRed, active.Name, active.DatabaseId, constants.ColorReset)
			return executePipelineErrorLogs([]string{"--no-cache"})
		}
		
		if active == nil {
			// All runs completed
			return executePipelineErrorLogs([]string{})
		}
		
		time.Sleep(pollInterval)
	}
}
```

---

## 5. Subsystem: Muse Master Prompt & Skills Modernization

### 5.1 Overview & Responsibilities
Modernizes the canonical Muse AI master prompt and repository skills to enforce native GitMap tool primacy:
1. **`01-prompts/27-muse-prompts/01-muse-master-prompt.md`:**
   - Mandate `gitmap aum search` over slow shell grep/Select-String.
   - Mandate `gitmap task` for SQLite multi-agent state tracking.
   - Mandate `gitmap run` for universal file execution.
   - Mandate `gitmap pe --ai` for CI/CD diagnostics without clipboard clobbering.
   - Mandate `gitmap sync` for fleet synchronization across 43 repositories.
   - Mandate `gitmap cpc` for chore commits.
2. **Skills Modernization (`.agents/skills/gitmap/SKILL.md`, `.cursor/skills/gitmap/skill.md`, `.agents/skills/muse-master-prompt/skill.md`, `.cursor/skills/muse-master-prompt/skill.md`):**
   - Incorporate `gitmap run`, `gitmap supabase`, `gitmap cpc`, `gitmap pe --ai`, `gitmap pe -ud`.
3. **Coding Guidelines Modernization (`02-spec/02-coding-guidelines/06-ai-optimization/`):**
   - Add rule `AH-GM1: Mandate gitmap pe --ai for AI Error Inspection`.
   - Add rule `AH-GM2: Ban Clipboard Side-Effects in Automated Pipelines`.
   - Update quick-reference validation checklists.

---

## 6. Acceptance Criteria Matrix

| Subsystem | Requirement | Acceptance Criteria |
| :--- | :--- | :--- |
| `cmdrun` | Universal execution | Resolves `.py`, `.ps1`, `.sh`, `.js`, `.ts`, `.go` to verified interpreters. |
| `cmdrun` | Audit logging | Tasks enqueued in `RunAuditTasks`, marked completed/failed on process termination. |
| `cmdrun` | Errors DB | Failures persisted in `RunErrors` with stack trace and exit code; viewable via `gitmap run errors`. |
| `cmdsupabase` | Multi-project vault | Registers Supabase projects via `gitmap supabase add <alias> <url> <anon> <service> [db_url]`. |
| `cmdsupabase` | Encryption | Zero cleartext secrets; AES-256-GCM encryption with machine nonce. |
| `cmd` | Chore commits | `gitmap cpc "<msg>"` creates commit with `Chore: ` prefix and pushes. |
| `cmd` | Full forms help | Help menus and markdown docs display full forms of `cpf`, `cpb`, `cpc`, `cpr`, `pcp`, `pas`. |
| `cmdpipeline` | `--ai` flag | `gitmap pe --ai` suppresses all clipboard copy operations and clipboard notices. |
| `cmdpipeline` | `pe -t` 2-min abort | Polls every 2 minutes; aborts immediately to show errors when log failures are detected. |
| `cmdpipeline` | `-ud` until done | Polls until all CI/CD jobs finish before displaying status. |
| `cmdpipeline` | Instant hash | Pulls instantly from SQLite if commit SHA runs are completed in database. |
| `cmdpipeline` | Source details | Emits logger name, workflow name, run URL, and relative log file path in failure cards. |
| Prompts/Skills | Muse & GitMap skills | Updated with native GitMap commands, `pe --ai` mandate, and coding guidelines alignment. |
