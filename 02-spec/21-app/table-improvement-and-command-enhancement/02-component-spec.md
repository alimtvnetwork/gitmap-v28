# Component Specification: Dual SQLite Databases (Cache & Backup), New Commands Discovery, Help Documentation & LLM Skills Synchronization

**Module:** `cli/cmdpending`, `cli/constants`, `cli/helpdoc`, `.agents/skills/gitmap`  
**Parent Plan:** `.ai-memory/plans/table-improvement-and-command-enhancement.md`  
**Specification File:** `02-spec/21-app/table-improvement-and-command-enhancement/02-component-spec.md`  
**Version:** 1.1.0  
**Status:** Approved for Implementation  

---

## 1. Executive Summary & Component Scope

This component specification establishes the technical requirements, database schemas, interfaces, CLI command catalogs, help documentation, and LLM agent skill integrations for four interconnected subsystems of the `table-improvement-and-command-enhancement` task:

1. **SQLite Status Cache Engine (`cli/cmdpending/pending_commits_cache.go`)**:
   - Ultra-fast, local SQLite-backed status caching for `gitmap pending-commits` (`gitmap pc`).
   - Replaces expensive multi-repo filesystem git subprocess executions (`git status --porcelain`, `git rev-parse`, `git rev-list`) with sub-millisecond cached telemetry lookups.
   - Strictly enforces a **90-second Time-To-Live (TTL)** (within the required 1–2 minute window).
   - Stale records exceeding 90 seconds are immediately purged from SQLite and never trusted.
   - First-class bypass flags: `--no-cache` (skip read and write) and `--refresh` (force immediate re-scan and cache update).

2. **SQLite Status Backup Database Engine (`cli/cmdpending/pending_commits_backup.go`)**:
   - Dedicated persistent historical backup database (`pending_commits_backup.db`).
   - Persists point-in-time snapshot records of repository statuses with full metadata (branch, commit SHA, version, uncommitted files count, unpushed count, dirty state, and complete JSON telemetry payload).
   - Unlike the cache DB, backup records are **not auto-purged on the 90-second TTL**, guaranteeing that historical fleet states can be retrieved and served again.
   - First-class retrieval flag: `--backup` (`gitmap pc --backup`) allows viewing the latest backed-up snapshot when working offline or reviewing historical states.

3. **New Commands Discovery Engine (`cli/cmdpending/new_commands_cmd.go`)**:
   - Dedicated command `gitmap new-commands` (alias: `gitmap nc`) cataloging the last 100 new commands introduced across recent GitMap milestones.
   - Flexible CLI filtering capabilities:
     - `--limit <n>` (alias `-n`, default: 100)
     - `--filter <str>` (aliases `-f` and `-q`) supporting keyword search across name, alias, description, and copy-pasteable example.
     - `--category <cat>` (alias `-c`) supporting category isolation (`commits`, `pipeline`, `fleet`, `ai`, `runner`, `status`, `storage`, `automation`).
     - `--json` (alias `-j`) emitting structured machine-readable payload.
     - `-h` / `--help` displaying comprehensive contextual usage help.
   - Rich terminal table formatting with category badges and copy-pasteable runnable examples.

4. **CLI Constants, Help Catalog & LLM Skills Synchronization**:
   - Centralized constants in `cli/constants/constants_cli.go`.
   - Comprehensive markdown documentation in `cli/helpdoc/pending-commits.md` and `cli/helpdoc/new-commands.md`.
   - Topic summary and alias routing in `cli/helpdoc/catalog.go` and `cli/helpdoc/print.go`.
   - LLM skill integration in `.agents/skills/gitmap/SKILL.md` enabling AI agents to discover dirty repos, parse consolidated `UNCOMMITTED` counts, apply hierarchical tree remedies, utilize cache flags safely, and explore new commands.

---

## 2. Component 1: SQLite Status Cache Engine

### 2.1 Problem & Motivation
In multi-repository workspaces (e.g., 40+ repositories), executing `gitmap pending-commits` triggers multiple subprocess git executions per repository:
- `git status --porcelain` (untracked, modified, staged files)
- `git rev-parse --abbrev-ref HEAD` (current branch)
- `git rev-parse --abbrev-ref @{u}` (upstream branch)
- `git rev-list --count @{u}..HEAD` (unpushed commits count)
- `git rev-list @{u}..HEAD` (unpushed commit hashes)

Across dozens of repositories on Windows or network filesystems, full inspection takes 2–5 seconds. Because developers and autonomous AI agents query workspace status in rapid continuous loops, an ultra-fast, local SQLite cache drastically accelerates responsiveness (<5ms) while guaranteeing freshness through a strict 90-second TTL.

### 2.2 Database File & Storage Location
Following GitMap's standard split-database architecture:
- **Location:** Resolved dynamically via `store.BinaryDataDir()`:
  - Relative logical path: `.gitmap/pending_commits_cache.db` (under user app data directory).
  - Concrete resolution helper: `filepath.Join(store.BinaryDataDir(), "pending_commits_cache.db")`.
- **Driver:** Standard Go SQLite driver (`sql.Open("sqlite", dbPath)` with `modernc.org/sqlite`).
- **Connection Flags:** Enabled `WAL` mode (`PRAGMA journal_mode=WAL;`) and busy timeout (`PRAGMA busy_timeout=5000;`) for non-blocking concurrent reads and resilient single-writer operations.

### 2.3 Table Schema: `pending_commits_cache`

```sql
CREATE TABLE IF NOT EXISTS pending_commits_cache (
    repo_path          TEXT PRIMARY KEY,
    head_sha           TEXT NOT NULL,
    uncommitted_count  INTEGER NOT NULL,
    unpushed_count     INTEGER NOT NULL,
    is_dirty           INTEGER NOT NULL,
    cached_at_unix     INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pending_cache_time 
ON pending_commits_cache (cached_at_unix);
```

#### Field Specifications:
| Column Name | SQLite Type | Go Type | Description & Invariants |
| :--- | :--- | :--- | :--- |
| `repo_path` | `TEXT PRIMARY KEY` | `string` | Normalized clean repository path or slug serving as primary key. |
| `head_sha` | `TEXT NOT NULL` | `string` | Full 40-character commit SHA at repository HEAD (`git rev-parse HEAD`). |
| `uncommitted_count` | `INTEGER NOT NULL` | `int` | Consolidated count of untracked, modified, and staged working files. |
| `unpushed_count` | `INTEGER NOT NULL` | `int` | Total number of local commits ahead of upstream branch. |
| `is_dirty` | `INTEGER NOT NULL` | `int` | Positive boolean flag (1 if `uncommitted_count > 0` or `unpushed_count > 0`, else 0). |
| `cached_at_unix` | `INTEGER NOT NULL` | `int64` | Unix epoch timestamp in seconds (`time.Now().Unix()`) recording entry creation. |

### 2.4 TTL Policy & Strict No-Trust Invariants
- **TTL Duration:** 90 seconds (1.5 minutes, squarely within the required 1–2 minute range).
- **Expiration Threshold:** An entry is considered expired when:
  $$\text{time.Now().Unix()} - \text{cached\_at\_unix} > 90$$
- **Immediate Purge Invariant:** Expired entries are **never trusted** and must be purged immediately:
  1. On individual cache query: If an entry exists for `repo_path` but `nowUnix - cached_at_unix > 90`, the entry is deleted immediately from the database (`DELETE FROM pending_commits_cache WHERE repo_path = ?`) and treated as a cache miss.
  2. On cache initialization: A housekeeping query runs to purge all expired entries in bulk:
     ```sql
     DELETE FROM pending_commits_cache WHERE cached_at_unix < (? - 90);
     ```
- **HEAD Commit Verification (Dual-Gate Validation):** Even if an entry is within the 90-second window, if the repository's current `HEAD` commit SHA differs from the stored `head_sha`, the cache entry is invalidated immediately, deleted, and re-evaluated via fresh git inspection.

### 2.5 Bypass Flags: `--no-cache` and `--refresh`
- **`--no-cache` Flag:**
  - Completely disables cache reading and writing.
  - Forces direct git command execution for all target repositories.
  - Useful for CI/CD gates, release pre-flights, and deterministic tests where cached data is disallowed.
- **`--refresh` Flag:**
  - Bypasses reading existing cache entries.
  - Runs fresh git status inspections across all target repositories.
  - Overwrites existing records in SQLite with updated metrics and a new `cached_at_unix` timestamp (`INSERT ... ON CONFLICT(repo_path) DO UPDATE SET ...`).

### 2.6 Data Structures & API Contract (`cli/cmdpending/pending_commits_cache.go`)

```go
package cmdpending

import (
	"database/sql"
)

// PendingCommitCacheRecord stores cached git status metrics for a repository.
type PendingCommitCacheRecord struct {
	RepoPath         string
	HeadSHA          string
	UncommittedCount int
	UnpushedCount    int
	IsDirty          bool
	CachedAtUnix     int64
}

// DefaultCacheTTLSeconds defines the standard 90-second TTL (1.5 minutes).
const DefaultCacheTTLSeconds int64 = 90

// OpenPendingCommitsCache opens or initializes the cache SQLite database.
func OpenPendingCommitsCache(dbPath string) (*sql.DB, error)

// GetCachedPendingStatus retrieves a valid, non-expired cache record for a repo.
// Returns (nil, false, nil) on cache miss or expired record (expired records are purged).
func GetCachedPendingStatus(conn *sql.DB, repoPath, currentHeadSHA string, nowUnix int64) (*PendingCommitCacheRecord, bool, error)

// SaveCachedPendingStatus writes or updates a repository status record in SQLite.
func SaveCachedPendingStatus(conn *sql.DB, record PendingCommitCacheRecord) error

// PurgeExpiredPendingCache removes all records older than ttlSeconds.
func PurgeExpiredPendingCache(conn *sql.DB, nowUnix, ttlSeconds int64) (int64, error)

// InvalidatePendingCache deletes an individual repository record by path.
func InvalidatePendingCache(conn *sql.DB, repoPath string) error
```

---

## 3. Component 2: SQLite Backup Database Engine

### 3.1 Purpose & Architectural Contrast with Cache DB
The pending commits subsystem operates with a dual-database design:
1. **Primary Cache DB (`pending_commits_cache.db`)**: Short-lived, ephemeral cache with a strict 90-second TTL. Stale records are purged immediately.
2. **Secondary Backup DB (`pending_commits_backup.db`)**: Long-lived historical snapshot store. It captures repository status states across scanning sessions. It is **never automatically purged** on the 90-second TTL. If network filesystems are offline, or if the developer or automated agent needs to review previously captured states, the backup DB serves this historical data via the `--backup` flag.

### 3.2 Database File & Storage Path
- **Location:** Resolved dynamically via `store.BinaryDataDir()`:
  - Relative logical path: `.gitmap/pending_commits_backup.db`.
  - Concrete resolution helper: `filepath.Join(store.BinaryDataDir(), "pending_commits_backup.db")`.
- **Driver:** SQLite (`modernc.org/sqlite`) with `WAL` mode and busy timeout enabled.

### 3.3 Database Schemas: Latest Snapshots & Historical Ledger

```sql
-- Table 1: Latest snapshot per repository (upserted on every fresh scan)
CREATE TABLE IF NOT EXISTS pending_commits_backup (
    repo_path          TEXT PRIMARY KEY,
    branch             TEXT NOT NULL,
    head_sha           TEXT NOT NULL,
    version            TEXT NOT NULL,
    uncommitted_count  INTEGER NOT NULL,
    unpushed_count     INTEGER NOT NULL,
    is_dirty           INTEGER NOT NULL,
    backed_up_at_unix  INTEGER NOT NULL,
    payload_json       TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pending_backup_time 
ON pending_commits_backup (backed_up_at_unix);

-- Table 2: Historical log recording audit snapshots over time
CREATE TABLE IF NOT EXISTS pending_commits_backup_history (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_path          TEXT NOT NULL,
    branch             TEXT NOT NULL,
    head_sha           TEXT NOT NULL,
    version            TEXT NOT NULL,
    uncommitted_count  INTEGER NOT NULL,
    unpushed_count     INTEGER NOT NULL,
    is_dirty           INTEGER NOT NULL,
    backed_up_at_unix  INTEGER NOT NULL,
    payload_json       TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_backup_hist_repo_time 
ON pending_commits_backup_history (repo_path, backed_up_at_unix);
```

#### Field Specifications:
| Column Name | SQLite Type | Go Type | Description & Invariants |
| :--- | :--- | :--- | :--- |
| `repo_path` | `TEXT PRIMARY KEY` / `TEXT` | `string` | Normalized repository path or workspace slug. |
| `branch` | `TEXT NOT NULL` | `string` | Current branch name at time of backup (e.g. `main`). |
| `head_sha` | `TEXT NOT NULL` | `string` | Commit SHA at HEAD when backup was recorded. |
| `version` | `TEXT NOT NULL` | `string` | Git tag version or SemVer string (e.g. `v6.524.0`). |
| `uncommitted_count` | `INTEGER NOT NULL` | `int` | Consolidated uncommitted files count. |
| `unpushed_count` | `INTEGER NOT NULL` | `int` | Local commits ahead of remote tracking branch. |
| `is_dirty` | `INTEGER NOT NULL` | `int` | Affirmative flag (`1` for changes, `0` for clean). |
| `backed_up_at_unix` | `INTEGER NOT NULL` | `int64` | Timestamp recording when the backup was saved. |
| `payload_json` | `TEXT NOT NULL` | `string` | Complete JSON telemetry representing `PendingCommitItem`. |

### 3.4 Backup CLI Flag: `--backup`
When the developer or agent runs:
```bash
gitmap pc --backup
```
The command:
1. Bypasses live filesystem git checks.
2. Reads the latest snapshot from `pending_commits_backup.db`.
3. Renders the table marked with a badge: `[SERVED FROM BACKUP DB: YYYY-MM-DD HH:MM:SS UTC]`.
4. Allows developers to review status when offline, or compare previous state before executing large batch changes.

### 3.5 Data Structures & API Contract (`cli/cmdpending/pending_commits_backup.go`)

```go
package cmdpending

import (
	"database/sql"
	"time"
)

// PendingCommitBackupRecord represents a persistent backup snapshot.
type PendingCommitBackupRecord struct {
	RepoPath         string    `json:"repoPath"`
	Branch           string    `json:"branch"`
	HeadSHA          string    `json:"headSha"`
	Version          string    `json:"version"`
	UncommittedCount int       `json:"uncommittedCount"`
	UnpushedCount    int       `json:"unpushedCount"`
	IsDirty          bool      `json:"isDirty"`
	BackedUpAtUnix   int64     `json:"backedUpAtUnix"`
	PayloadJSON      string    `json:"payloadJson"`
}

// OpenPendingCommitsBackup opens or initializes the backup SQLite database.
func OpenPendingCommitsBackup(dbPath string) (*sql.DB, error)

// SavePendingCommitBackup persists a repository status snapshot to both latest and history tables.
func SavePendingCommitBackup(conn *sql.DB, record PendingCommitBackupRecord) error

// GetLatestPendingBackup retrieves the most recent backup for a specific repository.
func GetLatestPendingBackup(conn *sql.DB, repoPath string) (*PendingCommitBackupRecord, bool, error)

// GetAllLatestPendingBackups retrieves the most recent backup snapshots across all repositories.
func GetAllLatestPendingBackups(conn *sql.DB) ([]PendingCommitBackupRecord, error)

// GetPendingBackupHistory retrieves recent historical snapshots for a repository up to limit.
func GetPendingBackupHistory(conn *sql.DB, repoPath string, limit int) ([]PendingCommitBackupRecord, error)
```

---

## 4. Component 3: New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`)

### 4.1 Overview & Motivation
Over the course of GitMap's recent 100 commits and multiple feature releases, dozens of enterprise commands and subcommands were added to support high-speed scanning, cluster management, CI/CD telemetry, and AI automation. However, discovering these commands is difficult without a dedicated discovery tool.

The `new-commands` command (alias: `nc`) provides a structured, searchable catalog of the last 100 new commands introduced in recent git history. It offers concise syntax descriptions, category tags, flags, and actionable copy-pasteable examples.

### 4.2 CLI Routing & Command Registration
- **Command Name:** `new-commands`
- **Short Alias:** `nc`
- **Registration Point:** Registered in `cli/cmd/roottooling.go` under `toolingDevEntries()`:
  ```go
  cli.Entry{
      Name:        constants.CmdNewCommands,
      Aliases:     []string{constants.CmdNewCommandsAlias},
      Description: constants.HelpNewCommands,
      Handler:     cmdpending.RunNewCommands,
  }
  ```

### 4.3 CLI Flags & Argument Parsing
The discovery engine supports flexible, POSIX-compliant command-line flags:

| Flag Name | Short Flags | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--limit` | `-n` | `int` | `100` | Caps the number of cataloged commands returned. Clamped to min 1. |
| `--filter` | `-f`, `-q` | `string` | `""` | Substring keyword filter matching command name, alias, description, or example. Both `-f` and `-q` must be supported. |
| `--category` | `-c` | `string` | `""` | Category filter (e.g. `commits`, `pipeline`, `fleet`, `ai`, `runner`, `status`, `storage`, `automation`). |
| `--json` | `-j` | `bool` | `false` | Emits structured JSON array conforming to `NewCommandsPayload`. |
| `--help` | `-h` | `bool` | `false` | Displays contextual command usage card. |

### 4.4 Data Structures (`cli/cmdpending/new_commands_types.go`)

```go
package cmdpending

import "time"

// NewCommandEntry represents a cataloged command entry.
type NewCommandEntry struct {
	Name        string   `json:"name"`
	Alias       string   `json:"alias,omitempty"`
	Version     string   `json:"version"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Example     string   `json:"example"`
	Flags       []string `json:"flags,omitempty"`
}

// NewCommandsPayload encapsulates structured output for new commands.
type NewCommandsPayload struct {
	Timestamp        time.Time         `json:"timestamp"`
	TotalCommands    int               `json:"totalCommands"`
	FilteredCommands int               `json:"filteredCommands"`
	Limit            int               `json:"limit"`
	Category         string            `json:"category,omitempty"`
	Filter           string            `json:"filter,omitempty"`
	Commands         []NewCommandEntry `json:"commands"`
}

// NewCommandsOptions encapsulates parsed CLI flags.
type NewCommandsOptions struct {
	Limit    int
	Category string
	Filter   string
	IsJSON   bool
	IsHelp   bool
}

// DefaultNewCommandsOptions returns options initialized with safe defaults.
func DefaultNewCommandsOptions() NewCommandsOptions {
	return NewCommandsOptions{
		Limit:    100,
		Category: "",
		Filter:   "",
		IsJSON:   false,
		IsHelp:   false,
	}
}
```

### 4.5 Core Catalog Architecture & Command Inventory
The catalog engine maintains a chronological list of commands added across recent milestones:

1. **Commit & Workspace Remediation (`commits`):**
   - `cpar`: Commit & push all repositories across workspace (`gitmap cpar "wip: save changes"`).
   - `cpf`: Stage, commit, and push feature with hyphen title (`gitmap cpf "auth - add jwt validation"`).
   - `cpb`: Stage, commit, and push bugfix (`gitmap cpb "cache - fix ttl expiry"`).
   - `cpr`: Stage, commit, and push release milestone (`gitmap cpr "v6.525.0"`).
   - `pc` / `pending-commits`: Discover uncommitted modifications & unpushed commits with tree remediation (`gitmap pc`).

2. **Pipeline Telemetry & Diagnostics (`pipeline`):**
   - `pe`: Pipeline error analyzer with dynamic runner extraction (`gitmap pe -t`).
   - `pe all`: Workspace-wide CI/CD pipeline health monitor (`gitmap pe all`).
   - `pipeline-ai status`: Smart timeout status monitor (`gitmap pipeline-ai status -t 120`).
   - `rerun` / `rr`: Rerun failed CI/CD pipeline steps (`gitmap rerun --step=lint`).
   - `sug`: Shutdown execution until quality gates are green (`gitmap sug`).

3. **Fleet & Cluster Operations (`fleet`):**
   - `commons` / `co`: Apply standardized repository configurations (`gitmap commons`).
   - `space backup-branch`: Create workspace backup branches (`gitmap space backup-branch "feat-login"`).
   - `nodes fs`: Distributed full status across cluster nodes (`gitmap nodes fs`).
   - `nodes pe all`: Distributed pipeline error scan across fleet (`gitmap nodes pe all`).
   - `sync`: High-speed multi-repo prompt, skill, and script sync (`gitmap sync --workers 8`).

4. **AI & Agent Systems (`ai`):**
   - `task init`: Initialize SQLite agent task manager (`gitmap task init --name "audit" --budget 10`).
   - `task claim`: Claim next available atomic subtask (`gitmap task claim --db task.db --agent "Worker 02"`).
   - `task log-action`: Record agent step for crash forensics (`gitmap task log-action --db task.db --subtask-id 2 --action "write_to_file"`).
   - `task complete`: Mark subtask complete with evidence (`gitmap task complete --db task.db --subtask-id 2 --evidence "PASS"`).
   - `llm train`: Generate 4-stage chained curriculum and Antigravity skill (`gitmap llm train`).
   - `new-commands` / `nc`: Discover recent commands added in the last 100 commits (`gitmap nc --limit 20`).

5. **Storage, Secrets & Runners (`storage`, `runner`):**
   - `rs text`: Save sensitive tokens to secure vault (`gitmap rs text "token-val" --slug "gh-pat"`).
   - `rc text`: Save reusable test harness to shared cache (`gitmap rc text "script" --slug "test" --ext .ps1`).
   - `py`: Execute Python script with auto-resolved interpreter (`gitmap py -c "print('hello')"`).
   - `pwsh` / `ps`: Execute PowerShell with deterministic `-NoProfile` (`gitmap ps -c "Get-Location"`).

### 4.6 Terminal Rendering Specification
When outputting to terminal (non-JSON), `gitmap new-commands` renders a structured 80-character boxed presentation:
```text
┌──────────────────────────────────────────────────────────────────────────────┐
│                    GITMAP NEW COMMANDS DISCOVERY (LAST 100)                  │
├──────────────────────────────────────────────────────────────────────────────┤
│ Total: 100 cataloged  │ Displaying: 100 commands  │ Filter: none             │
├──────────────────────────────────────────────────────────────────────────────┤
│ COMMAND / ALIAS           VERSION    CATEGORY     DESCRIPTION                │
├──────────────────────────────────────────────────────────────────────────────┤
│ cpar                      v6.520.0   commits      Commit & push all repos    │
│   Example: gitmap cpar "wip: batch backup all dirty repositories"            │
│                                                                              │
│ pending-commits (pc)      v6.523.0   commits      Status table & tree hints  │
│   Example: gitmap pc --refresh                                               │
│                                                                              │
│ new-commands (nc)         v6.525.0   ai           Discover recent commands   │
│   Example: gitmap nc --filter "cache" --limit 10                             │
│                                                                              │
│ pe all                    v6.518.0   pipeline     Fleet pipeline error check │
│   Example: gitmap pe all --json                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ Filter Tip: Use 'gitmap nc -f <keyword>' or 'gitmap nc -c <category>'        │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Component 4: CLI Constants, Help Catalog & Documentation Synchronization

### 5.1 Centralized Constants (`cli/constants/constants_cli.go`)

```go
package constants

const (
	// CmdPendingCommits discovers uncommitted changes and unpushed commits across repos.
	CmdPendingCommits      = "pending-commits"
	CmdPendingCommitsAlias = "pc"

	// CmdNewCommands catalogs and filters commands added in recent git history.
	CmdNewCommands         = "new-commands"
	CmdNewCommandsAlias    = "nc"

	// HelpPendingCommits describes pending-commits in root help.
	HelpPendingCommits = "  pending-commits (pc) Discover uncommitted changes & unpushed commits with tree remediation"

	// HelpNewCommands describes new-commands in root help.
	HelpNewCommands    = "  new-commands (nc)    Inspect last 100 commands added in recent commits with examples"

	// Flag constants for cache and backup control
	FlagNoCache     = "no-cache"
	FlagDescNoCache = "Bypass status cache and query git status directly from filesystem"

	FlagRefresh     = "refresh"
	FlagDescRefresh = "Force refresh status cache by re-evaluating git status and updating SQLite"

	FlagBackup      = "backup"
	FlagDescBackup  = "Serve repository statuses from persistent backup database without live git scan"

	FlagLimit       = "limit"
	FlagDescLimit   = "Maximum number of items to display"

	FlagFilter      = "filter"
	FlagDescFilter  = "Keyword filter matching command name, alias, description, or example"

	FlagCategory      = "category"
	FlagDescCategory  = "Category filter (e.g. commits, pipeline, fleet, ai, runner)"
)
```

### 5.2 Help Catalog & Aliases Synchronization

#### `cli/helpdoc/catalog.go`:
```go
var topicSummaries = map[string]string{
	// Existing topics...
	"pending-commits": "Discover uncommitted modifications and unpushed commits with tree view remediation and batch fix options.",
	"pc":              "Discover uncommitted modifications and unpushed commits with tree view remediation and batch fix options.",
	"new-commands":    "Catalog and search recent CLI commands added in the last 100 commits with syntax and usage examples.",
	"nc":              "Catalog and search recent CLI commands added in the last 100 commits with syntax and usage examples.",
}
```

#### `cli/helpdoc/print.go`:
```go
var helpAliases = map[string]string{
	// Existing aliases...
	"pc": "pending-commits",
	"nc": "new-commands",
}
```

### 5.3 Documentation Updates: `cli/helpdoc/pending-commits.md`
The help document must document:
1. Unified table layout: `REPOSITORY`, `VER/BRANCH`, `UNCOMMITTED`, `UNPUSHED`, `STATUS`.
2. Hierarchical tree-view remediation (`├── Option 1: ...`, `└── Option 2: ...`).
3. Global batch remediation command in table footer.
4. Cache & backup flags: `--no-cache`, `--refresh`, and `--backup`.
5. Concrete runnable CLI examples.

### 5.4 Documentation Updates: `cli/helpdoc/new-commands.md`
The new help document must document:
1. Command synopsis for `gitmap new-commands` and alias `gitmap nc`.
2. Flag specifications:
   - `--limit <n>` (alias `-n`)
   - `--filter <str>` (aliases `-f` and `-q`)
   - `--category <cat>` (alias `-c`)
   - `--json` (alias `-j`)
   - `-h` / `--help`
3. Description of discovery mechanism analyzing recent git commits.
4. Runnable copy-pasteable examples filtering by search keyword, category, limit cap, and JSON output.

---

## 6. Component 5: LLM Skills Synchronization (`.agents/skills/gitmap/SKILL.md`)

Autonomous AI coding agents rely on `.agents/skills/gitmap/SKILL.md` as their core operational manual. Synchronizing this skill file guarantees that LLMs pair-programming with the user will automatically:
1. Prefer `gitmap pc` over slow multi-repo shell status loops.
2. Accurately understand the new `UNCOMMITTED` consolidated column (`UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`) and `VER/BRANCH` column.
3. Parse the tree-view remediation options and choose the correct surgical fix (`Option 1: commit & push` vs `Option 2: stash`).
4. Utilize `--no-cache` and `--refresh` when verifying git state transitions, and understand the 90-second SQLite status cache TTL.
5. Use `--backup` to inspect historical snapshots when offline.
6. Use `gitmap new-commands` / `gitmap nc` (with `-f` / `-q` and `-c`) to discover newly introduced capabilities and examples.

### 6.1 Replacement Matrix Updates
Add new mandatory replacement rows to Section 2:

| Prohibited Legacy Tool / Pattern | GitMap Mandatory Command | Why GitMap is Required |
| :--- | :--- | :--- |
| ❌ Multi-repo `git status` loops or manual directory scanning | ✅ `gitmap pc` (`gitmap pending-commits`) | High-speed multi-repo status with 90s SQLite cache, consolidated `UNCOMMITTED` metrics, short branch display, tree-view surgical fixes, and batch footer remediation. |
| ❌ Guessing newly added CLI commands from raw git commit logs | ✅ `gitmap nc` (`gitmap new-commands [--limit 100]`) | Native cataloging of recently introduced commands, syntax, flags, and usage examples directly from git history. |

### 6.2 Essential Command Cheat Sheet Updates
Under Section 2 (**Workspace & Multi-Repository Status**):
- `gitmap pc` (alias `gitmap pending-commits`) — Fast multi-repo pending commits table with 90s SQLite cache, consolidated `UNCOMMITTED` column, and tree remediation.
- `gitmap pc --refresh` — Force-refreshes SQLite cache by re-evaluating live git status.
- `gitmap pc --no-cache` — Disables SQLite cache completely for strict pre-flight verification.
- `gitmap pc --backup` — Serves status from persistent backup DB without querying live filesystem.
- `gitmap pc --json` — Emits structured JSON telemetry of uncommitted and unpushed repositories.

Under Section 8 (**Autonomous Agent Onboarding & Curriculum (LLM)**):
- `gitmap new-commands` (alias `gitmap nc`) — Discovers and filters the last 100 commands added across recent git history with runnable examples.
- `gitmap nc --filter "<pattern>"` / `gitmap nc -f "<pattern>"` / `gitmap nc -q "<pattern>"` — Filters new commands by keyword or category.
- `gitmap nc --category "<cat>"` — Filters commands by functional category.

### 6.3 Operational Guardrail Rule 9
> **9. Pending Status Caching, Backup & Remediation Invariant:** When checking multi-repo status, agents should invoke `gitmap pc`. Respect the 90-second SQLite status cache. When validating state immediately after applying code modifications or git operations, pass `gitmap pc --refresh` or `gitmap pc --no-cache` to ensure live filesystem validation. If dirty repositories are reported, use the suggested tree remediation command or batch footer command (`gitmap cpar`). To review historical status snapshots, pass `gitmap pc --backup`. To discover recently added commands and usage examples, use `gitmap nc`.

---

## 7. Acceptance Criteria (AC1 Through AC8)

### AC1: Unified `UNCOMMITTED` Column
- The pending commits table replaces separate `DIRTY`, `UNTRACK`, `MODIF`, and `STAGE` columns with a single `UNCOMMITTED` column.
- The metric equals `UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`.
- Renders clean `-` or `0` for clean repos, and highlighted integer count for dirty repos.

### AC2: Repository Version and Branch (`VER/BRANCH`)
- The table displays `<version>/<branch>` if a SemVer tag exists (e.g. `v6.523.1/main`), or `<branch>` if unversioned.
- Column width is standardized and safely truncated with `…` to preserve 80-character box alignment.

### AC3: Batch Remediation Footer Command
- When dirty repositories exist, the table footer renders:
  `Fleet Remediation: gitmap cpar "wip: save changes"`
- Clean status output omits the fleet remediation line.

### AC4: Hierarchical Tree-View Remediation
- Every dirty repository row renders two indented tree branches:
  - `├── Option 1: git -C "<repo>" add -A && git commit -m "wip" && git push`
  - `└── Option 2: git -C "<repo>" stash -u`
- Clean repositories do not render tree lines.

### AC5: SQLite Status Cache with 90-Second TTL
- Primary cache stored in `store.BinaryDataDir()` / `pending_commits_cache.db`.
- Records younger than 90 seconds with matching `HEAD` SHA return cached counts instantly (<5ms).
- Records older than 90 seconds are deleted immediately on access and never trusted.
- Bypass flags `--no-cache` and `--refresh` work as specified.

### AC6: SQLite Status Backup Database
- Persistent backup DB stored in `store.BinaryDataDir()` / `pending_commits_backup.db`.
- Holds point-in-time snapshots of repository status and JSON payloads.
- Records are NOT auto-purged by the 90-second TTL.
- `--backup` flag allows serving backed-up status snapshots again.

### AC7: New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`)
- Both `gitmap new-commands` and `gitmap nc` invoke the catalog engine.
- Filters the last 100 new commands from recent git history by default (`--limit 100`).
- Supports `--filter` (both `-f` and `-q`), `--category`, and `--json`.
- Displays syntax, version, category, description, and concrete copy-pasteable examples.

### AC8: CLI Help Documentation & LLM Skills Sync
- `cli/constants/constants_cli.go` defines all command, alias, help, and flag constants.
- `cli/helpdoc/pending-commits.md` and `cli/helpdoc/new-commands.md` provide complete user documentation.
- `cli/helpdoc/catalog.go` and `cli/helpdoc/print.go` register help summaries and aliases.
- `.agents/skills/gitmap/SKILL.md` includes `gitmap pc` and `gitmap nc` in replacement matrix, cheat sheet, and operational guardrails.

---

## 8. Implementation Standards & Invariants

1. **Affirmative Booleans:** All boolean fields, variables, parameters, and return values MUST use affirmative naming (`isDirty`, `isClean`, `hasUncommitted`, `hasUpstream`, `isJSON`, `isHelpRequested`). Negative or inverted booleans (`notClean`, `noDirty`, `skipCache`) are strictly forbidden.
2. **Function Sizing Guard (8–15 LOC Cap):** Decompose all logic into small, single-purpose functions.
3. **Structured Error Handling:** All errors must be wrapped with `appfault.AppError` and standard error codes (`E9001`, `E9002`, `E9003`).
4. **Strict Relative Git Paths:** All file paths cited in documentation, error messages, and tests must be relative to repository root. Absolute paths and `file:///` URIs are forbidden.
5. **No Git Commands:** Never execute git commands directly during multi-agent subtask execution.
