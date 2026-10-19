# Subtask 01: Pull Caching (3s TTL), Redundant Notice Removal, Inline Untracked/Modified Display & Settings UI

- **Subtask ID**: `subtask-01-pull-caching-and-untracked-display`
- **Parent Task**: `gitmap-efficiency-and-coding-guidelines-enhancement`
- **Owned Files**:
  - `cli/cmdpull/pull_dispatch_routes.go`
  - `cli/cmdpull/pull_cwd.go`
  - `cli/cmdpull/pull_summary_single.go`
  - `cli/constants/constants_settings.go`
  - `src/pages/Settings.tsx`
  - `src/components/settings/PullCacheCard.tsx`
- **Assigned Worker**: Worker 01

---

## 1. Objectives
1. Remove redundant CWD print statement `-> cwd is a git repo — running plain git pull here` from `executePullCWDWithNotice` in `cli/cmdpull/pull_dispatch_routes.go`.
2. Define setting constant `SettingPullCacheWindowSeconds = "pull_cache_window_seconds"` and `DefaultPullCacheWindowSeconds = 3` in `cli/constants/constants_settings.go`.
3. In `cli/cmdpull/pull_cwd.go`:
   - Query last pull timestamp from `gitmap-pull.db` (`store.GetRecentRepoPullHistory(repoPath, 1)`).
   - If previous pull completed $\le 3$s ago (or configured TTL):
     - Inspect local dirty state via `gitutil.InspectDirtyState(repoPath)`.
     - If remote fetch is not required within window and dirty state is consistent, output cached status:
       `⚡ gitmap: cached (checked %ds ago) - up-to-date (%s) [%s] (dirty: ...)`
4. In `cli/cmdpull/pull_summary_single.go`:
   - When dirty files exist (`diag.IsDirty`), render modified and untracked files directly beneath the single-repo summary line:
     ```text
       M  .gitignore
       M  cli/cmd/roottooling.go
       ?? cli/cmdlemon/daemon_settings.go
     ```
     eliminating the requirement for the user to run a follow-up `git status -s`.
5. In `src/pages/Settings.tsx` & `src/components/settings/PullCacheCard.tsx`:
   - Add a configuration card allowing users to inspect and adjust the Git pull cache window (default 3 seconds, min 0, max 60).

---

## 2. Verification
- `go test ./cmdpull -v`
- Inspect output formatting matches screenshot 01 and 02 requirements.
