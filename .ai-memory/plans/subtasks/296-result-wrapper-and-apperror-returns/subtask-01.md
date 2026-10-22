# Subtask 01: Centralize Store Domain Types in types.go & Eradicate Stutter Anti-Patterns

> **Subtask ID:** `SUBTASK-296-01`  
> **Parent Task:** `296-result-wrapper-and-apperror-returns`  
> **Status:** Pending  
> **Target Files:**  
> - `cli/store/types.go` (new file)  
> - `cli/store/repository_types.go` (modify)  
> **Related Architecture Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/01-architecture-spec.md`  

---

## 1. Objective

Create `cli/store/types.go` to centralize all repository and heatmap store domain model slice collections and monadic Result wrapper aliases. Eliminate the banned anti-pattern `ScanRecordsResult = result.Result[[]model.ScanRecord]` in `cli/store/repository_types.go`, removing package stutter and raw generic bracket leakage.

---

## 2. Action Checklist

- [ ] 1. Author new file `cli/store/types.go`:
  - [ ] Package declaration `package store`.
  - [ ] Import `github.com/alimtvnetwork/gitmap-v28/cli/model` and `github.com/alimtvnetwork/gitmap-v28/cli/result`.
  - [ ] Define the 13 required canonical type aliases:
    - [ ] `ScanRecordSlice = []model.ScanRecord`
    - [ ] `ScanRecordWrap = result.Wrap[model.ScanRecord]`
    - [ ] `ScanRecordSliceWrap = result.WrapSlice[model.ScanRecord]`
    - [ ] `ScanRecordMapWrap = result.WrapMap[string, model.ScanRecord]`
    - [ ] `HeatmapActionLogRecordSliceWrap = result.WrapSlice[HeatmapActionLogRecord]`
    - [ ] `HeatmapExecutionTreeNodeSliceWrap = result.WrapSlice[HeatmapExecutionTreeNode]`
    - [ ] `HeatmapStatusSummaryWrap = result.Wrap[HeatmapStatusSummary]`
    - [ ] `HeatmapTaskSummaryWrap = result.Wrap[HeatmapTaskSummary]`
    - [ ] `FileTouchCountSliceWrap = result.WrapSlice[FileTouchCount]`
    - [ ] `DuplicateRepoGroupSliceWrap = result.WrapSlice[DuplicateRepoGroup]`
    - [ ] `RowsAffectedWrap = result.Wrap[int64]`
    - [ ] `RepoIDWrap = result.Wrap[int64]`
    - [ ] `RepoCountWrap = result.Wrap[int]`
    - [ ] `RepoPathWrap = result.Wrap[string]`
  - [ ] Ensure formatting follows strict Go guidelines (blank line after `}`, docstrings for all exported aliases).

- [ ] 2. Clean up `cli/store/repository_types.go`:
  - [ ] Remove duplicate type aliases lines 14-29:
    - Remove `ScanRecordWrap = result.Wrap[model.ScanRecord]`
    - Remove `ScanRecordSliceWrap = result.WrapSlice[model.ScanRecord]`
    - Remove `ScanRecordMapWrap = result.WrapMap[string, model.ScanRecord]`
    - Remove `RepoPathWrap = result.Wrap[string]`
    - **CRITICAL:** Remove banned line 28 `ScanRecordsResult = result.Result[[]model.ScanRecord]` and its doc comment.
  - [ ] Retain repository interfaces `IRepository`, `IPromptBackupRepository`, `PromptBackupBatch`, etc.
  - [ ] Verify that unused imports (e.g. `github.com/alimtvnetwork/gitmap-v28/cli/result` if no other types in `repository_types.go` use it) are cleaned up if applicable.

- [ ] 3. Invariant Verification:
  - [ ] Strictly relative Git paths only.
  - [ ] Zero compiler builds (`go build` strictly banned).
  - [ ] Zero tests (`go test` strictly banned).
  - [ ] Zero git commands.

---

## 3. Concrete Code Specifications

### Specification: `cli/store/types.go`

```go
// Package store — types.go defines centralized domain slice collections and canonical Result wrappers.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// ScanRecordSlice is the canonical dedicated model slice collection type.
	ScanRecordSlice = []model.ScanRecord

	// ScanRecordWrap is the canonical single reusable result container for ScanRecord.
	ScanRecordWrap = result.Wrap[model.ScanRecord]

	// ScanRecordSliceWrap is the canonical single reusable collection container for ScanRecord.
	ScanRecordSliceWrap = result.WrapSlice[model.ScanRecord]

	// ScanRecordMapWrap is the canonical single reusable keyed lookup container for ScanRecord.
	ScanRecordMapWrap = result.WrapMap[string, model.ScanRecord]

	// RepoPathWrap is the canonical single reusable result container for repository paths.
	RepoPathWrap = result.Wrap[string]

	// HeatmapActionLogRecordSliceWrap is the canonical single reusable collection container for HeatmapActionLogRecord.
	HeatmapActionLogRecordSliceWrap = result.WrapSlice[HeatmapActionLogRecord]

	// HeatmapExecutionTreeNodeSliceWrap is the canonical single reusable collection container for HeatmapExecutionTreeNode.
	HeatmapExecutionTreeNodeSliceWrap = result.WrapSlice[HeatmapExecutionTreeNode]

	// HeatmapStatusSummaryWrap is the canonical single reusable result container for HeatmapStatusSummary.
	HeatmapStatusSummaryWrap = result.Wrap[HeatmapStatusSummary]

	// HeatmapTaskSummaryWrap is the canonical single reusable result container for HeatmapTaskSummary.
	HeatmapTaskSummaryWrap = result.Wrap[HeatmapTaskSummary]

	// FileTouchCountSliceWrap is the canonical single reusable collection container for FileTouchCount.
	FileTouchCountSliceWrap = result.WrapSlice[FileTouchCount]

	// DuplicateRepoGroupSliceWrap is the canonical single reusable collection container for DuplicateRepoGroup.
	DuplicateRepoGroupSliceWrap = result.WrapSlice[DuplicateRepoGroup]

	// RowsAffectedWrap is the canonical single reusable result container for database rows affected.
	RowsAffectedWrap = result.Wrap[int64]

	// RepoIDWrap is the canonical single reusable result container for repository database IDs.
	RepoIDWrap = result.Wrap[int64]

	// RepoCountWrap is the canonical single reusable result container for repository counts.
	RepoCountWrap = result.Wrap[int]
)
```

### Remediation in `cli/store/repository_types.go`

Lines 14-29 currently read:
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

	// ScanRecordsResult is the canonical single reusable result container for []model.ScanRecord.
	ScanRecordsResult = result.Result[[]model.ScanRecord]
)
```

Remediation:
Remove the entire redundant block above from `cli/store/repository_types.go` because all aliases are now cleanly defined in `cli/store/types.go`, completely eliminating line 28 `ScanRecordsResult = result.Result[[]model.ScanRecord]`.

---

## 4. Verification & Done Criteria

- [ ] `cli/store/types.go` exists with package `store` and all 13 required aliases.
- [ ] Line 28 of `cli/store/repository_types.go` containing `ScanRecordsResult` is deleted.
- [ ] No compilation or test commands were run.
- [ ] Zero absolute paths or file URIs used.
