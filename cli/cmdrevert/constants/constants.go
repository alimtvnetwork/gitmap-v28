package constants

// Command names and runner identifiers.
const (
	CmdRevert       = "revert"
	CmdRevertRunner = "revert-runner"
)

// Flag names for CLI argument parsing.
const (
	FlagListTxn  = "list-txn"
	FlagShowTxn  = "show-txn"
	FlagPruneTxn = "prune-txn"
	FlagTxn      = "txn"
	FlagLastTxn  = "last-txn"
	FlagLastNTxn = "last-n-txn"
	FlagForce    = "force"
	FlagVerbose  = "verbose"
)

// Flag CLI arguments formatted with double hyphen.
const (
	FlagArgListTxn  = "--list-txn"
	FlagArgShowTxn  = "--show-txn"
	FlagArgPruneTxn = "--prune-txn"
	FlagArgTxn      = "--txn"
	FlagArgLastTxn  = "--last-txn"
	FlagArgLastNTxn = "--last-n-txn"
	FlagArgForce    = "--force"
	FlagArgVerbose  = "--verbose"
)

// Flag descriptions for help text and diagnostic reporting.
const (
	FlagDescListTxn  = "list the last 50 transactions and exit"
	FlagDescShowTxn  = "print one transaction (including every captured file) by id, then exit"
	FlagDescTxn      = "revert the named transaction id"
	FlagDescLastTxn  = "revert the most recent committed transaction"
	FlagDescLastNTxn = "revert the most recent N committed transactions, newest first"
	FlagDescPruneTxn = "force a transaction prune cycle now (drops everything beyond the 50-row cap)"
	FlagDescForce    = "skip the confirm prompt before reverting"
	FlagDescVerbose  = "enable verbose output logging"
)

// Retention limits and status tokens.
const (
	TxnRetentionCap    = 50
	TxnStatusCommitted = "committed"
	TxnStatusReverted  = "reverted"
	ConfirmTokenYes    = "yes"
)

// System binaries, arguments, and script names.
const (
	BinGit          = "git"
	GitFlagDir      = "-C"
	GitCmdCheckout  = "checkout"
	BinPowershell   = "powershell"
	PSArgExecPolicy = "-ExecutionPolicy"
	PSValBypass     = "Bypass"
	PSArgNoProfile  = "-NoProfile"
	PSArgNoLogo     = "-NoLogo"
	PSArgFile       = "-File"
	ScriptRunPS1    = "run.ps1"
)

// Output messages and status banners.
const (
	MsgRevertUsage         = "Usage: gitmap revert <version>\nExample: gitmap revert v2.9.0"
	MsgRevertCheckout      = "  → Checking out %s...\n"
	MsgRevertStarting      = "\n  Building reverted version...\n"
	MsgRevertDone          = "\n  ✓ Revert complete.\n"
	MsgUpdateActive        = "  → Active: %s\n  → Handoff: %s\n"
	MsgTxnPruned           = "  ✓ pruned %d transaction(s)\n"
	MsgTxnReverted         = "  ✓ reverted transaction #%d (%s)\n"
	MsgTxnNoCommitted      = "  no committed transactions to revert\n"
	MsgTxnConfirmRevert    = "About to revert transaction #%d (%s, %d file(s)).\n"
	MsgTxnConfirmPrompt    = "Are you sure? Type 'yes' to proceed: "
	MsgTxnAbortedByUser    = "  revert aborted by user\n"
	MsgTxnLastNHeader      = "About to revert %d transactions (newest first):\n"
	MsgTxnLastNRow         = "  #%-4d %-10s %s  %s\n"
	MsgTxnLastNDone        = "  ✓ reverted %d transaction(s)\n"
	MsgTxnLastNNoneFound   = "  no committed transactions found to revert\n"
	MsgTxnListRowFmt       = "#%-4d %-10s %-9s %s  %s\n"
	MsgTxnDetailHeaderFmt  = "Transaction #%d\n  kind   : %s\n  status : %s\n  argv   : %s\n  cwd    : %s\n"
	MsgTxnDetailReverseFmt = "  reverse: %s\n  files  : %d\n"
	MsgTxnDetailFileRowFmt = "    [%s] %s\n"
)

// PowerShell runner script templates.
const (
	RevertPSHeader = `# gitmap revert script (auto-generated)
Set-Location "%s"
`
	RevertPSBuild = `
Write-Host ""
Write-Host "  Building from checked-out version..." -ForegroundColor Cyan
& "%s"
$runExit = $LASTEXITCODE
if (($runExit -ne 0) -and ($runExit -ne $null)) {
    exit $runExit
}

`
	RevertPSPostActions = `
$cmdAfter = Get-Command gitmap -ErrorAction SilentlyContinue | Select-Object -First 1
if ($cmdAfter -and (Test-Path $cmdAfter.Source)) {
    $activeAfter = & $cmdAfter.Source version 2>&1
    Write-Host "  Active version: $activeAfter" -ForegroundColor DarkGray
}

# NOTE: cleanup is intentionally NOT invoked here, mirroring the same
# decision in UpdatePSPostActions. Phase 3 (scheduleDeployedCleanupHandoff
# in runRevertRunner) handles cleanup after this PowerShell script and
# the surrounding handoff process exit, so the deployed binary can
# safely remove the still-locked handoff exe + .old backup. See
# 02-spec/08-generic-update/06-cleanup.md.

Write-Host ""
exit 0
`
)

// Error codes, domain names, and operation tags.
const (
	ErrCodeE9000            = "E9000"
	ErrCodeE2012            = "E2012"
	ErrCodeE2013            = "E2013"
	ErrCodeE2014            = "E2014"
	ErrCodeE2015            = "E2015"
	DomainCmdRevert         = "cmd.revert"
	DomainCmdRevertScript   = "cmd.revertscript"
	DomainCmdRevertTxn      = "cmd.reverttxn"
	DomainCmdRevertTxnLastN = "cmd.reverttxn_lastn"
	OpRevertCheckout        = "revert.checkout"
	OpRevertRunner          = "revert.runner"
	OpRunRevertPS           = "runRevertPS"
	OpTxnRevert             = "txn.revert"
	OpTxnRevertMany         = "txn.revertMany"
)

// Error messages, format strings, and verbose log templates.
const (
	ErrTagNotFound           = "Error: tag %s not found locally. Run 'git fetch --tags' first.\n"
	ErrNoRepoPath            = "repo path is not set; cannot run revert"
	ErrRevertLastNBadCount   = "revert: --%s requires a positive integer, got %q\n"
	ErrTxnRowNotFound        = "transaction not found"
	ErrInvalidTxnIDFmt       = "revert: invalid transaction id %q: %v\n"
	ErrInvalidTxnIDNotPosFmt = "revert: invalid transaction id %q\n"
	RevertScriptLogExec      = "executing revert script: %s"
	RevertScriptLogExit      = "revert script exited: err=%v"
)
