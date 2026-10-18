# Subtask 02: Multi-Login Profile Fallback (Default First, Parallel Secondary)

- **Task-ID**: Task-02
- **Subtask**: 02-multi-login-profile-fallback
- **Owner**: Worker 02
- **Status**: PENDING

## Owned Files (Exclusive)
- `cli/cmdclone/clone_auth_resolve.go`
- `cli/cmdclone/gh_resolver.go`
- `cli/cmdclone/clone.go`
- `cli/cmdclone/flags.go`
- `cli/cmdclone/clone_profile_test.go`

## Read-Only Context Files
- `02-spec/21-app/goemailsender-connection-and-testing/02-component-spec.md`
- `cli/store/git_profile_store.go`
- `cli/usercontext/project_binder.go`
- `cli/model/git_profile.go`

## Objectives
1. Add `--profile` / `-p` flag support to `CloneFlags` and propagation in `cmdclone`.
2. Implement fallback hierarchy:
   - Primary profile probed first.
   - On 401/404/auth failure, spawn concurrent lightweight probes (`git ls-remote --heads`) across all secondary profiles in `02-git-profiles.json`.
   - Atomic winner selection and automatic project binding via `usercontext.BindProject`.
   - Unified diagnostic error if all secondary profiles fail.
3. Unit test verifying primary success, secondary fallback on 404, and all-fail diagnostics in `cli/cmdclone/clone_profile_test.go`.
