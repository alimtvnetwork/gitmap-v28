package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runReplace is the entrypoint for `gitmap replace`. It dispatches into
// literal mode, version mode (-N / all), or audit mode based on the
// shape of args. See spec/04-generic-cli/15-replace-command.md.
func runReplace(args []string) {
	checkHelp("replace", args)

	opts, positional, err := parseReplaceFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	mode := classifyReplaceMode(positional, opts)
	dispatchReplaceMode(mode, positional, opts)
}

// dispatchReplaceMode runs the right handler for a classified mode.
func dispatchReplaceMode(mode replaceMode, positional []string, opts replaceOpts) {
	switch mode {
	case replaceModeLiteral:
		runReplaceLiteral(positional[0], positional[1], opts)
	case replaceModeAudit:
		runReplaceAudit(opts)
	case replaceModeAll:
		runReplaceVersion(0, opts, true)
	case replaceModeVersionN:
		n := mustParseDashN(positional[0])
		runReplaceVersion(n, opts, false)
	default:
		fmt.Fprint(os.Stderr, constants.ErrReplaceNeedsArgs)
		os.Exit(1)
	}
}

// replaceMode enumerates the four invocation shapes the spec accepts.
type replaceMode int

const (
	replaceModeUnknown replaceMode = iota
	replaceModeLiteral
	replaceModeVersionN
	replaceModeAll
	replaceModeAudit
)

// classifyReplaceMode picks the operating mode from positional args and
// the audit flag captured during flag parsing.
func classifyReplaceMode(positional []string, opts replaceOpts) replaceMode {
	if opts.audit {
		return replaceModeAudit
	}
	if len(positional) == 1 && positional[0] == constants.ReplaceSubcmdAll {
		return replaceModeAll
	}
	if len(positional) == 1 && looksLikeDashN(positional[0]) {
		return replaceModeVersionN
	}
	if len(positional) == 2 {
		return replaceModeLiteral
	}
	return replaceModeUnknown
}

// looksLikeDashN matches strings of the form "-1", "-23", etc.
func looksLikeDashN(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// mustParseDashN converts "-N" to the integer N. Caller guarantees
// looksLikeDashN(s) is true; we still bail on overflow.
func mustParseDashN(s string) int {
	n := 0
	for i := 1; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
		if n > 1_000_000 {
			fmt.Fprintf(os.Stderr, constants.ErrReplaceBadN, s)
			os.Exit(1)
		}
	}
	if n < 1 {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceBadN, s)
		os.Exit(1)
	}
	return n
}
