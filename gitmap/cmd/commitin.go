package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/orchestrator"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runCommitIn is the top-level entry point for `gitmap commit-in` / `gitmap cin`.
func runCommitIn(args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			fmt.Println(commitin.PrintCommitInHelp())
			os.Exit(0)
		}
	}

	raw, perr := commitin.Parse(args)
	if perr != nil {
		fmt.Fprintf(os.Stderr, constants.CommitInErrBadArgs, perr.Message)
		os.Exit(constants.CommitInExitBadArgs)
	}
	exitCode := orchestrator.Run(raw, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
