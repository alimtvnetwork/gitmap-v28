package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	usageMsg      = "Usage: gitmap cluster [command]"
	unknownCmdMsg = "Unknown cluster command: %s\n"
)

// runCluster handles the "cluster" subcommand and routes to sub-handlers.
func runCluster(args []string) {
	checkHelp("cluster", args)
	hasNoArgs := len(args) == 0
	if hasNoArgs == true {
		fmt.Fprintln(os.Stderr, usageMsg)
		os.Exit(1)
	}

	sub := args[0]
	isStatus := sub == constants.CmdClusterStatus
	if isStatus == true {
		runClusterStatus(args[1:])
		return
	}

	fmt.Fprintf(os.Stderr, unknownCmdMsg, sub)
	os.Exit(1)
}
