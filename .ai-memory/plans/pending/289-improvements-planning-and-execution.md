# Master Plan 289: GitMap 30 Improvements Planning and Execution

## Status: COMPLETED
- **Plan ID:** 289
- **Slug:** `289-improvements-planning-and-execution`
- **Master Spec Reference:** [01-architecture-spec.md](../../../02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md)
- **Component Spec Reference:** [02-component-spec.md](../../../02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md)
- **Migration & Verification Spec Reference:** [03-migration-and-verification-spec.md](../../../02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md)
- **Scope:** 30 High-Impact Architectural, UI, Guidelines, and Operational Safety Improvements
- **Execution Mandate:** Zero-Build, Zero-Test, Zero Git Commands, Strictly Relative Lowercase Git Paths

---

## 1. Executive Summary & Wave Organization

This master plan coordinates the specification and staged execution of 30 discrete architectural improvements for GitMap. The work is structured into three execution waves partitioned by functional domain and file boundary isolation to ensure zero concurrent write collisions:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│ Wave 1: Error Sovereignty, Codebase Guidelines & Types (Subtasks 01–10)      │
│ Spec: 02-spec/21-app/289-improvements-planning-and-execution/01-architecture │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 2: Terminal UI, Lipgloss Theming & Modern Renderers (Subtasks 11–20)    │
│ Spec: 02-spec/21-app/289-improvements-planning-and-execution/02-component    │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 3: System UX, Operational Safety & DB Introspection (Subtasks 21–30)    │
│ Spec: 02-spec/21-app/289-improvements-planning-and-execution/03-migration    │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Complete 30-Subtask Registry & Traceability Matrix

### Wave 1: Error Sovereignty, Codebase Guidelines & Type Hardening (Subtasks 01–10)

- [x] **Subtask-01:** [subtask-01.md](../subtasks/289-improvements-planning-and-execution/subtask-01.md)
  - **Target:** `cli/cmdpull/pull_remediation_ssh.go`
  - **Objective:** Convert void `updateRepoTransportInDB` to return `error` and propagate upstream.
  - **Status:** COMPLETED
- [x] **Subtask-02:** [subtask-02.md](../subtasks/289-improvements-planning-and-execution/subtask-02.md)
  - **Target:** `cli/cmdpull/pull_remediation_ssh.go`
  - **Objective:** Convert void `clearResolvedPullError` to return `error` and propagate upstream.
  - **Status:** COMPLETED
- [x] **Subtask-03:** [subtask-03.md](../subtasks/289-improvements-planning-and-execution/subtask-03.md)
  - **Target:** `cli/cmdpull/pull_missing_remediate.go`
  - **Objective:** Eliminate swallowed error in `remediateMissingRepoTransport` on transport persistence.
  - **Status:** COMPLETED
- [x] **Subtask-04:** [subtask-04.md](../subtasks/289-improvements-planning-and-execution/subtask-04.md)
  - **Target:** `cli/cmdpe/pe_gh.go`
  - **Objective:** Convert void helper functions that ignore errors to return `*appfault.AppError`.
  - **Status:** COMPLETED
- [x] **Subtask-05:** [subtask-05.md](../subtasks/289-improvements-planning-and-execution/subtask-05.md)
  - **Target:** `cli/cmdhealth/health_checks.go`
  - **Objective:** Replace raw error discard with structured `appfault.Wrap`.
  - **Status:** COMPLETED
- [x] **Subtask-06:** [subtask-06.md](../subtasks/289-improvements-planning-and-execution/subtask-06.md)
  - **Target:** `cli/cmdpipeline/pipeline_target_resolve.go`
  - **Objective:** Convert void directory slug resolution to return typed Result wrapper.
  - **Status:** COMPLETED
- [x] **Subtask-07:** [subtask-07.md](../subtasks/289-improvements-planning-and-execution/subtask-07.md)
  - **Target:** `cli/store/repo_transport.go`
  - **Objective:** Enforce positive boolean `hasUpdated` instead of negative checks.
  - **Status:** COMPLETED
- [x] **Subtask-08:** [subtask-08.md](../subtasks/289-improvements-planning-and-execution/subtask-08.md)
  - **Target:** `cli/gitutil/git_cmd_exec.go`
  - **Objective:** Ensure zero nested ifs in `buildGitDirArgs` via early guard returns.
  - **Status:** COMPLETED
- [x] **Subtask-09:** [subtask-09.md](../subtasks/289-improvements-planning-and-execution/subtask-09.md)
  - **Target:** `cli/constants/constants_git.go`
  - **Objective:** Centralize remaining raw URL protocol literals into typed `GitURLPrefixType`.
  - **Status:** COMPLETED
- [x] **Subtask-10:** [subtask-10.md](../subtasks/289-improvements-planning-and-execution/subtask-10.md)
  - **Target:** `cli/constants/constants_store.go`
  - **Objective:** Add strongly typed query enums ending with `Type` suffix for repo transport lookups.
  - **Status:** COMPLETED

---

### Wave 2: Terminal UI, Lipgloss Theming & Modern Renderers (Subtasks 11–20)

- [x] **Subtask-11:** [subtask-11.md](../subtasks/289-improvements-planning-and-execution/subtask-11.md)
  - **Target:** `cli/termout/themes.go`
  - **Objective:** Tokyo Night Truecolor Palette Token Constants (`TokyoBase`, `TokyoPrimary`, etc.).
  - **Status:** COMPLETED
- [x] **Subtask-12:** [subtask-12.md](../subtasks/289-improvements-planning-and-execution/subtask-12.md)
  - **Target:** `cli/termout/theme_lipgloss.go`
  - **Objective:** Lipgloss TTY-Safe Styling Primitives (`SafeTokyoStyle`, `SafeAdaptiveStyle`).
  - **Status:** COMPLETED
- [x] **Subtask-13:** [subtask-13.md](../subtasks/289-improvements-planning-and-execution/subtask-13.md)
  - **Target:** `cli/termout/termhelp_types.go`
  - **Objective:** Sectional Badge Support in CLI Help System (`Badge string`, `BadgeColor string`).
  - **Status:** COMPLETED
- [x] **Subtask-14:** [subtask-14.md](../subtasks/289-improvements-planning-and-execution/subtask-14.md)
  - **Target:** `cli/termout/render.go`
  - **Objective:** Bracketed Help Section Badge Renderer (`RenderHelpSectionBadge`).
  - **Status:** COMPLETED
- [x] **Subtask-15:** [subtask-15.md](../subtasks/289-improvements-planning-and-execution/subtask-15.md)
  - **Target:** `cli/termpad/termpad.go`
  - **Objective:** Zero-Allocation Right-Padding Buffer Pool (`PadRightMax` using `sync.Pool`).
  - **Status:** COMPLETED
- [x] **Subtask-16:** [subtask-16.md](../subtasks/289-improvements-planning-and-execution/subtask-16.md)
  - **Target:** `cli/termtable/table.go`
  - **Objective:** Terminal Table Truncation with Visual String Width (`runewidth.StringWidth`).
  - **Status:** COMPLETED
- [x] **Subtask-17:** [subtask-17.md](../subtasks/289-improvements-planning-and-execution/subtask-17.md)
  - **Target:** `cli/cmdcg/cg_version.go`
  - **Objective:** Modernized Coding Guideline Version Visualizer with Tokyo Night badges.
  - **Status:** COMPLETED
- [x] **Subtask-18:** [subtask-18.md](../subtasks/289-improvements-planning-and-execution/subtask-18.md)
  - **Target:** `cli/cmdcg/cg_help.go`
  - **Objective:** Coding Guideline Help Banner & Section Badging alignment.
  - **Status:** COMPLETED
- [x] **Subtask-19:** [subtask-19.md](../subtasks/289-improvements-planning-and-execution/subtask-19.md)
  - **Target:** `cli/heatmap/heatmap_cmd_render.go`
  - **Objective:** Agent Heatmap Tree Visualizer Tokyo Night Styling.
  - **Status:** COMPLETED
- [x] **Subtask-20:** [subtask-20.md](../subtasks/289-improvements-planning-and-execution/subtask-20.md)
  - **Target:** `cli/cmdpull/pull_efficient_render_helpers.go`
  - **Objective:** Pull Renderer Theme Token Standardization using `termout.RoleSGR`.
  - **Status:** COMPLETED

---

### Wave 3: System UX, Operational Safety & Database Introspection (Subtasks 21–30)

- [x] **Subtask-21:** [subtask-21.md](../subtasks/289-improvements-planning-and-execution/subtask-21.md)
  - **Target:** `cli/store/repo.go`
  - **Objective:** Add `FindOneBySlug(slug string) (*model.ScanRecord, error)` to query single repo by slug.
  - **Status:** COMPLETED
- [x] **Subtask-22:** [subtask-22.md](../subtasks/289-improvements-planning-and-execution/subtask-22.md)
  - **Target:** `cli/store/store.go`
  - **Objective:** Store `path string` on `DB struct` and add `GetDatabasePath() string` getter.
  - **Status:** COMPLETED
- [x] **Subtask-23:** [subtask-23.md](../subtasks/289-improvements-planning-and-execution/subtask-23.md)
  - **Target:** `cli/store/heatmap_schema.go`
  - **Objective:** Extend `FileTouchCount` with `ReadCount`, `EditCount`, and `CreatedCount`.
  - **Status:** COMPLETED
- [x] **Subtask-24:** [subtask-24.md](../subtasks/289-improvements-planning-and-execution/subtask-24.md)
  - **Target:** `cli/store/heatmap_schema.go`
  - **Objective:** Define `FileTouchType` enum (`FileTouchRead`, `FileTouchEdit`, `FileTouchCreated`).
  - **Status:** COMPLETED
- [x] **Subtask-25:** [subtask-25.md](../subtasks/289-improvements-planning-and-execution/subtask-25.md)
  - **Target:** `cli/cmdagent/agent_ui_server.go`
  - **Objective:** Expose touch type telemetry (`deriveTouchType`, `touchType`) in agent action logs API.
  - **Status:** COMPLETED
- [x] **Subtask-26:** [subtask-26.md](../subtasks/289-improvements-planning-and-execution/subtask-26.md)
  - **Target:** `cli/gitutil/git_cmd_exec.go`
  - **Objective:** Add `GetRepoCurrentBranch(repoPath string) (string, error)` via `rev-parse`.
  - **Status:** COMPLETED
- [x] **Subtask-27:** [subtask-27.md](../subtasks/289-improvements-planning-and-execution/subtask-27.md)
  - **Target:** `cli/gitutil/git_cmd_exec.go`
  - **Objective:** Add `IsRepoClean(repoPath string) (bool, error)` via `status --porcelain` with positive boolean.
  - **Status:** COMPLETED
- [x] **Subtask-28:** [subtask-28.md](../subtasks/289-improvements-planning-and-execution/subtask-28.md)
  - **Target:** `cli/cmdclone/remotetransport.go`
  - **Objective:** Verify and enforce zero nested ifs with guard clauses and positive booleans (`isSSH`, `isHTTPS`).
  - **Status:** COMPLETED
- [x] **Subtask-29:** [subtask-29.md](../subtasks/289-improvements-planning-and-execution/subtask-29.md)
  - **Target:** `cli/cmdpull/pull_path_heal.go`
  - **Objective:** Verify and enforce zero nested ifs with guard clauses in path healing functions.
  - **Status:** COMPLETED
- [x] **Subtask-30:** [subtask-30.md](../subtasks/289-improvements-planning-and-execution/subtask-30.md)
  - **Target:** `cli/cmdpe/pe_gh.go`
  - **Objective:** Verify and implement graceful handling in `ghJSON` when `gh` CLI is absent.
  - **Status:** COMPLETED

---

## 3. Core Architectural Invariants

Every subtask implementation must strictly enforce:
1. **Rule 1 (Error Return Sovereignty):** Leaf functions return errors directly; zero `os.Exit` or `cliexit.HandleError` in leaf routines.
2. **Rule 2 (Function Limit):** Target <= 8 lines, hard limit <= 15 lines per function.
3. **Rule 3 (Zero Nested Ifs):** Flatten all branching using guard clauses and early returns.
4. **Rule 4 (Positive Boolean Naming):** All booleans must start with `is*` or `has*` only. No negative polarity booleans.
5. **Rule 5 (Enum Standards):** All enum types must end with `Type` suffix.
6. **Rule 6 (Vertical Whitespace):** Single blank line after block closure `}` and before return statements.
7. **Strictly Relative Paths:** Zero absolute filesystem paths or URI schemes in code or documentation.
8. **Zero Git / Zero Compiler Commands:** No git commands, builds, or tests during planning phase.
