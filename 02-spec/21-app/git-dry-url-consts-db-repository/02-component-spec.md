# Component Specification: Git URL Constants, Provider Enums & Subsystem Refactoring

- **Document ID:** 02-component-spec.md
- **Parent Task:** git-dry-url-consts-db-repository
- **Specification Group:** 02-spec/21-app/git-dry-url-consts-db-repository/
- **Related Plans:** 
  - .ai-memory/plans/subtasks/git-dry-url-consts-db-repository/02-git-url-constants-and-enums.md
  - .ai-memory/plans/subtasks/git-dry-url-consts-db-repository/04-pull-remediation-refactor.md
- **Status:** DRAFT / APPROVED FOR IMPLEMENTATION

---

## 1. Executive Summary & Architectural Scope

The GitMap codebase currently contains fragmented, duplicate implementations for Git URL handling, remote URL resolution, and repository transport persistence. Specifically:
1. **Hardcoded Git URL Strings & Host Matching:** Subsystems across `cli/cmdpull/`, `cli/cmdpe/`, `cli/cmdhealth/`, `cli/cmdpipeline/`, and `cli/cmdmerge/` repetitively define raw string literals (`"https://github.com/"`, `"http://github.com/"`, `"git@github.com:"`, `"ssh://git@github.com/"`, `"github.com/"`, `"github.com:"`, and `".git"`).
2. **Duplicate Subprocess Git Execution:** Functions such as `readRemoteOriginURL` in `cli/cmdpull/pull_remediation_ssh.go`, `githubSlugForDir` in `cli/cmdpe/pe_gh.go`, `healthRepoSlug` in `cli/cmdhealth/health_checks.go`, and `resolveLocalOrCwdSlug` in `cli/cmdpipeline/pipeline_target_resolve.go` each independently invoke `exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")`.
3. **Repository Pattern Violation (Raw SQL in Business Logic):** `cli/cmdpull/pull_remediation_ssh.go` executes raw SQL directly against `db.Conn().Exec()` rather than delegating data persistence to the domain repository layer (`cli/store/`).

This specification establishes:
- Standardized, collision-safe Git URL constants and provider enums in `cli/constants/constants_git.go`.
- A multi-package call-site refactoring architecture replacing ad-hoc string comparisons and raw git calls with centralized constants and `cli/gitutil` helpers.
- An end-to-end component specification for `cli/cmdpull/pull_remediation_ssh.go` to eliminate private git execution, leverage typed constants, and consume the repository pattern via `*store.DB.UpdateRepoTransport`.

---

## 2. Component 1: Typed Git URL Constants & Provider Enums

All constants and enums reside in package `constants` within `cli/constants/constants_git.go`. They conform to Go idiomatic typing, the project's enum naming standards (type name suffixed with `Type`), and zero magic strings.

### 2.1 Provider & Host Enums

```go
// GitProviderType defines recognized Git hosting providers.
type GitProviderType string

const (
	GitProviderGitHub      GitProviderType = "github"
	GitProviderGitLab      GitProviderType = "gitlab"
	GitProviderBitbucket   GitProviderType = "bitbucket"
	GitProviderAzureDevOps GitProviderType = "azuredevops"
	GitProviderUnknown     GitProviderType = "unknown"
)

// GitHostType defines normalized domain names for Git hosts.
type GitHostType string

const (
	GitHostGitHub      GitHostType = "github.com"
	GitHostGitLab      GitHostType = "gitlab.com"
	GitHostBitbucket   GitHostType = "bitbucket.org"
	GitHostAzureDevOps GitHostType = "dev.azure.com"
)
```

### 2.2 Typed Git URL Prefixes

```go
// GitURLPrefixType represents canonical protocol and scheme prefixes for Git remotes.
type GitURLPrefixType string

const (
	// GitHub remote prefixes
	GitHubHTTPSPrefix     GitURLPrefixType = "https://github.com/"
	GitHubHTTPPrefix      GitURLPrefixType = "http://github.com/"
	GitHubSSHPrefix       GitURLPrefixType = "git@github.com:"
	GitHubSSHSchemePrefix GitURLPrefixType = "ssh://git@github.com/"

	// GitLab remote prefixes
	GitLabHTTPSPrefix     GitURLPrefixType = "https://gitlab.com/"
	GitLabHTTPPrefix      GitURLPrefixType = "http://gitlab.com/"
	GitLabSSHPrefix       GitURLPrefixType = "git@gitlab.com:"
	GitLabSSHSchemePrefix GitURLPrefixType = "ssh://git@gitlab.com/"

	// Bitbucket remote prefixes
	BitbucketHTTPSPrefix     GitURLPrefixType = "https://bitbucket.org/"
	BitbucketHTTPPrefix      GitURLPrefixType = "http://bitbucket.org/"
	BitbucketSSHPrefix       GitURLPrefixType = "git@bitbucket.org:"
	BitbucketSSHSchemePrefix GitURLPrefixType = "ssh://git@bitbucket.org/"

	// Generic transport scheme prefixes
	GitPrefixHTTPS     GitURLPrefixType = "https://"
	GitPrefixHTTP      GitURLPrefixType = "http://"
	GitPrefixSSH       GitURLPrefixType = "git@"
	GitPrefixSSHScheme GitURLPrefixType = "ssh://"
)
```

### 2.3 Delimiters, Suffixes & Host Path Tokens

```go
const (
	// GitDotGitSuffix represents the canonical repository extension.
	GitDotGitSuffix = ".git"

	// Delimiters
	GitPathSeparator  = "/"
	GitColonSeparator = ":"

	// Host path tokens for substring and boundary tokenization
	GitHostTokenGitHubSlash = "github.com/"
	GitHostTokenGitHubColon = "github.com:"

	// Default organization/owner fallback
	DefaultGitHubOwner = "alimtvnetwork"
)
```

### 2.4 Collision Immunity Verification

Existing constants in `cli/constants/constants_git.go` already declare:
- `TransportSSH = "ssh"`
- `DotGit = ".git"` (Used for local directory detection `.git`)
- `CloneInstructionFmt = "git clone -b %s %s %s"`
- `HTTPSFromSSHFmt = "https://%s/%s"`
- `SSHFromHTTPSFmt = "git@%s:%s"`

The new constants do NOT collide with existing identifiers:
- `GitDotGitSuffix` (`.git`) is explicitly distinct from `DotGit` (`.git`) and intended for remote URL string trimming.
- `GitProviderType`, `GitHostType`, and `GitURLPrefixType` introduce new strongly-typed types.
- All new constant identifiers start with their respective provider or category prefix (`GitHub*`, `GitLab*`, `Bitbucket*`, `GitPrefix*`, `GitHostToken*`).

---

## 3. Component 2: Multi-Package Call-Site Refactoring Architecture

Across the command packages, five distinct call-site patterns must be refactored to consume the centralized constants and `cli/gitutil` helpers.

```
+-----------------------------------------------------------------------------------+
|                                  cli/constants/                                   |
|   GitProviderType  |  GitHostType  |  GitURLPrefixType  |  GitHostPathTokens      |
+-----------------------------------------------------------------------------------+
                                         ^
                                         | Consumes
+-----------------------------------------------------------------------------------+
|                                   cli/gitutil/                                    |
|   GetRemoteOriginURL(repoPath) | ParseGitHubSlug(url) | ExecGitInDir(repoPath,..) |
+-----------------------------------------------------------------------------------+
       ^                        ^                        ^                   ^
       |                        |                        |                   |
+---------------+        +---------------+        +---------------+  +---------------+
|  cli/cmdpull/ |        |  cli/cmdpe/   |        | cli/cmdhealth/|  |cli/cmdpipeline|
|  pull_remed.. |        |  pe_gh.go     |        | health_checks |  | pipeline_..   |
+---------------+        +---------------+        +---------------+  +---------------+
       |
       v Consumes
+-----------------------------------------------------------------------------------+
|                                   cli/store/                                      |
|            UpdateRepoTransport(repoPath, sshURL, transport string) error          |
+-----------------------------------------------------------------------------------+
```

### 3.1 `cli/cmdpull/` Call-Site Architecture

1. **`cli/cmdpull/pull_remediation_ssh.go`**:
   - **Current State:**
     - Defines private `readRemoteOriginURL(repoPath string) string` executing raw `exec.Command("git", "-C", repoPath, "config", "--get", "remote.origin.url")`.
     - Uses string literals `"https://github.com/"`, `"git@github.com:"`, `"alimtvnetwork"`.
     - Executes direct raw SQL via `db.Conn().Exec("UPDATE Repo SET ...")`.
   - **Target Architecture:**
     - Delete private `readRemoteOriginURL` and replace all calls with `gitutil.GetRemoteOriginURL(repoPath)` (or `gitutil.ReadRemoteOriginURL(repoPath)`).
     - Replace string prefixes with `constants.GitHubHTTPSPrefix`, `constants.GitHubSSHPrefix`, and `constants.DefaultGitHubOwner`.
     - Replace raw SQL with `db.UpdateRepoTransport(repoPath, newURL, constants.TransportSSH)`.
2. **`cli/cmdpull/pull_remediation_rebase.go`**:
   - **Current State:** Direct `exec.Command("git", "-C", repoPath, "pull", "--progress", "--rebase", "--autostash")`.
   - **Target Architecture:** Route rebase execution through `gitutil.ExecGitInDirCombined(repoPath, "pull", "--progress", "--rebase", "--autostash")`.
3. **`cli/cmdpull/pull_prompt_step.go`, `pull_prompt_conflict.go`, `pull_prompt_exec.go`**:
   - **Current State:** Direct `exec.Command("git", "-C", repoPath, "status", "--porcelain")`.
   - **Target Architecture:** Route git status inspection through `gitutil.GetStatusPorcelain(repoPath)`.

### 3.2 `cli/cmdpe/` Call-Site Architecture

1. **`cli/cmdpe/pe_gh.go`**:
   - **Function `githubSlugForDir(dir string) (string, bool)`:**
     - **Current State:** Invokes `exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")`.
     - **Target Architecture:** Invoke `origin, err := gitutil.GetRemoteOriginURL(dir); if err != nil || origin == "" { return "", false }`.
   - **Function `parseGitHubSlug(raw string) (string, bool)`:**
     - **Current State:** Uses hardcoded cases:
       `case strings.HasPrefix(lower, "https://github.com/"):`
       `case strings.HasPrefix(lower, "http://github.com/"):`
       `case strings.HasPrefix(lower, "git@github.com:"):`
       `case strings.HasPrefix(lower, "ssh://git@github.com/"):`
       `path = strings.TrimSuffix(path, ".git")`
     - **Target Architecture:** Match against typed constants:
       `case strings.HasPrefix(lower, string(constants.GitHubHTTPSPrefix)):`
       `case strings.HasPrefix(lower, string(constants.GitHubHTTPPrefix)):`
       `case strings.HasPrefix(lower, string(constants.GitHubSSHPrefix)):`
       `case strings.HasPrefix(lower, string(constants.GitHubSSHSchemePrefix)):`
       `path = strings.TrimSuffix(path, constants.GitDotGitSuffix)`
   - **Function `gitLastCommitAge(dir string) string`:**
     - **Current State:** Direct `exec.Command("git", "-C", dir, "log", "-1", "--format=%ar")`.
     - **Target Architecture:** Delegate to `gitutil.ExecGitInDirOutput(dir, "log", "-1", "--format=%ar")`.

### 3.3 `cli/cmdhealth/` Call-Site Architecture

1. **`cli/cmdhealth/health_checks.go`**:
   - **Function `healthRepoSlug(dir string) (string, bool)`:**
     - **Current State:** Exact copy of `pe_gh.go`'s `githubSlugForDir` and `parseGitHubSlug` with raw `exec.Command` and hardcoded `"https://github.com/"` and `"git@github.com:"`.
     - **Target Architecture:** Replace git execution with `gitutil.GetRemoteOriginURL(dir)` and slug parsing with centralized `gitutil.ParseGitHubSlug(raw)` or typed constants (`constants.GitHubHTTPSPrefix`, `constants.GitHubSSHPrefix`, `constants.GitDotGitSuffix`).
   - **Health Status Check:**
     - **Current State:** Line 58 invokes `exec.Command("git", "-C", repoPath, "status", "--porcelain")`.
     - **Target Architecture:** Replace with `gitutil.GetStatusPorcelain(repoPath)`.

### 3.4 `cli/cmdpipeline/` Call-Site Architecture

1. **`cli/cmdpipeline/pipeline_query_db.go`**:
   - **Function `resolveCurrentRepoSlug() string`:**
     - **Current State:** Direct `exec.Command("git", "config", "--get", "remote.origin.url")`.
     - **Target Architecture:** Replace with `gitutil.GetRemoteOriginURL(".")`. Default fallback uses `constants.DefaultGitHubOwner + "/gitmap-v28"`.
   - **Function `parseSlugFromGitURL(raw string) string`:**
     - **Current State:** Hardcoded `strings.Contains(clean, "github.com/")` and `"github.com:"` and trimming `".git"`.
     - **Target Architecture:** Use `constants.GitDotGitSuffix`, `constants.GitHostTokenGitHubSlash`, and `constants.GitHostTokenGitHubColon`.
2. **`cli/cmdpipeline/pipeline_target_resolve.go`**:
   - **Functions `resolveLocalOrCwdSlug` & `resolveDirectorySlug`:**
     - **Current State:** Direct `exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")`.
     - **Target Architecture:** Replace with `gitutil.GetRemoteOriginURL(dir)`.
   - **Function `isPipelineGitURL(s string) bool`:**
     - **Current State:** Hardcoded prefix checks (`"http://"`, `"https://"`, `"git@"`, `"ssh://"`, `".git"`, `"github.com/"`, `"github.com:"`).
     - **Target Architecture:** Use typed constants:
       `strings.HasPrefix(low, string(constants.GitPrefixHTTPS)) ||`
       `strings.HasPrefix(low, string(constants.GitPrefixHTTP)) ||`
       `strings.HasPrefix(low, string(constants.GitPrefixSSH)) ||`
       `strings.HasPrefix(low, string(constants.GitPrefixSSHScheme)) ||`
       `strings.HasSuffix(low, constants.GitDotGitSuffix) ||`
       `strings.Contains(low, constants.GitHostTokenGitHubSlash) ||`
       `strings.Contains(low, constants.GitHostTokenGitHubColon)`

### 3.5 `cli/cmdmerge/` Call-Site Architecture

1. **`cli/cmdmerge/merge_write_gh.go`**:
   - **Function `parseGitHubSlug(remote string) (owner, repo string, isGitHub bool)`:**
     - **Current State:**
       `rest := strings.TrimSuffix(remote, ".git")`
       `rest = strings.Replace(rest, "git@github.com:", "https://github.com/", 1)`
       `rest = strings.Replace(rest, "ssh://git@github.com/", "https://github.com/", 1)`
       `if !strings.HasPrefix(rest, "https://github.com/") && !strings.HasPrefix(rest, "http://github.com/")`
     - **Target Architecture:** Replace all raw literals with:
       `constants.GitDotGitSuffix`
       `constants.GitHubSSHPrefix`
       `constants.GitHubSSHSchemePrefix`
       `constants.GitHubHTTPSPrefix`
       `constants.GitHubHTTPPrefix`

---

## 4. Component 3: Detailed Specification for `cli/cmdpull/pull_remediation_ssh.go`

`cli/cmdpull/pull_remediation_ssh.go` is the focal remediation unit for SSH transport fallback. It will undergo a complete structural refactoring.

### 4.1 Structural Changes Breakdown

1. **Remove `readRemoteOriginURL`:**
   - **Action:** Delete lines 14–21 of `cli/cmdpull/pull_remediation_ssh.go`.
   - **Replacement:** Call `gitutil.GetRemoteOriginURL(repoPath)` across callers.
2. **Update `isSameOwnerHTTPSRepo`:**
   - **Action:** Replace `readRemoteOriginURL` call and use `constants.GitHubHTTPSPrefix`.
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
3. **Update `extractGitHubOwnerFromURL`:**
   - **Action:** Refactor branch conditions to use typed constants.
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
4. **Update `hasSSHRepoForSameOwner`:**
   - **Action:** Use `constants.DefaultGitHubOwner`.
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
   ```
5. **Update `hasSSHRepoOwner`:**
   - **Action:** Construct prefix using `constants.GitHubSSHPrefix` and `constants.GitPathSeparator`.
   ```go
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
6. **Refactor `updateRepoTransportInDB` to Consume Repository Pattern:**
   - **Action:** Eliminate direct SQL invocation `db.Conn().Exec(...)`.
   - **Target Repository Method:**
     In `cli/store/repo_transport.go`:
     ```go
     func (db *DB) UpdateRepoTransport(repoPath, sshURL, transport string) error
     ```
   - **Refactored `updateRepoTransportInDB`:**
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
7. **Refactor Credential Configuration Execution:**
   - **Action:** Replace direct `exec.Command("git", "config", ...)` with `gitutil.SetGlobalGitConfig`:
     ```go
     func remediateCredentialAndConvertToSSH(repoPath, repoName string) {
         storeType := gitutil.ResolveWindowsCredentialStore()
         _ = gitutil.SetGlobalGitConfig("credential.helper", "manager")
         _ = gitutil.SetGlobalGitConfig("credential.credentialStore", string(storeType))
         applySSHTransportRemediation(repoPath)
     }
     ```

---

## 5. Coding Guidelines & Invariant Compliance

Every refactored component and call site must adhere strictly to the repository's coding guidelines:
- **Error Return Sovereignty:** Leaf functions returning errors bubble them up directly; no exit handlers in leaf functions.
- **Positive Booleans Only:** All booleans must begin with `is*` or `has*` (`isGitHub`, `hasChanged`, `isConflict`). Banned prefixes: `can`, `should`, `was`, `will`, `no*`, `not*`.
- **Zero Raw SQL in Business Packages:** All SQL statements belong exclusively in `cli/store/` behind `*store.DB` repository methods.
- **Zero Direct Git Subprocess Invocations in Command Handlers:** All `git -C` and git command invocations belong exclusively in `cli/gitutil/`.
- **Function Size Limits:** <= 15 lines per function, <= 8 lines preferred.
- **Zero Nested Ifs:** Flatten with early returns and guard clauses.
- **File Length Limits:** Strict <= 300 lines per Go source file.

---

## 6. Verification & Acceptance Criteria

1. **Static Analysis & Inspection:**
   - Package `constants` exposes `GitProviderType`, `GitHostType`, `GitURLPrefixType`, and all specified constants in `cli/constants/constants_git.go`.
   - Zero occurrences of raw string literals for GitHub/GitLab URL prefixes in `cli/cmdpull/pull_remediation_ssh.go`, `cli/cmdpe/pe_gh.go`, `cli/cmdhealth/health_checks.go`, `cli/cmdpipeline/`, and `cli/cmdmerge/merge_write_gh.go`.
   - Zero occurrences of `db.Conn().Exec` in `cli/cmdpull/pull_remediation_ssh.go`.
   - Zero occurrences of `readRemoteOriginURL` in `cli/cmdpull/pull_remediation_ssh.go`.
2. **Behavioral Equivalence:**
   - SSH remediation behaves identically when encountering GitHub HTTPS repositories with owner matches.
   - Remote slug resolution in `cmdpe`, `cmdhealth`, and `cmdpipeline` parses identical owner/repo slugs across all supported URL forms.
   - Database records in `Repo` correctly update `SshUrl`, `DiscoveredUrl`, and `IdentifiedTransport`.
