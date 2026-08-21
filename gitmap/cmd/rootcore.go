package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// dispatchCore routes scan, clone, pull, and status commands.
func dispatchCore(command string) bool {
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
	return []dispatchEntry{
		{[]string{constants.CmdScan, constants.CmdScanAlias}, func() { runScan(argsTail()) }},
		{[]string{constants.CmdClone, constants.CmdCloneAlias}, func() { runClone(argsTail()) }},
		{[]string{constants.CmdPull, constants.CmdPullAlias}, func() { runPull(argsTail()) }},
		{[]string{constants.CmdPush, constants.CmdPushAlias}, func() { runPush(argsTail()) }},
		{[]string{constants.CmdPullAll, constants.CmdPullAllAlias}, func() { runPullAll(argsTail()) }},
		{[]string{constants.CmdStatus, constants.CmdStatusAlias}, func() { runStatus(argsTail()) }},
		{[]string{constants.CmdExec, constants.CmdExecAlias}, func() { runExec(argsTail()) }},
	}
}

func coreWorkflowEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdHasAnyUpdates, constants.CmdHasAnyUpdatesAlias, constants.CmdHasAnyChanges, constants.CmdHasAnyChangesAlias}, func() { runHasAnyUpdates(argsTail()) }},
		{[]string{constants.CmdHasChange, constants.CmdHasChangeAlias}, func() { runHasChange(argsTail()) }},
		{[]string{constants.CmdCloneNext, constants.CmdCloneNextAlias}, func() { runCloneNext(argsTail()) }},
		{[]string{constants.CmdAs, constants.CmdAsAlias}, func() { runAs(argsTail()) }},
		{[]string{constants.CmdCode, constants.CmdCodeAlias, constants.CmdCodeAlias2}, func() { runCode(argsTail()) }},
		{[]string{constants.CmdInject, constants.CmdInjectAlias}, func() { runInject(argsTail()) }},
		{[]string{constants.CmdOpen, constants.CmdOpenAlias}, func() { runOpen(argsTail()) }},
		{[]string{constants.CmdCloneFrom, constants.CmdCloneFromAlias}, func() { runCloneFrom(argsTail()) }},
	}
}

func coreCloneExtEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdCloneReclone, constants.CmdCloneRecloneAlias, constants.CmdCloneNow, constants.CmdCloneNowAlias, constants.CmdCloneRel, constants.CmdCloneRelAlias}, func() { runCloneNow(argsTail()) }},
		{[]string{constants.CmdClonePick, constants.CmdClonePickAlias}, func() { runClonePick(argsTail()) }},
		{[]string{constants.CmdCommitIn, constants.CmdCommitInAlias}, func() { runCommitIn(argsTail()) }},
		{[]string{constants.CmdCloneFixRepo, constants.CmdCloneFixRepoAlias}, func() { runCloneFixRepo(argsTail()) }},
		{[]string{constants.CmdCloneFixRepoPub, constants.CmdCloneFixRepoPubAlias}, func() { runCloneFixRepoPub(argsTail()) }},
		{[]string{constants.CmdVSCodePMSync, constants.CmdVSCodePMSyncAlias}, func() { runVSCodePMSync(argsTail()) }},
	}
}

func coreVisibilityActionEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdMakePublic}, func() { runMakePublic(argsTail()) }},
		{[]string{constants.CmdMakePrivate}, func() { runMakePrivate(argsTail()) }},
		{[]string{constants.CmdMakeAllPublic, constants.CmdMAPUB}, func() { runMakeAllPublic(argsTail()) }},
		{[]string{constants.CmdMakeAllPrivate, constants.CmdMAPRI}, func() { runMakeAllPrivate(argsTail()) }},
		{[]string{constants.CmdMakeAllPublicExceptLatest, constants.CmdMAPUBXL}, func() { runMakeAllPublicExceptLatest(argsTail()) }},
		{[]string{constants.CmdMakeAllPrivateExceptLatest, constants.CmdMAPRIXL}, func() { runMakeAllPrivateExceptLatest(argsTail()) }},
		{[]string{constants.CmdMakeLastPublic, constants.CmdMLPUB}, func() { runMakeLastPublic(argsTail()) }},
		{[]string{constants.CmdMakeLastPrivate, constants.CmdMLPRI}, func() { runMakeLastPrivate(argsTail()) }},
	}
}

func coreVisibilityHistoryEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdVisibilityUndo, constants.CmdVisibilityUndoAlias}, func() { runVisibilityUndo(argsTail()) }},
		{[]string{constants.CmdVisibilityRedo, constants.CmdVisibilityRedoAlias}, func() { runVisibilityRedo(argsTail()) }},
		{[]string{constants.CmdVisibilityHistory, constants.CmdVisibilityHistoryAlias}, func() { runVisibilityHistory(argsTail()) }},
	}
}

func coreClusterEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdServersClients, constants.CmdSC}, func() { dispatchServersClients(argsTail()) }},
		{[]string{constants.CmdClients}, func() { dispatchClients(argsTail()) }},
		{[]string{"servers"}, func() { dispatchServers(argsTail()) }},
		{[]string{constants.CmdCluster, constants.CmdClusterAlias}, func() { runCluster(argsTail()) }},
	}
}

func dispatchServersClients(args []string) {
	if len(args) == 0 {
		runClusterCommand(cluster.ServersClients, args)
		return
	}
	subCmd, rest := args[0], args[1:]
	if dispatchServersClientsPathCmd(subCmd, rest) || dispatchClusterReadWrite(cluster.ServersClients, subCmd, rest) || dispatchClusterMutate(cluster.ServersClients, subCmd, rest) {
		return
	}
	runClusterCommand(cluster.ServersClients, args)
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

func dispatchClusterReadWrite(selector cluster.TargetSelectorType, subCmd string, rest []string) bool {
	switch subCmd {
	case "ls":
		runClusterLS(selector, rest)
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

func dispatchClients(args []string) {
	if len(args) == 0 {
		runClusterCommand(cluster.ClientsOnly, args)
		return
	}
	subCmd, rest := args[0], args[1:]
	if dispatchClusterReadWrite(cluster.ClientsOnly, subCmd, rest) || dispatchClusterMutate(cluster.ClientsOnly, subCmd, rest) {
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
