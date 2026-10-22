# Subtask 06: Author cli/cmdpipeline/types.go and Refactor Callers in pipeline_history_render.go & pipeline_fallback.go

## Status: Completed
- **Parent Plan:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/01-architecture-spec.md`
- **Component Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/02-component-spec.md`
- **Target Files:**
  - `cli/cmdpipeline/types.go`
  - `cli/cmdpipeline/pipeline_history_render.go`
  - `cli/cmdpipeline/pipeline_fallback.go`
- **Objective:** Establish domain collection type aliases in `cli/cmdpipeline/types.go` (`PEFormatProfileSliceWrap`, `PipelineSyncResultWrap`) and modernize caller sites across `pipeline_history_render.go` and `pipeline_fallback.go` to use affirmative predicates (`res.IsFailure()`, `res.Items()`, `res.HasRecord()`, `res.IsEmpty()`).

---

## Action Items

1. **Author `cli/cmdpipeline/types.go`:**
   - Define dedicated slice collection type: `type PEFormatProfileSlice = []PEFormatProfile`.
   - Alias canonical single reusable collection container: `type PEFormatProfileSliceWrap = result.Wrap[PEFormatProfileSlice]`.
   - Alias pointer container: `type PipelineSyncResultWrap = result.Wrap[*PipelineSyncResult]`.
   - Prevent package stutter and eliminate raw generic bracket leakage across pipeline public APIs.

2. **Refactor Caller Sites in `cli/cmdpipeline/pipeline_history_render.go`:**
   - In `findLocalRepoDirForSlug(slug string)`:
     - Replace obsolete `if pathRes.HasError()` check with affirmative `if pathRes.IsFailure()`.
     - Log failure using `pathRes.AppError()`.
     - Return path safely via `pathRes.Data` (or `pathRes.UnwrapOr("")`).
   - In `renderDetailedOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64)`:
     - Check `if detailRes.HasRecord()`.
     - Pass `detailRes.Items()` and `legacyRes.Items()` to `renderCachedErrorsList` instead of accessing raw struct field `Data`.
   - In `renderCompactOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64)`:
     - Check `if compactRes.HasRecord()`.
     - Pass `compactRes.Items()` to `renderCachedCompactErrorsList` and `legacyRes.Items()` to `renderCachedErrorsList`.
   - In `queryCachedFailedRuns(db *pipelinedb.PipelineSplitDb, count int)`:
     - Guard on `if runRes.IsFailure()`, returning `nil, runRes.AppError()`.
     - Return slice using `runRes.Items(), nil`.

3. **Refactor Caller Sites in `cli/cmdpipeline/pipeline_fallback.go`:**
   - In `TryRetrievePreviousDbRun(repo string)`:
     - Replace `if runRes.IsFailed() || len(runRes.Data) == 0` with affirmative `if runRes.IsFailure() || runRes.IsEmpty()`.
     - Retrieve record safely using `runRes.Items()[0]`.
   - In `queryDbRunErrorText(db *pipelinedb.PipelineSplitDb, runId uint64)`:
     - Invert negative boolean check `!errRes.IsFailed() && len(errRes.Data) > 0` to affirmative `if errRes.HasRecord()`.
     - Pass `errRes.Items()` to `formatDbErrorRecords`.
     - Invert negative boolean check `!compactRes.IsFailed() && len(compactRes.Data) > 0` to affirmative `if compactRes.HasRecord()`.
     - Pass `compactRes.Items()` to `formatDbCompactRecords`.

4. **Verify Coding Invariants:**
   - Ensure every refactored function stays $\le 15$ lines.
   - Maintain zero nested `if` statements using guard clauses.
   - Enforce vertical whitespace rules: blank line before `return` and after `}`.
   - Eliminate all negative boolean checks (`!is...`, `!has...`, `!res.IsFailed()`).
   - Zero compiler builds or unit tests during authoring.

---

## Implementation Details

### 1. `cli/cmdpipeline/types.go`
```go
package cmdpipeline

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// PEFormatProfileSlice defines a collection of pipeline format profiles.
	PEFormatProfileSlice = []PEFormatProfile

	// PEFormatProfileSliceWrap is the canonical single reusable container for PEFormatProfile collections.
	PEFormatProfileSliceWrap = result.Wrap[PEFormatProfileSlice]

	// PipelineSyncResultWrap is the canonical single reusable container for PipelineSyncResult pointers.
	PipelineSyncResultWrap = result.Wrap[*PipelineSyncResult]
)
```

### 2. `cli/cmdpipeline/pipeline_history_render.go` (Surgical Updates)
```go
func findLocalRepoDirForSlug(slug string) string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	repo := store.NewRepoRepository(db.Conn())
	pathRes := repo.FindPathBySlug(slug)
	if pathRes.IsFailure() {
		logRepoLookupFailure(slug, pathRes.AppError())

		return ""
	}

	return pathRes.Data
}

func renderDetailedOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64) {
	detailRes := db.QueryDetailedErrorLogsByRunId(runId)
	if detailRes.HasRecord() {
		renderCachedErrorsList(detailRes.Items())

		return
	}

	legacyRes := db.QueryErrorLogsByRunId(runId)
	renderCachedErrorsList(legacyRes.Items())
}

func renderCompactOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64) {
	compactRes := db.QueryCompactErrorLogsByRunId(runId)
	if compactRes.HasRecord() {
		renderCachedCompactErrorsList(compactRes.Items())

		return
	}

	legacyRes := db.QueryErrorLogsByRunId(runId)
	renderCachedErrorsList(legacyRes.Items())
}

func queryCachedFailedRuns(db *pipelinedb.PipelineSplitDb, count int) ([]pipelinedb.PipelineRunRecord, error) {
	limit := resolveMaxSyncLimit(count)
	runRes := db.QueryLastFailedRuns(limit)
	if runRes.IsFailure() {
		return nil, runRes.AppError()
	}

	return runRes.Items(), nil
}
```

### 3. `cli/cmdpipeline/pipeline_fallback.go` (Surgical Updates)
```go
// TryRetrievePreviousDbRun retrieves the most recent failed run from pipelinedb.
func TryRetrievePreviousDbRun(repo string) (*pipelinedb.PipelineRunRecord, string, bool) {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil, "", false
	}
	defer db.Close()

	runRes := db.QueryLastFailedRuns(1)
	if runRes.IsFailure() || runRes.IsEmpty() {
		return nil, "", false
	}

	items := runRes.Items()
	rec := items[0]
	errLogs := queryDbRunErrorText(db, rec.RunId)

	return &rec, errLogs, true
}

func queryDbRunErrorText(db *pipelinedb.PipelineSplitDb, runId uint64) string {
	errRes := db.QueryDetailedErrorLogsByRunId(runId)
	if errRes.HasRecord() {
		return formatDbErrorRecords(errRes.Items())
	}

	compactRes := db.QueryCompactErrors(10)
	if compactRes.HasRecord() {
		return formatDbCompactRecords(compactRes.Items())
	}

	return ""
}
```

---

## Verification & Conformance Checklist

- [ ] `cli/cmdpipeline/types.go` created with `PEFormatProfileSliceWrap` and `PipelineSyncResultWrap`.
- [ ] No raw generic brackets leaked across `cmdpipeline` types.
- [ ] Replaced obsolete `pathRes.HasError()` with affirmative `pathRes.IsFailure()`.
- [ ] Replaced negative booleans `!errRes.IsFailed()` and `!compactRes.IsFailed()` with affirmative `HasRecord()`.
- [ ] Replaced raw slice access with `res.Items()`.
- [ ] All updated functions $\le 15$ lines with zero nested `if` statements.
- [ ] Zero compiler builds and zero unit tests executed during plan creation.
- [ ] Strictly relative Git paths utilized throughout all references.
