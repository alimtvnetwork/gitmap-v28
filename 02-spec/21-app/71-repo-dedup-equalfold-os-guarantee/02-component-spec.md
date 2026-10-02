# Component Specification: Redundant Repos, EqualFold, OS Path Sensitivity & SQLite Uniqueness

## Component Details

### 1. `cli/cmdignore/ignore_cmd.go` & `cli/cmdinstall/installantigravity_deploy_windows.go`
- Eliminates manual lowercase allocation; replaces equality check with `strings.EqualFold`.

### 2. `cli/cmdvscode/find_duplicates_vscode.go` & `cli/cmdchromeprofile/find_duplicates_chrome.go`
- Replaces raw `strings.ToLower(filepath.Clean(...))` with `fsutil.CanonicalPathKey(...)` to ensure correct OS path behavior.

### 3. `cli/repodb/repo_db.go` & `cli/repodb/repo_file.go`
- `IdxRepoFile_RelativePath` index uses `COLLATE NOCASE` on Windows and standard binary collation on Unix.
- `EnsureFileUniqueInDB` allows checking whether a file path is already indexed before insertion.
