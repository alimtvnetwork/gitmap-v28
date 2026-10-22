# Subtask 08: Modernize System Domain Types & Eliminate Result Stutter in OSUser, OSClean, and CommitTransfer

> **Subtask ID:** `SUBTASK-296-08`  
> **Parent Task:** `296-result-wrapper-and-apperror-returns`  
> **Wave:** Wave 3 (Stutter Elimination, FastGate Verification & Release)  
> **Status:** Completed  
> **Target Files:**  
> - `cli/osuser/types.go` (modify)  
> - `cli/osclean/types.go` (modify)  
> - `cli/committransfer/types.go` (modify)  
> **Related Architecture Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/01-architecture-spec.md`  
> **Related Component Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/02-component-spec.md`  
> **Related Migration Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/03-migration-and-verification-spec.md`  

---

## 1. Objective

Modernize OS and transfer domain types across `cli/osuser/types.go`, `cli/osclean/types.go`, and `cli/committransfer/types.go` by eliminating legacy package stutter (`result.Result[...]`) and introducing canonical monadic containers (`result.Wrap[...]`). Standardize type aliases ending in `Wrap` while preserving backward compatibility for existing callers.

---

## 2. Action Checklist

- [x] 1. Modernize `cli/osuser/types.go`:
  - [x] Update type alias block lines 67–79:
    - [x] Introduce canonical alias: `UserWrap = result.Wrap[UserInfo]`
    - [x] Forward legacy alias: `UserResult = UserWrap`
    - [x] Ensure scalar result aliases are clean:
      - `UserBoolResult = result.BoolResult`
      - `UserStringResult = result.StringResult`
      - `UserStringsResult = result.StringSliceResult`
  - [x] Enforce Go coding standards (blank line after `}`, doc comments for exported types).

- [x] 2. Modernize `cli/osclean/types.go`:
  - [x] Update line 23:
    - [x] Replace `CleanResult = result.Result[CleanStats]` with canonical `CleanWrap = result.Wrap[CleanStats]`
    - [x] Add backward-compatibility alias: `CleanResult = CleanWrap`
  - [x] Enforce Go formatting rules.

- [x] 3. Modernize `cli/committransfer/types.go`:
  - [x] Import `github.com/alimtvnetwork/gitmap-v28/cli/result`.
  - [x] Define canonical single-return monadic Result wrappers for replay operations:
    - [x] `ReplayResultWrap = result.Wrap[ReplayResult]`
    - [x] `ReplayPlanWrap = result.Wrap[ReplayPlan]`
  - [x] Document the wrappers for eliminating multi-value error tuples in commit replay flows.

- [x] 4. Invariant Verification:
  - [x] Strictly relative Git paths only.
  - [x] Zero compiler builds (`go build` strictly banned).
  - [x] Zero tests (`go test` strictly banned).
  - [x] Zero git commands.

---

## 3. Concrete Code Specifications

### 3.1 `cli/osuser/types.go`

```go
type (
	// UserWrap is the canonical single reusable result container for UserInfo (Task 296).
	UserWrap = result.Wrap[UserInfo]

	// UserResult is preserved as a backward-compatibility alias for UserWrap.
	UserResult = UserWrap

	// UserBoolResult encapsulates a boolean outcome.
	UserBoolResult = result.BoolResult

	// UserStringResult encapsulates a string outcome.
	UserStringResult = result.StringResult

	// UserStringsResult encapsulates a slice of strings outcome.
	UserStringsResult = result.StringSliceResult
)
```

---

### 3.2 `cli/osclean/types.go`

```go
type (
	// CleanWrap is the canonical single reusable result container for CleanStats (Task 296).
	CleanWrap = result.Wrap[CleanStats]

	// CleanResult is preserved as a backward-compatibility alias for CleanWrap.
	CleanResult = CleanWrap
)
```

---

### 3.3 `cli/committransfer/types.go`

```go
import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// ReplayResultWrap is the canonical single reusable result container for ReplayResult (Task 296).
	ReplayResultWrap = result.Wrap[ReplayResult]

	// ReplayPlanWrap is the canonical single reusable result container for ReplayPlan (Task 296).
	ReplayPlanWrap = result.Wrap[ReplayPlan]
)
```
