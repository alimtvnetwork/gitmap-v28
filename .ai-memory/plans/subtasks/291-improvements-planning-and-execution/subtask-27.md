# Subtask Plan 27: Git Detached HEAD Detection (`cli/gitutil/git_repo_query.go`)

> **Subtask ID:** `subtask-27`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/gitutil/git_repo_query.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `IsRepoDetached(repoPath string) (bool, error)` in `cli/gitutil/git_repo_query.go`.
2. Inspect whether the repository at `repoPath` is on a detached HEAD commit rather than an active branch reference.
3. Leverage `git symbolic-ref -q HEAD` which cleanly exits non-zero or produces empty output when detached.
4. Enforce Positive Booleans (Rule 4) with `isDetached`, Function Limit <= 15 lines (Rule 2), and Error Return Sovereignty (Rule 1).

---

## 2. Targeted Problem & Architectural Rationale
When GitMap executes automated operations (such as bisecting, testing tagged releases, or inspecting pipeline commits), repositories may be left in a detached HEAD state.

Subsequent operations (such as `gitmap pull`, `gitmap push`, or `gitmap commit`) will fail or create dangling orphaned commits if they assume an active tracking branch exists.

`IsRepoDetached` gives callers a fast, lightweight check to verify branch attachment before initiating write or sync flows.

---

## 3. Concrete Implementation Details

### 3.1 Function Addition in `cli/gitutil/git_repo_query.go`

```go
// IsRepoDetached reports whether the repository is in a detached HEAD state.
func IsRepoDetached(repoPath string) (bool, error) {
	out, err := RunGitInDir(repoPath, constants.GitSymbolicRef, "-q", constants.GitHEAD)
	if err != nil {
		return true, nil
	}

	isDetached := len(strings.TrimSpace(string(out))) == 0

	return isDetached, nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Direct return of boolean and error without exit handlers.
- **Rule 2 (Function Limit):** 10 lines of code, well within the 15-line limit.
- **Rule 3 (Zero Nested Ifs):** Linear check with early return on non-zero exit code.
- **Rule 4 (Positive Booleans):** Strictly positive naming `isDetached`.
- **Rule 6 (Vertical Whitespace):** Blank lines before `return` and after `}`.

---

## 5. Verification Criteria
- [x] `IsRepoDetached(repoPath string) (bool, error)` is declared in `cli/gitutil/git_repo_query.go`.
- [x] Returns `true` when HEAD is detached (e.g. non-zero exit code from `symbolic-ref`).
- [x] Returns `false` when HEAD points to a valid symbolic reference like `refs/heads/main`.
- [x] Function is <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
