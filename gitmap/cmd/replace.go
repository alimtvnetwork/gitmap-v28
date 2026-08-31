package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runReplace is the entrypoint for `gitmap replace`. It dispatches into
// literal mode, version mode (-N / all), or audit mode based on the
// shape of args. See spec/04-generic-cli/15-replace-command.md.
func runReplace(args []string) error {
	checkHelp(constants.CmdReplace, args)

	opts, positional, err := parseReplaceFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}

	mode := classifyReplaceMode(positional, opts)
	dispatchReplaceMode(mode, positional, opts)
	return nil
}

// dispatchReplaceMode runs the right handler for a classified mode.
func dispatchReplaceMode(mode replaceMode, positional []string, opts replaceOpts) {
	switch mode {
	case replaceModeLiteral:
		runReplaceLiteral(positional[0], positional[1], opts)
	case replaceModeAudit:
		runReplaceAudit(opts)
	case replaceModeAll, replaceModeVersionN:
		dispatchVersionMode(mode, positional, opts)
	case replaceModeUnknown:
		fmt.Fprint(os.Stderr, constants.ErrReplaceNeedsArgs)
		cliexit.HandleError(nil, constants.ExitCodeError)
	default:
		fmt.Fprint(os.Stderr, constants.ErrReplaceNeedsArgs)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
}

func dispatchVersionMode(mode replaceMode, positional []string, opts replaceOpts) {
	if mode == replaceModeAll {
		runReplaceVersion(constants.ReplaceAllVersionTarget, opts, true)
		return
	}
	n := mustParseDashN(positional[0])
	runReplaceVersion(n, opts, false)
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
	if len(positional) == constants.ReplaceMinPositionalArgs {
		return classifySingleArgMode(positional[0])
	}
	if len(positional) == constants.ReplaceLiteralArgsCount {
		return replaceModeLiteral
	}
	return replaceModeUnknown
}

func classifySingleArgMode(arg string) replaceMode {
	if arg == constants.ReplaceSubcmdAll {
		return replaceModeAll
	}
	if arg == "history" || arg == "audit" {
		return replaceModeAudit
	}
	if looksLikeDashN(arg) {
		return replaceModeVersionN
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
	n, ok := parseDashNDigits(s)
	if !ok || n < 1 {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceBadN, s)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
	return n
}

func parseDashNDigits(s string) (int, bool) {
	n := 0
	for i := 1; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
		if n > constants.ReplaceMaxDashN {
			return 0, false
		}
	}
	return n, true
}
