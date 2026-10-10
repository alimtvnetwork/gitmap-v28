# Subtask 01 — Architecture Spec and Backup DB Spec

> **Subtask Code:** `Task-01`  
> **Subtask Title:** Architecture Spec and Backup DB Spec  
> **Assigned Agent Role:** Worker 01 (Authoring Subagent)  
> **Parent Plan:** `.ai-memory/plans/table-improvement-and-command-enhancement.md`  
> **Target Task DB:** `.ai-memory/temp-agents/84-table-improvement-and-command-enhancement/agent-task.db`  
> **Status:** IN_PROGRESS  

---

## 1. Objective & Scope

Author and ground the foundational architectural specification and subtask execution plan for the `table-improvement-and-command-enhancement` initiative:
1. **Pending Commits Table Overhaul:**
   - Eliminate fragmented columns (`DIRTY`, `UNTRACK`, `MODIF`, `STAGE`) in favor of a single unified `UNCOMMITTED` column that strictly avoids double-counting `MM` or `AM` files.
   - Introduce combined short version and branch (`VER/BRANCH`) column with dynamic tag inspection.
   - Introduce two-line hierarchical tree-view remediation hints (`├── Option 1: ...`, `└── Option 2: ...`) beneath dirty repositories with untruncated executable commands.
   - Feature workspace fleet batch fix command in table footer (`gitmap cpar "wip: save changes"`).
2. **Dual-Database Status Caching & Persistent Backup Engine:**
   - **Cache DB (`pending_commits_cache.db`):** Transient high-performance cache enforcing strict 90-second TTL (1.5 min), purged immediately on expiration and never trusted if stale.
   - **Backup DB (`pending_commits_backup.db`):** Dedicated persistence engine storing historical status snapshots, allowing previous states to be served again via `--backup`.
3. **Downstream Implementation Contracts:**
   - Define exact Go structs, method signatures, database schemas, and SQLite pragmas for downstream implementation in `Task-03`.

---

## 2. Disjoint File Ownership & Strict Boundaries

To prevent file collision and git index locks across concurrent subagents:
- **Owned Files for Worker 01 (Task-01):**
  1. `02-spec/21-app/table-improvement-and-command-enhancement/01-architecture-spec.md`
  2. `.ai-memory/plans/subtasks/table-improvement-and-command-enhancement/01-table-and-backup-db-spec.md`
- **Strictly Prohibited Actions:**
  - **TOTAL BAN ON GIT COMMANDS:** Subagent MUST NOT execute ANY git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). Only the parent orchestrator executes git commits upon completion.
  - **No Absolute Paths:** All paths in documentation and code must be strictly relative to repository root (`02-spec/...`, `.ai-memory/...`, `cli/...`).
  - **No Build / No Test in Spec Phase:** Spec authoring requires zero build or test execution.

---

## 3. Technical Contracts & Architectural Specifications

### 3.1 Table Model Contract (`cli/cmdpending/pending_commits_types.go`)

```go
package cmdpending

// PendingRepoRow represents a single repository entry in the enhanced table.
type PendingRepoRow struct {
	RepoName         string   // Repository directory name or slug (e.g. "gitmap")
	RepoPath         string   // Relative repository path
	VersionBranch    string   // Formatted short version/branch (e.g. "v6.524.0/main")
	UncommittedCount int      // Untracked + Modified + Staged total
	UnpushedCount    int      // Commits ahead of upstream remote tracking branch
	IsDirty          bool     // True if UncommittedCount > 0 or UnpushedCount > 0
	StatusLabel      string   // "● PEND" or "○ CLEAN"
	Option1Command   string   // Fast commit & push: git -C "<path>" add -A && git commit -m "wip" && git push
	Option2Command   string   // Fast stash: git -C "<path>" stash -u
}

// PendingTableSummary aggregates totals for table header and footer.
type PendingTableSummary struct {
	TotalScanned     int
	TotalDirty       int
	TotalClean       int
	TotalUncommitted int
	TotalUnpushed    int
	HasDirtyRepos    bool
	FleetFixCommand  string // "gitmap cpar \"wip: save changes\""
}
```

### 3.2 Cache Database Contract (`cli/cmdpending/pending_commits_cache.go`)

- **Database File:** Resolved via `filepath.Join(store.BinaryDataDir(), "pending_commits_cache.db")`.
- **Driver & Pragmas:** `PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`.
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
  CREATE INDEX IF NOT EXISTS idx_pending_cache_time ON pending_commits_cache (cached_at_unix);
  ```
- **Go Function Contracts:**
  - `OpenPendingCommitsCache(dbPath string) (*sql.DB, error)`
  - `GetCachedPendingStatus(conn *sql.DB, repoPath, currentHeadSHA string, nowUnix int64) (*PendingCommitCacheRecord, bool, error)`
    - *Invariant:* Returns `(nil, false, nil)` if missing, if `nowUnix - cached_at_unix > 90`, or if `head_sha != currentHeadSHA`. Expired records are immediately deleted (`DELETE FROM pending_commits_cache WHERE repo_path = ?`).
  - `SaveCachedPendingStatus(conn *sql.DB, record PendingCommitCacheRecord) error`
    - *Invariant:* Upserts via `INSERT ... ON CONFLICT(repo_path) DO UPDATE SET ...`.
  - `PurgeExpiredPendingCache(conn *sql.DB, nowUnix, ttlSeconds int64) (int64, error)`
    - *Invariant:* Deletes entries where `cached_at_unix < (nowUnix - ttlSeconds)`.
  - `InvalidatePendingCache(conn *sql.DB, repoPath string) error`

### 3.3 Backup Database Contract (`cli/cmdpending/pending_commits_backup.go`)

- **Database File:** Resolved via `filepath.Join(store.BinaryDataDir(), "pending_commits_backup.db")`.
- **Driver & Pragmas:** `PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`.
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
  CREATE INDEX IF NOT EXISTS idx_backup_repo_time ON pending_commits_backup (repo_path, backed_up_at_unix DESC);
  CREATE INDEX IF NOT EXISTS idx_backup_time ON pending_commits_backup (backed_up_at_unix DESC);
  ```
- **Go Function Contracts:**
  - `OpenPendingCommitsBackup(dbPath string) (*sql.DB, error)`
  - `SaveBackupStatus(conn *sql.DB, record PendingCommitBackupRecord) error`
    - *Invariant:* Appends durable snapshot record on every live git inspection.
  - `GetLatestBackupStatus(conn *sql.DB, repoPath string) (*PendingCommitBackupRecord, bool, error)`
    - *Invariant:* Retrieves the newest historical snapshot for a single repository.
  - `GetAllLatestBackups(conn *sql.DB) ([]PendingCommitBackupRecord, error)`
    - *Invariant:* Retrieves the latest snapshot for every known repository, ordered by dirty status and name.
  - `PruneOldBackups(conn *sql.DB, maxPerRepo int) (int64, error)`
    - *Invariant:* Caps history to avoid unbounded storage consumption.

### 3.4 CLI Flag Routing & State Machine

```text
Invocation: gitmap pc [flags]
├── --backup:
│   └── Query GetAllLatestBackups(backupConn) -> Render Table with [BACKUP SNAPSHOT] indicator.
│       Live git commands are completely bypassed.
├── --no-cache:
│   └── Bypass Cache DB -> Run live git scan -> Dual-write to Backup DB -> Render Table.
├── --refresh:
│   └── Bypass Cache DB read -> Run live git scan -> Upsert Cache DB + Append Backup DB -> Render Table.
└── Default:
    └── Check Cache DB -> If HIT (<90s, SHA match) -> Use cached metrics.
        If MISS (>90s or SHA diff) -> Run live git scan -> Upsert Cache DB + Append Backup DB -> Render Table.
```

---

## 4. Execution Steps for Task-01

- [x] **Step 1:** Claim subtask `Task-01` in `.ai-memory/temp-agents/84-table-improvement-and-command-enhancement/agent-task.db`.
- [x] **Step 2:** Log action in Task DB for architecture spec authoring.
- [x] **Step 3:** Author and enhance `02-spec/21-app/table-improvement-and-command-enhancement/01-architecture-spec.md` with full table improvement, MM file deduplication, untruncated tree remediation, fleet footer, dual-database caching and persistent backup specifications.
- [x] **Step 4:** Author `.ai-memory/plans/subtasks/table-improvement-and-command-enhancement/01-table-and-backup-db-spec.md` detailing contracts, schemas, workflows, and acceptance criteria.
- [ ] **Step 5:** Complete subtask `Task-01` in Task DB with verification evidence.
- [ ] **Step 6:** Dispatch completed JSON report to parent agent.

---

## 5. Downstream Execution Blueprint for Task-03 (Worker 01)

When assigned to `Task-03` ("Backup DB Engine and Dual-DB Cache Implementation"):
1. Create `cli/cmdpending/pending_commits_cache.go` implementing `OpenPendingCommitsCache`, `GetCachedPendingStatus`, `SaveCachedPendingStatus`, `PurgeExpiredPendingCache`, `InvalidatePendingCache`.
2. Create `cli/cmdpending/pending_commits_cache_test.go` verifying 90-second TTL expiration, immediate deletion on read, HEAD SHA change invalidation, and `--no-cache` / `--refresh` behavior.
3. Create `cli/cmdpending/pending_commits_backup.go` implementing `OpenPendingCommitsBackup`, `SaveBackupStatus`, `GetLatestBackupStatus`, `GetAllLatestBackups`, `PruneOldBackups`.
4. Create `cli/cmdpending/pending_commits_backup_test.go` verifying snapshot insertion, retrieving latest backups per repo, `--backup` data fidelity, and historical pruning.
5. Ensure 100% compliance with 8–15 LOC cap per function, affirmative boolean naming (`is*`, `has*`), structured `apperror.AppError` error codes, and strict relative paths.

---

## 6. Acceptance Criteria Matrix

| Criterion | Description | Verification Method |
|---|---|---|
| **AC1** | Consolidated `UNCOMMITTED` column replacing untracked, modified, staged with zero MM double-counting | Architecture spec section 2.1 |
| **AC2** | Combined `VER/BRANCH` column showing tag and branch | Architecture spec section 2.2 |
| **AC3** | Hierarchical tree-view remediation hints beneath dirty repositories with untruncated commands | Architecture spec section 2.3 |
| **AC4** | Fleet batch fix footer command `gitmap cpar "wip: save changes"` | Architecture spec section 2.4 |
| **AC5** | SQLite Cache DB with 90s TTL, auto-purging expired entries, zero trust | Architecture spec section 3.1 |
| **AC6** | SQLite Backup DB storing durable snapshots, serving via `--backup` | Architecture spec section 3.2 |
| **AC7** | Affirmative boolean naming across all struct fields and signatures | Architecture spec section 6.1 |
| **AC8** | Strict relative paths and zero Git command execution | Architecture spec sections 6.4, 6.5 |

---

## 7. Verification Evidence & Completion Sign-Off

- **Spec Document 1:** `02-spec/21-app/table-improvement-and-command-enhancement/01-architecture-spec.md` (Version 1.1.0, 360+ lines, complete dual-db & table architecture).
- **Spec Document 2:** `.ai-memory/plans/subtasks/table-improvement-and-command-enhancement/01-table-and-backup-db-spec.md` (Comprehensive subtask plan & implementation contracts).
- **Git Hygiene:** Zero git commands executed; all paths relative.
