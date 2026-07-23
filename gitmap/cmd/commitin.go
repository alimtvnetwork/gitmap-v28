package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/orchestrator"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runCommitIn is the top-level entry point for `gitmap commit-in` /
// `gitmap cin`. The orchestration loop (workspace setup, walk, replay,
// runlog) lives in the commitin sub-packages; this wrapper handles
// only argv parsing and exit-code mapping.
//
// Spec: spec/03-commit-in/.
func runCommitIn(args []string) {
	raw, perr := commitin.Parse(args)
	if perr != nil {
		fmt.Fprintf(os.Stderr, constants.CommitInErrBadArgs, perr.Message)
		os.Exit(constants.CommitInExitBadArgs)
	}
	exitCode := orchestrator.Run(raw, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
