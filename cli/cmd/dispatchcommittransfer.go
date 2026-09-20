package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// dispatchCommitTransfer routes commit-left / commit-right / commit-both
// (and their aliases cml / cmr / cmb), as well as pr / pull-request / pr-clean / pr-rm / pr-list.
//
// Spec: 02-spec/01-app/106-commit-left-right-both.md
// Spec: 02-spec/21-app/129-pr-commit-engines-and-sqlite-split-db.md
func dispatchCommitTransfer(command string) (bool, error) {
	if isPRCleanCommand(command) {
		return true, runPRClean(argsTail())
	}

	if isPRListCommand(command) {
		return true, runPRList(argsTail())
	}

	spec, ok := commitTransferSpecFor(command)
	if !ok {
		return false, nil
	}

	return dispatchCommitTransferSpec(spec, argsTail())
}

func isPRCleanCommand(command string) bool {
	return command == constants.CmdPRClean || command == constants.CmdPRRm
}

func isPRListCommand(command string) bool {
	return command == constants.CmdPRList
}

func dispatchCommitTransferSpec(spec commitTransferSpec, args []string) (bool, error) {
	if isPRSubcommandClean(spec.Name, args) {
		return true, runPRClean(args[1:])
	}

	if isPRSubcommandList(spec.Name, args) {
		return true, runPRList(args[1:])
	}

	runCommitTransfer(spec, args)

	return true, nil
}

func isPRSubcommandClean(name string, args []string) bool {
	if name != constants.CmdPR && name != constants.CmdPullRequest {
		return false
	}
	if len(args) == 0 {
		return false
	}

	return args[0] == "clean" || args[0] == "rm"
}

func isPRSubcommandList(name string, args []string) bool {
	if name != constants.CmdPR && name != constants.CmdPullRequest {
		return false
	}
	if len(args) == 0 {
		return false
	}

	return args[0] == "list" || args[0] == "ls"
}
