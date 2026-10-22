# Subtask 01: Direct Clone Helper Error Return Sovereignty & Propagation

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-01.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope

1. Modernize `cli/cmdclone/clone_direct.go` by eliminating all leaf invocations of `cliexit.HandleError` and `os.Exit`.
2. Convert void helper routines (`cloneDirectURL`, `validateDirectURLTarget`, and execution wrappers in `cli/cmdclone/clone_direct.go` / `clonerun.go`) to return explicit `error` values.
3. Bubble errors directly up to the top-level command dispatcher (`runClone` / `cmd.Execute`) so that database telemetry and command audits via `finishCommandAudit` are guaranteed to capture failure details.
4. Enforce strict coding guidelines across modified files:
   - Function limit: $\le 15$ lines per function.
   - Conditionals: Zero nested `if` statements through guard clauses and early returns.
   - Positive booleans: Variables and checks prefixed strictly with `is` or `has` (e.g., `isDirectURL`, `hasExistingDir`).
   - Vertical spacing: Blank line after closing brace `}` and blank line before every `return`.

---

## 2. Detailed Technical Plan

### 2.1 Refactor Target Validation in `cli/cmdclone/clone_direct.go`
- Locate `validateDirectURLTarget(folderName, absPath string, params CloneFlags)`:
  - Remove direct calls to `cliexit.HandleError(nil, 1)` when `isDirectURL(folderName)` triggers the stale binary folder URL error.
  - Return `fmt.Errorf(constants.ErrCloneStaleBinaryFolderURL, folderName, constants.Version)` directly.
  - Remove direct calls to `cliexit.HandleError(nil, 1)` when `params.NoReplace && statErr == nil`.
  - Return `fmt.Errorf(constants.ErrCloneURLExists, absPath)` directly.
- Convert signature from `func validateDirectURLTarget(...)` to `func validateDirectURLTarget(...) error`.

### 2.2 Refactor Direct Clone Execution & Error Return
- Locate `cloneDirectURL(source, absPath, folderName string, params CloneFlags)`:
  - Update signature to `func cloneDirectURL(source, absPath, folderName string, params CloneFlags) error`.
  - Invert pre-condition checks with early returns if validation fails:
    ```go
    if err := validateDirectURLTarget(folderName, absPath, params); err != nil {
        return err
    }
    ```
  - On clone command failure:
    - Instead of calling `cliexit.HandleError(cloneErr, 1)`, emit diagnostic output to `os.Stderr` and return `fmt.Errorf(constants.ErrCloneURLFailed, url, cloneErr)`.
  - Return `nil` upon successful execution.

### 2.3 Refactor Dispatcher Integration in `cli/cmdclone/clonerun.go` / `clone.go`
- Ensure caller receiving `cloneDirectURL(...)` checks the returned error:
  ```go
  if err := cloneDirectURL(source, absPath, folderName, params); err != nil {
      return err
  }

  return nil
  ```
- Guarantee that no intermediate function discards errors with blank identifier `_ = cloneDirectURL(...)`.

---

## 3. Verification Criteria

- [ ] Zero invocations of `cliexit.HandleError` or `os.Exit` exist within `cli/cmdclone/clone_direct.go`.
- [ ] `cloneDirectURL` and `validateDirectURLTarget` declare explicit `error` return types.
- [ ] All functions in `cli/cmdclone/clone_direct.go` satisfy the $\le 15$ lines limit.
- [ ] No nested `if` statements remain; all conditionals use guard clauses.
- [ ] All boolean variables and functions use positive `is*` or `has*` prefixes.
- [ ] Strictly relative Git paths are preserved throughout documentation and comments.
