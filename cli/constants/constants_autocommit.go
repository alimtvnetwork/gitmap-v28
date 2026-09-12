package constants

// Auto-commit messages for post-release commit behavior.
const (
	// Extra leading newline creates a clear visual gap between the
	// release-complete block and the auto-commit scanner so the two
	// phases feel like distinct steps in the pipeline.
	MsgAutoCommitScanning    = "\n\n  " + ColorBlue + "🔍 Checking for uncommitted changes..." + ColorReset + "\n"
	MsgAutoCommitReleaseOnly = "  " + ColorGreen + "✓ Release metadata committed:" + ColorReset + " %s\n"
	MsgAutoCommitPushed      = "  " + ColorGreen + "✓ Pushed to" + ColorReset + " %s\n"
	MsgAutoCommitNone        = "  " + ColorGreen + "✓ Working tree clean" + ColorReset + " — nothing else to commit\n"
	MsgAutoCommitPrompt      = "  " + ColorYellow + "→ Uncommitted changes outside .gitmap/release/:" + ColorReset + "\n"
	MsgAutoCommitFile        = "      " + ColorDim + "•" + ColorReset + " %s\n"
	MsgAutoCommitAsk         = "  " + ColorYellow + "→ Auto-commit these alongside the release? [y/N]: " + ColorReset
	MsgAutoCommitAll         = "  " + ColorGreen + "✓ All changes committed:" + ColorReset + " %s\n"
	MsgAutoCommitPartial     = "  " + ColorGreen + "✓ Committed .gitmap/release/ changes only:" + ColorReset + " %s\n"
	MsgAutoCommitSkipped     = "  " + ColorDim + "→ Skipped auto-commit (--no-commit)" + ColorReset + "\n"
	MsgAutoCommitDryRun      = "  " + ColorDim + "[dry-run] Would auto-commit release changes" + ColorReset + "\n"
	MsgAutoCommitSyncRetry   = "  " + ColorYellow + "→ Remote %s moved; rebasing and retrying push..." + ColorReset + "\n"
	ErrAutoCommitFailed      = "  " + ColorRed + "✗ Auto-commit failed: %v" + ColorReset + "\n"
	ErrAutoCommitPush        = "  " + ColorRed + "✗ Push failed: %v" + ColorReset + "\n"
	AutoCommitMsgFmt         = "Release %s"
	FlagDescNoCommit         = "Skip post-release auto-commit and push"
	FlagDescYes              = "Auto-confirm all prompts (e.g. commit)"
	MsgAutoCommitAutoYes     = "  " + ColorDim + "→ Auto-confirmed via -y flag" + ColorReset + "\n"

	// Git diff arguments for detecting changes.
	GitDiff         = "diff"
	GitDiffNameOnly = "--name-only"
	GitDiffCached   = "--cached"
)
