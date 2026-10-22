# Subtask 02: Has-Change Error Propagation & Dirty Check Sovereignty

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-02.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Eliminate error swallowing in working tree dirty checks inside `cli/cmdhaschange/helpers_shim.go` by converting `isWorkingTreeDirty` from a bare boolean return to `(bool, error)`.
2. Modernize `cli/cmdhaschange/haschange.go` by removing all leaf invocations of exit handlers (`cliexit.HandleUsageError`, `cliexit.HandleError`).
3. Replace void reporting helpers (`printHasChangeOne`, `printHasChangeAll`) with functions returning explicit `error` / `*apperror.AppError`.
4. Bubble all flag parsing, git subprocess, and invalid mode errors up to `RunHasChange(args []string) error` for proper database audit logging.
5. Enforce strict coding guidelines:
   - Function limit: $\le 15$ lines per function.
   - Conditionals: Zero nested `if` statements through guard clauses.
   - Positive booleans: Variables named strictly with `is*` or `has*` (`isDirty`, `hasChange`, `hasUpstream`).
   - Vertical spacing: Blank line after `}` and blank line before every `return`.

---

## 2. Detailed Technical Plan

### 2.1 Refactor `isWorkingTreeDirty` in `cli/cmdhaschange/helpers_shim.go`
- Locate `isWorkingTreeDirty(target string)`:
  - Update signature from `func isWorkingTreeDirty(target string) bool` to `func isWorkingTreeDirty(target string) (bool, error)`.
  - Invert error check: If `cmd.Output()` fails, wrap and return `false, apperror.WrapSimple(err, "git status porcelain in "+target)`.
  - On success, return `len(strings.TrimSpace(string(out))) > 0, nil`.
  - Ensure vertical blank line before `return`.

### 2.2 Refactor Reporting Helpers in `cli/cmdhaschange/haschange.go`
- Locate `printHasChangeOne(target, mode string)`:
  - Update signature to `func printHasChangeOne(target, mode string) error`.
  - On `constants.HCModeDirty`, call `isWorkingTreeDirty(target)`. If error returned, propagate error directly.
  - On `default:` invalid mode branch, remove `cliexit.HandleUsageError(nil)` and return `apperror.NewUsageError(fmt.Sprintf(constants.ErrHCBadMode, mode))`.
  - Return `nil` after printing boolean string.
- Locate `printHasChangeAll(target string)`:
  - Update signature to `func printHasChangeAll(target string) error`.
  - Check error from `isWorkingTreeDirty(target)`:
    ```go
    dirty, err := isWorkingTreeDirty(target)
    if err != nil {
        return err
    }
    ```
  - Return `nil` upon successful formatting.

### 2.3 Refactor Flag Parsing & Dispatcher in `cli/cmdhaschange/haschange.go`
- Locate `parseHasChangeFlags(args []string)`:
  - Update signature to `func parseHasChangeFlags(args []string) (alias, mode string, all, fetch bool, err error)`.
  - If `fs.Parse` fails, return `"", "", false, false, apperror.WrapSimple(err, "parse has-change flags")` instead of invoking `cliexit.HandleError(err, 2)`.
- In `RunHasChange(args []string) error`:
  - Check returned error from `parseHasChangeFlags`.
  - Guard against empty alias without calling `cliexit.HandleUsageError(nil)`:
    ```go
    if len(alias) == 0 {
        return apperror.NewUsageError(constants.ErrHCUsage)
    }
    ```
  - Propagate errors returned from `printHasChangeAll(target)` and `printHasChangeOne(target, mode)`.

---

## 3. Verification Criteria

- [x] Zero calls to `cliexit.HandleUsageError`, `cliexit.HandleError`, or `os.Exit` exist in `cli/cmdhaschange/`.
- [x] `isWorkingTreeDirty` declares `(bool, error)` return and never silently swallows git output errors.
- [x] `printHasChangeOne` and `printHasChangeAll` return `error` directly.
- [x] `RunHasChange` bubbles all errors to the root command dispatcher.
- [x] All functions in `cli/cmdhaschange/` are $\le 15$ lines.
- [x] Zero nested `if` statements throughout.
- [x] Blank lines strictly maintained after `}` and before `return`.
- [x] Strictly relative Git paths used throughout all references.
