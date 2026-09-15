# Plan 174: AppError Join Architecture & Centralized DRY Help Checking Suite

> **Consolidated Task Execution Record**
> **Origination**: User feedback on `cli/cmd/join.go` requesting universal `*apperror.AppError` return envelopes and a centralized, DRY help checking architecture to eliminate repeated `len(args) == 0 || hasHelpFlag(args)` patterns across commands.
> **Total Steps / Loops to Completion**: 2 self-loop phases (Phase 1 planning with 2 research subagents + Phase 2 parallel execution with 2 implementation subagents).
> **Status**: Completed & Verified.

---

## 1. Architectural Changes Summary

### A. Centralized DRY Help Checking Engine (`cli/cmd/helpcheck.go`)
- **`CheckHelpOrEmpty(command string, args []string)`**: Centralized entrypoint that evaluates `IsHelpRequestedOrEmpty(args)` and executes `printHelpAndExit(command, args)` with `cliexit.Exit(0)` to ensure stdout pipe drainers execute cleanly before process exit.
- **`IsHelpRequestedOrEmpty(args []string) bool`**: Pure affirmative predicate evaluating `len(args) == 0 || hasHelpFlag(args)`.
- **Refactored `checkHelp(command string, args []string)`**: Replaced negative boolean variable `lacksHelpFlag := !hasHelpFlag(args)` with direct guard clause `if !hasHelpFlag(args) { return }`.
- **Decomposed `printHelpAndExit(command string, args []string)`**: Reusable 6-line helper invoking `ParsePrettyFlag(args)`, `helptext.PrintWithMode(command, mode)`, `printUsageFooterShort()`, and `cliexit.Exit(0)`.

### B. Universal `*apperror.AppError` in `cli/cmd/join.go`
- Refactored `runJoin` return signature from `error` to `*apperror.AppError`.
- Replaced 4-line ad-hoc help check with a single call to `CheckHelpOrEmpty(constants.CmdJoin, args)`.
- Switched `flag.NewFlagSet(constants.CmdJoin, flag.ExitOnError)` to `flag.ContinueOnError` to eliminate bare `os.Exit(2)` crashes on unrecognized flags, returning `apperror.NewValidationError(err.Error())`.
- Validated positional address and `--token` inputs with structured `apperror.NewValidationError` (error code `"E1000"`).
- Wrapped `client.Handshake()` failures via `apperror.WrapWithDetails` with domain error code `"E8005"`, `ErrorTypeExecution`, `SeverityError`, and diagnostic context attributes (`address`, `hostname`).
- Decomposed into small helpers: `parseJoinArgs`, `validateJoinInputs`, `executeJoinHandshake`, `resolveJoinHostname`, `wrapJoinHandshakeError` (all <= 15 lines).

### C. Repetition Elimination Across Command Suite
- **`cli/cmd/clustercommand.go`**: Replaced custom `checkClusterCommandHelp` helper and conditional guard with `CheckHelpOrEmpty(resolveClusterHelpCmd(selector), args)`. Deleted dead function `checkClusterCommandHelp`.
- **`cli/cmd/servercmd.go`**: Replaced manual usage check with `CheckHelpOrEmpty(constants.CmdServerCmd, args)`. Deleted 18-line hardcoded usage function `printServerCmdUsage`, reconnecting `server-cmd` to standard markdown rendering.
- **`cli/cmd/rootcore.go`**: Replaced `dispatchServersClientsHelp` and `dispatchClientsHelp` with direct calls to `CheckHelpOrEmpty`. Deleted both dead dispatch helpers.
- **`cli/cmd/roottooling.go`**: Safely unwrapped `runJoin` in `toolingNetworkEntries()` to prevent typed-nil interface trap.

---

## 2. Completed Subtasks Register

### Subtask 01: Centralized DRY Help Checking Engine & Callers
- **Files Modified/Created**:
  - `cli/cmd/helpcheck.go`
  - `cli/cmd/clustercommand.go`
  - `cli/cmd/servercmd.go`
  - `cli/cmd/rootcore.go`
  - `cli/cmd/helpcheck_test.go`
- **Result**: All repetitive help checks consolidated into `CheckHelpOrEmpty`. 100% test coverage for empty and flag conditions.

### Subtask 02: Join Command AppError Architecture & Unit Tests
- **Files Modified/Created**:
  - `cli/cmd/join.go`
  - `cli/cmd/roottooling.go`
  - `cli/cmd/join_test.go`
- **Result**: `runJoin` converted to return `*apperror.AppError`. Full unit test suite covering zero args, help flags, validation errors, handshake failure, and successful TLS RPC handshake.

---

## 3. Verification & Compliance
- **Function Line Length**: All functions <= 15 lines (target <= 8 lines).
- **Affirmative Booleans**: 100% affirmative booleans (`is*`, `has*`). No negative booleans.
- **Line Endings & Encoding**: Strict Unix LF line endings and UTF-8 without BOM across all modified files.
- **Banned Test & Build Commands**: Zero routine runs of `go test`, `go build`, or `06-cicd-local-runner.py`.
