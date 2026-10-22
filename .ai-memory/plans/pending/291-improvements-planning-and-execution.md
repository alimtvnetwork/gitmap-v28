# Master Plan 291: GitMap 30 Improvements Planning and Execution (Batch 3)

## Status: PENDING
- **Plan ID:** 291
- **Slug:** `291-improvements-planning-and-execution`
- **Master Spec Reference:** [01-architecture-spec.md](../../../02-spec/21-app/291-improvements-planning-and-execution/01-architecture-spec.md)
- **Component Spec Reference:** [02-component-spec.md](../../../02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md)
- **Migration & Verification Spec Reference:** [03-migration-and-verification-spec.md](../../../02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md)
- **Scope:** 30 High-Impact Architectural, UI, Guidelines, and Operational Safety Improvements (Batch 3)
- **Execution Mandate:** Zero-Build, Zero-Test, Zero Git Commands, Strictly Relative Lowercase Git Paths

---

## 1. Executive Summary & Wave Organization

This master plan coordinates the specification and staged execution of Batch 3 of the 30 discrete architectural improvements for GitMap. The work is structured into three execution waves partitioned by functional domain and file boundary isolation to ensure zero concurrent write collisions:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│ Wave 1: Error Sovereignty, Codebase Guidelines & Types (Subtasks 01–10)      │
│ Spec: 02-spec/21-app/291-improvements-planning-and-execution/01-architecture │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 2: Terminal UI Modernization, Lipgloss & Visual Rhythm (Subtasks 11–20) │
│ Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component    │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 3: System UX, Operational Safety & DB Introspection (Subtasks 21–30)    │
│ Spec: 02-spec/21-app/291-improvements-planning-and-execution/03-migration    │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Complete 30-Subtask Registry & Traceability Matrix

### Wave 1: Error Sovereignty, Codebase Guidelines & Type Hardening (Subtasks 01–10)

- [ ] **Subtask-01:** [subtask-01.md](../subtasks/291-improvements-planning-and-execution/subtask-01.md)
  - **Target:** `cli/cmdclone/clonerun.go`, `cli/cmdclone/clone_direct.go`
  - **Objective:** Convert void clone helpers to return `error`, bubble errors up to caller.
  - **Status:** Pending
- [ ] **Subtask-02:** [subtask-02.md](../subtasks/291-improvements-planning-and-execution/subtask-02.md)
  - **Target:** `cli/cmdclone/clone_manifest.go`
  - **Objective:** Convert `cloneManifestBatch` and helper error discards into structured error returns.
  - **Status:** Pending
- [ ] **Subtask-03:** [subtask-03.md](../subtasks/291-improvements-planning-and-execution/subtask-03.md)
  - **Target:** `cli/cmdclone/clonepretty.go`
  - **Objective:** Leaf error sovereignty, replace direct exit handler calls with returned errors.
  - **Status:** Pending
- [ ] **Subtask-04:** [subtask-04.md](../subtasks/291-improvements-planning-and-execution/subtask-04.md)
  - **Target:** `cli/cmdhelp/helpcheck.go`, `cli/cmdhelp/help_topic.go`
  - **Objective:** Convert void error discards in help lookup into direct error returns.
  - **Status:** Pending
- [ ] **Subtask-05:** [subtask-05.md](../subtasks/291-improvements-planning-and-execution/subtask-05.md)
  - **Target:** `cli/cmd/rootutility_help.go`
  - **Objective:** Bubble up errors from root utility help dispatch.
  - **Status:** Pending
- [ ] **Subtask-06:** [subtask-06.md](../subtasks/291-improvements-planning-and-execution/subtask-06.md)
  - **Target:** `cli/cmddb/cmddb_dispatch.go`, `cli/cmddb/db_vacuum.go`
  - **Objective:** Propagate vacuum and DB maintenance errors without swallowing.
  - **Status:** Pending
- [ ] **Subtask-07:** [subtask-07.md](../subtasks/291-improvements-planning-and-execution/subtask-07.md)
  - **Target:** `cli/cmdpipeline/pipeline_logs_cmd.go`
  - **Objective:** Handle stream closing and pipe write errors properly.
  - **Status:** Pending
- [ ] **Subtask-08:** [subtask-08.md](../subtasks/291-improvements-planning-and-execution/subtask-08.md)
  - **Target:** `cli/constants/constants_git.go`
  - **Objective:** Define `GitRevList = "rev-list"`, `GitCountFlag = "--count"`, `GitSymbolicRef = "symbolic-ref"`, etc. with typed `GitCommandType`.
  - **Status:** Pending
- [ ] **Subtask-09:** [subtask-09.md](../subtasks/291-improvements-planning-and-execution/subtask-09.md)
  - **Target:** `cli/constants/constants_store.go`
  - **Objective:** Add `TablePipelineRun = "PipelineRun"`, `TableHeatmapActionLog = "HeatmapActionLog"` with typed `DatabaseTableType`.
  - **Status:** Pending
- [ ] **Subtask-10:** [subtask-10.md](../subtasks/291-improvements-planning-and-execution/subtask-10.md)
  - **Target:** `cli/store/repo_query.go`
  - **Objective:** Clean value semantics, <= 15 lines, positive booleans.
  - **Status:** Pending

---

### Wave 2: Terminal UI Modernization, Lipgloss & Visual Rhythm (Subtasks 11–20)

- [ ] **Subtask-11:** [subtask-11.md](../subtasks/291-improvements-planning-and-execution/subtask-11.md)
  - **Target:** `cli/termout/theme_lipgloss.go`
  - **Objective:** Implement `TokyoTag(label, value string) string` with Tokyo Night tokens.
  - **Status:** Pending
- [ ] **Subtask-12:** [subtask-12.md](../subtasks/291-improvements-planning-and-execution/subtask-12.md)
  - **Target:** `cli/termout/theme_lipgloss.go`
  - **Objective:** Implement `TokyoMetricPill(label string, count int) string`.
  - **Status:** Pending
- [ ] **Subtask-13:** [subtask-13.md](../subtasks/291-improvements-planning-and-execution/subtask-13.md)
  - **Target:** `cli/termout/render.go`
  - **Objective:** Implement `FormatSectionDivider(title string, theme Theme) string` and `PrintSectionDivider(title string, theme Theme)`.
  - **Status:** Pending
- [ ] **Subtask-14:** [subtask-14.md](../subtasks/291-improvements-planning-and-execution/subtask-14.md)
  - **Target:** `cli/termout/termtable_types.go`
  - **Objective:** Add `AlignmentType` (`AlignLeft`, `AlignRight`, `AlignCenter`) to `ColumnConfig`.
  - **Status:** Pending
- [ ] **Subtask-15:** [subtask-15.md](../subtasks/291-improvements-planning-and-execution/subtask-15.md)
  - **Target:** `cli/termout/table.go`
  - **Objective:** Implement column alignment rendering using `AlignmentType` and `VisualStringWidth`.
  - **Status:** Pending
- [ ] **Subtask-16:** [subtask-16.md](../subtasks/291-improvements-planning-and-execution/subtask-16.md)
  - **Target:** `cli/termout/truncate.go`
  - **Objective:** Implement `TruncateMiddleWithEllipsis(text string, maxWidth int) string` using visual width and rune slicing.
  - **Status:** Pending
- [ ] **Subtask-17:** [subtask-17.md](../subtasks/291-improvements-planning-and-execution/subtask-17.md)
  - **Target:** `cli/cmdpipeline/pipeline_history_render.go`
  - **Objective:** Modernize failure cards with `TokyoBoxStyle` and `TokyoTag` metrics.
  - **Status:** Pending
- [ ] **Subtask-18:** [subtask-18.md](../subtasks/291-improvements-planning-and-execution/subtask-18.md)
  - **Target:** `cli/cmdcd/cd_table.go`
  - **Objective:** Create dedicated CD selection table renderer printing to `os.Stderr`.
  - **Status:** Pending
- [ ] **Subtask-19:** [subtask-19.md](../subtasks/291-improvements-planning-and-execution/subtask-19.md)
  - **Target:** `cli/cmdcd/cdops.go`
  - **Objective:** Wire CD selection table renderer to use `cd_table.go`.
  - **Status:** Pending
- [ ] **Subtask-20:** [subtask-20.md](../subtasks/291-improvements-planning-and-execution/subtask-20.md)
  - **Target:** `cli/cmdcg/cg_doctor.go`
  - **Objective:** Enhance doctor summary cards with `TokyoMetricPill` and section dividers.
  - **Status:** Pending

---

### Wave 3: System UX, Operational Safety & Database Introspection (Subtasks 21–30)

- [ ] **Subtask-21:** [subtask-21.md](../subtasks/291-improvements-planning-and-execution/subtask-21.md)
  - **Target:** `cli/store/repo.go`
  - **Objective:** Add `ListReposByBranch(branch string) ([]model.ScanRecord, error)` to filter repositories by branch.
  - **Status:** Pending
- [ ] **Subtask-22:** [subtask-22.md](../subtasks/291-improvements-planning-and-execution/subtask-22.md)
  - **Target:** `cli/store/repo.go`
  - **Objective:** Add `ListReposDirtyStatus() ([]model.ScanRecord, error)` to query repositories flagged dirty.
  - **Status:** Pending
- [ ] **Subtask-23:** [subtask-23.md](../subtasks/291-improvements-planning-and-execution/subtask-23.md)
  - **Target:** `cli/store/heatmap_store.go`
  - **Objective:** Add `PurgeStaleActionLogs(olderThanDays int) (int64, error)` safely guarding active tasks.
  - **Status:** Pending
- [ ] **Subtask-24:** [subtask-24.md](../subtasks/291-improvements-planning-and-execution/subtask-24.md)
  - **Target:** `cli/heatmap/heatmap_tracker.go`
  - **Objective:** Add `TrackBatchFileTouches(cwd string, touches []store.FileTouchCount) error` with single transaction and prepared statement.
  - **Status:** Pending
- [ ] **Subtask-25:** [subtask-25.md](../subtasks/291-improvements-planning-and-execution/subtask-25.md)
  - **Target:** `cli/cmdagent/agent_ui_logs.go`
  - **Objective:** Add `parseQueryOffset(val string) int` and wire `offset` query parameter into pagination.
  - **Status:** Pending
- [ ] **Subtask-26:** [subtask-26.md](../subtasks/291-improvements-planning-and-execution/subtask-26.md)
  - **Target:** `cli/gitutil/git_repo_query.go`
  - **Objective:** Add `GetRepoDirtyFiles(repoPath string) ([]string, error)` parsing `git status --porcelain`.
  - **Status:** Pending
- [ ] **Subtask-27:** [subtask-27.md](../subtasks/291-improvements-planning-and-execution/subtask-27.md)
  - **Target:** `cli/gitutil/git_repo_query.go`
  - **Objective:** Add `IsRepoDetached(repoPath string) (bool, error)` via `git symbolic-ref -q HEAD`.
  - **Status:** Pending
- [ ] **Subtask-28:** [subtask-28.md](../subtasks/291-improvements-planning-and-execution/subtask-28.md)
  - **Target:** `cli/cmdpe/pe_target.go`
  - **Objective:** Extract target directory resolution and filtering helpers from `pe_gh.go` and `pe_triage.go`.
  - **Status:** Pending
- [ ] **Subtask-29:** [subtask-29.md](../subtasks/291-improvements-planning-and-execution/subtask-29.md)
  - **Target:** `cli/cmdpe/pe_rerun.go`
  - **Objective:** Decompose oversized `RunPeRerun` into discrete helpers <= 15 lines.
  - **Status:** Pending
- [ ] **Subtask-30:** [subtask-30.md](../subtasks/291-improvements-planning-and-execution/subtask-30.md)
  - **Target:** `cli/cmdclone/remotetransport.go`
  - **Objective:** Verify and flatten remaining conditional branching with positive booleans and early returns.
  - **Status:** Pending
