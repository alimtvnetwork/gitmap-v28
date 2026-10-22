# Subtask 01: Centralized Git Execution Primitives & Internal GitUtil DRY Refactoring

> **Subtask ID:** `subtask-01`  
> **Parent Architecture Spec:** `02-spec/21-app/git-dry-url-consts-db-repository/01-architecture-spec.md`  
> **Target Path:** `.ai-memory/plans/subtasks/git-dry-url-consts-db-repository/01-git-exec-dry-helpers.md`  
> **Status:** `PENDING`  
> **Owned Files:**  
> - `cli/gitutil/git_cmd_exec.go` (new)  
> - `cli/gitutil/git_cmd_exec_test.go` (new)  
> - `cli/gitutil/gitutil.go`  
> - `cli/gitutil/active_branch.go`  
> - `cli/gitutil/case_collision.go`  
> - `cli/gitutil/checkout.go`  
> - `cli/gitutil/dirty_inspect.go`  
> - `cli/gitutil/last_sha.go`  
> - `cli/gitutil/latest_remote.go`  
> - `cli/gitutil/latest_tag.go`  
> - `cli/gitutil/latestbranch.go`  
> - `cli/gitutil/latestbranchresolve.go`  
> - `cli/gitutil/pr_detector.go`  
> - `cli/gitutil/remediation_case.go`  
> - `cli/gitutil/upstream_delta.go`  

---

## 1. Objectives & Scope

1. **Implement Centralized Execution Helpers in `cli/gitutil/git_cmd_exec.go`:**
   - `RunGitInDir(repoPath string, args ...string) (string, error)`: Standard trimmed stdout runner.
   - `RunGitCombinedInDir(repoPath string, args ...string) ([]byte, error)`: Combined stdout/stderr runner for porcelain and status parsing.
   - `GetRepoConfig(repoPath string, key string) (string, error)`: Typed wrapper for repository configuration lookup.
   - `ReadRemoteOriginURL(repoPath string) (string, error)`: Canonical reader for `remote.origin.url`.
   - `SetRemoteOriginURL(repoPath string, remoteURL string) error`: Canonical writer for `remote.origin.url`.
2. **Mandatory Subprocess Sanitization:**
   - Inject `BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)` in all helpers to eliminate headless stalls and interactive password prompts.
   - Enforce process directory confinement via `cmd.Dir = repoPath` instead of fragile `"-C"` flag manipulation.
3. **Internal `cli/gitutil/` DRY Refactoring:**
   - Migrate 13 internal utility files in `cli/gitutil/` away from ad-hoc `exec.Command` and `exec.Command("git", "-C", ...)` calls to the centralized helpers.
   - Replace the legacy unexported `runGit` in `cli/gitutil/gitutil.go` with delegation to `RunGitInDir`.
4. **Comprehensive Unit Testing:**
   - Author `cli/gitutil/git_cmd_exec_test.go` to verify environment injection, error wrapping, directory confinement, and configuration manipulation.

---

## 2. Technical Implementation Specifications

### 2.1 New Helper Primitives (`cli/gitutil/git_cmd_exec.go`)

Each helper MUST adhere strictly to the **15-line function limit**, **zero nested ifs**, and **positive boolean naming rules**.

#### 2.1.1 `RunGitInDir`
- **Signature:** `func RunGitInDir(repoPath string, args ...string) (string, error)`
- **Behavior:**
  - Instantiates `cmd := exec.Command(constants.GitBin, args...)`.
  - Sets `cmd.Dir = repoPath` if `repoPath != ""`.
  - Sets `cmd.Env = BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)`.
  - Invokes `out, err := cmd.Output()`.
  - If `err != nil`, returns `"", err`.
  - Returns `strings.TrimSpace(string(out)), nil`.
- **Line Count Budget:** 12 lines.

#### 2.1.2 `RunGitCombinedInDir`
- **Signature:** `func RunGitCombinedInDir(repoPath string, args ...string) ([]byte, error)`
- **Behavior:**
  - Instantiates `cmd := exec.Command(constants.GitBin, args...)`.
  - Sets `cmd.Dir = repoPath` if `repoPath != ""`.
  - Sets `cmd.Env = BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)`.
  - Returns `cmd.CombinedOutput()`.
- **Line Count Budget:** 10 lines.

#### 2.1.3 `GetRepoConfig`
- **Signature:** `func GetRepoConfig(repoPath string, key string) (string, error)`
- **Behavior:**
  - Calls `RunGitInDir(repoPath, constants.GitConfigCmd, constants.GitGetFlag, key)`.
  - Returns trimmed configuration value or error.
- **Line Count Budget:** 4 lines.

#### 2.1.4 `ReadRemoteOriginURL`
- **Signature:** `func ReadRemoteOriginURL(repoPath string) (string, error)`
- **Behavior:**
  - Calls `GetRepoConfig(repoPath, constants.GitRemoteOrigin)`.
  - Returns the URL string or error.
- **Line Count Budget:** 4 lines.

#### 2.1.5 `SetRemoteOriginURL`
- **Signature:** `func SetRemoteOriginURL(repoPath string, remoteURL string) error`
- **Behavior:**
  - Calls `_, err := RunGitInDir(repoPath, constants.GitConfigCmd, constants.GitRemoteOrigin, remoteURL)`.
  - Returns error.
- **Line Count Budget:** 5 lines.

---

### 2.2 Refactoring Inventory for `cli/gitutil/`

The following files currently instantiate ad-hoc `exec.Command` instances and MUST be updated:

| File | Current Pattern | Target Pattern |
| :--- | :--- | :--- |
| `cli/gitutil/active_branch.go` | `exec.Command("git", "-C", cleanPath, "rev-parse", ...)` | `RunGitInDir(cleanPath, constants.GitRevParse, ...)` |
| `cli/gitutil/case_collision.go` | `exec.Command("git", "-C", repoPath, "ls-files", "-s")` | `RunGitInDir(repoPath, "ls-files", "-s")` |
| `cli/gitutil/checkout.go` | `exec.Command(constants.GitBin, constants.GitCheckout, branch)` | `RunGitInDir("", constants.GitCheckout, branch)` |
| `cli/gitutil/dirty_inspect.go` | `exec.Command("git", "-C", repoPath, "status", "--porcelain")` | `RunGitInDir(repoPath, "status", "--porcelain")` |
| `cli/gitutil/last_sha.go` | `exec.Command("git", "-C", cleanPath, "rev-parse", ...)` | `RunGitInDir(cleanPath, constants.GitRevParse, ...)` |
| `cli/gitutil/latest_remote.go` | `exec.Command("git", "-C", repoPath, "for-each-ref", ...)` | `RunGitInDir(repoPath, "for-each-ref", ...)` |
| `cli/gitutil/latest_tag.go` | `exec.Command("git", "-C", repoPath, "describe", ...)` | `RunGitInDir(repoPath, "describe", ...)` |
| `cli/gitutil/latestbranch.go` | `exec.Command(constants.GitBin, constants.GitRevParse, ...)` | `RunGitInDir("", constants.GitRevParse, ...)` |
| `cli/gitutil/latestbranchresolve.go` | `exec.Command(constants.GitBin, constants.GitLog, ...)` | `RunGitInDir("", constants.GitLog, ...)` |
| `cli/gitutil/pr_detector.go` | `exec.Command("git", "-C", repoPath, "for-each-ref", ...)` | `RunGitInDir(repoPath, "for-each-ref", ...)` |
| `cli/gitutil/remediation_case.go` | `exec.Command("git", "-C", repoPath, "rm" / "add" / "commit")` | `RunGitInDir(repoPath, ...)` |
| `cli/gitutil/upstream_delta.go` | `exec.Command("git", "-C", repoPath, "rev-list", ...)` | `RunGitInDir(repoPath, "rev-list", ...)` |
| `cli/gitutil/gitutil.go` | unexported `runGit` with raw `exec.Command` | Delegate `runGit` to `RunGitInDir(dir, args...)` |

---

## 3. Coding Guidelines & Invariant Guardrails

1. **Positive Boolean Naming:**
   - All boolean variables and return checks MUST begin with `is` or `has` (`isSafeEnvInjected`, `hasDirectoryContext`, `isSuccessful`).
   - Strictly no negative booleans (`isNotDirty`, `unverified`).
2. **Function Line Budget (15 Lines Max):**
   - No function in `git_cmd_exec.go` or refactored files may exceed 15 lines. Helper decomposition is mandatory if branching logic expands.
3. **Zero Nested Conditionals:**
   - Use early returns and guard clauses exclusively.
4. **Vertical Whitespace Discipline:**
   - Blank line after closing curly brace `}`.
   - Blank line before `return` statement.
5. **Path Hygiene Invariant:**
   - Strictly relative repository paths in all comments and documentation.
   - Zero absolute paths or `file:///` URIs.

---

## 4. Actionable Task Checklist

- [ ] **Step 1: Author Primitives (`cli/gitutil/git_cmd_exec.go`)**
  - Implement `RunGitInDir`.
  - Implement `RunGitCombinedInDir`.
  - Implement `GetRepoConfig`.
  - Implement `ReadRemoteOriginURL`.
  - Implement `SetRemoteOriginURL`.
  - Verify compliance with line count budget (<= 15 lines per function) and positive booleans.

- [ ] **Step 2: Author Test Suite (`cli/gitutil/git_cmd_exec_test.go`)**
  - Test `RunGitInDir` with mock/valid parameters and error states.
  - Test `RunGitCombinedInDir` with output capture.
  - Test `GetRepoConfig` and `ReadRemoteOriginURL` with valid key and non-existent key.
  - Test `SetRemoteOriginURL` with mock updates.

- [ ] **Step 3: Refactor Core `cli/gitutil/gitutil.go`**
  - Replace unexported `runGit` body with `return RunGitInDir(dir, args...)`.
  - Update `RemoteURL` to delegate to `ReadRemoteOriginURL`.

- [ ] **Step 4: Refactor Inspection & Resolution Utilities**
  - Refactor `active_branch.go` to use `RunGitInDir`.
  - Refactor `last_sha.go` to use `RunGitInDir`.
  - Refactor `dirty_inspect.go` to use `RunGitInDir`.
  - Refactor `case_collision.go` to use `RunGitInDir`.
  - Refactor `upstream_delta.go` to use `RunGitInDir`.
  - Refactor `pr_detector.go` to use `RunGitInDir`.

- [ ] **Step 5: Refactor Remediation & Branch Utilities**
  - Refactor `latest_remote.go` to use `RunGitInDir`.
  - Refactor `latest_tag.go` to use `RunGitInDir`.
  - Refactor `checkout.go` to use `RunGitInDir`.
  - Refactor `latestbranch.go` to use `RunGitInDir`.
  - Refactor `latestbranchresolve.go` to use `RunGitInDir`.
  - Refactor `remediation_case.go` to use `RunGitInDir`.

- [ ] **Step 6: Final Cleanliness & Verification Review**
  - Verify zero remaining raw `exec.Command` calls in `cli/gitutil/` outside of `git_cmd_exec.go` and `git_exec_timeout.go`.
  - Confirm all file references use relative paths only.
  - Format all modified files with standard Go formatting rules.
