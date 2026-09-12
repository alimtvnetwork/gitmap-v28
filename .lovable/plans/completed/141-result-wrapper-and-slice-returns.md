# Plan 141: Result Wrapper Types, Collections & AppError Returns Architecture (Phase 2 - ResultSlice)

Trigger Keywords & Aliases: `cg-result-wrapper`, `cg-apperror-returns`, `cg-execute result-wrapper`, `audit result wrapper`, `fix map return error`, `fix slice return error`, `single return object audit`, `enforce apperror returns`, `enforce result map`, `fix multi-value returns`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)
> **Status:** COMPLETED

---

## 1. Executive Summary & Objective

Autonomously scan, discover, plan, refactor, and verify all Go functions returning multi-value error tuples (such as `(map[K]V, error)`, `([]T, error)`, or `(T, error)`), eliminating raw standard library error returns, replacing them with strongly-typed result wrappers (`ResultMap[K, V]`, `ResultSlice[T]`, `Result[T]`) and structured `*apperror.AppError` returns, guaranteeing a single return object, standardized outer-layer inspection predicates (`IsSuccess`, `IsFailure`, `HasError`, `IsEmptyError`, `IsEmpty`, `Data`, `AppError`, `Fault`, `Get`, `Has`, `Count`), and zero dual-handling across the entire codebase until 100% green without stopping.

Following Plan 138 (which successfully migrated all 13 `ResultMap` map return functions), Plan 141 addresses the `ResultSlice[T]` envelope across four primary domain subsystems:
1. `cli/macro/`: Macro import, parsing, and storage query operations.
2. `cli/pipelinedb/`: Pipeline split database execution, failure queries, and error log retrieval operations.
3. `cli/cmdprompt/`: Target repository discovery, workdir resolution, and child repo resolution.
4. `cli/db/` & `cli/cluster/`: Node path alias listing and cluster argument parsing.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/01-index.md`: Single return type mandate and micro-tasking.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: File and function size limits (functions <= 8–15 lines, files <= 100 coding lines).
- `spec/03-error-manage/01-index.md`: Universal AppError wrapping and error envelopes.
- `spec/03-error-manage/02-error-architecture/02-error-handling-reference.md`: Error handling architecture and Result wrappers.
- `spec/03-error-manage/02-error-architecture/05-response-envelope/05-response-envelope-reference.md`: Response envelope schemas.
- `spec/03-error-manage/02-error-architecture/06-apperror-package/03-go-apperror-linter-spec.md`: Go AppError implementation specifications.
- `.agents/skills/cg-result-wrapper/skill.md`: Canonical skill for result wrappers and AppError returns.

---

## 3. Domain-Specific Rules for Result Wrappers

1. **Rule 1 (Single Return Object Mandate):** Every domain query and collection retrieval function must return a single envelope: `result.ResultSlice[T]` for slices, `result.ResultMap[K, V]` for maps, `result.Result[T]` for scalars.
2. **Rule 2 (Structured AppError Returns):** Raw stdlib `error` returns are strictly banned in domain/service functions; use structured `*apperror.AppError`.
3. **Rule 3 (Outer-Layer Inspection Predicates):** Callers must use `res.IsSuccess()`, `res.IsFailure()`, `res.IsEmpty()`, `res.Data`, and `res.AppError()`.
4. **Rule 4 (No Inverted Predicates):** Never write `if !res.IsSuccess()`; always write `if res.IsFailure()`.
5. **Rule 5 (Function Size & Formatting):** Target <= 8–15 lines per function, mandatory braces, clean blank-line spacing before return.

---

## 4. Exhaustive Violation Ledger

| Id | File | Line | Function | Existing Return | Planned Refactoring | Status |
|:---|:---|:---:|:---|:---|:---|:---:|
| V-01 | `cli/macro/import_sqlite.go` | 15 | `ParseImportSQLite` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-02 | `cli/macro/import_sqlite.go` | 31 | `queryAllMacrosWithSteps` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-03 | `cli/macro/import_sqlite.go` | 48 | `scanMacroRows` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-04 | `cli/macro/import.go` | 107 | `ParseImportJSON` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-05 | `cli/macro/import.go` | 122 | `ParseImportYAML` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-06 | `cli/macro/import.go` | 137 | `ParseImportZIP` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-07 | `cli/macro/import.go` | 148 | `extractMacrosFromZip` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-08 | `cli/macro/import.go` | 204 | `ParseImportFile` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-09 | `cli/macro/storage.go` | 104 | `ListMacros` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-10 | `cli/pipelinedb/pipeline_split_ops.go` | 241 | `QueryRecentRuns` | `([]PipelineRunRecord, error)` | `result.ResultSlice[PipelineRunRecord]` | Fixed |
| V-11 | `cli/pipelinedb/pipeline_split_ops.go` | 224 | `collectRecentRuns` | `([]PipelineRunRecord, error)` | `result.ResultSlice[PipelineRunRecord]` | Fixed |
| V-12 | `cli/pipelinedb/pipeline_split_ops.go` | 282 | `QueryRecentErrorLogs` | `([]PipelineErrorRecord, error)` | `result.ResultSlice[PipelineErrorRecord]` | Fixed |
| V-13 | `cli/pipelinedb/pipeline_split_ops.go` | 263 | `collectRecentErrors` | `([]PipelineErrorRecord, error)` | `result.ResultSlice[PipelineErrorRecord]` | Fixed |
| V-14 | `cli/pipelinedb/pipeline_split_ops.go` | 323 | `QueryDetailedErrors` | `([]PipelineErrorRecord, error)` | `result.ResultSlice[PipelineErrorRecord]` | Fixed |
| V-15 | `cli/pipelinedb/pipeline_split_ops.go` | 336 | `QueryCompactErrors` | `([]PipelineCompactErrorRecord, error)` | `result.ResultSlice[PipelineCompactErrorRecord]` | Fixed |
| V-16 | `cli/pipelinedb/pipeline_split_ops.go` | 308 | `collectCompactErrors` | `([]PipelineCompactErrorRecord, error)` | `result.ResultSlice[PipelineCompactErrorRecord]` | Fixed |
| V-17 | `cli/pipelinedb/pipeline_split_ops.go` | 529 | `QueryCachedErrorRunIds` | `([]uint64, error)` | `result.ResultSlice[uint64]` | Fixed |
| V-18 | `cli/pipelinedb/pipeline_split_ops.go` | 592 | `QueryLastFailedRuns` | `([]PipelineRunRecord, error)` | `result.ResultSlice[PipelineRunRecord]` | Fixed |
| V-19 | `cli/pipelinedb/pipeline_split_ops.go` | 605 | `QueryLastNFailedRuns` | `([]PipelineRunRecord, error)` | `result.ResultSlice[PipelineRunRecord]` | Fixed |
| V-20 | `cli/pipelinedb/pipeline_split_ops.go` | 610 | `QueryErrorLogsByRunId` | `([]PipelineErrorRecord, error)` | `result.ResultSlice[PipelineErrorRecord]` | Fixed |
| V-21 | `cli/pipelinedb/pipeline_split_ops.go` | 622 | `QueryDetailedErrorLogsByRunId` | `([]PipelineErrorRecord, error)` | `result.ResultSlice[PipelineErrorRecord]` | Fixed |
| V-22 | `cli/pipelinedb/pipeline_split_ops.go` | 634 | `QueryCompactErrorLogsByRunId` | `([]PipelineCompactErrorRecord, error)` | `result.ResultSlice[PipelineCompactErrorRecord]` | Fixed |
| V-23 | `cli/cmdprompt/prompt_child_repos.go` | 8 | `DiscoverPromptChildRepos` | `([]string, error)` | `result.ResultSlice[string]` | Fixed |
| V-24 | `cli/cmdprompt/prompt_target_resolver.go` | 13 | `ResolvePromptTarget` | `([]string, error)` | `result.ResultSlice[string]` | Fixed |
| V-25 | `cli/cmdprompt/prompt_workdir_resolver.go` | 8 | `ResolveAllWorkDirPromptTargets` | `([]string, error)` | `result.ResultSlice[string]` | Fixed |
| V-26 | `cli/db/nodepath.go` | 16 | `ListPathAliases` | `([]NodePathAlias, error)` | `result.ResultSlice[NodePathAlias]` | Fixed |
| V-28 | `cli/macro/import.go` | 248 | `filterAndValidateImportTargets` | `([]Macro, error)` | `result.ResultSlice[Macro]` | Fixed |
| V-27 | `cli/cluster/pathalias.go` | 8 | `ParseSetPathAliasArg` | `([]AliasEntry, error)` | `result.ResultSlice[AliasEntry]` | Fixed |

---

## 5. Subtask Decomposition

- **Subtask 01:** `.lovable/plans/subtasks/141-result-wrapper/01-refactor-macro-slice-returns.md`
  - Targets: `cli/macro/import_sqlite.go`, `cli/macro/import.go`, `cli/macro/storage.go`, callers in `cli/cmdmacro/` and `cli/macro/`.
  - Scope: Refactor 9 macro functions to return `result.ResultSlice[Macro]` and modernize all call sites.
- **Subtask 02:** `.lovable/plans/subtasks/141-result-wrapper/02-refactor-pipelinedb-slice-returns.md`
  - Targets: `cli/pipelinedb/pipeline_split_ops.go`, callers in `cli/cmdpipeline/` and `cli/pipelinedb/`.
  - Scope: Refactor 13 pipeline DB functions to return `result.ResultSlice[...]` and modernize all call sites.
- **Subtask 03:** `.lovable/plans/subtasks/141-result-wrapper/03-refactor-prompt-and-nodepath-slice-returns.md`
  - Targets: `cli/cmdprompt/`, `cli/db/nodepath.go`, `cli/cluster/pathalias.go`, and callers in `cli/cmd/` and `cli/cmdcg/`.
  - Scope: Refactor 5 functions to return `result.ResultSlice[...]` and modernize all call sites.
- **Subtask 04:** `.lovable/plans/subtasks/141-result-wrapper/04-upgrade-auditor-and-verify-quality-gates.md`
  - Targets: `03-ai-scripts/35-result-wrapper-auditor.py`, `03-ai-scripts/01-index.md`, quality linters, test inventory lock, milestone consolidation.
  - Scope: Enhance auditor script with ResultSlice checks, run all linters, record modified files, atomic commit & push.

---

## 6. Verification Plan

1. **Go Compilation & Vet:** `go vet ./...` in `cli/` must exit 0 with zero errors.
2. **Result Wrapper Auditor:** `python 03-ai-scripts/35-result-wrapper-auditor.py` must exit 0.
3. **Boolean Guidelines Linter:** `python linter-scripts/check-boolean-guidelines.py` must pass.
4. **Nested If Linter:** `python linter-scripts/check-nested-ifs.py` must pass.
5. **Relative Paths Linter:** `python linter-scripts/check-relative-paths.py` must pass.
6. **Enum Guidelines Linter:** `python linter-scripts/check-enum-guidelines.py` must pass.

---

## 7. Verification & Consolidation Summary

All 28 targeted slice return functions across `cli/macro`, `cli/pipelinedb`, `cli/cmdprompt`, `cli/db`, and `cli/cluster` have been fully migrated to single-object `result.ResultSlice[T]` envelopes and structured `*apperror.AppError` returns. All callers have been modernized to use explicit inspection predicates (`res.IsFailure()`, `res.IsEmpty()`, `res.Data`, `res.AppError()`). The automated auditor `03-ai-scripts/35-result-wrapper-auditor.py` has been upgraded to permanently enforce ResultSlice compliance across these subsystems. All 6 quality gates passed with zero errors.
