# Component Specification: SQLite Status Cache, Help Documentation & LLM Skills Synchronization

**Module:** `cli/cmdpending`, `cli/constants`, `cli/helpdoc`, `.agents/skills/gitmap`  
**Parent Plan:** `.ai-memory/plans/table-improvement-and-command-enhancement.md`  
**Specification File:** `02-spec/21-app/table-improvement-and-command-enhancement/02-component-spec.md`  
**Version:** 1.0.0  
**Status:** Approved for Implementation  

---

## 1. Executive Summary & Component Scope

This component specification defines the architectural requirements, database schemas, interfaces, CLI help catalogs, and LLM agent skill integrations for two critical pillars of the `table-improvement-and-command-enhancement` task:

1. **SQLite Status Cache Engine (`cli/cmdpending/pending_commits_cache.go`)**:
   - High-performance, SQLite-backed status caching for `gitmap pending-commits` (`gitmap pc`).
   - Replaces repetitive, resource-intensive repository disk scans (`git status --porcelain` and `git rev-list`) with sub-millisecond cached telemetry lookups.
   - Strictly enforces a **90-second Time-To-Live (TTL)** (within the required 1–2 minute window).
   - Stale records past the TTL are immediately purged and never trusted.
   - First-class bypass flags: `--no-cache` (skip read and write) and `--refresh` (force immediate re-scan and cache update).

2. **CLI Constants, Help Documentation & LLM Skills Synchronization**:
   - Missing constants in `cli/constants/constants_cli.go` (`HelpPendingCommits`, `HelpNewCommands`, `CmdNewCommands`, `CmdNewCommandsAlias`, `FlagNoCache`, `FlagRefresh`).
   - Overhauled command documentation in `cli/helpdoc/pending-commits.md` and new comprehensive documentation in `cli/helpdoc/new-commands.md`.
   - Unified catalog registration in `cli/helpdoc/catalog.go` (`topicSummaries`) and alias routing in `cli/helpdoc/print.go` (`helpAliases`).
   - LLM skill updates in `.agents/skills/gitmap/SKILL.md` ensuring autonomous AI coding agents can discover dirty repositories, leverage tree-view remediation commands, and explore recent commands added across git history.

---

## 2. Component 1: SQLite Status Cache Engine

### 2.1 Problem & Motivation
In large multi-repository workspaces (e.g., 40+ repositories), running `gitmap pending-commits` triggers multiple subprocess git executions per repository:
- `git status --porcelain` (untracked, modified, staged files)
- `git rev-parse --abbrev-ref HEAD` (current branch)
- `git rev-parse --abbrev-ref @{u}` (upstream branch)
- `git rev-list --count @{u}..HEAD` (unpushed commits count)
- `git rev-list @{u}..HEAD` (unpushed commit hashes)

Across dozens of repositories on Windows or network filesystems, this inspection takes 2–5 seconds. Because developers and autonomous agents frequently check repository status in rapid feedback loops, an ultra-fast, local SQLite cache drastically improves responsiveness (<5ms) while guaranteeing freshness through a strict TTL.

### 2.2 Database File & Storage Location
Following GitMap's standard split-database architecture:
- **Location:** Resolved dynamically via `store.BinaryDataDir()`:
  - Relative logical path: `.gitmap/pending_commits_cache.db` (under user app data directory).
  - Concrete resolution helper: `filepath.Join(store.BinaryDataDir(), "pending_commits_cache.db")`.
- **Driver:** Standard Go SQLite driver (`sql.Open("sqlite", dbPath)` or `sqlite3`).
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
	"time"
)

// PendingCommitCacheRecord represents a single cached status row.
type PendingCommitCacheRecord struct {
	RepoPath         string
	HeadSHA          string
	UncommittedCount int
	UnpushedCount    int
	IsDirty          bool
	CachedAtUnix     int64
}

// DefaultCacheTTLSeconds defines the standard 90-second TTL.
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

### 2.7 Execution Pipeline Integration (`cli/cmdpending/pending_commits_cmd.go`)

```
                      [ CLI Invocation: gitmap pc ]
                                    │
                        Parse Flags & Options
                       (--no-cache, --refresh)
                                    │
               ┌────────────────────┴────────────────────┐
               │                                         │
       opts.IsNoCache == true                  opts.IsNoCache == false
               │                                         │
        Bypass DB Cache                     Open SQLite Status Cache
        Run Git Status Directly             Run Housekeeping Purge (now - 90s)
               │                                         │
               │                            For each workspace repo:
               │                            Check GetCachedPendingStatus
               │                                         │
               │                       ┌─────────────────┴─────────────────┐
               │                       │                                   │
               │                   Cache HIT                           Cache MISS
               │               (SHA match & age <= 90s)           (Expired, SHA diff,
               │                       │                          or opts.IsRefresh)
               │                 Load metrics                              │
               │                 from SQLite                       Run Git Status
               │                       │                                   │
               │                       │                        Save to SQLite
               │                       │                        (head_sha, counts,
               │                       │                         cached_at_unix)
               │                       │                                   │
               └───────────────────────┬───────────────────────────────────┘
                                       │
                      Build Consolidated Table Model
                        (UNCOMMITTED, Tree Remediation)
                                       │
                      Render Terminal UI or Output JSON
```

---

## 3. Component 2: CLI Constants & Help Catalog Synchronization

### 3.1 Constants Architecture (`cli/constants/constants_cli.go`)
The CLI command and flag vocabulary must be centralized to eliminate magic strings and ensure compile-time type safety across packages.

#### 1. Command & Alias Constants:
```go
const (
	// CmdPendingCommits discovers uncommitted changes and unpushed commits across repos.
	CmdPendingCommits      = "pending-commits"
	CmdPendingCommitsAlias = "pc"

	// CmdNewCommands catalogs and filters commands added in recent git history.
	CmdNewCommands         = "new-commands"
	CmdNewCommandsAlias    = "nc"
)
```

#### 2. Root Help Line Constants:
```go
const (
	HelpPendingCommits = "  pending-commits (pc) Discover uncommitted changes & unpushed commits with tree remediation"
	HelpNewCommands    = "  new-commands (nc)    Inspect last 100 commands added in recent commits with examples"
)
```

#### 3. Flag Constants:
```go
const (
	FlagNoCache     = "no-cache"
	FlagDescNoCache = "Bypass status cache and query git status directly from filesystem"

	FlagRefresh     = "refresh"
	FlagDescRefresh = "Force refresh status cache by re-evaluating git status and updating SQLite"

	FlagLimit       = "limit"
	FlagDescLimit   = "Maximum number of items to display"
)
```

### 3.2 Help Catalog Metadata (`cli/helpdoc/catalog.go`)
Add explicit topic summaries to `topicSummaries` map:
```go
var topicSummaries = map[string]string{
	// Existing topics...
	"pending-commits": "Discover uncommitted modifications and unpushed commits with tree view remediation and batch fix options.",
	"pc":              "Discover uncommitted modifications and unpushed commits with tree view remediation and batch fix options.",
	"new-commands":    "Catalog and search recent CLI commands added in the last 100 commits with syntax and usage examples.",
	"nc":              "Catalog and search recent CLI commands added in the last 100 commits with syntax and usage examples.",
}
```

### 3.3 Help Aliases Routing (`cli/helpdoc/print.go`)
Register aliases in `helpAliases` map:
```go
var helpAliases = map[string]string{
	// Existing aliases...
	"pc": "pending-commits",
	"nc": "new-commands",
}
```

---

## 4. Component 3: Markdown Help Documentation Updates

### 4.1 Updating `cli/helpdoc/pending-commits.md`
The help document must reflect:
1. The overhauled table layout where `DIRTY`, `UNTRACK`, `MODIF`, and `STAGE` are collapsed into a single `UNCOMMITTED` column.
2. Short branch names.
3. Tree-view remediation display (`├── Option 1: ...`, `└── Option 2: ...`).
4. Global batch fix footer command.
5. Cache control flags: `--no-cache` and `--refresh`.

```markdown
# gitmap pending-commits

Scans local workspace repositories or fleet cluster nodes over SSH to identify uncommitted working tree modifications and unpushed commits ahead of origin.

## Alias

pc

## Usage

```bash
gitmap pending-commits [repoName|all] [flags]
gitmap pc [repoName|all] [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--sort <mode>` | `-s` | `priority` | Sort order: `priority`, `name`, or `count` |
| `--detail <mode>` | `-d` | `summary` | Detail view: `summary`, `all`, `1`, or `none` |
| `--no-cache` | | `false` | Bypass SQLite status cache and query git status directly |
| `--refresh` | | `false` | Force refresh SQLite status cache with fresh repository metrics |
| `--ssh` | `-ssh` | `false` | Aggregate pending commits across all registered cluster SSH nodes |
| `--json` | `-j` | `false` | Output machine-readable JSON telemetry |
| `--dirty-only` | | `true` | Display only repositories with changes |
| `--all` | | `false` | Include completely clean repositories in output |
| `-h`, `--help` | | `false` | Display command help menu |

## Table Columns & Tree View Output

The output renders a clean table consolidating untracked, modified, and staged files:
- **REPOSITORY**: Repository slug or folder name.
- **BRANCH**: Active local branch (short version).
- **UNCOMMITTED**: Consolidated count of all uncommitted working files.
- **UNPUSHED**: Number of local commits ahead of remote tracking branch.
- **STATUS**: Clean or Dirty indicator.

For dirty repositories, a structured tree view provides instant surgical fix commands:
```text
├── Option 1: gitmap cpar "wip: save changes"
└── Option 2: gitmap stash
```

A global batch remediation command is rendered in the footer to resolve all dirty repositories simultaneously:
```text
Batch fix command for all dirty repositories:
  gitmap cpar "wip: batch backup all dirty repositories"
```

## Examples

### Example 1: Standard status overview (uses 90s SQLite cache)
```bash
gitmap pc
```

### Example 2: Force refresh cache with live git scan
```bash
gitmap pc --refresh
```

### Example 3: Completely bypass cache
```bash
gitmap pc --no-cache
```

### Example 4: Output machine-readable JSON
```bash
gitmap pc --json
```
```

### 4.2 Authoring `cli/helpdoc/new-commands.md`
New comprehensive help document for `gitmap new-commands`:

```markdown
# gitmap new-commands

Discovers, catalogs, and filters new CLI commands and subcommands added across the last 100 commits in GitMap's git history, complete with flags and usage examples.

## Alias

nc

## Usage

```bash
gitmap new-commands [flags]
gitmap nc [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit <n>` | `-n` | `100` | Maximum number of commands to inspect from git history |
| `--filter <str>` | `-f` | `""` | Filter commands by keyword or prefix |
| `--category <cat>` | `-c` | `""` | Filter by category (e.g. `workflow`, `cache`, `status`, `cluster`) |
| `--json` | `-j` | `false` | Output machine-readable JSON array of discovered commands |
| `-h`, `--help` | | `false` | Display command help menu |

## Description

The `new-commands` command analyzes commit messages, CLI constant registrations, and command router changes across recent history. It produces a formatted catalog showing:
- Command name and aliases.
- Concise summary of purpose.
- Available flags.
- Concrete, runnable examples.

This command empowers both human engineers and AI coding agents to discover recently introduced capabilities without parsing git logs manually.

## Examples

### Example 1: View recently added commands (up to 100)
```bash
gitmap nc
```

### Example 2: Filter recent commands by search term
```bash
gitmap nc --filter "cache"
```

### Example 3: Limit output to top 10 new commands
```bash
gitmap nc --limit 10
```

### Example 4: Machine-readable JSON output for automated agents
```bash
gitmap nc --json
```
```

---

## 5. Component 4: LLM Skills Synchronization (`.agents/skills/gitmap/SKILL.md`)

Autonomous AI coding agents rely on `.agents/skills/gitmap/SKILL.md` as their core operational manual. Synchronizing this skill file guarantees that LLMs pair-programming with the user will automatically:
1. Prefer `gitmap pc` over slow multi-repo shell status loops.
2. Accurately understand the new `UNCOMMITTED` consolidated column and short branch format.
3. Parse the tree-view remediation options and choose the correct surgical fix.
4. Utilize `--no-cache` and `--refresh` when verifying git state transitions.
5. Use `gitmap new-commands` / `gitmap nc` to discover newly introduced capabilities.

### 5.1 Updates to Non-Negotiable Command Replacement Matrix
Add new mandatory replacement rows to the table in Section 2:

| Prohibited Legacy Tool / Pattern | GitMap Mandatory Command | Why GitMap is Required |
| :--- | :--- | :--- |
| ❌ Multi-repo `git status` loops or manual directory scanning | ✅ `gitmap pc` (`gitmap pending-commits`) | High-speed multi-repo status with 90s SQLite cache, consolidated `UNCOMMITTED` metrics, short branch display, tree-view surgical fixes, and batch footer remediation. |
| ❌ Guessing newly added CLI commands from raw git commit logs | ✅ `gitmap nc` (`gitmap new-commands [--limit 100]`) | Native cataloging of recently introduced commands, syntax, flags, and usage examples directly from git history. |

### 5.2 Updates to Essential Command Cheat Sheet
Under Section 2 (**Workspace & Multi-Repository Status**):
- `gitmap pc` (alias `gitmap pending-commits`) — Fast multi-repo pending commits table with 90s SQLite cache, consolidated `UNCOMMITTED` column, and tree remediation.
- `gitmap pc --refresh` — Force-refreshes SQLite cache by re-evaluating live git status.
- `gitmap pc --no-cache` — Disables SQLite cache completely for strict pre-flight verification.
- `gitmap pc --json` — Emits structured JSON telemetry of uncommitted and unpushed repositories.

Under Section 8 (**Autonomous Agent Onboarding & Curriculum (LLM)**):
- `gitmap new-commands` (alias `gitmap nc`) — Discovers and filters the last 100 commands added across recent git history with runnable examples.
- `gitmap nc --filter "<pattern>"` — Filters new commands by keyword or category.

### 5.3 Updates to Operational Guardrails & Non-Negotiable Invariants
Add operational rule 9:
> **9. Pending Status Caching & Remediation Invariant:** When checking multi-repo status, agents should invoke `gitmap pc`. Respect the 90-second SQLite status cache. When validating state immediately after applying code modifications or git operations, pass `gitmap pc --refresh` or `gitmap pc --no-cache` to ensure live filesystem validation. If dirty repositories are reported, use the suggested tree remediation command or batch footer command.

---

## 6. Acceptance Criteria (AC3, AC5, AC6, AC7)

### 6.1 AC3: Batch Remediation Footer Command
- **Criterion AC3.1:** When one or more dirty repositories exist, `gitmap pc` terminal output MUST render a dedicated footer section following the table.
- **Criterion AC3.2:** The footer MUST present a valid, copy-pasteable batch fix command (e.g. `gitmap cpar "wip: batch backup all dirty repositories"` or equivalent batch remediation helper) that addresses all dirty repositories in a single invocation.
- **Criterion AC3.3:** When all repositories are completely clean, the footer batch remediation suggestion MUST NOT be displayed.

### 6.2 AC5: SQLite Status Cache with 90-Second TTL
- **Criterion AC5.1 (Schema & DB Location):** SQLite database initialized at `store.BinaryDataDir()` / `pending_commits_cache.db` with table `pending_commits_cache` having columns: `repo_path`, `head_sha`, `uncommitted_count`, `unpushed_count`, `is_dirty`, `cached_at_unix`.
- **Criterion AC5.2 (Cache Hit):** When querying a repository whose cache entry is younger than 90 seconds and whose current `HEAD` commit SHA matches `head_sha`, `gitmap pc` MUST return cached counts without executing git status subprocesses.
- **Criterion AC5.3 (Strict TTL & Immediate Purge):** If an entry's age exceeds 90 seconds ($\text{now} - \text{cached\_at\_unix} > 90$), the record MUST be immediately deleted from the database and treated as a cache miss. Expired data must never be returned.
- **Criterion AC5.4 (SHA Invalidation):** If a repository's `HEAD` SHA has changed even within 90 seconds, the entry MUST be invalidated, refreshed, and updated in SQLite.
- **Criterion AC5.5 (Bypass Flags):**
  - `--no-cache`: Neither reads nor writes to SQLite cache.
  - `--refresh`: Bypasses reading existing cache, executes fresh git inspection, and updates SQLite.

### 6.3 AC6: CLI Help Constants & Documentation
- **Criterion AC6.1 (Constants):** `cli/constants/constants_cli.go` defines and exports `HelpPendingCommits`, `CmdNewCommands`, `CmdNewCommandsAlias`, `HelpNewCommands`, `FlagNoCache`, `FlagRefresh`, `FlagDescNoCache`, `FlagDescRefresh`.
- **Criterion AC6.2 (Helpdoc):** `cli/helpdoc/pending-commits.md` contains updated table documentation (`UNCOMMITTED`), tree-view options, batch footer fix, and cache flags.
- **Criterion AC6.3 (New Commands Doc):** `cli/helpdoc/new-commands.md` exists and documents synopsis, flags (`--limit`, `--filter`, `--category`, `--json`), and examples.
- **Criterion AC6.4 (Help Catalog & Aliases):**
  - `cli/helpdoc/catalog.go` registers summaries for `"pending-commits"`, `"pc"`, `"new-commands"`, `"nc"`.
  - `cli/helpdoc/print.go` registers alias mappings `"pc": "pending-commits"` and `"nc": "new-commands"`.
  - Calling `gitmap help pc` or `gitmap help nc` resolves and renders the corresponding markdown documentation.

### 6.4 AC7: LLM Skills Synchronization
- **Criterion AC7.1 (Matrix):** `.agents/skills/gitmap/SKILL.md` includes `gitmap pc` and `gitmap nc` in the Non-Negotiable Command Replacement Matrix.
- **Criterion AC7.2 (Cheat Sheet):** `.agents/skills/gitmap/SKILL.md` includes syntax and flags for `gitmap pc` and `gitmap nc` in the Essential Command Cheat Sheet.
- **Criterion AC7.3 (Guardrails):** `.agents/skills/gitmap/SKILL.md` includes explicit operational rules governing cache freshness, `--refresh` usage, and tree-view remediation.

---

## 7. Verification Strategy & Implementation Sequence

1. **Subtask 02 Execution:**
   - Implement `cli/cmdpending/pending_commits_cache.go`.
   - Write comprehensive unit tests in `cli/cmdpending/pending_commits_cache_test.go` verifying TTL expiry, immediate deletion, SHA matching, `--no-cache`, and `--refresh`.
   - Wire caching into `cli/cmdpending/pending_commits_cmd.go`.

2. **Subtask 04 Execution:**
   - Update `cli/constants/constants_cli.go`.
   - Update `cli/helpdoc/pending-commits.md` and author `cli/helpdoc/new-commands.md`.
   - Update `cli/helpdoc/catalog.go` and `cli/helpdoc/print.go`.
   - Update `.agents/skills/gitmap/SKILL.md`.
   - Add unit test coverage in `cli/helpdoc/` confirming topic resolution and alias mappings.
