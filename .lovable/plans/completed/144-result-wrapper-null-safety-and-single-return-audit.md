# Plan 144: Result Wrapper Types, Collections & AppError Returns (Pointer Null Safety & Predicates)

Trigger Keywords & Aliases: `cg-result-wrapper`, `cg-apperror-returns`, `cg-execute result-wrapper`, `audit result wrapper`, `fix map return error`, `fix slice return error`, `single return object audit`, `enforce apperror returns`, `enforce result map`, `fix multi-value returns`, `is-count-other-than`, `has-record`, `is-defined`, `result-wrapper-null-safety`, `pointer-null-safety`

> **Prompt Version:** 2.3.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)
> **Status:** COMPLETED

---

## 1. Executive Summary & Objective

Autonomously scan, discover, plan, refactor, and verify all Go functions returning multi-value error tuples (such as `(map[K]V, error)`, `([]T, error)`, or `(T, error)`), eliminating raw standard library error returns, replacing them with strongly-typed result wrappers (`ResultMap[K, V]`, `ResultSlice[T]`, `Result[T]`) and structured `*apperror.AppError` returns, guaranteeing a single return object, pointer-attached null safety (`*Result[T]`, `*ResultSlice[T]`, `*ResultMap[K, V]`) with methods attached to pointer receivers (`(r *Result[T])`, `(rs *ResultSlice[T])`, `(rm *ResultMap[K, V])`) that verify `if r == nil` before dereferencing any fields or checking errors, standardized outer-layer inspection predicates (`IsSuccess`, `IsFailure`, `HasError`, `IsEmptyError`, `IsEmpty`, `HasRecord`, `IsDefined`, `IsCountOtherThan`, `Data`, `Items`, `AppError`, `Fault`, `Get`, `Has`, `Count`), eliminating dual-handling, and replacing verbose `if err != nil || len(...) != N` or `IsFailure() || Count() != N` conditions with fluent `if res.IsCountOtherThan(N)` across the entire codebase until 100% green without stopping.

Plan 144 executes across four primary architectural milestones:
1. **Pointer-Attached Null Safety & Core Predicates on Result Types (`cli/result/`):** Refactor all inspection methods on `Result[T]`, `ResultSlice[T]`, and `ResultMap[K, V]` from value receivers `func (r Result[...])` to pointer receivers `func (r *Result[...])` with immediate `if r == nil` guards returning canonical safe defaults. Implement the 4 core predicate methods (`IsCountOtherThan`, `IsEmpty`, `HasRecord`/`HasRecords`, `IsDefined`) and count/item accessors across all result containers.
2. **Caller Modernization & Fluent Predicate Adoption:** Replace clumsy compound checks (`IsFailure() || Count() != N` and `IsSuccess() && !IsEmpty()`) across `cli/cmdpipeline/`, `cli/macro/`, `cli/pipelinedb/`, and `cli/cmdprompt/` with fluent single-condition guards (`IsCountOtherThan(N)`, `HasRecord()`).
3. **Subsystem Slice Migration (`cli/cmdschedule/`):** Migrate 7 multi-value slice return functions in `cli/cmdschedule/schedule_export.go` and `cli/cmdschedule/schedule_import.go` from `([]T, error)` to `result.ResultSlice[T]` with structured `*apperror.AppError` instances, updating all callers and test assertions.
4. **Automated Linter Verification & Quality Gates:** Update `03-ai-scripts/35-result-wrapper-auditor.py` to enforce `cli/cmdschedule/` and detect value receiver declarations, verify zero regressions across all repository linters (`go vet`, `check-boolean-guidelines.py`, `check-nested-ifs.py`, `check-relative-paths.py`, `check-enum-guidelines.py`), record modified files under inventory lock, consolidate Plan 144, and push atomically via SSH.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/01-index.md`: Single return type mandate and micro-tasking.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: File and function size limits (functions <= 8–15 lines, files <= 100–300 lines).
- `spec/03-error-manage/01-index.md`: Universal AppError wrapping and error envelopes.
- `spec/03-error-manage/02-error-architecture/02-error-handling-reference.md`: Error handling architecture and Result wrappers.
- `spec/03-error-manage/02-error-architecture/05-response-envelope/05-response-envelope-reference.md`: Response envelope schemas.
- `spec/03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/04-result-types.md`: Result[T], ResultSlice[T], and ResultMap[K, V] method specifications and pointer null-safety rules.
- `.agents/skills/cg-result-wrapper/skill.md`: Canonical skill for result wrappers and AppError returns.
- `.lovable/coding-guidelines.md`: Master consolidated coding guidelines.

---

## 3. Domain-Specific Rules for Result Wrappers & Pointer Null Safety

1. **Rule 1 (Pointer Receivers for All Inspection Methods):** All inspection methods on `Result[T]`, `ResultSlice[T]`, and `ResultMap[K, V]` must be attached to pointer receivers (`*Result[T]`, `*ResultSlice[T]`, `*ResultMap[K, V]`). Value receivers are strictly prohibited to prevent runtime panics when calling methods on uninitialized/nil pointers.
2. **Rule 2 (Mandatory Immediate Nil Guard):** Every pointer receiver method must guard `if r == nil` on line 1 and return safe canonical defaults (`false` for `IsSuccess`, `true` for `IsFailure`, `0` for `Count`, `true` for `IsEmpty`, `true` for `IsCountOtherThan`, `false` for `HasRecord`/`IsDefined`, `nil` for `AppError`/`Items`/`Data`).
3. **Rule 3 (The 4 Core Predicate Methods):**
   - `IsCountOtherThan(number int) bool`: Returns `true` if operation failed (or `nil`) OR count != number.
   - `IsEmpty() bool`: Returns `true` if 0 items, payload empty/null/zero, or `nil`.
   - `HasRecord() bool` (and `HasRecords()`): Returns `true` if succeeded AND count > 0.
   - `IsDefined() bool`: Returns `true` if succeeded AND payload is non-zero / has records.
4. **Rule 4 (Single Return Object Mandate):** Functions must return a single envelope object (`Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`) or `*apperror.AppError` for pure side effects. Multi-value error tuples `(T, error)` are banned.
5. **Rule 5 (Structured AppError Returns):** Raw standard library `error` returns are strictly banned; always use structured `*apperror.AppError`.
6. **Rule 6 (Caller Expressiveness):** Callers must avoid compound disjunctions like `if res.IsFailure() || res.Count() != N` in favor of `if res.IsCountOtherThan(N)`.

---

## 4. Exhaustive Violation Ledger

| Id | File | Line | Identifier | Issue | Target Refactoring | Status |
|:---|:---|:---:|:---|---|---|:---:|
| V-01 | `cli/result/result.go` | 17 | `Result.IsSuccess` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return false }` | Completed |
| V-02 | `cli/result/result.go` | 22 | `Result.IsFailed` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return true }` | Completed |
| V-03 | `cli/result/result.go` | 27 | `Result.IsFailure` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return true }` | Completed |
| V-04 | `cli/result/result.go` | 32 | `Result.IsInvalid` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return true }` | Completed |
| V-05 | `cli/result/result.go` | 37 | `Result.HasError` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return true }` | Completed |
| V-06 | `cli/result/result.go` | 42 | `Result.IsEmptyError` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return false }` | Completed |
| V-07 | `cli/result/result.go` | 47 | `Result.HasNoError` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return false }` | Completed |
| V-08 | `cli/result/result.go` | 52 | `Result.HasValidError` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return false }` | Completed |
| V-09 | `cli/result/result.go` | 61 | `Result.IsEmpty` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return true }` | Completed |
| V-10 | `cli/result/result.go` | 68 | `Result.AppError` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return nil }` | Completed |
| V-11 | `cli/result/result.go` | 73 | `Result.Fault` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return nil }` | Completed |
| V-12 | `cli/result/result.go` | 78 | `Result.Unwrap` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { var z T; return z, nil }` | Completed |
| V-13 | `cli/result/result.go` | 83 | `Result.UnwrapOr` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return defaultVal }` | Completed |
| V-14 | `cli/result/result.go` | 156 | `Result.HandleError` | Value receiver `(r Result[T])` lacks nil-safety | Change to `(r *Result[T])` with `if r == nil { return }` | Completed |
| V-15 | `cli/result/result.go` | N/A | `Result[T]` Predicates | Missing `Count`, `IsCountOtherThan`, `HasRecord`, `IsDefined` | Implement pointer-attached core predicates | Completed |
| V-16 | `cli/result/result_slice.go` | 15 | `ResultSlice.IsSuccess` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return false }` | Completed |
| V-17 | `cli/result/result_slice.go` | 20 | `ResultSlice.IsFailed` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return true }` | Completed |
| V-18 | `cli/result/result_slice.go` | 25 | `ResultSlice.IsFailure` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return true }` | Completed |
| V-19 | `cli/result/result_slice.go` | 30 | `ResultSlice.HasError` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return true }` | Completed |
| V-20 | `cli/result/result_slice.go` | 35 | `ResultSlice.IsEmptyError` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return false }` | Completed |
| V-21 | `cli/result/result_slice.go` | 40 | `ResultSlice.HasNoError` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return false }` | Completed |
| V-22 | `cli/result/result_slice.go` | 45 | `ResultSlice.IsEmpty` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return true }` | Completed |
| V-23 | `cli/result/result_slice.go` | 50 | `ResultSlice.AppError` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return nil }` | Completed |
| V-24 | `cli/result/result_slice.go` | 55 | `ResultSlice.Fault` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return nil }` | Completed |
| V-25 | `cli/result/result_slice.go` | 60 | `ResultSlice.Count` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return 0 }` | Completed |
| V-26 | `cli/result/result_slice.go` | 65 | `ResultSlice.Get` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return zero, false }` | Completed |
| V-27 | `cli/result/result_slice.go` | 76 | `ResultSlice.Unwrap` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return nil, nil }` | Completed |
| V-28 | `cli/result/result_slice.go` | 81 | `ResultSlice.UnwrapOr` | Value receiver `(r ResultSlice[T])` lacks nil-safety | Change to `(r *ResultSlice[T])` with `if r == nil { return defaultVal }` | Completed |
| V-29 | `cli/result/result_slice.go` | N/A | `ResultSlice[T]` Predicates | Missing `IsCountOtherThan`, `HasRecord`, `IsDefined`, `Items` | Implement pointer-attached core predicates and `Items()` | Completed |
| V-30 | `cli/result/result_map.go` | 18 | `ResultMap.IsSuccess` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return false }` | Completed |
| V-31 | `cli/result/result_map.go` | 23 | `ResultMap.IsFailed` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return true }` | Completed |
| V-32 | `cli/result/result_map.go` | 28 | `ResultMap.IsFailure` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return true }` | Completed |
| V-33 | `cli/result/result_map.go` | 33 | `ResultMap.HasError` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return true }` | Completed |
| V-34 | `cli/result/result_map.go` | 38 | `ResultMap.IsEmptyError` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return false }` | Completed |
| V-35 | `cli/result/result_map.go` | 43 | `ResultMap.HasNoError` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return false }` | Completed |
| V-36 | `cli/result/result_map.go` | 48 | `ResultMap.IsEmpty` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return true }` | Completed |
| V-37 | `cli/result/result_map.go` | 53 | `ResultMap.AppError` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return nil }` | Completed |
| V-38 | `cli/result/result_map.go` | 58 | `ResultMap.Fault` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return nil }` | Completed |
| V-39 | `cli/result/result_map.go` | 63 | `ResultMap.Count` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return 0 }` | Completed |
| V-40 | `cli/result/result_map.go` | 72 | `ResultMap.Get` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return zero, false }` | Completed |
| V-41 | `cli/result/result_map.go` | 85 | `ResultMap.Has` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return false }` | Completed |
| V-42 | `cli/result/result_map.go` | 96 | `ResultMap.Keys` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return []K{} }` | Completed |
| V-43 | `cli/result/result_map.go` | 114 | `ResultMap.Values` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return []V{} }` | Completed |
| V-44 | `cli/result/result_map.go` | 125 | `ResultMap.Unwrap` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return nil, nil }` | Completed |
| V-45 | `cli/result/result_map.go` | 130 | `ResultMap.UnwrapOr` | Value receiver `(r ResultMap[K, V])` lacks nil-safety | Change to `(r *ResultMap[K, V])` with `if r == nil { return defaultVal }` | Completed |
| V-46 | `cli/result/result_map.go` | N/A | `ResultMap[K, V]` Predicates | Missing `IsCountOtherThan`, `HasRecord`, `IsDefined` | Implement pointer-attached core predicates | Completed |
| V-47 | `cli/cmdpipeline/pipeline_compact_test.go` | 199 | Call site | Compound check `detailRes.IsFailure() \|\| detailRes.Count() != 1` | Replace with `detailRes.IsCountOtherThan(1)` | Completed |
| V-48 | `cli/cmdpipeline/pipeline_compact_test.go` | 208 | Call site | Compound check `compactRes.IsFailure() \|\| compactRes.Count() != 1` | Replace with `compactRes.IsCountOtherThan(1)` | Completed |
| V-49 | `cli/macro/export_import_test.go` | 34 | Call site | Compound check `listRes.IsFailure() \|\| listRes.Count() != 2` | Replace with `listRes.IsCountOtherThan(2)` | Completed |
| V-50 | `cli/macro/export_import_test.go` | 44 | Call site | Compound check `singleRes.IsFailure() \|\| singleRes.Count() != 1` | Replace with `singleRes.IsCountOtherThan(1)` | Completed |
| V-51 | `cli/macro/export_import_test.go` | 58 | Call site | Compound check `yamlRes.IsFailure() \|\| yamlRes.Count() != 2` | Replace with `yamlRes.IsCountOtherThan(2)` | Completed |
| V-52 | `cli/macro/export_import_test.go` | 68 | Call site | Compound check `singleYamlRes.IsFailure() \|\| singleYamlRes.Count() != 1` | Replace with `singleYamlRes.IsCountOtherThan(1)` | Completed |
| V-53 | `cli/pipelinedb/pipeline_split_db_test.go` | 63 | Call site | Compound check `logRes.IsFailure() \|\| logRes.Count() != 1` | Replace with `logRes.IsCountOtherThan(1)` | Completed |
| V-54 | `cli/cmdpipeline/pipeline_history.go` | 254 | Call site | Verbose check `detailRes.IsSuccess() && !detailRes.IsEmpty()` | Replace with `detailRes.HasRecord()` | Completed |
| V-55 | `cli/cmdpipeline/pipeline_history.go` | 266 | Call site | Verbose check `compactRes.IsSuccess() && !compactRes.IsEmpty()` | Replace with `compactRes.HasRecord()` | Completed |
| V-56 | `cli/cmdprompt/prompt_workdir_resolver.go` | 26 | Call site | Verbose check `childRes.IsSuccess() && !childRes.IsEmpty()` | Replace with `childRes.HasRecord()` | Completed |
| V-57 | `cli/cmdschedule/schedule_export.go` | 112 | `collectExportBundles` | Multi-value tuple `([]scheduleExportBundle, error)` | Change to return `result.ResultSlice[scheduleExportBundle]` | Completed |
| V-58 | `cli/cmdschedule/schedule_import.go` | 48 | `parseImportFileBundles` | Multi-value tuple `([]scheduleExportBundle, error)` | Change to return `result.ResultSlice[scheduleExportBundle]` | Completed |
| V-59 | `cli/cmdschedule/schedule_import.go` | 62 | `parseImportJSON` | Multi-value tuple `([]scheduleExportBundle, error)` | Change to return `result.ResultSlice[scheduleExportBundle]` | Completed |
| V-60 | `cli/cmdschedule/schedule_import.go` | 81 | `parseImportYAML` | Multi-value tuple `([]scheduleExportBundle, error)` | Change to return `result.ResultSlice[scheduleExportBundle]` | Completed |
| V-61 | `cli/cmdschedule/schedule_import.go` | 100 | `parseImportSQLite` | Multi-value tuple `([]scheduleExportBundle, error)` | Change to return `result.ResultSlice[scheduleExportBundle]` | Completed |
| V-62 | `cli/cmdschedule/schedule_import.go` | 126 | `queryImportTasksFromDB` | Multi-value tuple `([]store.SchedulerTask, error)` | Change to return `result.ResultSlice[store.SchedulerTask]` | Completed |
| V-63 | `cli/cmdschedule/schedule_import.go` | 182 | `parseImportZIP` | Multi-value tuple `([]scheduleExportBundle, error)` | Change to return `result.ResultSlice[scheduleExportBundle]` | Completed |
| V-64 | `cli/cmdschedule/schedule_export_import_test.go` | 51 | Call site | Compound check `err != nil \|\| len(bundles) != 1` | Replace with `bundlesRes.IsCountOtherThan(1)` | Completed |

---

## 5. Subtask Decomposition

- **Subtask 01:** `.lovable/plans/subtasks/144-result-wrapper/01-result-wrapper-pointer-null-safety.md`
  - Scope: Refactor all methods in `cli/result/result.go`, `cli/result/result_slice.go`, `cli/result/result_map.go` to pointer receivers with `if r == nil` guards. Implement `IsCountOtherThan`, `IsEmpty`, `HasRecord`, `HasRecords`, `IsDefined`, `Count`, `Items`. Author comprehensive unit tests in `cli/result/null_safety_test.go`.
- **Subtask 02:** `.lovable/plans/subtasks/144-result-wrapper/02-caller-modernization-and-fluent-predicates.md`
  - Scope: Modernize callers across `cli/cmdpipeline/`, `cli/macro/`, `cli/pipelinedb/`, `cli/cmdprompt/` replacing `IsFailure() || Count() != N` with `IsCountOtherThan(N)` and `IsSuccess() && !IsEmpty()` with `HasRecord()`.
- **Subtask 03:** `.lovable/plans/subtasks/144-result-wrapper/03-cmdschedule-resultslice-migration.md`
  - Scope: Refactor 7 multi-value slice returns in `cli/cmdschedule/` (`collectExportBundles`, `parseImportFileBundles`, `parseImportJSON`, `parseImportYAML`, `parseImportSQLite`, `queryImportTasksFromDB`, `parseImportZIP`) to `result.ResultSlice[T]` with structured `*apperror.AppError`. Modernize `schedule_export_import_test.go:51` to `bundlesRes.IsCountOtherThan(1)`.
- **Subtask 04:** `.lovable/plans/subtasks/144-result-wrapper/04-automated-linter-and-quality-gates-verification.md`
  - Scope: Update `03-ai-scripts/35-result-wrapper-auditor.py` to enforce `cli/cmdschedule/` and audit value receiver / clumsy check regressions. Run all quality linters, record modified files under lock, consolidate Plan 144, update index, commit and push atomically.

---

## 6. Verification Plan

1. **Unit Tests for Null Safety:** `go test -v ./cli/result/...` verifying 100% pass on nil pointer receiver invocations.
2. **Go Compilation & Typing:** `go vet ./...` in `cli/` must pass with 0 errors.
3. **Boolean Guidelines Linter:** `python linter-scripts/check-boolean-guidelines.py` must pass.
4. **Nested If Linter:** `python linter-scripts/check-nested-ifs.py` must pass.
5. **Relative Paths Linter:** `python linter-scripts/check-relative-paths.py` must pass.
6. **Enum Guidelines Linter:** `python linter-scripts/check-enum-guidelines.py` must pass.
7. **Result Wrapper Auditor:** `python 03-ai-scripts/35-result-wrapper-auditor.py` must pass with 0 violations.
8. **Test Inventory Lock:** `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>` must track all modified files.
