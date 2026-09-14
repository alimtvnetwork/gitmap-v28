package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	usageMsg      = "Usage: gitmap cluster [command]"
	unknownCmdMsg = "Unknown cluster command: %s\n"
)

func routeClusterSSH(sub string, rest []string) (error, bool) {
	switch sub {
	case "exec", "run":
		return cmdssh.RunClusterExecCLI(rest), true
	case "run-script", "script":
		return cmdssh.RunClusterScriptCLI(rest), true
	case "node":
		return cmdssh.RunClusterNodeCLI(rest), true
	case "import", "import-config", "import-cluster":
		return cmdssh.RunClusterImportCLI(rest), true
	default:
		return nil, false
	}
}

func routeClusterLegacyOps(sub string, rest []string) (error, bool) {
	switch sub {
	case constants.CmdClusterStatus:
		return runClusterStatus(rest), true
	case "history", "hi":
		return runClusterHistory(rest), true
	case "export":
		return runClusterExport(rest), true
	case "stats":
		return runClusterStats(rest), true
	default:
		return nil, false
	}
}

func routeClusterNodeOps(sub string, rest []string) (error, bool) {
	switch sub {
	case "nodes", "ls":
		return runClusterNodes(rest), true
	case "remove", "rm":
		return runClusterRemove(rest), true
	case "set-password":
		return runClusterSetPassword(rest), true
	case "reset-password":
		return runClusterResetPassword(rest), true
	case "audit-clean":
		return runClusterAuditClean(rest), true
	default:
		return nil, false
	}
}

func routeClusterK8s(sub string, rest []string) (error, bool) {
	switch sub {
	case "k8s", "kube", "kubernetes":
		return cmdssh.RunClusterK8sCLI(rest), true
	default:
		return nil, false
	}
}

func dispatchClusterSubcommand(sub string, rest []string) (error, bool) {
	if err, isSSH := routeClusterSSH(sub, rest); isSSH {
		return err, true
	}
	if err, isK8s := routeClusterK8s(sub, rest); isK8s {
		return err, true
	}
	if err, isLegacy := routeClusterLegacyOps(sub, rest); isLegacy {
		return err, true
	}
	return routeClusterNodeOps(sub, rest)
}

// runCluster handles the "cluster" subcommand and routes to sub-handlers.
func runCluster(args []string) error {
	checkHelp("cluster", args)
	isEmptyArgs := len(args) == 0
	if isEmptyArgs {
		return apperror.NewSimple(usageMsg, "E9000")
	}
	err, isMatched := dispatchClusterSubcommand(args[0], args[1:])
	if isMatched {
		return err
	}

	return apperror.NewSimple("unknown command", "E9000")
}
