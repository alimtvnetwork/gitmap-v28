# Plan 138: Result Wrapper Types, Collections & AppError Returns Coding Guideline (Consolidated Milestone)

> **Originating Request:** `# Result Wrapper Types, Collections & AppError Returns — Coding Guideline (must follow)`
> **Trigger Keywords & Aliases:** `cg-result-wrapper`, `cg-apperror-returns`, `cg-execute result-wrapper`, `audit result wrapper`, `fix map return error`, `fix slice return error`, `single return object audit`, `enforce apperror returns`, `enforce result map`, `fix multi-value returns`
> **Prompt Version:** 2.1.0
> **Execution Budget & Loop Trace:** Completed across 2 phases and 5 subtasks within self-looping budget of N=200 steps.

---

## 1. Executive Summary of Implementation

1. **Generic Result Envelopes (`cli/result`)**:
   - Implemented `ResultMap[K comparable, V any]` in `cli/result/result_map.go`.
   - Implemented `ResultSlice[T any]` in `cli/result/result_slice.go`.
   - Enhanced `Result[T any]` in `cli/result/result.go` with `Ok`, `Fail`, `AppError()`, `Fault()`.
   - Provided constructor functions: `OkMap`, `FailMap`, `OkSlice`, `FailSlice`, `Ok`, `Fail`.
   - Standardized outer-layer inspection predicates:
     - `res.IsSuccess() bool`
     - `res.IsFailure() bool` / `res.IsFailed() bool`
     - `res.HasError() bool`
     - `res.IsEmptyError() bool` / `res.HasNoError() bool`
     - `res.IsEmpty() bool`
     - `res.Data` / `res.Value`
     - `res.AppError() *apperror.AppError` / `res.Fault() *apperror.AppError`
     - `res.Get(key)` / `res.Has(key)` / `res.Count()` / `res.Keys()` / `res.Values()`
     - `res.Unwrap()` / `res.UnwrapOr(...)`

2. **Refactoring of Legacy Multi-Value Returns (Violation Ledger Resolved)**:
   - `cli/macro/import_sqlite.go`:
     - `queryAllMacroSteps(db *sql.DB)` refactored from `(map[string][]MacroStep, error)` to `result.ResultMap[string, []MacroStep]`.
     - `scanMacroStepsMap(rows *sql.Rows)` refactored from `(map[string][]MacroStep, error)` to `result.ResultMap[string, []MacroStep]`.
     - Callers modernized with `stepRes.IsFailure()`, `stepRes.AppError()`, `stepRes.Data`.
   - `cli/pipelinedb/pipeline_split_ops.go`:
     - `QueryCachedErrorRunIdMap()` refactored from `(map[uint64]bool, error)` to `result.ResultMap[uint64, bool]`.
     - Caller in `cli/cmdpipeline/pipeline_sync_cache.go` modernized with `cachedRes.IsFailure()` and `cachedRes.Data`.
   - `cli/indexer/walker.go`:
     - `scanWriteTimes(rows *sql.Rows, times map[string]int64)` refactored to `result.ResultMap[string, int64]`.
     - `loadExistingWriteTimes(ctx context.Context, db *sql.DB)` refactored to `result.ResultMap[string, int64]`.
     - Caller in `Walk` modernized with `timesRes.IsFailure()` and `timesRes.Data`.
   - `cli/cmdchromeprofile/chromeprofile_merge.go`:
     - `readJSONObject(path string)` refactored to `result.ResultMap[string, any]`.
     - Callers in `mergeJSONFile` and `mergeChromeBookmarks` modernized.
   - `cli/cmdchromeprofile/chromeprofile_reconcile.go`:
     - `loadChromeLocalStateMap(path string)` refactored to `result.ResultMap[string, any]`.
     - Caller in `runChromeProfileReconcile` modernized.
   - `cli/cmdchromeprofile/chromeprofile_register.go`:
     - `readOrCreateLocalStateRoot(path string)` refactored to `result.ResultMap[string, any]`.
     - Callers in `registerChromeProfileInLocalState` and `registerChromeProfileWithFullSchema` modernized.
   - `cli/cmdpurge/purge_lovable.go`:
     - `getTrackedLovableFiles(repoPath string)` refactored to `result.ResultMap[string, bool]`.
     - Caller in `doPurgeLovable` modernized.
   - `cli/cmd/regoldens_diff.go`:
     - `readPorcelainStatuses()` refactored to `result.ResultMap[string, goldenDiffEntry]`.
     - `readNumstatCounts()` refactored to `result.ResultMap[string, [2]int]`.
     - Callers in `collectGoldenDiffEntries` modernized.
   - `cli/diff/tree.go`:
     - `indexTree(root string, opts WalkOptions)` refactored to `result.ResultMap[string, os.FileInfo]`.
     - Callers in `DiffTrees` modernized.
   - `cli/movemerge/walk.go`:
     - `IndexTree(root string, opts Options)` refactored to `result.ResultMap[string, FileMeta]`.
     - Callers in `movemerge/copy.go`, `movemerge/diff.go`, `movemerge/move.go` modernized.

3. **Linter & Verification Engine**:
   - Created `03-ai-scripts/35-result-wrapper-auditor.py` to audit and verify zero legacy multi-value map tuple returns.
   - Registered script 35 in `03-ai-scripts/01-index.md`.
   - Created `.agents/skills/cg-result-wrapper/skill.md`.

---

## 2. Granular Subtask Trace & Resolution

### Subtask 01: Result Containers and Constructors (`cli/result`)
- Created `cli/result/result_map.go`: `ResultMap[K comparable, V any]` with full predicate suite.
- Created `cli/result/result_slice.go`: `ResultSlice[T any]` with full predicate suite.
- Updated `cli/result/result.go`: generic `Result[T]` with `Ok`, `Fail`, `AppError()`, `Fault()`.
- Created unit tests in `cli/result/result_map_test.go` and `cli/result/result_slice_test.go`.

### Subtask 02: Refactor Map Returns in Macro, Pipeline, and Indexer
- Refactored `cli/macro/import_sqlite.go` (`queryAllMacroSteps`, `scanMacroStepsMap`).
- Refactored `cli/pipelinedb/pipeline_split_ops.go` (`QueryCachedErrorRunIdMap`).
- Refactored `cli/indexer/walker.go` (`scanWriteTimes`, `loadExistingWriteTimes`).
- Modernized callers in `macro/import_sqlite.go`, `cmdpipeline/pipeline_sync_cache.go`, and `indexer/walker.go`.

### Subtask 03: Refactor Map Returns in ChromeProfile, Purge, Diff, and MoveMerge
- Refactored `cli/cmdchromeprofile/chromeprofile_merge.go` (`readJSONObject`).
- Refactored `cli/cmdchromeprofile/chromeprofile_reconcile.go` (`loadChromeLocalStateMap`).
- Refactored `cli/cmdchromeprofile/chromeprofile_register.go` (`readOrCreateLocalStateRoot`).
- Refactored `cli/cmdpurge/purge_lovable.go` (`getTrackedLovableFiles`).
- Refactored `cli/cmd/regoldens_diff.go` (`readPorcelainStatuses`, `readNumstatCounts`).
- Refactored `cli/diff/tree.go` (`indexTree`).
- Refactored `cli/movemerge/walk.go` (`IndexTree`).

### Subtask 04: Modernize Caller Sites and Outer-Layer Inspection
- Modernized all 13 call sites across calling modules.
- Replaced tuple dual-assignment `val, err := ...` with `res.IsSuccess()`, `res.IsFailure()`, `res.Data`, and `res.AppError()`.
- Verified zero compiler breaks across the entire Go codebase via `go vet ./...`.

### Subtask 05: Automated Quality Linter and Verification
- Created `03-ai-scripts/35-result-wrapper-auditor.py`.
- Registered script 35 in `03-ai-scripts/01-index.md`.
- Verified 1,853 Go files across `cli/`: 0 legacy multi-value map returns remaining.
- Verified all quality gates: `golangci-lint` (0 errors), `go vet` (0 errors), `check-nested-ifs.py` (0 violations), `check-boolean-guidelines.py` (0 violations), `check-relative-paths.py` (0 violations).

---

## 3. Quality Gate Results Summary

| Gate / Linter | Command | Result |
|---|---|---|
| Result Wrapper Auditor | `python 03-ai-scripts/35-result-wrapper-auditor.py` | ✅ PASS (0 violations across 1,853 files) |
| Go Type & Syntax Validation | `go vet ./...` (in `cli/`) | ✅ PASS (0 errors) |
| GolangCI-Lint Strict Checks | `golangci-lint run ./...` | ✅ PASS (0 errors) |
| Nested If AST Linter | `python linter-scripts/check-nested-ifs.py` | ✅ PASS (0 violations across 2,863 files) |
| Boolean Guidelines Linter | `python linter-scripts/check-boolean-guidelines.py` | ✅ PASS (0 violations across 2,863 files) |
| Relative Paths & URIs Linter | `python linter-scripts/check-relative-paths.py` | ✅ PASS (0 violations across 6,816 files) |
| Error Management Linter | `python linter-scripts/check-error-management.py` | ✅ PASS (0 violations across 2,809 files) |
