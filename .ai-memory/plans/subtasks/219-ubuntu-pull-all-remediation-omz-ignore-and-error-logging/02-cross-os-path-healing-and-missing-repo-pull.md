# Subtask Plan 02: Cross-OS Path Normalization & Pre-Pull Auto-Healing

- **Subtask Slug:** `02-cross-os-path-healing-and-missing-repo-pull`
- **Parent Task:** `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging`
- **Target Files:**
  - `cli/constants/constants_pull.go`
  - `cli/store/repo_sanitize.go`
  - `cli/cmdpull/pull_dedup.go`
  - `cli/cmdpull/pull_worker.go`
  - `cli/cmdpull/helpers.go`

---

## 1. Overview & Objectives

In heterogenous environments where `gitmap.db` is transferred from Windows to Linux/Ubuntu nodes (or vice versa), repository records in the SQLite `Repo` table retain Windows backslash path separators (`D:\work\repo` or `C:\Users\admin\dev\tool`). 

On Linux/Ubuntu:
1. Backslash `\` is not a path separator; system calls treat `D:\work\repo` as a literal filename relative to the current directory, which fails immediately.
2. Direct execution of `git -C <path> pull` crashes with `cannot change to directory: No such file or directory`.
3. In-flight checks like `cloner.IsMissingRepo(rec.AbsolutePath)` report the repository as missing even though it exists at `/work/repo` or under normalized forward slashes.

This subtask implements:
1. Database-level cross-OS path normalization and auto-healing in `cli/store/repo_sanitize.go:HealCrossOSPaths`.
2. Automatic pre-pull invocation of `HealCrossOSPaths` inside `cli/cmdpull/pull_dedup.go`.
3. Real-time in-memory path repair fallback inside `cli/cmdpull/pull_worker.go` to salvage un-normalized records during active pull streaming.

---

## 2. Step-by-Step Implementation Instructions

### Step 1: Add Telemetry & Status Constants in `cli/constants/constants_pull.go`
1. Add the following user-facing messages:
   ```go
   const (
       MsgPullAutoHealedPaths = "  %s✓ SQLite: Auto-healed %d cross-OS repository path(s) to Unix slashes before pull.%s\n\n"
       MsgPullInFlightHealed  = "  %sℹ In-flight path repair: converted '%s' to '%s'%s\n"
   )
   ```

### Step 2: Implement `HealCrossOSPaths` in `cli/store/repo_sanitize.go`
1. Define the cross-OS path healing function:
   ```go
   // HealCrossOSPaths detects repositories with Windows backslashes on Unix systems
   // (or malformed separators) and auto-heals them in gitmap.db.
   func HealCrossOSPaths(db *sql.DB, isQuiet bool) (int, error) {
       if db == nil {
           return 0, nil
       }

       // Only execute backslash normalization on Unix systems
       if runtime.GOOS == "windows" {
           return 0, nil
       }

       rows, err := db.Query("SELECT RepoId, AbsolutePath FROM Repo WHERE AbsolutePath LIKE '%\\%'")
       if err != nil {
           return 0, err
       }
       defer rows.Close()

       type healTarget struct {
           repoID  int64
           oldPath string
           newPath string
       }

       var targets []healTarget
       for rows.Next() {
           var repoID int64
           var absPath string
           if err := rows.Scan(&repoID, &absPath); err == nil {
               normalized := strings.ReplaceAll(absPath, "\\", "/")
               // If Windows drive prefix exists (e.g., "D:/work/repo"), attempt WSL or project root resolution
               resolved := resolveUnixPathFromWindowsDrive(normalized)
               targets = append(targets, healTarget{
                   repoID:  repoID,
                   oldPath: absPath,
                   newPath: resolved,
               })
           }
       }

       healedCount := 0
       for _, t := range targets {
           if _, err := db.Exec("UPDATE Repo SET AbsolutePath = ? WHERE RepoId = ?", t.newPath, t.repoID); err == nil {
               healedCount++
           }
       }

       return healedCount, nil
   }

   func resolveUnixPathFromWindowsDrive(path string) string {
       // Match pattern like "C:/..." or "D:/..."
       if len(path) >= 3 && path[1] == ':' && path[2] == '/' {
           drive := strings.ToLower(string(path[0]))
           wslPath := "/mnt/" + drive + path[2:]
           if _, err := os.Stat(wslPath); err == nil {
               return wslPath
           }
           // Fallback: check if path without drive exists (e.g. "/work/repo")
           tailPath := path[2:]
           if _, err := os.Stat(tailPath); err == nil {
               return tailPath
           }
           return path[2:]
       }
       return path
   }
   ```

### Step 3: Wire Auto-Healing into Pre-Pull Sequence
1. In `cli/cmdpull/pull_dedup.go`:
   - Implement `AutoHealCrossOSPaths(db *store.DB, isQuiet bool) (int, error)`:
     ```go
     func AutoHealCrossOSPaths(db *store.DB, isQuiet bool) (int, error) {
         if db == nil {
             return 0, nil
         }
         healed, err := store.HealCrossOSPaths(db.UnderlyingDB(), isQuiet)
         if err != nil {
             return 0, err
         }
         if !isQuiet && healed > 0 {
             fmt.Fprintf(os.Stderr, constants.MsgPullAutoHealedPaths,
                 constants.ColorGreen, healed, constants.ColorReset)
         }
         return healed, nil
     }
     ```
2. In `cli/cmdpull/pull.go`:
   - In `runPullAll` and `runPull`, invoke `AutoHealCrossOSPaths(db, opts.quiet)` immediately before `OptimizeRedundantRepos(db, opts.quiet)`.

### Step 4: Real-Time In-Flight Path Repair in `cli/cmdpull/pull_worker.go`
1. In `ExecuteTrackedPullWithWorker`:
   - Before bailing out to `handleMissingRepoPull`:
     ```go
     if cloner.IsMissingRepo(rec.AbsolutePath) {
         // Attempt real-time in-flight path repair
         if runtime.GOOS != "windows" && strings.Contains(rec.AbsolutePath, "\\") {
             repaired := strings.ReplaceAll(rec.AbsolutePath, "\\", "/")
             if !cloner.IsMissingRepo(repaired) {
                 rec.AbsolutePath = repaired
                 return runTrackedPullLifecycle(rec, bar, workerID, start)
             }
         }
         return handleMissingRepoPull(rec, bar, start)
     }
     ```
2. In `sanitizeScanRecordIdentity(rec *model.ScanRecord)`:
   - Ensure `rec.AbsolutePath` has trailing slashes and backslashes trimmed uniformly.

---

## 3. CLI UX & Output Verification

When running on an Ubuntu node with synchronized database entries containing Windows paths:
```bash
gitmap pa
```
Output:
```text
  ✓ SQLite: Auto-healed 4 cross-OS repository path(s) to Unix slashes before pull.
  ✓ SQLite: Optimized 0 redundant repository record(s).
  
  Pulling 38 repositories across 8 workers...
  [████████████████████████████████████████] 100% (38/38)
```

---

## 4. Verification & Quality Gates

1. **Linter & Guideline Verification:**
   - `python linter-scripts/check-relative-paths.py`
   - `python 03-ai-scripts/05-guideline-autofixer.py cli/store cli/cmdpull --check-only`
2. **Behavioral Acceptance:**
   - Windows-style paths (`D:\foo\bar`) stored in `gitmap.db` are automatically healed to forward slashes on Linux before workers spawn.
   - Repositories whose paths were repaired pull successfully without triggering `missing repository directory`.
