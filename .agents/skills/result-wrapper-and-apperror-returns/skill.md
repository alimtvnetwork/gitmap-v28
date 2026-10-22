---
name: result-wrapper-and-apperror-returns
description: "Enforce enterprise Go single-return Result containers (result.Wrap[T], result.WrapSlice[T], result.WrapMap[K, V]), eliminating multi-value error tuples, centralizing types into types.go, and guaranteeing pointer-attached null safety with *apperror.AppError."
---

# Result Wrapper & AppError Return Sovereignty

Source prompt: `01-prompts/15-cg-execute/19-result-wrapper-and-apperror-returns.md`

## 1. Core Mandate

In enterprise Go architecture, functions must **NEVER** return multi-value error tuples such as `(T, error)`, `([]T, error)`, or `(map[K]V, error)`. Multi-value error returns violate single-return sovereignty, induce accidental `nil` dereferences, encourage swallowed errors, and cause dual-handling chaos.

All fallible operations MUST return a single strongly-typed monadic container from `cli/result` (or `pkg/appfault` / `pkg/result`):

| Container Type | Canonical Alias | Data Payload Representation | Constructor Helpers |
| :--- | :--- | :--- | :--- |
| **`result.Result[T]`** | **`result.Wrap[T]`** | Single scalar, model, or struct | `result.OkWrap(val)` / `result.FailWrap[T](err)` |
| **`result.ResultSlice[T]`** | **`result.WrapSlice[T]`** | Slice / collection of items | `result.OkWrapSlice(items)` / `result.FailWrapSlice[T](err)` |
| **`result.ResultMap[K, V]`** | **`result.WrapMap[K, V]`** | Keyed lookup dictionary or index | `result.OkWrapMap(data)` / `result.FailWrapMap[K, V](err)` |
| **`*apperror.AppError`** | **`*appfault.AppError`** | Structured enterprise error | `apperror.New(...)`, `apperror.WrapSimple(err, op)` |

---

## 2. Mandatory Rules

1. **Single Return Object:** Every function returning fallible data must return exactly one value: `result.Wrap[T]`, `result.WrapSlice[T]`, or `result.WrapMap[K, V]`.
2. **Zero Raw Error Returns:** Never return standard library `error`. All failures must wrap an underlying cause as `*apperror.AppError` (or `*appfault.AppError`).
3. **Pointer-Attached Null Safety:** All Result inspection methods (`IsSuccess`, `IsFailure`, `HasError`, `IsEmpty`, `HasRecord`, `IsCountOtherThan`, `Count`, `Items`) must attach to pointer receivers (`(r *Result[T])`, `(rs *ResultSlice[T])`, `(rm *ResultMap[K, V])`) with line-1 `if r == nil` guards returning safe defaults.
4. **Result Guard Rule:** Callers MUST check `if res.HasError()` (or `if res.IsFailure()`) before accessing payload fields (`res.Value`, `res.Data`, `res.Items()`). Accessing payload without an error check is a strict guideline violation.
5. **Centralization in `types.go`:** Never leak raw generic brackets like `result.Wrap[MyEntity]` or declare domain structs inline in implementation files. Every package MUST alias Result containers in a dedicated `types.go` file as single reusable named types:
   ```go
   type (
       ItemWrap      = result.Wrap[Item]
       ItemSliceWrap = result.WrapSlice[Item]
       ItemMapWrap   = result.WrapMap[string, Item]
   )
   ```

---

## 3. Concrete Implementation Examples

### Single Item Lookup (`result.Wrap[T]`)
```go
// Implementation in repository/service
func (r *RepoRepository) FindBySlug(slug string) result.Wrap[model.ScanRecord] {
	row := r.conn.QueryRow(SQLSelectBySlug, slug)
	var rec model.ScanRecord
	if err := row.Scan(&rec.RepoId, &rec.Slug, &rec.Branch); err != nil {
		return result.FailWrap[model.ScanRecord](apperror.WrapSimple(err, "FindBySlug"))
	}

	return result.OkWrap(rec)
}

// Caller consumption
res := repo.FindBySlug("gitmap")
if res.HasError() {
	return result.FailWrap[Output](res.Err)
}
if res.IsEmpty() {
	return result.FailWrap[Output](apperror.NewNotFound("slug not found"))
}
rec := res.Value
```

### Slice Collection Query (`result.WrapSlice[T]`)
```go
// Implementation in repository/service
func (r *RepoRepository) ListDirty() result.WrapSlice[model.ScanRecord] {
	rows, err := r.conn.Query(SQLSelectDirtyRepos)
	if err != nil {
		return result.FailWrapSlice[model.ScanRecord](apperror.WrapSimple(err, "ListDirty"))
	}
	defer rows.Close()

	records := make([]model.ScanRecord, 0)
	for rows.Next() {
		var rec model.ScanRecord
		if scanErr := rows.Scan(&rec.RepoId, &rec.Slug); scanErr != nil {
			return result.FailWrapSlice[model.ScanRecord](apperror.WrapSimple(scanErr, "ListDirty.Scan"))
		}
		records = append(records, rec)
	}

	return result.OkWrapSlice(records)
}

// Caller consumption
res := repo.ListDirty()
if res.HasError() {
	return result.FailWrapSlice[Output](res.Err)
}
for _, rec := range res.Items() {
	handleRecord(rec)
}
```

### Keyed Map Lookup (`result.WrapMap[K, V]`)
```go
// Implementation in repository/service
func (r *RepoRepository) IndexBySlug() result.WrapMap[string, model.ScanRecord] {
	recordsRes := r.ListAll()
	if recordsRes.HasError() {
		return result.FailWrapMap[string, model.ScanRecord](recordsRes.Err)
	}

	lookup := make(map[string]model.ScanRecord, recordsRes.Count())
	for _, rec := range recordsRes.Items() {
		lookup[rec.Slug] = rec
	}

	return result.OkWrapMap(lookup)
}

// Caller consumption
res := repo.IndexBySlug()
if res.HasError() {
	return result.FailWrap[Output](res.Err)
}
if res.Has("gitmap") {
	rec, isFound := res.Get("gitmap")
}
```

---

## 4. Verification Checklist

- [ ] All fallible Go functions return `result.Wrap[T]`, `result.WrapSlice[T]`, or `result.WrapMap[K, V]`.
- [ ] Zero multi-value error tuples `(T, error)`, `([]T, error)`, `(map[K]V, error)` remain.
- [ ] All error instances are typed `*apperror.AppError` (or `*appfault.AppError`).
- [ ] All domain Result aliases are declared in `types.go`.
- [ ] All callers enforce the Result Guard Rule with `if res.HasError()`.
