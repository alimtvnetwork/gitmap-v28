# Subtask 02: CLI Purge Command Decomposition & Zero-Swallow Error Handling

## Objective
Refactor `gitmap/cmd/purge.go` and extract `gitmap/cmd/purge_engine.go`:
- Keep `runPurge(args []string) error` in `gitmap/cmd/purge.go` (under 70 lines) for CLI argument parsing (`isRestore`, `isAutoConfirm`, `pattern`).
- Extract the core purge engine into `gitmap/cmd/purge_engine.go` (under 120 lines).
- Decompose the 134-line `doPurge` into focused functions <= 15 lines:
  - `verifyWorkTreeClean() error`
  - `fetchCurrentBranch() (string, error)`
  - `queryMatchingFiles(pattern string) ([]string, error)`
  - `createBackupBranch(branch string) error`
  - `backupFilesToTemp(tempDir string, files []string) ([]string, error)`
  - `executeFilterRepo(pattern string) error`
  - `appendGitignore(pattern string) error`
  - `savePurgeState(db *store.DB, log *store.PurgeHistoryLog) error`
- Fix all swallowed errors: check `os.MkdirAll`, check `copyPurgeFile`, check `sendToRecycleBin`, check `f.Close()`, check `git remote add`, and check `json.Marshal`.
- Normalize pattern path separators using `filepath.ToSlash(pattern)`.
- Use boolean variable `isAutoConfirm` instead of `autoConfirm`.

## Files Affected
- `gitmap/cmd/purge.go`
- `gitmap/cmd/purge_engine.go`
