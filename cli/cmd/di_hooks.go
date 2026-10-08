// Package cmd — di_hooks.go: cross-package dependency-injection wiring.
//
// Moved from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D).
// This init() wires cmd/ orchestrator functions into cmdX hook fields;
// it must stay in package cmd (thin dispatch glue).
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcfrppriorversion"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddesktopsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdimport"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinject"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsee"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdworkdir"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdadd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdamend"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdas"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbackup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdaudit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdauditlegacy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdauthor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbookmark"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchangelog"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommittransfer"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcompletion"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcreate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddashboard"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiff"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiffprofiles"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddir"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddocs"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddownload"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdenv"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdexec"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdexport"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfind"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfindnext"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfoldertree"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgitrm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgomod"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgroup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhaschange"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhygiene"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinteractive"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlatestbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlfscommon"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlist"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdnodes"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdopen"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdllm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmigrate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmultigroup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdorphans"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprobe"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprojectrepos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrecreate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreconcile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdregoldens"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreplace"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdresolver"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdselfinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdseowrite"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsequence"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservercmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsize"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstale"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstartup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstats"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstorage"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtask"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtemplates"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdversionhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvisibility"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtemprelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwatch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdjoin"
	"github.com/alimtvnetwork/gitmap-v28/cli/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdui"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprofiles"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
)

func init() {
	cmdpipeline.PipelineAgyFixRunner = cmdagy.RunPipelineFixAgyCLI
	cmdagy.SSHConnectionsFetcher = cmdssh.FetchAllSSHConnections
	cmdagy.SSHNodeDialer = cmdssh.DialSSHConnectionWithFallback
	cmdagy.RunAgySSHFn = cmdssh.RunSSHAgyCLI
	cmdagy.RunSSHUpdateFn = cmdssh.RunSSHUpdateCLI
	cmdssh.JoinRunner = cmdssh.RunSSHJoinCLI
	cmdssh.ProfileRunner = cmdprofiles.RunProfile
	cmdmacro.MacroSyncRunner = cmdssh.RunSSHMacroSyncCLI
	cmdmacro.FinishCommandAuditFn = cmdaudit.FinishCommandAudit

	cmdinstaller.ResolveProfileTreeFn = func(s string) (any, bool) {
		p, ok := cmdinstall.ResolveProfileTree(s)
		return p, ok
	}
	cmdinstaller.PrintProfileTreeFn = func(p any) {
		if prof, ok := p.(cmdinstall.ProfileComposition); ok {
			cmdinstall.PrintProfileTree(prof)
		}
	}
	cmdinstaller.PrintProfileInstallSummaryFn = cmdinstall.PrintProfileInstallSummary
	cmdinstaller.RunInstallAddFn = cmdinstall.RunInstallAdd

	cmdchrome.InstallChromeLinuxFn = func(isDryRun bool) error {
		return cmdinstall.RunInstallChromeLinux(cmdinstall.InstallOptions{Tool: constants.ToolChrome, DryRun: isDryRun})
	}
	cmdchrome.InstallToolFn = func(tool string, isDryRun bool) {
		cmdinstall.InstallTool(cmdinstall.InstallOptions{Tool: tool, DryRun: isDryRun})
	}
	cmdchrome.RunFindDuplicatesFn = cmdfind.RunFindDuplicates
	cmdchrome.CheckHelpFn = checkHelp
	cmdcargo.CheckHelpFn = checkHelp
	cmdinstall.CheckHelpFn = checkHelp
	cmdinstall.RemoteAgmUpdateFn = func(target string) error {
		return cmdssh.RunSSHUpdateCLI([]string{"agm", target})
	}
	cmdinstall.RemoteAgmUpdateFleetFn = func(target, except string) error {
		args := []string{"agm", target}
		if except != "" {
			args = append(args, "--except", except)
		}
		return cmdssh.RunSSHUpdateCLI(args)
	}
	cmdagy.RunAgySSHFn = cmdssh.RunSSHAgyCLI

	cmdupdate.RunPostUpdateMigrateFn = cmdupdate.RunPostUpdateMigrate
	cmdupdate.RequireOnlineFn = requireOnline
	cmdupdate.ResolveDeployedAndConfigPathsFn = cmdinstall.ResolveDeployedAndConfigPaths
	cmdupdate.PrintGitmapIdentityBlockLongFn = printGitmapIdentityBlockLong

	cmdpull.LoadAllRecordsDBFn = cmdstatus.LoadAllRecordsDB
	cmdexec.LoadAllRecordsDBFn = cmdstatus.LoadAllRecordsDB
	cmdpull.CreatePendingTaskFn = cmdaudit.CreatePendingTask
	cmdpull.CompletePendingTaskFn = cmdpending.CompletePendingTask
	cmdpull.FailPendingTaskFn = cmdpending.FailPendingTask
	cmdpull.RequireOnlineFn = requireOnline
	cmdpull.CheckHelpFn = checkHelp
	cmdpull.ApplyTransportFlagFn = cmdclone.ApplyTransportFlag
	cmdpull.RunRemoteSSHPullFn = cmdssh.RunSSHPullJSON
	cmdpull.RunRemoteSSHPullAllFleetFn = cmdssh.RunSSHPullAllFleet
	cmdssh.RunLocalPullAllJSONFn = cmdpull.RunPullAllJSON
	cmdignore.RunFixIgnoresAllSSHFn = func(args []string) *apperror.AppError {
		err := cmdssh.RunFleetPASCommand("fix-ignore-all", "gitmap fix-ignore-all -y", func() error {
			return cmdignore.RunFixIgnoreAll(args)
		})
		return apperror.WrapSimple(err, "fix-ignore-all")
	}
	cmdsee.HistoryRunnerFn = func(args []string) *apperror.AppError {
		return apperror.WrapSimple(cmdhistory.RunHistory(args), "history")
	}
	cmdpull.HasAliasFn = HasAlias
	cmdpull.GetAliasSlugFn = GetAliasSlug
	cmdpull.GetAliasPathFn = GetAliasPath
	cmdpull.RunStatusFn = cmdclone.RunStatus
	cmdpull.RunSpecialRepoProbeOnPullFn = RunSpecialRepoProbeForPull
	cmdpull.PrintRemediationSummaryNoPromptFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]cmdremediation.RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = cmdremediation.RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		cmdremediation.PrintRemediationSummaryNoPrompt(cmdItems)
	}
	cmdpull.PrintRemediationSummaryAutoFixFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]cmdremediation.RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = cmdremediation.RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		cmdremediation.PrintRemediationSummaryAutoFix(cmdItems)
	}
	cmdpull.PrintRemediationSummaryFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]cmdremediation.RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = cmdremediation.RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		cmdremediation.PrintRemediationSummary(cmdItems)
	}

	cmdclone.CreatePendingTaskFn = cmdaudit.CreatePendingTask
	cmdclone.CompletePendingTaskFn = cmdpending.CompletePendingTask
	cmdclone.FailPendingTaskFn = cmdpending.FailPendingTask
	cmdclone.RequireOnlineFn = requireOnline
	cmdclone.CheckHelpFn = checkHelp
	cmdclone.WriteShellHandoffFn = cmdclone.WriteShellHandoff
	cmdclone.EscapeCwdIfInsideFn = escapeCwdIfInside
	cmdclone.FinalizeErrorReportFn = cmdscan.FinalizeErrorReport
	cmdclone.RunCodingGuidelinesInstallFn = func(dir string) error {
		return cmdclone.RunCodingGuidelinesInstall(cmdclone.CodingGuidelinesOpts{WorkingDir: dir})
	}
	cmdclone.CommitCodingGuidelinesFn = func(dir string, isSkipCommit, isSkipPush bool) error {
		return cmdclone.CommitCodingGuidelines(cmdclone.CGCommitOpts{WorkingDir: dir, IsSkipCommit: isSkipCommit, IsSkipPush: isSkipPush})
	}
	cmdclone.RunCFRPPriorVersionPrivatizeFn = cmdcfrppriorversion.RunCFRPPriorVersionPrivatize
	cmdclone.RunGitHubDesktopOptimizeFn = cmdclone.RunGitHubDesktopOptimize
	cmdclone.ResolveEndpointStringFn = cmdclone.ResolveEndpointString
	cmdclone.ResolveReleaseAliasPathFn = cmdclone.ResolveReleaseAliasPath
	cmdclone.RunStatusFn = cmdclone.RunStatus

	cmdscan.CreatePendingTaskFn = cmdaudit.CreatePendingTask
	cmdscan.CompletePendingTaskFn = cmdpending.CompletePendingTask
	cmdscan.FailPendingTaskFn = cmdpending.FailPendingTask
	cmdscan.SyncRecordsToVSCodePMFn = cmdvscode.SyncRecordsToVSCodePM
	cmdscan.SyncRecordsToIDEsFn = cmdide.SyncScanRecordsToIDEs
	cmdscan.RunPruneStaleDBFn = cmdreconcile.RunPruneStaleDB
	cmdscan.CheckHelpFn = checkHelp
	cmdscan.CheckSpecialReposOnScanFn = func(workBaseDir string, isQuiet bool) {
		_, _ = CheckSpecialReposOnScan(workBaseDir, isQuiet)
	}

	cmdos.RunPowerNeverSleepFn = cmdos.RunPowerNeverSleep
	cmdos.RunPowerSetFn = cmdos.RunPowerSet
	cmdos.RunPowerResetFn = cmdos.RunPowerReset
	cmdos.CheckHelpFn = checkHelp

	cmdschedule.CheckHelpFn = checkHelp
	cmdmacro.RunScheduleFn = cmdschedule.RunSchedule

	cmdworkdir.CheckHelpFn = checkHelp

	// Wave 243b: hook shims for packages split out of cmd/ (Wave 2 moves).
	cmdadd.CheckHelpFn = checkHelp
	cmdamend.CheckHelpFn = checkHelp
	cmdas.CheckHelpFn = checkHelp
	cmdcluster.CheckHelpFn = checkHelp
	cmdcluster.HasHelpFlagFn = hasHelpFlag
	cmdcode.CheckHelpFn = checkHelp
	cmddesktopsync.HasGHDesktopInstallFlagFn = cmdvscode.HasGHDesktopInstallFlag
	cmdcluster.CheckHelpOrEmptyFn = CheckHelpOrEmpty
	cmdbackup.CheckHelpFn = checkHelp
	cmdcd.CheckHelpFn = checkHelp
	cmdcd.WriteShellHandoffFn = cmdclone.WriteShellHandoff
	cmdcd.RunCDSpecialRepoFn = runCDSpecialRepo
	cmdcd.HasAliasFn = HasAlias
	cmdcd.GetAliasPathFn = GetAliasPath
	cmdcd.DispatchFn = dispatch
	cmdaudit.CreatePendingTaskFn = cmdaudit.CreatePendingTask
	cmdauditlegacy.CheckHelpFn = checkHelp
	cmdauthor.CheckHelpFn = checkHelp
	cmdbookmark.CheckHelpFn = checkHelp
	cmdbookmark.DispatchFn = dispatch
	cmdbranch.CheckHelpFn = checkHelp
	cmdchangelog.CheckHelpFn = checkHelp
	cmdcommit.CheckHelpFn = checkHelp
	cmdcommitin.RunCommitPullBootstrapFn = cmdcommitpull.RunCommitPullBootstrap
	cmdcommittransfer.RunCommitPullFn = cmdcommitpull.RunCommitPull
	cmdcommittransfer.RunMigrateWizardFn = cmdmigrate.RunMigrateWizard
	cmdcommitin.PrintCommitPullTreeFn = cmdcommitpull.PrintCommitPullTree
	cmdcommitpull.RunCommitInFn = cmdcommitin.RunCommitIn
	cmdcommittransfer.CheckHelpFn = checkHelp
	cmdcommittransfer.ResolveEndpointStringFn = cmdclone.ResolveEndpointString
	cmdcompletion.CheckHelpFn = checkHelp
	cmdcreate.CheckHelpFn = checkHelp
	cmddashboard.CheckHelpFn = checkHelp
	cmddb.CheckHelpFn = checkHelp
	cmddiff.CheckHelpFn = checkHelp
	cmddiffprofiles.CheckHelpFn = checkHelp
	cmddir.CheckHelpFn = checkHelp
	cmddocs.CheckHelpFn = checkHelp
	cmddownload.CheckHelpFn = checkHelp
	cmdenv.CheckHelpFn = checkHelp
	cmderrors.CheckHelpFn = checkHelp
	cmdexec.CheckHelpFn = checkHelp
	cmdexec.CompletePendingTaskFn = cmdpending.CompletePendingTask
	cmdexec.CreatePendingTaskFn = cmdaudit.CreatePendingTask
	cmdexec.FailPendingTaskFn = cmdpending.FailPendingTask
	cmdexec.GetAliasSlugFn = GetAliasSlug
	cmdexec.HasAliasFn = HasAlias
	cmdexport.CheckHelpFn = checkHelp
	cmdfind.CheckHelpFn = checkHelp
	cmdfindnext.CheckHelpFn = checkHelp
	cmdfix.CheckHelpFn = checkHelp
	cmdfix.HasAliasFn = HasAlias
	cmdfix.GetAliasSlugFn = GetAliasSlug
	cmdfix.GetAliasPathFn = GetAliasPath
	cmdfix.LoadAllRecordsDBFn = cmdstatus.LoadAllRecordsDB
	cmdfix.LoadRecordsJSONFallbackFn = cmdstatus.LoadRecordsJSONFallback
	cmdfoldertree.CheckHelpFn = checkHelp
	cmdhistory.CheckHelpFn = checkHelp
	cmdimport.CheckHelpFn = checkHelp
	cmdui.CheckHelpFn = checkHelp
	cmdgitrm.CheckHelpFn = checkHelp
	cmdgomod.CheckHelpFn = checkHelp
	cmdgroup.CheckHelpFn = checkHelp
	cmdhaschange.CheckHelpFn = checkHelp
	cmdhaschange.ResolveReleaseAliasPathFn = cmdclone.ResolveReleaseAliasPath
	cmdinteractive.CheckHelpFn = checkHelp
	cmdlatestbranch.CheckHelpFn = checkHelp
	cmdlfscommon.CheckHelpFn = checkHelp
	cmdlist.CheckHelpFn = checkHelp
	cmdllm.CheckHelpFn = checkHelp
	cmdmerge.CheckHelpFn = checkHelp
	cmdmultigroup.CheckHelpFn = checkHelp
	cmdnodes.CheckHelpFn = checkHelp
	cmdopen.CheckHelpFn = checkHelp
	cmdopen.ResolveEndpointStringFn = cmdclone.ResolveEndpointString
	cmdopen.ConfigureDetachedProcessFn = cmdchrome.ConfigureDetachedProcess
	cmdorphans.CheckHelpFn = checkHelp
	cmdpending.CheckHelpFn = checkHelp
	cmdprobe.CheckHelpFn = checkHelp
	cmdprofiles.CheckHelpFn = checkHelp
	cmdprofiles.IsHelpFlagFn = isHelpFlag
	cmdprojectrepos.CheckHelpFn = checkHelp
	cmdprune.CheckHelpFn = checkHelp
	cmdrecreate.CheckHelpFn = checkHelp
	cmdreconcile.CheckHelpFn = checkHelp
	cmdregoldens.CheckHelpFn = checkHelp
	cmdreplace.CheckHelpFn = checkHelp
	cmdrepo.CheckHelpFn = checkHelp
	cmdrevert.CheckHelpFn = checkHelp
	cmdignore.CheckHelpFn = checkHelp
	cmdinject.CheckHelpFn = checkHelp
	cmdinject.WriteShellHandoffFn = cmdclone.WriteShellHandoff
	cmdinject.ExtractPositionalArgsFn = cmdrelease.ExtractPositionalArgs
	cmdrelease.CheckHelpFn = checkHelp
	cmdrelease.RequireOnlineFn = requireOnline
	cmdrelease.SplitSemverFn = cmdchangelog.SplitSemver
	cmdrepo.RunCreateFn = cmdcreate.RunCreate
	cmdrepo.RunRecreateFn = cmdrecreate.RunRecreate
	cmdselfinstall.CheckHelpFn = checkHelp
	cmdselfinstall.CleanCorruptedInstallDirsSilentFn = cmdmigrate.CleanCorruptedInstallDirsSilent
	cmdselfinstall.SetHiddenProcessAttrFn = cmdinstall.SetHiddenProcessAttr
	cmdrepo.RunCreateLocalFn = cmdcreate.RunCreateLocal
	cmdrm.CheckHelpFn = checkHelp
	cmdsearch.CheckHelpFn = checkHelp
	cmdseowrite.CheckHelpFn = checkHelp
	cmdsequence.CheckHelpFn = checkHelp
	cmdservercmd.CheckHelpFn = checkHelp
	cmdsize.CheckHelpFn = checkHelp
	cmdstale.CheckHelpFn = checkHelp
	cmdstartup.CheckHelpFn = checkHelp
	cmdstats.CheckHelpFn = checkHelp
	cmdstatus.CheckHelpFn = checkHelp
	cmdstatus.GetAliasPathFn = GetAliasPath
	cmdstatus.GetAliasSlugFn = GetAliasSlug
	cmdstorage.RunBackupCloudRestoreFn = cmdbackup.RunBackupCloudRestore
	cmdstatus.HasAliasFn = HasAlias
	cmdsync.CheckHelpFn = checkHelp
	cmdtask.CheckHelpFn = checkHelp
	cmdtemplates.CheckHelpFn = checkHelp
	cmdtemprelease.CheckHelpFn = checkHelp
	cmdundo.CheckHelpFn = checkHelp
	cmdversionhistory.CheckHelpFn = checkHelp
	cmdvisibility.CheckHelpFn = checkHelp
	cmdwatch.CheckHelpFn = checkHelp

	cmdbranch.PrintHelpAndExitFn = printHelpAndExit
	cmdtemplates.PrintHelpAndExitFn = printHelpAndExit
	cmdtemplates.IsHelpFlagFn = IsHelpFlag
	cmdenv.PrintGitmapIdentityBlockShortFn = printGitmapIdentityBlockShort
	cmddb.PopulateRepoAliasesWithCountFn = PopulateRepoAliasesWithCount
	cmdrm.ResolveMultiReposFn = cmdresolver.ResolveMultiRepos
	cmdhygiene.IsGitRepoFn = clonenext.IsGitRepo
	cmdexec.GetAliasPathFn = GetAliasPath
	cmdjoin.CheckHelpOrEmptyFn = CheckHelpOrEmpty
	cmdservercmd.CheckHelpOrEmptyFn = CheckHelpOrEmpty
}
