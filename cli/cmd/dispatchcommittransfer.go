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

	if command == constants.CmdPRIn || command == constants.CmdCommitInPR || command == constants.CmdCommitInPRA {
		return true, runCommitIn(append([]string{"--pr=feature-per-commit"}, argsTail()...))
	}

	if command == "commit-pull" || command == "cpull" || command == "pull-commits" {
		return true, runCommitPull(argsTail())
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

func dispatchPRSubcommands(spec *commitTransferSpec, args []string) (bool, error) {
	if isPRSubcommandClean(spec.Name, args) {
		return true, runPRClean(args[1:])
	}
	if isPRSubcommandList(spec.Name, args) {
		return true, runPRList(args[1:])
	}
	if isPRSubcommandIn(spec.Name, args) {
		return true, runCommitIn(args[1:])
	}

	return dispatchPRDirection(spec, args)
}

func dispatchPRDirection(spec *commitTransferSpec, args []string) (bool, error) {
	if isPRSubcommandLeft(spec.Name, args) {
		spec.Name = constants.CmdCommitLeft
		return true, runCommitTransfer(*spec, args[1:])
	}
	if isPRSubcommandRight(spec.Name, args) {
		spec.Name = constants.CmdCommitRight
		return true, runCommitTransfer(*spec, args[1:])
	}

	return false, nil
}

func dispatchCommitTransferSpec(spec commitTransferSpec, args []string) (bool, error) {
	if handled, err := dispatchPRSubcommands(&spec, args); handled {
		return true, err
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

func isPRSubcommandIn(name string, args []string) bool {
	if name != constants.CmdPR && name != constants.CmdPullRequest {
		return false
	}
	if len(args) == 0 {
		return false
	}

	return args[0] == "in"
}

func isPRSubcommandLeft(name string, args []string) bool {
	if name != constants.CmdPR && name != constants.CmdPullRequest {
		return false
	}
	if len(args) == 0 {
		return false
	}

	return args[0] == "left"
}

func isPRSubcommandRight(name string, args []string) bool {
	if name != constants.CmdPR && name != constants.CmdPullRequest {
		return false
	}
	if len(args) == 0 {
		return false
	}

	return args[0] == "right"
}
