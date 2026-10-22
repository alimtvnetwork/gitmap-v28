# Subtask 04: CPAR Batch Error Capture & Git Add Verification

> **Subtask ID:** Subtask-04  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-04.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Eliminate silent error discards in `cli/cmdcpar/cpar.go` where `git add -A` execution is ignored (`_ = exec.Command(...).Run()`).
2. Add dedicated error logging and state management for `git add` failures before commit attempts.
3. Propagate batch failure state from `dispatchCPARBatch` to `RunCPAR(args []string) *apperror.AppError` so that failed batch operations bubble up to the root dispatcher and are accurately recorded in the database audit log via `finishCommandAudit`.
4. Enforce strict coding guidelines across `cli/cmdcpar/cpar.go`:
   - Function limit: $\le 15$ lines per function ($\le 8$ lines preferred).
   - Conditionals: Zero nested `if` statements with early returns.
   - Positive booleans: Variables named strictly with `is` or `has` (`isSuccess`, `hasFailures`, `hasAddErr`).
   - Vertical whitespace: Blank line after closing brace `}` and blank line before every `return`.
   - File length: Keep `cli/cmdcpar/cpar.go` safely below the 500-line ratchet.

---

## 2. Detailed Technical Plan

### 2.1 Eliminate Swallowed Error on Git Add in `commitAndMaybePush`
- Locate `commitAndMaybePush(d DirtyRepoSummary, opts cparOptions) bool`:
  - Decompose into `executeGitAdd(d DirtyRepoSummary) error` and `executeGitCommit(d DirtyRepoSummary, opts cparOptions) (string, error)`.
  - In `executeGitAdd`:
    ```go
    func executeGitAdd(d DirtyRepoSummary) error {
        cmd := exec.Command("git", "-C", d.RepoPath, "add", "-A")

        return cmd.Run()
    }
    ```
  - In `commitAndMaybePush`:
    ```go
    if err := executeGitAdd(d); err != nil {
        logAddFailure(d, err)

        return false
    }
    ```
  - Implement `logAddFailure(d DirtyRepoSummary, err error)` to record an `InternalErrorRecord` with code `"E_CPAR_ADD"`.

### 2.2 Propagate Batch Execution State to `RunCPAR`
- Locate `dispatchCPARBatch(dirtyList []DirtyRepoSummary, opts cparOptions)`:
  - Update signature to `func dispatchCPARBatch(dirtyList []DirtyRepoSummary, opts cparOptions) *apperror.AppError`.
  - Capture `runState := executeCPARAcrossDirty(dirtyList, opts)`.
  - Finalize task queue and close DB.
  - Invert condition with guard clause:
    ```go
    if runState.HasFailures {
        return apperror.NewWithDetails(
            "cmd.cpar.batch",
            "E_CPAR_BATCH_FAILED",
            fmt.Sprintf("CPAR completed with %d failure(s)", runState.FailureCount),
            "cmd.cpar",
            apperror.ErrorTypeExecution,
            apperror.SeverityError,
            map[string]any{"failures": runState.FailureCount, "successes": runState.SuccessCount},
        )
    }

    return nil
    ```
- Update `RunCPAR(args []string) *apperror.AppError`:
  - Check return of `dispatchCPARBatch(dirtyList, opts)`:
    ```go
    return dispatchCPARBatch(dirtyList, opts)
    ```
  - When batch operations fail, `finishCommandAudit` accurately records non-zero audit status.

### 2.3 Whitespace Hygiene & Function Sizing
- Audit all functions in `cli/cmdcpar/cpar.go`:
  - Add missing blank line after closing brace `}`.
  - Add missing blank line before `return` statements across `updateRunState`, `handlePush`, `pushToRemote`.
  - Ensure every decomposed function does not exceed 15 lines.

---

## 3. Verification Criteria

- [ ] Blank identifier discard `_ = exec.Command(..., "add", "-A").Run()` is completely eliminated.
- [ ] `git add` failures are captured, logged to internal error store, and abort the commit step for that repository.
- [ ] `dispatchCPARBatch` returns `*apperror.AppError` and propagates failure counts.
- [ ] `RunCPAR` returns a non-nil `*apperror.AppError` if any repository fails during batch commit/push.
- [ ] All functions in `cli/cmdcpar/cpar.go` strictly comply with $\le 15$ lines limit.
- [ ] Zero nested `if` statements throughout.
- [ ] Vertical blank lines strictly present after `}` and before `return`.
- [ ] Strictly relative Git paths used throughout all references.
