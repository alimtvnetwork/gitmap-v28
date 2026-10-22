# Subtask Plan 07: Add `GetRepoHeadCommitSha` Helper with Direct Error Return (`cli/gitutil/git_cmd_exec.go`)

- **Subtask ID:** `subtask-07`
- **Subtask Name:** `07-gitutil-get-repo-head-commit-sha`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/gitutil/git_repo_query.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In GitMap, numerous commands (`cmdpipeline`, `cmdpull`, `cmdscan`) require retrieving the current HEAD commit hash of a repository at an arbitrary path.

Currently:
1. `cli/gitutil/git_cmd_exec.go` provides `GetRepoCurrentBranch(repoPath string) (string, error)` which calls `RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)`.
2. However, there is no canonical, dedicated helper for retrieving the active commit SHA (`git rev-parse HEAD`).
3. As a result, different modules construct ad-hoc git commands or parse output inconsistently, sometimes ignoring the error return.
4. Adding `GetRepoHeadCommitSha(repoPath string) (string, error)` establishes a single, reusable, error-sovereign utility function for querying the repository HEAD commit hash.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Implement `GetRepoHeadCommitSha`
- **File:** `cli/gitutil/git_cmd_exec.go`
- **Signature:** `func GetRepoHeadCommitSha(repoPath string) (string, error)`
- **Requirements:**
  - Execute `RunGitInDir(repoPath, constants.GitRevParse, constants.GitHEAD)`.
  - Check error; return `"", err` immediately if git command fails.
  - Return `strings.TrimSpace(string(out)), nil`.
  - Function length <= 10 lines.

```go
// GetRepoHeadCommitSha retrieves the full HEAD commit SHA of the repository at repoPath.
func GetRepoHeadCommitSha(repoPath string) (string, error) {
	out, err := RunGitInDir(repoPath, constants.GitRevParse, constants.GitHEAD)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
```

### 2.2 Sizing & Packaging Invariant
- **File:** `cli/gitutil/git_cmd_exec.go`
- **Rule:** Functions <= 15 lines, file size capped at <= 100 lines.
- If addition of `GetRepoHeadCommitSha` causes `cli/gitutil/git_cmd_exec.go` to exceed 100 lines, extract repository read query helpers (`GetRepoConfig`, `GetRepoCurrentBranch`, `GetRepoHeadCommitSha`, `IsRepoClean`) into `cli/gitutil/git_cmd_query.go` (preserving `package gitutil`).

---

## 3. Implementation Verification Checklist

- [x] `GetRepoHeadCommitSha(repoPath string) (string, error)` implemented in `cli/gitutil/`.
- [x] Returns error directly without swallow or exit call.
- [x] Uses `constants.GitRevParse` and `constants.GitHEAD`.
- [x] Function length <= 10 lines.
- [x] File size limit <= 100 lines maintained.
- [x] Blank line before `return` and after `}` preserved.
