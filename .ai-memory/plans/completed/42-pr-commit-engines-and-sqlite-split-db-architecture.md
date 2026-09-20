# Plan 42: PR Commit Engines, Split-DB Architecture, and Auto-Merge PR Release Suite

> **Execution Lifecycle & Header Summary:**
> - **Task Inception:** Initiated under the `[V2] Parent Task N-Step Continuous Loop & Multi-Agent Orchestration` workflow (`N=500` self-loop budget) based on **Spec 129** (`02-spec/21-app/129-pr-commit-engines-and-sqlite-split-db.md`), Spec 106, Spec 114, Spec 124, and user directives.
> - **Scope & Purpose:** Introduces the first-class `pr` (and `pull-request`) commit replay command family, standardizes the canonical SQLite Split-DB directory formula (`.gitmap/data/<section>/<slug>/sql.db`) across all GitMap subsystems with auto-migration, establishes the dedicated `prdb` storage engine, provides automated merge detection with feature branch forking (`feature/<slug>` or `pr/<id>`), rich Markdown PR description generation, mainline merge (`--no-ff`), rich `pterm` terminal UI, visual ASCII execution graphs, deterministic final snapshot synchronization, and merged PR branch pruning (`pr-clean` / `pr-rm`).
> - **Total Execution Steps / Loops:** 7 discrete subtasks executed across parallel sub-agents (Subtask 01: Split-DB Architecture Engine, Subtask 02: PR SQLite Storage Engine, Subtask 03: PR Command & Flags Engine, Subtask 04: PR Merge Detection & Feature Branch Engine, Subtask 05: Terminal UI & Execution Graph Engine, Subtask 06: Final Snapshot Sync & PR Clean Engine, Subtask 07: Spec 129, Help Text & Catalog Sync).
> - **Final Status:** 100% COMPLETED (All deliverables implemented, guideline verified, and consolidated).

---

## 1. Executive Summary & Full Architectural Scope

Plan 42 delivers a major evolution of GitMap's commit transfer and repository management capabilities:

```
                                  ┌────────────────────────────────────────────────────────┐
                                  │            GitMap PR & Commit Replay Engine            │
                                  └──────────────────────────┬─────────────────────────────┘
                                                             │
         ┌───────────────────────────────────────────────────┼───────────────────────────────────────────────────┐
         ▼                                                   ▼                                                   ▼
   ┌───────────┐                                       ┌───────────┐                                       ┌───────────┐
   │ Subtask 1 │                                       │ Subtask 2 │                                       │ Subtask 3 │
   │ Split-DB  │                                       │   prdb    │                                       │PR Commands│
   ├───────────┤                                       ├───────────┤                                       ├───────────┤
   │split_db_pa│                                       │types.go   │                                       │constants  │
   │automation │                                       │schema.go  │                                       │dispatch   │
   │pipeline   │                                       │conn.go    │                                       │committrans│
   │stores     │                                       │ops.go     │                                       │rootrelease│
   └───────────┘                                       └───────────┘                                       └───────────┘
         │                                                   │                                                   │
         ├───────────────────────────────────────────────────┼───────────────────────────────────────────────────┤
         ▼                                                   ▼                                                   ▼
   ┌───────────┐                                       ┌───────────┐                                       ┌───────────┐
   │ Subtask 4 │                                       │ Subtask 5 │                                       │ Subtask 6 │
   │Merge / PR │                                       │Terminal UI│                                       │Final Sync │
   ├───────────┤                                       ├───────────┤                                       ├───────────┤
   │pr_engine  │                                       │log.go     │                                       │replay.go  │
   │prdesc     │                                       │graph/     │                                       │pipeline.go│
   │replay.go  │                                       │run.go     │                                       │prclean/   │
   └───────────┘                                       └───────────┘                                       └───────────┘
```

---

## 2. Deliverables & Subtask Ledger

### Subtask 01: SQLite Split-DB Canonical Directory Structure Standardization
- Created `cli/store/split_db_path.go` implementing `ResolveSplitDbPath(section, slug, repoRoot)` and `ResolveSplitDbDir(section, slug, repoRoot)`.
- Standardized the database path formula across the entire repository to `.gitmap/data/<section>/<slug>/sql.db` (and `<BinaryDataDir>/<section>/<slug>/sql.db` fallback).
- Defined section constants: `SectionAutomation`, `SectionPipeline`, `SectionPR`, `SectionLogs`, `SectionInstallation`, `SectionStartup`, `SectionSites`, `SectionSchedule`, `SectionRepoSearch`, and `DbFileName = "sql.db"`.
- Added transparent auto-migration moving legacy files (`.gitmap/data/<slug>/automation/sql.db`, `pipeline.db`, etc.) to `sql.db` in canonical directories.
- Refactored `cmdautomation/runtime_db.go`, `cmdautomation/exclusions.go`, `pipelinedb/pipeline_split_db.go`, `pipelinedb/pipeline_split_db_fs.go`, `cmdpipeline/pipeline_sync_cache.go`, `store/installation_split_db.go`, `store/startup_split_db.go`, `store/sites_split_db.go`, `store/schedule_split_db.go`, and `store/storage_inventory.go`.
- Added unit tests in `cli/store/split_db_path_test.go`.

### Subtask 02: PR SQLite Storage Engine (`cli/prdb/`)
- Implemented `cli/prdb/types.go` with `PullRequestRecord`, `PrReleaseRecord`, `PrBranchRecord`, and monadic Result wrappers (`PullRequestResult`, `PrBranchListResult`, `PrReleaseResult`).
- Added DDL in `cli/prdb/pr_split_schema.go` and `cli/constants/constants_split_db_sql.go` for tables `PullRequest`, `PrRelease`, `PrBranch` with PascalCase columns, indexes, and backward-compatible snake_case views (`pull_requests`, `pr_releases`, `pr_branches`).
- Implemented connection opener `OpenPrSplitDb` in `cli/prdb/pr_split_conn.go` configuring WAL mode, 5000ms busy timeout, and `SetMaxOpenConns(1)`.
- Implemented monadic CRUD operations in `cli/prdb/pr_split_ops.go`: `CreatePullRequest`, `GetPullRequestByNumber`, `UpdatePullRequestStatus`, `AddPrRelease`, `UpsertPrBranch`, `ListActivePrBranches`, `ListMergedPrBranches`, and `MarkPrBranchDeleted`.
- Added unit tests in `cli/prdb/pr_split_db_test.go`.

### Subtask 03: PR Commit Transfer Command Family & Flags
- Defined `CmdPR = "pr"`, `CmdPullRequest = "pull-request"`, `CmdPRClean = "pr-clean"`, `CmdPRRm = "pr-rm"`, and `CmdPRList = "pr-list"` in `cli/constants/constants_committransfer.go`.
- Renamed legacy `CmdPR = "pull-requests"` to `CmdPullRequests` and changed `CmdReleasePullAlias` from `"pr"` to `"relp"` / `"rlp"` in `cli/constants/constants_cli.go` and `cli/cmd/rootrelease.go`, eliminating symbol and CLI collisions.
- Added `HelpGroupPR = "  Pull Request & Merge Automation (PR):"` and registered compact help strings in `cli/constants/constants_helpgroups.go`.
- Routed PR commands in `cli/cmd/dispatchcommittransfer.go`.
- Implemented `runPRClean` and `runPRList` in `cli/cmd/committransfer.go`, defaulting `opts.PRMode = "merges"` for `pr` commands.

### Subtask 04: PR Merge Detection & Feature Branch Engine
- Implemented `cli/committransfer/prdesc/generator.go` generating structured Markdown PR descriptions with metadata headers, executive summaries, component impact tables, commit logs, line deltas, and safety checklists.
- Added unit tests in `cli/committransfer/prdesc/generator_test.go`.
- Refactored `cli/committransfer/pr_engine.go` to simulate local Git PR workflows: detects merges (`git rev-parse <sha>^@` > 1), creates local feature branches (`feature/<slug>` or `pr/<id>`), replays commits onto feature branches, records PRs to SQLite, merges into mainline (`git merge --no-ff`), and closes PRs.
- Updated `cli/committransfer/replay.go` and `plan.go` to route merge commits through the PR simulation engine when PR mode is active.

### Subtask 05: Terminal UI & Visual Execution Graph
- Enhanced `cli/committransfer/log.go` with styled `pterm` headers (`pterm.DefaultHeader`), status badges (`[REPLAYING]`, `[PR CREATED]`, `[PR MERGED]`, `[SNAPSHOT SYNCED]`), and real-time commit counters.
- Implemented `cli/committransfer/graph/renderer.go` rendering an ASCII/Unicode visual execution graph illustrating mainline commits, feature branch tracks, PR merge nodes, and release tags.
- Added unit tests in `cli/committransfer/graph/renderer_test.go`.
- Wired visual execution graph rendering at the conclusion of commit replay runs in `cli/committransfer/run.go`.

### Subtask 06: Final Snapshot Sync & PR Clean Command
- Implemented `FinalizeSnapshotSync(sourceDir, targetDir, sourceHeadSha, cmdName)` in `cli/committransfer/replay.go` and wired it into `cli/cmd/commitin/orchestrator/pipeline.go`:
  - Prunes target files not present in source.
  - Copies all source files into target.
  - Creates a final synchronization commit (`chore(sync): synchronize final repository snapshot tree to match source <shortsha>`) if diffs exist, ensuring byte-for-byte and file-for-file equality with the latest commit/release.
- Implemented `cli/committransfer/prclean/pr_clean.go` and `pr_clean_test.go`:
  - Inspects merged PR branches in SQLite (`ListMergedPrBranches()`).
  - Renders a styled pre-flight table.
  - Prompts confirmation (or bypasses via `--yes`).
  - Deletes local git branches (`git branch -D`) and updates SQLite records.

### Subtask 07: Spec 129, Comprehensive Help & Catalog Sync
- Created Spec 129 (`02-spec/21-app/129-pr-commit-engines-and-sqlite-split-db.md`).
- Cross-referenced Spec 129 in Spec 106 (`02-spec/21-app/106-commit-left-right-both.md`), Spec 114 (`02-spec/21-app/114-committransfer-idempotence-and-merge-default.md`), and Spec 124 (`02-spec/21-app/124-polyglot-worker-orchestrator-and-automation-runner.md`).
- Added `pr`, `pr-clean`, and `pr-list` to `readme.md` and `what-to-read.md`.
- Expanded `cli/cmd/commitin/help.go` with full JSON examples, exclusion rules, and PR workflow documentation.
- Implemented `cli/helptext/docs/cmd/pr.go` and `pr_test.go`.

---

## 3. Verification & Compliance Matrix

| Rule / Guideline | Requirement | Status | Verification Evidence |
|---|---|---|---|
| Function Sizing | Functions <= 8–15 lines | PASS | Verified via `05-guideline-autofixer.py` across all modified files. |
| Boolean Prefixes | Affirmative `is*`, `has*` only | PASS | Verified across `split_db_path.go`, `prdb`, `prclean`, `graph`, `committransfer`. |
| Monadic Results | `result.Result[T]` or `*apperror.AppError` | PASS | Zero raw `(T, error)` tuples introduced. |
| Ban on Tests/Builds | No `go test` / `go build` during routine execution | PASS | Respected 100%; zero compiler or test commands run. |
| Path Standards | Strictly relative git paths | PASS | All internal links and file paths use relative git paths. |
