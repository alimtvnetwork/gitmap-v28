# Subtask 03: Multi-User Profile Store & Per-Project Auto-Switching Engine

- **Parent Task:** `242-multi-user-project-switching-gh-auth-status-and-git-config-manager`
- **Subtask Code:** `Task-03`
- **Assigned Worker:** Worker 01
- **Status:** COMPLETED

---

## Scope & Implementation Details
1. Extend `cli/model/git_profile.go`:
   - Add `ProjectBinding struct { RepoPath string, ProfileAlias string, BoundAt time.Time }`.
   - Add `ProjectBindings map[string]ProjectBinding` to `GitProfileConfig`.
2. Implement `cli/usercontext/project_binder.go`:
   - `ResolveProjectUser(repoDir string) (*model.GitProfile, bool)`:
     Normalizes repository directory path, looks up `ProjectBindings`, returns matching `GitProfile`.
   - `BindProject(repoDir, profileAlias string) error`:
     Associates repository directory with profile alias, writes to `git_profiles.json`.
   - `UnbindProject(repoDir string) error`:
     Clears repository binding.
   - `SyncProjectUser(repoDir string) error`:
     If repository is bound to a profile, automatically ensures local repository `git config user.name` and `user.email` match the profile!
3. Add unit test `cli/usercontext/project_binder_test.go`.

---

## Acceptance Criteria
- [ ] Project bindings correctly persist in `git_profiles.json`.
- [ ] `ResolveProjectUser` resolves exact and normalized directory matches.
- [ ] `SyncProjectUser` idempotently synchronizes local repository git config.
