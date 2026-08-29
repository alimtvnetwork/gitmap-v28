package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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
		return apperror.NewSimple(usageMsg, "E9000")
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

	return apperror.NewSimple("unknown command", "E9000")
	// 	return nil
}
