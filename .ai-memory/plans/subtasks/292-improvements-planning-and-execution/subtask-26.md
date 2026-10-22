# Subtask Plan 26: Subprocess Git Unpushed Commit Counter (`cli/gitutil/git_repo_commits.go`)

> **Subtask ID:** `subtask-26`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/gitutil/git_repo_commits.go`  
> **Status:** Pending  

---

## 1. Objectives & Scope
1. Implement `GetRepoUnpushedCommitCount(repoPath string) (int, error)` in `cli/gitutil/git_repo_commits.go`.
2. Execute `git rev-list --count @{u}..HEAD` in the target repository directory using `RunGitInDir`.
3. Parse the numeric stdout count into an integer.
4. House revision calculation helpers in `cli/gitutil/git_repo_commits.go` to prevent `cli/gitutil/git_repo_query.go` (84 lines) from exceeding the 100-line cap.
5. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
GitMap currently checks whether a repository has unpushed commits using `HasUnpushedCommits` via `git cherry -v` in `cli/gitutil/git_repo_query.go`. However, this only returns a boolean. Many CLI commands (such as `cmdstatus`, `cmdrelease`, and visual dashboard summaries) need the exact integer count of pending commits.

Furthermore, `git_repo_query.go` is 84 lines. Adding commit counting and merge base helpers directly to that file would push it past 100 lines.

Extracting commit revision helpers into `cli/gitutil/git_repo_commits.go`:
1. Provides exact unpushed commit metrics across repositories.
2. Keeps `git_repo_query.go` clean and within its 100-line limit.
3. Provides reusable revision metrics for status rendering and release gates.

---

## 3. Concrete Implementation Details

### 3.1 New File `cli/gitutil/git_repo_commits.go`
```go
package gitutil

import (
	"strconv"
	"strings"
)

// GetRepoUnpushedCommitCount returns the number of unpushed commits ahead of upstream tracking.
func GetRepoUnpushedCommitCount(repoPath string) (int, error) {
	out, err := RunGitInDir(repoPath, "rev-list", "--count", "@{u}..HEAD")
	if err != nil {
		return 0, err
	}

	clean := strings.TrimSpace(string(out))
	count, parseErr := strconv.Atoi(clean)
	if parseErr != nil {
		return 0, parseErr
	}

	return count, nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Subprocess errors and integer parse errors return directly as `(0, err)` or `(0, parseErr)`. No exit handlers invoked.
- **Rule 2 (Function Limit):** Function is exactly 12 lines, strictly compliant with <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Linear error checks with immediate early returns.
- **Rule 4 (Positive Booleans):** No boolean parameters or returns; error-checked flow.
- **Rule 6 (Vertical Whitespace):** Blank lines after `}` and before `return`.

---

## 5. Verification Criteria
- [ ] File `cli/gitutil/git_repo_commits.go` created in package `gitutil`.
- [ ] Method `GetRepoUnpushedCommitCount(repoPath string) (int, error)` declared.
- [ ] Invokes `RunGitInDir` with arguments `"rev-list"`, `"--count"`, `"@{u}..HEAD"`.
- [ ] Parses integer count from trimmed output via `strconv.Atoi`.
- [ ] Returns 0 and error when command fails or parsing fails.
- [ ] File size is ~30 lines, well within 100-line cap.
