package cmd

import (
	"os"
)

// dispatchCommitTransfer routes commit-left / commit-right / commit-both
// (and their aliases cml / cmr / cmb).
//
// Spec: spec/01-app/106-commit-left-right-both.md
func dispatchCommitTransfer(command string) (bool, error) {
	spec, ok := commitTransferSpecFor(command)
	if !ok {
		return false, nil
	}
	runCommitTransfer(spec, os.Args[2:])

	return true, nil
}
