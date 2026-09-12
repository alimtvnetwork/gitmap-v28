# RCA: Nested If Violations in `macro_export.go`, `macro_import.go`, and `macro/export.go`

## 1. Why it happened
Recent macro multi-format export/import features introduced nested conditional blocks (`depth 2`) inside error handling and argument disambiguation paths. The repository's AST linter (`linter-scripts/check-nested-ifs.py` and `linter-scripts/check-enum-and-boolean.py`) strictly forbids nested `if` blocks (max depth 1).

## 2. How it happened
- In `cmd/macro_export.go`, single-macro loading check nested an error check: `if !opts.IsAll && opts.TargetName != "" { if err != nil { ... } }`.
- In `cmd/macro_export.go`, default output resolution nested a target name check: `if outPath == "" { if opts.TargetName != "" && !opts.IsAll { ... } }`.
- In `cmd/macro_import.go`, extension disambiguation nested a known-extension check: `if !fileExists(...) && !fileExists(...) { if hasKnownImportExtension(...) { ... } }`.
- In `macro/export.go`, directory creation nested a mkdir check: `if dir := filepath.Dir(...); dir != "" && dir != "." { if err := os.MkdirAll(...) { ... } }`.

## 3. Root Cause
- Files:
  - `cli/cmd/macro_export.go:194, 226`
  - `cli/cmd/macro_import.go:89`
  - `cli/macro/export.go:160`
- Nested `if` blocks exceeding depth limit of 1.

## 4. Code Fix
- Extracted `collectSingleMacro(name string)` and `resolveSQLiteOutPath(opts macroExportOpts)` in `cmd/macro_export.go`.
- Extracted `shouldSwapImportPaths(opts *macroImportOpts) bool` in `cmd/macro_import.go`.
- Converted `ensureExportDirExists(filePath string)` to early guard return in `macro/export.go`.
