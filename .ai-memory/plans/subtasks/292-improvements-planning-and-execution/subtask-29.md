# Subtask Plan 29: Scope Resolution & DB Fallback Extraction (`cli/cmdstatus/status_resolve.go`)

> **Subtask ID:** `subtask-29`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdstatus/status_resolve.go`, `cli/cmdstatus/status.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Extract repository scope resolution, database query fallbacks, and JSON manifest loading from `cli/cmdstatus/status.go` into `cli/cmdstatus/status_resolve.go`.
2. Move functions:
   - `loadStatusByScope(groupName string, all bool) ([]model.ScanRecord, error)`
   - `loadRecordsByGroup(groupName string) ([]model.ScanRecord, error)`
   - `buildStatusDBOpenError(op, group string, err error) *apperror.AppError`
   - `loadAllRecordsDB() ([]model.ScanRecord, error)`
   - `LoadRecordsJSONFallback() ([]model.ScanRecord, error)`
   - `loadAllRecordsDBOrEmpty() ([]model.ScanRecord, error)`
   - `newStatusNoDataError(op, code string) *apperror.AppError`
   - `loadStatusRecords(path string) ([]model.ScanRecord, error)`
   - `handleStatusDBError(err error) error`
3. Drastically reduce `status.go` from 267 lines to ~75 lines, leaving it strictly responsible for flag parsing, progress bar presentation, and table rendering.
4. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
`cli/cmdstatus/status.go` currently mixes multiple divergent domains:
- Command argument and flag processing (`--dirty`, `--all`, `--group`).
- Scope resolution across aliases, group queries, full database scans, and JSON fallback files.
- Terminal table and summary statistics rendering.
- Database error handling and precondition error formulation.

Extracting scope resolution into `status_resolve.go`:
1. Reduces `status.go` to ~75 lines, achieving the <= 100-line modular design target.
2. Isolates the database and filesystem fallback retrieval mechanics in a dedicated file.
3. Aligns with Subtask 01 by converting void functions into direct error returns, ensuring leaf error return sovereignty.

---

## 3. Concrete Implementation Details

### 3.1 New File `cli/cmdstatus/status_resolve.go`
Create `cli/cmdstatus/status_resolve.go` in package `cmdstatus`:

```go
package cmdstatus

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// loadStatusByScope returns records filtered by alias, group, all DB repos, or JSON fallback.
func loadStatusByScope(groupName string, all bool) ([]model.ScanRecord, error) {
	if HasAlias() {
		return []model.ScanRecord{{
			RepoName:     GetAliasSlug(),
			Slug:         GetAliasSlug(),
			AbsolutePath: GetAliasPath(),
		}}, nil
	}

	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}

	if all {
		return loadAllRecordsDB()
	}

	return LoadRecordsJSONFallback()
}

func loadRecordsByGroup(groupName string) ([]model.ScanRecord, error) {
	db, err := openDB()
	if err != nil {
		return nil, buildStatusDBOpenError("loadGroup", groupName, err)
	}

	defer db.Close()

	records, err := db.ShowGroup(groupName)
	if err != nil {
		return nil, handleStatusDBError(err)
	}

	return records, nil
}

func buildStatusDBOpenError(op, group string, err error) *apperror.AppError {
	ctx := map[string]any(nil)
	if group != "" {
		ctx = map[string]any{"group": group}
	}

	return apperror.WrapWithDetails(
		err,
		"cmd.status."+op+".openDB",
		"E1083",
		"failed to open database for status group load",
		"cmd.status",
		apperror.ErrorTypeExecution,
		apperror.SeverityFatal,
		ctx,
	)
}

func loadAllRecordsDB() ([]model.ScanRecord, error) {
	db, err := openDB()
	if err != nil {
		return nil, apperror.WrapWithDetails(
			err,
			"cmd.status.loadAll.openDB",
			"E1084",
			"failed to open database for status load",
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
	}

	defer db.Close()

	records, err := db.ListRepos()
	if err != nil {
		return nil, handleStatusDBError(err)
	}

	return records, nil
}

func LoadRecordsJSONFallback() ([]model.ScanRecord, error) {
	jsonPath := filepath.Join(constants.DefaultOutputDir, constants.DefaultJSONFile)
	if _, statErr := os.Stat(jsonPath); os.IsNotExist(statErr) {
		return loadAllRecordsDBOrEmpty()
	}

	records, err := loadStatusRecords(jsonPath)
	if err != nil {
		return nil, apperror.WrapWithDetails(
			err,
			"cmd.status.loadJSON",
			"E1085",
			"failed to read status records from JSON",
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": jsonPath},
		)
	}

	return records, nil
}

func loadAllRecordsDBOrEmpty() ([]model.ScanRecord, error) {
	db, err := openDB()
	if err != nil {
		return nil, newStatusNoDataError("openDB", "E1086")
	}

	defer db.Close()

	records, err := db.ListRepos()
	if err != nil {
		return nil, handleStatusDBError(err)
	}

	if len(records) == 0 {
		return nil, newStatusNoDataError("noRepos", "E1087")
	}

	return records, nil
}

func newStatusNoDataError(op, code string) *apperror.AppError {
	return apperror.NewWithDetails(
		"cmd.status."+op,
		code,
		constants.MsgStatusNoData,
		"cmd.status",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		nil,
	)
}

func loadStatusRecords(path string) ([]model.ScanRecord, error) {
	return model.LoadStatusRecords(path)
}

func handleStatusDBError(err error) error {
	if isLegacyDataError(err) {
		return apperror.WrapWithDetails(
			err,
			"cmd.status.legacyData",
			"E1088",
			constants.MsgLegacyProjectData,
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
	}

	return apperror.WrapWithDetails(
		err,
		"cmd.status.dbError",
		"E1089",
		"database operation failed during status lookup",
		"cmd.status",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
}
```

### 3.2 Update `cli/cmdstatus/status.go`
Remove the extracted resolution functions from `cli/cmdstatus/status.go`. In `ExecuteStatus`, invoke `records, err := loadStatusByScope(groupName, all)` and handle errors at the root dispatch level.

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** All functions return structured `error` or `*apperror.AppError` directly. Zero calls to `cliexit.HandleGeneralError` inside leaf helpers.
- **Rule 2 (Function Limit):** Every function <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Linear guard clauses with early returns throughout.
- **Rule 4 (Positive Booleans):** All conditionals use positive checks (`HasAlias()`, `len(groupName) > 0`, `all`).
- **Rule 6 (Vertical Whitespace):** Blank lines after `}` and before `return`.

---

## 5. Verification Criteria
- [x] File `cli/cmdstatus/status_resolve.go` created in package `cmdstatus`.
- [x] All 9 scope and resolution helpers moved to `status_resolve.go`.
- [x] Leaf error sovereignty established: direct `error` returns replace `cliexit` calls.
- [x] `status.go` reduced from 267 lines to ~75 lines.
- [x] Zero build/test executions performed.
