package cmd

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
)

const (
	usageMsg      = "Usage: gitmap cluster [command]"
	unknownCmdMsg = "Unknown cluster command: %s\n"
)

func routeClusterBootstrapOrExec(sub string, rest []string) (error, bool) {
	switch sub {
	case "bootstrap", "bs":
		return cmdssh.RunClusterBootstrapCLI(rest), true
	case "exec", "run":
		return cmdssh.RunClusterExecCLI(rest), true
	default:
		return nil, false
	}
}

func routeClusterScriptOrNode(sub string, rest []string) (error, bool) {
	switch sub {
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

func routeClusterSSH(sub string, rest []string) (error, bool) {
	if err, isMatched := routeClusterBootstrapOrExec(sub, rest); isMatched {
		return err, true
	}
	return routeClusterScriptOrNode(sub, rest)
}

func routeClusterJoinOps(sub string, rest []string) (error, bool) {
	switch sub {
	case "add":
		return cmdssh.RunClusterAddCLI(rest), true
	case "join":
		return cmdssh.RunClusterJoinCLI(rest), true
	case "ping":
		return cmdssh.RunSJStatus(nil, rest, context.Background()), true
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

func routeClusterPasswordOps(sub string, rest []string) (error, bool) {
	switch sub {
	case "set-password":
		return runClusterSetPassword(rest), true
	case "reset-password":
		return runClusterResetPassword(rest), true
	default:
		return nil, false
	}
}

func routeClusterNodeOps(sub string, rest []string) (error, bool) {
	if err, isMatched := routeClusterPasswordOps(sub, rest); isMatched {
		return err, true
	}
	switch sub {
	case "nodes", "ls":
		return runClusterNodes(rest), true
	case "remove", "rm":
		return runClusterRemove(rest), true
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

func routeClusterCore(sub string, rest []string) (error, bool) {
	if err, isJoin := routeClusterJoinOps(sub, rest); isJoin {
		return err, true
	}
	return routeClusterSSH(sub, rest)
}

func routeClusterExt(sub string, rest []string) (error, bool) {
	if err, isK8s := routeClusterK8s(sub, rest); isK8s {
		return err, true
	}
	if err, isLegacy := routeClusterLegacyOps(sub, rest); isLegacy {
		return err, true
	}
	return routeClusterNodeOps(sub, rest)
}

func dispatchClusterSubcommand(sub string, rest []string) (error, bool) {
	if err, isCore := routeClusterCore(sub, rest); isCore {
		return err, true
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

func dispatchInvertedClusterHelp(args []string) (error, bool) {
	hasInverted := args[0] == "help" && len(args) > 1
	if hasInverted {
		err, isMatched := dispatchClusterSubcommand(args[1], append(args[2:], "--help"))
		if isMatched {
			return err, true
		}
		return apperror.NewSimple("unknown command", "E9000"), true
	}
	return nil, false
}

// runCluster handles the "cluster" subcommand and routes to sub-handlers.
func runCluster(args []string) error {
	if isClusterRootHelp(args) {
		helptext.PrintWithMode("cluster", render.PrettyAuto)
		return nil
	}
	if err, isHelp := dispatchInvertedClusterHelp(args); isHelp {
		return err
	}
	err, isMatched := dispatchClusterSubcommand(args[0], args[1:])
	if isMatched {
		return err
	}
	return apperror.NewSimple("unknown command", "E9000")
}
