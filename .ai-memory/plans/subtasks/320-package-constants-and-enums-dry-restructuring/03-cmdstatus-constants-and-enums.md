# Subtask 320-03: cmdstatus Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdstatus` (11 files)
- Leaf Packages to Create:
  - `cli/cmdstatus/enums/`
    - `status_enums.go`
    - `enums_test.go`
  - `cli/cmdstatus/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `StatusFilterModeType`: Query filtering (`all`, `dirty`, `ahead`, `behind`).
2. `StatusOutputModeType`: Rendering formats (`table`, `json`, `compact`).
3. `StatusTableColumnType`: Table column titles (`Repository`, `Branch`, `Status`, `Sync`, `Stash`, `Files`).
4. `StatusRepoStateType`: State tags (`clean`, `dirty`, `ahead`, `behind`, `stashed`, `missing`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (*Type, error)`, `All*() []...Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Subcommand names: `CmdStatus = "status"`, `CmdStatusAlias = "st"`
- Flags: `--all`, `--dirty`, `--only-dirty`, `--ahead`, `--behind`, `--table`, `--compact`, `--json`, `-g`, `--group`
- UI banners & messages: `StatusNotFoundText = "✖ not found"`, `StatusCleanSuccessMsg = "✨ All repositories are !clean"`, `StatusTableColGap = "   "`.

## Consuming Callers to Refactor
- `cli/cmdstatus/status.go`
- `cli/cmdstatus/status_render.go`
- `cli/cmdstatus/status_resolve.go`
- `cli/cmdstatus/status_help_menu.go`
- `cli/cmdstatus/helpers_shim.go`
- `cli/cmdstatus/*.go`

## Verification
- Unit tests in `cli/cmdstatus/enums/enums_test.go` and `cli/cmdstatus/constants/constants_test.go`.
- `go vet ./cli/cmdstatus/...`
- `go test ./cli/cmdstatus/...`
