# Subtask 02: Define Git URL Constants and Provider Enums & Eliminate Raw URL Literals

- **Task ID:** subtask-git-dry-02
- **Parent Task:** git-dry-url-consts-db-repository
- **Spec Reference:** 02-spec/21-app/git-dry-url-consts-db-repository/02-component-spec.md
- **Status:** READY FOR EXECUTION
- **Owner:** Worker 02
- **Stage List:**
  - cli/constants/constants_git.go
  - cli/constants/constants_git_test.go
  - cli/cmdmerge/merge_write_gh.go
  - cli/cmdpe/pe_gh.go
  - cli/cmdhealth/health_checks.go
  - cli/cmdpipeline/pipeline_query_db.go
  - cli/cmdpipeline/pipeline_target_resolve.go

---

## 1. Context & Objective

Multiple command packages across GitMap independently perform string manipulation and prefix matching on Git remote URLs (`https://github.com/`, `git@github.com:`, `http://github.com/`, `ssh://git@github.com/`, `github.com/`, `github.com:`, and `.git`). 

The objective of this subtask is to:
1. Define strongly-typed Git URL constants, provider enums, and host path tokens in `cli/constants/constants_git.go`.
2. Guarantee total collision immunity against existing declarations in package `constants`.
3. Eliminate hardcoded Git URL string literals across `cli/cmdmerge/`, `cli/cmdpe/`, `cli/cmdhealth/`, and `cli/cmdpipeline/`.
4. Author unit test verification in `cli/constants/constants_git_test.go` ensuring all constants and types remain stable and non-colliding.

---

## 2. Pre-Flight Collision Immunity Audit

Package `constants` is a large package containing multiple files (`constants_*.go`). Prior to adding new constants, we conduct an exhaustive collision analysis against symbols already declared in `cli/constants/`:

| Proposed Identifier | Proposed Type / Kind | Existing Symbol in `cli/constants/` | Collision Status | Resolution / Justification |
|---|---|---|---|---|
| `GitProviderType` | `type ... string` | None | Clean | New type name ending in `Type` |
| `GitHostType` | `type ... string` | None | Clean | New type name ending in `Type` |
| `GitURLPrefixType` | `type ... string` | None | Clean | New type name ending in `Type` |
| `GitProviderGitHub` | Constant (`GitProviderType`) | None | Clean | Fully qualified with provider prefix |
| `GitProviderGitLab` | Constant (`GitProviderType`) | None | Clean | Fully qualified with provider prefix |
| `GitProviderBitbucket` | Constant (`GitProviderType`) | None | Clean | Fully qualified with provider prefix |
| `GitProviderAzureDevOps` | Constant (`GitProviderType`) | None | Clean | Fully qualified with provider prefix |
| `GitProviderUnknown` | Constant (`GitProviderType`) | None | Clean | Standard unknown enum variant |
| `GitHostGitHub` | Constant (`GitHostType`) | None | Clean | `github.com` domain constant |
| `GitHostGitLab` | Constant (`GitHostType`) | None | Clean | `gitlab.com` domain constant |
| `GitHostBitbucket` | Constant (`GitHostType`) | None | Clean | `bitbucket.org` domain constant |
| `GitHostAzureDevOps` | Constant (`GitHostType`) | None | Clean | `dev.azure.com` domain constant |
| `GitHubHTTPSPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `GitHubHTTPPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `GitHubSSHPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `GitHubSSHSchemePrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `GitLabHTTPSPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `GitLabHTTPPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `GitLabSSHPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `GitLabSSHSchemePrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `BitbucketHTTPSPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `BitbucketHTTPPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `BitbucketSSHPrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `BitbucketSSHSchemePrefix` | Constant (`GitURLPrefixType`) | None | Clean | Unique URL prefix constant |
| `GitPrefixHTTPS` | Constant (`GitURLPrefixType`) | None | Clean | Generic scheme prefix |
| `GitPrefixHTTP` | Constant (`GitURLPrefixType`) | None | Clean | Generic scheme prefix |
| `GitPrefixSSH` | Constant (`GitURLPrefixType`) | None | Clean | Generic scheme prefix |
| `GitPrefixSSHScheme` | Constant (`GitURLPrefixType`) | None | Clean | Generic scheme prefix |
| `GitDotGitSuffix` | Constant (`string`) | `DotGit = ".git"` | Distinct | `DotGit` represents directory `.git`; `GitDotGitSuffix` represents URL extension `.git`. Distinct identifier ensures zero compilation collision. |
| `GitPathSeparator` | Constant (`string`) | None | Clean | URL path delimiter `/` |
| `GitColonSeparator` | Constant (`string`) | None | Clean | SCP syntax delimiter `:` |
| `GitHostTokenGitHubSlash`| Constant (`string`) | None | Clean | Boundary token `"github.com/"` |
| `GitHostTokenGitHubColon`| Constant (`string`) | None | Clean | Boundary token `"github.com:"` |
| `DefaultGitHubOwner` | Constant (`string`) | None | Clean | Organization fallback `"alimtvnetwork"` |

---

## 3. Detailed Step-by-Step Execution Plan

### Step 1: Implement Constants & Enums in `cli/constants/constants_git.go`
- Target file: `cli/constants/constants_git.go`
- Append the new typed enums, constants, delimiters, and host tokens at the bottom of `cli/constants/constants_git.go`.
- Ensure clean formatting:
  - Blank line after struct/type definitions.
  - Group related constants with descriptive Go doc comments.
  - Maintain the total line count of `cli/constants/constants_git.go` <= 220 lines (well under the 300-line limit).

### Step 2: Implement Test Suite `cli/constants/constants_git_test.go`
- Target file: `cli/constants/constants_git_test.go`
- Write table-driven unit tests verifying:
  - Value fidelity: Check that all constants evaluate to their exact specified string representations.
  - Prefix relationships: Ensure `GitHubHTTPSPrefix` has prefix `GitPrefixHTTPS`, `GitHubSSHPrefix` has prefix `GitPrefixSSH`, etc.
  - Enum exhaustiveness: Test that each provider has an associated host constant.
- File size constraint: <= 120 lines.

### Step 3: Refactor Call Sites in `cli/cmdmerge/merge_write_gh.go`
- Target file: `cli/cmdmerge/merge_write_gh.go`
- In function `parseGitHubSlug(remote string)`:
  - Replace `strings.TrimSuffix(remote, ".git")` with `strings.TrimSuffix(remote, constants.GitDotGitSuffix)`.
  - Replace `strings.Replace(rest, "git@github.com:", "https://github.com/", 1)` with:
    `strings.Replace(rest, string(constants.GitHubSSHPrefix), string(constants.GitHubHTTPSPrefix), 1)`
  - Replace `strings.Replace(rest, "ssh://git@github.com/", "https://github.com/", 1)` with:
    `strings.Replace(rest, string(constants.GitHubSSHSchemePrefix), string(constants.GitHubHTTPSPrefix), 1)`
  - Replace prefix checks with `constants.GitHubHTTPSPrefix` and `constants.GitHubHTTPPrefix`.
  - Ensure zero nested ifs and function length <= 15 lines.

### Step 4: Refactor Call Sites in `cli/cmdpe/pe_gh.go`
- Target file: `cli/cmdpe/pe_gh.go`
- In function `parseGitHubSlug(raw string)`:
  - Replace case branches:
    - `"https://github.com/"` -> `string(constants.GitHubHTTPSPrefix)`
    - `"http://github.com/"` -> `string(constants.GitHubHTTPPrefix)`
    - `"git@github.com:"` -> `string(constants.GitHubSSHPrefix)`
    - `"ssh://git@github.com/"` -> `string(constants.GitHubSSHSchemePrefix)`
  - Replace `strings.TrimSuffix(path, ".git")` with `strings.TrimSuffix(path, constants.GitDotGitSuffix)`.
  - Ensure positive boolean naming: `isGitHub` returned from caller.

### Step 5: Refactor Call Sites in `cli/cmdhealth/health_checks.go`
- Target file: `cli/cmdhealth/health_checks.go`
- In function `healthRepoSlug(dir string)`:
  - Replace case branches:
    - `"https://github.com/"` -> `string(constants.GitHubHTTPSPrefix)`
    - `"git@github.com:"` -> `string(constants.GitHubSSHPrefix)`
  - Replace `strings.TrimSuffix(path, ".git")` with `strings.TrimSuffix(path, constants.GitDotGitSuffix)`.

### Step 6: Refactor Call Sites in `cli/cmdpipeline/`
- Target file: `cli/cmdpipeline/pipeline_query_db.go`
  - In `resolveCurrentRepoSlug`:
    - Replace hardcoded `"alimtvnetwork/gitmap-v28"` with `constants.DefaultGitHubOwner + "/gitmap-v28"`.
  - In `parseSlugFromGitURL`:
    - Replace `strings.TrimSuffix(raw, ".git")` with `strings.TrimSuffix(raw, constants.GitDotGitSuffix)`.
    - Replace `"github.com/"` with `constants.GitHostTokenGitHubSlash`.
    - Replace `"github.com:"` with `constants.GitHostTokenGitHubColon`.
- Target file: `cli/cmdpipeline/pipeline_target_resolve.go`
  - In `isPipelineGitURL(s string)`:
    - Replace `"https://"` with `string(constants.GitPrefixHTTPS)`.
    - Replace `"http://"` with `string(constants.GitPrefixHTTP)`.
    - Replace `"git@"` with `string(constants.GitPrefixSSH)`.
    - Replace `"ssh://"` with `string(constants.GitPrefixSSHScheme)`.
    - Replace `".git"` with `constants.GitDotGitSuffix`.
    - Replace `"github.com/"` with `constants.GitHostTokenGitHubSlash`.
    - Replace `"github.com:"` with `constants.GitHostTokenGitHubColon`.
  - In `resolveDirectorySlug`:
    - Replace `"alimtvnetwork/" + base` with `constants.DefaultGitHubOwner + "/" + base`.

---

## 4. Quality Rules & Guidelines Checklist

- [ ] All new types end with `Type` (`GitProviderType`, `GitHostType`, `GitURLPrefixType`).
- [ ] Blank line after each `}` and before each `return`.
- [ ] No function exceeds 15 lines.
- [ ] No nested if blocks (flat guard clauses only).
- [ ] All booleans adhere to `is*` or `has*` positive prefixes.
- [ ] Zero compiler build or test commands executed during task.

---

## 5. Verification Evidence Requirements

Upon completion of changes by Worker 02:
1. `cli/constants/constants_git.go`: Show diff containing new enums and constants.
2. `cli/constants/constants_git_test.go`: Show authored unit test covering all constants and prefixes.
3. Code search verifying zero remaining literal `"https://github.com/"` and `"git@github.com:"` occurrences across the target files.
