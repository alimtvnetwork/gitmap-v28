# Component & Remediation Specification: Forward-Slash Paths & Interactive Remediation Engine

- **Feature Slug:** `243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt`
- **Specification Version:** 1.0.0
- **Status:** APPROVED

---

## 1. Forward-Slash Path Normalization

### A. `cli/cmdpull/pull_remediation_hint.go`
In `formatRepoGitCmd(repoDir, gitArgs string) string`:
```go
func formatRepoGitCmd(repoDir, gitArgs string) string {
	if repoDir == "" {
		return "git " + gitArgs
	}
	cleanDir := filepath.ToSlash(filepath.Clean(repoDir))
	return fmt.Sprintf("git -C \"%s\" %s", cleanDir, gitArgs)
}
```
**Why:**
- `fmt.Sprintf("git -C %q %s", repoDir, gitArgs)` on Windows uses Go's quoted string syntax which escapes backslashes: `"D:\\work\\gitmap"`.
- By converting `repoDir` with `filepath.ToSlash(...)` and formatting with `git -C "%s"`, the command becomes `git -C "D:/work/gitmap" stash`.
- This works identically on Windows, Linux, and macOS without double-backslash escapes.

### B. `cli/cmdpull/pull_efficient_render.go`
In `renderItemizedDirtyFiles(w io.Writer, diag gitutil.DirtyDiagnosis)`:
- Normalize each modified and untracked file via `filepath.ToSlash(f)`:
  ```go
  fmt.Fprintf(w, "        %smodified: %s%s\n", constants.ColorDim, filepath.ToSlash(f), constants.ColorReset)
  fmt.Fprintf(w, "        %suntracked: %s%s\n", constants.ColorDim, filepath.ToSlash(f), constants.ColorReset)
  ```

---

## 2. Interactive Remediation Prompt Engine

### A. Dispatch in `cli/cmdpull/pull.go`
In `handlePullRemediationForRecords(records []model.ScanRecord, opts pullOptions)`:
- Remove the suppression check `isConcisePullOutput(opts.all, opts.showStatus)`.
- When dirty or failed repositories exist:
  ```go
  func handlePullRemediationForRecords(records []model.ScanRecord, opts pullOptions) {
      if opts.isJSON || opts.noFix {
          return
      }
      var remItems []RemediationItem
      for _, rec := range records {
          diag := gitutil.InspectDirtyState(rec.AbsolutePath)
          if diag.IsDirty {
              remItems = append(remItems, buildRemediationItem(rec, diag))
          }
      }
      if len(remItems) == 0 {
          return
      }
      dispatchRemediationSummary(remItems, opts)
  }
  ```

### B. Interactive Menu in `cli/cmd/remediation_box.go`
In `promptForRemediation(items []RemediationItem)`:
```text
  Remediate dirty / failed repository(ies) now?
    [a/1] Fix all
    [s/2] Fix single / select repository
    [k/q] Skip / Exit
  Choice [a/s/k]: 
```
Choices handled:
- `a`, `1`, `all`, `y`, `yes`: Runs auto-remediation across all dirty repos (stash or wip commit, with pull-before-changes).
- `s`, `2`, `single`, `step`: Steps through each repository interactively or prompts for selection index.
- `k`, `q`, `skip`, `n`, `no`, `exit`: Exits cleanly.

---

## 3. Pull Before Changes Invariant

In `remediateSingleRepo(item RemediationItem)` / `runInteractiveRemediation(items []RemediationItem)`:
1. Ensure the remote branch is fetched/pulled before applying stashes or commits.
2. If the repository is currently behind, `git pull --rebase` or `git pull` is invoked as part of the remediation recipe.
3. This ensures that the workspace state is clean and current with upstream.
