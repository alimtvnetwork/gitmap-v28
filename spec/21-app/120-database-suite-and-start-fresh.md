# 120 — Database Management Suite & Start-Fresh Architecture

## Overview

**Module Number:** 120  
**Version:** 1.0.0  
**Updated:** 2026-09-03  
**Status:** Production-Ready  
**AI Confidence:** Production-Ready  
**Ambiguity Score:** None  

---

## Purpose

GitMap employs a hybrid database architecture consisting of a centralized **Primary Master SQLite Database** and per-repository **Split Search Databases** (`bin/data/repo_search/*.db`). This specification defines the database inspection, telemetry, maintenance commands (`gitmap db`), the full reset mechanism (`gitmap start-fresh`), and the database path display in the terminal help footer.

---

## Architecture Rationale: Why Split DBs Exist

GitMap separates repository file caches and search indices into per-repository split databases for three foundational reasons:

1. **Lock Elimination & Concurrency:**  
   When scanning dozens of repositories concurrently across multiple worker goroutines, writing file trees to a single SQLite master database creates lock contention on SQLite WAL journals. Split DBs allow each worker to open and write to an isolated database file with zero mutex bottlenecks.
2. **Fault Isolation:**  
   If an interrupted scan or crash corrupts a local repository index, only that repository's individual `.db` file in `repo_search/` is affected. The global master database remains entirely intact.
3. **Optimized Cache Locality:**  
   Fast search queries (`gitmap find-file`, `gitmap grep`) can open targeted split databases in parallel without loading massive unrelated indices into RAM.

---

## Command Signatures & Behaviors

### 1. `gitmap db ls` (or `gitmap db list`)
Displays an architectural breakdown of:
- **Primary Master Database**: Absolute path, file size, status, and responsibilities (repositories, tags, profiles, scans).
- **Split Databases**: Directory location (`repo_search/`), count of per-repo DB files, total disk size, and concurrency rationale.

### 2. `gitmap db repo-db list`
Scans all individual split databases in `bin/data/repo_search/` and renders a table:
- `SLUG`: Repository slug identifier.
- `DB FILE`: Split database file basename.
- `SIZE`: File size on disk.
- `FILES`: Row count in `RepoFile` table.
- `CACHES`: Row count in `SearchCache` table.
- `STATUS`: Tracking state in master database.

### 3. `gitmap db sizes list` (or `gitmap db sizes ls`)
Renders a consolidated disk utilization summary across:
- Master DB (`gitmap.db`, WAL, SHM).
- Split DB aggregate (`repo_search/`).
- Release metadata directory (`.gitmap/release/`).

### 4. `gitmap db reset` (or `gitmap db clear`)
Safely resets cached data while preserving repository configurations:
- Flags: `--all` (resets master repo records and split DBs), `--dry-run`, and `--yes`.
- Interactive confirmation prompt when `--yes` is omitted.

### 5. `gitmap start-fresh`
Full system re-initialization:
- **Irreversible Transaction Warning**: Displays a prominent warning modal in the terminal:
  ```text
  ⚠️ WARNING: IRREVERSIBLE TRANSACTION ⚠️
  This operation will permanently delete all GitMap databases, cached indices,
  and split search databases. Tracked repository paths will need to be re-scanned.
  Are you sure you want to proceed? [y/N]:
  ```
- **Purge Execution**: Deletes `gitmap.db`, `-wal`, `-shm`, and recursively purges `repo_search/`.
- **Reconstruction**: Executes clean schema initialization migrations and outputs verified readiness.

### 6. Terminal Help Footer Integration
The root `gitmap help` footer displays the active Primary Master SQLite database location:
```text
  ────────────────────────────────────────────────────────────
  gitmap binary
  ● Version:     v6.166.0
  ● Database:    D:\wp-work\riseup-asia\gitmap\bin\data\gitmap.db
```

---

## Acceptance Criteria

### Scenario 1: Database Inspection
- **Given** a populated GitMap installation with 10 tracked repositories
- **When** `gitmap db ls` is executed
- **Then** the terminal outputs the absolute path to `gitmap.db`, details the active split databases in `repo_search/`, and prints the architectural rationale block.

### Scenario 2: Irreversible Reset Safeguard
- **Given** an active GitMap database
- **When** `gitmap start-fresh` is invoked without `--yes`
- **Then** the command halts with the irreversible transaction warning, requires explicit `y` input, and only proceeds to wipe and rebuild upon confirmation.

---

## Cross-References

- Split DB Architecture: [`../05-split-db-architecture/00-overview.md`](../05-split-db-architecture/00-overview.md)
- Database ERD: [`gitmap-database-erd.mmd`](./gitmap-database-erd.mmd)
- Database Conventions: [`../04-database-conventions/00-overview.md`](../04-database-conventions/00-overview.md)
