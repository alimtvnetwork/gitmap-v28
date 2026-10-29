---
name: result-wrapper-and-apperror-returns
description: "Enforce enterprise Go single-return Result containers (result.Wrap[T], result.WrapSlice[T], result.WrapMap[K, V]), eliminating multi-value error tuples, centralizing types into types.go, and guaranteeing pointer-attached null safety with *apperror.AppError."
---

# Result Wrapper & AppError Return Sovereignty

> **Invocation:** `@[result-wrapper-and-apperror-returns](file;.cursor/skills/result-wrapper-and-apperror-returns)`

Source prompt: `01-prompts/15-cg-execute/19-result-wrapper-and-apperror-returns.md`

## 1. Core Mandate

In enterprise Go architecture, functions must **NEVER** return multi-value error tuples such as `(T, error)`, `([]T, error)`, or `(map[K]V, error)`. Multi-value error returns violate single-return sovereignty, induce accidental `nil` dereferences, encourage swallowed errors, and cause dual-handling chaos.

All fallible operations MUST return a single strongly-typed monadic container from `cli/result` (or `pkg/appfault` / `pkg/result`). The canonical container type is `result.Wrap[T]`.

> [!WARNING]
> **TOTAL BAN ON `result.Result` & ANONYMOUS SLICE WRAPPERS:**
> Writing `result.Result` creates package stutter and is strictly banned.
> Never write `result.Result[[]model.ScanRecord]`. This pattern violates three core rules:
> 1. **Package Stutter:** `result.Result` repeats the package name. Always use `result.Wrap`.
> 2. **Leaking Raw Generic Brackets:** Raw generic brackets `[...]` in function signatures create copy-paste boilerplate across implementation files.
> 3. **Wrapping Slices with `WrapSlice`:** Never wrap a slice inside a single-entity container (`result.Wrap[ScanRecordSlice]`). Always use the canonical slice container `result.WrapSlice[model.ScanRecord]` in `types.go`. For keyed lookups, always use `result.WrapMap[string, model.ScanRecord]`.

| Container Type | Primary Constructor Helpers | Data Payload Representation | Core Inspection Predicates |
| :--- | :--- | :--- | :--- |
| **`result.Wrap[T]`** | `result.Ok(val)` / `result.Fail[T](err)` | Single scalar, model, or struct | `res.IsSuccess()`, `res.IsFailure()`, `res.Data()`, `res.AppError()` |
| **`result.WrapSlice[T]`** | `result.OkSlice(slice)` / `result.FailSlice[T](err)` | Slice / collection of items | `res.IsSuccess()`, `res.IsFailure()`, `res.Data()`, `res.Count()`, `res.AppError()` |
| **`result.WrapMap[K, V]`** | `result.OkMap(mapData)` / `result.FailMap[K, V](err)` | Keyed lookup dictionary or index | `res.IsSuccess()`, `res.IsFailure()`, `res.Data()`, `res.Has(key)`, `res.Get(key)` |
| **`*appfault.AppError`** (`*apperror.AppError`) | `appfault.New(...)`, `appfault.Wrap(...)`, `apperror.WrapSimple(err, op)` | Structured enterprise error | `err.Message()`, `err.Code()`, `err.StackTrace()` |

---

## 2. Mandatory Rules

1. **Single Return Object:** Every function returning fallible data must return exactly one value: `result.Wrap[T]`, `result.WrapSlice[T]`, or `result.WrapMap[K, V]`.
2. **Zero Raw Error Returns:** Never return standard library `error`. All failures must wrap an underlying cause as `*appfault.AppError` (or `*apperror.AppError`).
3. **Zero Double-Wrapping (`FailWrap` + `WrapSimple`):** Never invoke constructors named `FailWrap` or `FailWrapSlice` that wrap an `apperror.WrapSimple` argument. `Wrap[T]` and `WrapSlice[T]` are purely **TYPE ALIASES**, never constructor functions. The constructor functions are strictly `result.Ok`, `result.Fail[T]`, `result.OkSlice`, `result.FailSlice[T]`, `result.OkMap`, and `result.FailMap[K, V]`.
4. **Pointer-Attached Null Safety:** All Result inspection methods (`IsSuccess`, `IsFailure`, `HasError`, `IsEmpty`, `HasRecord`, `IsCountOtherThan`, `Count`, `Data`) must attach to pointer receivers (`(r *Result[T])`, `(rs *ResultSlice[T])`, `(rm *ResultMap[K, V])`) with line-1 `if r == nil` guards returning safe defaults.
5. **Affirmative Predicate Inspection:** Callers MUST check `if res.IsFailure()` (or `if res.IsSuccess()`) before accessing payload data (`res.Data()`). Never access payload data without guarding against failure.
6. **Centralization in `types.go`:** Never leak raw generic brackets like `result.Wrap[MyEntity]` or declare domain structs inline in implementation files. Every package MUST alias Result containers in a dedicated `types.go` file as single reusable named types:
   ```go
   type (
       ItemWrap      = result.Wrap[Item]
       ItemSlice     = []Item
       ItemSliceWrap = result.WrapSlice[Item]
       ItemMap       = map[string]Item
       ItemMapWrap   = result.WrapMap[string, Item]
   )
   ```

---

## 3. Concrete Implementation Examples with DRY `types.go` Reuse

### 3.1 Single Item Lookup (`ScanRecordWrap`)

```go
// Step 1: In cli/store/types.go (in the package where it is defined for 100% DRY code)
// Declare the single reusable named type alias — NEVER leak raw generic brackets in function signatures!
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// ScanRecordWrap is the canonical single reusable result container for ScanRecord.
	ScanRecordWrap = result.Wrap[model.ScanRecord]
)
```

```go
// Step 2: In cli/store/repo_repository.go
// ✅ MANDATORY: Reuses ScanRecordWrap from types.go for 100% DRY code with value semantics
func (r *RepoRepository) FindBySlug(slug string) ScanRecordWrap {
	row := r.conn.QueryRow(SQLSelectBySlug, slug)
	var rec model.ScanRecord
	if err := row.Scan(&rec.RepoId, &rec.Slug, &rec.Branch); err != nil {
		return result.Fail[model.ScanRecord](apperror.WrapSimple(err, "FindBySlug"))
	}

	return result.Ok(rec)
}

// Caller consumption via affirmative predicates (IsFailure, IsSuccess)
res := repo.FindBySlug("gitmap")

// Affirmative failure guard:
if res.IsFailure() {
	return result.Fail[Output](res.AppError())
}

// Affirmative success branch:
if res.IsSuccess() {
	rec := res.Data()
	fmt.Println(rec.Slug)
}
```

---

### 3.2 Slice Collection Query (`ScanRecordSliceWrap`)

#### ❌ ANTI-PATTERN: The Screenshot Anti-Pattern (Package Stutter & Raw Generic Slice)
```go
// ❌ STRICTLY BANNED:
// 1. Package stutter: "result.Result" is banned; use "result.Wrap"
// 2. Leaking raw generic brackets in implementation signature
// 3. Wrapping an anonymous slice []model.ScanRecord instead of defining a dedicated model slice type
func (db *DB) ListReposByBranch(branch string) result.Result[[]model.ScanRecord] {
	rows, err := db.conn.Query(SQLSelectByBranch, branch)
	if err != nil {
		return result.Fail[[]model.ScanRecord](err) // ❌ Banned: bare error & unaliased generic
	}
	defer rows.Close()

	var repos []model.ScanRecord
	for rows.Next() {
		var rec model.ScanRecord
		_ = rows.Scan(&rec.RepoId, &rec.Slug) // ❌ Banned: orphaned error
		repos = append(repos, rec)
	}

	return result.Ok[[]model.ScanRecord](repos)
}

// ❌ BANNED: Outdated caller inspection
res := db.ListReposByBranch("main")
if res.HasError() { // ❌ Outdated: use affirmative IsFailure()
	return res.Err   // ❌ Banned: direct field access
}
for _, rec := range res.Value { // ❌ Banned: raw Value access
	process(rec)
}
```

#### ✅ MANDATORY: Dedicated Model Slice & Canonical Container in `types.go`

```go
// Step 1: In cli/store/types.go (in the package where it is defined for 100% DRY code)
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// ScanRecordSlice is the dedicated model slice collection type.
	ScanRecordSlice = []model.ScanRecord

	// ScanRecordSliceWrap is the canonical single reusable collection container.
	// Uses WrapSlice[model.ScanRecord] to provide slice monadic methods (.Items(), .Count()).
	ScanRecordSliceWrap = result.WrapSlice[model.ScanRecord]
)
```

```go
// Step 2: In cli/store/repo_repository.go
// ✅ MANDATORY: Reuses ScanRecordSliceWrap from types.go — clean, DRY single-type reuse
func (db *DB) ListReposByBranch(branch string) ScanRecordSliceWrap {
	rows, err := db.conn.Query(SQLSelectByBranch, branch)
	if err != nil {
		return result.FailSlice[model.ScanRecord](apperror.WrapSimple(err, "ListReposByBranch"))
	}
	defer rows.Close()

	records := make([]model.ScanRecord, 0)
	for rows.Next() {
		var rec model.ScanRecord
		if scanErr := rows.Scan(&rec.RepoId, &rec.Slug, &rec.Branch); scanErr != nil {
			return result.FailSlice[model.ScanRecord](apperror.WrapSimple(scanErr, "ListReposByBranch.Scan"))
		}
		records = append(records, rec)
	}

	return result.OkSlice(records)
}

// Caller consumption with affirmative predicates (IsFailure, IsSuccess)
res := db.ListReposByBranch("main")

// Affirmative failure guard:
if res.IsFailure() {
	return result.FailSlice[Output](res.AppError())
}

// Affirmative success branch:
if res.IsSuccess() {
	for _, rec := range res.Items() {
		processRecord(rec)
	}
}
```

---

### 3.3 Keyed Map Lookup (`ScanRecordMapWrap`)

```go
// Step 1: In cli/store/types.go (in the package where it is defined for 100% DRY code)
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// ScanRecordMap is the dedicated keyed lookup dictionary type.
	ScanRecordMap = map[string]model.ScanRecord

	// ScanRecordMapWrap is the canonical single reusable keyed container.
	// Uses WrapMap[string, model.ScanRecord] to provide map monadic methods (.Has(), .Get()).
	ScanRecordMapWrap = result.WrapMap[string, model.ScanRecord]
)
```

```go
// Step 2: In cli/store/repo_repository.go
// ✅ MANDATORY: Reuses ScanRecordMapWrap from types.go — clean, DRY single-type reuse
func (r *RepoRepository) IndexBySlug() ScanRecordMapWrap {
	recordsRes := r.ListDirty()
	if recordsRes.IsFailure() {
		return result.FailMap[string, model.ScanRecord](recordsRes.AppError())
	}

	lookup := make(ScanRecordMap, recordsRes.Count())
	for _, rec := range recordsRes.Items() {
		lookup[rec.Slug] = rec
	}

	return result.OkMap(lookup)
}

// Caller consumption with affirmative predicates (IsFailure, IsSuccess)
res := repo.IndexBySlug()

// Affirmative failure guard:
if res.IsFailure() {
	return result.Fail[Output](res.AppError())
}

// Affirmative success branch:
if res.IsSuccess() {
	if rec, isFound := res.Data()["gitmap"]; isFound {
		processRecord(rec)
	}
}
```

---

## 4. Verification Checklist

- [ ] All fallible Go functions return `result.Wrap[T]`, `result.WrapSlice[T]`, or `result.WrapMap[K, V]`.
- [ ] Zero usage of banned `result.Result` package stutter or `result.Result[[]model.ScanRecord]`.
- [ ] Dedicated model slice types (e.g. `ScanRecordSlice = []model.ScanRecord`) declared for all collection returns.
- [ ] Zero multi-value error tuples `(T, error)`, `([]T, error)`, `(map[K]V, error)` remain.
- [ ] All error instances are typed `*appfault.AppError` (or `*apperror.AppError`).
- [ ] All domain Result aliases are declared in `types.go` in the package where defined.
- [ ] All callers enforce the Result Guard Rule using affirmative predicates `if res.IsFailure()` / `if res.IsSuccess()`.
