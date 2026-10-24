# Subtask 324-02: cmdlfs Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdlfs` (2 files: `lfs_convert.go`, `lfs_exec.go`)
- Leaf Packages to Create:
  - `cli/cmdlfs/enums/`
    - `lfs_enums.go`
    - `enums_test.go`
  - `cli/cmdlfs/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `LFSStatusType`: Candidate status (`pending add`, `already tracked`, `converted`, `failed`).
2. `LFSModeType`: Execution mode (`convert`, `dry-run`).
3. `LFSCommandType`: Subcommand tokens (`convert`, `conv`).
4. `LFSUnitType`: Size units (`GB`, `G`, `MB`, `M`, `KB`, `K`, `B`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies.

## Specific Constants Required
- Command names & aliases: `CmdLFSConvert = "lfs-convert"`, `CmdLFSConvertAlias = "lfsconvert"`, `SubcmdConvert = "convert"`, `SubcmdConv = "conv"`.
- Flags: `--dry-run`, `-y`, `--yes`, `--force`, `--size`, `-s`, `--help`, `-h`.
- Threshold defaults: `DefaultSizeRaw = "5MB"`, `DefaultThresholdBytes = 5242880`, `UnitsScaleChars = "KMGTPE"`.
- Files & attributes: `FileGitAttributes = ".gitattributes"`, `FileModeGitAttributes = 0o644`, `CommentPrefix = "#"`, `LFSFilterToken = "filter=lfs"`, `LFSAttributeRule = "filter=lfs diff=lfs merge=lfs -text"`, `LFSAttributeRuleFmt = "%s filter=lfs diff=lfs merge=lfs -text\n"`.
- Git CLI: `BinGit = "git"`, `GitCmdRevParse = "rev-parse"`, `GitFlagInsideWorkTree = "--is-inside-work-tree"`, `GitValTrue = "true"`, `GitCmdLFS = "lfs"`, `GitCmdVersion = "version"`, `GitCmdLsFiles = "ls-files"`, `GitCmdInstall = "install"`, `GitFlagLocal = "--local"`, `GitCmdMigrate = "migrate"`, `GitCmdImport = "import"`, `GitFlagEverything = "--everything"`, `GitFlagIncludePrefix = "--include="`, `GitCmdAdd = "add"`, `GitCmdRm = "rm"`, `GitFlagCached = "--cached"`.
- ANSI colors & table headers: `ColHeaderFile`, `ColHeaderSize`, `ColHeaderLFSRule`, `ColHeaderStatus`, `ColorReset`, `ColorGreen`, `ColorYellow`, `ColorCyan`, `ColorWhite`.
- Messages & error codes: `MsgNoCandidatesFmt`, `MsgDryRunHeaderFmt`, `MsgSummaryFmt`, `MsgDryRunNotice`, `PromptConfirmWarning`, `PromptConfirmQuestion`, `MsgSuccessSummaryFmt`, `DomainLFSConvert = "lfs_convert"`, `ErrCodeE9001`, `ErrCodeE9002`, `ErrMsgNotInsideGitRepo`, `ErrMsgLFSNotInstalled`.

## Consuming Callers to Refactor
- `cli/cmdlfs/lfs_convert.go`:
  - Decouple completely from monolithic `cli/constants`.
  - Fix boolean field naming: `dryRun` -> `isDryRun`, `autoConfirm` -> `isAutoConfirm`, `force` -> `isForce`.
  - Replace subcommands using `enums.ParseLFSCommandType`.
  - Replace byte size parsing to use `enums.ResolveSizeMultiplier`.
- `cli/cmdlfs/lfs_exec.go`:
  - Fix boolean naming: `AlreadyInAttr` -> `isAlreadyInAttr`.
  - Replace candidate status strings with `enums.LFSStatusPendingAdd.String()`, etc.
  - Replace table columns, ANSI colors, and Git commands with `constants.*`.
  - Replace error codes with `constants.ErrCodeE9001`, `constants.ErrCodeE9002`, `constants.Op*`.
- Add `cli/cmdlfs/lfs_test.go`: Unit tests for size parsing, flag validation, command normalization, and dry-run table rendering.

## Verification
- Unit tests in `cli/cmdlfs/enums/enums_test.go` and `cli/cmdlfs/constants/constants_test.go`.
- Unit tests in `cli/cmdlfs/lfs_test.go`.
- `go vet ./cli/cmdlfs/...`
- `go test -count=1 ./cli/cmdlfs/...`
