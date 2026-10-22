# Subtask 07: Release Leaf Error Sovereignty & Subcommand Return Checking

> **Subtask ID:** Subtask-07  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-07.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Eliminate leaf exit handler invocations (`cliexit.HandleUsageError(nil)`) in `cli/cmdrelease/releasepending.go`.
2. Convert `rejectVersionArgOnPending` from a void exit-calling routine to a clean `error` return function and propagate errors up to `RunReleasePending`.
3. Eliminate swallowed execution errors in `cli/cmdrelease/releaserecentclone.go` where `RunRelease(args)` return value is discarded, converting `tryRunReleaseInRecentClone` to return `(bool, error)`.
4. Update callers in `cli/cmdrelease/release.go` to handle and bubble errors from the recent clone fallback routine.
5. Enforce strict coding guidelines:
   - Function limit: $\le 15$ lines per function.
   - Conditionals: Zero nested `if` statements with guard clauses.
   - Positive booleans: Variables named strictly with `is` or `has` (`isRunningProjectsAlias`, `isInsideRepo`).
   - Vertical spacing: Blank line after closing brace `}` and blank line before every `return`.
   - File length: Keep files safely below the 500-line ceiling.

---

## 2. Detailed Technical Plan

### 2.1 Refactor Version Argument Rejection in `cli/cmdrelease/releasepending.go`
- Locate `rejectVersionArgOnPending(args []string)`:
  - Update signature from `func rejectVersionArgOnPending(args []string)` to `func rejectVersionArgOnPending(args []string) error`.
  - Remove invocation of `cliexit.HandleUsageError(nil)`.
  - Return formatted usage error directly:
    ```go
    func rejectVersionArgOnPending(args []string) error {
        for _, a := range args {
            if strings.HasPrefix(a, "-") || !versionLikeArgPattern.MatchString(a) {
                continue
            }

            return apperror.NewUsageError(fmt.Sprintf(
                "gitmap release-pending takes no version argument (got %q); did you mean: gitmap pr %s or gitmap release %s",
                a, a, a,
            ))
        }

        return nil
    }
    ```
- In `RunReleasePending(args []string) error`:
  - Check error returned by `rejectVersionArgOnPending(args)`:
    ```go
    if err := rejectVersionArgOnPending(args); err != nil {
        return err
    }

    return executeReleasePendingRun(args)
    ```

### 2.2 Refactor Recent Clone Fallback in `cli/cmdrelease/releaserecentclone.go`
- Locate `tryRunReleaseInRecentClone(args []string)`:
  - Update signature from `func tryRunReleaseInRecentClone(args []string) bool` to `func tryRunReleaseInRecentClone(args []string) (bool, error)`.
  - Invert pre-condition checks with guard clauses:
    ```go
    target, ok := lookupRecentClone()
    if !ok {
        return false, nil
    }
    ```
  - After directory switch and git repo validation:
    - Capture error from `RunRelease(args)`:
      ```go
      err := RunRelease(args)

      return true, err
      ```
- Update caller in `cli/cmdrelease/release.go`:
  - Check returned `(bool, error)` tuple:
    ```go
    isHandled, err := tryRunReleaseInRecentClone(args)
    if err != nil {
        return err
    }

    if isHandled {
        return nil
    }
    ```

---

## 3. Verification Criteria

- [ ] Zero calls to `cliexit.HandleUsageError`, `cliexit.HandleError`, or `os.Exit` exist in `cli/cmdrelease/releasepending.go`.
- [ ] `rejectVersionArgOnPending` returns an explicit `error` on invalid positional version arguments.
- [ ] `tryRunReleaseInRecentClone` returns `(bool, error)` and propagates `RunRelease` execution errors.
- [ ] Callers in `cli/cmdrelease/release.go` capture and bubble errors to root dispatcher.
- [ ] All functions in touched files strictly comply with $\le 15$ lines ceiling.
- [ ] Zero nested `if` statements throughout.
- [ ] Mandatory vertical blank lines after `}` and before `return`.
- [ ] Strictly relative Git paths used throughout all references.
