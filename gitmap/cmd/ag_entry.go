package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/ag"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runAg(args []string) {
	checkHelp(constants.CmdAg, args)
	if err := ag.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
