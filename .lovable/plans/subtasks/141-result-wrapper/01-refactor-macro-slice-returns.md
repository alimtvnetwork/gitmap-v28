# Subtask 01: Refactor Macro Subsystem Slice Returns to ResultSlice[Macro]

Parent Plan: [141-result-wrapper-and-slice-returns.md](../../pending/141-result-wrapper-and-slice-returns.md)

## Goals
1. Refactor `cli/macro/import_sqlite.go`:
   - `ParseImportSQLite(dbPath string) result.ResultSlice[Macro]`
   - `queryAllMacrosWithSteps(db *sql.DB) result.ResultSlice[Macro]`
   - `scanMacroRows(rows *sql.Rows, stepMap map[string][]MacroStep) result.ResultSlice[Macro]`
2. Refactor `cli/macro/import.go`:
   - `ParseImportJSON(data []byte) result.ResultSlice[Macro]`
   - `ParseImportYAML(data []byte) result.ResultSlice[Macro]`
   - `ParseImportZIP(filePath string) result.ResultSlice[Macro]`
   - `extractMacrosFromZip(zr *zip.Reader) result.ResultSlice[Macro]`
   - `ParseImportFile(filePath string, format string) result.ResultSlice[Macro]`
3. Refactor `cli/macro/storage.go`:
   - `ListMacros() result.ResultSlice[Macro]`
4. Modernize callers in:
   - `cli/cmdmacro/macro_import.go`
   - `cli/cmdmacro/macro_cmd.go`
   - `cli/cmdmacro/macro_add_history_seed.go`
   - `cli/cmdmacro/macro_export.go`
   - `cli/macro/macro_test.go`
   - `cli/macro/export_import_test.go`

## Acceptance Criteria
- [x] All 9 macro slice return functions use `result.ResultSlice[Macro]`.
- [x] Callers modernized with `res.IsSuccess()`, `res.IsFailure()`, `res.Data`, `res.AppError()`.
- [x] Subtask completed and logged.
