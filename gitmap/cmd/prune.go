package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runPrune handles the prune command.
func runPrune(args []string) {
	checkHelp("prune", args)
	dryRun, confirm, remote := parsePruneFlags(args)

	branches := listReleaseBranches()
	stale := filterStaleBranches(branches)

	if len(stale) == 0 {
		fmt.Print(constants.MsgPruneNone)

		return
	}

	printStaleBranches(stale)

	if dryRun {
		fmt.Print(constants.MsgPruneDryRunHint)

		return
	}

	if shouldProceedWithPrune(confirm, len(stale)) {
		deleteStaleBranches(stale, remote)
	}
}

// parsePruneFlags parses prune command flags.
func parsePruneFlags(args []string) (bool, bool, bool) {
	fs := flag.NewFlagSet(constants.CmdPrune, flag.ExitOnError)
	dryRun := fs.Bool(constants.PruneFlagDryRun, false, constants.FlagDescDryRun)
	confirm := fs.Bool(constants.PruneFlagConfirm, false, constants.FlagDescConfirm)
	remote := fs.Bool(constants.PruneFlagRemote, false, "Also delete remote release branches")
	fs.Parse(args)

	return *dryRun, *confirm, *remote
}

// shouldProceedWithPrune checks confirmation or prompts the user.
func shouldProceedWithPrune(confirmed bool, count int) bool {
	if confirmed {
		return true
	}

	fmt.Printf(constants.MsgPrunePrompt, count)
	var answer string
	fmt.Scanln(&answer)

	if answer == "y" || answer == "Y" {
		return true
	}

	fmt.Print(constants.MsgPruneAborted)

	return false
}

// deleteStaleBranches deletes each stale branch and prints a summary.
func deleteStaleBranches(stale []staleBranch, remote bool) {
	fmt.Print(constants.MsgPruneDeleting)
	deleted := 0

	for _, sb := range stale {
		deleted += deleteSingleBranch(sb, remote)
	}

	kept := len(stale) - deleted
	fmt.Printf(constants.MsgPruneSummary, deleted, kept)
}

// deleteSingleBranch deletes one branch locally and optionally remotely.
func deleteSingleBranch(sb staleBranch, remote bool) int {
	err := deleteLocalBranch(sb.name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPruneDeleteBranch, sb.name, err)

		return 0
	}

	fmt.Printf(constants.MsgPruneDeleted, sb.name)

	if remote {
		deleteRemoteBranch(sb.name)
	}

	return 1
}
