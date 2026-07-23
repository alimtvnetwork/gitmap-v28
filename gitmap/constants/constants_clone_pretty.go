// Package constants — constants_clone_pretty.go: format strings for
// the unified colorful clone runner (clonepretty.go). Centralized
// per the no-magic-strings rule. All ANSI escapes pull from
// constants_terminal.go so theme changes flow everywhere at once.
package constants

const (
	// MsgClonePrettyHeader prints before every git clone invocation.
	// Format slots: url, dest, gitBin, gitClone, url, dest.
	MsgClonePrettyHeader = "\n" +
		"  " + ColorCyan + "▸ git clone" + ColorReset + "  " + ColorWhite + "%s" + ColorReset + "\n" +
		"  " + ColorDim + "  target  " + ColorReset + "%s\n" +
		"  " + ColorDim + "  exec    " + ColorReset + "%s %s %s %s\n"

	// MsgClonePrettyOK renders after a successful clone. Format slots:
	// dest, elapsed.
	MsgClonePrettyOK = "  " + ColorGreen + "✓ cloned" + ColorReset + " %s " +
		ColorDim + "(in %s)" + ColorReset + "\n"

	// MsgClonePrettyFail renders the failure panel on stderr. Format
	// slots: cmdline, exitCode, err, elapsed, hints.
	MsgClonePrettyFail = "\n" +
		"  " + ColorRed + "✗ git clone failed" + ColorReset + "\n" +
		"  " + ColorDim + "  command  " + ColorReset + "%s\n" +
		"  " + ColorDim + "  exit     " + ColorReset + "%d\n" +
		"  " + ColorDim + "  error    " + ColorReset + "%v\n" +
		"  " + ColorDim + "  elapsed  " + ColorReset + "%s\n" +
		"  " + ColorYellow + "→ try one of:" + ColorReset + "\n%s\n\n"

	// MsgCloneDryRunNoop is emitted in --dry-run mode in place of
	// the actual git invocation.
	MsgCloneDryRunNoop = "  " + ColorYellow + "[dry-run] would execute the command above — nothing cloned." + ColorReset

	// MsgCloneSpinnerLabel is the steady-state spinner label.
	MsgCloneSpinnerLabel = ColorWhite + "cloning…" + ColorReset
)

// Dry-run flag identifiers shared by clone, cfr, cfrp.
const (
	FlagCloneDryRun      = "dry-run"
	FlagCloneDryRunShort = "n"
	FlagDescCloneDryRun  = "Print the exact `git clone` command(s) and target path(s) without cloning"
)
