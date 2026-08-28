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
func runCluster(args []string) error {
	checkHelp("cluster", args)
	isEmptyArgs := len(args) == 0
	if isEmptyArgs {
		fmt.Fprintln(os.Stderr, usageMsg)
		os.Exit(1)
	}

	sub := args[0]
	switch sub {
	case constants.CmdClusterStatus:
		runClusterStatus(args[1:])
		return nil
	case "history", "hi":
		runClusterHistory(args[1:])
		return nil
	case "export":
		runClusterExport(args[1:])
		return nil
	case "import":
		runClusterImport(args[1:])
		return nil
	case "set-password":
		runClusterSetPassword(args[1:])
		return nil
	case "reset-password":
		runClusterResetPassword(args[1:])
		return nil
	case "nodes", "ls":
		runClusterNodes(args[1:])
		return nil
	case "remove", "rm":
		runClusterRemove(args[1:])
		return nil
	case "audit-clean":
		runClusterAuditClean(args[1:])
		return nil
	case "stats":
		runClusterStats(args[1:])
		return nil
	}

	fmt.Fprintf(os.Stderr, unknownCmdMsg, sub)
	os.Exit(1)
	return nil
}
