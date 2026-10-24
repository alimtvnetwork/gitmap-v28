// Package cmdbranch — `gitmap branch <subcommand>` dispatcher.
package cmdbranch

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	branchconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch/constants"
	branchenums "github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch/enums"
)

// RunBranch dispatches `gitmap branch <subcommand>` (alias `b`).
func RunBranch(args []string) *apperror.AppError {
	checkHelp(branchconst.CmdBranch, args)
	if len(args) == 0 {
		printHelpAndExit(branchconst.CmdBranch, args)

		return nil
	}

	return dispatchBranchSubcommand(args[0], args[1:])
}

func dispatchBranchSubcommand(sub string, rest []string) *apperror.AppError {
	action, _ := branchenums.ParseBranchActionType(sub)
	if action.HasDefaultAction() {
		return runBranchDefault(rest)
	}

	errMsg := fmt.Sprintf(branchconst.ErrBranchUnknownSubcommandFmt, sub)

	return apperror.NewSimple(errMsg, branchconst.ErrCodeE9000)
}

// runBranchDefault implements `gitmap branch default` / `b def`.
func runBranchDefault(args []string) *apperror.AppError {
	checkHelp(branchconst.CmdBranch, args)
	if !isInsideWorkTree() {
		fmt.Fprint(os.Stderr, branchconst.ErrBranchNotRepo)

		return apperror.NewSimple(branchconst.ErrFatalMessage, branchconst.ErrCodeE9000)
	}

	return executeBranchDefault()
}

func executeBranchDefault() *apperror.AppError {
	target := resolveDefaultBranchName(branchconst.SentinelCurrentDir)
	fmt.Printf(branchconst.MsgBranchDefaultSwitchingFmt, target)
	out, err := checkoutBranch(branchconst.SentinelCurrentDir, target)
	if len(out) > 0 {
		fmt.Println(out)
	}

	if err != nil {
		return apperror.WrapSimple(err, branchconst.ErrBranchDefaultFailed)
	}

	return nil
}
