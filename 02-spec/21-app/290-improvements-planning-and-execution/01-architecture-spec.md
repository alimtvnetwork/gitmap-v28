# Specification 290: GitMap 30 Improvements Architecture & Execution Protocol (Batch 2)

- **Sequence:** 290
- **Slug:** `290-improvements-planning-and-execution`
- **File:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:** `02-spec/21-app/288-create-prompts-for-codebase-improvements/01-architecture-spec.md`, `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Author:** Spec Author 01

---

## 1. Executive Summary & System Vision

The objective of Specification 290 is to architect, standardize, and govern the implementation of **Batch 2 of the 30 High-Impact Architectural Improvements** for GitMap. Following the successful completion and verification of Specification 289, Batch 2 continues systematically addressing architectural debt, hardening error propagation across leaf functions, elevating terminal user experience, and expanding database schema and fleet diagnostic capabilities.

These improvements span three canonical improvement pillars established in Specification 288 and `01-prompts/30-improvements/`:

```
                       GitMap 30 Improvements Architecture (Batch 2)
                                             │
         ┌───────────────────────────────────┼───────────────────────────────────┐
         ▼                                   ▼                                   ▼
   Pillar I (01-10)                    Pillar II (11-20)                   Pillar III (21-30)
  Error Sovereignty,                 Terminal UI Modernization,             Whole-System UX,
Type Hardening & Void Elimination     Lipgloss & Visual Rhythm         Split-DB & Fleet Reliability
```

1. **Pillar I: Error Return Sovereignty, Type Hardening & Void Elimination (Improvements 01–10):**
   - Eliminates silent error swallowing and unhandled I/O failures across repository matching, aliasing, export workflows, database inspection, and pipeline queries.
   - Enforces strict leaf error return sovereignty: leaf functions performing database operations, filesystem mutations, or subprocess calls MUST return `error` or `*appfault.AppError` directly rather than invoking CLI exit handlers or swallowing errors.
   - Eliminates `cliexit.Reportf` calls that terminate or swallow errors inside leaf helpers and return `nil`.
   - Restructures conditional branching with early guard clauses (zero nested ifs).
   - Enforces positive boolean semantics (`is*`, `has*` only; zero negative polarity booleans).
   - Introduces strongly typed constants and enums ending with the mandatory `Type` suffix.

2. **Pillar II: Terminal UI Modernization, Lipgloss & Visual Rhythm (Improvements 11–20):**
   - Standardizes terminal hierarchy, ANSI escape decoupling, and Lipgloss styling across status views and command outputs.
   - Replaces raw formatted strings with structured terminal blocks, box drawing, and adaptive color tokens.
   - Unifies progress bars, spinners, and branch glyphs (`constants.TreeBranch`, `constants.TreeCorner`).

3. **Pillar III: Whole-System UX, Split-DB & Fleet Reliability (Improvements 21–30):**
   - Fortifies the Split-DB architecture (`gitmap.db`, `installation.db`, `pipeline.db`) with strongly typed table names and transaction safety.
   - Enhances multi-node cluster synchronization, SSH delegation, and background pipeline diagnostics.
   - Guarantees backward compatibility and idempotency across all database migrations and CLI workflows.

---

## 2. Core Architectural Principles & Invariants

All implementations across the 30 improvements MUST strictly adhere to the following architectural rules:

### 2.1 Error Return Sovereignty (Rule 1)
- **Leaf Functions Return Errors Directly:** Leaf functions that declare an `error` return type MUST return the error directly (`return err` or `return appfault.Wrap(...)`). Leaf functions MUST NEVER call exit handlers (`cliexit.HandleError`, `cliexit.Reportf`) or `os.Exit` and then return `nil`.
- **Outer Caller Handles Exits:** Only the top-level orchestrator or root command dispatcher (`cli/cmd/root.go`) receives the error and decides how to terminate, render exit notices, or report.
- **Database Audit Capture:** Errors bubble up to `dispatch` handlers in `cli/cmd/root.go` so they are durably captured by `finishCommandAudit` in `installation.db` and `errors.db`.
- **No Magic Integer Exit Codes:** Never pass raw integers (`1`, `2`) to exit handlers. Always use strongly typed enums ending in `Type` (`ExitCodeType`, `constants.ExitGeneralError`, `constants.ExitUsageError`).

### 2.2 Function Length Limit (Rule 2)
- Target: <= 8 lines preferred.
- Maximum: <= 15 lines hard cap.
- Any routine exceeding 15 lines MUST be decomposed into small, focused private helper functions with single responsibilities.

### 2.3 Zero Nested Ifs (Rule 3)
- Nested conditionals (`if { if { ... } }`) are strictly forbidden.
- All conditional branching must be flattened using guard clauses, early returns, inverted conditions, and discrete predicate helpers.

### 2.4 Positive Boolean Naming (Rule 4)
- All boolean variables, fields, and functions MUST begin with `is` or `has` ONLY.
- Prefixes like `can`, `should`, `was`, `will` are banned.
- Negative booleans (e.g. `isNotUpdated`, `hasNoRemote`, `untracked`) are strictly forbidden; all logic must use positive polarity (e.g. `hasUpdated`, `hasRemote`, `isTracked`).

### 2.5 Strongly Typed Enums (Rule 5)
- All enum types MUST end with the `Type` suffix (e.g. `GitRefCommandType`, `GitRemoteNameType`, `TableNameType`).
- Raw magic strings and raw literals are banned. All constants must be declared in `cli/constants/`.

### 2.6 Return Styling & Vertical Whitespace
- Exactly one blank line before `return` statements (unless the function body is a single return line).
- Exactly one blank line after closing brace `}` before the next statement.

---

## 3. User Non-Negotiables & Execution Safety

During the specification and execution phases of Specification 290:
1. **Strictly Relative Git Paths:** All documentation, plans, specs, and source references MUST use strictly relative lowercase Git paths (e.g. `cli/cmdcd/cdops.go`, `02-spec/21-app/290-improvements-planning-and-execution/`). Absolute filesystem paths and file URI schemes are strictly forbidden.
2. **Strictly Lowercase Paths:** All directory names, filenames, and references must be strictly lowercase.
3. **Zero Git Commands:** AI agents and subagents must not invoke git commands (`git commit`, `git status`, `git add`, etc.) during spec authoring.
4. **Zero Compiler Builds and Zero Tests:** No `go build`, `go test`, or runner scripts shall be executed during this specification phase.

---

## 4. Disjoint Subtask Decomposition & Ownership Boundaries

To guarantee multi-agent execution safety and avoid concurrent file write collisions, the 30 improvements are partitioned into three decoupled 10-subtask batches:

```
┌────────────────────────────────────────────────────────────────────────┐
│                      Batch 1: Improvements 01 - 10                     │
│               Assigned: Spec Author 01 / Worker 01                     │
│  Files: cli/cmdcd/cdops.go, cli/cmd/alias.go,                         │
│         cli/cmdconfig/config_export.go, cli/cmddb/db_doctor.go,        │
│         cli/cmdpipeline/pipeline_history.go, pipeline_query.go,        │
│         cli/cmdpipeline/pipeline_segments.go,                          │
│         cli/gitutil/git_cmd_exec.go,                                   │
│         cli/constants/constants_git.go, constants_store.go             │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                      Batch 2: Improvements 11 - 20                     │
│               Assigned: Spec Author 02 / Worker 02                     │
│  Files: cli/termout/, cli/termtable/, cli/cmdpull/pull_render*,        │
│         cli/constants/constants_terminal.go, cli/tui/                  │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                      Batch 3: Improvements 21 - 30                     │
│               Assigned: Spec Author 03 / Worker 03                     │
│  Files: cli/cluster/, cli/cmdcluster/, cli/cmdnodes/,                  │
│         cli/pipelinedb/, cli/store/, cli/cmddb/                        │
└────────────────────────────────────────────────────────────────────────┘
```

### Detailed Mapping for Improvements 01 - 10 (Spec Author 01):

| Improvement ID | Target File(s) | Primary Architectural Focus |
|---|---|---|
| **Subtask 01** | `cli/cmdcd/cdops.go` | Convert void/error-ignoring repo matching helpers (`findCDRecords`) to return `([]model.ScanRecord, error)` and bubble errors. |
| **Subtask 02** | `cli/cmd/alias.go` | Eliminate swallowed DB errors in `dispatchAlias` by propagating errors from `runAliasRemove`, `runAliasList`, `runAliasShow`, and `runAliasSuggest`. |
| **Subtask 03** | `cli/cmdconfig/config_export.go` | Eradicate `cliexit.Reportf` calls returning `nil` in leaf functions; return `*apperror.AppError` directly. |
| **Subtask 04** | `cli/cmddb/db_doctor.go` | Refactor void/swallowed doctor checks (`dirDoctorSize`, `statDoctorItem`) to propagate directory scanning and I/O errors. |
| **Subtask 05** | `cli/cmdpipeline/pipeline_history.go` | Check and propagate clipboard (`writeClipboard`) and stream write errors rather than silently ignoring. |
| **Subtask 06** | `cli/cmdpipeline/pipeline_query.go`, `pipeline_segments.go` | Handle JSON unmarshal errors without silent swallowing (`parseGhJobsJSON`, `parseOpenPRsCount`). |
| **Subtask 07** | `cli/gitutil/git_cmd_exec.go` | Add `GetRepoHeadCommitSha(repoPath string) (string, error)` helper with direct error return. |
| **Subtask 08** | `cli/gitutil/git_cmd_exec.go` | Add `GetRepoRemoteURL(repoPath, remoteName string) (string, error)` helper with direct error return. |
| **Subtask 09** | `cli/constants/constants_git.go` | Add strongly typed constants ending in `Type` for `GitSymbolicRef`, `GitRemoteOrigin`, and Git ref subcommands. |
| **Subtask 10** | `cli/constants/constants_store.go` | Add strongly typed table constants `TableRepoVersionHistory = "RepoVersionHistory"` and `TableDetectedProject = "DetectedProject"`. |

---

## 5. Verification & Audit Protocol

1. **Static Analysis Compliance:** Each subtask plan enforces function line constraints (<= 15 lines), zero nested ifs, positive boolean prefixes (`is*`, `has*`), and typed enums ending in `Type`.
2. **Review Docket Readiness:** All 10 subtask plan files (`subtask-01.md` through `subtask-10.md`) are placed under `.ai-memory/plans/subtasks/290-improvements-planning-and-execution/` with explicit input/output contracts, concrete code diffs, and verification checklists.
