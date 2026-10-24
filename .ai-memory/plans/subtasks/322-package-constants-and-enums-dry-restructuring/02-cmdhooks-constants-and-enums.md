# Subtask 322-02: cmdhooks Constants & Enums DRY Restructuring

## Status: COMPLETED

## Target Scope
- Target Package: `cli/cmdhooks` (3 files refactored + new test suite)
- Leaf Packages Created:
  - `cli/cmdhooks/enums/`
    - `hooks_enums.go`
    - `enums_test.go`
  - `cli/cmdhooks/constants/`
    - `constants.go`
    - `constants_test.go`

## Enums Implemented
1. `GitHookType`: Hook types (`commit-msg`, `pre-commit`, `pre-push`, `post-checkout`, `post-merge`, `pre-rebase`, `prepare-commit-msg`, `post-commit`).
2. `HookActionType`: CLI actions (`install`, `uninstall`, `run`, `check`, `list`, `status`, `help`, `-h`, `--help`).
3. `HookPlanActionType`: Plan execution actions (`install`, `already`, `skip`).
4. `HookScopeType`: Target repository scopes (`all`, `repo`, `global`, `workdir`).

All enums implement:
- `...Type` naming suffix.
- `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive boolean helpers (`is*`, `has*`).
- Function lengths <= 15 lines (<= 8 lines preferred), blank lines after `}` and before `return`.

## Constants Implemented
- Command names & aliases: `CmdHooks = "hooks"`, `CmdHooksAlias = "hkin"`, `SubCmdInstall = "install"`, etc.
- CLI flags & descriptions: `--all`, `--workdir`, `--dry-run`, `--yes`, `--help`, `-h`.
- Directory & path tokens: `.git`, `hooks`, `commit-msg`, `HookFilePerm = 0755`.
- Guard marker & template: `GuardMarker = "gitmap commit-msg guard"`, `CommitMsgGuardTemplate`.
- Format strings, summary templates, error messages, and usage text.

## Consuming Callers Refactored
- `cli/cmdhooks/hooks_dispatch.go`: Map dispatch converted to use `enums.HookActionType` and `constants.Err*`.
- `cli/cmdhooks/hooks_install.go`: Decomposed monolithic functions into <= 15 line functions with positive booleans and zero magic strings.
- `cli/cmdhooks/hooks_template.go`: Bound to constants package.
- `cli/cmdhooks/hooks_test.go`: Added test coverage for dispatch, flag validation, dry-run, and live guard installation.

## Verification
- Unit tests: `go test -v ./cmdhooks/...` PASS (100% pass)
- Lint check: `go vet ./cmdhooks/...` PASS (0 warnings)
