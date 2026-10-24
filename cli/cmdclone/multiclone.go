package cmdclone

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	cloneconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdclone/constants"
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
	hasURLs := len(urls) > 0
	if !hasURLs {
		fmt.Fprint(os.Stderr, cloneconstants.MsgMultiCloneNoURLs)
		return apperror.NewValidationError("no valid repository URLs or slugs found in input")
	}

	items := BuildMultiCloneItems(urls, opts.TargetDir)
	if opts.IsDryRun {
		RenderMultiCloneDryRun(items, opts.TargetDir)

		return nil
	}

	summary := ExecuteMultiCloneBatch(items, opts)
	if summary.Failed > 0 {
		return apperror.NewWithDetails("clone_multi", "MULTI_CLONE_PARTIAL_FAIL", fmt.Sprintf("%d repositories failed to clone", summary.Failed), "processMultiCloneExecution", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}
