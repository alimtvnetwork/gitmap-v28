# Subtask 323-03: cmdstale Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdstale` (3 files: `stale.go`, `helpers_shim.go`, `hook_shim.go`)
- Leaf Packages to Create:
  - `cli/cmdstale/enums/`
    - `stale_enums.go`
    - `enums_test.go`
  - `cli/cmdstale/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `StaleOutputFormatType`: Output format (`table`, `json`, `csv`).
2. `StaleActionType`: Action mode (`inspect`, `archive`, `preview`).
3. `StaleFilterType`: Filtering criteria (`all`, `stale`, `active`, `archived`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies.

## Specific Constants Required
- Command names & aliases: `CmdStale = "stale"`, `CmdStaleAlias = "sta"`, `SubCmdHelp = "help"`.
- Flags: `--days`, `--root`, `--archive`, `--dry-run`, `--format`, `--help`, `-h`.
- Default thresholds: `DefaultStaleDays = 90`, `DefaultScanRoot = "."`, `DefaultOutputFormat = "table"`, `MinHygieneWorkers = 2`, `MaxHygieneWorkers = 8`, `HoursPerDay = 24`, `DefaultArchiveDirPerm = 0o755`.
- Directory paths: `DirDotGit = ".git"`, `DirGitmap = ".gitmap"`, `DirArchive = "archive"`.
- Date/time formats: `DateFormatStale = "2006-01-02"`, `TimestampFormatArchive = "20060102T150405Z"`, `JSONIndent = "  "`.
- CSV headers: `HeaderPath = "path"`, `HeaderLastCommitUTC = "last_commit_utc"`, `HeaderAgeDays = "age_days"`.
- Git CLI: `GitBin = "git"`, `GitFlagDir = "-C"`, `GitSubCmdLog = "log"`, `GitArgLimitOne = "-1"`, `GitFormatCommitEpoch = "--format=%ct"`.
- ANSI colors: `ColorCyan`, `ColorYellow`, `ColorGreen`, `ColorRed`, `ColorReset`.
- Messages & formats: `MsgNoStaleReposFmt`, `MsgStaleHeaderFmt`, `MsgStaleRowFmt`, `MsgWouldArchiveFmt`, `MsgArchivedFmt`, `MsgFailedFmt`, `ErrJSONEncodeFmt`, `ErrArchiveMkdirFmt`, `ErrInvalidFormatFmt`.

## Consuming Callers to Refactor
- `cli/cmdstale/stale.go`:
  - Enforce error return sovereignty: replace `cliexit.HandleError(err, 2)` and `cliexit.Fail(...)` with direct returns of `*appfault.AppError` or `error`.
  - Replace magic strings and formats with `constants.*` and `enums.*`.
  - Invert conditionals to flatten nested blocks.
- `cli/cmdstale/helpers_shim.go`:
  - Adopt `constants.MinHygieneWorkers`, `constants.MaxHygieneWorkers`, CSV headers, and formatting tokens.
- Add `cli/cmdstale/stale_test.go`: Unit tests for stale argument parsing, format resolution, dry-run simulation, and JSON/CSV emission.

## Verification
- Unit tests in `cli/cmdstale/enums/enums_test.go` and `cli/cmdstale/constants/constants_test.go`.
- Unit tests in `cli/cmdstale/stale_test.go`.
- `go vet ./cli/cmdstale/...`
- `go test -count=1 ./cli/cmdstale/...`

## Status
- State: COMPLETED
- Verification: Passed (`go vet ./cmdstale/...`, `go test -count=1 -v ./cmdstale/...`)
- Artifacts:
  - `cli/cmdstale/enums/stale_enums.go`
  - `cli/cmdstale/enums/enums_test.go`
  - `cli/cmdstale/constants/constants.go`
  - `cli/cmdstale/constants/constants_test.go`
  - `cli/cmdstale/stale.go`
  - `cli/cmdstale/helpers_shim.go`
  - `cli/cmdstale/hook_shim.go`
  - `cli/cmdstale/stale_test.go`

