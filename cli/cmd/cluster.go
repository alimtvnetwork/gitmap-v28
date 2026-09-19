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

func routeClusterBootstrapOrExec(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "bootstrap", "bs":
		return result.FailureWrapperErr(cmdssh.RunClusterBootstrapCLI(rest))
	case "exec", "run":
		return result.FailureWrapperErr(cmdssh.RunClusterExecCLI(rest))
	case "compare", "matrix":
		return result.FailureWrapperErr(cmdssh.RunSSHCompareCLI(rest))
	case "schedule", "schedules":
		return result.FailureWrapperErr(cmdssh.RunClusterExecCLI(append([]string{"schedule"}, rest...)))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterScriptOrNode(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "run-script", "script":
		return result.FailureWrapperErr(cmdssh.RunClusterScriptCLI(rest))
	case "node":
		return cmdssh.RouteClusterNodeCLI(rest)
	case "import", "import-config", "import-cluster":
		return result.FailureWrapperErr(cmdssh.RunClusterImportCLI(rest))
	case "init", "template", "init-config", "schema":
		return result.FailureWrapperErr(cmdssh.RunClusterInitCLI(rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterSSH(sub string, rest []string) result.ErrorWrapper {
	res := routeClusterBootstrapOrExec(sub, rest)
	if res.IsMatched() {
		return res
	}

	return routeClusterScriptOrNode(sub, rest)
}

func routeClusterJoinOps(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "add":
		return result.FailureWrapperErr(cmdssh.RunClusterAddCLI(rest))
	case "join":
		return result.FailureWrapperErr(cmdssh.RunClusterJoinCLI(rest))
	case "ping":
		return result.FailureWrapperErr(cmdssh.RunSJStatus(nil, rest, context.Background()))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterLegacyOps(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case constants.CmdClusterStatus:
		return result.FailureWrapperErr(runClusterStatus(rest))
	case "history", "hi":
		return result.FailureWrapperErr(runClusterHistory(rest))
	case "export":
		return result.FailureWrapperErr(runClusterExport(rest))
	case "stats":
		return result.FailureWrapperErr(runClusterStats(rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterPasswordOps(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "set-password":
		return result.FailureWrapperErr(runClusterSetPassword(rest))
	case "reset-password":
		return result.FailureWrapperErr(runClusterResetPassword(rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterNodeOps(sub string, rest []string) result.ErrorWrapper {
	resPass := routeClusterPasswordOps(sub, rest)
	if resPass.IsMatched() {
		return resPass
	}

	switch sub {
	case "nodes", "ls":
		return result.FailureWrapperErr(runClusterNodes(rest))
	case "remove", "rm":
		return result.FailureWrapperErr(runClusterRemove(rest))
	case "audit-clean":
		return result.FailureWrapperErr(runClusterAuditClean(rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterK8s(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "k8s", "kube", "kubernetes":
		return cmdssh.RouteClusterK8sCLI(rest)
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterInstall(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "install":
		return result.FailureWrapperErr(cmdssh.RunClusterInstallCLI(rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterCore(sub string, rest []string) result.ErrorWrapper {
	resInstall := routeClusterInstall(sub, rest)
	if resInstall.IsMatched() {
		return resInstall
	}

	resJoin := routeClusterJoinOps(sub, rest)
	if resJoin.IsMatched() {
		return resJoin
	}

	return routeClusterSSH(sub, rest)
}

func routeClusterExt(sub string, rest []string) result.ErrorWrapper {
	resK8s := routeClusterK8s(sub, rest)
	if resK8s.IsMatched() {
		return resK8s
	}

	resLegacy := routeClusterLegacyOps(sub, rest)
	if resLegacy.IsMatched() {
		return resLegacy
	}

	return routeClusterNodeOps(sub, rest)
}

func dispatchClusterSubcommand(sub string, rest []string) result.ErrorWrapper {
	resCore := routeClusterCore(sub, rest)
	if resCore.IsMatched() {
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

func dispatchInvertedClusterHelp(args []string) result.ErrorWrapper {
	hasInverted := len(args) > 1 && args[0] == "help"
	if !hasInverted {
		return result.UnmatchedWrapper()
	}

	res := dispatchClusterSubcommand(args[1], append(args[2:], "--help"))
	if res.IsMatched() {
		return res
	}

	return result.FailureWrapper(apperror.NewSimple("unknown command", "E9000"))
}

// routeCluster routes the cluster subcommand and returns ErrorWrapper.
func routeCluster(args []string) result.ErrorWrapper {
	if isClusterRootHelp(args) {
		helptext.PrintWithMode("cluster", render.PrettyAuto)
		return result.SuccessWrapper()
	}

	resHelp := dispatchInvertedClusterHelp(args)
	if resHelp.IsMatched() {
		return resHelp
	}

	resSub := dispatchClusterSubcommand(args[0], args[1:])
	if resSub.IsMatched() {
		return resSub
	}

	return result.FailureWrapper(apperror.NewSimple("unknown command", "E9000"))
}

// runCluster handles the "cluster" subcommand and routes to sub-handlers.
func runCluster(args []string) error {
	return result.AsError(routeCluster(args))
}
