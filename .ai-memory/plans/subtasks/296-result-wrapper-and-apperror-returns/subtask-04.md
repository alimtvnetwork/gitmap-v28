# Subtask 04: Author cli/gitutil/types.go and Migrate git_repo_query.go & git_repo_commits.go to Result Containers

## Status: Pending
- **Parent Plan:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/01-architecture-spec.md`
- **Component Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/02-component-spec.md`
- **Target Files:**
  - `cli/gitutil/types.go`
  - `cli/gitutil/git_repo_query.go`
  - `cli/gitutil/git_repo_commits.go`
- **Objective:** Establish domain collection type aliases in `cli/gitutil/types.go` (`CaseCollisionSliceWrap`, `RemoteBranchSliceWrap`) and migrate all multi-value error return tuples `(T, error)` in `git_repo_query.go` and `git_repo_commits.go` to strongly-typed single-return Result containers (`result.StringResult`, `result.BoolResult`, `result.IntResult`, `result.StringSliceResult`).

---

## Action Items

1. **Author `cli/gitutil/types.go`:**
   - Define dedicated slice collection types: `type CaseCollisionSlice = []CaseCollision` and `type RemoteBranchSlice = []RemoteBranchInfo`.
   - Alias single reusable Result containers: `type CaseCollisionSliceWrap = result.Wrap[CaseCollisionSlice]` and `type RemoteBranchSliceWrap = result.Wrap[RemoteBranchSlice]`.
   - Prevent package stutter and eliminate raw generic bracket leakage across caller signatures.

2. **Refactor `cli/gitutil/git_repo_query.go`:**
   - Update `GetRepoHeadCommitSha(repoPath string)` to return `result.StringResult`. Guard on `res.IsFailure()`, returning `result.Fail[string](res.AppError())`.
   - Update `GetRepoRemoteURL(repoPath, remoteName string)` to return `result.StringResult`, delegating to `GetRepoConfig`.
   - Update `GetRepoTrackingRemote(repoPath string)` to return `result.StringResult`. Guard on `res.IsFailure()`, returning `result.Fail[string](res.AppError())`.
   - Update `HasUnpushedCommits(repoPath string)` to return `result.BoolResult`. Guard on `res.IsFailure()`, returning `result.Fail[bool](res.AppError())`.
   - Update `GetRepoDirtyFiles(repoPath string)` to return `result.StringSliceResult`. Guard on `res.IsFailure()`, returning `result.FailSlice[string](res.AppError())`.
   - Update `IsRepoDetached(repoPath string)` to return `result.BoolResult`. Preserve fallback logic returning `result.Ok(true)` on symbolic-ref miss.

3. **Refactor `cli/gitutil/git_repo_commits.go`:**
   - Update `GetRepoUnpushedCommitCount(repoPath string)` to return `result.IntResult`. Guard on command failure with `result.Fail[int](res.AppError())` and parse errors with `result.Fail[int](apperror.WrapSimple(parseErr, "GetRepoUnpushedCommitCount"))`.
   - Update `GetRepoMergeBase(repoPath, ref1, ref2 string)` to return `result.StringResult`. Guard on `res.IsFailure()`, returning `result.Fail[string](res.AppError())`.

4. **Verify Coding Invariants:**
   - Ensure every refactored function stays $\le 15$ lines.
   - Maintain zero nested `if` statements by using guard clauses and early returns.
   - Enforce vertical whitespace rules: blank line before `return` and after `}`.
   - Strictly positive boolean names (`is...`, `has...`).
   - Zero compiler builds or unit tests during authoring.

---

## Implementation Details

### 1. `cli/gitutil/types.go`
```go
package gitutil

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// CaseCollisionSlice defines a collection of case collision records.
	CaseCollisionSlice = []CaseCollision

	// CaseCollisionSliceWrap is the canonical single reusable container for CaseCollision collections.
	CaseCollisionSliceWrap = result.Wrap[CaseCollisionSlice]

	// RemoteBranchSlice defines a collection of remote branch metadata.
	RemoteBranchSlice = []RemoteBranchInfo

	// RemoteBranchSliceWrap is the canonical single reusable container for remote branch collections.
	RemoteBranchSliceWrap = result.Wrap[RemoteBranchSlice]
)
```

### 2. `cli/gitutil/git_repo_query.go`
```go
package gitutil

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// GetRepoHeadCommitSha retrieves the full HEAD commit SHA of the repository at repoPath.
func GetRepoHeadCommitSha(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitHEAD)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// GetRepoRemoteURL retrieves the configured URL for the specified remoteName at repoPath.
func GetRepoRemoteURL(repoPath, remoteName string) result.StringResult {
	key := fmt.Sprintf("remote.%s.url", remoteName)

	return GetRepoConfig(repoPath, key)
}

// GetRepoTrackingRemote retrieves the upstream tracking reference for the active branch.
func GetRepoTrackingRemote(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, "@{u}")
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// HasUnpushedCommits reports whether the repository has commits not yet pushed to upstream.
func HasUnpushedCommits(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, "cherry", "-v")
	if res.IsFailure() {
		return result.Fail[bool](res.AppError())
	}

	hasCommits := len(strings.TrimSpace(string(res.Data))) > 0

	return result.Ok(hasCommits)
}

// GetRepoDirtyFiles returns relative paths of uncommitted, modified, or untracked files in repoPath.
func GetRepoDirtyFiles(repoPath string) result.StringSliceResult {
	res := RunGitInDir(repoPath, constants.GitStatus, constants.GitPorcelainFlag)
	if res.IsFailure() {
		return result.FailSlice[string](res.AppError())
	}

	files := parsePorcelainDirtyLines(string(res.Data))

	return result.OkSlice(files)
}

func parsePorcelainDirtyLines(output string) []string {
	var files []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		clean := strings.TrimSpace(line)
		if len(clean) > 3 {
			files = append(files, strings.TrimSpace(clean[3:]))
		}
	}

	return files
}

// IsRepoDetached reports whether the repository is in a detached HEAD state.
func IsRepoDetached(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, constants.GitSymbolicRef, constants.GitQuietFlag, constants.GitHEAD)
	if res.IsFailure() {
		return result.Ok(true)
	}

	isDetached := len(strings.TrimSpace(string(res.Data))) == 0

	return result.Ok(isDetached)
}
```

### 3. `cli/gitutil/git_repo_commits.go`
```go
package gitutil

import (
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// GetRepoUnpushedCommitCount returns the number of unpushed commits ahead of upstream tracking.
func GetRepoUnpushedCommitCount(repoPath string) result.IntResult {
	res := RunGitInDir(repoPath, "rev-list", "--count", "@{u}..HEAD")
	if res.IsFailure() {
		return result.Fail[int](res.AppError())
	}

	clean := strings.TrimSpace(string(res.Data))
	count, parseErr := strconv.Atoi(clean)
	if parseErr != nil {
		return result.Fail[int](apperror.WrapSimple(parseErr, "GetRepoUnpushedCommitCount"))
	}

	return result.Ok(count)
}

// GetRepoMergeBase finds the best common ancestor commit between two references.
func GetRepoMergeBase(repoPath, ref1, ref2 string) result.StringResult {
	res := RunGitInDir(repoPath, "merge-base", ref1, ref2)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}
```

---

## Verification & Conformance Checklist

- [ ] `cli/gitutil/types.go` declares `CaseCollisionSliceWrap` and `RemoteBranchSliceWrap`.
- [ ] No raw generic brackets leaked across `git_repo_query.go` or `git_repo_commits.go` public APIs.
- [ ] All functions return single Result containers (`result.StringResult`, `result.BoolResult`, `result.IntResult`, `result.StringSliceResult`).
- [ ] Failure paths return `res.AppError()` through `result.Fail[T]` or `result.FailSlice[T]`.
- [ ] Zero compiler builds and zero unit tests executed during plan creation.
- [ ] Strictly relative Git paths utilized throughout all references.
