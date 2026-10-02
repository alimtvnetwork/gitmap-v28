# Subtask 01: Refactor String Comparisons to EqualFold and CanonicalPathKey in CLI Commands

- **Task Code:** Task-01
- **Owner:** Worker 01
- **Status:** PENDING
- **Owned Files:**
  - `cli/cmdignore/ignore_cmd.go`
  - `cli/cmdinstall/installantigravity_deploy_windows.go`
  - `cli/cmdvscode/find_duplicates_vscode.go`
  - `cli/cmdchromeprofile/find_duplicates_chrome.go`

## Instructions:
1. In `cli/cmdignore/ignore_cmd.go`:
   - Replace `if strings.ToLower(args[0]) == "interval" && len(args) > 1 {` with `if strings.EqualFold(args[0], "interval") && len(args) > 1 {`.
2. In `cli/cmdinstall/installantigravity_deploy_windows.go`:
   - In `hasDirInPathString(pathStr, dir string) bool`:
     - Clean `dir` once outside the loop (`cleanDir := filepath.Clean(dir)`).
     - Inside loop: replace `if strings.ToLower(filepath.Clean(strings.TrimSpace(part))) == targetLow` with `if strings.EqualFold(filepath.Clean(strings.TrimSpace(part)), cleanDir)`.
3. In `cli/cmdvscode/find_duplicates_vscode.go`:
   - In `groupVSCodeDuplicates(entries []vscodepm.Entry)`:
     - Replace `norm := strings.ToLower(filepath.Clean(e.RootPath))` with `norm := fsutil.CanonicalPathKey(e.RootPath)`. Import `github.com/alimtvnetwork/gitmap-v28/cli/fsutil` if needed.
4. In `cli/cmdchromeprofile/find_duplicates_chrome.go`:
   - In `groupChromeDuplicates(profiles []chromeProfileDetail)`:
     - Replace `key := strings.ToLower(filepath.Clean(p.SourcePath))` with `key := fsutil.CanonicalPathKey(p.SourcePath)`. Import `github.com/alimtvnetwork/gitmap-v28/cli/fsutil` if needed.
5. Log actions in SQLite task DB before touching files, mark completed with evidence.
