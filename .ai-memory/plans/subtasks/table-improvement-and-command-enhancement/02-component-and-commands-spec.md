# Subtask Plan 02: Component Specification, Dual SQLite Schemas, New Commands Discovery & Documentation

## Overview
Author and ground the technical component specifications and subtask implementation requirements for the `table-improvement-and-command-enhancement` initiative. This defines the dual SQLite database architecture (ephemeral 90s status cache vs. persistent status backup), the `gitmap new-commands` (`gitmap nc`) discovery engine filtering the last 100 commands from recent git history with runnable examples, full CLI flag specifications (`--limit`, `--filter` with both `-f` and `-q`, `--category`, `--json`, `-h/--help`), comprehensive CLI help documentation, and autonomous LLM skill synchronization in `.agents/skills/gitmap/SKILL.md`.

---

## Target Specification Files & Owned Artifacts
- `02-spec/21-app/table-improvement-and-command-enhancement/02-component-spec.md` (Component Specification)
- `.ai-memory/plans/subtasks/table-improvement-and-command-enhancement/02-component-and-commands-spec.md` (Implementation Subtask Plan)

---

## Architectural Pillars & Requirements

### 1. Dual-Database Status Architecture (Cache vs. Backup)
- **Primary Status Cache DB (`pending_commits_cache.db`)**:
  - Path: `store.BinaryDataDir()` / `pending_commits_cache.db`.
  - Schema: `pending_commits_cache` with columns `repo_path`, `head_sha`, `uncommitted_count`, `unpushed_count`, `is_dirty`, `cached_at_unix`.
  - Strict 90-second TTL (1.5 minutes). Records older than 90s are immediately purged on access and never trusted.
  - Bypass flags: `--no-cache` (skip read/write) and `--refresh` (force immediate re-scan and cache update).
- **Secondary Status Backup DB (`pending_commits_backup.db`)**:
  - Path: `store.BinaryDataDir()` / `pending_commits_backup.db`.
  - Schemas: `pending_commits_backup` (latest snapshot per repo) and `pending_commits_backup_history` (audit ledger).
  - Columns: `repo_path`, `branch`, `head_sha`, `version`, `uncommitted_count`, `unpushed_count`, `is_dirty`, `backed_up_at_unix`, `payload_json`.
  - Independent of 90s TTL (not auto-purged). Persists snapshots so they can be served again.
  - Retrieval flag: `--backup` (`gitmap pc --backup`) serves the latest backed-up snapshot when offline or reviewing previous state.

### 2. New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`)
- Catalogs the last 100 new commands introduced across recent GitMap milestones and commits.
- Flags:
  - `--limit <n>` (alias `-n`, default: 100, clamped).
  - `--filter <str>`: Text keyword filter matching name, alias, description, or example. **MUST support both `-f` and `-q` flags**.
  - `--category <cat>` (alias `-c`): Case-insensitive category filter (`commits`, `pipeline`, `fleet`, `ai`, `runner`, `status`, `storage`, `automation`).
  - `--json` (alias `-j`): Output machine-readable JSON array conforming to `NewCommandsPayload`.
  - `-h`, `--help`: Contextual help display.
- Every entry includes command name, alias, introduced version, category, concise description, and a copy-pasteable runnable example.
- Terminal rendering formatted with 80-character box borders and category styling.

### 3. CLI Constants & Help Documentation Updates
- Centralize all commands, aliases, flags, and help texts in `cli/constants/constants_cli.go`:
  - `CmdPendingCommits`, `CmdPendingCommitsAlias`, `HelpPendingCommits`.
  - `CmdNewCommands`, `CmdNewCommandsAlias`, `HelpNewCommands`.
  - `FlagNoCache`, `FlagDescNoCache`, `FlagRefresh`, `FlagDescRefresh`, `FlagBackup`, `FlagDescBackup`.
  - `FlagLimit`, `FlagDescLimit`, `FlagFilter`, `FlagDescFilter`, `FlagCategory`, `FlagDescCategory`.
- Update `cli/helpdoc/pending-commits.md` detailing:
  - Unified table columns: `REPOSITORY`, `VER/BRANCH`, `UNCOMMITTED`, `UNPUSHED`, `STATUS`.
  - Indented tree-view surgical fix options (`├── Option 1: ...`, `└── Option 2: ...`).
  - Fleet batch remediation command in footer (`gitmap cpar "wip: save changes"`).
  - Cache and backup flags (`--no-cache`, `--refresh`, `--backup`).
- Author/Update `cli/helpdoc/new-commands.md` detailing:
  - Command overview, aliases, flags (`-f` and `-q`, `-c`, `-n`, `-j`), and usage examples.
- Register help topics and aliases in `cli/helpdoc/catalog.go` (`topicSummaries`) and `cli/helpdoc/print.go` (`helpAliases`).

### 4. LLM Skills Integration (`.agents/skills/gitmap/SKILL.md`)
- Update Non-Negotiable Command Replacement Matrix:
  - Replace manual `git status` loops with `gitmap pc`.
  - Replace guessing new commands with `gitmap nc`.
- Update Essential Command Cheat Sheet:
  - Add `gitmap pc` with flags `--refresh`, `--no-cache`, `--backup`, `--json`.
  - Add `gitmap nc` with flags `--filter`, `-f`, `-q`, `--category`, `--limit`.
- Add Operational Guardrail Rule 9 instructing LLM agents on:
  - Parsing the unified `UNCOMMITTED` metric (`UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`).
  - Applying hierarchical tree-view remedies (`Option 1` vs `Option 2`) or fleet footer fix (`gitmap cpar`).
  - Respecting the 90-second SQLite status cache TTL and using `--refresh` / `--no-cache` after modifications.
  - Using `gitmap nc` to discover recent tools with copy-pasteable examples.

---

## Detailed Implementation Breakdown for Dependent Tasks

### Task-03 Implementation Responsibilities (Worker 01):
1. Create `cli/cmdpending/pending_commits_backup.go`:
   - Implement `OpenPendingCommitsBackup(dbPath string) (*sql.DB, error)`.
   - Implement `SavePendingCommitBackup(conn *sql.DB, record PendingCommitBackupRecord) error`.
   - Implement `GetLatestPendingBackup(conn *sql.DB, repoPath string) (*PendingCommitBackupRecord, bool, error)`.
   - Implement `GetAllLatestPendingBackups(conn *sql.DB) ([]PendingCommitBackupRecord, error)`.
   - Implement `GetPendingBackupHistory(conn *sql.DB, repoPath string, limit int) ([]PendingCommitBackupRecord, error)`.
2. Author unit tests in `cli/cmdpending/pending_commits_backup_test.go` verifying snapshot insertion, retrieval, and history logging.
3. Verify TTL expiry, immediate deletion, and dual-gate SHA validation in `cli/cmdpending/pending_commits_cache.go` and `cli/cmdpending/pending_commits_cache_test.go`.

### Task-04 Implementation Responsibilities (Worker 02):
1. Update `cli/cmdpending/pending_commits_cmd.go` & `pending_commits_types.go`:
   - Replace old columns (`DIRTY`, `UNTRACK`, `MODIF`, `STAGE`) with unified `UNCOMMITTED` column.
   - Add `VER/BRANCH` column displaying `<version>/<branch>` or `<branch>`.
   - Render indented tree-view remediation hints under dirty repositories.
   - Render fleet batch remediation command (`gitmap cpar "wip: save changes"`) in table footer when dirty repos exist.
   - Wire `--no-cache`, `--refresh`, and `--backup` flags.
2. Update `cli/cmdpending/new_commands_cmd.go`:
   - Update `isFilterFlag` to accept both `-f` and `-q` flags in addition to `--filter`.
   - Verify terminal boxed rendering and JSON output.
3. Update `cli/constants/constants_cli.go`:
   - Centralize all commands, aliases, flags, and help text constants.
4. Update `cli/helpdoc/pending-commits.md` & `cli/helpdoc/new-commands.md`:
   - Full documentation, flags table, and copy-pasteable examples.
5. Register help catalog in `cli/helpdoc/catalog.go` and `cli/helpdoc/print.go`.
6. Update `.agents/skills/gitmap/SKILL.md`:
   - Replacement matrix, cheat sheet, and operational guardrail rule 9.

---

## Acceptance Criteria Checklist
- [x] Full specification for `gitmap new-commands` (`gitmap nc`) filtering the last 100 commands from recent git history with examples.
- [x] CLI flags fully specified: `--limit`, `--filter` (both `-f` and `-q`), `--category`, `--json`, `-h/--help`.
- [x] SQLite cache DB schema (`pending_commits_cache.db`) and Backup DB schemas (`pending_commits_backup.db`) specified with exact DDL.
- [x] Documentation updates for `cli/helpdoc/pending-commits.md` and `cli/helpdoc/new-commands.md` specified.
- [x] LLM skills integration in `.agents/skills/gitmap/SKILL.md` specified for dirty repo discovery, uncommitted parsing, tree remedies, and new command discovery.
- [x] Strict adherence to relative Git paths, lowercase filenames, and zero git commands executed.
