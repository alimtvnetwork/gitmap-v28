package constants

// History-rewrite (`history-purge` / `history-pin`) flag names. See
// spec/04-generic-cli/16-history-rewrite.md.
const (
	HistoryFlagYes         = "yes"
	HistoryFlagYesShort    = "y"
	HistoryFlagNoPush      = "no-push"
	HistoryFlagDryRun      = "dry-run"
	HistoryFlagMessage     = "message"
	HistoryFlagKeepSandbox = "keep-sandbox"
	HistoryFlagQuiet       = "quiet"
	HistoryFlagQuietShort  = "q"
	HistoryPathSep         = ","
	HistorySandboxPrefix   = "gitmap-history-rewrite-"
	HistoryFilterRepoBin   = "git-filter-repo"
	HistoryGitBin          = "git"
	HistoryRemoteOrigin    = "origin"
	HistoryPushRefSpec     = "--mirror"
	HistoryForceWithLease  = "--force-with-lease"
)

// History-rewrite flag descriptions (rendered in --help).
const (
	HistoryDescYes         = "Skip the push confirmation prompt; force-push immediately on success"
	HistoryDescNoPush      = "Stop after verification; print the manual git push command"
	HistoryDescDryRun      = "Run the rewrite + verification in the sandbox, then exit without pushing"
	HistoryDescMessage     = "Rewrite the commit message of ONLY commits that touch a requested path; other commit messages are left unchanged"
	HistoryDescKeepSandbox = "Don't delete the temp mirror-clone on exit"
	HistoryDescQuiet       = "Suppress per-phase progress lines; only print errors and the final summary"
)

// History-rewrite exit codes. Match spec §6 exactly.
const (
	HistoryExitOk           = 0
	HistoryExitNotInRepo    = 2
	HistoryExitNoFilterRepo = 3
	HistoryExitBadArgs      = 4
	HistoryExitFilterFailed = 5
	HistoryExitVerifyFailed = 6
	HistoryExitPushFailed   = 7
)

// History-rewrite user-facing messages. All literals live here (no
// magic strings rule).
const (
	HistoryMsgPhaseIdentify    = "▸ history-rewrite: identifying origin remote\n"
	HistoryMsgPhaseClone       = "▸ history-rewrite: mirror-cloning %s into %s\n"
	HistoryMsgPhaseFilterPurge = "▸ history-rewrite: running filter-repo (purge) for %d path(s)\n"
	HistoryMsgPhaseFilterPin   = "▸ history-rewrite: running filter-repo (pin) for %d path(s)\n"
	HistoryMsgPhaseVerify      = "▸ history-rewrite: verifying sandbox\n"
	HistoryMsgPhasePush        = "▸ history-rewrite: pushing to %s with --force-with-lease\n"
	HistoryMsgVerifyOk         = "✓ history-rewrite: verification passed\n"
	HistoryMsgPushOk           = "✓ history-rewrite: push complete\n"
	HistoryMsgDryRunDone       = "✓ history-rewrite: dry-run complete; sandbox at %s\n"
	HistoryMsgKeepSandbox      = "▣ history-rewrite: sandbox kept at %s\n"
	HistoryMsgManualPush       = "▣ history-rewrite: to push manually, run:\n    git -C %s push --force-with-lease --mirror %s\n"
	HistoryMsgConfirmBanner    = "\n────────────────────────────────────────────────────────────\n  ✓ Verification PASSED\n  Mode      : %s\n  Paths     : %d\n  Sandbox   : %s\n  Remote    : %s\n  Action    : git push --mirror --force-with-lease\n  WARNING   : This rewrites published history. Downstream\n              clones will need to re-clone or hard-reset.\n────────────────────────────────────────────────────────────\n"
	HistoryMsgConfirmPush      = "Type 'yes' to force-push to %s (anything else aborts): "
	HistoryMsgUserAborted      = "▣ history-rewrite: aborted by user; sandbox at %s\n"
	HistoryMsgInstallHintLinux = "history-rewrite: install with: pip install --user git-filter-repo\n"
	HistoryMsgInstallHintMac   = "history-rewrite: install with: brew install git-filter-repo\n"
	HistoryMsgInstallHintWin   = "history-rewrite: install with: scoop install git-filter-repo (or pip install --user git-filter-repo)\n"
	HistoryMsgSummary          = "history-rewrite: %d path(s), %d commit(s) rewritten, repo size %s -> %s\n"
)

// History-rewrite errors.
const (
	HistoryErrNotInRepo       = "history-rewrite: not inside a git repository: %v\n"
	HistoryErrNoOrigin        = "history-rewrite: cannot read origin remote: %v\n"
	HistoryErrNoFilterRepo    = "history-rewrite: `git filter-repo` is not installed.\n"
	HistoryErrBadArgs         = "history-rewrite: %s\n"
	HistoryErrNoPaths         = "expected at least one path argument"
	HistoryErrConflictFlags   = "--yes and --no-push are mutually exclusive"
	HistoryErrPathNotReadable = "pin: cannot read %s from working tree: %v"
	HistoryErrMirrorClone     = "history-rewrite: mirror-clone failed: %v\n"
	HistoryErrFilterRepo      = "history-rewrite: filter-repo failed (exit %d): %s\n"
	HistoryErrVerifyPurge     = "history-rewrite: verification failed: path %s still has %d commit(s) in history\n"
	HistoryErrVerifyPin       = "history-rewrite: verification failed: path %s has %d distinct content hashes across history (expected 1)\n"
	HistoryErrPush            = "history-rewrite: push failed: %v\n"
	HistoryErrSandbox         = "history-rewrite: cannot create sandbox: %v\n"
	HistoryErrManifest        = "history-rewrite: cannot write blob manifest: %v\n"
)
