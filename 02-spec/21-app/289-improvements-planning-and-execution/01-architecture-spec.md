# Specification 289: GitMap 30 Improvements Architecture & Execution Protocol

- **Sequence:** 289
- **Slug:** `289-improvements-planning-and-execution`
- **File:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:** `02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/01-architecture-spec.md`, `02-spec/21-app/288-create-prompts-for-codebase-improvements/01-architecture-spec.md`
- **Assigned Author:** Spec Author 01

---

## 1. Executive Summary & System Vision

The objective of Specification 289 is to architect, standardize, and govern the implementation of **30 High-Impact Architectural Improvements** for GitMap. These improvements synthesize findings across three canonical improvement pillars established in Specification 288 and `01-prompts/30-improvements/`:

```
                       GitMap 30 Improvements Architecture
                                      │
         ┌────────────────────────────┼────────────────────────────┐
         ▼                            ▼                            ▼
   Pillar I (01-10)             Pillar II (11-20)            Pillar III (21-30)
  Codebase Guidelines,         Terminal UI, Lipgloss         Whole-System UX,
 Error Sovereignty & Types     & Visual Modernization     Automation & Reliability
```

1. **Pillar I: Codebase Guidelines, Error Sovereignty & Type Hardening (Improvements 01–10):**
   - Eliminates void functions that mutate persistent state or perform I/O.
   - Enforces strict leaf error return sovereignty: functions with error return types propagate failures directly (`return err` or `return appfault.Wrap(...)`) rather than invoking CLI exit handlers or swallowing errors.
   - Restructures conditionals with early guard clauses (zero nested ifs).
   - Enforces positive boolean semantics (`is*`, `has*` only; zero negative polarity booleans).
   - Introduces strongly typed constants and enums ending with `Type` suffix.

2. **Pillar II: Terminal UI, Lipgloss & Visual Modernization (Improvements 11–20):**
   - Harmonizes terminal layout, hierarchy, and color tokens across all CLI commands.
   - Migrates hardcoded ANSI escape sequences to declarative Lipgloss styles and adaptive color tokens.
   - Introduces modern tree-branch visual glyphs (`constants.TreeBranch`, `constants.TreeCorner`) and consistent progress indicators.
   - Ensures high visual contrast in both dark and light terminal viewports.

3. **Pillar III: Whole-System UX, Automation & Fleet Reliability (Improvements 21–30):**
   - Elevates multi-node fleet orchestration, SSH delegation, and cluster telemetry.
   - Hardens SQLite Split-DB operations with thread-safe connection pooling, automated recovery, and positive state flags.
   - Provides non-destructive, backwards-compatible improvements that preserve user working trees, stored configurations, and database integrity.

---

## 2. Core Architectural Principles & Invariants

All implementations across the 30 improvements MUST strictly adhere to the following architectural rules:

### 2.1 Error Return Sovereignty (Rule 1)
- **Leaf Functions Return Errors Directly:** Leaf functions performing database writes, filesystem mutations, or subprocess calls MUST return `error` or `*appfault.AppError`. They MUST NEVER call `cliexit.HandleError` or `os.Exit` and then return `nil`.
- **Outer Caller Exits & Audit Capture:** Only the top-level command handler or root dispatcher (`cli/cmd/root.go`) decides how to terminate or display error states. This guarantees that all command failures bubble up to `finishCommandAudit` and are durably recorded in `installation.db` / `errors.db`.
- **Structured Error Wrapping:** Internal package errors must be wrapped using `appfault.Wrap(err, code, msg)` or `apperror.WrapSimple(err, context)` to maintain traceable call context.

### 2.2 Function Length Limit (Rule 2)
- Target: <= 8 lines preferred.
- Maximum: <= 15 lines hard cap.
- Monolithic routines must be decomposed into small, cohesive private helpers with distinct responsibilities.

### 2.3 Zero Nested Ifs (Rule 3)
- Nested conditionals (`if { if { ... } }`) are strictly forbidden.
- All conditional branching must be flattened using guard clauses, early returns, and discrete predicate functions.

### 2.4 Positive Boolean Naming (Rule 4)
- All boolean variables, fields, and functions MUST begin with `is` or `has` only.
- Prefixes like `can`, `should`, `was`, `will` are forbidden.
- Negative booleans (e.g. `isNotUpdated`, `hasNoRemote`, `untracked`) are strictly banned; all logic must use positive polarity (e.g. `hasUpdated`, `hasRemote`, `isTracked`).

### 2.5 Strongly Typed Enums (Rule 5)
- All enum types MUST end with the `Type` suffix (e.g. `CredentialStoreType`, `GitURLPrefixType`, `RepoTransportQueryType`).
- Raw magic strings and integer exit codes are banned. Constants must be declared in `cli/constants/`.

### 2.6 Return Styling & Vertical Whitespace
- Exactly one blank line before `return` statements (unless the function body is a single return line).
- Exactly one blank line after closing brace `}` before the next statement.

---

## 3. Disjoint Subtask Decomposition & Ownership Boundaries

To guarantee multi-agent execution safety and avoid concurrent file write collisions, the 30 improvements are partitioned into three decoupled 10-subtask batches:

```
┌────────────────────────────────────────────────────────────────────────┐
│                      Batch 1: Improvements 01 - 10                     │
│               Assigned: Spec Author 01 / Worker 01                     │
│  Files: cli/cmdpull/, cli/cmdpe/, cli/cmdhealth/, cli/cmdpipeline/,   │
│         cli/store/repo_transport.go, cli/gitutil/git_cmd_exec.go,      │
│         cli/constants/constants_git.go, constants_store.go             │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                      Batch 2: Improvements 11 - 20                     │
│               Assigned: Spec Author 02 / Worker 02                     │
│  Files: cli/termout/, cli/termtable/, cli/cmdpull/pull_render*,        │
│         cli/constants/constants_terminal.go, cli/macro/                │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                      Batch 3: Improvements 21 - 30                     │
│               Assigned: Spec Author 03 / Worker 03                     │
│  Files: cli/cmdssh/, cli/cluster/, cli/store/split_db*,                │
│         cli/scanpipe/, cli/clonenow/                                   │
└────────────────────────────────────────────────────────────────────────┘
```

### Batch 1 Inventory (Subtasks 01 through 10):
- **Subtask 01 (`subtask-01.md`):** `cli/cmdpull/pull_remediation_ssh.go` - convert void `updateRepoTransportInDB` to return `error` and propagate to callers (`remediateCredentialAndConvertToSSH`, `applySSHTransportRemediation`).
- **Subtask 02 (`subtask-02.md`):** `cli/cmdpull/pull_remediation_ssh.go` - convert void `clearResolvedPullError` to return `error` and propagate upstream.
- **Subtask 03 (`subtask-03.md`):** `cli/cmdpull/pull_missing_remediate.go` - eliminate swallowed error in `remediateMissingRepoTransport` when recording updated transport.
- **Subtask 04 (`subtask-04.md`):** `cli/cmdpe/pe_gh.go` - convert void helper functions that ignore errors to return `*appfault.AppError`.
- **Subtask 05 (`subtask-05.md`):** `cli/cmdhealth/health_checks.go` - replace raw error discard with structured `appfault.Wrap`.
- **Subtask 06 (`subtask-06.md`):** `cli/cmdpipeline/pipeline_target_resolve.go` - convert void directory slug resolution to return typed Result wrapper.
- **Subtask 07 (`subtask-07.md`):** `cli/store/repo_transport.go` - enforce positive boolean `hasUpdated` instead of negative checks.
- **Subtask 08 (`subtask-08.md`):** `cli/gitutil/git_cmd_exec.go` - ensure zero nested ifs in `buildGitDirArgs` via early guard returns.
- **Subtask 09 (`subtask-09.md`):** `cli/constants/constants_git.go` - centralize remaining raw URL protocol literals into typed `GitURLPrefixType`.
- **Subtask 10 (`subtask-10.md`):** `cli/constants/constants_store.go` - add strongly typed query enums for repo transport lookups.

---

## 4. Execution Invariants & Non-Negotiables

During the specification authoring and planning phase, the following invariants are strictly binding:

1. **Strictly Relative Git Paths Only:**
   - Every file path referenced in specifications, plans, and documentation MUST be relative to the repository root (e.g. `02-spec/21-app/...`, `.ai-memory/...`, `cli/...`).
   - Absolute filesystem paths (e.g. workspace paths) and file URIs are strictly forbidden.

2. **All Filenames and Paths Strictly Lowercase:**
   - All directories and files created under `02-spec/21-app/289-improvements-planning-and-execution/` and `.ai-memory/plans/subtasks/289-improvements-planning-and-execution/` must be 100% lowercase.

3. **Zero Git Commands:**
   - Agents must not run `git add`, `git commit`, `git push`, `git status`, or `git checkout`. Version control ceremonies are reserved for the release orchestrator.

4. **Zero Compiler Builds and Zero Tests:**
   - No `go build`, `go test`, or binary compilation during this planning and specification phase.

5. **Non-Destructive Backwards Compatibility:**
   - All refactorings must preserve external CLI flag signatures, existing SQLite database schemas, and backwards-compatible JSON outputs.
