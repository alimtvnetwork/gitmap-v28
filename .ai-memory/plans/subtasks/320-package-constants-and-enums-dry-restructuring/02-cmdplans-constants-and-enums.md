# Subtask 320-02: cmdplans Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdplans` (7 files)
- Leaf Packages to Create:
  - `cli/cmdplans/enums/`
    - `plan_enums.go`
    - `enums_test.go`
  - `cli/cmdplans/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `PlanStatusType`: Lifecycle status (`pending`, `active`, `in_progress`, `completed`, `archived`, `failed`).
2. `PlansFilterModeType`: Query filtering (`all`, `active`, `pending`, `completed`).
3. `PlansOutputFormatType`: Output rendering format (`table`, `json`, `markdown`, `csv`).
4. `PlanMetaFileType`: Known metadata filenames (`readme.md`, `summary.md`, `checklist.md`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (*Type, error)`, `All*() []...Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Subcommand names: `CmdPlans = "plans"`
- Flags: `--all`, `--status`, `--format`, `--json`, `--dir`
- Directory tokens: `.ai-memory/plans`, `pending/`, `completed/`, `subtasks/`
- ANSI colors, table column headers.

## Consuming Callers to Refactor
- `cli/cmdplans/plans.go`
- `cli/cmdplans/plans_types.go`
- `cli/cmdplans/*.go`

## Verification
- Unit tests in `cli/cmdplans/enums/enums_test.go` and `cli/cmdplans/constants/constants_test.go`.
- `go vet ./cli/cmdplans/...`
- `go test ./cli/cmdplans/...`
