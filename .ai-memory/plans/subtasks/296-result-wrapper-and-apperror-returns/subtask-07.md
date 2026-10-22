# Subtask 07: Modernize Database Domain Types & Eliminate Result Stutter in PRDB, PipelineDB, and DBEngine

> **Subtask ID:** `SUBTASK-296-07`  
> **Parent Task:** `296-result-wrapper-and-apperror-returns`  
> **Wave:** Wave 3 (Stutter Elimination, FastGate Verification & Release)  
> **Status:** Pending  
> **Target Files:**  
> - `cli/prdb/types.go` (modify)  
> - `cli/pipelinedb/types.go` (modify)  
> - `cli/dbengine/result_types.go` (modify)  
> **Related Architecture Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/01-architecture-spec.md`  
> **Related Component Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/02-component-spec.md`  
> **Related Migration Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/03-migration-and-verification-spec.md`  

---

## 1. Objective

Modernize database domain types across `cli/prdb/types.go`, `cli/pipelinedb/types.go`, and `cli/dbengine/result_types.go` by replacing legacy package stutter (`result.Result[...]`) and raw slice generic brackets (`result.Result[[]T]`) with canonical monadic container aliases (`result.Wrap[...]`, `result.WrapSlice[...]`, and `result.WrapMap[...]`). Maintain 100% backwards compatibility for existing callers through type alias forwarding.

---

## 2. Action Checklist

- [ ] 1. Modernize `cli/prdb/types.go`:
  - [ ] Replace package stutter in lines 52–57:
    - [ ] `PullRequestWrap = result.Wrap[*PullRequestRecord]`
    - [ ] `PrBranchSliceWrap = result.WrapSlice[PrBranchRecord]`
    - [ ] `PrReleaseWrap = result.Wrap[*PrReleaseRecord]`
  - [ ] Preserve backward-compatibility aliases:
    - [ ] `PullRequestResult = PullRequestWrap`
    - [ ] `PrBranchListResult = PrBranchSliceWrap`
    - [ ] `PrReleaseResult = PrReleaseWrap`
  - [ ] Ensure formatting strictly complies with Go conventions (blank line after `}`, docstrings for exported types).

- [ ] 2. Modernize `cli/pipelinedb/types.go`:
  - [ ] Modernize type alias block lines 8–29:
    - [ ] `PipelineRunSliceWrap = result.WrapSlice[PipelineRunRecord]`
    - [ ] `PipelineErrorSliceWrap = result.WrapSlice[PipelineErrorRecord]`
    - [ ] `PipelineCompactErrorSliceWrap = result.WrapSlice[PipelineCompactErrorRecord]`
    - [ ] `PipelineRunIdSliceWrap = result.WrapSlice[uint64]`
    - [ ] `PipelineRunIdMapWrap = result.WrapMap[uint64, bool]`
    - [ ] `PipelineDbStatsWrap = result.Wrap[PipelineDbStats]`
    - [ ] `PipelineDatabaseInfoWrap = result.Wrap[PipelineDatabaseInfo]`
  - [ ] Preserve backward-compatibility aliases pointing to the new `*Wrap` definitions:
    - [ ] `PipelineRunSliceResult = PipelineRunSliceWrap`
    - [ ] `PipelineErrorSliceResult = PipelineErrorSliceWrap`
    - [ ] `PipelineCompactErrorSliceResult = PipelineCompactErrorSliceWrap`
    - [ ] `PipelineRunIdSliceResult = PipelineRunIdSliceWrap`
    - [ ] `PipelineRunIdMapResult = PipelineRunIdMapWrap`
    - [ ] `PipelineDbStatsResult = PipelineDbStatsWrap`
    - [ ] `PipelineDatabaseInfoResult = PipelineDatabaseInfoWrap`

- [ ] 3. Modernize `cli/dbengine/result_types.go`:
  - [ ] Eliminate package stutter and raw generic slice brackets in lines 14–16:
    - [ ] `type EntityResult[T any] = result.Wrap[*T]`
    - [ ] `type ListResult[T any] = result.WrapSlice[T]`
    - [ ] `type CompiledQueryWrap = result.Wrap[CompiledQuery]`
    - [ ] `type CompiledQueryResult = CompiledQueryWrap`
  - [ ] Modernize constructor helper functions lines 67–100:
    - [ ] `SuccessEntity[T any](val *T) EntityResult[T]` -> delegates to `result.Ok(val)`
    - [ ] `FailureEntity[T any](err *apperror.AppError) EntityResult[T]` -> delegates to `result.Fail[*T](err)`
    - [ ] `SuccessList[T any](val []T) ListResult[T]` -> delegates to `result.OkSlice(val)`
    - [ ] `FailureList[T any](err *apperror.AppError) ListResult[T]` -> delegates to `result.FailSlice[T](err)`
    - [ ] `SuccessCompiledQuery(val CompiledQuery) CompiledQueryResult` -> delegates to `result.Ok(val)`

- [ ] 4. Invariant Verification:
  - [ ] Strictly relative Git paths only.
  - [ ] Zero compiler builds (`go build` strictly banned).
  - [ ] Zero tests (`go test` strictly banned).
  - [ ] Zero git commands.

---

## 3. Concrete Code Specifications

### 3.1 `cli/prdb/types.go`

```go
// Monadic Result wrappers for PR database operations (Spec 129 / Task 296).
type (
	// PullRequestWrap is the canonical single reusable result container for PullRequestRecord.
	PullRequestWrap = result.Wrap[*PullRequestRecord]

	// PrBranchSliceWrap is the canonical single reusable collection container for PrBranchRecord.
	PrBranchSliceWrap = result.WrapSlice[PrBranchRecord]

	// PrReleaseWrap is the canonical single reusable result container for PrReleaseRecord.
	PrReleaseWrap = result.Wrap[*PrReleaseRecord]

	// Backward-compatibility aliases for legacy callers.
	PullRequestResult  = PullRequestWrap
	PrBranchListResult = PrBranchSliceWrap
	PrReleaseResult    = PrReleaseWrap
)
```

---

### 3.2 `cli/pipelinedb/types.go`

```go
type (
	// Canonical single reusable monadic containers for pipeline records (Task 296).
	PipelineRunSliceWrap          = result.WrapSlice[PipelineRunRecord]
	PipelineErrorSliceWrap        = result.WrapSlice[PipelineErrorRecord]
	PipelineCompactErrorSliceWrap = result.WrapSlice[PipelineCompactErrorRecord]
	PipelineRunIdSliceWrap        = result.WrapSlice[uint64]
	PipelineRunIdMapWrap          = result.WrapMap[uint64, bool]
	PipelineDbStatsWrap           = result.Wrap[PipelineDbStats]
	PipelineDatabaseInfoWrap      = result.Wrap[PipelineDatabaseInfo]

	// Backward-compatibility aliases for legacy callers.
	PipelineRunSliceResult          = PipelineRunSliceWrap
	PipelineErrorSliceResult        = PipelineErrorSliceWrap
	PipelineCompactErrorSliceResult = PipelineCompactErrorSliceWrap
	PipelineRunIdSliceResult        = PipelineRunIdSliceWrap
	PipelineRunIdMapResult          = PipelineRunIdMapWrap
	PipelineDbStatsResult           = PipelineDbStatsWrap
	PipelineDatabaseInfoResult      = PipelineDatabaseInfoWrap
)
```

---

### 3.3 `cli/dbengine/result_types.go`

```go
// Exact typed result envelopes wrapping canonical result types.
type (
	Uint64Result       = result.Uint64Result
	Int64Result        = result.Int64Result
	StringResult       = result.StringResult
	BoolResult         = result.BoolResult
	RowsAffectedResult = result.Int64Result
	EntityResult[T any] = result.Wrap[*T]
	ListResult[T any]   = result.WrapSlice[T]
	CompiledQueryWrap  = result.Wrap[CompiledQuery]
	CompiledQueryResult = CompiledQueryWrap
)

// SuccessEntity wraps an entity pointer in a successful result envelope.
func SuccessEntity[T any](val *T) EntityResult[T] {
	return result.Ok(val)
}

// FailureEntity wraps an AppError in a failed EntityResult envelope.
func FailureEntity[T any](err *apperror.AppError) EntityResult[T] {
	return result.Fail[*T](err)
}

// SuccessList wraps an entity slice in a successful result envelope.
func SuccessList[T any](val []T) ListResult[T] {
	return result.OkSlice(val)
}

// FailureList wraps an AppError in a failed ListResult envelope.
func FailureList[T any](err *apperror.AppError) ListResult[T] {
	return result.FailSlice[T](err)
}

// SuccessCompiledQuery wraps a CompiledQuery in a successful result envelope.
func SuccessCompiledQuery(val CompiledQuery) CompiledQueryResult {
	return result.Ok(val)
}
```
