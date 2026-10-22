# Subtask Plan 08: Add `GetRepoRemoteURL` Helper with Direct Error Return (`cli/gitutil/git_cmd_exec.go`)

- **Subtask ID:** `subtask-08`
- **Subtask Name:** `08-gitutil-get-repo-remote-url`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/gitutil/git_repo_query.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/gitutil/git_cmd_exec.go`, remote URL queries are rigidly tied to `origin`:

1. `GetRemoteOriginURL(repoPath string) string` silently swallows errors and returns `""`.
2. `ReadRemoteOriginURL(repoPath string) (string, error)` queries `constants.GitRemoteOrigin` ("remote.origin.url") directly, but cannot query alternative remotes (such as `upstream`, `backup`, or custom remotes used in forks).
3. To support multi-remote workflows and adhere to **Rule 1 (Error Return Sovereignty)**:
   - A parameterized helper `GetRepoRemoteURL(repoPath, remoteName string) (string, error)` must be introduced.
   - It must return the remote URL or propagate git errors directly without swallowing.
   - `ReadRemoteOriginURL` can delegate to `GetRepoRemoteURL(repoPath, constants.GitOrigin)` for code DRYness.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Implement `GetRepoRemoteURL`
- **File:** `cli/gitutil/git_cmd_exec.go`
- **Signature:** `func GetRepoRemoteURL(repoPath, remoteName string) (string, error)`
- **Requirements:**
  - Construct configuration key `fmt.Sprintf("remote.%s.url", remoteName)`.
  - Delegate to `GetRepoConfig(repoPath, key)`.
  - Directly propagate `(string, error)` return.
  - Function length <= 6 lines.

```go
// GetRepoRemoteURL retrieves the configured URL for the specified remoteName at repoPath.
func GetRepoRemoteURL(repoPath, remoteName string) (string, error) {
	key := fmt.Sprintf("remote.%s.url", remoteName)

	return GetRepoConfig(repoPath, key)
}
```

### 2.2 Refactor `ReadRemoteOriginURL`
- **File:** `cli/gitutil/git_cmd_exec.go`
- **Signature:** `func ReadRemoteOriginURL(repoPath string) (string, error)`
- **Requirements:**
  - Delegate to `GetRepoRemoteURL(repoPath, constants.GitOrigin)`.
  - Function length <= 4 lines.

```go
// ReadRemoteOriginURL retrieves the configured remote.origin.url for repoPath.
func ReadRemoteOriginURL(repoPath string) (string, error) {
	return GetRepoRemoteURL(repoPath, constants.GitOrigin)
}
```

### 2.3 Sizing & Packaging Invariant
- **File:** `cli/gitutil/git_cmd_exec.go`
- **Rule:** Functions <= 15 lines, file size capped at <= 100 lines.
- If addition of `GetRepoRemoteURL` pushes `cli/gitutil/git_cmd_exec.go` over 100 lines, extract query helpers into `cli/gitutil/git_cmd_query.go`.

---

## 3. Implementation Verification Checklist

- [x] `GetRepoRemoteURL(repoPath, remoteName string) (string, error)` implemented in `cli/gitutil/`.
- [x] Returns error directly without swallow.
- [x] `GetRepoTrackingRemote` and `HasUnpushedCommits` implemented with error sovereignty.
- [x] Function lengths <= 15 lines.
- [x] File size limit <= 100 lines maintained.
- [x] Blank line before `return` and after `}` preserved.
