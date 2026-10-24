# Subtask 322-03: cmdprune Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdprune` (3 files)
- Leaf Packages to Create:
  - `cli/cmdprune/enums/`
    - `prune_enums.go`
    - `enums_test.go`
  - `cli/cmdprune/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `PruneScopeType`: Pruning scopes (`local`, `remote`, `all`, `stale`).
2. `PruneTargetType`: Prune targets (`branches`, `tags`, `remotes`, `objects`, `worktrees`).
3. `PruneActionType`: Execution actions (`delete`, `dry-run`, `prompt`, `keep`, `skip`).
4. `PruneStatusType`: Candidate outcomes (`deleted`, `kept`, `skipped`, `failed`, `aborted`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & aliases: `CmdPrune = "prune"`, `CmdPruneAlias = "prn"`.
- CLI flags & descriptions: `--dry-run`, `--confirm`, `--remote`.
- Git CLI tokens: `GitBin = "git"`, `GitBranch = "branch"`, `GitBranchDeleteFlag = "-D"`, `ReleaseBranchPrefix = "release/"`, `GitWildcard = "*"`.
- Prompt answers: `y`, `Y`, `n`, `N`.
- Output messages, summary formats, and error strings.
- Decouple completely from monolithic root `cli/constants`.

## Consuming Callers to Refactor
- `cli/cmdprune/prune.go`
- `cli/cmdprune/pruneops.go`
- `cli/cmdprune/hook_shim.go`

## Verification
- Unit tests in `cli/cmdprune/enums/enums_test.go` and `cli/cmdprune/constants/constants_test.go`.
- `go vet ./cli/cmdprune/...`
- `go test ./cli/cmdprune/...`
