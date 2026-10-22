# Subtask 05: Temp-Release Void Discard Elimination & Error Bubbling

> **Subtask ID:** Subtask-05  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-05.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Eliminate leaf exit handler calls (`cliexit.HandleError(err, 1)`) in `cli/cmdtemprelease/temprelease.go` and replace with direct `error` returns.
2. Refactor `dispatchTempRelease` from a void function into an explicit `error` return function that captures and propagates errors from subcommands (`runTempReleaseList`, `runTempReleaseRemove`, `runTempReleaseCreate`).
3. Eliminate swallowed JSON encoding errors in `cli/cmdtemprelease/tempreleaselist.go` by bubbling `encodeTempReleaseListJSON` errors instead of printing to stderr and returning `nil`.
4. Ensure all subcommands bubble failures up to `RunTempRelease(args []string) error` for proper database audit logging via `finishCommandAudit`.
5. Enforce strict coding guidelines:
   - Function limit: $\le 15$ lines per function.
   - Conditionals: Zero nested `if` statements with guard clauses.
   - Positive booleans: Variables prefixed strictly with `is` or `has` (`isJSONOutput`, `hasArgs`).
   - Vertical spacing: Blank line after `}` and blank line before every `return`.
   - File length ratchet: Keep files well under the 500-line ceiling.

---

## 2. Detailed Technical Plan

### 2.1 Refactor Argument Validation in `cli/cmdtemprelease/temprelease.go`
- Locate `RunTempRelease(args []string) error`:
  - When `len(args) == 0`:
    - Remove `cliexit.HandleError(err, 1)`.
    - Return the created validation error directly:
      ```go
      if len(args) == 0 {
          return apperror.NewWithDetails(
              "cmd.temprelease.run",
              "E1116",
              "missing required arguments for temp-release",
              "cmd.temprelease",
              apperror.ErrorTypeValidation,
              apperror.SeverityError,
              nil,
          )
      }
      ```
  - Delegate execution and return error directly:
    ```go
    return dispatchTempRelease(args[0], args)
    ```

### 2.2 Refactor `dispatchTempRelease`
- Locate `dispatchTempRelease(sub string, args []string)`:
  - Update signature to `func dispatchTempRelease(sub string, args []string) error`.
  - Invert dispatch branches into guard clauses:
    ```go
    if sub == constants.SubCmdTRList {
        return runTempReleaseList(args[1:])
    }

    if sub == constants.SubCmdTRRemove {
        return runTempReleaseRemove(args[1:])
    }

    return runTempReleaseCreate(args)
    ```

### 2.3 Eliminate Swallowed JSON Encoding Error in `cli/cmdtemprelease/tempreleaselist.go`
- Locate `runTempReleaseList(args []string) error`:
  - Check error returned by `encodeTempReleaseListJSON(os.Stdout, releases)`:
    ```go
    if err := encodeTempReleaseListJSON(os.Stdout, releases); err != nil {
        return apperror.WrapSimple(err, "encode temp release json")
    }

    return nil
    ```
  - Invert non-JSON branch to guard clause:
    ```go
    if !jsonOutput {
        printTRList(releases)

        return nil
    }
    ```

---

## 3. Verification Criteria

- [ ] Zero calls to `cliexit.HandleError` or `os.Exit` exist in `cli/cmdtemprelease/temprelease.go`.
- [ ] `dispatchTempRelease` declares an `error` return and never discards subcommand errors.
- [ ] `runTempReleaseList` returns JSON encoding errors directly rather than swallowing them.
- [ ] `RunTempRelease` bubbles all validation and dispatch errors to root dispatcher.
- [ ] All functions in touched files strictly comply with $\le 15$ lines ceiling.
- [ ] Zero nested `if` statements throughout.
- [ ] Mandatory vertical blank lines after `}` and before `return`.
- [ ] Strictly relative Git paths used throughout all references.
