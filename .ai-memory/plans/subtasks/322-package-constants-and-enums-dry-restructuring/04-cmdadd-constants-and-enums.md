# Subtask 322-04: cmdadd Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdadd` (2 files) + `cli/cmd/add` delegate
- Leaf Packages to Create:
  - `cli/cmdadd/enums/`
    - `add_enums.go`
    - `enums_test.go`
  - `cli/cmdadd/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `AddTargetType`: Targets (`repo`, `folder`, `lfs`, `hook`, `common-attr`, `common-ignore`, `attributes`, `ignore`, `lfs-install`).
2. `AddActionType`: Actions (`create`, `append`, `merge`, `overwrite`, `diff`).
3. `AddModeType`: Execution mode (`normal`, `dry-run`, `force`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & aliases: `CmdAdd = "add"`, `CmdAddAlias = "ad"`.
- Subcommands: `ignore`, `attributes`, `lfs-install`, `common-attr`, `common-ignore`.
- Files & permissions: `.gitattributes`, `.gitignore`, `0644`.
- Builtin file templates: `TemplateCommonAttr`, `TemplateCommonIgnore`.
- Error codes & messages: `ErrCodeAddMissingArg`, `ErrCodeAddUnknownTarget`, `ErrPrefixError`.
- Operation & component names: `OpAddRun`, `OpCmdAdd`, `OpWriteCommonAttr`.

## Consuming Callers to Refactor
- `cli/cmdadd/add_entry.go`
- `cli/cmdadd/hook_shim.go`
- `cli/cmd/add/add.go`

## Verification
- Unit tests in `cli/cmdadd/enums/enums_test.go` and `cli/cmdadd/constants/constants_test.go`.
- `go vet ./cli/cmdadd/...`
- `go test ./cli/cmdadd/...`
