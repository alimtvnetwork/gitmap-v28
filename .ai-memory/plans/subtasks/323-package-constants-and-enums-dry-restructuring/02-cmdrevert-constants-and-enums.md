# Subtask 323-02: cmdrevert Constants & Enums DRY Restructuring

**Status:** COMPLETED
**Completed At:** 2026-10-24T21:55:00+08:00

## Target Scope
- Target Package: `cli/cmdrevert` (6 files: `revert.go`, `revertscript.go`, `reverttxn.go`, `reverttxn_lastn.go`, `helpers_shim.go`, `hook_shim.go`)
- Leaf Packages Created:
  - `cli/cmdrevert/enums/`
    - `revert_enums.go`
    - `enums_test.go`
  - `cli/cmdrevert/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Implemented
1. `RevertModeType`: Modes/subcommands (`version`, `revert-runner`, `txn`, `last-txn`, `last-n-txn`, `list-txn`, `show-txn`, `prune-txn`).
2. `RevertPromptResponseType`: Affirmative responses (`yes`, `no`).
3. `RevertExitCodeType`: Strongly-typed exit codes (`0` success, `1` general, `2` usage).

Each enum follows:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies (standard library only).

## Specific Constants Implemented
- Command names: `CmdRevert = "revert"`, `CmdRevertRunner = "revert-runner"`.
- Flags: `FlagListTxn`, `FlagShowTxn`, `FlagPruneTxn`, `FlagTxn`, `FlagLastTxn`, `FlagLastNTxn`, `FlagForce`, `FlagVerbose`.
- Flag CLI argument strings: `FlagArgListTxn`, `FlagArgShowTxn`, `FlagArgPruneTxn`, `FlagArgTxn`, `FlagArgLastTxn`, `FlagArgLastNTxn`, `FlagArgForce`, `FlagArgVerbose`.
- Flag descriptions: `FlagDescListTxn`, `FlagDescShowTxn`, `FlagDescTxn`, `FlagDescLastTxn`, `FlagDescLastNTxn`, `FlagDescPruneTxn`, `FlagDescForce`, `FlagDescVerbose`.
- Retention & tokens: `TxnRetentionCap = 50`, `TxnStatusCommitted = "committed"`, `TxnStatusReverted = "reverted"`, `ConfirmTokenYes = "yes"`.
- System binaries & args: `BinGit = "git"`, `GitFlagDir = "-C"`, `GitCmdCheckout = "checkout"`, `BinPowershell = "powershell"`, `PSArgExecPolicy = "-ExecutionPolicy"`, `PSValBypass = "Bypass"`, `PSArgNoProfile = "-NoProfile"`, `PSArgNoLogo = "-NoLogo"`, `PSArgFile = "-File"`, `ScriptRunPS1 = "run.ps1"`.
- Messages & templates: `MsgRevertUsage`, `MsgRevertCheckout`, `MsgRevertStarting`, `MsgRevertDone`, `MsgUpdateActive`, `MsgTxnPruned`, `MsgTxnReverted`, `MsgTxnNoCommitted`, `MsgTxnConfirmRevert`, `MsgTxnConfirmPrompt`, `MsgTxnAbortedByUser`, `MsgTxnLastNHeader`, `MsgTxnLastNRow`, `MsgTxnLastNDone`, `MsgTxnLastNNoneFound`, `MsgTxnListRowFmt`, `MsgTxnDetailHeaderFmt`, `MsgTxnDetailReverseFmt`, `MsgTxnDetailFileRowFmt`.
- PowerShell runner templates: `RevertPSHeader`, `RevertPSBuild`, `RevertPSPostActions`.
- Error codes & domains: `ErrCodeE9000`, `ErrCodeE2012`, `ErrCodeE2013`, `ErrCodeE2014`, `ErrCodeE2015`, `DomainCmdRevert`, `DomainCmdRevertScript`, `DomainCmdRevertTxn`, `DomainCmdRevertTxnLastN`, `OpRevertCheckout`, `OpRevertRunner`, `OpRunRevertPS`, `OpTxnRevert`, `OpTxnRevertMany`.
- Error messages & logs: `ErrTagNotFound`, `ErrNoRepoPath`, `ErrRevertLastNBadCount`, `ErrTxnRowNotFound`, `ErrInvalidTxnIDFmt`, `ErrInvalidTxnIDNotPosFmt`, `RevertScriptLogExec`, `RevertScriptLogExit`.

## Consuming Callers Refactored
- `cli/cmdrevert/revert.go`: Decoupled completely from root `cli/constants`. Eliminated magic integer exit code `1` by returning `error` or `*apperror.AppError` directly. Handled `validateRevertVersion`, `checkoutRevertTag`, and `launchRevertHandoff` with error return sovereignty.
- `cli/cmdrevert/revertscript.go`: Used `constants.BinPowershell`, `constants.ScriptRunPS1`, `constants.RevertPS*` templates, and package constants. Leaf functions return errors directly.
- `cli/cmdrevert/reverttxn.go`: Refactored flag checking, list/show/prune handlers to use `constants.*` and `enums.*`. Replaced magic exit codes with `enums.ExitCodeGeneralError`. Leaf functions preserve error sovereignty.
- `cli/cmdrevert/reverttxn_lastn.go`: Used `constants.ConfirmTokenYes`, `constants.MsgTxnLastN*`, `constants.ErrRevertLastNBadCount`. Positive booleans only.
- `cli/cmdrevert/helpers_shim.go` and `hook_shim.go`: Clean DI hooks for `RepoPathFn` and `openDBFn`.
- Added `cli/cmdrevert/revert_test.go`: Comprehensive unit tests for version validation, mode resolution, flag parsing, script generation, and dry-run execution.

## Verification
- Unit tests in `cli/cmdrevert/enums/enums_test.go` and `cli/cmdrevert/constants/constants_test.go` (100% PASS).
- Unit tests in `cli/cmdrevert/revert_test.go` (100% PASS).
- `go vet ./cmdrevert/...` (PASS, 0 issues).
- `go test -count=1 -v ./cmdrevert/...` (PASS, 100% pass across all packages).
