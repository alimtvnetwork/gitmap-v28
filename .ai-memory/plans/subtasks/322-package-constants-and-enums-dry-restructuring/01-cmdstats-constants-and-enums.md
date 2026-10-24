# Subtask 322-01: cmdstats Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdstats` (7 files)
- Leaf Packages to Create:
  - `cli/cmdstats/enums/`
    - `stats_enums.go`
    - `enums_test.go`
  - `cli/cmdstats/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `StatsOutputFormatType`: Output format (`table`, `json`, `tui`, `csv`).
2. `StatsMetricType`: Aggregation metric (`total_runs`, `success_count`, `fail_count`, `fail_rate`, `avg_duration_ms`, `min_duration_ms`, `max_duration_ms`, `total_commands`, `unique_commands`, `last_used`).
3. `StatsTimeRangeType`: Time range (`all`, `today`, `7d`, `30d`, `90d`).
4. `StatsCategoryType`: Category (`all`, `repository`, `remote`, `ssh`, `automation`, `maintenance`).
5. `StatsChartModeType`: Chart presentation (`none`, `bar`, `sparkline`, `histogram`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & aliases: `CmdStats = "stats"`, `CmdStatsAlias = "ss"`.
- Flags: `--command`, `--json`, `--help`.
- Wire contract keys: `WireKeyTotalCommands`, `WireKeyUniqueCommands`, `WireKeyTotalSuccess`, `WireKeyTotalFail`, etc.
- SSH forward candidate strings (`AllSSHForwardCandidates()`).
- Error codes & pattern substrings: `ErrCodeE1154`, `ErrCodeE1155`, `ErrCodeE1156`, `ErrPatternScanError`, `ErrPatternDriverValue`.
- Exit code: `ExitGeneralError = 1`.
- Zero magic strings.

## Consuming Callers to Refactor
- `cli/cmdstats/stats.go`
- `cli/cmdstats/statsrender.go`
- `cli/cmdstats/helpers_shim.go`
- `cli/cmdstats/hook_shim.go`
- `cli/cmdstats/stats_test.go`
- Contract tests in `cli/cmdstats/`

## Verification
- Unit tests in `cli/cmdstats/enums/enums_test.go` and `cli/cmdstats/constants/constants_test.go`.
- `go vet ./cli/cmdstats/...`
- `go test ./cli/cmdstats/...`
