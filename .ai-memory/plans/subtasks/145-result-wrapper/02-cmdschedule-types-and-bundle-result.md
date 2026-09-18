# Subtask 02: Cmdschedule Types Centralization & Reusable Envelopes

Parent Plan: [145-result-wrapper-and-types-go-centralization-audit.md](../../completed/145-result-wrapper-and-types-go-centralization-audit.md)

## Goals
1. Create `cli/cmdschedule/types.go`:
   - Declare exported domain struct `ScheduleExportBundle`:
     ```go
     type ScheduleExportBundle struct {
         Task store.SchedulerTask       `json:"task" yaml:"task"`
         Runs []store.ScheduleRunRecord `json:"runs" yaml:"runs"`
     }
     ```
   - Declare exported CLI options struct `ScheduleExportOpts`:
     ```go
     type ScheduleExportOpts struct {
         TargetName string
         FilePath   string
         Format     string
         ExceptList []string
         IsAll      bool
     }
     ```
   - Declare single reusable Result aliases:
     `type ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]`
     `type ScheduleExportBundleSingleResult = result.Result[ScheduleExportBundle]`
     `type SchedulerTaskSliceResult = result.ResultSlice[store.SchedulerTask]`
2. Refactor `cli/cmdschedule/schedule_export.go`:
   - Remove inline `scheduleExportBundle` and `scheduleExportOpts`.
   - Update `collectExportBundles` return type to `ScheduleExportBundleResult`.
   - Update `runScheduleExport` and `parseScheduleExportOpts`.
3. Refactor `cli/cmdschedule/schedule_import.go`:
   - Update `parseImportFileBundles` return type to `ScheduleExportBundleResult`.
   - Update `parseImportJSON` return type to `ScheduleExportBundleResult`.
   - Update `parseImportYAML` return type to `ScheduleExportBundleResult`.
   - Update `parseImportSQLite` return type to `ScheduleExportBundleResult`.
   - Update `queryImportTasksFromDB` return type to `SchedulerTaskSliceResult`.
   - Update `parseImportZIP` return type to `ScheduleExportBundleResult`.
   - Update `runScheduleImport`.
4. Update `cli/cmdschedule/schedule_export_import_test.go` to use `ScheduleExportBundle`.
5. Run `go test -v ./cmdschedule/...` in `cli/` to verify 100% pass.

## Acceptance Criteria
- [x] `cli/cmdschedule/types.go` created with exported domain models and Result aliases.
- [x] 7 functions updated to return single reusable types from `types.go`.
- [x] `schedule_export.go` and `schedule_import.go` cleaned of inline structs.
- [x] Unit tests pass with 0 errors.
