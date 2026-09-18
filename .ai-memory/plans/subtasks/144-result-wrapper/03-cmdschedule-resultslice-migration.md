# Subtask 03: CmdSchedule ResultSlice Migration & Structured AppErrors

Parent Plan: [144-result-wrapper-null-safety-and-single-return-audit.md](../../completed/144-result-wrapper-null-safety-and-single-return-audit.md)

## Goals
1. In `cli/cmdschedule/schedule_export.go`:
   - Refactor `collectExportBundles(db *store.DB, opts scheduleExportOpts) result.ResultSlice[scheduleExportBundle]`.
   - Replace raw error returns with structured `*apperror.AppError`.
   - Update call sites in `schedule_export.go`.
2. In `cli/cmdschedule/schedule_import.go`:
   - Refactor `parseImportFileBundles(filePath string) result.ResultSlice[scheduleExportBundle]`.
   - Refactor `parseImportJSON(filePath string) result.ResultSlice[scheduleExportBundle]`.
   - Refactor `parseImportYAML(filePath string) result.ResultSlice[scheduleExportBundle]`.
   - Refactor `parseImportSQLite(filePath string) result.ResultSlice[scheduleExportBundle]`.
   - Refactor `queryImportTasksFromDB(conn *sql.DB) result.ResultSlice[store.SchedulerTask]`.
   - Refactor `parseImportZIP(filePath string) result.ResultSlice[scheduleExportBundle]`.
   - Replace raw stdlib errors with structured `*apperror.AppError`.
   - Update call sites in `schedule_import.go`.
3. In `cli/cmdschedule/schedule_export_import_test.go`:
   - Update call site at line 51 to use `bundlesRes.IsCountOtherThan(1)` instead of `err != nil || len(bundles) != 1`.
   - Verify all test assertions adapt to `ResultSlice[T]` envelopes.

## Acceptance Criteria
- [x] All 7 slice functions in `cli/cmdschedule/` return `result.ResultSlice[T]`.
- [x] Zero raw `error` returns in `schedule_export.go` and `schedule_import.go`.
- [x] `go vet ./...` in `cli/` passes with 0 errors.
