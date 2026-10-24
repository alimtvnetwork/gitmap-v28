package cmd

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/cluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdnodes"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservercmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func coreClusterEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"nodes", "node", "allnodes", "all-nodes", "fleet-nodes"}, func() error { return cmdnodes.RunUnifiedNodesCLI(argsTail()) }},
		{[]string{"nodes-running-only", "nodesrunningonly", "nodes-running", "nro", "nodes-online", "nodesonline", "nodes-open", "nodes-active"}, func() error { return cmdnodes.RunNodesRunningOnlyCLI(argsTail()) }},
		{[]string{"nodes-clone", "node-clone", "fleet-clone"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"clone"}, argsTail()...)) }},
		{[]string{"nodes-cfr", "node-cfr", "fleet-cfr"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"cfr"}, argsTail()...)) }},
		{[]string{"nodes-cfrp", "node-cfrp", "fleet-cfrp"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"cfrp"}, argsTail()...)) }},
		{[]string{"ping", "nodes-ping", "nodeping", "fleet-ping"}, func() error { return cmdnodes.RunUnifiedNodesPingCLI(argsTail()) }},
		{[]string{"nodes-history", "nodes-histories", "node-history"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"history"}, argsTail()...)) }},
		{[]string{"history-ssh"}, func() error { return cmdhistory.RunHistory([]string{"ssh"}) }},
		{[]string{"nodes-agy-ui", "nodes-agy", "nodes-ui", "agy-ui"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"agy", "ui"}, argsTail()...)) }},
		{[]string{"nodes-push-settings", "push-settings"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"push-settings"}, argsTail()...)) }},
		{[]string{"nodes-sync-settings", "sync-settings"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"sync-settings"}, argsTail()...)) }},
		{[]string{"nodes-send-projects", "send-projects", "sync-projects"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"send-projects"}, argsTail()...)) }},
		{[]string{"nodes-deploy-agm-accounts", "sync-agm-accounts", "deploy-agm-accounts"}, func() error {
			return cmdnodes.RunUnifiedNodesCLI(append([]string{"deploy", "agm-accounts"}, argsTail()...))
		}},
		{[]string{"nodes-deploy-repo", "deploy-repo", "node-deploy-repo"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"deploy", "repo"}, argsTail()...)) }},
		{[]string{"nodes-deploy-repos", "deploy-repos", "node-deploy-repos"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"deploy", "repos"}, argsTail()...)) }},
		{[]string{"nodes-scan", "fleet-scan", "node-scan"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"scan"}, argsTail()...)) }},
		{[]string{"nodes-rescan", "fleet-rescan", "node-rescan"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"rescan"}, argsTail()...)) }},
		{[]string{constants.CmdServersClients, constants.CmdServersClientsAlias, constants.CmdSC}, func() error { dispatchServersClients(argsTail()); return nil }},
		{[]string{constants.CmdClients, constants.CmdClientsAlias}, func() error { dispatchClients(argsTail()); return nil }},
		{[]string{"servers"}, func() error { dispatchServers(argsTail()); return nil }},
		{[]string{constants.CmdCluster, constants.CmdClusterAlias}, func() error { return cmdcluster.RunCluster(argsTail()) }},
		{[]string{constants.CmdServerCmd, constants.CmdServerCmds, constants.CmdServerCmdAlias}, func() error { return cmdservercmd.RunServerCmd(argsTail()) }},
	}
}

func dispatchServersClients(args []string) {
	if IsHelpRequestedOrEmpty(args) {
		cmdschedule.RenderSCHelp()

		return
	}

	subCmd, rest := args[0], args[1:]
	isHandled := dispatchSCNodeOps(subCmd, rest)
	if isHandled {
		return
	}
	isCustom := dispatchSCCustomOps(subCmd, rest)
	if isCustom {
		return
	}

	cmdcluster.RunClusterCommand(cluster.ServersClients, args)
}

func dispatchSCCustomOps(subCmd string, rest []string) bool {
	isPath := dispatchServersClientsPathCmd(subCmd, rest)
	if isPath {
		return true
	}
	isRW := dispatchClusterReadWrite(cluster.ServersClients, subCmd, rest)
	if isRW {
		return true
	}

	return dispatchClusterMutate(cluster.ServersClients, subCmd, rest)
}

func dispatchSCClusterRecipes(subCmd string, rest []string) bool {
	switch subCmd {
	case "node":
		_ = cmdssh.RouteClusterNodeCLI(rest)
		return true
	case "bootstrap", "bs":
		_ = cmdssh.RunClusterBootstrapCLI(rest)
		return true
	case "k8s", "kube", "kubernetes":
		_ = cmdssh.RouteClusterK8sCLI(rest)
		return true
	default:
		return false
	}
}

func dispatchSCExecOps(subCmd string, rest []string) bool {
	switch subCmd {
	case "exec", "run":
		_ = cmdssh.RunClusterExecCLI(rest)
		return true
	case "schedule", "schedules":
		_ = cmdssh.RunClusterExecCLI(append([]string{"schedule"}, rest...))
		return true
	default:
		return false
	}
}

func dispatchSCNetworkOps(subCmd string, rest []string) bool {
	switch subCmd {
	case "compare", "matrix":
		_ = cmdssh.RunSSHCompareCLI(rest)
		return true
	case "join", "add", "enroll":
		_ = cmdssh.RunClusterJoinCLI(rest)
		return true
	case "ping", "health":
		_ = cmdssh.RunSJStatus(nil, rest, context.Background())
		return true
	default:
		return false
	}
}

func dispatchSCAuthOps(subCmd string, rest []string) bool {
	switch subCmd {
	case "auth-key", "copy-id", "fix-auth":
		_ = cmdssh.RunSSHAuthKeyDeployCLI(rest)
		return true
	default:
		return false
	}
}

func dispatchSCMetaOps(subCmd string, rest []string) bool {
	if dispatchSCExecOps(subCmd, rest) {
		return true
	}
	if dispatchSCNetworkOps(subCmd, rest) {
		return true
	}
	if dispatchSCDeployOps(subCmd, rest) {
		return true
	}
	return dispatchSCAuthOps(subCmd, rest)
}

func dispatchSCDeployOps(subCmd string, rest []string) bool {
	switch subCmd {
	case "deploy", "deploy-right", "deploy-left":
		_ = cmdssh.RunSSHDeployCLI(subCmd, rest)
		return true
	case "deploy-config", "deploy-config-ssh", "deploy-ssh-config":
		_ = cmdssh.RunSSHDeployConfigSSHCLI(rest)
		return true
	default:
		return false
	}
}

func dispatchSCInventoryOps(subCmd string, rest []string) bool {
	switch subCmd {
	case "nodes", "list", "machines", "joined":
		_ = cmdcluster.RunClusterNodes(rest)
		return true
	case "rm", "remove", "delete":
		_ = cmdcluster.RunClusterRemove(rest)
		return true
	default:
		return false
	}
}

func dispatchSCNodeOps(subCmd string, rest []string) bool {
	isMeta := dispatchSCMetaOps(subCmd, rest)
	if isMeta {
		return true
	}
	isRecipe := dispatchSCClusterRecipes(subCmd, rest)
	if isRecipe {
		return true
	}

	return dispatchSCInventoryOps(subCmd, rest)
}

func dispatchServersClientsPathCmd(subCmd string, rest []string) bool {
	switch subCmd {
	case "set-default-path":
		cmdcluster.RunClusterSetDefaultPath(cluster.ServersClients, rest)

		return true
	case "set-path-alias":
		cmdcluster.RunClusterSetPathAlias(cluster.ServersClients, rest)

		return true
	default:
		return false
	}
}

func checkClusterLSHelp(selector cluster.TargetSelectorType, rest []string) {
	if selector == cluster.ClientsOnly {
		checkHelp(constants.CmdClientsLS, rest)
		return
	}

	checkHelp(constants.CmdSCLS, rest)
}

func dispatchClusterLS(selector cluster.TargetSelectorType, rest []string) {
	if hasHelpFlag(rest) {
		checkClusterLSHelp(selector, rest)
		return
	}

	_ = cmdcluster.RunClusterNodes(rest)
}

func dispatchClusterReadWrite(
	selector cluster.TargetSelectorType,
	subCmd string,
	rest []string,
) bool {
	switch subCmd {
	case "ls":
		dispatchClusterLS(selector, rest)
	case "cat":
		cmdcluster.RunClusterCat(selector, rest)
	case "write":
		cmdcluster.RunClusterWrite(selector, rest)
	default:
		return false
	}

	return true
}

func dispatchClusterMutate(selector cluster.TargetSelectorType, subCmd string, rest []string) bool {
	switch subCmd {
	case "update":
		cmdcluster.RunClusterUpdate(selector, false, rest)
	case "update-all":
		cmdcluster.RunClusterUpdate(selector, true, rest)
	case "clone", "cfr", "cfrp":
		cmdcluster.RunClusterClone(selector, subCmd, rest)
	default:
		return false
	}

	return true
}

func dispatchClientsCustom(subCmd string, rest []string) bool {
	isRW := dispatchClusterReadWrite(cluster.ClientsOnly, subCmd, rest)
	if isRW {
		return true
	}

	return dispatchClusterMutate(cluster.ClientsOnly, subCmd, rest)
}

func dispatchClients(args []string) {
	CheckHelpOrEmpty(constants.CmdClients, args)

	subCmd, rest := args[0], args[1:]
	isNodeOp := dispatchSCNodeOps(subCmd, rest)
	if isNodeOp {
		return
	}
	isCustom := dispatchClientsCustom(subCmd, rest)
	if isCustom {
		return
	}

	cmdcluster.RunClusterCommand(cluster.ClientsOnly, args)
}

func dispatchServers(args []string) {
	if len(args) == 0 {
		return
	}

	if args[0] == "ls" {
		cmdcluster.RunClusterLS(cluster.ServersOnly, args[1:])

		return
	}

	dispatchServersUpdate(args[0], args[1:])
}

func dispatchServersUpdate(subCmd string, rest []string) {
	if subCmd == "update" {
		cmdcluster.RunClusterUpdate(cluster.ServersOnly, false, rest)
	} else if subCmd == "update-all" {
		cmdcluster.RunClusterUpdate(cluster.ServersOnly, true, rest)
	}
}
