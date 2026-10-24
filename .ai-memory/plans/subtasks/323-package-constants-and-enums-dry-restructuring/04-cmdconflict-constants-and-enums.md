# Subtask 323-04: cmdconflict Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdconflict` (3 files: `conflicts_cmd.go`, `ide_launcher.go`, `resolve_cmd.go`)
- Leaf Packages to Create:
  - `cli/cmdconflict/enums/`
    - `conflict_enums.go`
    - `enums_test.go`
  - `cli/cmdconflict/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `ConflictActionType`: Resolution action (`theirs`, `ours`, `ide`, `abort`, `rerere`).
2. `ConflictToolType`: IDE tools (`code`, `cursor`, `windsurf`, `codium`).
3. `ConflictStatusCodeType`: Unmerged status codes (`UU`, `AA`, `UD`, `DU`, `DD`, `AU`, `UA`).
4. `ConflictScopeType`: Inspection scope (`target`, `current`, `all`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies.

## Specific Constants Required
- Command names & aliases: `CmdConflicts = "conflicts"`, `CmdConflictsAlias = "conflict"`, `CmdResolve = "resolve"`, `CmdResolveAlias = "res"`.
- Flags: `--theirs`, `-t`, `--ours`, `-o`, `--ide`, `-i`, `--code`, `--abort`, `-a`, `--rerere`, `--help`, `-h`.
- Strategy tokens: `TokenTheirs`, `TokenOurs`, `TokenIDE`, `TokenCode`, `TokenResolve`, `TokenAbort`, `TokenRerere`.
- Git CLI: `GitBin = "git"`, `GitFlagDir = "-C"`, `GitSubCmdDiff = "diff"`, `GitArgNameOnly = "--name-only"`, `GitDiffFilterUnmerged = "--diff-filter=U"`, `GitSubCmdStatus = "status"`, `GitArgPorcelain = "--porcelain"`, `GitSubCmdCheckout = "checkout"`, `GitArgTheirs = "--theirs"`, `GitArgOurs = "--ours"`, `GitSubCmdMerge = "merge"`, `GitSubCmdRebase = "rebase"`, `GitSubCmdCherryPick = "cherry-pick"`, `GitArgAbort = "--abort"`, `GitConfigRerereEnabled = "rerere.enabled"`.
- Colors & glyphs: `ColorGreen`, `ColorYellow`, `ColorReset`, `GlyphCheck = "✓"`, `GlyphWarn = "⚠"`, `TreePrefixBranch = "├──"`, `TreePrefixLeaf = "└──"`.
- Help & guides: `HelpSectionQuickActions`, `HelpActionIDE`, `HelpActionTheirs`, `HelpActionOurs`, `HelpActionAbort`, `GuideDivider`, `GuideTitle`, `GuideStep1` .. `GuideStep6`.
- Messages & errors: `MsgNoConflictsDetectedFmt`, `MsgAcceptedTheirsFmt`, `MsgAcceptedOursFmt`, `MsgNextCommitInstruction`, `MsgIDELaunchedFmt`, `MsgMergeAbortedFmt`, `ErrNoSupportedIDE`, `ErrNoAbortInProgressFmt`, `ErrResolveUsage`.

## Consuming Callers to Refactor
- `cli/cmdconflict/conflicts_cmd.go`: Decouple completely from monolithic `cli/constants`. Replace magic status strings and porcelain codes with `enums.ConflictStatusCodeType` and `constants.*`.
- `cli/cmdconflict/ide_launcher.go`: Use `enums.ConflictToolType` (`BinaryName()`, `DisplayName()`) and `constants.*`.
- `cli/cmdconflict/resolve_cmd.go`: Use `enums.ConflictActionType` and `constants.*`.
- Add `cli/cmdconflict/conflict_test.go`: Unit tests covering status code parsing, IDE tool lookup, resolution action parsing, and quick action help formatting.

## Verification
- Unit tests in `cli/cmdconflict/enums/enums_test.go` and `cli/cmdconflict/constants/constants_test.go`.
- Unit tests in `cli/cmdconflict/conflict_test.go`.
- `go vet ./cli/cmdconflict/...`
- `go test -count=1 ./cli/cmdconflict/...`

## Status
- State: COMPLETED
- Verification: Passed (`go vet -C cli ./cmdconflict/...`, `go test -C cli -count=1 -v ./cmdconflict/...`)
- Artifacts:
  - `cli/cmdconflict/enums/conflict_enums.go`
  - `cli/cmdconflict/enums/enums_test.go`
  - `cli/cmdconflict/constants/constants.go`
  - `cli/cmdconflict/constants/constants_test.go`
  - `cli/cmdconflict/conflicts_cmd.go`
  - `cli/cmdconflict/ide_launcher.go`
  - `cli/cmdconflict/resolve_cmd.go`
  - `cli/cmdconflict/conflict_test.go`

