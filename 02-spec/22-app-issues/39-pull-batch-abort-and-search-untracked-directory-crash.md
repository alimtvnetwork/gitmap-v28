# Issue 39: Pull Batch Abort & Search Untracked Directory Crash

> **Issue Status:** Resolved  
> **Traceability IDs:** Task-01, Task-02, Task-08  
> **Canonical Path:** `02-spec/22-app-issues/39-pull-batch-abort-and-search-untracked-directory-crash.md`  
> **Parent Spec:** `02-spec/22-app-issues/01-index.md`

---

## 1. Reproduction & Symptoms

### Symptom A: Pull Batch Interactive Quit Dumps E9000 Stack Trace
When running `gitmap pull all` and hitting a dirty repository remediation prompt:
```text
Pick [1=stash, 2=wip, 3=discard, s=skip, a=all-stash, q=quit]: q
gitmap: [E9000:EXECUTION] execution: pull batch failed with 1 failure(s)
  origin: cmdpull/pull.go:498
  stack trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdpull.finalizePullBatchTask (cmdpull/pull.go:498)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdpull.executePullBatchLifecycle (cmdpull/pull.go:242)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdpull.runPullBatch (cmdpull/pull.go:207)
```
Selecting `q` (Quit) is an intentional user cancellation, not a catastrophic runtime defect.

### Symptom B: Search in Untracked Directory Dumps E9000 Stack Trace
Running `gitmap search <filename>` from a parent directory (e.g. `D:\work>`):
```text
PS D:\work> gitmap search vmpass.txt
gitmap search: execute failed: [E9000:EXECUTION] search.getRepoDB: current directory is not a tracked gitmap repository. run 'gitmap scan' first (at=cmd/search.go:69)
Stack Trace:
  at github.com/alimtvnetwork/gitmap-v28/cli/cmd.executeSearchCommand (cmd/search.go:69)
```

---

## 2. Root Cause Analysis (RCA)

1. **Pull Batch Finalization:** `finalizePullBatchTask` unconditionally calls `cliexit.HandleError(apperror.NewExecutionError(errMsg), 1)`. In the Go AppError pattern, `HandleError` prints full multiline stack traces intended for unexpected internal errors, terrifying users when they simply chose to abort.
2. **Search Directory Resolution:** `cmd/cmd_db.go:getRepoDB()` queries `mainDB.FindByPath(cwd)`. If `cwd` is not a registered repository, it returns an error wrapped with `search.getRepoDB` AppError, which propagates to `cliexit.HandleError`. It completely lacked a fallback to global indexed files or filesystem discovery.

---

## 3. Grounded Code Fix

1. **Pull Batch:** In `finalizePullBatchTask`, when `failCount > 0`, log the failure cleanly to task DB and print a clean red warning message:
   `fmt.Fprintf(os.Stderr, "pull batch finished with %d failure(s)\n", failCount)`
   Exit with code 1 without invoking `cliexit.HandleError`.
2. **Search Fallback:** In `executeSearchCommand`, if `getRepoDB` fails due to untracked repository, fall back to searching across the workspace filesystem or querying `mainDB.GetAll()` without throwing an AppError stack trace.

---

## 4. Prevention & Invariants

- User cancellations (`q`, `Ctrl+C`) and expected operational outcomes (e.g. some repositories having dirty trees) MUST NEVER trigger `cliexit.HandleError`.
- Commands that search or inspect files MUST provide graceful fallbacks when invoked from untracked folders.
