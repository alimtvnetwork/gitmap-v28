# Subtask 01: OS-Aware Store Indexes & Path Normalization

- **Subtask ID:** Task-01
- **Assigned Worker:** Worker 01
- **Owned Files:**
  - `cli/store/store.go`
  - `cli/store/storage_inventory.go`
  - `cli/vscodepm/path_filter.go`
  - `cli/workspacesync/path_guard.go`

## Instructions
1. In `cli/store/store.go`, replace static constants `constants.SQLCreateAbsPathIndex` and `constants.SQLCreateScanFolderPathIndex` with functions `getRepoAbsPathIndexQuery()` and `getScanFolderPathIndexQuery()` that check `runtime.GOOS == "windows"`.
   - On Windows: use `COLLATE NOCASE`.
   - On Unix: use standard binary collation.
2. In `cli/store/storage_inventory.go`, replace `filepath.ToSlash(strings.ToLower(clean))` with `fsutil.CanonicalPathKey(clean)`.
3. In `cli/vscodepm/path_filter.go` and `cli/workspacesync/path_guard.go`, ensure path comparisons respect OS case sensitivity via `fsutil.IsPathCaseInsensitive`.
