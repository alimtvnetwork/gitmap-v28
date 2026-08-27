package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/add"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runAdd(args []string) {
	checkHelp(constants.CmdAdd, args)
	if err := add.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
