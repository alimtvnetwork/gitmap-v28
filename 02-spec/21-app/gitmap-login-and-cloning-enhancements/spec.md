# Specification: gitmap-login-and-cloning-enhancements

## 1. Overview & Objective
Resolve critical authentication and cloning friction points across GitMap:
1. Fix GitMap cloning when using SSH URLs (`git@gitlab.com:...`, `git@github.com:...`) or `--ssh` flag by eliminating false-negative SSH connectivity probes and bypassing inappropriate Personal Access Token (PAT) / OAuth prompts.
2. Eliminate the Windows Credential Manager / GitHub Desktop sign-in modal that unexpectedly pops up on fresh GitMap execution.
3. Ensure Git identity reporting strictly queries global Git config (`user.name`, `user.email`) and never invents placeholder usernames (`"default"`, `"github-user"`).
4. Implement a unified, comprehensive `gitmap login` subsystem supporting:
   - Interactive login method chooser.
   - `gitmap login ssh [alias] [-y]` with automated key generation, existing key reuse, clipboard copying, and service settings links.
   - `gitmap login browser-flow` for OAuth/device code flows.
   - `gitmap login token` / `pat` for secure token storage.
   - Multi-account / multi-profile management with per-repository SSH key binding via `core.sshCommand`.

---

## 2. Root Cause & Architectural Remediation

### 2.1 SSH Probe Flag Bug (`cli/cmdclone/clone_auth_probe.go`)
- **Root Cause:** `executeSSHProbe` executes `git ls-remote --exit-code -h <url> HEAD`. The `-h` flag filters strictly for `refs/heads/*`. Because `HEAD` is a root symbolic ref, `git ls-remote` returns 0 matches and exits with code 2 on *every repository on earth*.
- **Fix:** Remove `-h` so the command is `git ls-remote --exit-code <url> HEAD`.

### 2.2 Unconditional Auth Probe on Explicit SSH URLs (`cli/cmdclone/clone.go`, `clonemulti.go`, `clone_auth_resolve.go`)
- **Root Cause:** `executeDirectClone` forces SSH URLs into `rec.HTTPSUrl` and calls `ResolveRepoAuth`. When `ProbeSSHAccess` fails, it falls back to `resolveHTTPSAuth` and demands a PAT for an SSH URL!
- **Fix:**
  - In `ResolveRepoAuth`: if `isSSHCloneURL(rawURL)`, immediately return `rec.Transport = "ssh"`, `rec.SSHUrl = rawURL`, and `nil` error.
  - In `executeDirectClone` and `executeDirectCloneOne`: if `isSSHCloneURL(url)` or `params.UseSSH`, bypass `ResolveRepoAuth` and directly invoke `git clone`.
  - Only if `git clone` exits with a genuine authentication error (`Permission denied (publickey)` or `Authentication failed`) should GitMap display auth guidance.

### 2.3 Windows Credential Manager Popup on Startup (`cli/secrets/ghtoken.go`)
- **Root Cause:** `tokenFromGitCredential()` executes `git credential fill` with only `GIT_TERMINAL_PROMPT=0`. On Windows, Git Credential Manager (GCM) launches an interactive GUI popup window ("Connect to GitHub") when credentials are not found.
- **Fix:** Add `GCM_INTERACTIVE=never` to environment variables and pass `-c credential.interactive=never` to `git`.

### 2.4 Active Git Global Identity Enforcement (`cli/cmd/rootusagecompact.go`)
- **Root Cause:** `detectActiveGitUser()` runs `git config user.name` without `--global`, picking up repository-local or invented author names.
- **Fix:** Strictly execute `git config --global user.name` and `git config --global user.email`. If either is unset, display `✖ Not configured / Not logged in` and show login suggestions. Remove all placeholder fallbacks (`"default"`, `"github-user"`).

### 2.5 Unified `gitmap login` Subsystem (`cli/cmdlogin/`)
- **Interactive:** `gitmap login` prompts:
  1) SSH Key Setup & Copy
  2) Access Token (PAT)
  3) Browser Device Flow
  4) Authentication Status
- **Direct SSH Command:** `gitmap login ssh [alias] [-y]`
  - Checks if `~/.ssh/id_ed25519` or `~/.ssh/id_ed25519_<alias>` exists.
  - If exists: reuses existing key unless `-y` or user confirms overwrite.
  - If not exists or regenerating: creates modern Ed25519 key.
  - Reads public key and copies to OS clipboard (`clip.exe` on Windows, `pbcopy` on macOS, `xclip`/`wl-copy` on Linux).
  - Displays public key in terminal with direct settings link:
    - GitHub: `https://github.com/settings/ssh/new`
    - GitLab: `https://gitlab.com/-/user_settings/ssh_keys`
  - Registers profile in `git_profiles.json`. If in a Git repo, configures `core.sshCommand`.

---

## 3. Invariants & Constraints
- Relative Git paths only. Zero absolute paths and zero `file:///` URIs.
- All modified and new Go files must remain strictly <= 300 lines with nesting depth <= 1.
- Zero heavy builds or test suites (`go build`, `npm run build`, `go test`, `pytest`).
- FastGate pre-commit validation must achieve 100% PASS.
