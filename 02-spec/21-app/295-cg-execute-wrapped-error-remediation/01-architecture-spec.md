# Architecture Specification 295: Repository-Wide Monadic Wrapped Result Architecture (CG Execute)

## Status: APPROVED
- **Spec ID:** 295-01
- **Slug:** `295-cg-execute-wrapped-error-remediation`
- **Domain:** Error Sovereignty & Monadic Wrapped Result Design
- **Cross-Spec References:**
  - `02-spec/03-error-manage/`
  - `02-spec/02-coding-guidelines/01-cross-language/13-strict-typing.md`
  - `02-spec/02-coding-guidelines/01-cross-language/40-database-repository-pattern-and-error-sovereignty.md`
- **Invariants:** Zero-Build, Zero-Test, Strictly Relative Lowercase Git Paths, Bounded 5–8 File Micro-Batches

---

## 1. Executive Summary & Architectural Motivation

In the recent implementation of the store repository pattern (`cli/store/repo_repository.go`), database query methods returned two-value tuples: `([]model.ScanRecord, *apperror.AppError)` and `(string, *apperror.AppError)`.

While substituting raw `error` with strongly-typed `*apperror.AppError` resolved untyped error risks, returning multi-value error tuples directly violates the repository coding standard:
> **Single Return Value & DRY types.go Mandate:** Methods that perform database queries, calculations, or lookups MUST return a single monadic container (`result.Wrap[T]`, `result.WrapSlice[T]`, `result.WrapMap[K, V]`) rather than multi-value tuples `(T, error)` or `(T, *apperror.AppError)`. All domain result envelopes MUST be declared as single named type aliases in the package's `types.go` file (e.g. `ScanRecordSliceWrap`, `ScanRecordWrap`, `RepoPathWrap`) to guarantee 100% DRY code reuse across implementations and callers without leaking generic brackets.

This specification formalizes the transformation of `cli/store/repo_repository.go`, `cli/store/repo_query_filter.go`, `cli/store/repo.go`, and dependent pipeline query renderers to use `ScanRecordSliceWrap`, `ScanRecordWrap`, and `RepoPathWrap` containers.

---

## 2. Monadic Container Contracts & DRY types.go Standard

### 2.1 The `result.Wrap[T]` and `result.WrapSlice[T]` Containers
Defined in `cli/result/types.go`, `cli/result/result.go`, and `cli/result/result_slice.go`:
```go
type Result[T any] struct {
    Value     T
    Data      T
    Err       *apperror.AppError
    isDefined bool
}

type ResultSlice[T any] struct {
    Value []T
    Data  []T
    Err   *apperror.AppError
}

// Canonical Aliases
type Wrap[T any] = Result[T]
type WrapSlice[T any] = ResultSlice[T]
```

### 2.2 Package-Level Type Aliases in `cli/store/repository_types.go`
```go
type (
    // ScanRecordWrap is the canonical single reusable result container for ScanRecord.
    ScanRecordWrap = result.Wrap[model.ScanRecord]

    // ScanRecordSliceWrap is the canonical single reusable collection container for ScanRecord.
    ScanRecordSliceWrap = result.WrapSlice[model.ScanRecord]

    // ScanRecordMapWrap is the canonical single reusable keyed lookup container for ScanRecord.
    ScanRecordMapWrap = result.WrapMap[string, model.ScanRecord]

    // RepoPathWrap is the canonical single reusable result container for repository paths.
    RepoPathWrap = result.Wrap[string]
)
```

Key constructors and receiver methods:
- `result.OkWrap(val)` / `result.OkWrapSlice(items)`: Constructs a successful envelope with `.Value` populated.
- `result.Fail[T](appErr)` / `result.FailSlice[T](appErr)`: Constructs a failed envelope with `.Err = appErr`. Eliminates redundant `FailWrap` / `FailWrapSlice` double-wrap calls.
- `res.IsSuccess() bool`: Returns `true` if `res.Err == nil || res.Err.HasNoError()`.
- `res.HasError() bool`: Returns `true` if `res.Err != nil && res.Err.HasError()`.
- `res.Value` / `res.Items()`: Direct access to payload after `!res.HasError()` guard.
- `res.Err`: Direct access to `*apperror.AppError`.

### 2.3 Repository Method Signatures
All methods of `RepoRepository` and corresponding delegating methods on `*store.DB` are refactored from tuple returns to single DRY wrapped returns:

| Prior Signature | Monadic Wrapped Signature |
| :--- | :--- |
| `ListByBranch(branch string) ([]model.ScanRecord, *apperror.AppError)` | `ListByBranch(branch string) ScanRecordSliceWrap` |
| `ListDirty() ([]model.ScanRecord, *apperror.AppError)` | `ListDirty() ScanRecordSliceWrap` |
| `ListByTag(tag string) ([]model.ScanRecord, *apperror.AppError)` | `ListByTag(tag string) ScanRecordSliceWrap` |
| `ListUnpushed() ([]model.ScanRecord, *apperror.AppError)` | `ListUnpushed() ScanRecordSliceWrap` |
| `FindPathBySlug(slug string) (string, *apperror.AppError)` | `FindPathBySlug(slug string) RepoPathWrap` |
| `ListAll() ([]model.ScanRecord, error)` | `ListAll() ScanRecordSliceWrap` |
| `FindBySlug(slug string) (*model.ScanRecord, error)` | `FindBySlug(slug string) ScanRecordWrap` |

---

## 3. Caller Consumption Pattern

Callers inspect `res.HasError()` via guard clauses before accessing payload values:

```go
repo := store.NewRepoRepository(db.Conn())
pathRes := repo.FindPathBySlug(slug)
if pathRes.HasError() {
    logRepoLookupFailure(slug, pathRes.Err)

    return ""
}

return pathRes.Value
```

This ensures complete leaf error sovereignty without unhandled nil pointer panics or tuple destructuring boilerplate.
