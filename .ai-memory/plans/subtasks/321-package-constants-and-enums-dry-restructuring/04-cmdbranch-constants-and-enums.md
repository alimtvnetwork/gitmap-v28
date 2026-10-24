# Subtask 321-04: cmdbranch Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdbranch` (3 files)
- Leaf Packages to Create:
  - `cli/cmdbranch/enums/`
    - `branch_enums.go`
    - `enums_test.go`
  - `cli/cmdbranch/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `BranchFilterType`: Filter criteria (`all`, `local`, `remote`, `merged`, `unmerged`).
2. `BranchActionType`: Actions/subcommands (`default`, `def`, `list`, `current`, `latest-branch`, `lb`, `help`).
3. `BranchSortOrderType`: Sorting (`recent`, `committerdate`, `authordate`, `alphabetical`).
4. `BranchOutputFormatType`: Output format (`terminal`, `json`, `csv`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & aliases: `CmdBranch = "branch"`, `CmdBranchAlias = "b"`, `CmdSubDefault = "default"`, `CmdSubDefaultAlias = "def"`, `CmdSubLatest = "latest-branch"`.
- CLI flags: `-j, --json`, `-h, --help`.
- Help text, usage lines, tips, error codes (`ErrCodeE9000`), error format strings, current dir sentinel (`"."`).

## Consuming Callers to Refactor
- `cli/cmdbranch/branch.go`
- `cli/cmdbranch/branch_help_menu.go`
- `cli/cmdbranch/hook_shim.go`

## Verification
- Unit tests in `cli/cmdbranch/enums/enums_test.go` and `cli/cmdbranch/constants/constants_test.go`.
- `go vet ./cli/cmdbranch/...`
- `go test ./cli/cmdbranch/...`
