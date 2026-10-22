# Subtask 26: Current Branch Inspection Helper

> **Subtask ID:** Subtask-26  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/gitutil/git_cmd_exec.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `GetRepoCurrentBranch(repoPath string) (string, error)` in `cli/gitutil/git_cmd_exec.go`.
2. Execute `git rev-parse --abbrev-ref HEAD` in `repoPath` using safe directory execution primitives (`RunGitInDir`).
3. Return trimmed branch name and propagate any execution error directly to the caller.
4. Adhere strictly to error return sovereignty, function length limits (<= 15 lines), and zero nested conditionals.

---

## 2. Detailed Technical Plan

### 2.1 Implementation in `cli/gitutil/git_cmd_exec.go`
- Add the exported function:
  ```go
  // GetRepoCurrentBranch returns the active branch name of the repository at repoPath.
  func GetRepoCurrentBranch(repoPath string) (string, error) {
  	out, err := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
  	if err != nil {
  		return "", err
  	}

  	return strings.TrimSpace(string(out)), nil
  }
  ```

### 2.2 Reusable Constant Integration
- Leverages existing typed constants in `cli/constants/constants_git.go`:
  - `constants.GitRevParse` (`"rev-parse"`)
  - `constants.GitAbbrevRef` (`"--abbrev-ref"`)
  - `constants.GitHEAD` (`"HEAD"`)

---

## 3. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Error is returned directly to caller (`return "", err`).
- **Rule 2 (Function Limit):** 8 lines total (well within <= 15 line limit).
- **Rule 3 (Zero Nested Ifs):** Single guard check on `err != nil`.
- **Rule 6 (Vertical Whitespace):** Single blank line after error block and before return.

---

## 4. Verification Criteria
- [x] `GetRepoCurrentBranch` compiles cleanly under Go 1.22+.
- [x] Returns the branch name (e.g. `"main"`, `"master"`, `"develop"`) for a valid git repository.
- [x] Propagates error when `repoPath` is not a git repository or HEAD is detached.
- [x] Does not invoke CLI exit handlers or discard errors.
