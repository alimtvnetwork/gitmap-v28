# Subtask 324-04: cmdsize Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdsize` (3 files: `size.go`, `helpers_shim.go`, `hook_shim.go`)
- Leaf Packages to Create:
  - `cli/cmdsize/enums/`
    - `size_enums.go`
    - `enums_test.go`
  - `cli/cmdsize/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `SizeOutputFormatType`: Output layout (`table`, `json`, `csv`).
2. `SizeActionType`: Operational mode (`report`, `prune`, `dry-run`).
3. `SizeExitCodeType`: Exit codes (`0` success, `1` general, `2` usage).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies.

## Specific Constants Required
- Command names & aliases: `CmdSize = "size"`, `CmdSizeAlias = "sz"`, `SubCmdHelp = "help"`.
- Flags: `--root`, `--top`, `--prune`, `--dry-run`, `--format`, `--help`, `-h`.
- Defaults & limits: `DefaultRoot = "."`, `DefaultTopN = 0`, `DefaultPrune = false`, `DefaultDryRun = false`, `DefaultFormat = "table"`, `MinHygieneWorkers = 2`, `MaxHygieneWorkers = 8`, `ByteUnitScale = 1024.0`.
- Formats & keys: `FormatTable = "table"`, `FormatJSON = "json"`, `FormatCSV = "csv"`, `JSONIndent = "  "`, `KeyPath = "path"`, `KeyBytes = "bytes"`.
- Units: `UnitBytes = "B"`, `UnitKB = "KB"`, `UnitMB = "MB"`, `UnitGB = "GB"`, `UnitTB = "TB"`, `UnitPB = "PB"`, `UnitEB = "EB"`, `UnitFloatFmt = "%.1f %s"`, `UnitIntFmt = "%d B"`.
- Git CLI: `GitBin = "git"`, `GitDirDotGit = ".git"`, `GitFlagDir = "-C"`, `GitSubCmdGC = "gc"`, `GitFlagAggressive = "--aggressive"`, `GitFlagPruneNow = "--prune=now"`.
- ANSI colors: `ColorCyan`, `ColorYellow`, `ColorGreen`, `ColorRed`, `ColorReset`.
- Messages & templates: `MsgNoReposFound`, `MsgRepoCountHeaderFmt`, `MsgRepoRowFmt`, `MsgWouldGCFmt`, `MsgGCStartFmt`, `MsgGCFailedFmt`, `MsgGCDoneFmt`, `ErrJSONEncodeFmt`, `ErrInvalidFormatFmt`, `ErrParseFlags = "size: failed to parse flags"`.

## Consuming Callers to Refactor
- `cli/cmdsize/size.go`:
  - CRITICAL: Enforce error return sovereignty! Remove `cliexit.HandleError(err, 2)` and `cliexit.Fail(...)`. Return `error` / `*appfault.AppError` directly.
  - Decouple from `cmdhygiene.HygieneFormatType`, using `enums.SizeOutputFormatType`.
  - Replace magic strings, numbers, and formats with `constants.*` and `enums.*`.
  - Invert conditionals to flatten nested blocks (guard clauses).
- `cli/cmdsize/helpers_shim.go`:
  - Refactor `parseHygieneFormat` to `parseSizeOutputFormat` returning `enums.SizeOutputFormatType`.
  - Replace worker limits with `constants.MinHygieneWorkers` and `constants.MaxHygieneWorkers`.
  - Replace `1024.0` with `constants.ByteUnitScale`.
- `cli/cmdsize/hook_shim.go`: Clean DI hook.
- Add `cli/cmdsize/size_test.go`: Unit tests for size argument parsing, format resolution, dry-run simulation, and JSON/CSV emission.

## Verification
- Unit tests in `cli/cmdsize/enums/enums_test.go` and `cli/cmdsize/constants/constants_test.go`.
- Unit tests in `cli/cmdsize/size_test.go`.
- `go vet ./cli/cmdsize/...`
- `go test -count=1 ./cli/cmdsize/...`
