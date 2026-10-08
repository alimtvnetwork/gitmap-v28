# Component & CLI Specification: User Context, Project Switching & Git Config

- **Specification ID:** 242-multi-user-project-switching-gh-auth-status-and-git-config-manager-02
- **Status:** Authoritative
- **Version:** 1.0.0
- **Scope:** CLI Routing, Domain Model Extensions, Terminal Renderers, and Testing Matrix

---

## 1. Domain Model Extensions (`cli/model/git_profile.go`)

Extend `GitProfileConfig` and `GitProfile` to model per-project contextual bindings:

```go
// ProjectBinding associates a local repository directory with a Git user profile.
type ProjectBinding struct {
	RepoSlug     string    `json:"repoSlug"`     // Normalized repository slug or relative directory
	ProfileAlias string    `json:"profileAlias"` // Name of the bound GitProfile
	BoundAt      time.Time `json:"boundAt"`
}

// GitProfileConfig holds all configured Git profiles, project bindings, and active defaults.
type GitProfileConfig struct {
	Profiles        []GitProfile              `json:"profiles"`
	ProjectBindings map[string]ProjectBinding `json:"projectBindings,omitempty"`
	Active          string                    `json:"active"`
	Default         string                    `json:"default"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
}
```

---

## 2. User Context Inspector Engine (`cli/usercontext/`)

Create dedicated leaf package `cli/usercontext/` containing:
1. `types.go`:
   - `UserSummary`: Encapsulates GitHub CLI auth, active Git user (local/global), current profile, and project binding.
   - `AuthStatus`: Enum (`Authorized`, `Unauthenticated`, `ToolMissing`).
2. `gh_auth.go`:
   - `DetectGhAuthStatus() (username string, status AuthStatus, err error)`:
     Probes `gh auth status` or `gh api user --jq .login`.
     Falls back to `ghtoken.Resolve()` if `gh` is unauthenticated or missing.
3. `git_identity.go`:
   - `GetGitIdentity(isGlobal bool) (name, email string, err error)`:
     Queries `git config user.name` and `git config user.email` with or without `--global`.
   - `SetGitIdentity(name, email string, isGlobal bool) error`:
     Applies git config commands and verifies exit code.
4. `project_binder.go`:
   - `ResolveProjectUser(repoPath string) (*model.GitProfile, bool)`:
     Matches the current repository working directory against `ProjectBindings`.
   - `BindProjectUser(repoPath, profileAlias string) error`:
     Updates `ProjectBindings` in `git_profiles.json`.
5. `render.go`:
   - Renders Catppuccin Macchiato cards for `gitmap user info` and startup summary.

---

## 3. CLI Command Matrix

| Command | Arguments / Flags | Description |
| :--- | :--- | :--- |
| `gitmap user info` (or `user-info`) | `[--json]` | Displays comprehensive public user information, GitHub CLI auth status, and Git identities. |
| `gitmap user list` (or `ls`) | `[--json]` | Lists all configured Git profiles, marking default and active project user. |
| `gitmap user add <alias>` | `--name "<name>" --email "<email>" [--org] [--provider github\|gitlab]` | Registers a Git user profile for project/global use. (Note: `--password` triggers OS user creation). |
| `gitmap user switch <alias>` (or `use`) | `[--global] [--project]` | Switches active user. In a Git repo, sets local identity and binds project; with `--global`, updates global config. |
| `gitmap user project bind <alias>` | `[path]` | Binds a Git profile to a target repository path (defaults to current dir). |
| `gitmap user project unbind` | `[path]` | Removes the project-specific user binding. |
| `gitmap user config global` | `[--name "<name>"] [--email "<email>"]` | Inspects or sets global git configuration. |
| `gitmap user config local` | `[--name "<name>"] [--email "<email>"]` | Inspects or sets local repository git configuration. |
| `gitmap user sync` | None | Applies bound user identity to the current repository if divergent. |

---

## 4. Startup Readout Integration (`cli/cmd/rootusagecompact.go`)

When `gitmap` is launched with zero arguments, the compact readout includes:
```
  ────────────────────────────────────────────────────────────
  Git & Provider Status
  ● GitHub CLI:     ✔ Authorized as alimtvnetwork
  ● Active Git:     MD ALIM UL KARIM <devorg.bd@gmail.com> (global)
  ● Project User:   (inherited from global)
  ────────────────────────────────────────────────────────────
```
If GitHub CLI is unauthenticated:
```
  ● GitHub CLI:     ✖ Not logged in (run 'gh auth login' or 'gitmap login')
```

---

## 5. Acceptance Criteria & Test Plan

1. **GitHub CLI Auth Detection:**
   - Detects authorized username when `gh` is authenticated.
   - Gracefully reports unauthenticated when `gh` returns non-zero.
2. **User Info Command:**
   - `gitmap user info` succeeds in and outside Git repositories.
   - Emits clean JSON when `--json` flag is provided.
   - Zero secrets or token strings leaked in stdout/stderr.
3. **Multi-User Switching:**
   - `gitmap user add work --name "Work User" --email "work@example.com"` creates profile.
   - `gitmap user switch work` inside a repo changes local `git config user.name/email`.
   - `gitmap user list` displays `work` profile with usage counters.
4. **Backward Compatibility:**
   - `gitmap user add <name> --password <pwd>` still creates an OS user on Linux/Windows.
   - `gitmap user create-root` and `gitmap user kill` remain unaffected.
