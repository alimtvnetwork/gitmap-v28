# Subtask 323-01: cmdundo Constants & Enums DRY Restructuring

**Status:** COMPLETED
**Completed At:** 2026-10-24T21:48:30+08:00

## Target Scope
- Target Package: `cli/cmdundo` (3 files: `undo.go`, `undo_flatten.go`, `hook_shim.go`)
- Leaf Packages Created:
  - `cli/cmdundo/enums/`
    - `undo_enums.go`
    - `enums_test.go`
  - `cli/cmdundo/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Implemented
1. `UndoModeType`: Execution mode (`restore`, `list`, `dry-run`, `flatten`).
2. `UndoActionType`: Restored entry status (`restored`, `[dry-run]`, `failed`).
3. `UndoExitCodeType`: Strongly-typed exit codes (`0` success, `1` general, `6` bad flag, `7` write failed, `8` bad config).

Each enum follows:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies (standard library only).

## Specific Constants Implemented
- Command names & aliases: `CmdUndo = "undo"`, `CmdUndoAlias = "ud"`, `CmdUndoFlatten = "undo-flatten"`, `CmdRepoUndoFlatten = "repo-undo-flatten"`, `CmdFixRepo = "fix-repo"`.
- Flags: `--list`, `-list`, `-l`, `--dry-run`, `-dry-run`, `-DryRun`, `--snapshot`, `-snapshot`, `-s`, `--help`, `-h`, `help`.
- Paths & suffixes: `GitMapDir = ".gitmap"`, `SubdirBackup = "backup"`, `SubdirFiles = "files"`, `ManifestName = "manifest.json"`, `SuffixMarkedForDel = "-marked-for-deletion"`, `SentinelCurrentDir = "."`.
- Markers: `MarkerLatest = "*"`, `MarkerOlder = " "`, `ModeWrite = "write"`, `ModeDryRun = "dry-run"`.
- Output messages & formats: `MsgNoSnapshotsFmt`, `MsgListHeaderFmt`, `MsgListRowFmt`, `MsgRestoreHeaderFmt`, `MsgRestoreRowFmt`, `MsgRestoreSummaryFmt`, `MsgRestoreErrFmt`, `MsgFlattenRestoringFmt`, `MsgFlattenRestoredSuccessFmt`, `MsgFlattenRestoredPathFmt`, `MsgFlattenRemovingTreeFmt`, `MsgFlattenUsage`.
- Error messages & operations: `ErrCodeE9000`, `ErrNoSnapshotFmt`, `ErrBadFlagFmt`, `ErrManifestMissingFmt`, `ErrManifestBadFmt`, `ErrArchiveNotFoundFmt`, `OpResolveRepoPathFmt`, `OpRestoreArchiveFmt`, `OpResolveCwd`, `OpRemoveExistingFlattenedFmt`.
- Exit codes: `ExitOk`, `ExitBadFlag`, `ExitWriteFailed`, `ExitBadConfig`.

## Consuming Callers Refactored
- `cli/cmdundo/undo.go`: Decoupled completely from root `cli/constants`. Converted magic strings to package-local `constants.*` and `enums.*`. Preserved error sovereignty: returned `error` or `*apperror.AppError` directly from leaf functions, removed all 5 `cliexit.HandleError` calls inside `restoreUndoSnapshot` and `parseUndoArgs`.
- `cli/cmdundo/undo_flatten.go`: Replaced magic strings with `constants.*`. Guard clauses for early returns, decomposed all functions into <= 15 lines.
- `cli/cmdundo/hook_shim.go`: Preserved clean DI hook.
- Added `cli/cmdundo/undo_test.go`: Comprehensive unit tests for command dispatch, arguments parsing, flag validation, snapshot listing/picking, and mock snapshot recovery.

## Verification
- Unit tests in `cli/cmdundo/enums/enums_test.go` and `cli/cmdundo/constants/constants_test.go` (100% PASS).
- Unit tests in `cli/cmdundo/undo_test.go` (100% PASS).
- `go vet ./cmdundo/...` (PASS, 0 issues).
- `go test -count=1 -v ./cmdundo/...` (PASS, 100% pass across all packages).
