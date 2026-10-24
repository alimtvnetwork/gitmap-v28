# Subtask 324-03: cmdgitrm Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdgitrm` (2 files: `gitrm_entry.go`, `hook_shim.go`)
- Leaf Packages to Create:
  - `cli/cmdgitrm/enums/`
    - `gitrm_enums.go`
    - `enums_test.go`
  - `cli/cmdgitrm/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `GitRmInputFormatType`: Input format (`json`, `csv`, `text`, `list`, `single`).
2. `GitRmActionType`: Lifecycle stage (`parse`, `backup`, `rewrite`, `cleanup`).
3. `GitRmExitCodeType`: Exit codes (`0` success, `1` general, `2` usage).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies.

## Specific Constants Required
- Command names & aliases: `CmdGitRm = "git-rm"`, `CmdGitRmAlias = "rm"`, `SubCmdHelp = "help"`.
- Flags: `--help`, `-h`, `FlagHelpEntryHelp = "-h, --help"`.
- Directories, paths & refs: `DirGitmap = ".gitmap"`, `DirBackups = "backups"`, `DirGitRm = "git-rm"`, `RefsOriginal = "refs/original/"`.
- File extensions & delimiters: `ExtJSON = ".json"`, `ExtCSV = ".csv"`, `ExtTXT = ".txt"`, `DelimiterComma = ","`, `DelimiterSpace = " "`, `QuoteSingle = "'"`, `Newline = "\n"`.
- Permissions: `DefaultDirPerm = 0o755`, `DefaultFilePerm = 0o644`.
- Git CLI: `GitBin = "git"`, `GitSubCmdFilterBranch = "filter-branch"`, `GitFlagForce = "--force"`, `GitFlagIndexFilter = "--index-filter"`, `GitFilterCmdFmt = "git rm --cached --ignore-unmatch %s"`, `GitFlagPruneEmpty = "--prune-empty"`, `GitFlagTagNameFilter = "--tag-name-filter"`, `GitArgCat = "cat"`, `GitSubCmdForEachRef = "for-each-ref"`, `GitSubCmdReflog = "reflog"`, `GitSubCmdGC = "gc"`, `GitFlagPruneNow = "--prune=now"`, `GitFlagAggressive = "--aggressive"`.
- Error codes & messages: `ErrScopeGitRmRun = "GitRmRun"`, `ErrCodeMissingInput = "E_GITRM_MISSING_INPUT"`, `ErrCodeNoPaths = "E_GITRM_NO_PATHS"`, `ErrWrapPrefix = "Error:"`, `ErrParseInput`, `ErrCreateBackupDir`, `ErrHistoryRewrite`.
- Status messages: `MsgRemovingFilesFmt = "Removing %d files from git history...\n"`, `MsgSuccessRewritten = "Git history successfully rewritten."`.

## Consuming Callers to Refactor
- `cli/cmdgitrm/gitrm_entry.go`:
  - Decouple completely from monolithic `github.com/alimtvnetwork/gitmap-v28/cli/constants`.
  - Replace `constants.CmdGitRm` with local `constants.CmdGitRm`.
  - Replace magic string `"Error:"` with `constants.ErrWrapPrefix`.
  - Maintain error return sovereignty.
- `cli/cmdgitrm/hook_shim.go`: Clean DI hook.
- Add `cli/cmdgitrm/gitrm_test.go`: Unit tests for `RunGitRm` help execution, argument forwarding, and error wrapping.

## Verification
- Unit tests in `cli/cmdgitrm/enums/enums_test.go` and `cli/cmdgitrm/constants/constants_test.go`.
- Unit tests in `cli/cmdgitrm/gitrm_test.go`.
- `go vet ./cli/cmdgitrm/...`
- `go test -count=1 ./cli/cmdgitrm/...`
