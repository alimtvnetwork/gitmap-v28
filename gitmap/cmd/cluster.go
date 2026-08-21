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
	case "history", "hi":
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
	case "reset-password":
		runClusterResetPassword(args[1:])
		return
	case "nodes", "ls":
		runClusterNodes(args[1:])
		return
	case "remove", "rm":
		runClusterRemove(args[1:])
		return
	case "audit-clean":
		runClusterAuditClean(args[1:])
		return
	case "stats":
		runClusterStats(args[1:])
		return
	}

	fmt.Fprintf(os.Stderr, unknownCmdMsg, sub)
	os.Exit(1)
}
