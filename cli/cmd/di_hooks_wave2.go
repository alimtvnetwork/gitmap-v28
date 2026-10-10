package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdadd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdamend"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdas"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdaudit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdauditlegacy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdauthor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbackup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbookmark"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchangelog"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommitpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommittransfer"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcompletion"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcreate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddashboard"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddesktopsync"
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
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhygiene"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdimport"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinject"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinteractive"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdjoin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlatestbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlfscommon"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlist"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdllm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmigrate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmultigroup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdnodes"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdopen"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdorphans"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprobe"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprofiles"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprojectrepos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreconcile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrecreate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdregoldens"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreplace"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdresolver"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdselfinstall"
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
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtemprelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdui"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdversionhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvisibility"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwatch"
)

func wireWave2Hooks() {
	// Wave 243b: hook shims for packages split out of cmd/ (Wave 2 moves).
	cmdadd.CheckHelpFn = checkHelp
	cmdamend.CheckHelpFn = checkHelp
	cmdas.CheckHelpFn = checkHelp
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
	cmdaudit.CreatePendingTaskFn = cmdpending.CreatePendingTask
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
