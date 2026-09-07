# Subtask 04: OS Power SQLite Persistence & CLI Command Integration

## Status
Completed

## Context & Objectives
1. **SQLite Persistence**:
   - `gitmap/store/power.go`: CRUD operations for `PowerSetting` (profiles: `baseline`, `current`, `previous`) and `PowerSettingHistory` (audit trail).
   - SQL DDL in `gitmap/store/power_schema.go` (kept internal to avoid ERD drift or registered cleanly).
2. **CLI Commands & Dispatcher**:
   - Top-level constants in `gitmap/constants/constants_cli.go`: `CmdPower = "power"`, `CmdPowerAlias = "pw"`, `CmdPowerAlias2 = "pwr"`.
   - Synchronize `cmd_constants_test.go:topLevelCmds()`.
   - Dispatch entry in `gitmap/cmd/rootutility.go`.
   - `gitmap/cmd/power.go` (<= 200 lines): Subcommand router.
   - `gitmap/cmd/power_ops.go` (<= 200 lines): `status`, `never-sleep`, `set`, `reset`, `history`.

## Verification Steps
- AST parity test passes: `go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1`.
- `gitmap power --help` and unit tests pass.
