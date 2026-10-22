# Subtask 01: Status DB Helper Error Return Sovereignty & Propagation

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-01.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Eliminate all leaf exit calls (`cliexit.HandleGeneralError`) in `cli/cmdstatus/status.go` and `cli/cmdstatus/helpers_shim.go`.
2. Convert void and error-swallowing status helper routines (`loadRecordsByGroup`, `loadAllRecordsDB`, `handleStatusDBError`, `loadStatusByScope`) into direct `error` / `*apperror.AppError` returns.
3. Bubble all database access and parsing errors directly up to the top-level command handler `RunStatus(args []string) error` so that `finishCommandAudit` accurately logs failed executions into the database audit log.
4. Enforce strict coding guidelines across modified files:
   - Function limit: $\le 15$ lines per function ($\le 8$ lines preferred).
   - Conditionals: Zero nested `if` statements through early returns and guard clauses.
   - Positive booleans: Variables and parameters prefixed strictly with `is` or `has` (e.g., `hasRecords`, `isDirty`).
   - Vertical spacing: Mandatory blank line after closing brace `}` and blank line before every `return`.
   - File length ratchet: Keep all files strictly under the 500-line limit.

---

## 2. Detailed Technical Plan

### 2.1 Refactor DB Error Factory in `cli/cmdstatus/status.go`
- Locate `handleStatusDBError(err error)`:
  - Remove direct invocations of `cliexit.HandleGeneralError(appErr)`.
  - Update signature to `func handleStatusDBError(err error) *apperror.AppError`.
  - Invert legacy error branch into a guard clause:
    ```go
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
    ```

### 2.2 Refactor `loadRecordsByGroup` & `loadAllRecordsDB`
- Locate `loadRecordsByGroup(groupName string)`:
  - Change signature to `func loadRecordsByGroup(groupName string) ([]model.ScanRecord, error)`.
  - Remove `cliexit.HandleGeneralError(appErr)`.
  - If `openDB()` fails, return `nil, appErr`.
  - If `db.ShowGroup(groupName)` fails, return `nil, handleStatusDBError(err)`.
  - Return `records, nil` on success.
- Locate `loadAllRecordsDB()`:
  - Change signature to `func loadAllRecordsDB() ([]model.ScanRecord, error)`.
  - Remove `cliexit.HandleGeneralError(appErr)`.
  - If `openDB()` fails, return `nil, appErr`.
  - If `db.ListRepos()` fails, return `nil, handleStatusDBError(err)`.
  - Return `records, nil` on success.

### 2.3 Refactor `loadStatusByScope` & `RunStatus`
- Locate `loadStatusByScope(groupName string, all bool)`:
  - Update signature to `func loadStatusByScope(groupName string, all bool) ([]model.ScanRecord, error)`.
  - Check returned error when delegating to `loadRecordsByGroup` or `loadAllRecordsDB`:
    ```go
    if len(groupName) > 0 {
        return loadRecordsByGroup(groupName)
    }

    if all {
        return loadAllRecordsDB()
    }
    ```
- In `RunStatus(args []string) error`:
  - Capture records and error from `loadStatusByScope`:
    ```go
    records, err := loadStatusByScope(groupName, all)
    if err != nil {
        return err
    }
    ```
  - Bubble `err` directly to caller so the top-level dispatcher handles exit and audit logging.

---

## 3. Verification Criteria

- [ ] Zero invocations of `cliexit.HandleGeneralError`, `cliexit.HandleError`, or `os.Exit` exist within `cli/cmdstatus/status.go` or `cli/cmdstatus/helpers_shim.go`.
- [ ] `loadRecordsByGroup`, `loadAllRecordsDB`, and `loadStatusByScope` declare explicit `([]model.ScanRecord, error)` return signatures.
- [ ] `handleStatusDBError` returns `*apperror.AppError` instead of calling exit handlers.
- [ ] `RunStatus` captures and returns errors from `loadStatusByScope`.
- [ ] All functions in touched files strictly comply with the $\le 15$ lines ceiling.
- [ ] Zero nested `if` statements; all conditions flattened with guard clauses.
- [ ] Vertical blank lines strictly placed after `}` and before `return`.
- [ ] Strictly relative Git paths used throughout all references.
