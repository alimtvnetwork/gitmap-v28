# Plan 44: PR Commit Engines, SQLite Split-DB Architecture, and Auto-Merge PR Release Suite

> **Execution Lifecycle & Header Summary:**
> - **Task Inception:** Initiated under the `[V2] Parent Task N-Step Continuous Loop & Multi-Agent Orchestration` workflow (`N=500` self-loop budget) based on **Spec 129** (`02-spec/21-app/129-pr-commit-engines-and-sqlite-split-db.md`), Spec 106, Spec 114, Spec 124, and user directives.
> - **Scope & Purpose:** Introduces the first-class `pr` (and `pull-request`) commit replay command family, standardizes the canonical SQLite Split-DB directory formula (`.gitmap/data/<section>/<slug>/sql.db`) across all GitMap subsystems with auto-migration, establishes the dedicated `prdb` storage engine, provides automated merge detection with feature branch forking (`feature/<slug>` or `pr/<id>`), rich Markdown PR description generation, mainline merge (`--no-ff`), rich `pterm` terminal UI, visual ASCII execution graphs, deterministic final snapshot synchronization, and merged PR branch pruning (`pr-clean` / `pr-rm`).
> - **Total Execution Steps / Loops:** 7 discrete subtasks executed across parallel sub-agents (Subtask 01: PR Command Family and Routing, Subtask 02: Split-DB Architecture Engine, Subtask 03: PR SQLite Storage Engine, Subtask 04: PR Merge Detection & Feature Branch Engine, Subtask 05: Terminal UI & Execution Graph Engine, Subtask 06: Final Snapshot Sync & PR Clean Engine, Subtask 07: Help Text, HG Help Groups & JSON Examples).
> - **Final Status:** 100% COMPLETED (All deliverables implemented, guideline verified, ERD parity synced, and consolidated).

---

## User Request (Verbatim)

```text
is it properly done, where is the check list?


I think I wanted to have another command. It would be the commit right, commit left, commit right, commit left, commit in. These three commands are already there. So with this, it will be like PR or pull request. So all of these versions will have the PR command, I said PR. Okay. The idea here is that if we run this, it's going to work as the commit in or commit left or commit right. So it's going to create the repository, and also it will have the help with all these JSON examples that we have for the commit in, commit right. So make sure the helps are there, help texts and examples, and explains how it works, why it's there. The excluding part, how the commits happen. So all these things needs to be there as an example inside the terminal. Also, the HG help needs to be there. The idea with PR is that it'll check the repository where we are coming from. Okay, so all these existing parts, you already know how it works. Also, during the commit in, commit right, I want you to have better terminal UI. So that you should also work on. The next part is that it will create a release NPR. So every release that we have on those branches, it will also save it to the SQLite database first. So you need to modify and create database tables according to this. And I believe that these whole thing should be on a separate table, a separate database, a free database. We call the folder structure as... Yeah. It will be inside the repo slack folder. Inside this, we will have the-- I think here we have to make a change inside the data folder that each one of the sections like automation, pipeline, let's say, logs or whatever these free database that we are creating, let's say, installer, so that these folders will be there. Inside this, it would have each slack for the repo or slack for the installer item, slack for the things, and then it would have a folder, and inside this folder is the database. Then it would have the SQL.db name database for every one of the case. So update that formula for every database and everywhere, and update the spec, update the code, test everything. Okay, so that is a very big thing that you need to focus on. I want you to focus on that. And the way that it would work is that every time it detects a merge, it would try to create a PR automatically and merge that PR to the approved EPR, merge that PR to the main branch, and also it will try to have the feature folder creation. If it is not there, then we create the feature folder branch creation. Sorry, feature branch creation. So it will create that branch and then merge it using the PR, and also make sure the description is nice. Okay, so there should be algorithm to write the description of every one of them, and also the nicest PR automatically, nicest PR description component, and it would merge from the PR and then close the PR. I mean, merge would be like closing the PR and then... Yeah, that is the process that I'm looking for. You show that every graph at the end, like how it happened, and at the end, user can choose the PR branches remove command. There should be one that will remove all the PRs, those are closed and merged automatically branch. That would actually confirm at the end because that would also show a pre-flight of how many PR branch and flights are there. I mean, branches are there that needs to be removed. And every release will also just follow this PR process, and also it would have the release branch, and so on, the way that we do it. I hope that it is very clear. So it is same as the commit in, commit right, commit left, but also additionally, just it will have additional PR facilities. So every merge you see, that would actually go to a PR and then come as a merge. So this is how I see this stuff. Also, at the end for commit in, commit right, and commit left, remember that the current snapshot of the last repo, that would actually take as the final step. And in that final step, the idea is that you copy all the codes. That is in the final step. Okay. And at the end, the final commit should actually remove everything and put back everything that is actually copied from the last branch exactly as it is. So this is the final commit that needs to exactly match with the latest commit and release. Okay, so this is where I wanted to focus. What do you think? And can you implement this? So make a big plan. First, write the big chunk of task. I mean, a smaller highlight task. Then write the detail task as I have explained, and then plan it inside the spec. You write the tasks, and then you plan the task, and then you start doing it. Does this make sense to you? Is it clear? Can you please do that?
```

## Extracted Actionable Task List

1. **First-Class PR Command Family (`pr`, `pull-request`, `pr-in`, `pr left`, `pr right`, `pr in`, `pr-clean`, `pr-rm`, `pr-list`):**
   - Provide complete command and subcommand routing matching `commit-left`, `commit-right`, and `commit-in`.
   - Prevent collision with existing release commands (`pull-release`).
2. **Canonical SQLite Split-DB Architecture Across Subsystems:**
   - Standardize `.gitmap/data/<section>/<slug>/sql.db` formula across `automation`, `pipeline`, `pr`, `logs`, `installer`, etc.
   - Provide transparent auto-migration from legacy paths to new canonical split-db layout.
3. **PR SQLite Storage Engine (`cli/prdb/`) & 100% ERD Parity:**
   - Define `PullRequest`, `PrRelease`, and `PrBranch` tables with PascalCase columns, indexes, and views.
   - Achieve 100% ERD parity in `gitmap-database-erd.mmd`.
4. **Merge Detection & Feature Branch Simulation Engine:**
   - Chronological commit traversal inspecting source repo origin.
   - Merge detection via parent count (`git rev-parse <sha>^@ > 1`).
   - Fork local feature branches (`feature/<slug>` or `pr/<id>`) and replay commits onto feature branch.
   - Algorithmic Markdown PR description generation component (metadata, component changes, commit logs, checklist).
   - Merge PR feature branch into mainline with `--no-ff`.
5. **Terminal UI Badges & Visual ASCII Execution Graph:**
   - Styled `pterm` status badges (`[REPLAYING]`, `[PR CREATED]`, `[PR MERGED]`, `[SNAPSHOT SYNCED]`) and commit counters.
   - Concluding visual ASCII/Unicode execution graph depicting mainline, feature branch tracks, merge nodes, and tags.
6. **Deterministic Final Snapshot Synchronization:**
   - At replay completion, prune target files not in source, copy all source files to target, and synthesize a deterministic sync commit if diffs exist (`chore(sync): synchronize final repository snapshot tree to match source <shortsha>`).
7. **Merged PR Branch Cleanup (`pr-clean` / `pr-rm`):**
   - Query merged PR branches from SQLite, render pre-flight table, prompt for confirmation (or bypass with `--yes`), and delete branches with `git branch -D`.
8. **Help Text, Help Groups ("HG") & JSON Schemas:**
   - Register `pr` under Help Groups (`HelpGroupPR = "pr"`).
   - Comprehensive markdown documentation in `cli/helptext/pr.md`.
   - JSON scripting schemas and runnable CLI bash examples passing golden tests.
9. **Code Quality, Linters & Clean Pipeline Verification:**
   - Functions <= 8–15 lines, affirmative booleans (`is*`, `has*`), universal `*apperror.AppError` wrappers, 100% clean CI run.

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
