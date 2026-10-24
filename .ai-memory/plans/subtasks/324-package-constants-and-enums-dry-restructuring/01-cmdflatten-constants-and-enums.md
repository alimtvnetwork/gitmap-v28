# Subtask 324-01: cmdflatten Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdflatten` (2 files: `flatten_cmd.go`, `flatten_copy.go`)
- Leaf Packages to Create:
  - `cli/cmdflatten/enums/`
    - `flatten_enums.go`
    - `enums_test.go`
  - `cli/cmdflatten/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `FlattenModeType`: Execution mode (`execute`, `dry-run`).
2. `FlattenStepType`: Pipeline stages (`validate`, `stage`, `copy`, `init`, `commit`, `archive`, `promote`).
3. `FlattenExcludedDirType`: Default excluded directories (`.git`, `node_modules`, `dist`, `temp`, `.pnpm-store`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies.

## Specific Constants Required
- Command names & aliases: `CmdFlatten = "flatten"`, `CmdRepoFlatten = "repo-flatten"`, `CmdUndoFlatten = "undo-flatten"`.
- Flags: `--dry-run`, `-dry-run`, `-y`, `--yes`, `-yes`, `--target`, `-target`, `--target=`, `--help`, `-h`.
- Directory names & patterns: `DefaultTargetDir = "."`, `DirGit = ".git"`, `SuffixMarkedForDeletion = "-marked-for-deletion"`, `SuccessorNamePattern = "%s-flat-v%d"`, `DefaultDirPerm = 0755`.
- Git tokens: `BinGit = "git"`, `GitFlagDir = "-C"`, `GitCmdRevList = "rev-list"`, `GitFlagCount = "--count"`, `GitRefHead = "HEAD"`, `GitCmdInit = "init"`, `GitBranchMain = "main"`, `GitCmdAdd = "add"`, `GitFlagAll = "-A"`, `GitCmdCommit = "commit"`, `GitDefaultCommitMsg = "flat: initial squashed repository snapshot"`.
- Units & colors: `ByteUnitDivisor = 1024`, `ByteUnitB`, `ByteUnitKB`, `ByteUnitMB`, `ByteUnitGB`, `ByteUnitTB`, `ColorReset`, `ColorBold`, `ColorGreen`, `ColorYellow`, `ColorCyan`.
- Messages & prompts: `PromptConfirmWarning`, `PromptConfirmQuestion`, `MsgFlattenAborted`, `MsgCopyingWorkingTreeFmt`, `MsgArchivingBloatedRepoFmt`, `MsgPromotingSuccessorFmt`, `MsgDryRunHeaderFmt`, `MsgSuccessSummaryFmt`, `MsgUndoHintFmt`.
- Operations & error messages: `ErrTargetNotDirFmt`, `ErrTargetNotGitRepoFmt`, `ErrMissingTargetFlagValue`, `ErrUserAborted`, `OpResolveTargetPathFmt`, `OpCalcWorkingTreeSize`, `OpCreateSuccessorDirFmt`, `OpCopyWorkingTree`, `OpInitGitRepo`, `OpCommitFlatSnapshot`, `OpArchiveOriginalRepo`, `OpPromoteSuccessorDir`.

## Consuming Callers to Refactor
- `cli/cmdflatten/flatten_cmd.go`: Decouple completely from monolithic `cli/constants`. Replace raw flag strings, colors, prompts, messages, and error operations with `constants.*` and `enums.*`. Flatten nested conditionals.
- `cli/cmdflatten/flatten_copy.go`: Replace Git command tokens, excluded dir checks (`enums.IsExcludedDir`), and byte units with `constants.*` and `enums.*`.
- Add `cli/cmdflatten/flatten_test.go`: Unit tests for command argument parsing, dry-run evaluation, confirmation logic, and excluded directory filtering.

## Verification
- Unit tests in `cli/cmdflatten/enums/enums_test.go` and `cli/cmdflatten/constants/constants_test.go`.
- Unit tests in `cli/cmdflatten/flatten_test.go`.
- `go vet ./cli/cmdflatten/...`
- `go test -count=1 ./cli/cmdflatten/...`
