# Master Plan 295: Repository-Wide Wrapped Error & AppError Remediation (CG Execute)

## Status: PENDING
- **Plan ID:** 295
- **Slug:** `295-cg-execute-wrapped-error-remediation`
- **Component Spec Reference:** [02-spec/03-error-manage/](../../../02-spec/03-error-manage/)
- **Prompt Reference:** [01-prompts/15-cg-execute/02-error-management.md](../../../01-prompts/15-cg-execute/02-error-management.md)
- **Scope:** Complete elimination of unwrapped Go errors, uncontextualized `fmt.Errorf` calls, and raw `error` returns across remaining `cli/` packages.
- **Execution Mandate:** Zero-Build, Zero-Test, Strictly Relative Lowercase Git Paths, Bounded 5–8 File Micro-Batches

---

## 1. Executive Summary & Problem Context

During recent code audits and coding guideline reviews (Rule Family 02: Error Management), several packages across `cli/` were discovered to contain non-compliant error management practices:
1. Functions returning raw, untyped `error` rather than strongly-typed `*appfault.AppError`.
2. Direct calls to `fmt.Errorf(...)` without `appfault.Wrap` or `appfault.NewSimple`.
3. Leaf functions invoking exit handlers (`cliexit.HandleError`) or `os.Exit` directly and returning `nil`, thereby violating leaf error sovereignty and bypassing database command audit logging (`finishCommandAudit`).
4. Swallowed errors (`_ = ...`) in resource cleanup and database operations.

This master plan structures the repository-wide remediation across discrete, bounded waves.

---

## 2. Target Subsystems & Wave Breakdown

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Wave 1: Core CLI Root Dispatchers & Tooling Commands (cli/cmd/)             │
│ Focus: root.go, rootutility.go, command error audit capture                 │
├─────────────────────────────────────────────────────────────────────────────┤
│ Wave 2: Storage Engine & Split-DB Repositories (cli/store/)                 │
│ Focus: sqlite_*.go, query filters, repository error returns                 │
├─────────────────────────────────────────────────────────────────────────────┤
│ Wave 3: Pipeline Execution & Diagnostics (cli/cmdpe/, cli/cmdpipeline/)     │
│ Focus: pe_*.go, pipeline database queries, history renderers                │
├─────────────────────────────────────────────────────────────────────────────┤
│ Wave 4: Pull Coordination & Scanner Filtering (cli/cmdpull/, cli/cmdexclude/)│
│ Focus: pull coordination, exclusion rules, target resolvers                 │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Wave Execution Details

### Wave 1: Core CLI Root Dispatchers & Tooling (`cli/cmd/`)
- **Target Files:**
  - `cli/cmd/root.go`
  - `cli/cmd/rootutility.go`
  - `cli/cmd/root_helpers.go`
  - `cli/cmd/root_dispatch.go`
- **Objectives:**
  - Ensure all command dispatchers capture `*appfault.AppError` and pass into `finishCommandAudit`.
  - Replace raw `fmt.Errorf` with `appfault.NewSimple` or `appfault.Wrap`.
  - Enforce strongly-typed exit codes ending in `Type` (`constants.ExitGeneralError`, `constants.ExitUsageError`).

### Wave 2: Storage Engine & Split-DB Repositories (`cli/store/`)
- **Target Files:**
  - `cli/store/sqlite_config.go`
  - `cli/store/sqlite_db_helpers.go`
  - `cli/store/tasks_split_db.go`
  - `cli/store/agent_task_sqlite.go`
  - `cli/store/internal_error_db.go`
- **Objectives:**
  - Replace all generic `error` return signatures with `*appfault.AppError`.
  - Eliminate all swallowed `_ = rows.Close()` and `_ = db.Close()` errors via structured error logging.
  - Wrap all SQLite driver failures with `appfault.WrapSimple` preserving root cause.

### Wave 3: Pipeline Execution & Diagnostics (`cli/cmdpe/`, `cli/cmdpipeline/`)
- **Target Files:**
  - `cli/cmdpe/pe_run.go`
  - `cli/cmdpe/pe_diagnostics.go`
  - `cli/cmdpipeline/pipeline_db.go`
  - `cli/cmdpipeline/pipeline_history.go`
- **Objectives:**
  - Convert unwrapped pipeline errors to structured `appfault.AppError` with error codes.
  - Guarantee leaf error sovereignty: leaf diagnostic extractors return errors directly without calling `os.Exit`.

### Wave 4: Pull Coordination & Exclusion Rules (`cli/cmdpull/`, `cli/cmdexclude/`)
- **Target Files:**
  - `cli/cmdpull/pull_coordinator.go`
  - `cli/cmdpull/pull_target.go`
  - `cli/cmdexclude/exclude_rules.go`
  - `cli/cmdexclude/exclude_filter.go`
- **Objectives:**
  - Eliminate raw string error formatting in pull and exclusion subsystems.
  - Return `*appfault.AppError` through target resolvers.

---

## 4. Verification & Invariants

1. **Leaf Error Sovereignty:** Leaf functions must return errors directly; only root dispatchers invoke exit handlers.
2. **Micro-Batching:** 5–8 files per wave to avoid context blowup.
3. **Zero Builds & Tests:** Static AST analysis and FastGate verification only (`python 03-ai-scripts/50-fastgate.py`).
4. **Relative Paths:** Strictly relative lowercase Git paths across all documentation.
