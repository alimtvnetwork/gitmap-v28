package cmdnodes

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpullerror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

func isNodesPingCommand(arg string) bool {
	low := strings.ToLower(arg)
	return low == "ping" || low == "probe"
}

func isNodesCloneRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	_, isClone := IsNodesCloneCommand(args[0])
	return isClone
}

func runNodesHistoryCLI(args []string) error {
	return cmdssh.RunFleetPASCommand("nodes history", "gitmap history --limit 10", func() error {
		return cmdhistory.RunHistory(args)
	})
}

func runNodesErrorsCLI(args []string) error {
	cmdStr := "gitmap errors " + strings.Join(args, " ")
	cmdStr = strings.TrimSpace(cmdStr)

	return cmdssh.RunFleetPASCommand("nodes errors", cmdStr, func() error {
		return cmderrors.RunErrorsCLI(args)
	})
}

func runNodesPullErrorsCLI(args []string) error {
	cmdStr := "gitmap pull errors " + strings.Join(args, " ")
	cmdStr = strings.TrimSpace(cmdStr)

	return cmdssh.RunFleetPASCommand("nodes pull errors", cmdStr, func() error {
		return cmdpullerror.RunPullErrorCLI(args)
	})
}

func isNodesAgyRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	tok := strings.ToLower(args[0])
	return tok == "agy" || tok == "agy-ui" || tok == "nodes-agy-ui" || tok == "studio" || tok == "dashboard"
}

func runNodesAgyDispatch(args []string) error {
	if len(args) == 0 {
		return RunNodesAgyUI(nil)
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "query":
		return RunNodesAgyQuery(args[1:])
	case "prompt":
		return RunNodesAgyPrompt(args[1:])
	case "push-settings", "pushsettings":
		return RunNodesPushSettings(args[1:])
	case "sync-settings", "syncsettings":
		return RunNodesSyncSettings(args[1:])
	case "send-projects", "sendprojects", "sync-projects", "syncprojects":
		return RunNodesSendProjects(args[1:])
	case "ui", "dashboard", "studio":
		return RunNodesAgyUI(args[1:])
	default:
		return RunNodesAgyUI(args)
	}
}

func isNodesDeployAGMRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "sync-agm-accounts" || first == "sync-agm" || first == "deploy-agm-accounts" || first == "deploy-agm" {
		return true
	}
	if first != "deploy" || len(args) < 2 {
		return false
	}
	sub := strings.ToLower(args[1])
	return sub == "agm-accounts" || sub == "agm" || sub == "accounts" || sub == "agm-account"
}

func runNodesDeployAGMDispatch(args []string) error {
	if len(args) == 0 {
		return RunNodesDeployAGMAccounts(nil)
	}
	first := strings.ToLower(args[0])
	if first == "deploy" && len(args) > 1 {
		return RunNodesDeployAGMAccounts(args[2:])
	}
	return RunNodesDeployAGMAccounts(args[1:])
}

func isNodesDeployRepoRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "deploy-repo" || first == "deploy-repos" || first == "deployrepo" || first == "deployrepos" {
		return true
	}
	if first != "deploy" || len(args) < 2 {
		return false
	}
	sub := strings.ToLower(args[1])
	return sub == "repo" || sub == "repos" || sub == "repository" || sub == "repositories"
}

func runNodesDeployRepoDispatch(args []string) error {
	if len(args) == 0 {
		return RunNodesDeployRepo(nil)
	}
	first := strings.ToLower(args[0])
	if first == "deploy-repos" || first == "deployrepos" {
		return RunNodesDeployRepos(args[1:])
	}
	if first == "deploy-repo" || first == "deployrepo" {
		return RunNodesDeployRepo(args[1:])
	}
	if first == "deploy" && len(args) > 1 {
		return dispatchDeploySubcommand(args)
	}
	return RunNodesDeployRepo(args[1:])
}

func dispatchDeploySubcommand(args []string) error {
	sub := strings.ToLower(args[1])
	if sub == "repos" || sub == "repositories" {
		return RunNodesDeployRepos(args[2:])
	}
	return RunNodesDeployRepo(args[2:])
}

func isNodesScanRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	return first == "scan" || first == "rescan"
}

func runNodesScanDispatch(args []string) error {
	first := strings.ToLower(args[0])
	if first == "rescan" {
		return RunNodesRescan(args[1:])
	}
	return RunNodesScan(args[1:])
}

// runUnifiedNodesCLI handles `gitmap nodes [flags] [target]`.
func RunUnifiedNodesCLI(args []string) error {
	for len(args) > 0 && (strings.EqualFold(args[0], "nodes") || strings.EqualFold(args[0], "node")) {
		args = args[1:]
	}
	if len(args) == 0 {
		return runUnifiedNodesList(args)
	}

	switch {
	case isNodesUpdateSoftwareRequest(args):
		return runNodesUpdateSoftwareDispatch(args[1:])
	case isNodesRunningOnlyRequest(args):
		return RunNodesRunningOnlyCLI(args[1:])
	case hasNodesRunningOnlyFlag(args):
		return RunNodesRunningOnlyCLI(args)
	case isNodesAgyRequest(args):
		return runNodesAgyDispatch(args[1:])
	case strings.EqualFold(args[0], "push-settings") || strings.EqualFold(args[0], "pushsettings"):
		return RunNodesPushSettings(args[1:])
	case strings.EqualFold(args[0], "sync-settings") || strings.EqualFold(args[0], "syncsettings"):
		return RunNodesSyncSettings(args[1:])
	case isNodesDeployAGMRequest(args):
		return runNodesDeployAGMDispatch(args)
	case isNodesDeployRepoRequest(args):
		return runNodesDeployRepoDispatch(args)
	case isNodesScanRequest(args):
		return runNodesScanDispatch(args)
	case strings.EqualFold(args[0], "send-projects") || strings.EqualFold(args[0], "sendprojects") ||
		strings.EqualFold(args[0], "sync-projects") || strings.EqualFold(args[0], "syncprojects"):
		return RunNodesSendProjects(args[1:])
	case strings.EqualFold(args[0], "history") || strings.EqualFold(args[0], "histories"):
		return runNodesHistoryCLI(args[1:])
	case strings.EqualFold(args[0], "errors"):
		return runNodesErrorsCLI(args[1:])
	case strings.EqualFold(args[0], "pull") && len(args) > 1 && strings.EqualFold(args[1], "errors"):
		return runNodesPullErrorsCLI(args[2:])
	case IsNodesPullAllRequest(args):
		return RunNodesPullAll(args[1:])
	case isNodesPingCommand(args[0]):
		return RunUnifiedNodesPingCLI(args[1:])
	case strings.EqualFold(args[0], "pending-commits") || strings.EqualFold(args[0], "pc") || strings.EqualFold(args[0], "pendingcommits"):
		return RunNodesPendingCommits(args[1:])
	case strings.EqualFold(args[0], "commits"):
		return RunNodesCommits(args[1:])
	case strings.EqualFold(args[0], "cpf"):
		return RunNodesCommitPushFeature(args[1:])
	case strings.EqualFold(args[0], "cpb"):
		return RunNodesCommitPushBug(args[1:])
	case strings.EqualFold(args[0], "cpr"):
		return RunNodesCommitPushRelease(args[1:])
	case strings.EqualFold(args[0], "commit-fix") || strings.EqualFold(args[0], "commitfix"):
		return RunNodesCommitFix(args[1:])
	case isNodesFullSummary(args):
		return RunNodesFullSummary(args)
	case strings.EqualFold(args[0], "summary"):
		return RunNodesSummary(args[1:])
	case isNodesPipelineAll(args):
		return RunNodesPipelineErrorsAll(args)
	case isNodesCloneRequest(args):
		return RunNodesClone(args)
	case isNodesResourcesRequest(args):
		return RunNodesResourcesCLI(args[1:])
	case isNodesTasksRequest(args):
		return RunNodesTasksCLI(args[1:])
	case isNodesRemoteRunRequest(args):
		return RunNodesRemoteCLI(args[1:])
	case isNodesHelpRequest(args):
		return printUnifiedNodesHelp()
	}

	return runUnifiedNodesList(args)
}

func isNodesFullSummary(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "fs" || first == "fs+pe" || first == "fspe" ||
		first == "full-status" || first == "full-summary" ||
		first == "full-status+pe" || first == "full-summary+pe" {
		return true
	}
	if first == "full" && len(args) > 1 {
		second := strings.ToLower(args[1])
		return second == "status" || second == "summary" || second == "status+pe" || second == "summary+pe"
	}
	return false
}

func isNodesPipelineAll(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if (first == "pe" || first == "pipe" || first == "pipeline") && len(args) > 1 && strings.EqualFold(args[1], "all") {
		return true
	}
	if (first == "pipe" || first == "pipeline") && len(args) > 2 &&
		(strings.EqualFold(args[1], "error") || strings.EqualFold(args[1], "errors")) &&
		strings.EqualFold(args[2], "all") {
		return true
	}
	return first == "pe-all" || first == "pe_all" ||
		first == "pipe-error-all" || first == "pipe-errors-all" ||
		first == "pipeline-error-all" || first == "pipeline-errors-all"
}

func isNodesResourcesRequest(args []string) bool {
	hasArgs := len(args) > 0
	if hasArgs {
		tok := strings.ToLower(args[0])
		return tok == "resources" || tok == "res" || tok == "top" || tok == "metrics" || tok == "usage"
	}
	return false
}

func isNodesTasksRequest(args []string) bool {
	hasArgs := len(args) > 0
	if hasArgs {
		tok := strings.ToLower(args[0])
		return tok == "tasks" || tok == "task" || tok == "task-queue" || tok == "tasks-queue" || tok == "queue" || tok == "jobs" || tok == "delegation"
	}
	return false
}

func isNodesRemoteRunRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	return first == "run" || first == "exec" || first == "r"
}

func isNodesUpdateSoftwareRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])

	return first == "update" || first == "up" || first == "upgrade" || first == "update-software"
}

func runNodesUpdateSoftwareDispatch(args []string) error {
	if len(args) == 0 {
		printNodesUpdateSoftwareHelp()
		return nil
	}
	app := args[0]

	return RunFleetSoftwareUpdate(app, args[1:])
}

func printNodesUpdateSoftwareHelp() {
	fmt.Println("\nUsage: gitmap nodes update <app> [node|all] [flags]")
	fmt.Println("\nSupported software applications:")
	fmt.Println("  cursor              - Cursor IDE")
	fmt.Println("  vscode / code       - Visual Studio Code")
	fmt.Println("  agm                 - Antigravity Manager")
	fmt.Println("  agy / antigravity   - Antigravity Tools / IDE")
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap nodes update cursor all")
	fmt.Println("  gitmap nodes update vscode u1")
	fmt.Println("  gitmap nodes update agm all")
	fmt.Println("  gitmap nodes update agy all")
	fmt.Println()
}
