package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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
	return []dispatchEntry{
		{[]string{constants.CmdScan, constants.CmdScanAlias}, func() error { return runScan(argsTail()) }},
		{[]string{constants.CmdClone, constants.CmdCloneAlias}, func() error { return runClone(argsTail()) }},
		{[]string{constants.CmdCloneSync, constants.CmdCloneSyncAlias}, func() error { return runCloneSync() }},
		{[]string{constants.CmdPull, constants.CmdPullAlias}, func() error { return runPull(argsTail()) }},
		{[]string{constants.CmdPush, constants.CmdPushAlias}, func() error { return runPush(argsTail()) }},
		{[]string{constants.CmdPullAll, constants.CmdPullAllAlias}, func() error { return runPullAll(argsTail()) }},
		{[]string{constants.CmdStatus, constants.CmdStatusAlias}, func() error { return runStatus(argsTail()) }},
		{[]string{"fix"}, func() error { return runFix(argsTail(), "") }},
		{[]string{"stash"}, func() error { return runFix(argsTail(), "stash") }},
		{[]string{"wip"}, func() error { return runFix(argsTail(), "wip") }},
		{[]string{"discard"}, func() error { return runFix(argsTail(), "discard") }},

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
		{[]string{constants.CmdCloneReclone, constants.CmdCloneRecloneAlias, constants.CmdCloneNow, constants.CmdCloneNowAlias, constants.CmdCloneRel, constants.CmdCloneRelAlias}, func() error { return runCloneNow(argsTail()) }},
		{[]string{constants.CmdClonePick, constants.CmdClonePickAlias}, func() error { return runClonePick(argsTail()) }},
		{[]string{constants.CmdCommitIn, constants.CmdCommitInAlias}, func() error { return runCommitIn(argsTail()) }},
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
		{[]string{constants.CmdServersClients, constants.CmdSC}, func() error { dispatchServersClients(argsTail()); return nil }},
		{[]string{constants.CmdClients}, func() error { dispatchClients(argsTail()); return nil }},
		{[]string{"servers"}, func() error { dispatchServers(argsTail()); return nil }},
		{[]string{constants.CmdCluster, constants.CmdClusterAlias}, func() error { return runCluster(argsTail()) }},
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
