package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdas"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcache"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcpar"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcreate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcursor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddoctor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdexec"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhaschange"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinject"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmigrate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdopen"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpullerror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpurge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpushfix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreconcile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus"
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
		{[]string{"gitignore", "gitignore-agm", "gitignore-agy"}, func() error { return gitignoreagm.RunCLI(argsTail()) }},
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
		{[]string{"pull-all-ssh-fix", "pas-fix", "pasf", "pa-fleet-fix"}, func() error { return cmdssh.RunSSHPullFleetRemediate(argsTail()) }},
		{[]string{"paswh", "pas-wh"}, func() error { return runPASWH(argsTail()) }},
		{[]string{"pull-error", "pull-errors", "pulle", "pull-e"}, func() error { return cmdpullerror.RunPullErrorCLI(argsTail()) }},
		{[]string{"commit-push-all-repos", "cpar"}, func() error { return cmdcpar.RunCPAR(argsTail()) }},
		{[]string{"see", "c"}, func() error { return RunSeeCLI(argsTail()) }},
		{[]string{"ses", "see-errors-ssh"}, func() error { return RunSeeErrorsSSH(argsTail()) }},
		{[]string{"repo-manage", "repo-manage-ui"}, func() error { return RunRepoManageUI() }},
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
		{[]string{"commit-info", "ci"}, func() error { return cmdcommit.RunCommitInfo(argsTail()) }},
	}
}

func coreCloneExtEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdMultiClone, constants.CmdMultiCloneAlias, constants.CmdMutliCloneAlias}, func() error { return cmdclone.RunMultiCloneCommand(argsTail()) }},
		{[]string{constants.CmdCloneReclone, constants.CmdCloneRecloneAlias, constants.CmdCloneNow, constants.CmdCloneNowAlias, constants.CmdCloneRel, constants.CmdCloneRelAlias}, func() error { return cmdclone.RunCloneNow(argsTail()) }},
		{[]string{constants.CmdClonePick, constants.CmdClonePickAlias}, func() error { return cmdclone.RunClonePick(argsTail()) }},
		{[]string{constants.CmdCommitIn, constants.CmdCommitInAlias, "commitin", "merge-in", "mergein"}, func() error { return cmdcommitin.RunCommitIn(argsTail()) }},
		{[]string{"commit-pull", "cpull", "pull-commits"}, func() error { return cmdcommitpull.RunCommitPull(argsTail()) }},
		{[]string{"migrate", "wizard", "migration-wizard"}, func() error { return cmdmigrate.RunMigrateWizard(argsTail()) }},
		{[]string{constants.CmdCloneFixRepo, constants.CmdCloneFixRepoShort}, func() error { return cmdclone.RunCloneFixRepo(argsTail()) }},
		{[]string{constants.CmdCloneFixRepoPub, constants.CmdCloneFixRepoPubAlias, constants.CmdCloneFixRepoPubShort}, func() error { return cmdclone.RunCloneFixRepoPub(argsTail()) }},
		{[]string{constants.CmdVSCodePMSync, constants.CmdVSCodePMSyncAlias}, func() error { return cmdvscode.RunVSCodePMSync(argsTail()) }},
	}
}
