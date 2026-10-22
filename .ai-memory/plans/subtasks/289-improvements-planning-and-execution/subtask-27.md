# Subtask 27: Positive Repository Clean State Inspection

> **Subtask ID:** Subtask-27  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/gitutil/git_cmd_exec.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `IsRepoClean(repoPath string) (bool, error)` in `cli/gitutil/git_cmd_exec.go`.
2. Inspect repository working tree status via `git status --porcelain`.
3. Enforce positive boolean naming convention (`isClean`) with zero negative polarity booleans.
4. Return `(true, nil)` when working tree has zero untracked or modified files, and `(false, nil)` when changes exist.
5. Adhere strictly to error return sovereignty, function length limits (<= 15 lines), and zero nested conditionals.

---

## 2. Detailed Technical Plan

### 2.1 Implementation in `cli/gitutil/git_cmd_exec.go`
- Add the exported function:
  ```go
  // IsRepoClean reports whether the working tree at repoPath has no uncommitted changes.
  func IsRepoClean(repoPath string) (bool, error) {
  	out, err := RunGitInDir(repoPath, "status", "--porcelain")
  	if err != nil {
  		return false, err
  	}

  	isClean := len(strings.TrimSpace(string(out))) == 0

  	return isClean, nil
  }
  ```

### 2.2 Constant Constants Integration
- If constants for `"status"` and `"--porcelain"` are added to `cli/constants/constants_git.go`:
  - `GitStatusCmd = "status"`
  - `GitPorcelainFlag = "--porcelain"`
  Use them or standard literals consistently.

---

## 3. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Subprocess error is returned directly (`return false, err`).
- **Rule 2 (Function Limit):** 10 lines total (<= 15 line limit).
- **Rule 3 (Zero Nested Ifs):** Single guard clause on error.
- **Rule 4 (Positive Boolean):** Variable and function named `isClean` / `IsRepoClean` (banned: `isDirty`, `hasUncommitted`).
- **Rule 6 (Vertical Whitespace):** Single blank line after error check and before return.

---

## 4. Verification Criteria
- [x] `IsRepoClean` compiles cleanly under Go 1.22+.
- [x] Returns `(true, nil)` for a clean working tree.
- [x] Returns `(false, nil)` when working tree has modified or untracked files.
- [x] Propagates error when `repoPath` is invalid.
- [x] Verified zero negative boolean identifiers.
