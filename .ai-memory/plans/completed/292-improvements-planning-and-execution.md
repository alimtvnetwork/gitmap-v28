# Master Plan 292: GitMap 30 Improvements Planning and Execution (Batch 4)

## Status: COMPLETED
- **Plan ID:** 292
- **Slug:** `292-improvements-planning-and-execution`
- **Master Spec Reference:** [01-architecture-spec.md](../../../02-spec/21-app/292-improvements-planning-and-execution/01-architecture-spec.md)
- **Component Spec Reference:** [02-component-spec.md](../../../02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md)
- **Migration & Verification Spec Reference:** [03-migration-and-verification-spec.md](../../../02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md)
- **Scope:** 30 High-Impact Architectural, UI, Guidelines, and Operational Safety Improvements (Batch 4 / Improvements 91–120)
- **Execution Mandate:** Zero-Build, Zero-Test, Zero Git Commands, Strictly Relative Lowercase Git Paths

---

## 1. Executive Summary & Wave Organization

This master plan coordinates the specification and staged execution of Batch 4 of the 30 discrete architectural improvements for GitMap. The work is structured into three execution waves partitioned by functional domain and file boundary isolation to ensure zero concurrent write collisions:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│ Wave 1: Error Sovereignty, Codebase Guidelines & Constants (Subtasks 01–10) │
│ Spec: 02-spec/21-app/292-improvements-planning-and-execution/01-architecture │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 2: Terminal UI Modernization, Lipgloss & Visual Rhythm (Subtasks 11–20) │
│ Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component    │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 3: System UX, Operational Safety & DB Introspection (Subtasks 21–30)    │
│ Spec: 02-spec/21-app/292-improvements-planning-and-execution/03-migration    │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Complete 30-Subtask Registry & Traceability Matrix

### Wave 1: Error Sovereignty, Codebase Guidelines & Constants (Subtasks 01–10)

- [x] **Subtask-01:** [subtask-01.md](../subtasks/292-improvements-planning-and-execution/subtask-01.md)
  - **Target:** `cli/cmdstatus/status.go`, `cli/cmdstatus/helpers_shim.go`
  - **Objective:** Convert void status helpers that swallow errors (`loadRecordsByGroup`, `loadAllRecordsDB`, `handleStatusDBError`) into direct `error` returns.
  - **Status:** Completed
- [x] **Subtask-02:** [subtask-02.md](../subtasks/292-improvements-planning-and-execution/subtask-02.md)
  - **Target:** `cli/cmdhaschange/haschange.go`, `cli/cmdhaschange/helpers_shim.go`
  - **Objective:** Convert `hasChange` and `isWorkingTreeDirty` from bare boolean returns swallowing errors to `(bool, error)`.
  - **Status:** Completed
- [x] **Subtask-03:** [subtask-03.md](../subtasks/292-improvements-planning-and-execution/subtask-03.md)
  - **Target:** `cli/cmd/cmd_see.go`
  - **Objective:** Replace swallowed errors in `RunRepoManageUI` with returned `*apperror.AppError`; bubble up.
  - **Status:** Completed
- [x] **Subtask-04:** [subtask-04.md](../subtasks/292-improvements-planning-and-execution/subtask-04.md)
  - **Target:** `cli/cmdcpar/cpar.go`
  - **Objective:** Eliminate swallowed errors in batch commit/push/add/remote; check `git add` error and propagate batch failure state.
  - **Status:** Completed
- [x] **Subtask-05:** [subtask-05.md](../subtasks/292-improvements-planning-and-execution/subtask-05.md)
  - **Target:** `cli/cmdtemprelease/temprelease.go`, `cli/cmdtemprelease/tempreleaselist.go`
  - **Objective:** Eliminate void discards; bubble errors up to caller.
  - **Status:** Completed
- [x] **Subtask-06:** [subtask-06.md](../subtasks/292-improvements-planning-and-execution/subtask-06.md)
  - **Target:** `cli/cmdfixauth/fixauth.go`, `cli/cmdfixauth/fixcredential.go`
  - **Objective:** Convert void helpers (`unsetLegacyCredentialHelper`, `applyCredentialStoreType`) to return `*apperror.AppError`.
  - **Status:** Completed
- [x] **Subtask-07:** [subtask-07.md](../subtasks/292-improvements-planning-and-execution/subtask-07.md)
  - **Target:** `cli/cmdrelease/releasepending.go`, `cli/cmdrelease/releaserecentclone.go`
  - **Objective:** Replace leaf `cliexit.HandleUsageError` with returned errors; check `RunRelease` return.
  - **Status:** Completed
- [x] **Subtask-08:** [subtask-08.md](../subtasks/292-improvements-planning-and-execution/subtask-08.md)
  - **Target:** `cli/constants/constants_git.go`
  - **Objective:** Add `GitInspectCommandType` (`GitCmdCherry`, `GitCmdDiff`, `GitCmdMergeBase`) and `GitFlagType` (`GitFlagStat`).
  - **Status:** Completed
- [x] **Subtask-09:** [subtask-09.md](../subtasks/292-improvements-planning-and-execution/subtask-09.md)
  - **Target:** `cli/constants/constants_store.go`
  - **Objective:** Add `TableNameHeatmapLocalTracking` and `TableNamePipelineErrorLog` to `TableNameType` and string constants.
  - **Status:** Completed
- [x] **Subtask-10:** [subtask-10.md](../subtasks/292-improvements-planning-and-execution/subtask-10.md)
  - **Target:** `cli/store/repo_query_filter.go`
  - **Objective:** Design clean value semantics and extraction for tag and unpushed repository queries.
  - **Status:** Completed

---

### Wave 2: Terminal UI Modernization, Lipgloss & Visual Rhythm (Subtasks 11–20)

- [x] **Subtask-11:** [subtask-11.md](../subtasks/292-improvements-planning-and-execution/subtask-11.md)
  - **Target:** `cli/termout/theme_lipgloss.go`
  - **Objective:** Implement `TokyoBadge(label, bgHex, fgHex string) string` and `TokyoHeader(title string) string`.
  - **Status:** Completed
- [x] **Subtask-12:** [subtask-12.md](../subtasks/292-improvements-planning-and-execution/subtask-12.md)
  - **Target:** `cli/termout/render.go`
  - **Objective:** Implement `FormatKeyValuePair` and `PrintKeyValuePair` with visual column alignment via `VisualStringWidth`.
  - **Status:** Completed
- [x] **Subtask-13:** [subtask-13.md](../subtasks/292-improvements-planning-and-execution/subtask-13.md)
  - **Target:** `cli/termout/termtable_types.go`
  - **Objective:** Add `WrapModeType` enum (`WrapClip`, `WrapEllipsis`, `WrapWord`) and `RowStripe` to `TableConfig`.
  - **Status:** Completed
- [x] **Subtask-14:** [subtask-14.md](../subtasks/292-improvements-planning-and-execution/subtask-14.md)
  - **Target:** `cli/termout/table.go`
  - **Objective:** Implement cell wrap modes and alternating row striping with subtle Tokyo surface background tinting.
  - **Status:** Completed
- [x] **Subtask-15:** [subtask-15.md](../subtasks/292-improvements-planning-and-execution/subtask-15.md)
  - **Target:** `cli/termout/truncate.go`
  - **Objective:** Implement `TruncateLeftWithEllipsis(text string, maxWidth int) string` using rune slicing and visual width.
  - **Status:** Completed
- [x] **Subtask-16:** [subtask-16.md](../subtasks/292-improvements-planning-and-execution/subtask-16.md)
  - **Target:** `cli/cmdstatus/status_render.go`
  - **Objective:** Extract status table renderer using `termout.PrintTable` and `BorderRounded`.
  - **Status:** Completed
- [x] **Subtask-17:** [subtask-17.md](../subtasks/292-improvements-planning-and-execution/subtask-17.md)
  - **Target:** `cli/cmdhaschange/haschange.go`
  - **Objective:** Add TokyoTag formatting for `--all` while keeping single-mode stdout pure for scripting.
  - **Status:** Completed
- [x] **Subtask-18:** [subtask-18.md](../subtasks/292-improvements-planning-and-execution/subtask-18.md)
  - **Target:** `cli/cmdpending/pending_commits_render.go`
  - **Objective:** Modernize table header/footer with `BorderRounded` and `TokyoBoxStyle`.
  - **Status:** Completed
- [x] **Subtask-19:** [subtask-19.md](../subtasks/292-improvements-planning-and-execution/subtask-19.md)
  - **Target:** `cli/cmd/cmd_see.go`
  - **Objective:** Modernize `printPendingCommits` and help menu with Tokyo Night styling.
  - **Status:** Completed
- [x] **Subtask-20:** [subtask-20.md](../subtasks/292-improvements-planning-and-execution/subtask-20.md)
  - **Target:** `cli/cmdpending/pending_commits_filter.go`
  - **Objective:** Isolate interactive candidate selection prompts strictly to `os.Stderr`.
  - **Status:** Completed

---

### Wave 3: System UX, Operational Safety & Database Introspection (Subtasks 21–30)

- [x] **Subtask-21:** [subtask-21.md](../subtasks/292-improvements-planning-and-execution/subtask-21.md)
  - **Target:** `cli/store/repo_query_filter.go`
  - **Objective:** Add `ListReposByTag(tag string) ([]model.ScanRecord, error)` to query repositories filtered by metadata tags or notes.
  - **Status:** Completed
- [x] **Subtask-22:** [subtask-22.md](../subtasks/292-improvements-planning-and-execution/subtask-22.md)
  - **Target:** `cli/store/repo_query_filter.go`
  - **Objective:** Add `ListReposUnpushed() ([]model.ScanRecord, error)` to query repositories flagged with unpushed commits in metadata.
  - **Status:** Completed
- [x] **Subtask-23:** [subtask-23.md](../subtasks/292-improvements-planning-and-execution/subtask-23.md)
  - **Target:** `cli/store/heatmap_store.go`, `cli/store/heatmap_schema.go`
  - **Objective:** Define `HeatmapTaskSummary` struct and implement `GetHeatmapTaskSummary(taskSlug string) (*HeatmapTaskSummary, error)`.
  - **Status:** Completed
- [x] **Subtask-24:** [subtask-24.md](../subtasks/292-improvements-planning-and-execution/subtask-24.md)
  - **Target:** `cli/heatmap/heatmap_tracker_agent.go`
  - **Objective:** Add `TrackAgentAction(taskSlug, actionType, summary string) error` in dedicated companion file.
  - **Status:** Completed
- [x] **Subtask-25:** [subtask-25.md](../subtasks/292-improvements-planning-and-execution/subtask-25.md)
  - **Target:** `cli/cmdagent/agent_ui_logs.go`, `cli/cmdagent/agent_ui_server.go`
  - **Objective:** Add `SortBy` and `Order` to `ActionLogFilter`; implement whitelisted sort parsing and dynamic SQL ordering.
  - **Status:** Completed
- [x] **Subtask-26:** [subtask-26.md](../subtasks/292-improvements-planning-and-execution/subtask-26.md)
  - **Target:** `cli/gitutil/git_repo_commits.go`
  - **Objective:** Add `GetRepoUnpushedCommitCount(repoPath string) (int, error)` via `git rev-list --count @{u}..HEAD`.
  - **Status:** Completed
- [x] **Subtask-27:** [subtask-27.md](../subtasks/292-improvements-planning-and-execution/subtask-27.md)
  - **Target:** `cli/gitutil/git_repo_commits.go`
  - **Objective:** Add `GetRepoMergeBase(repoPath, ref1, ref2 string) (string, error)` via `git merge-base ref1 ref2`.
  - **Status:** Completed
- [x] **Subtask-28:** [subtask-28.md](../subtasks/292-improvements-planning-and-execution/subtask-28.md)
  - **Target:** `cli/cmdcpar/cpar_target.go`
  - **Objective:** Extract target discovery and dirty repository collection from `cli/cmdcpar/cpar.go`.
  - **Status:** Completed
- [x] **Subtask-29:** [subtask-29.md](../subtasks/292-improvements-planning-and-execution/subtask-29.md)
  - **Target:** `cli/cmdstatus/status_resolve.go`
  - **Objective:** Extract scope resolution and DB/JSON record loading from `cli/cmdstatus/status.go`.
  - **Status:** Completed
- [x] **Subtask-30:** [subtask-30.md](../subtasks/292-improvements-planning-and-execution/subtask-30.md)
  - **Target:** `cli/cmd/see_ops.go`
  - **Objective:** Extract operational handlers (`RunRepoManageUI`, `runSeeCommitPending`, `runSeeHistory`) from `cli/cmd/cmd_see.go`.
  - **Status:** Completed

> **Closed 2026-10-23 (task 190, subtask 06 drain):** plan complete per its own status record and the spec 02 audit (substantively complete with matching commits); moved from pending/ to completed/ so the pending folder reflects only genuinely open work.
