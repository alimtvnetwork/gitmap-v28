package cmd

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func routeClusterBootstrapOrExec(sub string, rest []string) result.Result[bool] {
	switch sub {
	case "bootstrap", "bs":
		return result.RouteMatched(cmdssh.RunClusterBootstrapCLI(rest))
	case "exec", "run":
		return result.RouteMatched(cmdssh.RunClusterExecCLI(rest))
	default:
		return result.RouteUnmatched()
	}
}

func routeClusterScriptOrNode(sub string, rest []string) result.Result[bool] {
	switch sub {
	case "run-script", "script":
		return result.RouteMatched(cmdssh.RunClusterScriptCLI(rest))
	case "node":
		return result.RouteMatched(cmdssh.RunClusterNodeCLI(rest))
	case "import", "import-config", "import-cluster":
		return result.RouteMatched(cmdssh.RunClusterImportCLI(rest))
	default:
		return result.RouteUnmatched()
	}
}

func routeClusterSSH(sub string, rest []string) result.Result[bool] {
	res := routeClusterBootstrapOrExec(sub, rest)
	if res.Data {
		return res
	}

	return routeClusterScriptOrNode(sub, rest)
}

func routeClusterJoinOps(sub string, rest []string) result.Result[bool] {
	switch sub {
	case "add":
		return result.RouteMatched(cmdssh.RunClusterAddCLI(rest))
	case "join":
		return result.RouteMatched(cmdssh.RunClusterJoinCLI(rest))
	case "ping":
		return result.RouteMatched(cmdssh.RunSJStatus(nil, rest, context.Background()))
	default:
		return result.RouteUnmatched()
	}
}

func routeClusterLegacyOps(sub string, rest []string) result.Result[bool] {
	switch sub {
	case constants.CmdClusterStatus:
		return result.RouteMatched(runClusterStatus(rest))
	case "history", "hi":
		return result.RouteMatched(runClusterHistory(rest))
	case "export":
		return result.RouteMatched(runClusterExport(rest))
	case "stats":
		return result.RouteMatched(runClusterStats(rest))
	default:
		return result.RouteUnmatched()
	}
}

func routeClusterPasswordOps(sub string, rest []string) result.Result[bool] {
	switch sub {
	case "set-password":
		return result.RouteMatched(runClusterSetPassword(rest))
	case "reset-password":
		return result.RouteMatched(runClusterResetPassword(rest))
	default:
		return result.RouteUnmatched()
	}
}

func routeClusterNodeOps(sub string, rest []string) result.Result[bool] {
	resPass := routeClusterPasswordOps(sub, rest)
	if resPass.Data {
		return resPass
	}

	switch sub {
	case "nodes", "ls":
		return result.RouteMatched(runClusterNodes(rest))
	case "remove", "rm":
		return result.RouteMatched(runClusterRemove(rest))
	case "audit-clean":
		return result.RouteMatched(runClusterAuditClean(rest))
	default:
		return result.RouteUnmatched()
	}
}

func routeClusterK8s(sub string, rest []string) result.Result[bool] {
	switch sub {
	case "k8s", "kube", "kubernetes":
		return result.RouteMatched(cmdssh.RunClusterK8sCLI(rest))
	default:
		return result.RouteUnmatched()
	}
}

func routeClusterCore(sub string, rest []string) result.Result[bool] {
	resJoin := routeClusterJoinOps(sub, rest)
	if resJoin.Data {
		return resJoin
	}

	return routeClusterSSH(sub, rest)
}

func routeClusterExt(sub string, rest []string) result.Result[bool] {
	resK8s := routeClusterK8s(sub, rest)
	if resK8s.Data {
		return resK8s
	}

	resLegacy := routeClusterLegacyOps(sub, rest)
	if resLegacy.Data {
		return resLegacy
	}

	return routeClusterNodeOps(sub, rest)
}

func dispatchClusterSubcommand(sub string, rest []string) result.Result[bool] {
	resCore := routeClusterCore(sub, rest)
	if resCore.Data {
		return resCore
	}

	return routeClusterExt(sub, rest)
}

func isClusterHelpToken(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}

func isClusterRootHelp(args []string) bool {
	isZero := len(args) == 0
	isSingleHelp := len(args) == 1 && isClusterHelpToken(args[0])

	return isZero || isSingleHelp
}

func dispatchInvertedClusterHelp(args []string) result.Result[bool] {
	hasInverted := len(args) > 1 && args[0] == "help"
	if !hasInverted {
		return result.RouteUnmatched()
	}

	res := dispatchClusterSubcommand(args[1], append(args[2:], "--help"))
	if res.Data {
		return res
	}

	return result.RouteMatchedAppErr(apperror.NewSimple("unknown command", "E9000"))
}

// runCluster handles the "cluster" subcommand and routes to sub-handlers.
func runCluster(args []string) error {
	if isClusterRootHelp(args) {
		helptext.PrintWithMode("cluster", render.PrettyAuto)
		return nil
	}

	resHelp := dispatchInvertedClusterHelp(args)
	if resHelp.Data {
		return resHelp.AppError()
	}

	resSub := dispatchClusterSubcommand(args[0], args[1:])
	if resSub.Data {
		return resSub.AppError()
	}

	return apperror.NewSimple("unknown command", "E9000")
}
