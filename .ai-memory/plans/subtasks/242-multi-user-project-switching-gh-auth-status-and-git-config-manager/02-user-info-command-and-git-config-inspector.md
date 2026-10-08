# Subtask 02: User Info Command & Git Configuration Inspector

- **Parent Task:** `242-multi-user-project-switching-gh-auth-status-and-git-config-manager`
- **Subtask Code:** `Task-02`
- **Assigned Worker:** Worker 02
- **Status:** PENDING

---

## Scope & Implementation Details
1. Implement `cli/usercontext/git_identity.go`:
   - `GetGitIdentity(isGlobal bool) (name, email string, isSet bool)`:
     Reads `git config user.name` and `user.email`.
   - `SetGitIdentity(name, email string, isGlobal bool) error`:
     Configures git `user.name` and `user.email`.
2. Implement `cli/usercontext/summary.go`:
   - `CollectUserSummary() UserSummary`:
     Assembles GitHub CLI auth status, local git user, global git user, active GitMap profile, and bound project user into a consolidated model.
3. Implement `cli/cmd/user_info.go`:
   - Handles `gitmap user info` and `gitmap user-info`.
   - Supports `--json` flag emitting structured JSON.
   - Renders Catppuccin terminal card showing:
     - GitHub CLI authorization status and username.
     - Active Git Identity (Local repository if present, Global fallback).
     - Global Git Identity.
     - Active GitMap Profile.
     - Current project binding.
4. Unit tests in `cli/usercontext/git_identity_test.go`.

---

## Acceptance Criteria
- [ ] `gitmap user info` runs cleanly and displays public identity without security leaks.
- [ ] `--json` produces valid parseable JSON.
- [ ] Functions satisfy $\le 15$ line caps and positive booleans.
