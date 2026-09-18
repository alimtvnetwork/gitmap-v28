# Subtask 03: Auxiliary Command Routers ErrorWrapper Refactor

## 1. Objectives
Refactor remaining auxiliary command routing and dispatch functions across `cli/cmdpipeline/`, `cli/cmddb/`, `cli/cmdvscode/`, `cli/cmdzsh/`, `cli/cmdmacro/`, and `cli/cmd/` to return `result.ErrorWrapper`, eliminating all intermediate `.AsError()` conversions.

## 2. Disjoint Target Files
- `cli/cmdpipeline/pipeline_db_dispatch.go`
- `cli/cmd/profiles_cmd.go`
- `cli/cmd/startup_cmd.go`
- `cli/cmd/storage_cmd.go`
- `cli/cmd/cd.go`
- `cli/cmddb/cmddb_dispatch.go`
- `cli/cmdvscode/vscode_cmd.go`
- `cli/cmdzsh/zsh_cmd.go`
- `cli/cmdmacro/macro_cmd.go`

## 3. Implementation Details

### A. `cli/cmdpipeline/pipeline_db_dispatch.go`
1. `routePipelineDB(args []string) result.ErrorWrapper`:
   - If empty args, returns `result.MatchWrapper(runPipelineDBStatusWithTelemetry(nil))`.
   - Calls `res := dispatchPipelineDBSubcmd(sub, args[1:])`. If matched, returns `res` directly.
   - Prints usage and returns `result.FailureWrapper(apperror.NewValidationError("unknown pipeline db subcommand: " + sub))`.
2. `handlePipelineDB(args []string) error`:
   - Returns `routePipelineDB(args).AsError()`.

### B. `cli/cmd/profiles_cmd.go`
1. `routeProfilesSub(sub string, args []string) result.ErrorWrapper`:
   - Matches subcommands (`ls`, `set-default`, `switch`, `add`, `rm`, `status`, default), wrapping calls in `result.MatchWrapper`.
2. `runProfiles(args []string) error`:
   - Returns `routeProfilesSub(sub, rest).AsError()`.

### C. `cli/cmd/startup_cmd.go`
1. `routeStartupSubcommand(sub string, tail []string) result.ErrorWrapper`:
   - Matches subcommands (`ls`, `add`, `rm`, `run`, `logs`), wrapping calls in `result.MatchWrapper`.
   - For `help`, returns `result.SuccessWrapper()`.
   - Default returns `result.FailureWrapper(...)`.
2. `runStartup(args []string) error`:
   - Returns `routeStartupSubcommand(sub, tail).AsError()`.

### D. `cli/cmd/storage_cmd.go`
1. `routeStorageSubcommand(args []string) result.ErrorWrapper`:
   - Matches list, clean, and drive report subcommands, wrapping calls in `result.MatchWrapper`.
2. `runStorage(args []string) error`:
   - Returns `routeStorageSubcommand(args).AsError()`.

### E. `cli/cmd/cd.go`
1. `routeCDSub(sub string, args []string) result.ErrorWrapper`:
   - Matches `CmdCDRepos`, `CmdCDSetDefault`, `CmdCDClearDefault`, default `runCDLookup`, wrapping calls in `result.MatchWrapper`.
2. `runCD(args []string) error`:
   - Returns `routeCDSub(sub, rest).AsError()`.

### F. `cli/cmddb/cmddb_dispatch.go`
1. `routeDbSubcommand(sub string, tail []string) result.ErrorWrapper`:
   - Matches subcommands (`status`, `optimize`, `ls`, `help`, `repo-db`, `sizes`, `reset`, `clear`), wrapping calls in `result.MatchWrapper`.
2. `runDb(tail []string) error`:
   - If empty, returns `runDBStatus(nil)`.
   - Returns `routeDbSubcommand(tail[0], tail[1:]).AsError()`.

### G. `cli/cmdvscode/vscode_cmd.go`
1. `routeVSCodeProjectAction(sub string, args []string) result.ErrorWrapper`:
   - Wraps calls in `result.MatchWrapper` or returns `result.SuccessWrapper()`.
2. `routeVSCodeMaintenanceAction(sub string, args []string) result.ErrorWrapper`:
   - Wraps calls in `result.MatchWrapper` or returns `result.SuccessWrapper()`.
3. `routeVSCodeSubcommand(sub string, args []string) result.ErrorWrapper`:
   - Routes between project actions and maintenance actions.
4. `dispatchVSCodeAction(args []string) error`:
   - Returns `routeVSCodeSubcommand(sub, args).AsError()`.

### H. `cli/cmdzsh/zsh_cmd.go`
1. `routeZshSubcommand(subcmd string, args []string) result.ErrorWrapper`:
   - Dispatches via router map, wrapping result in `result.MatchWrapper`.
   - Fallback returns `routeFallback(subcmd, args)`.
2. `routeFallback(subcmd string, args []string) result.ErrorWrapper`:
   - Wraps help and unknown handler in `result.MatchWrapper`.
3. `runZsh(args []string) error`:
   - Returns `routeZshSubcommand(subcmd, subArgs).AsError()`.

### I. `cli/cmdmacro/macro_cmd.go`
1. `routeMacroSubcommand(sub string, rest []string) result.ErrorWrapper`:
   - Dispatches to `routeExecSubcommand`, `routeExportImportSubcommand`, or `routeManagementSubcommand`.
2. `routeExportImportSubcommand(sub string, rest []string) result.ErrorWrapper`:
   - Wraps export/import actions in `result.MatchWrapper`.
3. `routeExecSubcommand(sub string, rest []string) result.ErrorWrapper`:
   - Wraps retry, run-until, execute in `result.MatchWrapper`.
4. `routeManagementSubcommand(sub string, rest []string) result.ErrorWrapper`:
   - Wraps startup, schedule, etc. in `result.MatchWrapper`.
5. `runMacroCmd(args []string) error`:
   - Returns `routeMacroSubcommand(args[0], args[1:]).AsError()`.

## 4. Verification
- Functions strictly <= 15 lines.
- Affirmative booleans only.
- Strict Unix LF line endings.
- Targeted linters only.
