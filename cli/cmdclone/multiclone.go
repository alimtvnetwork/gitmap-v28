package cmdclone

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunMultiCloneCommand executes the batch multiclone command workflow.
func RunMultiCloneCommand(args []string) error {
	isHelp := checkIsHelpRequested(args)
	if isHelp {
		RenderMultiCloneHelp()
		return nil
	}

	opts, optErr := ResolveMultiCloneOptions(args)
	if optErr != nil {
		return optErr
	}

	return processMultiCloneExecution(opts)
}

func checkIsHelpRequested(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := args[0]
	return first == "-h" || first == "--help" || first == "help"
}

func processMultiCloneExecution(opts MultiCloneOptions) error {
	urls := ParseMultiCloneText(opts.RawInput)
	hasNoURLs := len(urls) == 0
	if hasNoURLs {
		fmt.Fprint(os.Stderr, "  ⚠ No valid repository URLs or slugs found in input.\n")
		fmt.Fprint(os.Stderr, "  Run 'gitmap mc --help' to see supported input formats.\n")
		cliexit.HandleError(nil, constants.ExitCloneMultiAllInvalid)
		return nil
	}

	items := BuildMultiCloneItems(urls, opts.TargetDir)
	if opts.IsDryRun {
		RenderMultiCloneDryRun(items, opts.TargetDir)
		return nil
	}

	summary := ExecuteMultiCloneBatch(items, opts)
	if summary.Failed > 0 {
		cliexit.HandleError(nil, constants.ExitCloneMultiPartialFail)
	}

	return nil
}
