# Architecture Specification: Multi-User Git Project Switching & GitHub CLI Auth Detection

- **Specification ID:** 242-multi-user-project-switching-gh-auth-status-and-git-config-manager-01
- **Status:** Authoritative
- **Version:** 1.0.0
- **Scope:** GitMap Core, User Context Subsystem, Launch Banner, and Cross-Project Git Configuration

---

## 1. Executive Summary & Problem Statement

Engineers and automated AI agents operating on multi-account development workstations frequently interact with multiple Git hosting providers (GitHub, GitLab), multiple identity tiers (work vs. personal vs. client identities), and diverse project boundaries.

Currently:
1. When GitMap launches (`gitmap` with no arguments), it displays binary deployment locations, commit SHA, and version, but does **not** indicate which GitHub CLI identity is authorized or which Git user name/email is active.
2. The existing `gitmap user` CLI command was strictly scoped to OS-level user management (`add`, `rm`, `create-root`, `kill`, `add-ssh-key`), returning errors when developers attempt `gitmap user info` or `gitmap user switch`.
3. Managing multiple Git users across distinct repositories requires manual `git config user.name` and `git config user.email` mutations per repository, with no persistent project-to-user binding or automated contextual switching.
4. Developers lack a privacy-safe, unified command to inspect logged-in public user details (username, email, provider, auth source) without exposing private authentication tokens or sensitive credentials.

This specification designs a **unified User Context & Multi-Project Git Configuration Architecture** inside GitMap.

---

## 2. Core Architectural Pillars

### Pillar 1: GitHub CLI Authorization Status & Startup Readout
- Probes `gh auth status` and `ghtoken.Resolve()` in non-blocking, fail-safe mode (<10ms).
- Identifies whether `gh` CLI is authenticated and extracts the active public username (e.g., `alimtvnetwork`).
- Integrates into `cli/cmd/binarylocations.go` and `cli/cmd/rootusagecompact.go` to display authorization status whenever GitMap launches without arguments.
- Privacy Invariant: NEVER displays raw authentication tokens, secret keys, or private credential payloads.

### Pillar 2: Safe `user info` Command & Inspector
- Canonical Command: `gitmap user info` (and aliases `gitmap user-info`, `gitmap user status`).
- Inspects and renders:
  1. **GitHub CLI Authorization:** Status (`Authorized as <user>` or `Not logged in`), auth method/source (`gh-cli`, `Git Credential Manager`, `GH_TOKEN`).
  2. **Active Local Git Identity:** `user.name` and `user.email` from the current working tree git config.
  3. **Active Global Git Identity:** `user.name` and `user.email` from `~/.gitconfig`.
  4. **Active GitMap Profile:** Currently active profile alias and provider.
  5. **Project Binding Context:** If the current project repository is mapped to a dedicated profile alias, indicates the binding and sync status.

### Pillar 3: Multi-User Profile Management & Project Bindings
- Storage Layer: Extends `cli/model/git_profile.go` and `cli/store/git_profile_store.go` (`.gitmap/data/git_profiles.json`) with:
  ```go
  type ProjectBinding struct {
      RepoPath     string    `json:"repoPath"`     // Relative or normalized path
      ProfileAlias string    `json:"profileAlias"` // Bound profile name
      BoundAt      time.Time `json:"boundAt"`
  }
  ```
- Subcommands under `gitmap user`:
  - `gitmap user list` (alias `ls`): Displays all configured Git user profiles, marking global default and current project active profile.
  - `gitmap user add <alias> --name "<name>" --email "<email>" [--signingkey "<key>"] [--provider github]`: Enrolls a user profile. (Retains `--password` for OS user creation).
  - `gitmap user switch <alias>` / `gitmap user use <alias>`:
    - Default (inside Git repo): Configures local repository git identity (`git config user.name` and `email`) and binds the project.
    - With `--global`: Configures global `~/.gitconfig` and updates global active profile.
  - `gitmap user project bind <alias> [path]`: Explicitly maps repository to user profile.
  - `gitmap user project unbind [path]`: Clears contextual binding.
  - `gitmap user sync`: Reads the current repository's binding and applies matching git config if out of sync.

### Pillar 4: Global & Multi-Project Git Configuration
- Direct CLI commands:
  - `gitmap user config global [--name "<name>"] [--email "<email>"]`
  - `gitmap user config local [--name "<name>"] [--email "<email>"]`
- Allows setting or viewing git user configuration with affirmative feedback and validation.

---

## 3. Non-Negotiable Quality & Safety Invariants
1. **Zero Secret Leakage:** Authentication tokens are strictly masked or omitted; only public identities (`username`, `email`, provider name) are displayed.
2. **Backward Compatibility:** All existing OS user subcommands (`create-root`, `kill`, `rm`, `add-ssh-key`, `add <user> --password`) remain 100% operational.
3. **Strict Relative Git Paths:** All stored repository paths, specs, and plans use relative Git paths.
4. **Execution Performance:** Startup probe executes within $\le 15$ms via cached LookPath and non-blocking process timeout.
