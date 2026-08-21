# Subtask 05: Go Wrapped Result Structs & Single Return Pattern

Slug: 05-wrapped-booleans-and-single-return
Parent Plan: 01-coding-guideline-fixes
Status: pending

## Objective
Enforce the Go wrapped result return convention: functions returning status should return a single generic/typed wrapped `Result` struct containing mutually exclusive `IsSuccess` and `IsFailed` flags along with payload/error, avoiding bare boolean return signatures.

## Concrete Execution Steps (20 Steps)

1. `gitmap/cloner/cloner.go:141`: Refactor `hasExistingRepos(records, targetDir) bool` to return `ExistingReposResult{IsFound: bool, IsMissing: bool, FoundPaths: []string}`.
2. `gitmap/cloner/cache.go:45`: Refactor cache check `hasCacheEntry` to return `CacheLookupResult{IsHit: bool, IsMiss: bool, Entry: CloneCacheEntry}`.
3. `gitmap/cloner/strategy.go:30`: Refactor `canFastForward` to return `FastForwardResult{CanFastForward: bool, IsBlocked: bool, Reason: string}`.
4. `gitmap/scanner/scanner_sniff.go`: Refactor `isGitDir` to return `GitSniffResult{IsGitDir: bool, IsRegularFile: bool, IsWorktree: bool}`.
5. `gitmap/scanner/sort.go:25`: Refactor sort predicates to use typed comparison result structs.
6. `gitmap/store/lock.go:40`: Refactor `isLockActive` to return `LockStatusResult{IsLocked: bool, IsUnlocked: bool, OwnerPID: int}`.
7. `gitmap/store/pendingtask.go:55`: Refactor `hasPendingTask` to return `PendingTaskLookupResult{IsFound: bool, IsEmpty: bool, Task: PendingTask}`.
8. `gitmap/store/installedtool.go:35`: Refactor `isToolInstalled` to return `ToolInstalledResult{IsInstalled: bool, IsMissing: bool, ToolPath: string}`.
9. `gitmap/cmd/installdetect.go:40`: Refactor `isCommandAvailable` to return `CommandAvailabilityResult{IsAvailable: bool, IsNotFound: bool, ResolvedPath: string}`.
10. `gitmap/cmd/sshgen.go:60`: Refactor `hasSSHKey` to return `SSHKeyStatusResult{HasKey: bool, IsMissing: bool, KeyPath: string}`.
11. `gitmap/cmd/visibilityundo.go:45`: Refactor `canUndoVisibility` to return `VisibilityUndoResult{CanUndo: bool, IsBlocked: bool, SnapshotID: int64}`.
12. `gitmap/cmd/reclone_validate.go:50`: Refactor validation routines to return `ValidationStatusResult{IsValid: bool, IsInvalid: bool, ValidationErrors: []error}`.
13. `gitmap/completion/cdfunction.go:40`: Refactor `isShellActive` to return `ShellStatusResult{IsActive: bool, IsInactive: bool, ShellName: string}`.
14. `gitmap/config/validate_shape.go:65`: Refactor shape validation to return `ShapeValidationResult{IsValid: bool, IsInvalid: bool, SchemaIssues: []string}`.
15. `gitmap/gitutil/gitutil.go:50`: Refactor `isCleanWorkingTree` to return `WorkingTreeResult{IsClean: bool, IsDirty: bool, ModifiedCount: int}`.
16. `gitmap/gitutil/gitutil.go:80`: Refactor `isBranchMerged` to return `BranchMergeResult{IsMerged: bool, IsUnmerged: bool, TargetBranch: string}`.
17. `gitmap/release/semver.go:45`: Refactor `isValidSemver` to return `SemverValidationResult{IsValid: bool, IsInvalid: bool, ParsedVersion: Semver}`.
18. `gitmap/release/assets.go:70`: Refactor `hasAsset` to return `AssetLookupResult{HasAsset: bool, IsMissing: bool, AssetURL: string}`.
19. `gitmap/startup/winregistry_remove_windows.go:30`: Refactor registry key verification to return `RegistryKeyResult{Exists: bool, IsAbsent: bool}`.
20. `gitmap/vscodeworkspace/vscodeworkspace.go:40`: Refactor `isWorkspaceRegistered` to return `WorkspaceStatusResult{IsRegistered: bool, IsUnregistered: bool, ConfigPath: string}`.

## Target Verification Files
- `gitmap/cloner/*`
- `gitmap/scanner/*`
- `gitmap/store/*`
- `gitmap/cmd/*`
