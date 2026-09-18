# Subtask 01: Core CLI Routers ErrorWrapper Refactor (env, group, profile, bookmark)

## 1. Objectives
Refactor intermediate routers in `cli/cmd/env.go`, `cli/cmd/group.go`, `cli/cmd/profile.go`, and `cli/cmd/bookmark.go` from `error` to `result.ErrorWrapper`, eliminating premature `.AsError()` unpacking and preventing Go typed nil interface hazards.

## 2. Disjoint Target Files
- `cli/cmd/env.go`
- `cli/cmd/group.go`
- `cli/cmd/profile.go`
- `cli/cmd/bookmark.go`

## 3. Implementation Details

### A. `cli/cmd/env.go`
1. `routeEnvSub(sub string, args []string) result.ErrorWrapper`:
   - Calls `resVar := routeEnvVariableSub(sub, args)`. If `resVar.IsMatched()`, returns `resVar` directly (NO `.AsError()`).
   - If `sub == constants.CmdEnvPathAdd`, returns `routeEnvPath(args)`.
   - Else returns `result.FailureWrapper(apperror.NewSimple(constants.ErrEnvSubcommand, "E9000"))`.
2. `routeEnvPath(args []string) result.ErrorWrapper`:
   - If `len(args) < 1`, returns `result.MatchWrapper(runEnvPathList())`.
   - Else returns `dispatchEnvPath(args[0], args[1:])`.
3. `dispatchEnvPath(sub string, rest []string) result.ErrorWrapper`:
   - Matches `CmdEnvPathSub`, `CmdEnvPathRemove`, `CmdEnvPathList` wrapped via `result.MatchWrapper`.
   - Default returns `result.FailureWrapper(apperror.NewSimple("constants.ErrEnvSubcommand "+"path "+sub, "E9000"))`.
4. `runEnv(args []string) error`:
   - Calls `res := routeEnvSub(sub, rest)`.
   - Returns `res.AsError()`.

### B. `cli/cmd/group.go`
1. `dispatchGroup(sub string, args []string) result.ErrorWrapper`:
   - Calls `resCRUD := dispatchGroupCRUD(sub, args)`. If `resCRUD.IsMatched()`, returns `resCRUD` directly.
   - For `CmdGroupShow`, returns `result.MatchWrapper(runGroupShow(args))`.
   - For `CmdGroupDelete`, returns `result.MatchWrapper(runGroupDelete(args))`.
   - Calls `resScoped := dispatchGroupScoped(sub, args)`. If `resScoped.IsMatched()`, returns `resScoped` directly.
   - Returns `result.MatchWrapper(activateGroup(sub))`.
2. `showActiveGroup() result.ErrorWrapper`:
   - Returns `result.FailureWrapper(...)` on DB error.
   - Returns `result.SuccessWrapper()` on success.
3. `runGroup(args []string) error`:
   - If `len(args) == 0`, returns `showActiveGroup().AsError()`.
   - Returns `dispatchGroup(args[0], args[1:]).AsError()`.

### C. `cli/cmd/profile.go`
1. `routeProfileSub(subCmd string, tailArgs []string) result.ErrorWrapper`:
   - Checks `routeGitProfileSub`, `routeDBProfileSub`, `routeChromeProfileSub`, `routeInstallProfileSub`.
   - If matched, returns the matched `result.ErrorWrapper` directly (NO `.AsError()`).
   - Default prints usage and returns `result.FailureWrapper(apperror.NewSimple("fatal error", "E9000"))`.
2. `runProfile(args []string) error`:
   - Returns `routeProfileSub(subCmd, tailArgs).AsError()`.

### D. `cli/cmd/bookmark.go`
1. `routeBookmarkSub(sub string, args []string) result.ErrorWrapper`:
   - For `CmdBookmarkSave`, returns `result.MatchWrapperAppErr(runBookmarkSave(args))`.
   - For `CmdBookmarkList`, returns `result.MatchWrapper(runBookmarkList(args))`.
   - For `CmdBookmarkRun`, returns `result.MatchWrapperAppErr(runBookmarkRun(args))`.
   - For `CmdBookmarkDelete`, returns `result.MatchWrapper(runBookmarkDelete(args))`.
   - Default prints usage and returns `result.FailureWrapper(apperror.NewSimple("fatal error", "E9000"))`.
2. `runBookmark(args []string) error`:
   - Returns `routeBookmarkSub(sub, rest).AsError()`.

## 4. Verification
- Functions strictly <= 15 lines.
- Affirmative booleans only.
- Strict Unix LF line endings.
- Targeted linters only.
