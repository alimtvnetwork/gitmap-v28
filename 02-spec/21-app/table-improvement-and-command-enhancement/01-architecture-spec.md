# Architecture Specification: Pending Commits Table Improvement & New Commands Discovery Engine

> **Spec Version:** 1.1.0  
> **Status:** Approved / Grounded  
> **Task Slug:** `table-improvement-and-command-enhancement`  
> **Target Path:** `02-spec/21-app/table-improvement-and-command-enhancement/01-architecture-spec.md`  

---

## 1. System Architecture Context & Problem Statement

### 1.1 Existing Pending Commits UX Shortcomings
GitMap provides the `gitmap pending-commits` (alias: `pc`) command implemented in `cli/cmdpending/pending_commits_cmd.go`. While functional, the current terminal table rendering exhibits several critical architectural and usability deficiencies:

1. **Fragmented Column Allocation & Screen Real Estate Waste:**
   - The current table allocates four separate columns for file status: `DIRTY` (5 chars), `UNTRACK` (7 chars), `MODIF` (6 chars), and `STAGE` (6 chars), totaling 24+ characters of terminal width.
   - For repository health assessment, developers and automated agents primarily need to know whether uncommitted files exist and the aggregate change count. Fragmenting these numbers forces aggressive truncation on repository names (capped at 20 chars) and branch names (capped at 8 chars).

2. **Absence of Repository Version Identification:**
   - The table displays only `BRANCH` (truncated to 8 characters, e.g. `main`, `feat-x…`), giving zero visibility into which release milestone or SemVer tag the local repository is anchored to (e.g., `v6.523.1` or `v6.524.0`).
   - When managing multi-repo fleets or versioned submodules, developers must run secondary commands to identify whether a dirty repository is tracking a released tag or an unversioned feature branch.

3. **Missing In-Place Tree-View Remediation Hints:**
   - When a repository is identified as dirty or pending, developers must manually formulate remediation commands.
   - There are no inline contextual suggestions for quick resolution (e.g., fast commit and push vs. stash or discard).

4. **Lack of Workspace Fleet Batch Guidance:**
   - When 5, 10, or 20 repositories have pending uncommitted changes across a workspace, fixing them repo-by-repo is slow and error-prone.
   - GitMap already possesses the enterprise fleet commit-and-push engine `gitmap cpar` (`commit-push-all-repos`), but `gitmap pending-commits` never surfaces this batch capability in its table footer.

5. **Repeated Subprocess Git Overhead & Latency:**
   - Scanning 40+ repositories spawns multiple subprocess calls (`git status`, `git rev-parse`, `git rev-list`) per repository, causing 2–5 second delays on disk-heavy or Windows environments. Rapid checks require sub-millisecond responses without returning stale data.

6. **Absence of Offline or Historic Status Recovery:**
   - If repositories are temporarily inaccessible or when reviewing previous states during troubleshooting, developers have no mechanism to retrieve previously observed statuses without running a live scan.

### 1.2 The Need for `gitmap new-commands` (`gitmap nc`)
GitMap has grown through dozens of minor and major releases, accumulating over 100 enterprise commands spanning fleet SSH delegation, split-DB indexing, CI/CD diagnostics, automation units, and workspace sync. However:
- Discoverability of recently added commands is low; developers rely on memory or fragmented help text.
- Standard help lists all commands alphabetically or in monolithic clusters without highlighting recent additions.
- Developers require a focused command: `gitmap new-commands` (and short alias `gitmap nc`) that filters the last 100 new commands, provides concise descriptions, version metadata, and practical copy-pasteable CLI examples.

---

## 2. Pending Commits Table Improvement Architecture

### 2.1 Unified `UNCOMMITTED` Column
Instead of splitting file states into `DIRTY`, `UNTRACK`, `MODIF`, and `STAGE`, the table replaces these four columns with a single unified `UNCOMMITTED` column.
- **Metric Calculation:** `TotalUncommitted = UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`.
- **Display Formatting:**
  - When `TotalUncommitted == 0`: Formatted as `-` or `0` in dimmed or green text.
  - When `TotalUncommitted > 0`: Formatted as clean numeric count `N` (colored yellow).
- **Space Savings:** Reclaims 15+ columns of horizontal space to expand repository names and version strings.

### 2.2 Repository Short Version and Branch (`VER/BRANCH`)
The improved table introduces a combined `VER/BRANCH` column:
- **Version Discovery Logic:**
  1. Inspect local repository Git tags via `git -C <dir> describe --tags --abbrev=0`.
  2. Fall back to reading `.gitmap/version.json` or `version.json` if present.
  3. If no tag or version file is found, version is empty string `""`.
- **Formatting Algorithm (`formatShortVersionBranch`):**
  - If version is present: `<version>/<branch>` (e.g. `v6.523/main`, `v1.4.0/develop`).
  - If version is absent: `<branch>` (e.g. `main`, `master`, `feat-login`).
  - If combined string exceeds column width (16–18 characters), apply smart abbreviation keeping the version intact and truncating the branch with `…` (e.g. `v6.523/feat-auth…`).

### 2.3 Hierarchical Tree-View Remediation Hints
Beneath each dirty repository row in the terminal table, GitMap renders a two-line hierarchical tree view showing concrete, actionable remediation options:
- **Branch Symbols:**
  - Line 1: `├── Option 1: <command>`
  - Line 2: `└── Option 2: <command>`
- **Option 1 (Commit & Push Local WIP):**
  - Targets committing and pushing the uncommitted working tree.
  - Command: `git -C "<repoPath>" add -A && git commit -m "wip: save changes" && git push` or `gitmap cpf "wip: save changes"`.
- **Option 2 (Stash / Isolate Local Changes):**
  - Targets stashing uncommitted changes cleanly.
  - Command: `git -C "<repoPath>" stash -u`.
- **Indentation & Alignment:** Tree lines are indented beneath the table border using the table border style (`│   ├── Option 1: ...`), ensuring the table bounding box remains rectangular, visually aligned, and clean. Clean repositories do not render tree-view hints.

### 2.4 Overall Fleet Batch Command in Table Footer
In the summary footer box of the pending commits table, when dirty repositories exist, GitMap prominently features the workspace fleet remediation command:
```
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ Fleet Remediation: gitmap cpar "wip: save changes"                           │
  └──────────────────────────────────────────────────────────────────────────────┘
```
This informs the developer immediately that they can batch-commit and push all pending changes across all dirty repositories in a single atomic operation without manual per-repo intervention. When all repositories are clean, the fleet remediation line is omitted.

### 2.5 Table Layout & Column Alignment Specifications
The bounding box width is standardized at 78 inner characters (80 characters total including border edges `│`), preserving full cross-platform terminal compatibility:

```
  ┌──────────────────────────────────────────────────────────────────────────────┐
  │                    GITMAP PENDING COMMITS SUMMARY                            │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ Scanned: 42 repos    │ Dirty: 2 repos   │ Uncommitted: 15 files │ Unpushed: 1 │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ REPOSITORY             VER/BRANCH         UNCOMMITTED   UNPUSHED   STATUS    │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ gitmap                 v6.523.1/main               12          1   ● PEND    │
  │   ├── Option 1: git -C "gitmap" add -A && git commit -m "wip"                │
  │   └── Option 2: git -C "gitmap" stash -u                                     │
  │ movie-cli-v8           v8.0.0/main                  3          0   ● PEND    │
  │   ├── Option 1: git -C "movie-cli-v8" add -A && git commit -m "wip"          │
  │   └── Option 2: git -C "movie-cli-v8" stash -u                               │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ 40 clean repos         -                            0          0   ○ CLEAN   │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ Fleet Remediation: gitmap cpar "wip: save changes"                           │
  └──────────────────────────────────────────────────────────────────────────────┘
```

Width specifications:
- `REPOSITORY`: 22 characters (`%-22s`)
- `VER/BRANCH`: 16 characters (`%-16s`)
- `UNCOMMITTED`: 12 characters (`%12d`)
- `UNPUSHED`: 10 characters (`%10d`)
- `STATUS`: 9 characters (`%-9s`)

---

## 3. Dual-Database Status Caching & Persistent Backup Architecture

To achieve rapid response times (<5ms) while guaranteeing absolute data integrity and offline recoverability, GitMap adopts a dual-database architectural pattern:

```
                       ┌──────────────────────────────────────────────┐
                       │        CLI Invocation: gitmap pc             │
                       │        Flags: --no-cache, --refresh,         │
                       │               --backup, --json               │
                       └──────────────────────┬───────────────────────┘
                                              │
                      ┌───────────────────────┴───────────────────────┐
                      │                                               │
             IsBackupMode == true                            IsBackupMode == false
                      │                                               │
        ┌─────────────┴─────────────┐                         ┌───────┴───────┐
        │  Read from Backup DB      │                         │ Check Cache   │
        │  pending_commits_backup.db│                         │ DB or Scan    │
        └─────────────┬─────────────┘                         └───────┬───────┘
                      │                                               │
                      │                           ┌───────────────────┴───────────────────┐
                      │                           │                                       │
                      │                   IsNoCache == false                      IsNoCache == true
                      │                           │                                       │
                      │                 Query Cache DB                            Bypass Cache DB
                      │                 pending_commits_cache.db                  Run Live Git Status
                      │                 (90s TTL & SHA Gate)                              │
                      │                           │                                       │
                      │                 ┌─────────┴─────────┐                             │
                      │                 │                   │                             │
                      │             Cache HIT           Cache MISS                        │
                      │          (Age <= 90s &        (Expired >90s,                      │
                      │           SHA match)          SHA mismatch,                       │
                      │                 │             or --refresh)                       │
                      │                 │                   │                             │
                      │            Use cached               │                             │
                      │            metrics                  ├─────────────────────────────┤
                      │                 │                   │ Run Live Git Status Subproc │
                      │                 │                   │                             │
                      │                 │                   │ 1. Purge & Save to Cache DB │
                      │                 │                   │    (strict 90s TTL)         │
                      │                 │                   │ 2. Dual-Write to Backup DB  │
                      │                 │                   │    (persistent snapshots)   │
                      │                 │                   └──────────────┬──────────────┘
                      │                 │                                  │
                      └─────────────────┴─────────────────┬────────────────┘
                                                          │
                                         Render Table or Emit JSON
```

### 3.1 Primary Cache DB: `pending_commits_cache.db`
- **Purpose:** Ultra-fast, transient query cache for multi-repo workspace status.
- **Storage Location:** Dynamic resolution via `filepath.Join(store.BinaryDataDir(), "pending_commits_cache.db")`.
- **Database Engine:** SQLite in `WAL` mode (`PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`).
- **Table Schema:**
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

#### 3.1.1 Strict 90-Second TTL & Zero-Trust Invariant
- **TTL Duration:** 90 seconds (1.5 minutes, strictly adhering to the 1–2 minute requirement).
- **Expiration Threshold:** An entry is expired when $\text{time.Now().Unix()} - \text{cached\_at\_unix} > 90$.
- **Immediate Purge Invariant:** Expired entries are **never trusted** and must be purged immediately:
  1. **Query-Time Invalidation:** If an entry for `repo_path` is found but $\text{nowUnix} - \text{cached\_at\_unix} > 90$, the record is immediately deleted (`DELETE FROM pending_commits_cache WHERE repo_path = ?`) and treated as a cache miss.
  2. **Startup Housekeeping:** On opening the cache database, a housekeeping purge query deletes all stale records:
     ```sql
     DELETE FROM pending_commits_cache WHERE cached_at_unix < (? - 90);
     ```
  3. **Dual-Gate HEAD Commit Verification:** Even if $\text{age} \le 90\text{s}$, if the repository's live `HEAD` commit SHA differs from `head_sha`, the entry is invalidated, purged, and recomputed.

#### 3.1.2 Cache Bypass Flags
- `--no-cache`: Completely disables reading and writing to `pending_commits_cache.db`. Useful for deterministic CI/CD quality gates and release pre-flights.
- `--refresh`: Bypasses reading existing cache records, forces fresh live git evaluation, updates SQLite cache with new timestamp, and writes an archive snapshot to the Backup DB.

---

### 3.2 Secondary Backup DB: `pending_commits_backup.db`
- **Purpose:** Dedicated, durable persistence store that retains historical status backups and snapshots of repositories across scans, allowing status states to be served again via `--backup`.
- **Storage Location:** Dynamic resolution via `filepath.Join(store.BinaryDataDir(), "pending_commits_backup.db")`.
- **Database Engine:** SQLite in `WAL` mode (`PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`).
- **Persistence Policy:** Unlike `pending_commits_cache.db`, the Backup DB **does NOT purge records on 90s TTL**. It preserves historical snapshots with timestamps, branch names, versions, and payload JSON for auditing, offline recovery, and comparison.
- **Table Schema:**
  ```sql
  CREATE TABLE IF NOT EXISTS pending_commits_backup (
      id                 INTEGER PRIMARY KEY AUTOINCREMENT,
      repo_path          TEXT NOT NULL,
      branch             TEXT NOT NULL,
      version            TEXT NOT NULL,
      head_sha           TEXT NOT NULL,
      uncommitted_count  INTEGER NOT NULL,
      unpushed_count     INTEGER NOT NULL,
      is_dirty           INTEGER NOT NULL,
      status_label       TEXT NOT NULL,
      details_json       TEXT NOT NULL,
      backed_up_at_unix  INTEGER NOT NULL
  );

  CREATE INDEX IF NOT EXISTS idx_backup_repo_time 
  ON pending_commits_backup (repo_path, backed_up_at_unix DESC);

  CREATE INDEX IF NOT EXISTS idx_backup_time 
  ON pending_commits_backup (backed_up_at_unix DESC);
  ```

#### 3.2.1 Field Specifications:
| Column Name | SQLite Type | Go Type | Description & Invariants |
| :--- | :--- | :--- | :--- |
| `id` | `INTEGER PRIMARY KEY AUTOINCREMENT` | `int64` | Monotonically increasing unique record ID. |
| `repo_path` | `TEXT NOT NULL` | `string` | Normalized repository directory path or slug. |
| `branch` | `TEXT NOT NULL` | `string` | Git branch name at the time of backup (e.g. `main`). |
| `version` | `TEXT NOT NULL` | `string` | Short tag or version string (e.g. `v6.524.0`). |
| `head_sha` | `TEXT NOT NULL` | `string` | Commit SHA at HEAD at time of backup. |
| `uncommitted_count` | `INTEGER NOT NULL` | `int` | Total uncommitted files (untracked + modified + staged). |
| `unpushed_count` | `INTEGER NOT NULL` | `int` | Commits ahead of upstream remote tracking branch. |
| `is_dirty` | `INTEGER NOT NULL` | `int` | Positive boolean flag (1 if uncommitted > 0 or unpushed > 0, else 0). |
| `status_label` | `TEXT NOT NULL` | `string` | Status label string (`● PEND` or `○ CLEAN`). |
| `details_json` | `TEXT NOT NULL` | `string` | Serialized JSON containing detailed file change lists. |
| `backed_up_at_unix` | `INTEGER NOT NULL` | `int64` | Unix epoch timestamp in seconds recording snapshot creation. |

#### 3.2.2 Serving Backups via `--backup`
When the `--backup` CLI flag is passed to `gitmap pending-commits` (`gitmap pc --backup`):
1. **Subprocess Git Scan Bypassed:** Live git subprocesses are completely bypassed.
2. **Snapshot Retrieval:** The engine queries `pending_commits_backup.db` for the latest recorded backup entry for each repository in the workspace:
   ```sql
   SELECT b.* FROM pending_commits_backup b
   INNER JOIN (
       SELECT repo_path, MAX(backed_up_at_unix) AS max_time
       FROM pending_commits_backup
       GROUP BY repo_path
   ) latest ON b.repo_path = latest.repo_path AND b.backed_up_at_unix = latest.max_time
   ORDER BY b.is_dirty DESC, b.repo_path ASC;
   ```
3. **Snapshot Notice Rendered:** The table header or footer indicates that results represent a restored backup snapshot:
   ```
   │ Source: BACKUP SNAPSHOT (from 2026-10-10 10:25:00 UTC)                      │
   ```
4. **Offline Resilience:** If Git is unavailable, network shares are down, or developer is in air-gapped environment, `--backup` reliably surfaces the most recent repository status snapshot.

#### 3.2.3 Dual-Write Synchronization
Whenever a live scan executes (during normal cache misses or with `--refresh` / `--no-cache`):
1. **Step 1:** Status metrics are calculated from live git inspections.
2. **Step 2:** If caching is active, record is written into `pending_commits_cache.db` with `cached_at_unix = now`.
3. **Step 3:** Simultaneously, a persistent snapshot record is appended into `pending_commits_backup.db` with full metadata and `backed_up_at_unix = now`.
4. **Pruning Safeguard:** To prevent unbounded disk growth over months, `PruneOldBackups` retains the last 10 snapshots per repository or entries from the last 30 days.

---

### 3.3 Data Structures & Go API Contracts

#### Cache Engine API (`cli/cmdpending/pending_commits_cache.go`):
```go
package cmdpending

import (
	"database/sql"
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

// OpenPendingCommitsCache opens or initializes the cache SQLite database in WAL mode.
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

#### Backup Engine API (`cli/cmdpending/pending_commits_backup.go`):
```go
package cmdpending

import (
	"database/sql"
)

// PendingCommitBackupRecord represents a persistent historical status backup snapshot.
type PendingCommitBackupRecord struct {
	ID               int64  `json:"id"`
	RepoPath         string `json:"repoPath"`
	Branch           string `json:"branch"`
	Version          string `json:"version"`
	HeadSHA          string `json:"headSha"`
	UncommittedCount int    `json:"uncommittedCount"`
	UnpushedCount    int    `json:"unpushedCount"`
	IsDirty          bool   `json:"isDirty"`
	StatusLabel      string `json:"statusLabel"`
	DetailsJSON      string `json:"detailsJson"`
	BackedUpAtUnix   int64  `json:"backedUpAtUnix"`
}

// OpenPendingCommitsBackup opens or initializes the backup SQLite database in WAL mode.
func OpenPendingCommitsBackup(dbPath string) (*sql.DB, error)

// SaveBackupStatus saves a repository status record to the backup database.
func SaveBackupStatus(conn *sql.DB, record PendingCommitBackupRecord) error

// GetLatestBackupStatus returns the most recent backup record for a specific repository.
func GetLatestBackupStatus(conn *sql.DB, repoPath string) (*PendingCommitBackupRecord, bool, error)

// GetAllLatestBackups retrieves the most recent backup records across all known repositories.
func GetAllLatestBackups(conn *sql.DB) ([]PendingCommitBackupRecord, error)

// PruneOldBackups prunes historical records per repository beyond maxPerRepo count.
func PruneOldBackups(conn *sql.DB, maxPerRepo int) (int64, error)
```

---

## 4. `gitmap new-commands` (`gitmap nc`) Architectural Design

### 4.1 Overview & Command Routing
The `new-commands` command provides a curated discovery engine for the last 100 newly introduced commands across recent GitMap releases.
- **Primary Command:** `gitmap new-commands`
- **Short Alias:** `gitmap nc`
- **CLI Registration:** Registered in `cli/cmd/roottooling.go` inside `toolingDevEntries()` and constants defined in `cli/constants/constants_cli.go`.

### 4.2 Data Models (`cli/cmdpending/new_commands_types.go`)
```go
package cmdpending

import "time"

// NewCommandEntry represents a single cataloged command entry.
type NewCommandEntry struct {
	Name        string   `json:"name"`
	Alias       string   `json:"alias,omitempty"`
	Version     string   `json:"version"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Example     string   `json:"example"`
	Flags       []string `json:"flags,omitempty"`
}

// NewCommandsPayload defines the top-level structured telemetry for JSON output.
type NewCommandsPayload struct {
	Timestamp      time.Time         `json:"timestamp"`
	TotalCataloged int               `json:"totalCataloged"`
	ReturnedCount  int               `json:"returnedCount"`
	CategoryFilter string            `json:"categoryFilter,omitempty"`
	Commands       []NewCommandEntry `json:"commands"`
}

// NewCommandsOptions encapsulates CLI flags for the new-commands command.
type NewCommandsOptions struct {
	Limit          int
	FilterQuery    string
	CategoryFilter string
	IsJSON         bool
}
```

### 4.3 Command Catalog Engine (`cli/cmdpending/new_commands_cmd.go`)
The catalog maintains a structured, chronologically sorted registry of the last 100 new commands introduced across GitMap milestones:
1. **Fleet & Workspace Batch Operations:**
   - `cpar` (Commit and push all repositories in workspace): `gitmap cpar "wip: save changes"`
   - `commons` / `co` (Apply standardized repo configurations): `gitmap commons`
   - `space` (Workspace backup branches and baselines): `gitmap space backup-branch "feat-login"`
2. **Pipeline Telemetry & Diagnostics:**
   - `pe` (Pipeline error analyzer with dynamic runner extraction): `gitmap pe -t`
   - `sug` (Shutdown execution until quality gates are green): `gitmap sug`
   - `rerun` / `rr` (Rerun failed pipeline jobs): `gitmap rerun --step=lint`
3. **Commit & Pull Modernization:**
   - `pending-commits` / `pc` (Inspect dirty working trees and unpushed commits): `gitmap pc`
   - `cpf` (Atomic commit, push, feature): `gitmap cpf "feat: user authentication"`
   - `cpb` (Atomic commit, push, bugfix): `gitmap cpb "fix: buffer overflow"`
   - `cpr` (Commit, push, release workflow): `gitmap cpr "v6.524.0"`
   - `cin` (Commit in replay): `gitmap cin --source repo-a --dest repo-b`
4. **AI & Automation Systems:**
   - `agy` (Antigravity developer tools and queue): `gitmap agy deploy`
   - `agm` (Antigravity fleet agent manager): `gitmap agm status`
   - `aum` (Automation Unit Manager): `gitmap aum run 06-cicd-local-runner.py`
5. **System & OS Integration:**
   - `os dock` (Cross-platform taskbar/dock positioner): `gitmap os dock bottom`
   - `apps` (Installed application auditor and uninstaller): `gitmap apps list`
   - `fix-link` (Symlink and junction repairs): `gitmap fix-link --repair`

### 4.4 CLI Flags & Filtering Capabilities
- `--limit N` (default: 100): Cap output to the most recent `N` commands.
- `--category <cat>`: Filter by functional category (e.g., `commits`, `pipeline`, `fleet`, `ai`, `os`).
- `--filter <query>` / `-q <query>`: Text substring search matching name, alias, description, or example.
- `--json`: Emit complete structured machine-readable payload.
- `-h, --help`: Display contextual help card with examples and flag guide.

---

## 5. Acceptance Criteria

### AC1: Pending Commits Unified Column
- The table output replaces the four disparate columns (`DIRTY`, `UNTRACK`, `MODIF`, `STAGE`) with a single consolidated `UNCOMMITTED` column.
- The value represents `UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`.
- The summary row accurately aggregates total uncommitted files across all inspected repositories.

### AC2: Version/Branch Identification
- The table displays the combined `VER/BRANCH` column for every repository row.
- Repositories with valid tags display `<tag>/<branch>` (e.g. `v6.523.1/main`).
- Repositories without tags display `<branch>`.
- The column gracefully truncates long names with `…` without breaking box border alignment.

### AC3: Hierarchical Tree-View Remediation & Fleet Footer
- Every dirty repository row displays two hierarchical tree branches immediately below:
  - `├── Option 1: <commit-and-push-command>`
  - `└── Option 2: <stash-or-discard-command>`
- Clean repositories do not render tree-view hints.
- The table footer displays the fleet-wide batch command:
  `Fleet Remediation: gitmap cpar "wip: save changes"`

### AC4: SQLite Status Cache with Strict 90s TTL
- Status cache database initialized at `store.BinaryDataDir()` / `pending_commits_cache.db`.
- Query checks age against 90-second TTL. If age $\le 90\text{s}$ and HEAD SHA matches, cached metrics are returned without executing git subprocesses.
- If age $> 90\text{s}$, record is purged immediately (`DELETE`) and treated as cache miss. Expired cache is never trusted.
- Bypass flags:
  - `--no-cache`: Neither reads nor writes to status cache DB.
  - `--refresh`: Bypasses reading existing cache, evaluates live git status, and updates cache DB.

### AC5: Persistent Backup DB & `--backup` Serving
- Backup database initialized at `store.BinaryDataDir()` / `pending_commits_backup.db`.
- Backup records are NOT auto-purged on 90s TTL; they persist snapshot history across sessions.
- On live scan, status records are dual-written to both `pending_commits_cache.db` (with TTL) and `pending_commits_backup.db` (persistent snapshot).
- Passing `--backup` flag loads the latest backup records from SQLite, bypassing git subprocesses entirely.
- Renders table with backup snapshot indicator.

### AC6: CLI Help Constants & Documentation
- Constants registered in `cli/constants/constants_cli.go`: `HelpPendingCommits`, `HelpNewCommands`, `CmdNewCommands`, `CmdNewCommandsAlias`, `FlagNoCache`, `FlagRefresh`, `FlagBackup`.
- Help documents updated: `cli/helpdoc/pending-commits.md` and authored: `cli/helpdoc/new-commands.md`.
- Help catalog in `cli/helpdoc/catalog.go` and aliases in `cli/helpdoc/print.go` resolve `pc` and `nc`.

### AC7: LLM Skills Synchronization
- `.agents/skills/gitmap/SKILL.md` includes `gitmap pc` and `gitmap nc` in Non-Negotiable Replacement Matrix.
- Essential Command Cheat Sheet lists `gitmap pc`, `gitmap pc --refresh`, `gitmap pc --no-cache`, `gitmap pc --backup`, and `gitmap nc`.
- Operational guardrails specify status caching and tree remediation protocols for AI agents.

### AC8: New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`)
- Both `gitmap new-commands` and `gitmap nc` invoke the command catalog engine.
- Supports filtering the last 100 commands by default.
- Every entry displays command name, alias, version introduced, description, and practical CLI example.
- Supports `--json` flag producing valid `NewCommandsPayload` JSON.
- Supports `--limit`, `--category`, and `--filter` flags.

---

## 6. Coding Guidelines & Implementation Invariants

### 6.1 Affirmative Boolean Naming
All boolean fields, variables, parameters, and return values MUST strictly use affirmative naming with `is*` or `has*` prefixes:
- Permitted: `isDirty`, `isClean`, `hasUncommitted`, `hasUpstream`, `isJSON`, `isHelpRequested`, `hasRemediation`, `hasBackup`, `isNoCache`, `isRefresh`, `isBackupMode`.
- Strictly Forbidden: `notClean`, `noDirty`, `uncommitted`, `disableTree`, `skipFooter`, `noCache` (as raw bool var).

### 6.2 Structured Error Handling
All operations returning errors MUST wrap them with `apperror.AppError` and standard error codes:
- `E9001`: Directory inspection failure.
- `E9002`: Repository target resolution error.
- `E9003`: Invalid command filter or limit argument.
- `E9004`: Cache database open or query failure.
- `E9005`: Backup database open or query failure.
- `E9006`: Snapshot serialization failure.

### 6.3 Function Sizing Guard (8–15 LOC Cap)
Every function MUST strictly adhere to the 8–15 line length cap:
- Decompose rendering logic into discrete helpers (`renderTableHeader`, `renderRepoRow`, `renderTreeHints`, `renderTableFooter`).
- Separate CLI argument parsing, model transformation, cache checking, backup persistence, and output rendering into dedicated micro-functions.

### 6.4 Strict Relative Paths
All file references and error messages MUST use relative paths relative to repository root (`02-spec/...`, `.ai-memory/...`, `cli/...`). Absolute filesystem paths (e.g. `C:\...`) and `file:///` URIs are strictly prohibited.

### 6.5 Total Git Command Ban for Subagents
Autonomous subagents are strictly prohibited from executing any Git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). Only the parent orchestrator executes git operations upon completion of all subtasks.
