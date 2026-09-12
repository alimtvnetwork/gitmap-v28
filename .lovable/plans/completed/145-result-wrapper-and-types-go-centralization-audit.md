# Plan 145: Result Wrapper Types, Collections & AppError Returns (Types.go Centralization & Single Reusable Types)

Trigger Keywords & Aliases: `cg-result-wrapper`, `cg-apperror-returns`, `cg-execute result-wrapper`, `audit result wrapper`, `fix map return error`, `fix slice return error`, `single return object audit`, `enforce apperror returns`, `enforce result map`, `fix multi-value returns`, `is-count-other-than`, `has-record`, `is-defined`, `result-wrapper-null-safety`, `pointer-null-safety`, `types-go-single-type`, `types-go-result-reuse`, `centralize-types-go`

> **Prompt Version:** 2.4.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)
> **Status:** COMPLETED

---

## 1. Executive Summary & Objective

Autonomously scan, discover, plan, refactor, and verify all Go functions returning multi-value error tuples (such as `(map[K]V, error)`, `([]T, error)`, or `(T, error)`), eliminating raw standard library error returns, centralizing all domain payload structs and Result type aliases into `types.go` within each package as single reusable types everywhere rather than scattering inline structs or raw generic Result declarations across implementation files, replacing multi-value returns with strongly-typed result wrappers (`ResultMap[K, V]`, `ResultSlice[T]`, `Result[T]`) and structured `*apperror.AppError` returns, guaranteeing a single return object, pointer-attached null safety (`*Result[T]`, `*ResultSlice[T]`, `*ResultMap[K, V]`) with methods attached to pointer receivers (`(r *Result[T])`, `(rs *ResultSlice[T])`, `(rm *ResultMap[K, V])`) that verify `if r == nil` before dereferencing any fields or checking errors, standardized outer-layer inspection predicates (`IsSuccess`, `IsFailure`, `HasError`, `IsEmptyError`, `IsEmpty`, `HasRecord`, `IsDefined`, `IsCountOtherThan`, `Data`, `Items`, `AppError`, `Fault`, `Get`, `Has`, `Count`), eliminating dual-handling, and replacing verbose `if err != nil || len(...) != N` or `IsFailure() || Count() != N` conditions with fluent `if res.IsCountOtherThan(N)` across the entire codebase until 100% green without stopping.

Plan 145 focuses on the `types.go` single reusable type mandate and Result alias architecture across 4 milestones:
1. **Result Package Types Centralization & Affirmative Field Naming (`cli/result/types.go`):** Create `cli/result/types.go` centralizing core container structs (`Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`), core aliases (`Wrap[T]`, `Slice[T]`, `Map[K, V]`), and inspector/verifier interfaces. Fix the non-affirmative field `defined bool` -> `isDefined bool` in `Result[T]`. Fix clumsy checks in `null_safety_test.go` to use `res.IsCountOtherThan(N)`.
2. **Cmdschedule Types Centralization & Reusable Envelopes (`cli/cmdschedule/types.go`):** Create `cli/cmdschedule/types.go` centralizing exported domain structs (`ScheduleExportBundle`, `ScheduleExportOpts`) and single reusable Result envelopes (`ScheduleExportBundleResult`, `ScheduleExportBundleSingleResult`, `SchedulerTaskSliceResult`). Refactor 7 functions in `schedule_export.go` and `schedule_import.go` from raw generics to single reusable types from `types.go`.
3. **Macro & PipelineDB Types Centralization (`cli/macro/types.go`, `cli/pipelinedb/types.go`):** Add single reusable Result aliases (`MacroSliceResult`, `MacroResult`, `MacroStepsMapResult`) to `cli/macro/types.go` and update 10 functions in `import.go`, `import_sqlite.go`, and `storage.go`. Create `cli/pipelinedb/types.go` defining single reusable Result aliases (`PipelineRunSliceResult`, `PipelineErrorSliceResult`, `PipelineCompactErrorSliceResult`, `PipelineRunIdSliceResult`, `PipelineRunIdMapResult`) and update 14 functions in `pipeline_split_ops.go`.
4. **Automated Auditor Verification & Quality Gates:** Update `03-ai-scripts/35-result-wrapper-auditor.py` to enforce `types.go` presence and ban raw generic Result returns in non-types files across enforced subsystems. Verify zero regressions across all quality linters (`go vet`, `check-boolean-guidelines.py`, `check-nested-ifs.py`, `check-relative-paths.py`, `check-enum-guidelines.py`), lock modified files in test inventory, consolidate Plan 145, and push atomically via SSH.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/01-index.md`: Single return type mandate and micro-tasking.
- `spec/02-coding-guidelines/01-cross-language/27-types-folder-convention.md`: Types.go and single type definitions.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: File and function size limits (functions <= 8–15 lines, files <= 100–300 lines).
- `spec/03-error-manage/01-index.md`: Universal AppError wrapping and error envelopes.
- `spec/03-error-manage/02-error-architecture/02-error-handling-reference.md`: Error handling architecture and Result wrappers.
- `spec/03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/04-result-types.md`: Result[T], ResultSlice[T], and ResultMap[K, V] method specifications and pointer null-safety rules.
- `.agents/skills/cg-result-wrapper/skill.md`: Canonical skill for result wrappers and AppError returns.
- `.lovable/coding-guidelines.md`: Master consolidated coding guidelines.

---

## 3. Domain-Specific Rules for Result Wrappers & Types.go Centralization

1. **Rule 1 (Dedicated `types.go` for Single Reusable Types):** Every package managing domain models, payloads, or Result envelopes MUST define them inside a dedicated `types.go` file within the package directory as a single reusable named type. Never declare unexported domain structs or raw generic Result returns inline across implementation files.
2. **Rule 2 (Affirmative Boolean Field Naming):** Struct fields and local variables MUST carry affirmative prefixes (`is*` or `has*`). In Result wrappers, definition state MUST be named `isDefined bool` (TOTAL BAN on bare `defined bool`).
3. **Rule 3 (Pointer Receivers for All Inspection Methods):** All inspection methods on `Result[T]`, `ResultSlice[T]`, and `ResultMap[K, V]` must be attached to pointer receivers (`*Result[T]`, `*ResultSlice[T]`, `*ResultMap[K, V]`).
4. **Rule 4 (Immediate Nil Guard):** Every pointer receiver method must guard `if r == nil` on line 1 and return safe canonical defaults.
5. **Rule 5 (The 4 Core Predicate Methods):** Enforce `IsCountOtherThan(number int)`, `IsEmpty()`, `HasRecord()`/`HasRecords()`, and `IsDefined()`.
6. **Rule 6 (Single Return Object Mandate):** Functions must return a single envelope object or `*apperror.AppError` for pure side effects. Multi-value error tuples `(T, error)` are banned.

---

## 4. Exhaustive Violation Ledger

| Id | File | Line | Identifier | Issue | Target Refactoring | Status |
|:---|:---|:---:|:---|---|---|:---:|
| V-01 | `cli/result/types.go` | 1 | `types.go` missing | Core container structs scattered in implementation files | Create `cli/result/types.go` declaring `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]` | Completed |
| V-02 | `cli/result/result.go` | 14 | `Result.defined` | Non-affirmative field `defined bool` | Rename to `isDefined bool` | Completed |
| V-03 | `cli/result/result.go` | 115 | `r.defined` | Field reference in `Count()` | Update to `r.isDefined` | Completed |
| V-04 | `cli/result/result.go` | 202 | `defined: true` | Field initialization in `Ok()` | Update to `isDefined: true` | Completed |
| V-05 | `cli/result/types.go` | N/A | Type Aliases | Missing single reusable aliases | Define `Wrap[T] = Result[T]`, `Slice[T] = ResultSlice[T]`, `Map[K, V] = ResultMap[K, V]` | Completed |
| V-06 | `cli/result/null_safety_test.go` | 292 | Call site | Clumsy check `sliceOk.IsFailure() \|\| sliceOk.Count() != 2` | Replace with `sliceOk.IsCountOtherThan(2)` | Completed |
| V-07 | `cli/result/null_safety_test.go` | 304 | Call site | Clumsy check `sliceEmpty.IsFailure() \|\| sliceEmpty.Count() != 0` | Replace with `sliceEmpty.IsCountOtherThan(0)` | Completed |
| V-08 | `cli/result/null_safety_test.go` | 313 | Call site | Clumsy check `mapOk.IsFailure() \|\| mapOk.Count() != 1` | Replace with `mapOk.IsCountOtherThan(1)` | Completed |
| V-09 | `cli/cmdschedule/types.go` | 1 | `types.go` missing | Domain models and Result envelopes scattered inline | Create `cli/cmdschedule/types.go` with domain models & Result aliases | Completed |
| V-10 | `cli/cmdschedule/schedule_export.go` | 22 | `scheduleExportBundle` | Unexported struct declared inline | Move to `types.go` as exported `ScheduleExportBundle` | Completed |
| V-11 | `cli/cmdschedule/schedule_export.go` | 27 | `scheduleExportOpts` | Unexported struct declared inline | Move to `types.go` as exported `ScheduleExportOpts` | Completed |
| V-12 | `cli/cmdschedule/types.go` | N/A | Result Alias | Missing single reusable Result alias | Define `ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]` | Completed |
| V-13 | `cli/cmdschedule/types.go` | N/A | Result Alias | Missing single reusable single-bundle alias | Define `ScheduleExportBundleSingleResult = result.Result[ScheduleExportBundle]` | Completed |
| V-14 | `cli/cmdschedule/types.go` | N/A | Result Alias | Missing single reusable task slice alias | Define `SchedulerTaskSliceResult = result.ResultSlice[store.SchedulerTask]` | Completed |
| V-15 | `cli/cmdschedule/schedule_export.go` | 113 | `collectExportBundles` | Returns raw generic `result.ResultSlice[scheduleExportBundle]` | Update signature to return `ScheduleExportBundleResult` | Completed |
| V-16 | `cli/cmdschedule/schedule_import.go` | 49 | `parseImportFileBundles` | Returns raw generic `result.ResultSlice[scheduleExportBundle]` | Update signature to return `ScheduleExportBundleResult` | Completed |
| V-17 | `cli/cmdschedule/schedule_import.go` | 63 | `parseImportJSON` | Returns raw generic `result.ResultSlice[scheduleExportBundle]` | Update signature to return `ScheduleExportBundleResult` | Completed |
| V-18 | `cli/cmdschedule/schedule_import.go` | 82 | `parseImportYAML` | Returns raw generic `result.ResultSlice[scheduleExportBundle]` | Update signature to return `ScheduleExportBundleResult` | Completed |
| V-19 | `cli/cmdschedule/schedule_import.go` | 101 | `parseImportSQLite` | Returns raw generic `result.ResultSlice[scheduleExportBundle]` | Update signature to return `ScheduleExportBundleResult` | Completed |
| V-20 | `cli/cmdschedule/schedule_import.go` | 127 | `queryImportTasksFromDB` | Returns raw generic `result.ResultSlice[store.SchedulerTask]` | Update signature to return `SchedulerTaskSliceResult` | Completed |
| V-21 | `cli/cmdschedule/schedule_import.go` | 183 | `parseImportZIP` | Returns raw generic `result.ResultSlice[scheduleExportBundle]` | Update signature to return `ScheduleExportBundleResult` | Completed |
| V-22 | `cli/cmdschedule/schedule_export.go` | 35 | `runScheduleExport` | Uses raw inline types | Modernize with `ScheduleExportOpts` & `ScheduleExportBundleResult` | Completed |
| V-23 | `cli/cmdschedule/schedule_import.go` | 30 | `runScheduleImport` | Uses raw inline types | Modernize with `ScheduleExportBundle` | Completed |
| V-24 | `cli/cmdschedule/schedule_export_import_test.go` | 1 | Call site | References unexported inline types | Update to `ScheduleExportBundle` | Completed |
| V-25 | `cli/macro/types.go` | 1 | `types.go` Result Aliases | Missing single reusable Result aliases | Define `MacroSliceResult`, `MacroResult`, `MacroStepsMapResult` | Completed |
| V-26 | `cli/macro/import.go` | 108 | `ParseImportJSON` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-27 | `cli/macro/import.go` | 123 | `ParseImportYAML` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-28 | `cli/macro/import.go` | 138 | `ParseImportZIP` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-29 | `cli/macro/import.go` | 149 | `extractMacrosFromZip` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-30 | `cli/macro/import.go` | 205 | `ParseImportFile` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-31 | `cli/macro/import.go` | 248 | `filterAndValidateImportTargets` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-32 | `cli/macro/import_sqlite.go` | 15 | `ParseImportSQLite` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-33 | `cli/macro/import_sqlite.go` | 31 | `queryAllMacrosWithSteps` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-34 | `cli/macro/import_sqlite.go` | 48 | `scanMacroRows` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-35 | `cli/macro/import_sqlite.go` | 86 | `queryAllMacroSteps` | Returns raw generic `result.ResultMap[string, []MacroStep]` | Update signature to return `MacroStepsMapResult` | Completed |
| V-36 | `cli/macro/import_sqlite.go` | 100 | `scanMacroStepsMap` | Returns raw generic `result.ResultMap[string, []MacroStep]` | Update signature to return `MacroStepsMapResult` | Completed |
| V-37 | `cli/macro/storage.go` | 106 | `ListMacros` | Returns raw generic `result.ResultSlice[Macro]` | Update signature to return `MacroSliceResult` | Completed |
| V-38 | `cli/pipelinedb/types.go` | 1 | `types.go` missing | Domain Result aliases scattered | Create `cli/pipelinedb/types.go` with single reusable Result aliases | Completed |
| V-39 | `cli/pipelinedb/pipeline_split_ops.go` | 222 | `collectRecentRuns` | Returns raw generic `result.ResultSlice[PipelineRunRecord]` | Update signature to return `PipelineRunSliceResult` | Completed |
| V-40 | `cli/pipelinedb/pipeline_split_ops.go` | 241 | `QueryRecentRuns` | Returns raw generic `result.ResultSlice[PipelineRunRecord]` | Update signature to return `PipelineRunSliceResult` | Completed |
| V-41 | `cli/pipelinedb/pipeline_split_ops.go` | 263 | `collectRecentErrors` | Returns raw generic `result.ResultSlice[PipelineErrorRecord]` | Update signature to return `PipelineErrorSliceResult` | Completed |
| V-42 | `cli/pipelinedb/pipeline_split_ops.go` | 282 | `QueryRecentErrorLogs` | Returns raw generic `result.ResultSlice[PipelineErrorRecord]` | Update signature to return `PipelineErrorSliceResult` | Completed |
| V-43 | `cli/pipelinedb/pipeline_split_ops.go` | 304 | `collectRecentCompactErrors` | Returns raw generic `result.ResultSlice[PipelineCompactErrorRecord]` | Update signature to return `PipelineCompactErrorSliceResult` | Completed |
| V-44 | `cli/pipelinedb/pipeline_split_ops.go` | 323 | `QueryDetailedErrors` | Returns raw generic `result.ResultSlice[PipelineErrorRecord]` | Update signature to return `PipelineErrorSliceResult` | Completed |
| V-45 | `cli/pipelinedb/pipeline_split_ops.go` | 336 | `QueryCompactErrors` | Returns raw generic `result.ResultSlice[PipelineCompactErrorRecord]` | Update signature to return `PipelineCompactErrorSliceResult` | Completed |
| V-46 | `cli/pipelinedb/pipeline_split_ops.go` | 510 | `collectRunIdList` | Returns raw generic `result.ResultSlice[uint64]` | Update signature to return `PipelineRunIdSliceResult` | Completed |
| V-47 | `cli/pipelinedb/pipeline_split_ops.go` | 529 | `QueryCachedErrorRunIds` | Returns raw generic `result.ResultSlice[uint64]` | Update signature to return `PipelineRunIdSliceResult` | Completed |
| V-48 | `cli/pipelinedb/pipeline_split_ops.go` | 541 | `QueryCachedErrorRunIdMap` | Returns raw generic `result.ResultMap[uint64, bool]` | Update signature to return `PipelineRunIdMapResult` | Completed |
| V-49 | `cli/pipelinedb/pipeline_split_ops.go` | 590 | `QueryLastFailedRuns` | Returns raw generic `result.ResultSlice[PipelineRunRecord]` | Update signature to return `PipelineRunSliceResult` | Completed |
| V-50 | `cli/pipelinedb/pipeline_split_ops.go` | 603 | `QueryLastNFailedRuns` | Returns raw generic `result.ResultSlice[PipelineRunRecord]` | Update signature to return `PipelineRunSliceResult` | Completed |
| V-51 | `cli/pipelinedb/pipeline_split_ops.go` | 608 | `QueryErrorLogsByRunId` | Returns raw generic `result.ResultSlice[PipelineErrorRecord]` | Update signature to return `PipelineErrorSliceResult` | Completed |
| V-52 | `cli/pipelinedb/pipeline_split_ops.go` | 620 | `QueryDetailedErrorLogsByRunId` | Returns raw generic `result.ResultSlice[PipelineErrorRecord]` | Update signature to return `PipelineErrorSliceResult` | Completed |
| V-53 | `cli/pipelinedb/pipeline_split_ops.go` | 632 | `QueryCompactErrorLogsByRunId` | Returns raw generic `result.ResultSlice[PipelineCompactErrorRecord]` | Update signature to return `PipelineCompactErrorSliceResult` | Completed |
| V-54 | `03-ai-scripts/35-result-wrapper-auditor.py` | 1 | Auditor Script | Lacks verification for types.go centralization and raw generic Result returns | Add automated types.go and single reusable type checks | Completed |

---

## 5. Subtask Decomposition

- **Subtask 01:** `.lovable/plans/subtasks/145-result-wrapper/01-result-types-centralization-and-affirmative-field.md`
  - Scope: Create `cli/result/types.go` declaring core container types `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`, aliases `Wrap[T]`, `Slice[T]`, `Map[K, V]`. Rename `defined bool` to `isDefined bool` in `cli/result/result.go`. Modernize clumsy checks in `cli/result/null_safety_test.go` to use `IsCountOtherThan(N)`.
- **Subtask 02:** `.lovable/plans/subtasks/145-result-wrapper/02-cmdschedule-types-and-bundle-result.md`
  - Scope: Create `cli/cmdschedule/types.go` declaring `ScheduleExportBundle`, `ScheduleExportOpts`, `ScheduleExportBundleResult`, `ScheduleExportBundleSingleResult`, `SchedulerTaskSliceResult`. Refactor 7 functions in `schedule_export.go` and `schedule_import.go` to return single reusable types. Modernize callers and tests.
- **Subtask 03:** `.lovable/plans/subtasks/145-result-wrapper/03-macro-and-pipelinedb-types-centralization.md`
  - Scope: Add `MacroSliceResult`, `MacroResult`, `MacroStepsMapResult` to `cli/macro/types.go` and update 10 functions in `import.go`, `import_sqlite.go`, `storage.go`. Create `cli/pipelinedb/types.go` defining single reusable Result aliases and update 14 functions in `pipeline_split_ops.go`.
- **Subtask 04:** `.lovable/plans/subtasks/145-result-wrapper/04-auditor-linter-and-quality-gates-verification.md`
  - Scope: Update `03-ai-scripts/35-result-wrapper-auditor.py` to enforce `types.go` presence and audit raw generic returns. Run all quality linters, record modified files under test inventory lock, consolidate Plan 145, update index, and commit and push atomically via SSH.

---

## 6. Verification Plan

1. **Unit Tests for Result Package:** `go test -v ./cli/result/...`
2. **Subsystem Unit Tests:** `go test -v ./cli/cmdschedule/... ./cli/macro/... ./cli/pipelinedb/...`
3. **Go Type Safety & Vet:** `go vet ./...` in `cli/` must pass with 0 errors.
4. **Boolean Guidelines Linter:** `python linter-scripts/check-boolean-guidelines.py` must pass.
5. **Nested If Linter:** `python linter-scripts/check-nested-ifs.py` must pass.
6. **Relative Paths Linter:** `python linter-scripts/check-relative-paths.py` must pass.
7. **Enum Guidelines Linter:** `python linter-scripts/check-enum-guidelines.py` must pass.
8. **Result Wrapper Auditor:** `python 03-ai-scripts/35-result-wrapper-auditor.py` must pass with 0 violations.
9. **Test Inventory Lock:** `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>` must track all modified files.
