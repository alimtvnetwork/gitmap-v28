package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/folder"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runFolder(args []string) {
	checkHelp(constants.CmdFolder, args)
	if err := folder.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
