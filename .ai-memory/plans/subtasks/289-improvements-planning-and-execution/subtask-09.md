# Subtask Plan 09: Centralize Remaining Raw URL Protocol Literals into `GitURLPrefixType`

- **Subtask ID:** `subtask-09`
- **Subtask Name:** `09-constants-git-url-prefix-type`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/constants/constants_git.go`
  - `cli/archive/source.go`
  - `cli/clonefrom/validate.go`
  - `cli/cmdcommitin/commitin.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/constants/constants_git.go`, `GitURLPrefixType` standardizes canonical protocol prefixes for GitHub, GitLab, and Bitbucket (`GitHubHTTPSPrefix`, `GitPrefixHTTPS`, `GitPrefixSSH`).

However, multiple packages continue to embed raw string literals:
1. `cli/archive/source.go`: `strings.HasPrefix(s, "git://")`
2. `cli/clonefrom/validate.go`: `prefixes := []string{"https://", "http://", "ssh://", "git://", "file://"}`
3. `cli/cmdcommitin/commitin.go`: `prefixes := []string{"http://", "https://", "git@", "ssh://", "git://"}`
4. `cli/cmdscan/scan_rc_merge.go`: `strings.TrimPrefix(u, "git://")`
5. Azure DevOps remote patterns: raw string `"https://dev.azure.com/"` and `"git@ssh.dev.azure.com:v3/"`.

Under **Rule 5 (Strongly Typed Enums)**:
- All URL schemes and remote prefixes must be consolidated under `GitURLPrefixType`.
- Enforce the `Type` suffix naming rule.
- Replace scattered raw string literals across the codebase with their strongly typed constant counterparts.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Extend `GitURLPrefixType` in `cli/constants/constants_git.go`
- **File:** `cli/constants/constants_git.go`
- **Additions to `GitURLPrefixType`:**
  ```go
  const (
  	// Azure DevOps remote prefixes.
  	AzureHTTPSPrefix GitURLPrefixType = "https://dev.azure.com/"
  	AzureSSHPrefix   GitURLPrefixType = "git@ssh.dev.azure.com:v3/"

  	// Additional generic transport scheme prefixes.
  	GitPrefixGit       GitURLPrefixType = "git://"
  	GitPrefixFile      GitURLPrefixType = "file://"
  	GitPrefixSSHGitAt  GitURLPrefixType = "ssh://git@"
  )
  ```

### 2.2 Refactor Consumer Packages to Use Typed Constants
- **`cli/archive/source.go`:**
  - Replace `strings.HasPrefix(s, "git://")` with:
    `strings.HasPrefix(s, string(constants.GitPrefixGit))`
- **`cli/clonefrom/validate.go`:**
  - Replace raw literal slice with:
    ```go
    prefixes := []string{
    	string(constants.GitPrefixHTTPS),
    	string(constants.GitPrefixHTTP),
    	string(constants.GitPrefixSSHScheme),
    	string(constants.GitPrefixGit),
    	string(constants.GitPrefixFile),
    }
    ```
- **`cli/cmdcommitin/commitin.go`:**
  - Replace raw literals with typed constants:
    ```go
    prefixes := []string{
    	string(constants.GitPrefixHTTP),
    	string(constants.GitPrefixHTTPS),
    	string(constants.GitPrefixSSH),
    	string(constants.GitPrefixSSHScheme),
    	string(constants.GitPrefixGit),
    }
    ```

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Strongly Typed Enums:** All constants declared under `GitURLPrefixType` with `Type` suffix.
- [ ] **No Magic Strings:** Raw `"git://"`, `"file://"`, `"ssh://git@"` eradicated from callers.
- [ ] **15-Line Limit:** Function bodies consuming prefixes remain concise and within limits.
- [ ] **Zero Nested Ifs:** Preserved across consuming modules.
- [ ] **Return Styling:** Standard vertical whitespace before returns.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.
