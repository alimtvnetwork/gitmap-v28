package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservercmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdnodes"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmigrate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvisibility"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdas"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdopen"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinject"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdexec"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhaschange"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreconcile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcreate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfix"
	"context"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcache"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcpar"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcursor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddoctor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpullerror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpurge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpushfix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsee"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
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
	entries := make([]dispatchEntry, 0, 18)
	entries = append(entries, coreBasicMaintenanceEntries()...)

	return append(entries, coreBasicOpEntries()...)
}

func coreBasicMaintenanceEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"clean-corrupted", "clean-corrupted-dirs"}, func() error { return cmddoctor.RunCleanCorrupted(argsTail()) }},
		{[]string{"purge", "purge-history"}, func() error { return cmdpurge.RunPurge(argsTail()) }},
		{[]string{"gitignore", "gitignore-agm", "gitignore-agy", "agm"}, func() error { return gitignoreagm.RunCLI(argsTail()) }},
		{[]string{"stash"}, func() error { return cmdfix.RunFix(argsTail(), "stash") }},
		{[]string{"wip"}, func() error { return cmdfix.RunFix(argsTail(), "wip") }},
		{[]string{"discard"}, func() error { return cmdfix.RunFix(argsTail(), "discard") }},
		{[]string{"fix-ignore-all", "fix-ignores-all", "fia"}, func() error { return cmdignore.RunFixIgnoreAll(argsTail()) }},
		{[]string{"fix-ignores-all-ssh", "fix-ignore-all-ssh", "fias"}, func() error { return cmdignore.RunFixIgnoresAllSSH(argsTail()) }},
		{[]string{"ignore", "ig"}, func() error { return cmdignore.RunIgnoreCLI(argsTail()) }},
		{[]string{"cache"}, func() error { return cmdcache.RunCacheCLI(argsTail()) }},
		{[]string{constants.CmdReconcile, constants.CmdReconcileAlias}, func() error { return cmdreconcile.RunReconcileCmd(argsTail()) }},
		{[]string{
			"vscode-optimize-projects", "vscode-optimize",
			"vsc-optimize-projects", "vsc-optimize",
			"vpm-optimize", "vpm-optimize-projects",
			"optimize-projects",
		}, func() error { return cmdvscode.RunVSCode(append([]string{"optimize-projects"}, argsTail()...)) }},
	}
}

func coreBasicOpEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdScan, constants.CmdScanAlias}, func() error { return runScan(argsTail()) }},
		{[]string{constants.CmdClone, constants.CmdCloneAlias}, func() error { return cmdclone.RunClone(argsTail()) }},
		{[]string{constants.CmdCloneOnlyMissing, constants.CmdCloneOnlyMissingAlias}, func() error { return cmdclone.RunCloneOnlyMissing(argsTail()) }},
		{[]string{
			constants.CmdCreate, constants.CmdCreateAlias,
			constants.CmdRepoCreate, constants.CmdRepoCreateAlias,
			constants.CmdCreateRepo, constants.CmdCreateRepoAlias, "cr",
		}, func() error { return cmdcreate.RunCreate(argsTail()) }},
		{[]string{
			constants.CmdCreateLocalRepo, constants.CmdCreateLocalRepoAlias,
			constants.CmdCreateRepoLocal, constants.CmdRepoCreateLocal,
		}, func() error { return cmdcreate.RunCreateLocal(argsTail()) }},
		{[]string{
			constants.CmdRecreateRepo, constants.CmdRecreateRepoAlias,
		}, func() error { return cmdrepo.RunRecreateRepo(argsTail()) }},
		{[]string{constants.CmdCloneSync, constants.CmdCloneSyncAlias}, cmdclone.RunCloneSync},
		{[]string{constants.CmdPull, constants.CmdPullAlias}, func() error { return runPull(argsTail()) }},
		{[]string{constants.CmdPush, constants.CmdPushAlias}, func() error { return runPush(argsTail()) }},
		{[]string{
			constants.CmdPushFix, constants.CmdPushFixAlias,
			constants.CmdPushFixSolid, constants.CmdPushFixInvert,
		}, func() error { return cmdpushfix.RunPushFix(argsTail()) }},
		{[]string{constants.CmdPullAll, constants.CmdPullAllAlias, "ta"}, func() error { return cmdpull.RunPullAll(argsTail()) }},
		{[]string{"pull-all-table", "pat"}, func() error { return cmdpull.RunPullAll(append([]string{"--status"}, argsTail()...)) }},
		{[]string{
			constants.CmdPullAllEfficient, constants.CmdPullAllEfficientAlias,
			constants.CmdPullAE,
		}, func() error {
			alias := subcommandName()
			isShort := alias == constants.CmdPullAllEfficientAlias || alias == constants.CmdPullAE
			return cmdpull.RunPullAllEfficient(argsTail(), false, alias, isShort)
		}},
		{[]string{
			constants.CmdPullAllEfficientTable, constants.CmdPullAllEfficientTableAlias,
		}, func() error {
			alias := subcommandName()
			isShort := alias == constants.CmdPullAllEfficientTableAlias
			return cmdpull.RunPullAllEfficient(argsTail(), true, alias, isShort)
		}},
		{[]string{"pull-all-ssh", "pas"}, func() error { return cmdpull.RunPullAll(append([]string{"--ssh"}, argsTail()...)) }},
		{[]string{"paswh", "pas-wh"}, func() error { return runPASWH(argsTail()) }},
		{[]string{"pull-error", "pull-errors", "pulle", "pull-e"}, func() error { return cmdpullerror.RunPullErrorCLI(argsTail()) }},
		{[]string{"commit-push-all-repos", "cpar"}, func() error { return cmdcpar.RunCPAR(argsTail()) }},
		{[]string{"see", "c"}, func() error { return cmdsee.RunSeeCLI(argsTail()) }},
		{[]string{"ses", "see-errors-ssh"}, func() error { return cmdsee.RunSeeErrorsSSH(argsTail()) }},
		{[]string{"repo-manage", "repo-manage-ui"}, func() error { return cmdsee.RunRepoManageUI() }},
		{[]string{constants.CmdStatus, constants.CmdStatusAlias}, func() error { return cmdstatus.RunStatus(argsTail()) }},
		{[]string{constants.CmdCommit, constants.CmdCommitAlias, constants.CmdCommitAlias2, constants.CmdCommitAlias3}, func() error { return cmdcommit.RunCommit(argsTail()) }},
		{[]string{"git"}, func() error { return runGitSubcommand(argsTail()) }},
		{[]string{constants.CmdExec, constants.CmdExecAlias}, func() error { return cmdexec.RunExec(argsTail()) }},
		{[]string{"py", "python"}, func() error { return cmdpy.RunPy(argsTail()) }},
		{[]string{constants.CmdSync, constants.CmdSyncAlias}, func() error { return cmdsync.RunSync(argsTail()) }},
	}
}

func coreWorkflowEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdHasAnyUpdates, constants.CmdHasAnyUpdatesAlias, constants.CmdHasAnyChanges, constants.CmdHasAnyChangesAlias}, func() error { return cmdhaschange.RunHasAnyUpdates(argsTail()) }},
		{[]string{constants.CmdHasChange, constants.CmdHasChangeAlias}, func() error { return cmdhaschange.RunHasChange(argsTail()) }},
		{[]string{constants.CmdCloneNext, constants.CmdCloneNextAlias}, func() error { return cmdclone.RunCloneNext(argsTail()) }},
		{[]string{constants.CmdAs, constants.CmdAsAlias}, func() error { return cmdas.RunAs(argsTail()) }},
		{[]string{constants.CmdCode, constants.CmdCodeAlias, constants.CmdCodeAlias2}, func() error { return cmdcode.RunCode(argsTail()) }},
		{[]string{constants.CmdInject, constants.CmdInjectAlias}, func() error { return cmdinject.RunInject(argsTail()) }},
		{[]string{constants.CmdOpen, constants.CmdOpenAlias}, func() error { return cmdopen.RunOpen(argsTail()) }},
		{[]string{constants.CmdCloneFrom, constants.CmdCloneFromAlias}, func() error { return cmdclone.RunCloneFrom(argsTail()) }},
		{[]string{"cursor", "cur"}, func() error { return cmdcursor.RunCursor(argsTail()) }},
	}
}

func coreCloneExtEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdMultiClone, constants.CmdMultiCloneAlias, constants.CmdMutliCloneAlias}, func() error { return cmdclone.RunMultiCloneCommand(argsTail()) }},
		{[]string{constants.CmdCloneReclone, constants.CmdCloneRecloneAlias, constants.CmdCloneNow, constants.CmdCloneNowAlias, constants.CmdCloneRel, constants.CmdCloneRelAlias}, func() error { return cmdclone.RunCloneNow(argsTail()) }},
		{[]string{constants.CmdClonePick, constants.CmdClonePickAlias}, func() error { return cmdclone.RunClonePick(argsTail()) }},
		{[]string{constants.CmdCommitIn, constants.CmdCommitInAlias, "commitin"}, func() error { return cmdcommitin.RunCommitIn(argsTail()) }},
		{[]string{"commit-pull", "cpull", "pull-commits"}, func() error { return cmdcommitpull.RunCommitPull(argsTail()) }},
		{[]string{"migrate", "wizard", "migration-wizard"}, func() error { return cmdmigrate.RunMigrateWizard(argsTail()) }},
		{[]string{constants.CmdCloneFixRepo, constants.CmdCloneFixRepoAlias}, func() error { return cmdclone.RunCloneFixRepo(argsTail()) }},
		{[]string{constants.CmdCloneFixRepoPub, constants.CmdCloneFixRepoPubAlias}, func() error { return cmdclone.RunCloneFixRepoPub(argsTail()) }},
		{[]string{constants.CmdVSCodePMSync, constants.CmdVSCodePMSyncAlias}, func() error { return cmdvscode.RunVSCodePMSync(argsTail()) }},
	}
}

func coreVisibilityActionEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdMakePublic}, func() error { return cmdvisibility.RunMakePublic(argsTail()) }},
		{[]string{constants.CmdMakePrivate}, func() error { return cmdvisibility.RunMakePrivate(argsTail()) }},
		{[]string{constants.CmdMakeAllPublic, constants.CmdMAPUB}, func() error { return cmdvisibility.RunMakeAllPublic(argsTail()) }},
		{[]string{constants.CmdMakeAllPrivate, constants.CmdMAPRI}, func() error { return cmdvisibility.RunMakeAllPrivate(argsTail()) }},
		{[]string{constants.CmdMakeAllPublicExceptLatest, constants.CmdMAPUBXL}, func() error { return cmdvisibility.RunMakeAllPublicExceptLatest(argsTail()) }},
		{[]string{constants.CmdMakeAllPrivateExceptLatest, constants.CmdMAPRIXL}, func() error { return cmdvisibility.RunMakeAllPrivateExceptLatest(argsTail()) }},
		{[]string{constants.CmdMakeLastPublic, constants.CmdMLPUB}, func() error { return cmdvisibility.RunMakeLastPublic(argsTail()) }},
		{[]string{constants.CmdMakeLastPrivate, constants.CmdMLPRI}, func() error { return cmdvisibility.RunMakeLastPrivate(argsTail()) }},
	}
}

func coreVisibilityHistoryEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdVisibilityUndo, constants.CmdVisibilityUndoAlias}, func() error { return cmdvisibility.RunVisibilityUndo(argsTail()) }},
		{[]string{constants.CmdVisibilityRedo, constants.CmdVisibilityRedoAlias}, func() error { return cmdvisibility.RunVisibilityRedo(argsTail()) }},
		{[]string{constants.CmdVisibilityHistory, constants.CmdVisibilityHistoryAlias}, func() error { return cmdvisibility.RunVisibilityHistory(argsTail()) }},
	}
}

func coreClusterEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"nodes", "node", "allnodes", "all-nodes", "fleet-nodes"}, func() error { return cmdnodes.RunUnifiedNodesCLI(argsTail()) }},
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
		{[]string{"nodes-deploy-agm-accounts", "sync-agm-accounts", "deploy-agm-accounts"}, func() error { return cmdnodes.RunUnifiedNodesCLI(append([]string{"deploy", "agm-accounts"}, argsTail()...)) }},
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

func runPASWH(args []string) error {
	if isPASWHHelp(args) {
		checkHelp("paswh", args)
		return nil
	}
	n, y, rest := parsePASWHArgs(args)
	pasArgs := append([]string{"--w", strconv.Itoa(n), "--hand", strconv.Itoa(y)}, rest...)
	return cmdssh.RunSSHPASFleet(pasArgs)
}

func parsePositiveIntArg(arg string, fallback int) int {
	parsed, err := strconv.Atoi(arg)
	if err == nil && parsed > 0 {
		return parsed
	}
	return fallback
}

func parsePASWHArgs(args []string) (int, int, []string) {
	n, y := 1, 1
	var rest []string
	posIdx := 0

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			rest = append(rest, arg)
			continue
		}
		switch posIdx {
		case 0:
			n = parsePositiveIntArg(arg, n)
			posIdx++
		case 1:
			y = parsePositiveIntArg(arg, y)
			posIdx++
		default:
			rest = append(rest, arg)
		}
	}

	return n, y, rest
}

func isPASWHHelp(args []string) bool {
	for i, a := range args {
		if a == "--help" || a == "help" {
			return true
		}
		if a == "-h" && !isNextArgNumeric(args, i) {
			return true
		}
	}
	return false
}

func isNextArgNumeric(args []string, i int) bool {
	if i+1 >= len(args) {
		return false
	}
	_, err := strconv.Atoi(args[i+1])
	return err == nil
}
