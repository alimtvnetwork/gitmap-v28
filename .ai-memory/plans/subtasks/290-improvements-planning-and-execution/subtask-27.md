# Subtask Plan 27: Unpushed Commits Status Helper (`cli/gitutil/git_cmd_exec.go`)

> **Subtask ID:** `subtask-27`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/gitutil/git_cmd_exec.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `HasUnpushedCommits(repoPath string) (bool, error)` in `cli/gitutil/git_cmd_exec.go`.
2. Determine whether the current local branch has committed changes that have not yet been pushed to the upstream remote using `git cherry -v`.
3. Enforce positive boolean naming (`hasCommits`) in compliance with Rule 4.
4. Adhere strictly to Error Return Sovereignty (Rule 1) and Function Length Limits (<= 15 lines, Rule 2).

---

## 2. Targeted Problem & Architectural Rationale
GitMap operations (such as multi-repo pull, branch switching, release tagging, and health checks) need to know if a repository has unpushed commits before executing disruptive operations.

Previously, subsystems either:
1. Checked only `IsRepoClean` (working tree dirtiness), which reports clean even when 10 unpushed commits exist locally.
2. Swallowed errors and guessed push state based on branch name comparisons.

`git cherry -v` compares the current branch HEAD against the upstream tracking branch and outputs lines beginning with `+ <hash> <subject>` for unpushed commits. If the output is non-empty, unpushed commits exist.

---

## 3. Concrete Implementation Details

### 3.1 Implementation in `cli/gitutil/git_cmd_exec.go`
Add `HasUnpushedCommits`:

```go
// HasUnpushedCommits reports whether the local branch contains commits not yet pushed to upstream.
func HasUnpushedCommits(repoPath string) (bool, error) {
	out, err := RunGitInDir(repoPath, "cherry", "-v")
	if err != nil {
		return false, err
	}

	hasCommits := len(strings.TrimSpace(string(out))) > 0

	return hasCommits, nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Errors from `RunGitInDir` are returned directly (`return false, err`).
- **Rule 2 (Function Limit):** 10 lines of code (well under the 15-line limit).
- **Rule 3 (Zero Nested Ifs):** Single flat guard clause.
- **Rule 4 (Positive Boolean Naming):** The boolean variable is named `hasCommits` (`has` prefix, positive polarity). Banned prefixes (`can`, `should`, `isNot*`, `unpushed`) are completely avoided.
- **Rule 6 (Vertical Whitespace):** Single blank line after `}` and before final `return`.

---

## 5. Verification Criteria
- [x] Exported function `HasUnpushedCommits(repoPath string) (bool, error)` exists in `cli/gitutil/git_cmd_exec.go`.
- [x] Executes `RunGitInDir(repoPath, "cherry", "-v")`.
- [x] Uses positive boolean variable `hasCommits`.
- [x] Returns `(true, nil)` when cherry output contains lines, `(false, nil)` when empty.
- [x] 10 lines of code.
- [x] All file paths are strictly relative lowercase Git paths.
