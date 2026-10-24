# Subtask Plan 02: Autonomous Git & SSH Operation Interceptor with Longest-Prefix Resolution

- **Subtask Slug:** `subtask-02-git-ssh-interceptor`
- **Parent Task:** `317-system-for-folder-specific-work-directories-and-ssh-key-automation`
- **Spec Path:** `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/02-git-ssh-interceptor-spec.md`
- **Status:** Pending Execution
- **Target Files:**
  - `cli/constants/constants_git.go`
  - `cli/constants/constants_ssh.go`
  - `cli/model/folder_ssh.go`
  - `cli/gitutil/git_ssh_interceptor.go`
  - `cli/gitutil/git_ssh_cmd.go`
  - `cli/gitutil/git_env.go`
  - `cli/gitutil/git_cmd_exec.go`
  - `cli/store/folder_ssh_schema.go`
  - `cli/store/folder_ssh_repo.go`
  - `cli/cloner/safe_pull.go`
  - `cli/cloner/safe_push.go`
  - `cli/cmdpull/push_exec.go`
  - `cli/gitutil/git_ssh_interceptor_test.go`
  - `cli/gitutil/git_ssh_cmd_test.go`

---

## 1. Overview & Objectives

In multi-client and multi-organization setups, developers and AI swarms execute Git commands across various directories where each directory or workgroup is bound to a distinct SSH key, author email, and working directory anchor.

This subtask implements the **Autonomous Git & SSH Operation Interceptor**:
1. **Longest-Prefix Folder Resolution:** Given any repository path or CWD, find the most specific directory binding in SQLite using boundary-safe prefix matching.
2. **Deterministic SSH Command Synthesis:** Construct an isolated `GIT_SSH_COMMAND` containing `-F /dev/null -o IdentitiesOnly=yes -i <key> -o StrictHostKeyChecking=accept-new -o BatchMode=yes -o ConnectTimeout=15`.
3. **Subprocess Environment Injection:** Inject zero-prompt parameters (`GIT_TERMINAL_PROMPT=0`, `GIT_ASKPASS=""`, `SSH_ASKPASS=""`, `GCM_INTERACTIVE=never`, `GCM_NO_PERSIST=1`) into all Git invocations.
4. **Transparent Git Operation Interception:** Wire the interceptor into `git pull`, `git push`, `git fetch`, and fleet sync operations across both single-repo and batch parallel executions.
5. **Zero-Prompt Persistence:** Read bindings directly from SQLite (`gitmap.db`), guaranteeing that once registered, GitMap never prompts the user or agent again.

---

## 2. Step-by-Step Implementation Instructions

### Phase 1: Define Constants in `cli/constants/constants_ssh.go` & `cli/constants/constants_git.go`

1. Update `cli/constants/constants_ssh.go` to add OpenSSH flags and options:
   ```go
   // SSH isolation and security flags.
   const (
       SSHFlagConfigFile       = "-F"
       DevNullPath             = "/dev/null"
       SSHFlagIdentity         = "-i"
       SSHIdentitiesOnlyYes    = "IdentitiesOnly=yes"
       SSHBatchModeYes         = "BatchMode=yes"
       SSHConnectTimeout15     = "ConnectTimeout=15"
       TableFolderSSHBinding   = "folder_ssh_bindings"
   )
   ```

2. Update `cli/constants/constants_git.go` with interceptor environment keys:
   ```go
   const (
       EnvGitAuthorName     = "GIT_AUTHOR_NAME"
       EnvGitAuthorEmail    = "GIT_AUTHOR_EMAIL"
       EnvGitCommitterName  = "GIT_COMMITTER_NAME"
       EnvGitCommitterEmail = "GIT_COMMITTER_EMAIL"
   )
   ```

---

### Phase 2: Domain Models & Types in `cli/model/folder_ssh.go`

1. Create `cli/model/folder_ssh.go`:
   ```go
   package model

   import "time"

   // FolderSSHBinding models a persistent binding between a folder and an SSH key.
   type FolderSSHBinding struct {
       ID          int64     `json:"id" db:"id"`
       FolderPath  string    `json:"folderPath" db:"folder_path"`
       WorkgroupID int64     `json:"workgroupId,omitempty" db:"workgroup_id"`
       SSHKeyID    int64     `json:"sshKeyId" db:"ssh_key_id"`
       KeyPath     string    `json:"keyPath,omitempty" db:"-"`
       AuthorName  string    `json:"authorName,omitempty" db:"author_name"`
       AuthorEmail string    `json:"authorEmail,omitempty" db:"author_email"`
       Priority    int       `json:"priority" db:"priority"`
       IsActive    bool      `json:"isActive" db:"is_active"`
       CreatedAt   time.Time `json:"createdAt" db:"created_at"`
       UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
   }

   // ResolvedGitContext holds resolved credentials and environment for Git execution.
   type ResolvedGitContext struct {
       RepoPath       string   `json:"repoPath"`
       WorkdirAnchor  string   `json:"workdirAnchor"`
       BoundKeyPath   string   `json:"boundKeyPath"`
       AuthorName     string   `json:"authorName"`
       AuthorEmail    string   `json:"authorEmail"`
       GitSSHCommand  string   `json:"gitSSHCommand"`
       HasInjectedSSH bool     `json:"hasInjectedSSH"`
       Environment    []string `json:"-"`
   }
   ```

---

### Phase 3: Longest-Prefix Matcher in `cli/gitutil/git_ssh_interceptor.go`

1. Implement `FindLongestPrefixBinding` with boundary validation:
   ```go
   package gitutil

   import (
       "strings"

       "github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
       "github.com/alimtvnetwork/gitmap-v28/cli/model"
   )

   // FindLongestPrefixBinding identifies the most specific folder binding for targetPath.
   func FindLongestPrefixBinding(targetPath string, bindings []model.FolderSSHBinding) (model.FolderSSHBinding, bool) {
       targetKey := fsutil.CanonicalPathKey(fsutil.NormalizeToForwardSlashes(targetPath))
       var bestMatch model.FolderSSHBinding
       hasMatch := false
       maxLen := -1

       for _, b := range bindings {
           if !b.IsActive {
               continue
           }

           folderKey := fsutil.CanonicalPathKey(fsutil.NormalizeToForwardSlashes(b.FolderPath))
           isCandidate := isPathPrefixMatch(targetKey, folderKey)
           if !isCandidate {
               continue
           }

           matchLen := len(folderKey)
           if matchLen > maxLen || (matchLen == maxLen && b.Priority > bestMatch.Priority) {
               maxLen = matchLen
               bestMatch = b
               hasMatch = true
           }
       }

       return bestMatch, hasMatch
   }

   // isPathPrefixMatch checks if target matches folder directly or at directory boundary.
   func isPathPrefixMatch(targetKey, folderKey string) bool {
       if targetKey == folderKey {
           return true
       }

       if !strings.HasPrefix(targetKey, folderKey) {
           return false
       }

       return strings.HasPrefix(targetKey[len(folderKey):], "/")
   }
   ```

---

### Phase 4: Deterministic SSH Command Builder in `cli/gitutil/git_ssh_cmd.go`

1. Create `cli/gitutil/git_ssh_cmd.go`:
   ```go
   package gitutil

   import (
       "fmt"

       "github.com/alimtvnetwork/gitmap-v28/cli/constants"
       "github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
       "github.com/alimtvnetwork/gitmap-v28/cli/model"
   )

   // BuildDeterministicSSHCommand synthesizes an isolated OpenSSH invocation string.
   func BuildDeterministicSSHCommand(privateKeyPath string) string {
       cleanPath := fsutil.NormalizeToForwardSlashes(privateKeyPath)
       quotedKey := "\"" + cleanPath + "\""

       return fmt.Sprintf(
           "%s %s %s %s %s %s %s %s %s %s %s %s",
           constants.SSHBin,
           constants.SSHFlagConfigFile,
           constants.DevNullPath,
           constants.SSHOptionFlag,
           constants.SSHIdentitiesOnlyYes,
           constants.SSHFlagIdentity,
           quotedKey,
           constants.SSHOptionFlag,
           constants.SSHStrictHostKeyAcceptNew,
           constants.SSHOptionFlag,
           constants.SSHBatchModeYes,
           constants.SSHConnectTimeout15,
       )
   }

   // BuildSafeGitEnvForRepo constructs a subprocess environment with resolved SSH credentials.
   func BuildSafeGitEnvForRepo(ctx model.ResolvedGitContext, extra ...string) []string {
       hasSSH := ctx.HasInjectedSSH && len(ctx.GitSSHCommand) > 0
       if !hasSSH {
           return BuildSafeGitEnv(extra...)
       }

       sshVar := constants.EnvGitSSHCommand + constants.EnvAssignmentSeparator + ctx.GitSSHCommand
       combinedExtra := make([]string, 0, len(extra)+1)
       combinedExtra = append(combinedExtra, sshVar)
       combinedExtra = append(combinedExtra, extra...)

       return BuildSafeGitEnv(combinedExtra...)
   }
   ```

---

### Phase 5: Wire Git Execution Helpers in `cli/gitutil/git_cmd_exec.go`

1. Add `CreateGitCmdInDirWithContext` to `cli/gitutil/git_cmd_exec.go`:
   ```go
   // CreateGitCmdInDirWithContext creates an exec.Cmd configured with directory and resolved context.
   func CreateGitCmdInDirWithContext(repoPath string, ctx model.ResolvedGitContext, args ...string) *exec.Cmd {
       full := buildGitDirArgs(repoPath, args...)
       cmd := exec.Command(constants.GitBin, full...)
       cmd.Env = BuildSafeGitEnvForRepo(ctx)

       return cmd
   }
   ```

2. Add context resolution helper `ResolveRepoGitContext`:
   ```go
   // ResolveRepoGitContext queries SQLite bindings and constructs the ResolvedGitContext for repoPath.
   func ResolveRepoGitContext(repoPath string) model.ResolvedGitContext {
       ctx := model.ResolvedGitContext{RepoPath: repoPath}
       bindings, err := store.ListActiveFolderSSHBindings()
       if err != nil || len(bindings) == 0 {
           return ctx
       }

       binding, hasMatch := FindLongestPrefixBinding(repoPath, bindings)
       if !hasMatch || len(binding.KeyPath) == 0 {
           return ctx
       }

       ctx.WorkdirAnchor = binding.FolderPath
       ctx.BoundKeyPath = binding.KeyPath
       ctx.AuthorName = binding.AuthorName
       ctx.AuthorEmail = binding.AuthorEmail
       ctx.GitSSHCommand = BuildDeterministicSSHCommand(binding.KeyPath)
       ctx.HasInjectedSSH = true

       return ctx
   }
   ```

---

### Phase 6: Cloner & Pull/Push Integration

1. Update `cli/cloner/safe_pull.go` & `cli/cloner/safe_push.go`:
   - Replace static `buildSafePushEnv()` with `gitutil.BuildSafeGitEnvForRepo(gitutil.ResolveRepoGitContext(repoDir))`
   - Ensure parallel workers inherit repository-specific injected environments without race conditions.

2. Update `cli/cmdpull/push_exec.go`:
   - Inject resolved context into `runGitCapturingStderr` and `runGitInherit`.

---

### Phase 7: SQLite Schema & Store Queries

1. Create `cli/store/folder_ssh_schema.go`:
   ```go
   package store

   const (
       SQLCreateFolderSSHBindingsTable = `
       CREATE TABLE IF NOT EXISTS folder_ssh_bindings (
           id INTEGER PRIMARY KEY AUTOINCREMENT,
           folder_path TEXT UNIQUE NOT NULL,
           workgroup_id INTEGER DEFAULT 0,
           ssh_key_id INTEGER NOT NULL,
           author_name TEXT DEFAULT '',
           author_email TEXT DEFAULT '',
           priority INTEGER DEFAULT 100,
           is_active INTEGER DEFAULT 1,
           created_at TEXT DEFAULT CURRENT_TIMESTAMP,
           updated_at TEXT DEFAULT CURRENT_TIMESTAMP
       );`

       SQLSelectActiveFolderSSHBindings = `
       SELECT b.id, b.folder_path, b.workgroup_id, b.ssh_key_id, k.PrivatePath,
              b.author_name, b.author_email, b.priority, b.is_active, b.created_at, b.updated_at
       FROM folder_ssh_bindings b
       JOIN SshKey k ON b.ssh_key_id = k.SshKeyId
       WHERE b.is_active = 1
       ORDER BY length(b.folder_path) DESC, b.priority DESC;`
   )
   ```

2. Create `cli/store/folder_ssh_repo.go` implementing `ListActiveFolderSSHBindings`, `UpsertFolderSSHBinding`, and `DeleteFolderSSHBinding`.

---

## 3. Coding Guidelines Compliance Invariants

Before considering this subtask complete, verify all modifications strictly uphold the GitMap coding standards:

- [ ] **Positive Booleans Only:** All booleans prefixed with `is` or `has` (`hasBoundKey`, `hasInjectedSSH`, `isCandidate`, `isActive`). Banned: `can`, `should`, `was`, `will`, and negative booleans.
- [ ] **Function Line Limits:** Every function <= 15 lines max (<= 8 lines preferred). Split longer functions into discrete helper functions.
- [ ] **Zero Nested Conditionals:** Zero nested `if` statements. Guard clauses and early returns are used everywhere.
- [ ] **Return Spacing:** Empty line after `}` and empty line before `return`.
- [ ] **Leaf Error Sovereignty:** No leaf functions call `cliexit.HandleError` or `os.Exit`. All errors return `error` or `result.Result[T]`.

---

## 4. Acceptance Criteria & Test Verification Checklist

- [ ] **AC 1 (Longest-Prefix Accuracy):** Target `D:/work/clientA/repo1` selects folder `D:/work/clientA` even when `D:/work` also exists in database.
- [ ] **AC 2 (Boundary Safety):** Target `D:/work/clientA-extra/repo1` does NOT match `D:/work/clientA`.
- [ ] **AC 3 (Case-Insensitive Normalization):** On Windows, `d:\Work\ClientA\Repo` correctly matches `D:/work/clienta`.
- [ ] **AC 4 (Deterministic Command Synthesis):** Output of `BuildDeterministicSSHCommand` contains:
  - `-F /dev/null`
  - `-o IdentitiesOnly=yes`
  - `-i "<path>"`
  - `-o StrictHostKeyChecking=accept-new`
  - `-o BatchMode=yes`
  - `-o ConnectTimeout=15`
- [ ] **AC 5 (Subprocess Environment):** `BuildSafeGitEnvForRepo` outputs `GIT_TERMINAL_PROMPT=0`, `GIT_ASKPASS=""`, `SSH_ASKPASS=""`, `GCM_INTERACTIVE=never`, `GCM_NO_PERSIST=1`, and the synthesized `GIT_SSH_COMMAND`.
- [ ] **AC 6 (Zero-Prompt Persistence):** Running `git pull` or `git push` on a mapped folder never triggers interactive terminal prompts or questions.
- [ ] **AC 7 (Unit Test Coverage):**
  - `cli/gitutil/git_ssh_interceptor_test.go` verifies prefix matching, priority sorting, and boundary checks.
  - `cli/gitutil/git_ssh_cmd_test.go` verifies OpenSSH flag string formatting and space escaping.
