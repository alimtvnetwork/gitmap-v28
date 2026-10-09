package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspace"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdisplay"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

// CheckHelpOrEmpty prints help and exits if args is empty or help flag is present.
func CheckHelpOrEmpty(command string, args []string) {
	if IsHelpRequestedOrEmpty(args) {
		printHelpAndExit(command, args)
	}
}

// IsHelpRequestedOrEmpty returns true if args is empty or contains help flags.
func IsHelpRequestedOrEmpty(args []string) bool {
	return len(args) == 0 || hasHelpFlag(args)
}

// checkHelp prints embedded help and exits if --help or -h is present.
// Honors --pretty / --no-pretty so users can force-enable rendering for
// pagers (`gitmap foo --help --pretty | less -R`) or strip ANSI for
// scripting (`gitmap foo --help --no-pretty > help.txt`).
//
// Uses cliexit.Exit so theme/glyphs pipe drainers run before the
// process teardown.
func checkHelp(command string, args []string) {
	if !hasHelpFlag(args) {
		return
	}

	printHelpAndExit(command, args)
}

// printHelpAndExit prints embedded help with parsed pretty mode and exits with code 0.
func printHelpAndExit(command string, args []string) {
	if tryPrintDisplayerHelp(command, args) {
		return
	}
	if tryRenderRichTopic(command) {
		cliexit.Exit(0)
	}
	if helpdoc.HasTopic(command) {
		_, mode := ParsePrettyFlag(args)
		helpdoc.PrintWithMode(command, mode)
		printUsageFooterShort()
		cliexit.Exit(0)
	}
	if RenderDynamicCommandHelp(command) {
		cliexit.Exit(0)
	}
	cliexit.Exit(0)
}

// helpDisplayEntry binds a command to its HelpDisplay builder plus the
// sub-command builders rendered on `--help <sub>`.
type helpDisplayEntry struct {
	build func() helpdisplay.Displayer
	subs  map[string]func() helpdisplay.Displayer
}

// helpDisplayRegistry maps pilot commands (spec 243.3) to their HelpDisplay
// builders. A registered command renders through the displayer on --help;
// unregistered commands keep today's behavior.
var helpDisplayRegistry = map[string]helpDisplayEntry{
	constants.CmdSpace: {
		build: cmdspace.SpaceHelpDisplay,
		subs: map[string]func() helpdisplay.Displayer{
			"common":                       cmdspace.SpaceCommonHelpDisplay,
			constants.CmdSpaceBackupBranch: cmdspace.SpaceBackupBranchHelpDisplay,
		},
	},
	constants.CmdScan:    {build: cmdscan.ScanHelpDisplay},
	constants.CmdPullAll: {build: cmdpull.PullAllHelpDisplay},
}

// TryPrintDisplayerHelp renders a registered command's HelpDisplay and exits.
// It is the exported entry point for dispatch surfaces that short-circuit
// before printHelpAndExit (e.g. root.go's tryInterceptCommandHelp, which
// intercepts rich topics like `scan --help` first). Returns false when the
// command is not registered.
func TryPrintDisplayerHelp(command string, args []string) bool {
	return tryPrintDisplayerHelp(command, args)
}

// tryPrintDisplayerHelp renders a registered command's HelpDisplay and exits.
// It honors `--help <sub>`: the token after the help flag selects the
// sub-display when the command's helper carries that sub. Returns false
// when the command is not registered.
func tryPrintDisplayerHelp(command string, args []string) bool {
	entry, isRegistered := helpDisplayRegistry[resolveHelpCommand(command)]
	if !isRegistered {
		return false
	}
	builder := displayerBuilderForSub(entry, args)
	builder().Print(helpdisplay.NewRenderContext(nil))
	cliexit.Exit(0)

	return true
}

// helpDisplayAlias maps pilot-command aliases to their canonical names so
// `--help` works identically through every alias (aliases are a kept feature).
var helpDisplayAlias = map[string]string{
	constants.CmdScanAlias:    constants.CmdScan,
	constants.CmdPullAllAlias: constants.CmdPullAll,
}

// resolveHelpCommand maps a pilot alias to its canonical command name.
func resolveHelpCommand(command string) string {
	if canonical, isAlias := helpDisplayAlias[command]; isAlias {
		return canonical
	}

	return command
}

// displayerBuilderForSub picks the sub-command builder when `--help <sub>`
// names a known sub, else the command's own builder.
func displayerBuilderForSub(entry helpDisplayEntry, args []string) func() helpdisplay.Displayer {
	sub := helpSubToken(args)
	subBuilder, isKnownSub := entry.subs[sub]
	if len(sub) > 0 && isKnownSub {
		return subBuilder
	}

	return entry.build
}

// helpSubToken returns the token following the help flag (`--help <sub>`),
// or "" when there is none.
func helpSubToken(args []string) string {
	for i, arg := range args {
		if IsHelpFlag(arg) && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// IsHelpFlag reports whether token is a help request indicator.
func IsHelpFlag(token string) bool {
	return token == "--help" || token == "-h" || token == "help"
}

// hasHelpFlag scans args for the standard help triggers.
func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if IsHelpFlag(a) {
			return true
		}
	}

	return false
}

// requireOnline checks network connectivity and exits if offline.
func requireOnline() {
	if gitutil.IsOnline() {
		return
	}

	gitutil.PrintOfflineWarning()
	cliexit.HandleGeneralError(apperror.NewSimple("network offline", "E9000"))
}
