# Subtask 03: Clone Pretty Runner Error Handling & Function Length Enforcement

> **Subtask ID:** Subtask-03  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-03.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Audit and refactor `cli/cmdclone/clonepretty.go` to strictly adhere to leaf error sovereignty across all clone execution paths.
2. Ensure that `runCloneCommandPretty` and associated helper routines never swallow errors or invoke exit handlers (`cliexit.HandleError`, `os.Exit`).
3. Decompose functions in `clonepretty.go` that approach or exceed the 15-line limit (e.g., `buildClonePrettyHints` and `printClonePrettyFailure`).
4. Enforce strict boolean conventions (`isCloneAssumeYes`, `isSSHCloneURL`) and positive polarity checks.
5. Guarantee vertical whitespace hygiene (blank line after `}` and blank line before `return`).

---

## 2. Detailed Technical Plan

### 2.1 Enforce Leaf Error Sovereignty in `runCloneCommandPretty`
- In `cli/cmdclone/clonepretty.go`:
  - Verify `runCloneCommandPretty(url, dest string) error` returns `runErr` directly upon command failure.
  - Ensure failure logging in `printClonePrettyFailure` does not trigger process termination.
  - Return `nil` when dry-run mode is active with explicit blank lines before return.

### 2.2 Function Decomposition ($\le 15$ Lines Limit)
- Locate `buildClonePrettyHints(url, dest string) string` (currently 19 lines):
  - Decompose into two targeted helpers:
    1. `appendTransportCloneHints(hints []string, url string) []string`: handles HTTPS, SSH, and host-key acceptance suggestions ($\le 12$ lines).
    2. `buildClonePrettyHints(url, dest string) string`: initializes default cleanup/retry hints and joins lines ($\le 10$ lines).
- Locate `printClonePrettyFailure(url, dest string, runErr error, elapsed time.Duration)`:
  - Extract exit code resolution into a separate helper: `resolveExitCodeFromError(err error) int` ($\le 8$ lines).
  - Keep `printClonePrettyFailure` strictly focused on formatting the red banner and hint panel ($\le 12$ lines).

### 2.3 Positive Boolean Naming & Spacing
- Ensure `isSSHCloneURL` and `isCloneAssumeYes` maintain positive polarity.
- Add blank line before every `return` statement:
  ```go
  if runErr != nil {
      printClonePrettyFailure(url, dest, runErr, elapsed)

      return runErr
  }
  ```

---

## 3. Verification Criteria

- [ ] Zero calls to `cliexit.HandleError` or `os.Exit` exist within `cli/cmdclone/clonepretty.go`.
- [ ] Every function in `cli/cmdclone/clonepretty.go` is $\le 15$ lines.
- [ ] `buildClonePrettyHints` is cleanly decomposed into focused sub-functions.
- [ ] Subprocess exit codes are resolved via dedicated helper without mutating outer process exit state.
- [ ] Blank line precedes every `return` and follows every closing brace `}`.
