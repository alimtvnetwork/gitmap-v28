# Master Plan 290: GitMap 30 Improvements Planning and Execution (Batch 2)

## Status: COMPLETED
- **Plan ID:** 290
- **Slug:** `290-improvements-planning-and-execution`
- **Master Spec Reference:** [01-architecture-spec.md](../../../02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md)
- **Component Spec Reference:** [02-component-spec.md](../../../02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md)
- **Migration & Verification Spec Reference:** [03-migration-and-verification-spec.md](../../../02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md)
- **Scope:** 30 High-Impact Architectural, UI, Guidelines, and Operational Safety Improvements (Batch 2)
- **Execution Mandate:** Zero-Build, Zero-Test, Zero Git Commands, Strictly Relative Lowercase Git Paths

---

## 1. Executive Summary & Wave Organization

This master plan coordinates the specification and staged execution of Batch 2 of the 30 discrete architectural improvements for GitMap. The work is structured into three execution waves partitioned by functional domain and file boundary isolation to ensure zero concurrent write collisions:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│ Wave 1: Error Sovereignty, Codebase Guidelines & Types (Subtasks 01–10)      │
│ Spec: 02-spec/21-app/290-improvements-planning-and-execution/01-architecture │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 2: Terminal UI Modernization, Lipgloss & Visual Rhythm (Subtasks 11–20) │
│ Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component    │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 3: System UX, Operational Safety & DB Introspection (Subtasks 21–30)    │
│ Spec: 02-spec/21-app/290-improvements-planning-and-execution/03-migration    │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Complete 30-Subtask Registry & Traceability Matrix

### Wave 1: Error Sovereignty, Codebase Guidelines & Type Hardening (Subtasks 01–10)

- [x] **Subtask-01:** [subtask-01.md](../subtasks/290-improvements-planning-and-execution/subtask-01.md)
  - **Target:** `cli/cmdcd/cdops.go`
  - **Objective:** Convert void/error-ignoring repo matching helpers (`findCDRecords`) to return `([]model.ScanRecord, error)` and bubble errors.
  - **Status:** Completed
- [x] **Subtask-02:** [subtask-02.md](../subtasks/290-improvements-planning-and-execution/subtask-02.md)
  - **Target:** `cli/cmd/alias.go`
  - **Objective:** Eliminate swallowed DB errors in `dispatchAlias` by propagating errors from `runAliasRemove`, `runAliasList`, `runAliasShow`, and `runAliasSuggest`.
  - **Status:** Completed
- [x] **Subtask-03:** [subtask-03.md](../subtasks/290-improvements-planning-and-execution/subtask-03.md)
  - **Target:** `cli/cmdconfig/config_export.go`
  - **Objective:** Eradicate `cliexit.Reportf` calls returning `nil` in leaf functions; return `*apperror.AppError` directly.
  - **Status:** Completed
- [x] **Subtask-04:** [subtask-04.md](../subtasks/290-improvements-planning-and-execution/subtask-04.md)
  - **Target:** `cli/cmddb/db_doctor.go`
  - **Objective:** Refactor void/swallowed doctor checks (`dirDoctorSize`, `statDoctorItem`) to propagate directory scanning and I/O errors.
  - **Status:** Completed
- [x] **Subtask-05:** [subtask-05.md](../subtasks/290-improvements-planning-and-execution/subtask-05.md)
  - **Target:** `cli/cmdpipeline/pipeline_history.go`
  - **Objective:** Check and propagate clipboard (`writeClipboard`) and stream write errors rather than silently ignoring.
  - **Status:** Completed
- [x] **Subtask-06:** [subtask-06.md](../subtasks/290-improvements-planning-and-execution/subtask-06.md)
  - **Target:** `cli/cmdpipeline/pipeline_query.go`, `cli/cmdpipeline/pipeline_segments.go`
  - **Objective:** Handle JSON unmarshal errors without silent swallowing (`parseGhJobsJSON`, `parseOpenPRsCount`).
  - **Status:** Completed
- [x] **Subtask-07:** [subtask-07.md](../subtasks/290-improvements-planning-and-execution/subtask-07.md)
  - **Target:** `cli/gitutil/git_cmd_exec.go`
  - **Objective:** Add `GetRepoHeadCommitSha(repoPath string) (string, error)` helper with direct error return.
  - **Status:** Completed
- [x] **Subtask-08:** [subtask-08.md](../subtasks/290-improvements-planning-and-execution/subtask-08.md)
  - **Target:** `cli/gitutil/git_cmd_exec.go`
  - **Objective:** Add `GetRepoRemoteURL(repoPath, remoteName string) (string, error)` helper with direct error return.
  - **Status:** Completed
- [x] **Subtask-09:** [subtask-09.md](../subtasks/290-improvements-planning-and-execution/subtask-09.md)
  - **Target:** `cli/constants/constants_git.go`
  - **Objective:** Add strongly typed constants ending in `Type` for `GitSymbolicRef`, `GitRemoteOrigin`, and Git ref subcommands.
  - **Status:** Completed
- [x] **Subtask-10:** [subtask-10.md](../subtasks/290-improvements-planning-and-execution/subtask-10.md)
  - **Target:** `cli/constants/constants_store.go`
  - **Objective:** Add strongly typed table constants `TableRepoVersionHistory = "RepoVersionHistory"` and `TableDetectedProject = "DetectedProject"`.
  - **Status:** Completed

---

### Wave 2: Terminal UI Modernization, Lipgloss & Visual Rhythm (Subtasks 11–20)

- [x] **Subtask-11:** [subtask-11.md](../subtasks/290-improvements-planning-and-execution/subtask-11.md)
  - **Target:** `cli/termout/themes.go`
  - **Objective:** Nord, Dracula, and Solarized Truecolor Palette Tokens (`NordPolarNight`, `DraculaBackground`, etc.).
  - **Status:** Completed
- [x] **Subtask-12:** [subtask-12.md](../subtasks/290-improvements-planning-and-execution/subtask-12.md)
  - **Target:** `cli/termout/theme_lipgloss.go`
  - **Objective:** Lipgloss High-Contrast Adaptive Borders and Box-Drawing Primitives.
  - **Status:** Completed
- [x] **Subtask-13:** [subtask-13.md](../subtasks/290-improvements-planning-and-execution/subtask-13.md)
  - **Target:** `cli/termout/termhelp_types.go`
  - **Objective:** Command Help Status Badges and Glyph Metadata (`Badge string`, `BadgeColor string`).
  - **Status:** Completed
- [x] **Subtask-14:** [subtask-14.md](../subtasks/290-improvements-planning-and-execution/subtask-14.md)
  - **Target:** `cli/termout/render.go`
  - **Objective:** Contextual Section Headers and Card Box Renderers (`RenderCardBox`).
  - **Status:** Completed
- [x] **Subtask-15:** [subtask-15.md](../subtasks/290-improvements-planning-and-execution/subtask-15.md)
  - **Target:** `cli/termtable/table.go`
  - **Objective:** Terminal Table Cell Padding and Visual Column Auto-Sizing with East Asian Width.
  - **Status:** Completed
- [x] **Subtask-16:** [subtask-16.md](../subtasks/290-improvements-planning-and-execution/subtask-16.md)
  - **Target:** `cli/termtable/style.go`
  - **Objective:** Unified Table Border Styles (Rounded, Double, Minimalist ANSI Frames).
  - **Status:** Completed
- [x] **Subtask-17:** [subtask-17.md](../subtasks/290-improvements-planning-and-execution/subtask-17.md)
  - **Target:** `cli/cmdpull/pull_render_helpers.go`
  - **Objective:** Pull Summary Status Table and Visual Tree Glyphs Alignment.
  - **Status:** Completed
- [x] **Subtask-18:** [subtask-18.md](../subtasks/290-improvements-planning-and-execution/subtask-18.md)
  - **Target:** `cli/cmdpull/pull_summary_remediation.go`
  - **Objective:** Pull Failure Remediation Interactive Menu Styling and Modern Tokyo Night/Nord Badges.
  - **Status:** Completed
- [x] **Subtask-19:** [subtask-19.md](../subtasks/290-improvements-planning-and-execution/subtask-19.md)
  - **Target:** `cli/heatmap/heatmap_cmd_render.go`
  - **Objective:** Heatmap Execution Depth Visualizer and Action Status Color Tokens.
  - **Status:** Completed
- [x] **Subtask-20:** [subtask-20.md](../subtasks/290-improvements-planning-and-execution/subtask-20.md)
  - **Target:** `cli/tui/prompt.go`
  - **Objective:** Interactive Terminal Confirmation Prompts and Unified Theming Primitives.
  - **Status:** Completed

---

### Wave 3: System UX, Operational Safety & Database Introspection (Subtasks 21–30)

- [x] **Subtask-21:** [subtask-21.md](../subtasks/290-improvements-planning-and-execution/subtask-21.md)
  - **Target:** `cli/store/repo.go`
  - **Objective:** Add `ListReposByTransport(transport string) ([]model.ScanRecord, error)` to query repositories by transport protocol.
  - **Status:** Completed
- [x] **Subtask-22:** [subtask-22.md](../subtasks/290-improvements-planning-and-execution/subtask-22.md)
  - **Target:** `cli/store/repo.go`
  - **Objective:** Add `CountRepos() (int, error)` fast repository counter executing SQLite `COUNT(*)`.
  - **Status:** Completed
- [x] **Subtask-23:** [subtask-23.md](../subtasks/290-improvements-planning-and-execution/subtask-23.md)
  - **Target:** `cli/store/heatmap_store.go`
  - **Objective:** Update `populateTopTouchedFiles` with SQL SUM aggregation for `ReadCount`, `EditCount`, and `CreatedCount`.
  - **Status:** Completed
- [x] **Subtask-24:** [subtask-24.md](../subtasks/290-improvements-planning-and-execution/subtask-24.md)
  - **Target:** `cli/heatmap/heatmap_tracker.go`
  - **Objective:** Create `cli/heatmap/heatmap_tracker.go` and implement `TrackFileTouch(filePath string, touchType store.FileTouchType) error`.
  - **Status:** Completed
- [x] **Subtask-25:** [subtask-25.md](../subtasks/290-improvements-planning-and-execution/subtask-25.md)
  - **Target:** `cli/cmdagent/agent_ui_logs.go`
  - **Objective:** Support `"touchType"` filter in `appendLogFilterConditions` and `parseLogFilter`.
  - **Status:** Completed
- [x] **Subtask-26:** [subtask-26.md](../subtasks/290-improvements-planning-and-execution/subtask-26.md)
  - **Target:** `cli/gitutil/git_cmd_exec.go`
  - **Objective:** Add `GetRepoTrackingRemote(repoPath string) (string, error)` via `git rev-parse --abbrev-ref @{u}`.
  - **Status:** Completed
- [x] **Subtask-27:** [subtask-27.md](../subtasks/290-improvements-planning-and-execution/subtask-27.md)
  - **Target:** `cli/gitutil/git_cmd_exec.go`
  - **Objective:** Add `HasUnpushedCommits(repoPath string) (bool, error)` via `git cherry -v` with positive boolean `hasCommits`.
  - **Status:** Completed
- [x] **Subtask-28:** [subtask-28.md](../subtasks/290-improvements-planning-and-execution/subtask-28.md)
  - **Target:** `cli/cmdfixauth/fixauth.go`
  - **Objective:** Fix negative boolean naming (`isNonGitRepoCWD` -> `isGitRepo`), rename `keyExists` -> `hasKey`, and return errors directly from leaf functions.
  - **Status:** Completed
- [x] **Subtask-29:** [subtask-29.md](../subtasks/290-improvements-planning-and-execution/subtask-29.md)
  - **Target:** `cli/cmdpull/pull_missing_remediate.go`
  - **Objective:** Enforce vertical whitespace (blank line after `}`, blank line before `return`) and flat guard clauses.
  - **Status:** Completed
- [x] **Subtask-30:** [subtask-30.md](../subtasks/290-improvements-planning-and-execution/subtask-30.md)
  - **Target:** `cli/cmdpe/pe_triage.go`
  - **Objective:** Decompose oversized functions (`RunPeTriage`) and flatten nested loops in `worstFailureDetail`.
  - **Status:** Completed

> **Closed 2026-10-23 (task 190, subtask 06 drain):** plan complete per its own status record and the spec 02 audit (substantively complete with matching commits); moved from pending/ to completed/ so the pending folder reflects only genuinely open work.
