# Subtask 05: Refactor git_cmd_exec.go and gitutil.go to Return Monadic Result Containers

## Status: Completed
- **Parent Plan:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/01-architecture-spec.md`
- **Component Spec:** `02-spec/21-app/296-result-wrapper-and-apperror-returns/02-component-spec.md`
- **Target Files:**
  - `cli/gitutil/git_cmd_exec.go`
  - `cli/gitutil/gitutil.go`
- **Objective:** Eliminate multi-value error tuples `([]byte, error)`, `(string, error)`, and `(bool, error)` across process execution and high-level git metadata utilities, migrating them to strongly-typed monadic Result containers (`result.ByteSliceResult`, `result.StringResult`, `result.BoolResult`).

---

## Action Items

1. **Refactor Low-Level Execution in `cli/gitutil/git_cmd_exec.go`:**
   - Update `RunGitInDir(repoPath string, args ...string)` to return `result.ByteSliceResult`. Wrap process execution failures with `result.Fail[[]byte](apperror.WrapSimple(err, "RunGitInDir"))` and success with `result.Ok(out)`.
   - Update `RunGitCombinedInDir(repoPath string, args ...string)` to return `result.ByteSliceResult`. Wrap combined output failures with `result.Fail[[]byte](apperror.WrapSimple(err, "RunGitCombinedInDir"))` and success with `result.Ok(out)`.
   - Update `GetRepoConfig(repoPath, key string)` to return `result.StringResult`. Check `res.IsFailure()`, propagate `res.AppError()`, and return trimmed output in `result.Ok`.
   - Update `ReadRemoteOriginURL(repoPath string)` to return `result.StringResult`, directly delegating to `GetRepoConfig(repoPath, constants.GitRemoteOrigin)`.
   - Update `GetRemoteOriginURL(repoPath string)` to consume `GetRepoConfig` via `res.UnwrapOr("")`.
   - Update `GetRepoCurrentBranch(repoPath string)` to return `result.StringResult`. Guard on `res.IsFailure()`, returning `result.Fail[string](res.AppError())`.
   - Update `IsRepoClean(repoPath string)` to return `result.BoolResult`. Guard on `res.IsFailure()`, returning `result.Fail[bool](res.AppError())`.

2. **Refactor Metadata Utilities in `cli/gitutil/gitutil.go`:**
   - Update `RemoteURL(repoPath string)` to return `result.StringResult`, delegating to `GetRepoConfig(repoPath, constants.GitRemoteOrigin)`.
   - Update `RepoRoot(startPath string)` to return `result.StringResult`, invoking `RunGitInDir(startPath, constants.GitRevParse, "--show-toplevel")`. Guard on `res.IsFailure()`, returning `result.Fail[string](res.AppError())`.
   - Update `CurrentBranch(repoPath string)` to return `result.StringResult`, invoking `RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)`. Guard on `res.IsFailure()`, returning `result.Fail[string](res.AppError())`.
   - Update internal helper `runGit(dir string, args ...string)` to adapt `RunGitInDir` results with `if res.IsFailure() { return "", res.AppError() }`.

3. **Verify Coding Invariants:**
   - Ensure every refactored function stays $\le 15$ lines.
   - Maintain zero nested `if` statements using early returns.
   - Enforce vertical whitespace rules: blank line before `return` and after `}`.
   - Strictly positive boolean names (`is...`, `has...`).
   - Zero compiler builds or unit tests during authoring.

---

## Implementation Details

### 1. `cli/gitutil/git_cmd_exec.go`
```go
package gitutil

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunGitInDir executes a git command with directory arguments and returns stdout output.
func RunGitInDir(repoPath string, args ...string) result.ByteSliceResult {
	cmd := createGitCmd(repoPath, args...)
	out, err := cmd.Output()
	if err != nil {
		return result.Fail[[]byte](apperror.WrapSimple(err, "RunGitInDir"))
	}

	return result.Ok(out)
}

// RunGitCombinedInDir executes a git command with directory arguments and returns combined stdout and stderr.
func RunGitCombinedInDir(repoPath string, args ...string) result.ByteSliceResult {
	cmd := createGitCmd(repoPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return result.Fail[[]byte](apperror.WrapSimple(err, "RunGitCombinedInDir"))
	}

	return result.Ok(out)
}

// GetRepoConfig retrieves a git configuration value for key in repoPath.
func GetRepoConfig(repoPath, key string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitConfigCmd, constants.GitGetFlag, key)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// GetRemoteOriginURL retrieves the origin remote URL or an empty string on error.
func GetRemoteOriginURL(repoPath string) string {
	res := GetRepoConfig(repoPath, constants.GitRemoteOrigin)

	return res.UnwrapOr("")
}

// ReadRemoteOriginURL retrieves the configured remote.origin.url for repoPath.
func ReadRemoteOriginURL(repoPath string) result.StringResult {
	return GetRepoConfig(repoPath, constants.GitRemoteOrigin)
}

// GetRepoCurrentBranch returns the active branch name of the repository at repoPath.
func GetRepoCurrentBranch(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// IsRepoClean reports whether the working tree at repoPath has no uncommitted changes.
func IsRepoClean(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, "status", "--porcelain")
	if res.IsFailure() {
		return result.Fail[bool](res.AppError())
	}

	isClean := len(strings.TrimSpace(string(res.Data))) == 0

	return result.Ok(isClean)
}
```

### 2. `cli/gitutil/gitutil.go`
```go
// RemoteURL returns the origin remote URL for a repo at the given path.
func RemoteURL(repoPath string) result.StringResult {
	return GetRepoConfig(repoPath, constants.GitRemoteOrigin)
}

// RepoRoot returns the absolute path of the working-tree root for the
// repo containing `startPath`.
func RepoRoot(startPath string) result.StringResult {
	res := RunGitInDir(startPath, constants.GitRevParse, "--show-toplevel")
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// CurrentBranch returns the current branch name for a repo.
func CurrentBranch(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// runGit executes a git command in the given directory and returns stdout.
func runGit(dir string, args ...string) (string, error) {
	res := RunGitInDir(dir, args...)
	if res.IsFailure() {
		return "", res.AppError()
	}

	return string(res.Data), nil
}
```

---

## Verification & Conformance Checklist

- [ ] `RunGitInDir` and `RunGitCombinedInDir` return `result.ByteSliceResult`.
- [ ] Process failures wrapped in `*apperror.AppError` with clear operation context.
- [ ] `GetRepoConfig`, `ReadRemoteOriginURL`, `GetRepoCurrentBranch`, `RemoteURL`, `RepoRoot`, `CurrentBranch` return `result.StringResult`.
- [ ] `IsRepoClean` returns `result.BoolResult`.
- [ ] `GetRemoteOriginURL` safely consumes `res.UnwrapOr("")`.
- [ ] Internal `runGit` correctly adapts `RunGitInDir` results without leaking bare errors.
- [ ] All functions $\le 15$ lines with zero nested `if` statements.
- [ ] Zero compiler builds and zero unit tests executed during plan creation.
- [ ] Strictly relative Git paths utilized throughout all references.
