# Architecture Specification: Code Refactoring and Coding Guideline Update

> **Task Slug:** `294-code-refactoring-and-coding-guideline-update`  
> **Specification Version:** 1.0.0  
> **Status:** Approved / Ready for Implementation  
> **Target Path:** `02-spec/21-app/294-code-refactoring-and-coding-guideline-update/01-architecture-spec.md`  

---

## 1. Executive Summary & Architectural Goals

The GitMap v28 codebase operates as a high-velocity developer CLI and autonomous agent orchestration platform. To preserve system maintainability, runtime stability, and automated agent reliability, this specification establishes comprehensive architectural refactorings targeting several core domains:

1. **Elimination of Magic Strings & Token Normalization:**
   - Centralize IDE token registration and string joining in `cli/cmdide/ide_add.go` using strongly typed target constants and a dedicated `strutil.JoinCSV` utility.
   - Eliminate hardcoded environment keys (`"GITMAP_WORKDIR"`) and fallback directories (`"workspace"`, `"repos"`) in `cli/cmdpe/pe_target.go` by creating `cli/constants/constants_workdir.go` and `cli/envutil/envutil.go`.

2. **Repository Pattern Encapsulation & Error Sovereignty in `cli/store`:**
   - Decouple inline raw SQL strings from `cli/store/repo.go` and `cli/store/repo_query_filter.go` into centralized SQL constants in `cli/constants/constants_store.go`.
   - Implement `cli/store/repo_repository.go` as a dedicated repository component providing typed query execution and returning `*apperror.AppError` instead of generic untyped errors.

3. **SQL Decoupling & Parsing Hygiene in `cli/cmdpipeline`:**
   - Eliminate raw ad-hoc database queries (`"SELECT AbsolutePath FROM Repo..."`) and manual `sql.Open` connection cycles in `cli/cmdpipeline/pipeline_history_render.go`.
   - Route path resolution through the store repository layer and sanitize Git tag parsing using `strutil` helpers.

4. **Temporal Standardization via `timeutil`:**
   - Standardize UTC ISO8601/RFC3339 timestamp generation across CLI components by introducing `cli/timeutil/timeutil.go`.

5. **Machine-Readable Agent Telemetry:**
   - Extend agent subtask lifecycle commands (`subtask complete`, `subtask start`, `subtask fail`) with structured JSON telemetry flags and payloads.

---

## 2. Invariants & Architecture Constraints

All implementations guided by this specification MUST adhere strictly to the following invariants:

### 2.1 Coding Standards & Invariants
- **15-Line Function Limit:** Every function must be <= 15 lines of executable code (<= 8 lines preferred). Decompose larger operations into focused helper functions.
- **Zero Nested Conditionals:** Guard clauses and early returns must be applied universally. Nested `if` blocks are strictly forbidden.
- **Positive Booleans:** Boolean identifiers and parameters must begin with `is` or `has` only (e.g. `isValid`, `hasError`, `isTargeted`). Negative booleans (e.g. `isNotValid`, `unmatched`) and alternative prefixes (`can`, `should`, `was`, `will`) are prohibited.
- **Leaf Error Sovereignty:** Leaf functions returning an error MUST bubble errors up (`return err` or `return apperror.WrapSimple(err, ...)`). Leaf functions must never call `os.Exit` or `cliexit.HandleError`. Outer orchestrators alone decide termination.
- **Whitespace & Formatting:** Blank line after closing brace `}` and blank line before `return` statements.

### 2.2 Operational & Tooling Constraints
- **Relative Git Paths Only:** All file references in documentation, code, and specifications must use strictly relative repository paths (e.g. `cli/strutil/strutil.go`). Absolute paths and URI schemes are strictly prohibited.
- **Strictly Lowercase Paths:** All directory and file paths must remain entirely lowercase.
- **Zero Git Commands:** Autonomous refactoring agents must not invoke Git CLI commands directly during code modification phases.
- **Zero Compiler Builds & Zero Tests:** In compliance with repository safety gates, `go build` and `go test` invocations are strictly prohibited during autonomous authoring and refactoring phases.

---

## 3. Detailed Component Architecture

### 3.1 IDE Tokens & String Utility (`cli/strutil` and `cli/cmdide`)

#### 3.1.1 Problem Statement
In `cli/cmdide/ide_add.go`, `formatRegisteredTokens` builds comma-delimited registration tokens by manually accumulating strings and calling standard library `strings.Join(tokens, ",")`. Furthermore, inconsistent target representation mixes enum values (`TargetDesktop`) and literal mappings.

#### 3.1.2 CSV Joiner Utility (`cli/strutil/strutil.go`)
Extend `cli/strutil/strutil.go` with a deterministic, allocation-efficient CSV joiner:
```go
// JoinCSV joins a slice of strings with a comma separator, skipping empty items.
func JoinCSV(tokens []string) string
```
- Skips empty or whitespace-only tokens.
- Preserves token order.
- Returns empty string `""` when input slice contains zero non-empty tokens.
- Max function length: <= 12 lines.

#### 3.1.3 Token Accumulator Refactoring (`cli/cmdide/ide_add.go`)
Refactor `formatRegisteredTokens(res IDEActionResult) string`:
- Accumulate target tokens using strongly typed constants (`TargetVSCode`, `TargetCursor`, `TargetAntigravity`, `TargetGitHubDesktop`).
- Delegate joining to `strutil.JoinCSV(tokens)`.
- Enforce early guard clauses:
```go
func formatRegisteredTokens(res IDEActionResult) string {
	tokens := collectRegisteredTokens(res)

	return strutil.JoinCSV(tokens)
}
```

---

### 3.2 Workdir Constants & Environment Utility (`cli/constants` and `cli/envutil`)

#### 3.2.1 Problem Statement
In `cli/cmdpe/pe_target.go`, `ResolveWorkdir` inspects the environment using raw string literals:
```go
if envVal := strings.TrimSpace(os.Getenv("GITMAP_WORKDIR")); envVal != "" {
    return envVal
}
```
In `fallbackWorkdir`:
```go
return filepath.Join(home, "workspace", "repos")
```
Hardcoded directory segments `"workspace"` and `"repos"` bypass configuration registries and prevent unified overrides.

#### 3.2.2 Central Workdir Constants (`cli/constants/constants_workdir.go`)
Create `cli/constants/constants_workdir.go`:
```go
package constants

// Environment variable keys for workspace resolution.
const (
	EnvWorkdir = "GITMAP_WORKDIR"
)

// Default folder names for workspace and repository storage.
const (
	DefaultWorkdirWorkspace = "workspace"
	DefaultWorkdirRepos     = "repos"
)
```

#### 3.2.3 Environment Utility Package (`cli/envutil/envutil.go`)
Create a dedicated package `cli/envutil` containing whitespace-trimmed environment helpers:
- `GetTrimmed(key string) string`: returns `strings.TrimSpace(os.Getenv(key))`.
- `GetOrDefault(key, fallback string) string`: returns `GetTrimmed(key)` if non-empty, otherwise returns `fallback`.

#### 3.2.4 Target Workdir Resolution Refactoring (`cli/cmdpe/pe_target.go`)
Refactor `ResolveWorkdir` and `fallbackWorkdir`:
- Utilize `envutil.GetTrimmed(constants.EnvWorkdir)`.
- Replace inline directory strings with `constants.DefaultWorkdirWorkspace` and `constants.DefaultWorkdirRepos`.
- Add documentation comments detailing fallback semantics for enterprise workspace layouts.

---

### 3.3 Centralized SQL & Store Repository Pattern (`cli/store`)

#### 3.3.1 Problem Statement
Currently, SQL query strings are scattered across implementation files:
- `cli/store/repo.go` defines `sqlSelectReposByBranch` and `sqlSelectDirtyRepos`.
- `cli/store/repo_query_filter.go` defines `sqlSelectReposByTag` and `sqlSelectUnpushedRepos`.
- Database query errors are wrapped via `fmt.Errorf(constants.ErrDBQuery, err)`, breaking typed error propagation (`*apperror.AppError`).

#### 3.3.2 SQL Constants Consolidation (`cli/constants/constants_store.go`)
Add the following SQL constants to `cli/constants/constants_store.go`:
- `SQLSelectReposByBranch`: Query for selecting repos by branch with case-insensitive collation.
- `SQLSelectDirtyRepos`: Query for selecting dirty repos joined with local tracking.
- `SQLSelectReposByTag`: Query for selecting repos matching metadata notes tag.
- `SQLSelectUnpushedRepos`: Query for selecting repos with unpushed/ahead status.

#### 3.3.3 Dedicated Repository Component (`cli/store/repo_repository.go`)
Implement a dedicated repository struct and interface encapsulating repository database queries:
```go
// RepoRepository encapsulates repository query operations.
type RepoRepository interface {
	FindByBranch(branch string) ([]model.ScanRecord, *apperror.AppError)
	FindDirty() ([]model.ScanRecord, *apperror.AppError)
	FindByTag(tag string) ([]model.ScanRecord, *apperror.AppError)
	FindUnpushed() ([]model.ScanRecord, *apperror.AppError)
	FindPathBySlug(slug string) (string, *apperror.AppError)
}
```
All repository methods must return `*apperror.AppError`, wrapping SQL errors via `apperror.WrapSimple` or `apperror.NewDatabaseError`.

#### 3.3.4 Refactoring `cli/store/repo.go` and `cli/store/repo_query_filter.go`
- Remove local `const sqlSelect...` declarations from both files.
- Re-route queries through `constants.SQLSelect...` or delegate directly to `RepoRepository`.
- Maintain backward-compatible methods on `*DB` while ensuring errors bubble up cleanly.

---

### 3.4 Decoupling SQL & Parsing in Pipeline History Render (`cli/cmdpipeline`)

#### 3.4.1 Problem Statement
In `cli/cmdpipeline/pipeline_history_render.go`:
- `findLocalRepoDirForSlug` calls `sql.Open("sqlite", ...)` directly, manually constructing connection strings and bypassing the store abstraction.
- `queryRepoPathFromConn` executes raw SQL `"SELECT AbsolutePath FROM Repo WHERE Slug = ? OR Slug = ? LIMIT 1"`.
- Errors are swallowed: `if err := row.Scan(&absPath); err == nil`.
- `parseFirstVTag` uses raw newline splitting with index-based checks instead of reusable string helpers.

#### 3.4.2 Architectural Remedy
1. **Decouple Database Access:**
   - Refactor `findLocalRepoDirForSlug` to query via `store.DefaultDB()` or `store.RepoRepository`.
   - Remove direct `sql.Open` SQLite connection instantiation.
   - Return structured error or log diagnostic warnings without silent swallowing.
2. **Refactor Tag Parsing:**
   - Utilize `strutil` line iteration to locate valid version tags.
   - Enforce guard clauses to flatten iteration logic.

---

### 3.5 DateTime Utility Standard (`cli/timeutil`)

#### 3.5.1 Problem Statement
Multiple modules format UTC timestamps with ad-hoc invocations of `time.Now().UTC().Format(time.RFC3339)` or custom date layouts, risking inconsistent millisecond precision or timezone offsets.

#### 3.5.2 DateTime Utility Specification (`cli/timeutil/timeutil.go`)
Create `cli/timeutil/timeutil.go`:
```go
package timeutil

import "time"

// NowUTCISO returns the current UTC time in ISO8601/RFC3339 format.
func NowUTCISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// FormatUTCISO formats a given time in UTC ISO8601/RFC3339 format.
func FormatUTCISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
```

---

### 3.6 Agent Subtask JSON Telemetry (`cli/cmdagent`)

#### 3.6.1 Problem Statement
When autonomous agents execute subtask transitions (`start`, `complete`, `fail`), `cli/cmdagent/agent_subtask.go` prints human-readable strings (e.g. `fmt.Printf("Subtask %s completed by %s\n", subID, opts.AgentRole)`). Automated calling agents require structured JSON output to record execution outcomes and error telemetry reliably.

#### 3.6.2 Telemetry Architecture
- Add `--json` flag to `subtaskStartCmd`, `subtaskCompleteCmd`, and `subtaskFailCmd`.
- When `--json` is enabled, emit a structured JSON payload:
```json
{
  "status": "completed",
  "subtaskId": "subtask-01",
  "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
  "taskCode": "Subtask-01",
  "agentRole": "Worker 01",
  "timestamp": "2026-10-22T12:00:00Z",
  "evidence": "Verification complete",
  "hasError": false
}
```
- Use `timeutil.NowUTCISO()` to timestamp telemetry events.

---

## 4. Subtask Decomposition & Execution Plan

| Subtask ID | Focus Area | Key Files Affected |
| :--- | :--- | :--- |
| `subtask-01` | IDE Tokens, Workdir Constants & Environment Utility | `cli/strutil/strutil.go`<br>`cli/cmdide/ide_add.go`<br>`cli/constants/constants_workdir.go`<br>`cli/envutil/envutil.go`<br>`cli/cmdpe/pe_target.go` |
| `subtask-02` | Store Repository Pattern, SQL Constants & Pipeline Decoupling | `cli/constants/constants_store.go`<br>`cli/store/repo_repository.go`<br>`cli/store/repo.go`<br>`cli/store/repo_query_filter.go`<br>`cli/cmdpipeline/pipeline_history_render.go` |
| `subtask-03` | DateTime Standardization & Agent Subtask JSON Telemetry | `cli/timeutil/timeutil.go`<br>`cli/cmdagent/agent_subtask.go` |

---

## 5. Verification Protocol (Zero-Build Safety)

All changes authored according to this specification must be verified using static code inspection and syntactic analysis without invoking `go test` or `go build`:
1. Verify function length <= 15 lines across all modified functions.
2. Verify zero nested `if` statements using guard inversion.
3. Verify all booleans use `is*` or `has*` positive prefixes.
4. Verify leaf functions return `*apperror.AppError` and never call exit handlers.
5. Verify relative path hygiene across all documentation and comments.
