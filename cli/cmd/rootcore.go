package cmd

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/cluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpurge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// dispatchCore routes scan, clone, pull, and status commands.
func dispatchCore(command string) (bool, error) {
	return runDispatchTable(command, coreDispatchEntries())
}

// coreDispatchEntries returns the routing table for core commands.
func coreDispatchEntries() []dispatchEntry {
	entries := make([]dispatchEntry, 0, 36)
	entries = append(entries, coreBasicEntries()...)
	entries = append(entries, coreWorkflowEntries()...)
	entries = append(entries, coreCloneExtEntries()...)
	entries = append(entries, coreVisibilityActionEntries()...)
	entries = append(entries, coreVisibilityHistoryEntries()...)
	entries = append(entries, coreClusterEntries()...)

	return entries
}

func coreBasicEntries() []dispatchEntry {
	entries := make([]dispatchEntry, 0, 17)
	entries = append(entries, coreBasicMaintenanceEntries()...)

	return append(entries, coreBasicOpEntries()...)
}

func coreBasicMaintenanceEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"clean-corrupted", "clean-corrupted-dirs"}, func() error { return runCleanCorrupted(argsTail()) }},
		{[]string{"purge", "purge-history"}, func() error { return cmdpurge.RunPurge(argsTail()) }},
		{[]string{"fix"}, func() error { return runFix(argsTail(), "") }},
		{[]string{"stash"}, func() error { return runFix(argsTail(), "stash") }},
		{[]string{"wip"}, func() error { return runFix(argsTail(), "wip") }},
		{[]string{"discard"}, func() error { return runFix(argsTail(), "discard") }},
		{[]string{constants.CmdReconcile, constants.CmdReconcileAlias}, func() error { return RunReconcileCmd(argsTail()) }},
	}
}

func coreBasicOpEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdScan, constants.CmdScanAlias}, func() error { return runScan(argsTail()) }},
		{[]string{constants.CmdClone, constants.CmdCloneAlias}, func() error { return runClone(argsTail()) }},
		{[]string{constants.CmdCloneOnlyMissing, constants.CmdCloneOnlyMissingAlias}, func() error { return runCloneOnlyMissing(argsTail()) }},
		{[]string{
			constants.CmdCreate, constants.CmdCreateAlias,
			constants.CmdRepoCreate, constants.CmdRepoCreateAlias,
			constants.CmdCreateRepo, constants.CmdCreateRepoAlias, "cr",
		}, func() error { return runCreate(argsTail()) }},
		{[]string{
			constants.CmdCreateLocalRepo, constants.CmdCreateLocalRepoAlias,
			constants.CmdCreateRepoLocal, constants.CmdRepoCreateLocal,
		}, func() error { return runCreateLocal(argsTail()) }},
		{[]string{
			constants.CmdRecreateRepo, constants.CmdRecreateRepoAlias,
		}, func() error { return runRecreateRepo(argsTail()) }},
		{[]string{constants.CmdCloneSync, constants.CmdCloneSyncAlias}, runCloneSync},
		{[]string{constants.CmdPull, constants.CmdPullAlias}, func() error { return runPull(argsTail()) }},
		{[]string{constants.CmdPush, constants.CmdPushAlias}, func() error { return runPush(argsTail()) }},
		{[]string{constants.CmdPullAll, constants.CmdPullAllAlias}, func() error { return runPullAll(argsTail()) }},
		{[]string{
			constants.CmdPullAllEfficient, constants.CmdPullAllEfficientAlias,
			constants.CmdPullAE,
		}, func() error {
			alias := subcommandName()
			isShort := alias == constants.CmdPullAllEfficientAlias || alias == constants.CmdPullAE
			return runPullAllEfficient(argsTail(), false, alias, isShort)
		}},
		{[]string{
			constants.CmdPullAllEfficientTable, constants.CmdPullAllEfficientTableAlias,
		}, func() error {
			alias := subcommandName()
			isShort := alias == constants.CmdPullAllEfficientTableAlias
			return runPullAllEfficient(argsTail(), true, alias, isShort)
		}},
		{[]string{constants.CmdStatus, constants.CmdStatusAlias}, func() error { return runStatus(argsTail()) }},
		{[]string{"git"}, func() error { return runGitSubcommand(argsTail()) }},
		{[]string{constants.CmdExec, constants.CmdExecAlias}, func() error { return runExec(argsTail()) }},
	}
}

func coreWorkflowEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdHasAnyUpdates, constants.CmdHasAnyUpdatesAlias, constants.CmdHasAnyChanges, constants.CmdHasAnyChangesAlias}, func() error { return runHasAnyUpdates(argsTail()) }},
		{[]string{constants.CmdHasChange, constants.CmdHasChangeAlias}, func() error { return runHasChange(argsTail()) }},
		{[]string{constants.CmdCloneNext, constants.CmdCloneNextAlias}, func() error { return runCloneNext(argsTail()) }},
		{[]string{constants.CmdAs, constants.CmdAsAlias}, func() error { return runAs(argsTail()) }},
		{[]string{constants.CmdCode, constants.CmdCodeAlias, constants.CmdCodeAlias2}, func() error { return runCode(argsTail()) }},
		{[]string{constants.CmdInject, constants.CmdInjectAlias}, func() error { return runInject(argsTail()) }},
		{[]string{constants.CmdOpen, constants.CmdOpenAlias}, func() error { return runOpen(argsTail()) }},
		{[]string{constants.CmdCloneFrom, constants.CmdCloneFromAlias}, func() error { return runCloneFrom(argsTail()) }},
	}
}

func coreCloneExtEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdMultiClone, constants.CmdMultiCloneAlias, constants.CmdMutliCloneAlias}, func() error { return runMultiClone(argsTail()) }},
		{[]string{constants.CmdCloneReclone, constants.CmdCloneRecloneAlias, constants.CmdCloneNow, constants.CmdCloneNowAlias, constants.CmdCloneRel, constants.CmdCloneRelAlias}, func() error { return runCloneNow(argsTail()) }},
		{[]string{constants.CmdClonePick, constants.CmdClonePickAlias}, func() error { return runClonePick(argsTail()) }},
		{[]string{constants.CmdCommitIn, constants.CmdCommitInAlias}, func() error { return runCommitIn(argsTail()) }},
		{[]string{"commit-pull", "cpull", "pull-commits"}, func() error { return runCommitPull(argsTail()) }},
		{[]string{"migrate", "wizard", "migration-wizard"}, func() error { return runMigrateWizard(argsTail()) }},
		{[]string{constants.CmdCloneFixRepo, constants.CmdCloneFixRepoAlias}, func() error { return runCloneFixRepo(argsTail()) }},
		{[]string{constants.CmdCloneFixRepoPub, constants.CmdCloneFixRepoPubAlias}, func() error { return runCloneFixRepoPub(argsTail()) }},
		{[]string{constants.CmdVSCodePMSync, constants.CmdVSCodePMSyncAlias}, func() error { return runVSCodePMSync(argsTail()) }},
	}
}

func coreVisibilityActionEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdMakePublic}, func() error { return runMakePublic(argsTail()) }},
		{[]string{constants.CmdMakePrivate}, func() error { return runMakePrivate(argsTail()) }},
		{[]string{constants.CmdMakeAllPublic, constants.CmdMAPUB}, func() error { return runMakeAllPublic(argsTail()) }},
		{[]string{constants.CmdMakeAllPrivate, constants.CmdMAPRI}, func() error { return runMakeAllPrivate(argsTail()) }},
		{[]string{constants.CmdMakeAllPublicExceptLatest, constants.CmdMAPUBXL}, func() error { return runMakeAllPublicExceptLatest(argsTail()) }},
		{[]string{constants.CmdMakeAllPrivateExceptLatest, constants.CmdMAPRIXL}, func() error { return runMakeAllPrivateExceptLatest(argsTail()) }},
		{[]string{constants.CmdMakeLastPublic, constants.CmdMLPUB}, func() error { return runMakeLastPublic(argsTail()) }},
		{[]string{constants.CmdMakeLastPrivate, constants.CmdMLPRI}, func() error { return runMakeLastPrivate(argsTail()) }},
	}
}

func coreVisibilityHistoryEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdVisibilityUndo, constants.CmdVisibilityUndoAlias}, func() error { return runVisibilityUndo(argsTail()) }},
		{[]string{constants.CmdVisibilityRedo, constants.CmdVisibilityRedoAlias}, func() error { return runVisibilityRedo(argsTail()) }},
		{[]string{constants.CmdVisibilityHistory, constants.CmdVisibilityHistoryAlias}, func() error { return runVisibilityHistory(argsTail()) }},
	}
}

func coreClusterEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdServersClients, constants.CmdServersClientsAlias, constants.CmdSC}, func() error { dispatchServersClients(argsTail()); return nil }},
		{[]string{constants.CmdClients, constants.CmdClientsAlias}, func() error { dispatchClients(argsTail()); return nil }},
		{[]string{"servers"}, func() error { dispatchServers(argsTail()); return nil }},
		{[]string{constants.CmdCluster, constants.CmdClusterAlias}, func() error { return runCluster(argsTail()) }},
		{[]string{constants.CmdServerCmd, constants.CmdServerCmds, constants.CmdServerCmdAlias}, func() error { return runServerCmd(argsTail()) }},
	}
}

func dispatchServersClients(args []string) {
	if IsHelpRequestedOrEmpty(args) {
		RenderSCHelp()

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

	runClusterCommand(cluster.ServersClients, args)
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
	return dispatchSCAuthOps(subCmd, rest)
}

func dispatchSCInventoryOps(subCmd string, rest []string) bool {
	switch subCmd {
	case "nodes", "list", "machines", "joined":
		_ = runClusterNodes(rest)
		return true
	case "rm", "remove", "delete":
		_ = runClusterRemove(rest)
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
		runClusterSetDefaultPath(cluster.ServersClients, rest)

		return true
	case "set-path-alias":
		runClusterSetPathAlias(cluster.ServersClients, rest)

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

	_ = runClusterNodes(rest)
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
		runClusterCat(selector, rest)
	case "write":
		runClusterWrite(selector, rest)
	default:
		return false
	}

	return true
}

func dispatchClusterMutate(selector cluster.TargetSelectorType, subCmd string, rest []string) bool {
	switch subCmd {
	case "update":
		runClusterUpdate(selector, false, rest)
	case "update-all":
		runClusterUpdate(selector, true, rest)
	case "clone", "cfr", "cfrp":
		runClusterClone(selector, subCmd, rest)
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

	runClusterCommand(cluster.ClientsOnly, args)
}

func dispatchServers(args []string) {
	if len(args) == 0 {
		return
	}

	if args[0] == "ls" {
		runClusterLS(cluster.ServersOnly, args[1:])

		return
	}

	dispatchServersUpdate(args[0], args[1:])
}

func dispatchServersUpdate(subCmd string, rest []string) {
	if subCmd == "update" {
		runClusterUpdate(cluster.ServersOnly, false, rest)
	} else if subCmd == "update-all" {
		runClusterUpdate(cluster.ServersOnly, true, rest)
	}
}
