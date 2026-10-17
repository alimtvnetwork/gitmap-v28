# Subtask 02: GCM Modal Suppression, Global Git Identity Enforcement, and Unified `gitmap login` Subsystem

- **Task ID:** subtask-login-clone-02
- **Parent Task:** gitmap-login-and-cloning-enhancements
- **Status:** DONE
- **Owner:** Worker 02

## Objectives
1. In `cli/secrets/ghtoken.go`:
   - In `tokenFromGitCredential()`: add `GCM_INTERACTIVE=never` to `cmd.Env` and pass `-c credential.interactive=never` in args so GCM never spawns a GUI modal on Windows.
2. In `cli/cmd/rootusagecompact.go`:
   - In `detectActiveGitUser()`: strictly query `git config --global user.name` and `git config --global user.email`. If either is unset, return empty and display `✖ Not configured / Not logged in`.
   - In `printCompactGuidanceFooter()`: display login guidance when Git is not configured or GitHub CLI is not logged in.
3. In `cli/store/git_profile_store.go` and `cli/usercontext/gh_auth.go`:
   - Remove hardcoded `"default"` and `"github-user"` placeholders.
4. In `cli/model/git_profile.go`:
   - Add `SSHKeyPath string` and `SSHKeyName string` to `GitProfile`.
5. In `cli/usercontext/project_binder.go`:
   - When binding a project to an account or syncing, set `core.sshCommand = "ssh -i <key> -F /dev/null -o IdentitiesOnly=yes"`.
6. In `cli/cmdlogin/`:
   - Implement `gitmap login` interactive flow (SSH, Token, Browser flow, Status).
   - Implement `gitmap login ssh [alias] [-y]` (in `login_ssh.go`, <= 300 lines):
     - Check/generate SSH key (`id_ed25519` or `id_ed25519_<alias>`).
     - Copy public key to OS clipboard (`clip.exe` / `pbcopy` / `xclip`).
     - Print public key and direct settings URLs (GitHub/GitLab).
     - Register profile in `git_profiles.json` and configure `core.sshCommand` if inside a git repository.
   - Support `gitmap login browser-flow` and `gitmap login token [PAT]`.
7. Author `cli/helpdoc/login.md`.
8. Verify all modified Go files remain strictly <= 300 lines with nesting depth <= 1.
