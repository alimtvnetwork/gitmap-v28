package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/gitrm"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runGitRm(args []string) {
	checkHelp(constants.CmdGitRm, args)
	if err := gitrm.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
