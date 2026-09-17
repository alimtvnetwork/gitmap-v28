package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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

func handleReplaceUnknown() {
	fmt.Fprint(os.Stderr, constants.ErrReplaceNeedsArgs)
	cliexit.HandleError(nil, constants.ExitCodeError)
}

// dispatchReplaceMode runs the right handler for a classified mode.
func dispatchReplaceMode(mode ReplaceModeType, positional []string, opts replaceOpts) {
	switch mode {
	case ReplaceModeTypeLiteral:
		runReplaceLiteral(positional[0], positional[1], opts)
	case ReplaceModeTypeAudit:
		runReplaceAudit(opts)
	case ReplaceModeTypeAll, ReplaceModeTypeVersionN:
		dispatchVersionMode(mode, positional, opts)
	case ReplaceModeTypeUnknown:
		handleReplaceUnknown()
	default:
		handleReplaceUnknown()
	}
}

func dispatchVersionMode(mode ReplaceModeType, positional []string, opts replaceOpts) {
	if mode == ReplaceModeTypeAll {
		runReplaceVersion(constants.ReplaceAllVersionTarget, opts, true)

		return
	}

	n := mustParseDashN(positional[0])
	runReplaceVersion(n, opts, false)
}
