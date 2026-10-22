# Subtask Plan 27: Common Ancestor Merge Base Resolution (`cli/gitutil/git_repo_commits.go`)

> **Subtask ID:** `subtask-27`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/gitutil/git_repo_commits.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `GetRepoMergeBase(repoPath, ref1, ref2 string) (string, error)` in `cli/gitutil/git_repo_commits.go`.
2. Execute `git merge-base ref1 ref2` via `RunGitInDir` to discover the best common ancestor commit between two branches or references.
3. Trim whitespace and return the resulting SHA commit hash.
4. Co-locate merge base discovery with commit counting in `cli/gitutil/git_repo_commits.go`.
5. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
Commands handling multi-branch synchronization, divergence inspection, or branch reconciliation frequently need to determine the common ancestor between two revisions (such as `origin/main` and a local feature branch).

Currently, there is no centralized helper for `git merge-base`, resulting in ad-hoc subprocess invocations or duplicate shell wrappers across different CLI packages.

Adding `GetRepoMergeBase`:
1. Provides a single, tested, error-sovereign helper for discovering branch divergence points.
2. Supports branch safety checks before automated rebasing or merging.
3. Keeps git query logic organized within the `gitutil` package without bloating `git_repo_query.go`.

---

## 3. Concrete Implementation Details

### 3.1 Function Implementation in `cli/gitutil/git_repo_commits.go`
Add `GetRepoMergeBase` following `GetRepoUnpushedCommitCount`:

```go
// GetRepoMergeBase finds the best common ancestor commit between two references.
func GetRepoMergeBase(repoPath, ref1, ref2 string) (string, error) {
	out, err := RunGitInDir(repoPath, "merge-base", ref1, ref2)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Subprocess execution failure bubbles directly as `("", err)` to the caller.
- **Rule 2 (Function Limit):** Function is exactly 7 lines, well below the 8-line preferred target.
- **Rule 3 (Zero Nested Ifs):** Single linear error guard with early return.
- **Rule 4 (Positive Booleans):** No boolean parameters or returns.
- **Rule 6 (Vertical Whitespace):** Blank lines after `}` and before `return`.

---

## 5. Verification Criteria
- [x] Method `GetRepoMergeBase(repoPath, ref1, ref2 string) (string, error)` declared in `cli/gitutil/git_repo_commits.go`.
- [x] Invokes `RunGitInDir` with `"merge-base"`, `ref1`, `ref2`.
- [x] Trims output string to return clean SHA string.
- [x] Returns error if references do not share a common ancestor or command fails.
- [x] Function is 7 lines with zero nested ifs.
