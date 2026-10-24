# Subtask 320-01: cmdsync Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdsync` (11 files)
- Leaf Packages to Create:
  - `cli/cmdsync/enums/`
    - `sync_enums.go` (or `sync_target_enum.go`, `sync_status_enum.go`, `sync_action_enum.go`)
    - `enums_test.go`
  - `cli/cmdsync/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `SyncTargetType`: Target synchronization scope (`all`, `dirty`, `clean`, `stashed`, `missing`).
2. `SyncRepoStatusType`: Status outcomes (`synced`, `modified`, `conflict`, `skipped`, `error`).
3. `SyncActionType`: Discrete sync actions (`fetch`, `pull`, `merge`, `rebase`, `checkout`).
4. `SyncModeType`: Operational modes (`fast-forward`, `rebase`, `dry-run`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (*Type, error)`, `All*() []...Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Subcommand names: `CmdSync = "sync"`
- Flags: `--all`, `--dirty`, `--dry-run`, `--rebase`, `--force`, `--group`
- File tokens: `.git`, `gitmap.config.json`
- Progress headers, UI banners, error codes.

## Consuming Callers to Refactor
- `cli/cmdsync/sync.go`
- `cli/cmdsync/types.go`
- `cli/cmdsync/worker.go`
- `cli/cmdsync/*.go`

## Verification
- Unit tests in `cli/cmdsync/enums/enums_test.go` and `cli/cmdsync/constants/constants_test.go`.
- `go vet ./cli/cmdsync/...`
- `go test ./cli/cmdsync/...`
