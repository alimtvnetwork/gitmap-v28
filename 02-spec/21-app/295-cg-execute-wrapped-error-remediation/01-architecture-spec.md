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
> **Single Return Value Mandate:** Methods that perform database queries, calculations, or lookups MUST return a single monadic container (`result.Result[T]`, `store.RepoResult[T]`, `store.RepoListResult[T]`) rather than multi-value tuples `(T, error)` or `(T, *apperror.AppError)`.

This specification formalizes the transformation of `cli/store/repo_repository.go`, `cli/store/repo_query_filter.go`, `cli/store/repo.go`, and dependent pipeline query renderers to use `result.Result[T]` containers.

---

## 2. Monadic Container Contracts

### 2.1 The `result.Result[T]` Container
Defined in `cli/result/types.go` and `cli/result/result.go`:
```go
type Result[T any] struct {
    Value     T
    Data      T
    Err       *apperror.AppError
    isDefined bool
}
```

Key constructors and receiver methods:
- `result.Ok(val)`: Constructs a successful envelope with `.Value` populated and `.isDefined = true`.
- `result.Fail[T](appErr)`: Constructs a failed envelope with `.Err = appErr`.
- `res.IsSuccess() bool`: Returns `true` if `res.Err == nil || res.Err.HasNoError()`.
- `res.HasError() bool`: Returns `true` if `res.Err != nil && res.Err.HasError()`.
- `res.Value`: Direct access to payload after `!res.HasError()` guard.
- `res.Err`: Direct access to `*apperror.AppError`.

### 2.2 Repository Method Signatures
All methods of `RepoRepository` and corresponding delegating methods on `*store.DB` are refactored from tuple returns to single `result.Result[T]` returns:

| Prior Signature | Monadic Wrapped Signature |
| :--- | :--- |
| `ListByBranch(branch string) ([]model.ScanRecord, *apperror.AppError)` | `ListByBranch(branch string) result.Result[[]model.ScanRecord]` |
| `ListDirty() ([]model.ScanRecord, *apperror.AppError)` | `ListDirty() result.Result[[]model.ScanRecord]` |
| `ListByTag(tag string) ([]model.ScanRecord, *apperror.AppError)` | `ListByTag(tag string) result.Result[[]model.ScanRecord]` |
| `ListUnpushed() ([]model.ScanRecord, *apperror.AppError)` | `ListUnpushed() result.Result[[]model.ScanRecord]` |
| `FindPathBySlug(slug string) (string, *apperror.AppError)` | `FindPathBySlug(slug string) result.Result[string]` |

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
