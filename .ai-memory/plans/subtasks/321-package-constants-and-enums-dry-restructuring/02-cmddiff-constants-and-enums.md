# Subtask 321-02: cmddiff Constants & Enums DRY Restructuring

## Status
- **Status:** COMPLETED
- **Completed At:** 2026-10-24T20:53:00+08:00
- **Verification:** All tests in `cli/cmddiff/...`, `cli/cmddiff/constants/...`, `cli/cmddiff/enums/...` pass cleanly with zero vet errors.

## Target Scope
- Target Package: `cli/cmddiff` (6 files)
- Leaf Packages to Create:
  - `cli/cmddiff/enums/`
    - `diff_enums.go`
    - `enums_test.go`
  - `cli/cmddiff/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `DiffViewModeType`: Presentation modes (`unified`, `side-by-side`, `stat`, `summary`, `compact`).
2. `DiffOutputFormatType`: Output format (`text`, `json`, `raw`).
3. `DiffScopeType`: Evaluation scope (`working-tree`, `staged`, `branch`, `commit`).
4. `DiffFilterType`: Filter criteria (`all`, `only-conflicts`, `only-missing`, `include-identical`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & aliases: `CmdDiff = "diff"`, `CmdDiffAlias = "df"`, `CmdCompare = "compare"`, `CmdDiffProfiles = "diff-profiles"`.
- CLI flags & descriptions: `--stat`, `--cached`, `--staged`, `--name-only`, `--json`, `--help`.
- Layout tokens, help menu strings, error wrap tokens, and positional argument limits.
- Zero magic integer exit codes: fix `cliexit.HandleError(err, 2)` in `diff.go:63` to use typed constant `constants.ExitUsageError`.

## Consuming Callers to Refactor
- `cli/cmddiff/diff.go`
- `cli/cmddiff/diff_help_menu.go`
- `cli/cmddiff/dispatchcompare.go`
- `cli/cmddiff/dispatchdiff.go`

## Verification
- Unit tests in `cli/cmddiff/enums/enums_test.go` and `cli/cmddiff/constants/constants_test.go`.
- `go vet ./cli/cmddiff/...`
- `go test ./cli/cmddiff/...`
