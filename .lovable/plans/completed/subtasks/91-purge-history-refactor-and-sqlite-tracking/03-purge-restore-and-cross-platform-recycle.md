# Subtask 03: Purge Restore Decomposition & Cross-Platform Recycle Isolation

## Objective
Extract restore logic and cross-platform recycle bin support:
- Extract `gitmap/cmd/purge_restore.go`:
  - Decompose `doRestore` (34 lines) into functions <= 15 lines:
    - `fetchActivePurgeLog(db *store.DB, repoPath string) (*store.PurgeHistoryLog, error)`
    - `resetToBranch(branch string) error`
    - `restoreFilesFromTemp(tempDir string) error`
    - `markPurgeRestored(db *store.DB, id int64) error`
  - In `restoreFilesFromTemp`, check errors from `filepath.WalkDir`, propagate walk errors, check `os.MkdirAll`, and check `copyPurgeFile`.
  - In `markPurgeRestored`, pass `log.PurgeHistoryLogId` or `log.Id` and handle the returned error (eliminate swallow).
- Create `gitmap/cmd/purge_recycle_windows.go`:
  - Tagged with `//go:build windows`.
  - Houses `sendToRecycleBin(path string) error` using `shell32.dll` `SHFileOperationW`.
- Create `gitmap/cmd/purge_recycle_other.go`:
  - Tagged with `//go:build !windows`.
  - Houses portable `sendToRecycleBin(path string) error` using `os.RemoveAll` or standard removal fallback so builds and execution on Linux/macOS never fail.

## Files Affected
- `gitmap/cmd/purge_restore.go`
- `gitmap/cmd/purge_recycle_windows.go`
- `gitmap/cmd/purge_recycle_other.go`
