# Plan 196: Storage List UI Cleanup & Relative Paths

## 1. Overview
Clean up and enhance the `gitmap storage ls` table presentation. Previously, the table dumped full absolute paths (e.g. `C:\Users\Administrator\AppData\Local\gitmap-cli\data\installation.db`) directly into the `PATH` column, resulting in terminal line wrapping, column collisions, and cramped presentation.

This task enhances the UI by:
1. Displaying the root database folder above the inventory table.
2. Calculating clean relative paths (`installation.db`, `gitmap.db`, etc.) relative to the root storage directory for all table entries.
3. Leveraging the modular `termtable.PrintTable` framework with proper column alignments (numeric right-aligned, text left-aligned).
4. Decomposing database listing from `cli/cmd/storage_display.go` into `cli/cmd/storage_ls.go` to strictly adhere to the 100-line canonical file limit.

## 2. Changes Made
- **[NEW] `cli/cmd/storage_ls.go`** (79 lines):
  - Implements `runStorageListDatabases()`, `printStorageInventoryHeader()`, `buildStorageTableConfig()`, `storageTableColumns()`, `buildStorageTableRow()`, and `resolveRelativeDBPath()`.
  - Safely falls back to full path if `filepath.Rel` encounters errors or escapes the parent directory (`..`).
- **[MODIFY] `cli/cmd/storage_display.go`** (reduced from 95 to 69 lines):
  - Extracted database listing, table headers, and row formatting functions to `storage_ls.go`.
- **[MODIFY] `cli/cmd/storage_cmd_test.go`** (75 lines):
  - Added `TestResolveRelativeDBPath` verifying correct relative path resolution and fallback on external paths.

## 3. Verification
- `check-file-sizes.py`: All modified files strictly <= 100 lines.
- `check-nested-ifs.py`: 0 nested if violations.
- `check-enum-and-boolean.py`: 0 violations.
- `06-cicd-local-runner.py --filter "Compile"`: Go Compile Gate passed 100%.
- `06-cicd-local-runner.py --no-tests --filter "Linter"`: All 5 linters passed 100%.
- Interactive CLI validation: `bin/gitmap.exe storage ls` renders cleanly with root directory display and relative paths.
