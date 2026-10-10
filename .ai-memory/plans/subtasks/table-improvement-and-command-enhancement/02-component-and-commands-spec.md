# Subtask Plan 02: Component Specification, Dual SQLite Schemas, New Commands Discovery & Documentation

## Overview
Author and ground the technical component specifications, API contracts, table rendering rules, backup deserialization protocols, and subtask implementation requirements for the `table-improvement-and-command-enhancement` initiative.

This plan defines:
1. Component interfaces for `cli/cmdpending/pending_commits_cmd.go`, `cli/cmdpending/pending_commits_types.go`, `cli/cmdpending/pending_commits_cache.go`, `cli/cmdpending/pending_commits_backup.go`, and `cli/cmdpending/new_commands_cmd.go`.
2. Table rendering contracts: unified `UNCOMMITTED` column (`TotalUncommitted`), `VER/BRANCH` column, hierarchical tree-view remediation hints formatted without border clipping, and fleet batch footer layout.
3. Backup record deserialization contract: unpacking `payload_json` when serving status snapshots to preserve file change and commit lists without loss.
4. New commands discovery engine (`gitmap new-commands` / `gitmap nc`): cataloging 100 new commands across 10 functional categories (`commits`, `diagnostics`, `scanner`, `fleet`, `ai`, `os`, `spec`, `storage`, `sync`, `tooling`) with filtering flags (`--limit` / `-n`, `--filter` / `-f` / `-q`, `--category` / `-c`, `--json` / `-j`) and runnable copy-pasteable examples for every command.
5. Documentation & LLM skill synchronization in `cli/helpdoc/pending-commits.md`, `cli/helpdoc/new-commands.md`, and `.agents/skills/gitmap/SKILL.md`.

---

## Target Specification Files & Owned Artifacts
- `02-spec/21-app/table-improvement-and-command-enhancement/02-component-spec.md` (Component Specification)
- `.ai-memory/plans/subtasks/table-improvement-and-command-enhancement/02-component-and-commands-spec.md` (Implementation Subtask Plan)

---

## Architectural Pillars & Requirements

### 1. Table Rendering Contracts & Formatting Architecture
- **Unified `UNCOMMITTED` Column**:
  - Replaces four separate columns (`DIRTY`, `UNTRACK`, `MODIF`, `STAGE`).
  - Metric: `TotalUncommitted = UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`.
  - Format: `-` or `0` for clean repos; highlighted integer for dirty repos.
- **Combined `VER/BRANCH` Column**:
  - Format: `<version>/<branch>` if release tag present (e.g. `v6.523.1/main`), otherwise `<branch>`.
  - Width: 16 characters max (`%-16s`), truncated with `…` if longer.
- **Hierarchical Tree View Remediation Without Border Clipping**:
  - Rendered beneath dirty repositories (`IsDirty == true`):
    - `│   ├── Option 1: git -C "<repo>" add -A && git commit -m "wip: save changes" && git push`
    - `│   └── Option 2: git -C "<repo>" stash -u`
  - Strict 80-character bounding box: right-padded to 78 inner characters followed by border `│` to prevent terminal line wraps or clipped borders.
  - Omitted for clean repositories.
- **Fleet Batch Command in Footer**:
  - When dirty repositories exist, render batch box in table footer:
    `│ Fleet Remediation: gitmap cpar "wip: save changes"                           │`
  - Omitted when all repositories are clean.

### 2. Dual-Database Status Architecture (Cache vs. Backup)
- **Primary Status Cache DB (`pending_commits_cache.db`)**:
  - Path: `filepath.Join(store.BinaryDataDir(), "pending_commits_cache.db")`.
  - Schema: `pending_commits_cache` with columns `repo_path`, `head_sha`, `uncommitted_count`, `unpushed_count`, `is_dirty`, `cached_at_unix`.
  - Strict 90-second TTL (1.5 minutes). Records older than 90s are purged immediately upon access and never trusted.
  - Dual-gate validation: Requires both age $\le 90$s and `HEAD` commit SHA match.
  - Bypass flags: `--no-cache` (skip read/write) and `--refresh` (force immediate re-scan and cache update).
- **Secondary Status Backup DB (`pending_commits_backup.db`)**:
  - Path: `filepath.Join(store.BinaryDataDir(), "pending_commits_backup.db")`.
  - Schema: `pending_commits_backup` with columns `repo_path`, `repo_name`, `current_branch`, `version`, `short_version_branch`, `head_sha`, `uncommitted_count`, `unpushed_count`, `is_dirty`, `backed_up_at_unix`, `payload_json`.
  - Independent of 90s TTL (not auto-purged). Persists snapshots across sessions for offline viewing and troubleshooting.
  - Retrieval flag: `--backup` (`gitmap pc --backup`) serves the latest backed-up snapshot without filesystem git scans.
- **Backup Record Deserialization Contract**:
  - When serving backups via `--backup`, unpack `payload_json` into `RepoPendingCommitRecord`.
  - Restores `PendingFiles`, `UnpushedCommitSHAs`, `UntrackedFilesCount`, `ModifiedFilesCount`, `StagedFilesCount`, and `RemediationOptions`.
  - Graceful fallback to column metrics if `payload_json` is empty or malformed.

### 3. New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`)
- Catalogs 100 new commands introduced across recent GitMap milestones and commits.
- Structured into 10 categories (10 commands each):
  1. `commits` (10 commands): `cpar`, `cpf`, `cpb`, `cpr`, `cin`, `pending-commits` (`pc`), `backup-branch`, `fix all`, `pcp`, `commons`.
  2. `diagnostics` (10 commands): `pe`, `pe all`, `te all`, `sug`, `rerun`, `regoldens`, `doctor`, `audit-legacy`, `fastgate`, `pipeline-ai`.
  3. `scanner` (10 commands): `scan export`, `scan merge`, `rescan`, `rescan-subtree`, `clone-only-missing`, `clone-sync`, `repo-create`, `dedupe`, `size`, `orphans`.
  4. `fleet` (10 commands): `ssh-bind`, `ssh-deploy`, `nodes`, `cluster`, `sc`, `deploy-keys`, `ssh-join`, `cluster-run-script`, `cluster-bootstrap`, `ssh-scan`.
  5. `ai` (10 commands): `agent task`, `agy`, `agm`, `aum`, `ai-analysis`, `macro`, `llm-docs`, `llm train`, `prompt show --copy`, `ai-fix`.
  6. `os` (10 commands): `os dock`, `os panel`, `apps`, `fix-link`, `vpm`, `which-format`, `os-dns`, `os-theme`, `power`, `autologin`.
  7. `spec` (10 commands): `spec issue`, `spec-author`, `spec-audit`, `user`, `changelog`, `release-notes`, `seowrite`, `help-json`, `prompt-template`, `prompt list`.
  8. `storage` (10 commands): `rs file`, `rs folder`, `rs text`, `rc file`, `rc folder`, `rc text`, `storage`, `clean-dev`, `clear-terminal`, `purge-cache`.
  9. `sync` (10 commands): `sync`, `sync --repo`, `sync --dry-run`, `sync --list`, `pull-all-efficient`, `reconcile`, `has-any-updates`, `latest-branch`, `desktop-sync`, `watch`.
  10. `tooling` (10 commands): `new-commands` (`nc`), `which-format`, `find-files`, `find-files-any`, `find-files-startswith`, `find-files-endswith`, `list-files`, `replace`, `replace-regex`, `cat`.
- Flags:
  - `--limit <n>` (alias `-n`, default: 100, clamped min 1).
  - `--filter <str>`: Text keyword filter matching name, alias, description, or example. **MUST support both `-f` and `-q` flags**.
  - `--category <cat>` (alias `-c`): Category filter.
  - `--json` (alias `-j`): Output machine-readable JSON array conforming to `NewCommandsPayload`.
  - `-h`, `--help`: Contextual help display.
- Every entry includes command name, alias, introduced version, category, concise description, and a copy-pasteable runnable example.
- Terminal rendering formatted with 80-character box borders and category styling.

### 4. CLI Constants & Help Documentation Updates
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
- Update `cli/helpdoc/new-commands.md` detailing:
  - Command overview, aliases, flags (`-f` and `-q`, `-c`, `-n`, `-j`), and usage examples.
- Register help topics and aliases in `cli/helpdoc/catalog.go` (`topicSummaries`) and `cli/helpdoc/print.go` (`helpAliases`).

### 5. LLM Skills Integration (`.agents/skills/gitmap/SKILL.md`)
- Update Non-Negotiable Command Replacement Matrix:
  - Replace manual `git status` loops with `gitmap pc`.
  - Replace guessing new commands with `gitmap nc`.
- Update Essential Command Cheat Sheet:
  - Add `gitmap pc` with flags `--refresh`, `--no-cache`, `--backup`, `--json`.
  - Add `gitmap nc` with flags `--filter`, `-f`, `-q`, `--category`, `--limit`, `--json`.
- Add Operational Guardrail Rule 9 instructing LLM agents on:
  - Parsing the unified `UNCOMMITTED` metric (`UntrackedFilesCount + ModifiedFilesCount + StagedFilesCount`).
  - Applying hierarchical tree-view remedies (`Option 1` vs `Option 2`) or fleet footer fix (`gitmap cpar`).
  - Respecting the 90-second SQLite status cache TTL and using `--refresh` / `--no-cache` after modifications.
  - Using `--backup` for offline status reviews.
  - Using `gitmap nc` to discover recent tools with copy-pasteable examples.

---

## Detailed Implementation Breakdown for Dependent Tasks

### Task-03 Implementation Responsibilities:
1. `cli/cmdpending/pending_commits_backup.go`:
   - Implement `OpenPendingCommitsBackup(dbPath string) (*sql.DB, error)`.
   - Implement `SaveBackupPendingStatus(conn *sql.DB, record RepoPendingCommitRecord, headSHA string, nowUnix int64) error`.
   - Implement `GetBackupPendingStatus(conn *sql.DB, repoPath string) (*PendingCommitBackupRecord, bool, error)`.
   - Implement `ListAllBackupPendingStatuses(conn *sql.DB) ([]PendingCommitBackupRecord, error)`.
   - Implement `RestoreBackupToCache(backupConn, cacheConn *sql.DB, repoPath string, nowUnix int64) error`.
   - Implement `PurgeBackupOlderThan(conn *sql.DB, cutoffUnix int64) (int64, error)`.
   - Enforce the backup record deserialization contract unpacking `payload_json`.
2. Unit tests in `cli/cmdpending/pending_commits_backup_test.go` verifying snapshot insertion, retrieval, and payload unpacking.
3. Verify TTL expiry, immediate deletion, and dual-gate SHA validation in `cli/cmdpending/pending_commits_cache.go` and `cli/cmdpending/pending_commits_cache_test.go`.

### Task-04 Implementation Responsibilities:
1. `cli/cmdpending/pending_commits_cmd.go` & `pending_commits_types.go`:
   - Replace old columns (`DIRTY`, `UNTRACK`, `MODIF`, `STAGE`) with unified `UNCOMMITTED` column.
   - Add `VER/BRANCH` column displaying `<version>/<branch>` or `<branch>` with smart truncation.
   - Render indented tree-view remediation hints under dirty repositories without border clipping (exact 78 inner chars padding).
   - Render fleet batch remediation command (`gitmap cpar "wip: save changes"`) in table footer when dirty repos exist.
   - Wire `--no-cache`, `--refresh`, and `--backup` flags.
   - Unpack `payload_json` in `convertBackupToPendingRecords` to preserve file and commit lists.
2. `cli/cmdpending/new_commands_cmd.go`:
   - Ensure `isFilterFlag` accepts both `-f` and `-q` flags in addition to `--filter`.
   - Verify terminal boxed rendering and JSON output for all 100 commands across 10 categories.
3. `cli/constants/constants_cli.go`:
   - Centralize all commands, aliases, flags, and help text constants.
4. `cli/helpdoc/pending-commits.md` & `cli/helpdoc/new-commands.md`:
   - Full documentation, flags table, and copy-pasteable examples.
5. `cli/helpdoc/catalog.go` and `cli/helpdoc/print.go`:
   - Register help summaries and alias mappings.
6. `.agents/skills/gitmap/SKILL.md`:
   - Replacement matrix, cheat sheet, and operational guardrail rule 9.

---

## Acceptance Criteria Checklist
- [x] Component interfaces specified for `pending_commits_cmd.go`, `pending_commits_types.go`, `pending_commits_cache.go`, `pending_commits_backup.go`, and `new_commands_cmd.go`.
- [x] Table rendering contracts, tree view formatting without border clipping (80-char box layout), and footer layout specified.
- [x] Backup record deserialization contract specified for unpacking `payload_json` without data loss.
- [x] Full specification for `gitmap new-commands` (`gitmap nc`) cataloging 100 new commands across 10 functional categories with copy-pasteable examples for every command.
- [x] CLI flags fully specified: `--limit` / `-n`, `--filter` (both `-f` and `-q`), `--category` / `-c`, `--json` / `-j`, `-h` / `--help`.
- [x] SQLite cache DB schema (`pending_commits_cache.db`) and Backup DB schema (`pending_commits_backup.db`) specified with exact DDL.
- [x] Documentation updates for `cli/helpdoc/pending-commits.md` and `cli/helpdoc/new-commands.md` specified.
- [x] LLM skills integration in `.agents/skills/gitmap/SKILL.md` specified for dirty repo discovery, uncommitted parsing, tree remedies, and new command discovery.
- [x] Strict adherence to relative Git paths, lowercase filenames, and zero git commands executed.
