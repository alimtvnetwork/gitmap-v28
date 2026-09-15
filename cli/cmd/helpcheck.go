package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
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
	_, mode := ParsePrettyFlag(args)
	helptext.PrintWithMode(command, mode)
	printUsageFooterShort()
	cliexit.Exit(0)
}

// hasHelpFlag scans args for the standard help triggers.
func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
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
