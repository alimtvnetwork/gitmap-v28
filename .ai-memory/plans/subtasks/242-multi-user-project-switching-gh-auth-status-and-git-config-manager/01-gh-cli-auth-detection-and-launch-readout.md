# Subtask 01: GitHub CLI Authorization Status Detection & Launch Readout

- **Parent Task:** `242-multi-user-project-switching-gh-auth-status-and-git-config-manager`
- **Subtask Code:** `Task-01`
- **Assigned Worker:** Worker 01
- **Status:** COMPLETED

---

## Scope & Implementation Details
1. Implement `cli/usercontext/types.go` declaring:
   - `AuthStatus` enum (`StatusAuthorized`, `StatusUnauthenticated`, `StatusToolMissing`).
   - `GhAuthInfo struct { Username string, Status AuthStatus, Source string, Scope string }`.
2. Implement `cli/usercontext/gh_auth.go`:
   - `DetectGhAuth() GhAuthInfo`:
     Probes `gh auth status` or `gh api user --jq .login`.
     Falls back to `ghtoken.Resolve()` if token is stored in Git Credential Manager or env.
     Extracts authenticated GitHub login username in non-blocking manner ($\le 15$ms).
3. Update `cli/cmd/rootusagecompact.go` and `cli/cmd/binarylocations.go`:
   - On bare `gitmap` launch, print clean Catppuccin block:
     `● GitHub CLI:     ✔ Authorized as <username>` or `✖ Not logged in`
     `● Active Git:     <Name> <<email>>`
4. Add unit test `cli/usercontext/gh_auth_test.go`.

---

## Acceptance Criteria
- [ ] `DetectGhAuth()` correctly returns `GhAuthInfo` without crashing if `gh` is missing.
- [ ] Bare `gitmap` readout displays GitHub CLI status and active Git user.
- [ ] Unit tests pass in <50ms.
