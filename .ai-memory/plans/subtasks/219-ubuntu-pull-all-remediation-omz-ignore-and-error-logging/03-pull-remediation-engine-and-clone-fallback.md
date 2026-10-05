# Subtask Plan 03: Enhanced Pull Remediation Engine, Clone Fallback & DB-Only Removal

- **Subtask Slug:** `03-pull-remediation-engine-and-clone-fallback`
- **Parent Task:** `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging`
- **Target Files:**
  - `cli/constants/constants_pull.go`
  - `cli/cmdpull/pull_remediation.go`
  - `cli/cmdpull/pull_worker.go`
  - `cli/cmdpull/pull_not_found.go`
  - `cli/store/repo.go`

---

## 1. Overview & Objectives

When a repository directory tracked in `gitmap.db` is missing on disk (`cloner.IsMissingRepo` returns true):
1. `gitmap pull` and `gitmap pa` mark the repository as skipped or failed.
2. The current pull remediation engine (`cli/cmdpull/pull_remediation.go`) blindly invokes:
   `git -C <missing_path> pull --progress --no-rebase --no-edit --autostash`
   which crashes immediately with `fatal: cannot change to '...': No such file or directory`.
3. Users are left with no automated mechanism to:
   - Re-clone the repository from its known remote URL (`rec.HTTPSUrl` or `rec.SSHUrl`).
   - Cleanly prune the dead/stale repository record from `gitmap.db` without affecting physical filesystem trees.

This subtask implements:
1. Intelligent branching in `remediateSingleRepo`: detecting missing repositories and preventing invalid `git -C <missing_dir> pull` invocations.
2. Automatic clone fallback using known remote URLs from `gitmap.db`.
3. Safe database-only removal (`DELETE FROM Repo WHERE RepoId = ?`) for permanently deleted or inaccessible repositories.
4. Interactive remediation menu additions supporting `[c]lone`, `[d]elete record`, and `[s]kip`.

---

## 2. Step-by-Step Implementation Instructions

### Step 1: Enhance `PullFailureSummary` with Remote & ID Metadata
1. In `cli/cmdpull/pull_remediation.go`:
   ```go
   // PullFailureSummary holds concise failure information for a repository.
   type PullFailureSummary struct {
       RepoID    int64
       RepoName  string
       RepoPath  string
       RemoteURL string
       ErrorText string
       ErrorType string
   }
   ```
2. In `ExtractPullFailures(states []*PullRepoState)`:
   - Extract and populate `RemoteURL` from `rec.HTTPSUrl` or `rec.SSHUrl`.
3. In `loadRecentPullFailuresFromDB()`:
   - Query `Repo` table in `gitmap.db` using `RepoPath` to populate `RepoID` and `RemoteURL` (`COALESCE(HTTPSUrl, SSHUrl, DiscoveredURL)`).

### Step 2: Implement Missing Repository Remediation in `cli/cmdpull/pull_remediation.go`
1. Update `remediateSingleRepo(f PullFailureSummary) bool`:
   ```go
   func remediateSingleRepo(f PullFailureSummary) bool {
       if f.RepoPath == "" {
           return false
       }

       // Gate 1: Check if repository directory is missing
       if cloner.IsMissingRepo(f.RepoPath) {
           return remediateMissingRepo(f)
       }

       // Gate 2: Normal git pull / auto-merge for existing dirty or conflict repos
       fmt.Printf("  %s→%s Remediating %s%s%s...\n",
           constants.ColorCyan, constants.ColorReset,
           constants.ColorBold, f.RepoName, constants.ColorReset)
       out, err := execPullAutoMerge(f.RepoPath)
       if err == nil {
           fmt.Printf("  %s✓%s [%s] Successfully updated.\n",
               constants.ColorGreen, constants.ColorReset, f.RepoName)
           return true
       }

       return handleRemediationFailure(f.RepoPath, f.RepoName, string(out))
   }
   ```
2. Implement `remediateMissingRepo(f PullFailureSummary) bool`:
   ```go
   func remediateMissingRepo(f PullFailureSummary) bool {
       fmt.Printf("  %s⚠%s [%s] Directory missing on disk: %s\n",
           constants.ColorYellow, constants.ColorReset, f.RepoName, f.RepoPath)

       if f.RemoteURL == "" {
           fmt.Printf("  %s✗%s [%s] Cannot clone: no remote URL recorded in database.\n",
               constants.ColorRed, constants.ColorReset, f.RepoName)
           return false
       }

       fmt.Printf("  %s→%s Attempting clone fallback from %s...\n",
           constants.ColorCyan, constants.ColorReset, f.RemoteURL)

       if err := execCloneRepo(f.RemoteURL, f.RepoPath); err != nil {
           fmt.Printf("  %s✗%s [%s] Clone fallback failed: %v\n",
               constants.ColorRed, constants.ColorReset, f.RepoName, err)
           return false
       }

       fmt.Printf("  %s✓%s [%s] Successfully restored repository via clone fallback.\n",
           constants.ColorGreen, constants.ColorReset, f.RepoName)
       return true
   }

   func execCloneRepo(remoteURL, destination string) error {
       if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
           return err
       }
       cmd := exec.Command("git", "clone", "--progress", remoteURL, destination)
       cmd.Env = gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)
       return cmd.Run()
   }
   ```

### Step 3: Implement Database-Only Removal for Dead Repos
1. In `cli/store/repo.go` (or `cli/cmdpull/pull_remediation.go`):
   - Implement `DeleteRepoByPath(absPath string) error`:
     ```go
     func (s *Store) DeleteRepoByPath(absPath string) error {
         _, err := s.db.Exec("DELETE FROM Repo WHERE AbsolutePath = ?", absPath)
         return err
     }
     ```
2. Wire removal action into interactive remediation menu:
   - When a missing repo is prompted in `promptAndRemediateSingle`:
     ```text
     [!] Repository 'old-project' missing at /home/ubuntu/old-project
         Remote: https://github.com/org/old-project.git
         Options:
           [1/c] Clone from remote URL
           [2/d] Delete record from database (clean up stale tracking)
           [3/s] Skip
     ```
   - If user chooses `d` / `delete`:
     * Execute `DeleteRepoByPath(f.RepoPath)`
     * Report: `✓ [old-project] Stale record removed from database.`

### Step 4: Batch Remediation Automation
1. In `RunBatchRemediateAll(failures []PullFailureSummary)`:
   - For missing repositories:
     * If remote URL is present: automatically execute `execCloneRepo(f.RemoteURL, f.RepoPath)`.
     * If remote URL is absent: skip with clear diagnostic rather than invoking `git -C pull`.

---

## 3. CLI UX & Output Verification

### Interactive Remediation Flow:
```bash
gitmap pull --remediate
```
Output:
```text
  Found 2 failed repository(ies) from recent runs.

  [1/2] Repository: backend-service
    Path:  /home/ubuntu/work/backend-service
    Error: missing repository directory
    Remote: git@github.com:myorg/backend-service.git

  ? Remediate this repository? [c=clone / d=delete / s=skip / q=quit]: c
  → Attempting clone fallback from git@github.com:myorg/backend-service.git...
  ✓ [backend-service] Successfully restored repository via clone fallback.

  [2/2] Repository: scratch-experiment
    Path:  /home/ubuntu/work/scratch-experiment
    Error: missing repository directory
    Remote: (none)

  ? Remediate this repository? [c=clone / d=delete / s=skip / q=quit]: d
  ✓ [scratch-experiment] Stale record removed from database.

  ✓ Step-by-step remediation complete: 2/2 resolved.
```

---

## 4. Verification & Quality Gates

1. **Linter & Guideline Verification:**
   - `python linter-scripts/check-relative-paths.py`
   - `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdpull cli/store --check-only`
2. **Behavioral Acceptance:**
   - Nonexistent directories never cause subprocess failures via `git -C <missing> pull`.
   - Repositories with known remotes are restored via clone fallback upon user or batch request.
   - Stale records can be deleted from SQLite without leaving dangling metadata.
