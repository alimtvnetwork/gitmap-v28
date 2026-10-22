# Subtask Plan 08: Ensure Zero Nested Ifs in `buildGitDirArgs` via Early Guard Returns

- **Subtask ID:** `subtask-08`
- **Subtask Name:** `08-gitutil-zero-nested-ifs-buildgitdirargs`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/gitutil/git_cmd_exec.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/gitutil/git_cmd_exec.go`, `buildGitDirArgs` constructs the argument slice for `exec.Cmd` execution when invoking Git across discrete repository paths:

1. Under **Rule 3 (Zero Nested Ifs)**, all conditional paths must be flattened using early guard returns.
2. Under **Rule 4 (Positive Boolean Naming)**, condition checks must evaluate positive booleans (`hasRepoPath`) rather than negative polarity.
3. Under **Rule 2 (Function Line Limits)**, functions must strictly remain <= 15 lines (target <= 8 lines).
4. Under return styling rules, every closing brace `}` must be followed by a blank line, and a blank line must precede `return`.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `buildGitDirArgs`
- **File:** `cli/gitutil/git_cmd_exec.go`
- **Signature:** `func buildGitDirArgs(repoPath string, args ...string) []string`
- **Requirements:**
  - Evaluate positive boolean: `hasRepoPath := len(strings.TrimSpace(repoPath)) > 0`.
  - Guard return: `if !hasRepoPath { return args }`.
  - Allocate combined slice with pre-sized capacity: `full := make([]string, 0, len(args)+2)`.
  - Append `constants.GitDirFlag, repoPath`, followed by `args...`.
  - Return `full`.
  - Total length: 11 lines.

```go
// buildGitDirArgs constructs arguments prepended with the git directory flag when repoPath is provided.
func buildGitDirArgs(repoPath string, args ...string) []string {
	hasRepoPath := len(strings.TrimSpace(repoPath)) > 0
	if !hasRepoPath {
		return args
	}

	full := make([]string, 0, len(args)+2)
	full = append(full, constants.GitDirFlag, repoPath)
	full = append(full, args...)

	return full
}
```

### 2.2 Verify Accompanying Git Execution Helpers in `cli/gitutil/git_cmd_exec.go`
- **`createGitCmd(repoPath string, args ...string) *exec.Cmd`:**
  - Builds arguments via `buildGitDirArgs`.
  - Applies safe environment via `BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)`.
  - Length: 8 lines.
- **`GetRepoConfig(repoPath, key string) (string, error)`:**
  - Checks error from `RunGitInDir`; returns `("", err)` immediately if non-nil.
  - Returns `(strings.TrimSpace(string(out)), nil)`.
  - Length: 10 lines.
- **`SetRemoteOriginURL(repoPath, url string) error`:**
  - Invokes `RunGitCombinedInDir(repoPath, "remote", "set-url", constants.GitOrigin, url)`.
  - Returns error directly; length: 6 lines.

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Error Sovereignty:** Subprocess execution errors propagated directly.
- [ ] **15-Line Limit:** Every routine is strictly between 6 and 11 lines.
- [ ] **Zero Nested Ifs:** All branches eliminated via early guard returns.
- [ ] **Positive Booleans:** Uses `hasRepoPath` positive flag.
- [ ] **Return Styling:** Blank line before `return` and after `}`.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.
