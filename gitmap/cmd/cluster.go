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
	switch sub {
	case constants.CmdClusterStatus:
		runClusterStatus(args[1:])
		return
	case "history":
		runClusterHistory(args[1:])
		return
	case "export":
		runClusterExport(args[1:])
		return
	case "import":
		runClusterImport(args[1:])
		return
	case "set-password":
		runClusterSetPassword(args[1:])
		return
	}

	fmt.Fprintf(os.Stderr, unknownCmdMsg, sub)
	os.Exit(1)
}

func runClusterHistory(args []string) {
	fmt.Printf("runClusterHistory(args=%v)\n", args)
}

func runClusterExport(args []string) {
	fmt.Printf("runClusterExport(args=%v)\n", args)
}

func runClusterImport(args []string) {
	fmt.Printf("runClusterImport(args=%v)\n", args)
}

func runClusterSetPassword(args []string) {
	fmt.Printf("runClusterSetPassword(args=%v)\n", args)
}
