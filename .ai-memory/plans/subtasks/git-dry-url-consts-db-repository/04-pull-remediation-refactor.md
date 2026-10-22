# Subtask 04: Pull Remediation Subsystem Refactor (Eliminate Inline Git Execution, Hardcoded URLs, and Raw SQL)

- **Task ID:** subtask-git-dry-04
- **Parent Task:** git-dry-url-consts-db-repository
- **Spec Reference:** 02-spec/21-app/git-dry-url-consts-db-repository/02-component-spec.md
- **Status:** READY FOR EXECUTION
- **Owner:** Worker 04
- **Stage List:**
  - cli/cmdpull/pull_remediation_ssh.go
  - cli/cmdpull/pull_remediation_rebase.go
  - cli/cmdpull/pull_missing_remediate.go
  - cli/cmdpull/pull_prompt_step.go
  - cli/cmdpull/pull_prompt_conflict.go
  - cli/cmdpull/pull_prompt_exec.go
  - cli/cmdpull/pull_remediation_ssh_test.go

---

## 1. Context & Architectural Deficiencies

The `cli/cmdpull` package handles interactive pull workflows, automated rebase fixes, and SSH transport remediation. Over successive iterations, multiple anti-patterns accumulated across pull remediation files:

1. **Inline Subprocess Spawning (`git -C`):** Six separate functions in `cli/cmdpull/` bypass `cli/gitutil` and directly execute `exec.Command("git", "-C", repoPath, ...)` or `exec.Command("git", ...)`:
   - `pull_remediation_ssh.go:15` (`readRemoteOriginURL` running `git -C repoPath config --get remote.origin.url`)
   - `pull_remediation_ssh.go:99-100` (`git config --global credential...`)
   - `pull_remediation_rebase.go:50` (`git -C repoPath pull --progress --rebase --autostash`)
   - `pull_missing_remediate.go:116` (`git clone --progress remoteURL destination`)
   - `pull_prompt_step.go:57` (`git -C repoPath status --short`)
   - `pull_prompt_conflict.go:189` & `pull_prompt_exec.go:123` (`git -C item.RepoPath stash`)
2. **Hardcoded URLs & Organizations:** `pull_remediation_ssh.go` defines magic strings `"https://github.com/"`, `"git@github.com:"`, and `"alimtvnetwork"`.
3. **Repository Pattern Breach (Raw SQL in Business Logic):** `pull_remediation_ssh.go:121` directly runs `db.Conn().Exec("UPDATE Repo SET SshUrl = ?, DiscoveredUrl = ?, IdentifiedTransport = ? WHERE AbsolutePath = ?")`.

This subtask provides the concrete implementation plan to refactor `cli/cmdpull/pull_remediation_ssh.go` and companion remediation files, achieving zero inline git executions, zero raw URLs, and zero raw SQL in business logic.

---

## 2. Refactoring Target Matrix

| Source File | Function / Location | Current Anti-Pattern | Target Architecture / Replacement |
|---|---|---|---|
| `cli/cmdpull/pull_remediation_ssh.go` | `readRemoteOriginURL` (lines 14–21) | Private function with raw `exec.Command("git", "-C", ...)` | Delete function; replace calls with `gitutil.GetRemoteOriginURL(repoPath)` |
| `cli/cmdpull/pull_remediation_ssh.go` | `isSameOwnerHTTPSRepo` (lines 23–35) | Hardcoded `"https://github.com/"` | Use `constants.GitHubHTTPSPrefix` |
| `cli/cmdpull/pull_remediation_ssh.go` | `extractGitHubOwnerFromURL` (lines 37–49) | Hardcoded `"https://github.com/"` and `"git@github.com:"` | Use `constants.GitHubHTTPSPrefix` and `constants.GitHubSSHPrefix` |
| `cli/cmdpull/pull_remediation_ssh.go` | `hasSSHRepoForSameOwner` (lines 60–71) | Hardcoded `"alimtvnetwork"` | Use `constants.DefaultGitHubOwner` |
| `cli/cmdpull/pull_remediation_ssh.go` | `hasSSHRepoOwner` (lines 85–95) | Hardcoded `"git@github.com:"` | Use `string(constants.GitHubSSHPrefix) + strings.ToLower(owner) + constants.GitPathSeparator` |
| `cli/cmdpull/pull_remediation_ssh.go` | `remediateCredentialAndConvertToSSH` (lines 97–102) | Inline `exec.Command("git", "config", "--global", ...)` | Use `gitutil.SetGlobalGitConfig` |
| `cli/cmdpull/pull_remediation_ssh.go` | `updateRepoTransportInDB` (lines 114–125) | Raw `db.Conn().Exec(...)` in business logic | Call domain repository method `db.UpdateRepoTransport(repoPath, newURL, constants.TransportSSH)` |
| `cli/cmdpull/pull_remediation_rebase.go` | `execPullRebase` (lines 49–54) | Inline `exec.Command("git", "-C", repoPath, "pull", ...)` | Call `gitutil.ExecGitInDirCombined(repoPath, "pull", "--progress", "--rebase", "--autostash")` |
| `cli/cmdpull/pull_missing_remediate.go` | `ExecCloneRepo` (lines 109–120) | Inline `exec.Command("git", "clone", ...)` | Call `gitutil.CloneRepo(remoteURL, destination)` |
| `cli/cmdpull/pull_prompt_step.go` | `renderRepoShortDiff` (lines 56–70) | Inline `exec.Command("git", "-C", repoPath, "status", "--short")` | Call `gitutil.GetStatusShort(repoPath)` |
| `cli/cmdpull/pull_prompt_conflict.go` | `runStash` (lines 188–198) | Inline `exec.Command("git", "-C", item.RepoPath, "stash")` | Call `gitutil.StashWorkingTree(item.RepoPath)` |
| `cli/cmdpull/pull_prompt_exec.go` | `stashSingleRepoItem` (lines 121–130) | Inline `exec.Command("git", "-C", item.RepoPath, "stash")` | Call `gitutil.StashWorkingTree(item.RepoPath)` |

---

## 3. Step-by-Step Refactoring Instructions

### Phase 1: Refactor `cli/cmdpull/pull_remediation_ssh.go`

1. **Delete `readRemoteOriginURL`:**
   - Remove lines 14–21 completely.
   - Any calls in `isSameOwnerHTTPSRepo` now call `gitutil.GetRemoteOriginURL(repoPath)`.
2. **Refactor `isSameOwnerHTTPSRepo`:**
   ```go
   func isSameOwnerHTTPSRepo(repoPath string) bool {
       origin, err := gitutil.GetRemoteOriginURL(repoPath)
       if err != nil || !strings.HasPrefix(strings.ToLower(origin), string(constants.GitHubHTTPSPrefix)) {
           return false
       }
   
       owner := extractGitHubOwnerFromURL(origin)
       if owner == "" {
           return false
       }
   
       return hasSSHRepoForSameOwner(owner)
   }
   ```
3. **Refactor `extractGitHubOwnerFromURL`:**
   ```go
   func extractGitHubOwnerFromURL(url string) string {
       trimmed := strings.TrimSpace(url)
       lower := strings.ToLower(trimmed)
       if strings.HasPrefix(lower, string(constants.GitHubHTTPSPrefix)) {
           return extractFirstSlashSegment(trimmed[len(constants.GitHubHTTPSPrefix):])
       }
   
       if strings.HasPrefix(lower, string(constants.GitHubSSHPrefix)) {
           return extractFirstSlashSegment(trimmed[len(constants.GitHubSSHPrefix):])
       }
   
       return ""
   }
   ```
4. **Refactor `hasSSHRepoForSameOwner` & `hasSSHRepoOwner`:**
   ```go
   func hasSSHRepoForSameOwner(owner string) bool {
       if strings.EqualFold(owner, constants.DefaultGitHubOwner) {
           return true
       }
   
       repos := loadStoreRepos()
       if len(repos) == 0 {
           return false
       }
   
       return hasSSHRepoOwner(repos, owner)
   }
   
   func hasSSHRepoOwner(repos []model.ScanRecord, owner string) bool {
       prefix := string(constants.GitHubSSHPrefix) + strings.ToLower(owner) + constants.GitPathSeparator
       for _, r := range repos {
           if strings.HasPrefix(strings.ToLower(r.SSHUrl), prefix) ||
               strings.HasPrefix(strings.ToLower(r.DiscoveredURL), prefix) {
               return true
           }
       }
   
       return false
   }
   ```
5. **Refactor `remediateCredentialAndConvertToSSH`:**
   - Eliminate direct `exec.Command` calls for `credential.helper` and `credential.credentialStore`.
   - Route through `gitutil.SetGlobalGitConfig("credential.helper", "manager")` and `gitutil.SetGlobalGitConfig("credential.credentialStore", string(storeType))`.
6. **Refactor `updateRepoTransportInDB`:**
   - Remove raw `UPDATE Repo SET ...` SQL.
   - Consume `*store.DB.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)`:
   ```go
   func updateRepoTransportInDB(repoPath, sshURL string) {
       db, err := store.OpenDefault()
       if err != nil {
           return
       }
   
       defer db.Close()
       _ = db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)
   }
   ```
7. **Clean Imports:**
   - Remove `"os/exec"` from `cli/cmdpull/pull_remediation_ssh.go` entirely! No process execution remains in this file.

---

### Phase 2: Refactor Companion Remediation Files

1. **`cli/cmdpull/pull_remediation_rebase.go`**:
   - Refactor `execPullRebase(repoPath string) ([]byte, error)`:
     ```go
     func execPullRebase(repoPath string) ([]byte, error) {
         return gitutil.ExecGitInDirCombined(repoPath, "pull", "--progress", "--rebase", "--autostash")
     }
     ```
   - Remove `"os/exec"` import if no other exec calls exist.
2. **`cli/cmdpull/pull_missing_remediate.go`**:
   - Refactor `ExecCloneRepo(remoteURL, destination string) error`:
     ```go
     func ExecCloneRepo(remoteURL, destination string) error {
         parentDir := filepath.Dir(destination)
         if err := os.MkdirAll(parentDir, 0755); err != nil {
             return err
         }
     
         return gitutil.CloneRepo(remoteURL, destination)
     }
     ```
   - Remove `"os/exec"` import.
3. **`cli/cmdpull/pull_prompt_step.go`**:
   - Refactor `renderRepoShortDiff(repoPath string)`:
     - Replace raw status exec with `statusOut, err := gitutil.GetStatusShort(repoPath)`.
4. **`cli/cmdpull/pull_prompt_conflict.go` & `cli/cmdpull/pull_prompt_exec.go`**:
   - Refactor `runStash` and `stashSingleRepoItem`:
     - Replace `exec.Command("git", "-C", repoPath, "stash")` with `gitutil.StashWorkingTree(repoPath)`.

---

### Phase 3: Unit Testing in `cli/cmdpull/pull_remediation_ssh_test.go`

Author comprehensive test coverage in `cli/cmdpull/pull_remediation_ssh_test.go`:
- Test `extractGitHubOwnerFromURL`:
  - `https://github.com/alimtvnetwork/gitmap-v28` -> `alimtvnetwork`
  - `git@github.com:alimtvnetwork/gitmap-v28.git` -> `alimtvnetwork`
  - `https://gitlab.com/owner/repo` -> `""` (ignores non-GitHub URLs)
  - `invalid-url` -> `""`
- Test `hasSSHRepoForSameOwner`:
  - Returns `true` for `alimtvnetwork` without querying DB.
- Test `hasSSHRepoOwner`:
  - Verifies prefix matching with lowercased owner and SSH prefix.

---

## 4. Coding Guidelines & Invariant Rules

- **Strict Error Return Sovereignty:** In all leaf functions, bubble errors directly.
- **Positive Booleans Only:** Ensure `isSameOwnerHTTPSRepo`, `hasSSHRepoForSameOwner`, `hasSSHRepoOwner` follow positive prefixes.
- **15-Line Function Limit:** Every function must be <= 15 lines (ideally <= 8 lines).
- **Blank Line Spacing:** A blank line must precede every `return` statement and follow every closing brace `}`.
- **Zero Package Violations:** Zero SQL queries in `cmdpull`; zero `exec.Command` for git operations in `cmdpull`.

---

## 5. Verification Checklist & Evidence Requirements

1. Verify that `cli/cmdpull/pull_remediation_ssh.go` has:
   - Zero imports of `"os/exec"`.
   - Zero occurrences of `db.Conn().Exec`.
   - Zero occurrences of `"https://github.com/"` or `"git@github.com:"` string literals.
2. Verify that `cli/cmdpull/pull_remediation_rebase.go`, `pull_missing_remediate.go`, `pull_prompt_step.go`, `pull_prompt_conflict.go`, and `pull_prompt_exec.go` contain zero direct `git -C` subprocess executions.
3. Verify all file lengths are <= 300 lines with zero nested if statements.
