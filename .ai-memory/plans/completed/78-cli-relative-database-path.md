# Plan 78: Anchor Pipeline Database to CLI Data Directory with Relative Display Path

## Overview
- **Goal:** Anchor GitMap's pipeline database (and global split databases when no `repoRoot` is provided) directly to the CLI's installation directory inside its `data/` folder (`BinaryDataDir()`), and format displayed database paths as clean relative paths (`data/pipeline/<repo-slug>/sql.db`).
- **Status:** Completed
- **Target Release:** v6.307.0

## Root Cause Analysis
- `resolveBaseDataDir(repoRoot string)` in `cli/store/split_db_path.go` previously called `resolveTargetRoot(repoRoot)`. When `repoRoot == ""`, `resolveTargetRoot` searched up the directory tree from CWD looking for `.git` or `.gitmap`.
- If the command was executed inside any git repository where `.gitmap` existed, `resolveBaseDataDir` returned `filepath.Join(root, ".gitmap", "data")`.
- This caused `gitmap pe` and pipeline cache sync to place `sql.db` inside `.gitmap/data/pipeline/<repo-slug>/sql.db` within the user's workspace repo rather than co-locating it with the CLI binary in its own `data/` directory.
- Furthermore, `FormatRelativeDbPath` in `cli/cmdpipeline/pipeline_sync_cache.go` only stripped `.gitmap/` prefixes, returning full absolute paths for databases located outside the workspace.

## Key Changes
1. **`cli/store/split_db_path.go`**:
   - Updated `resolveBaseDataDir(repoRoot string)` to immediately return `BinaryDataDir()` when `repoRoot == ""`.
   - Preserved explicit `repoRoot` paths (e.g. In unit tests) while ensuring runtime CLI database calls are always anchored next to `gitmap.exe`.
2. **`cli/cmdpipeline/pipeline_sync_cache.go`**:
   - Added `github.com/alimtvnetwork/gitmap-v28/cli/store` import.
   - Updated `FormatRelativeDbPath` with `resolveBinaryOrRelPath(slashPath, target)` helper to format paths located in `BinaryDataDir()` or containing `/data/pipeline/` as clean relative paths: `data/pipeline/<repo-slug>/sql.db`.
3. **`cli/store/split_db_path_test.go`**:
   - Added `TestResolveSplitDbDirEmptyRoot` and `TestResolveSplitDbPathEmptyRoot` verifying that empty `repoRoot` correctly resolves to `BinaryDataDir()`.

## Verification
- Verified code passes `go vet ./...` with 0 warnings/errors.
- Verified coding guidelines (functions <= 15 lines, positive booleans, no nested conditionals).
