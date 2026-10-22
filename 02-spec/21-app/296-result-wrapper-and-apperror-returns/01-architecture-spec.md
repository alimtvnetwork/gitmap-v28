# Architecture Specification: Result Wrapper & AppError Return Sovereignty (Task 296)

> **Spec ID:** `SPEC-21-296-ARCH`  
> **Subsystem:** `cli/store`, `cli/heatmap`, `cli/result`, `cli/apperror`  
> **Status:** Active / Specification  
> **Date:** 2026-10-23  
> **Task:** 296-result-wrapper-and-apperror-returns  
> **Related Specs:**  
> - `02-spec/02-coding-guidelines/`  
> - `02-spec/03-error-manage/`  
> - `01-prompts/15-cg-execute/19-result-wrapper-and-apperror-returns.md`  
> **Related Plans:**  
> - `.ai-memory/plans/subtasks/296-result-wrapper-and-apperror-returns/subtask-01.md`  
> - `.ai-memory/plans/subtasks/296-result-wrapper-and-apperror-returns/subtask-02.md`  
> - `.ai-memory/plans/subtasks/296-result-wrapper-and-apperror-returns/subtask-03.md`  

---

## 1. Executive Summary & Root Cause Analysis (RCA)

In enterprise Go architecture, functions must return a single deterministic outcome rather than legacy multi-value error tuples such as `(T, error)`, `([]T, error)`, or `(map[K]V, error)`. Multi-value error returns induce severe software quality degradation:
1. **Accidental `nil` Dereferences:** Callers frequently neglect error checks or access the value tuple when an error occurred, leading to panics.
2. **Dual-Handling Chaos:** When a function returns both a partial value and an error, callers are forced into ambiguous state recovery logic.
3. **Double-Wrapping Anti-Patterns:** Developers attempt to fix multi-return signatures by wrapping constructor calls inside type aliases (e.g. `FailWrap(apperror.WrapSimple(err, op))`), creating broken syntax and bloated stack traces.
4. **Package Stutter & Leaky Generics:** Direct usage of `result.Result[[]model.ScanRecord]` creates stutter (`result.Result`), exposes raw generic brackets in domain signatures, and wraps unaliased anonymous slices.

### 4-Part Root Cause Analysis

```
+---------------------------------------------------------------------------------------------------+
| 1. SYMPTOM                                                                                        |
| Functions in cli/store/ and cli/heatmap/ return legacy multi-value error tuples (T, error) or     |
| bare error types. Banned anti-pattern ScanRecordsResult = result.Result[[]model.ScanRecord]       |
| lingers in cli/store/repository_types.go with package stutter and raw generic brackets.          |
+---------------------------------------------------------------------------------------------------+
                                                  |
                                                  v
+---------------------------------------------------------------------------------------------------+
| 2. MULTI-VALUE TUPLE ROOT CAUSE                                                                  |
| Methods in cli/store/repo.go (ListRepos, FindBySlug, FindOneBySlug, FindByPath, CountRepos,       |
| DeleteByPath, DeleteBySlug) rely on (T, error) or (int64, error). This breaks single-return       |
| sovereignty, prevents monadic composition, and encourages callers to swallow or bypass errors.    |
+---------------------------------------------------------------------------------------------------+
                                                  |
                                                  v
+---------------------------------------------------------------------------------------------------+
| 3. BARE ERROR & DOUBLE-WRAPPING ROOT CAUSE                                                        |
| Heatmap tracking APIs (TrackFileTouch, TrackBatchFileTouches, TrackAgentAction) return standard   |
| library error without structured categorization or audit metadata. Attempted fixes inadvertently  |
| invoke non-existent FailWrap constructors or double-wrap errors with redundant WrapSimple calls.  |
+---------------------------------------------------------------------------------------------------+
                                                  |
                                                  v
+---------------------------------------------------------------------------------------------------+
| 4. TYPE CENTRALIZATION GAP ROOT CAUSE                                                             |
| Store domain types and Result containers are scattered or defined inline in implementation files. |
| Without a centralized cli/store/types.go file, developers invent ad-hoc slice signatures and      |
| duplicate aliases, resulting in the banned ScanRecordsResult stutter anti-pattern.                |
+---------------------------------------------------------------------------------------------------+
```

---

## 2. Core Architectural Pillars

```
+-------------------------------------------------------------------------------------------------------+
|                                    SPEC-21-296 ARCHITECTURE PILLARS                                   |
+-----------------------------------+-----------------------------------+-------------------------------+
| 1. Monadic Return Sovereignty     | 2. Pointer Null Safety            | 3. Structured Error Types     |
| - Single return value per func    | - (r *Result[T]) pointer receiver | - *apperror.AppError returns  |
| - result.Wrap[T] for scalars      | - Line-1 if r == nil guard        | - Zero bare error interfaces  |
| - result.WrapSlice[T] for slices  | - Zero panic on nil dereference   | - Standardized error codes    |
| - result.WrapMap[K, V] for maps   | - Safe defaults returned          | - Audit trail integration     |
+-----------------------------------+-----------------------------------+-------------------------------+
| 4. Zero Double-Wrapping           | 5. types.go Centralization        | 6. Stutter Eradication        |
| - Total ban on FailWrap           | - All domain aliases in types.go  | - Ban result.Result stutter   |
| - Wrap[T] is alias, not func      | - Zero leaky generic brackets     | - Ban Result[[]model.X]       |
| - Use result.Fail[T](appErr)      | - Clean named collection types    | - Dedicated XSlice types      |
+-----------------------------------+-----------------------------------+-------------------------------+
```

### Pillar 1: Single-Return Monadic Containers
Every fallible operation that produces data MUST return exactly one value:
- Scalar / Entity Lookup: `result.Wrap[T]`
- Slice Collection Query: `result.WrapSlice[T]`
- Keyed Dictionary Lookup: `result.WrapMap[K, V]`

The container encapsulates both the successful payload and any failure state (`*apperror.AppError`), guaranteeing that a function returns a complete, self-contained result object.

### Pillar 2: Pointer-Attached Null Safety
All container inspection and extraction methods must attach to pointer receivers (`*Result[T]`, `*ResultSlice[T]`, `*ResultMap[K, V]`). Every method implements a line-1 `if r == nil` guard that returns safe defaults without panicking:
- `IsSuccess()` returns `false` when receiver is nil.
- `IsFailure()` returns `true` when receiver is nil.
- `IsEmpty()` returns `true` when receiver is nil.
- `Count()` returns `0` when receiver is nil.
- `Data()` / `Value` returns zero-value when receiver is nil.
- `AppError()` returns `nil` (or a fallback error) when receiver is nil.

### Pillar 3: Structured Enterprise Error Returns
Pure mutation and telemetry functions that return no payload (such as tracking and logging in `cli/heatmap/`) must return structured `*apperror.AppError` instead of standard library `error`. This enforces:
- Explicit error classification (`ErrorTypeExecution`, `ErrorTypeValidation`, `ErrorTypeDatabase`).
- Originating caller tagging (e.g. `"heatmap.TrackFileTouch"`).
- Rich context for CLI exit handlers and SQLite database audit capture.

### Pillar 4: Total Ban on Double-Wrapping (`FailWrap` + `WrapSimple`)
- `result.Wrap[T]`, `result.WrapSlice[T]`, and `result.WrapMap[K, V]` are **TYPE ALIASES**, NOT constructor functions.
- Calling `FailWrap` or `FailWrapSlice` is strictly forbidden.
- The only permitted constructor functions are:
  - Success: `result.Ok(val)`, `result.OkSlice(slice)`, `result.OkMap(mapData)`
  - Failure: `result.Fail[T](appErr)`, `result.FailSlice[T](appErr)`, `result.FailMap[K, V](appErr)`
- Callers must never double-wrap an existing `*apperror.AppError` with redundant `apperror.WrapSimple` calls.

### Pillar 5: Centralization in `cli/store/types.go` & Eradication of Package Stutter
All domain models, slice collections, and Result aliases used across `cli/store/` must be centralized in `cli/store/types.go`.
- **Anti-Pattern Banned:** `ScanRecordsResult = result.Result[[]model.ScanRecord]` in `cli/store/repository_types.go`.
  - Violates rule 1: Package stutter (`result.Result`).
  - Violates rule 2: Leaks raw generic brackets (`[...]`) into domain signatures.
  - Violates rule 3: Wraps an anonymous slice `[]model.ScanRecord` rather than a named slice type.
- **Mandatory Replacement:**
  ```go
  type (
      ScanRecordSlice     = []model.ScanRecord
      ScanRecordWrap      = result.Wrap[model.ScanRecord]
      ScanRecordSliceWrap = result.WrapSlice[model.ScanRecord]
      ScanRecordMapWrap   = result.WrapMap[string, model.ScanRecord]
  )
  ```

---

## 3. Centralized Domain Type Catalog (`cli/store/types.go`)

The following 13 type definitions must be authored in `cli/store/types.go`:

| Type Alias | Target Underlying Type | Category | Constructor Function | Inspection Predicates |
| :--- | :--- | :--- | :--- | :--- |
| `ScanRecordSlice` | `[]model.ScanRecord` | Model Slice Collection | N/A (standard slice) | `len(s)` |
| `ScanRecordWrap` | `result.Wrap[model.ScanRecord]` | Scalar Result Monad | `result.Ok(rec)`, `result.Fail[model.ScanRecord](err)` | `res.IsSuccess()`, `res.IsFailure()`, `res.Data()` |
| `ScanRecordSliceWrap` | `result.WrapSlice[model.ScanRecord]` | Collection Result Monad | `result.OkSlice(recs)`, `result.FailSlice[model.ScanRecord](err)` | `res.IsSuccess()`, `res.Count()`, `res.Data()` |
| `ScanRecordMapWrap` | `result.WrapMap[string, model.ScanRecord]` | Keyed Lookup Monad | `result.OkMap(m)`, `result.FailMap[string, model.ScanRecord](err)` | `res.IsSuccess()`, `res.Has(key)`, `res.Get(key)` |
| `HeatmapActionLogRecordSliceWrap` | `result.WrapSlice[HeatmapActionLogRecord]` | Heatmap Action Monad | `result.OkSlice(recs)`, `result.FailSlice[HeatmapActionLogRecord](err)` | `res.IsSuccess()`, `res.Count()`, `res.Data()` |
| `HeatmapExecutionTreeNodeSliceWrap` | `result.WrapSlice[HeatmapExecutionTreeNode]` | Heatmap Tree Monad | `result.OkSlice(trees)`, `result.FailSlice[HeatmapExecutionTreeNode](err)` | `res.IsSuccess()`, `res.Count()`, `res.Data()` |
| `HeatmapStatusSummaryWrap` | `result.Wrap[HeatmapStatusSummary]` | Heatmap Summary Monad | `result.Ok(summary)`, `result.Fail[HeatmapStatusSummary](err)` | `res.IsSuccess()`, `res.Data()` |
| `HeatmapTaskSummaryWrap` | `result.Wrap[HeatmapTaskSummary]` | Heatmap Task Monad | `result.Ok(summary)`, `result.Fail[HeatmapTaskSummary](err)` | `res.IsSuccess()`, `res.Data()` |
| `FileTouchCountSliceWrap` | `result.WrapSlice[FileTouchCount]` | File Touch Monad | `result.OkSlice(touches)`, `result.FailSlice[FileTouchCount](err)` | `res.IsSuccess()`, `res.Count()`, `res.Data()` |
| `DuplicateRepoGroupSliceWrap` | `result.WrapSlice[DuplicateRepoGroup]` | Duplicates Monad | `result.OkSlice(groups)`, `result.FailSlice[DuplicateRepoGroup](err)` | `res.IsSuccess()`, `res.Count()`, `res.Data()` |
| `RowsAffectedWrap` | `result.Wrap[int64]` | DB Mutation Monad | `result.Ok(rows)`, `result.Fail[int64](err)` | `res.IsSuccess()`, `res.Data()` |
| `RepoIDWrap` | `result.Wrap[int64]` | Repository ID Monad | `result.Ok(id)`, `result.Fail[int64](err)` | `res.IsSuccess()`, `res.Data()` |
| `RepoCountWrap` | `result.Wrap[int]` | Repository Count Monad | `result.Ok(count)`, `result.Fail[int](err)` | `res.IsSuccess()`, `res.Data()` |

---

## 4. Heatmap Subsystem Error Sovereignty (`cli/heatmap/`)

### Current State vs. Target State

```
+----------------------------------------------------+----------------------------------------------------+
| CURRENT SIGNATURE (cli/heatmap/)                   | TARGET SIGNATURE (cli/heatmap/)                    |
+----------------------------------------------------+----------------------------------------------------+
| TrackFileTouch(path, touchType) error              | TrackFileTouch(path, touchType) *apperror.AppError |
| recordFileTouchAtRoot(cwd, path, touchType) error  | recordFileTouchAtRoot(cwd, path, t) *apperror.AppError|
| TrackBatchFileTouches(cwd, touches) error          | TrackBatchFileTouches(cwd, touches) *apperror.AppError|
| insertBatchTouchesTx(conn, touches) error          | insertBatchTouchesTx(conn, touches) *apperror.AppError|
| TrackAgentAction(slug, action, summary) error      | TrackAgentAction(slug, action, sum) *apperror.AppError|
| recordAgentActionAtRoot(cwd, slug, act, sum) error | recordAgentActionAtRoot(cwd, ...) *apperror.AppError|
+----------------------------------------------------+----------------------------------------------------+
```

### Architectural Contract for `cli/heatmap/heatmap_tracker.go`
```go
package heatmap

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TrackFileTouch records a file access operation into the root heatmap database.
func TrackFileTouch(filePath string, touchType store.FileTouchType) *apperror.AppError {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimSpace(filePath)))
	isValidPath := len(cleanPath) > 0 && cleanPath != "."
	if !isValidPath {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.TrackFileTouch.cwd")
	}

	return recordFileTouchAtRoot(cwd, cleanPath, touchType)
}
```

### Invariants:
1. Zero bare `error` returns in `cli/heatmap/heatmap_tracker.go` and `cli/heatmap/heatmap_tracker_agent.go`.
2. Guard against missing or blank arguments early, returning `nil` when no operation needs recording.
3. Every external call (`os.Getwd`, `store.OpenHeatmapRootDB`, `hdb.LogAction`) that fails MUST be wrapped once with `apperror.WrapSimple(err, op)`. Never double-wrap.

---

## 5. Store Repository Monadic Refactoring (`cli/store/repo.go`)

### Current State vs. Target State

```
+----------------------------------------------------+----------------------------------------------------+
| CURRENT SIGNATURE (cli/store/repo.go)              | TARGET SIGNATURE (cli/store/repo.go)               |
+----------------------------------------------------+----------------------------------------------------+
| ListRepos() ([]model.ScanRecord, error)            | ListRepos() ScanRecordSliceWrap                    |
| FindBySlug(slug string) ([]model.ScanRecord, error)| FindBySlug(slug string) ScanRecordSliceWrap        |
| FindOneBySlug(slug string) (*model.ScanRecord, err)| FindOneBySlug(slug string) ScanRecordWrap          |
| FindByPath(absPath string) ([]model.ScanRecord, err| FindByPath(absPath string) ScanRecordSliceWrap     |
| CountRepos() (int, error)                          | CountRepos() RepoCountWrap                         |
| DeleteByPath(absPath string) (int64, error)        | DeleteByPath(absPath string) RowsAffectedWrap      |
| DeleteBySlug(slug string) (int64, error)           | DeleteBySlug(slug string) RowsAffectedWrap         |
+----------------------------------------------------+----------------------------------------------------+
```

### Architectural Contract Examples

#### `ListRepos() ScanRecordSliceWrap`
```go
// ListRepos returns all tracked repositories ordered by slug.
func (db *DB) ListRepos() ScanRecordSliceWrap {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllRepos).Destruct()
	if err != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "store.ListRepos.query"))
	}
	defer rows.Close()

	records, scanErr := scanRows(rows)
	if scanErr != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(scanErr, "store.ListRepos.scan"))
	}

	return result.OkSlice(records)
}
```

#### `FindOneBySlug(slug string) ScanRecordWrap`
```go
// FindOneBySlug returns the first repository matching the given slug.
func (db *DB) FindOneBySlug(slug string) ScanRecordWrap {
	sliceRes := db.FindBySlug(slug)
	if sliceRes.IsFailure() {
		return result.Fail[model.ScanRecord](sliceRes.AppError())
	}

	if sliceRes.IsEmpty() || sliceRes.Count() == 0 {
		return result.Fail[model.ScanRecord](NewRepoNotFoundError("FindOneBySlug", "Repo", slug))
	}

	return result.Ok(sliceRes.Data()[0])
}
```

#### `DeleteByPath(absPath string) RowsAffectedWrap`
```go
// DeleteByPath removes the repo row whose AbsolutePath matches.
func (db *DB) DeleteByPath(absPath string) RowsAffectedWrap {
	cleanPath := NormalizeStoragePath(absPath)
	res, err := ExecWrapper(db.conn, constants.SQLDeleteRepoByPath, cleanPath).Destruct()
	if err != nil {
		return result.Fail[int64](apperror.WrapSimple(err, "store.DeleteByPath"))
	}

	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		return result.Fail[int64](apperror.WrapSimple(rowsErr, "store.DeleteByPath.rows"))
	}

	return result.Ok(rows)
}
```

#### `CountRepos() RepoCountWrap`
```go
// CountRepos returns the total number of repositories tracked in the database.
func (db *DB) CountRepos() RepoCountWrap {
	var total int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM Repo").Scan(&total)
	if err != nil {
		return result.Fail[int](apperror.WrapSimple(err, "store.CountRepos"))
	}

	return result.Ok(total)
}
```

---

## 6. Caller Consumption & Affirmative Guard Patterns

Callers consuming single-return monadic Result containers MUST strictly follow affirmative predicate inspection before accessing data:

### Standard Affirmative Guard Pattern
```go
res := db.ListRepos()

// Guard against failure early:
if res.IsFailure() {
	return apperror.WrapSimple(res.AppError(), "caller.context")
}

// Access data payload safely:
records := res.Data()
for _, rec := range records {
	process(rec)
}
```

### Scalar Lookup Pattern
```go
res := db.FindOneBySlug("gitmap")

if res.IsFailure() {
	return res.AppError()
}

rec := res.Data()
fmt.Println("Found repo:", rec.RepoName)
```

---

## 7. Subtask Decomposition & Execution Plan

```
+-------------------------------------------------------------------------------------------------------+
| TASK 296 EXECUTION ROADMAP                                                                            |
+-------------------------------------------------------------------------------------------------------+
| Subtask 01: Centralize Store Domain Types in cli/store/types.go & Eradicate Stutter                   |
| - Target: cli/store/types.go (new), cli/store/repository_types.go (modify)                            |
| - Actions: Declare 13 aliases; purge ScanRecordsResult = result.Result[[]model.ScanRecord]            |
+-------------------------------------------------------------------------------------------------------+
| Subtask 02: Heatmap Tracker Error Sovereignty Refactoring                                             |
| - Target: cli/heatmap/heatmap_tracker.go, cli/heatmap/heatmap_tracker_agent.go                        |
| - Actions: Convert TrackFileTouch, TrackBatchFileTouches, TrackAgentAction from error to *AppError    |
+-------------------------------------------------------------------------------------------------------+
| Subtask 03: Store Repository Monadic Container Refactoring                                            |
| - Target: cli/store/repo.go                                                                           |
| - Actions: Refactor ListRepos, FindBySlug, FindOneBySlug, FindByPath, CountRepos, DeleteByPath/Slug  |
+-------------------------------------------------------------------------------------------------------+
```

---

## 8. Non-Negotiable Invariants

1. **Strictly Relative Git Paths:** All documentation, code comments, references, and logs must use clean relative Git paths (e.g. `cli/store/types.go`). No absolute paths, drive letters (`D:\`), or URI schemes (`file:///`).
2. **Zero Compiler Builds & Zero Tests:** No `go build`, `go test`, or execution runs may occur during specification and planning.
3. **Zero Git Commands:** All tracking is managed via `.ai-memory/plans/subtasks/` without invoking git CLI commands.
4. **Zero Double-Wrapping:** Never invoke non-existent `FailWrap` or nest `apperror.WrapSimple` within an existing `*apperror.AppError`.
5. **Strict Lowercase File Paths:** All files created or updated must use lowercase filenames and paths.
