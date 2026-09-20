# Plan 44: PR Commit Engines, SQLite Split-DB Architecture, and Auto-Merge PR Release Suite

> **Execution Lifecycle & Header Summary:**
> - **Task Inception:** Initiated under the `[V2] Parent Task N-Step Continuous Loop & Multi-Agent Orchestration` workflow (`N=500` self-loop budget) based on **Spec 129** (`02-spec/21-app/129-pr-commit-engines-and-sqlite-split-db.md`), Spec 106, Spec 114, Spec 124, and user directives.
> - **Scope & Purpose:** Introduces the first-class `pr` (and `pull-request`) commit replay command family, standardizes the canonical SQLite Split-DB directory formula (`.gitmap/data/<section>/<slug>/sql.db`) across all GitMap subsystems with auto-migration, establishes the dedicated `prdb` storage engine, provides automated merge detection with feature branch forking (`feature/<slug>` or `pr/<id>`), rich Markdown PR description generation, mainline merge (`--no-ff`), rich `pterm` terminal UI, visual ASCII execution graphs, deterministic final snapshot synchronization, and merged PR branch pruning (`pr-clean` / `pr-rm`).
> - **Total Execution Steps / Loops:** 7 discrete subtasks executed across parallel sub-agents (Subtask 01: PR Command Family and Routing, Subtask 02: Split-DB Architecture Engine, Subtask 03: PR SQLite Storage Engine, Subtask 04: PR Merge Detection & Feature Branch Engine, Subtask 05: Terminal UI & Execution Graph Engine, Subtask 06: Final Snapshot Sync & PR Clean Engine, Subtask 07: Help Text, HG Help Groups & JSON Examples).
> - **Final Status:** 100% COMPLETED (All deliverables implemented, guideline verified, ERD parity synced, and consolidated).

---

## 1. Executive Summary & Full Architectural Scope

Plan 44 delivers a major evolution of GitMap's commit transfer and repository management capabilities:

```text
                                  ┌────────────────────────────────────────────────────────┐
                                  │            GitMap PR & Commit Replay Engine            │
                                  └──────────────────────────┬─────────────────────────────┘
                                                             │
         ┌───────────────────────────────────────────────────┼───────────────────────────────────────────────────┐
         ▼                                                   ▼                                                   ▼
   ┌───────────┐                                       ┌───────────┐                                       ┌───────────┐
   │ Subtask 1 │                                       │ Subtask 2 │                                       │ Subtask 3 │
   │PR Routing │                                       │ Split-DB  │                                       │   prdb    │
   ├───────────┤                                       ├───────────┤                                       ├───────────┤
   │dispatch   │                                       │split_db_pa│                                       │types.go   │
   │constants  │                                       │automation │                                       │schema.go  │
   │root.go    │                                       │pipeline   │                                       │ERD parity │
   │rootrelease│                                       │migration  │                                       │WAL pragmas│
   └───────────┘                                       └───────────┘                                       └───────────┘
         │                                                   │                                                   │
         ├───────────────────────────────────────────────────┼───────────────────────────────────────────────────┤
         ▼                                                   ▼                                                   ▼
   ┌───────────┐                                       ┌───────────┐                                       ┌───────────┐
   │ Subtask 4 │                                       │ Subtask 5 │                                       │ Subtask 6 │
   │Merge / PR │                                       │Terminal UI│                                       │Final Sync │
   ├───────────┤                                       ├───────────┤                                       ├───────────┤
   │pr_engine  │                                       │log.go     │                                       │replay.go  │
   │prdesc     │                                       │graph/     │                                       │pr_clean   │
   │replay.go  │                                       │renderer.go│                                       │pipeline.go│
   └───────────┘                                       └───────────┘                                       └───────────┘
                                                             │
                                                             ▼
                                                       ┌───────────┐
                                                       │ Subtask 7 │
                                                       │Help & HG  │
                                                       ├───────────┤
                                                       │rootusage  │
                                                       │pr.md      │
                                                       │JSON schema│
                                                       └───────────┘
```

---

## 2. Deliverables & Subtask Ledger

### Subtask 01: PR Command Family & Conflict-Free Dispatch Routing
- Defined `CmdPR = "pr"`, `CmdPullRequest = "pull-request"`, `CmdPRIn = "pr-in"`, `CmdPRClean = "pr-clean"`, `CmdPRRm = "pr-rm"`, and `CmdPRList = "pr-list"` in `cli/constants/constants_committransfer.go`.
- Dispatched `dispatchCommitTransfer` prior to `dispatchRelease` in `cli/cmd/root.go`, eliminating route collision with `pull-release`.
- Preserved all 5 canonical aliases in `cli/cmd/rootrelease.go` (`pull-release`, `pr`, `release-pull`, `relp`, `rlp`), passing `TestPullReleaseAliasesAreUnique` and `TestPullReleaseAliasesShareOneDispatchEntry`.
- Wired directional subcommands (`pr left`, `pr right`, `pr in`, `pr clean`, `pr list`) in `cli/cmd/dispatchcommittransfer.go` using decomposed <=11 line functions.

### Subtask 02: SQLite Split-DB Canonical Directory Structure Standardization
- Created `cli/store/split_db_path.go` implementing `ResolveSplitDbPath(section, slug, repoRoot)` and `ResolveSplitDbDir(section, slug, repoRoot)`.
- Standardized the database path formula across the entire repository to `.gitmap/data/<section>/<slug>/sql.db` (and `<BinaryDataDir>/<section>/<slug>/sql.db` fallback).
- Defined section constants: `SectionAutomation`, `SectionPipeline`, `SectionPR`, `SectionLogs`, `SectionInstallation`, `SectionStartup`, `SectionSites`, `SectionSchedule`, `SectionRepoSearch`, and `DbFileName = "sql.db"`.
- Added transparent auto-migration moving legacy files (`.gitmap/data/<slug>/automation/sql.db`, `pipeline.db`, etc.) to `sql.db` in canonical directories.
- Refactored `cmdautomation/runtime_db.go`, `cmdautomation/exclusions.go`, `pipelinedb/pipeline_split_db.go`, `cmdpipeline/pipeline_sync_cache.go`, `store/installation_split_db.go`, `store/startup_split_db.go`, `store/sites_split_db.go`, `store/schedule_split_db.go`, and `store/storage_inventory.go`.

### Subtask 03: PR SQLite Storage Engine (`cli/prdb/`) & ERD Parity
- Implemented `cli/prdb/types.go` with `PullRequestRecord`, `PrReleaseRecord`, `PrBranchRecord`, and monadic Result wrappers (`PullRequestResult`, `PrBranchListResult`, `PrReleaseResult`).
- Added DDL in `cli/prdb/pr_split_schema.go` and `cli/constants/constants_split_db_sql.go` for tables `PullRequest`, `PrRelease`, `PrBranch` with PascalCase columns, indexes, and backward-compatible snake_case views (`pull_requests`, `pr_releases`, `pr_branches`).
- Implemented connection opener `OpenPrSplitDb` in `cli/prdb/pr_split_conn.go` configuring WAL mode, 5000ms busy timeout, and `SetMaxOpenConns(1)`.
- Implemented monadic CRUD operations in `cli/prdb/pr_split_ops.go`: `CreatePullRequest`, `GetPullRequestByNumber`, `UpdatePullRequestStatus`, `AddPrRelease`, `UpsertPrBranch`, `ListActivePrBranches`, `ListMergedPrBranches`, and `MarkPrBranchDeleted`.
- Updated `02-spec/21-app/gitmap-database-erd.mmd` to achieve 100% table-name parity, verified by `TestERDMatchesSQLCreate`.

### Subtask 04: PR Merge Detection & Feature Branch Engine
- Implemented `cli/committransfer/prdesc/generator.go` generating structured Markdown PR descriptions with metadata headers, executive summaries, component impact tables, commit logs, line deltas, and safety checklists.
- Refactored `cli/committransfer/pr_engine.go` to simulate local Git PR workflows: detects merges (`git rev-parse <sha>^@` > 1), creates local feature branches (`feature/<slug>` or `pr/<id>`), replays commits onto feature branches, records PRs to SQLite, merges into mainline (`git merge --no-ff`), and closes PRs.
- Updated `cli/committransfer/replay.go` and `plan.go` to route merge commits through the PR simulation engine when PR mode is active.

### Subtask 05: Terminal UI & Visual Execution Graph
- Enhanced `cli/committransfer/log.go` with styled `pterm` headers, status badges (`[REPLAYING]`, `[PR CREATED]`, `[PR MERGED]`, `[SNAPSHOT SYNCED]`), and real-time commit counters.
- Implemented `cli/committransfer/graph/renderer.go` rendering an ASCII/Unicode visual execution graph illustrating mainline commits, feature branch tracks, PR merge nodes, and release tags.
- Wired visual execution graph rendering at the conclusion of commit replay runs in `cli/committransfer/run.go`.

### Subtask 06: Final Snapshot Sync & PR Clean Command
- Implemented `FinalizeSnapshotSync(sourceDir, targetDir, sourceHeadSha, cmdName)` in `cli/committransfer/replay.go` and wired it into `cli/cmd/commitin/orchestrator/pipeline.go`:
  - Prunes target files not present in source.
  - Copies all source files into target.
  - Creates a final synchronization commit (`chore(sync): synchronize final repository snapshot tree to match source <shortsha>`) if diffs exist, ensuring byte-for-byte and file-for-file equality with the latest commit/release.
- Implemented `cli/committransfer/prclean/pr_clean.go`:
  - Inspects merged PR branches in SQLite (`ListMergedPrBranches()`).
  - Renders a styled pre-flight table.
  - Prompts confirmation (or bypasses via `--yes`).
  - Deletes local git branches (`git branch -D`) and updates SQLite records.

### Subtask 07: Help Text, HG (Help Groups) & JSON Examples
- Created comprehensive documentation in `cli/helptext/pr.md`.
- Restored `FuncIntel` and `SEO Commit Scheduling` in `cli/cmd/commitin/help.go`, passing `TestCommitInHelp`.
- Wired `HelpGroupPR` in root usage, compact usage, and filter rows (`rootusage.go`, `rootusage_groups.go`, `rootusagecompact.go`, `rootusagefilter_rows.go`).
- Added `pr-in` to `CompactPR` in `cli/constants/constants_helpgroups.go`.
- Linked in `02-spec/21-app/129-pr-commit-engines-and-sqlite-split-db.md` and `.ai-memory/what-to-read.md`.

---

## 3. Verification & Quality Matrix

| Rule / Guideline | Requirement | Status | Verification Evidence |
|---|---|---|---|
| Function Sizing | Functions <= 8–15 lines | PASS | All functions refactored and decomposed to <= 11 lines. |
| Boolean Prefixes | Affirmative `is*`, `has*` only | PASS | Verified via `check-enum-and-boolean.py` (2521 files scanned). |
| Nested Ifs | Zero nested if statements | PASS | Verified via `check-nested-ifs.py` (3345 files scanned). |
| Error Management | Universal AppError wrappers | PASS | Verified via `check-error-management.py` (3391 files scanned). |
| Unit Tests | Package test suites pass | PASS | `go test ./cmd`, `go test ./cmd/commitin`, `go test ./store -run TestERDMatchesSQLCreate` all passed. |
| CI Scripts | Python CI test suite | PASS | `.github/scripts/tests/test_ci_scripts.py` ran 18 tests in 50.4s: OK. |
