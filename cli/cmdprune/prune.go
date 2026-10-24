package cmdprune

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune/enums"
)

var promptReaderFn = readPromptInput

// RunPrune handles the prune command.
func RunPrune(args []string) error {
	checkHelp(constants.CmdPrune, args)
	isDryRun, isConfirmed, isRemote := parsePruneFlags(args)
	stale := filterStaleBranches(listBranchesFn())
	if len(stale) == 0 {
		fmt.Print(constants.MsgPruneNone)

		return nil
	}

	executePruneWorkflow(stale, isDryRun, isConfirmed, isRemote)

	return nil
}

func executePruneWorkflow(stale []staleBranch, isDryRun, isConfirmed, isRemote bool) {
	printStaleBranches(stale)
	action := resolvePruneAction(isDryRun, isConfirmed)
	if action.IsDryRun() {
		fmt.Print(constants.MsgPruneDryRunHint)

		return
	}

	if isPruneConfirmed(isConfirmed, len(stale)) {
		deleteStaleBranches(stale, isRemote)
	}
}

func resolvePruneAction(isDryRun, isConfirmed bool) enums.PruneActionType {
	if isDryRun {
		return enums.PruneActionDryRun
	}
	if isConfirmed {
		return enums.PruneActionDelete
	}

	return enums.PruneActionPrompt
}

func resolvePruneScope(isRemote bool) enums.PruneScopeType {
	if isRemote {
		return enums.PruneScopeAll
	}

	return enums.PruneScopeLocal
}

// parsePruneFlags parses prune command flags.
func parsePruneFlags(args []string) (bool, bool, bool) {
	fs := flag.NewFlagSet(constants.CmdPrune, flag.ExitOnError)
	isDryRun := fs.Bool(constants.FlagDryRun, false, constants.FlagDescDryRun)
	isConfirmed := fs.Bool(constants.FlagConfirm, false, constants.FlagDescConfirm)
	isRemote := fs.Bool(constants.FlagRemote, false, constants.FlagDescRemote)
	fs.Parse(args)

	return *isDryRun, *isConfirmed, *isRemote
}

// isPruneConfirmed checks confirmation or prompts the user.
func isPruneConfirmed(isConfirmed bool, count int) bool {
	if isConfirmed {
		return true
	}
	if promptPruneConfirmation(count) {
		return true
	}

	fmt.Print(constants.MsgPruneAborted)

	return false
}

func readPromptInput() string {
	var answer string
	fmt.Scanln(&answer)

	return answer
}

func promptPruneConfirmation(count int) bool {
	fmt.Printf(constants.MsgPrunePrompt, count)
	answer := promptReaderFn()

	return strings.EqualFold(answer, constants.PromptAnswerY)
}

// deleteStaleBranches deletes each stale branch and prints a summary.
func deleteStaleBranches(stale []staleBranch, isRemote bool) {
	fmt.Print(constants.MsgPruneDeleting)
	deleted := 0
	for _, sb := range stale {
		deleted += deleteSingleBranch(sb, isRemote)
	}

	kept := len(stale) - deleted
	fmt.Printf(constants.MsgPruneSummary, deleted, kept)
}

// deleteSingleBranch deletes one branch locally and optionally remotely.
func deleteSingleBranch(sb staleBranch, isRemote bool) int {
	if err := deleteLocalBranchFn(sb.name); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPruneDeleteBranch, sb.name, err)

		return 0
	}

	fmt.Printf(constants.MsgPruneDeleted, sb.name)
	pruneRemoteIfRequested(sb.name, isRemote)

	return 1
}

func pruneRemoteIfRequested(name string, isRemote bool) {
	if isRemote {
		deleteRemoteBranchFn(name)
	}
}
