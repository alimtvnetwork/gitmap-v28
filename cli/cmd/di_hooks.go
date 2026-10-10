// Package cmd — di_hooks.go: cross-package dependency-injection wiring.
//
// Moved from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D).
// This init() wires cmd/ orchestrator functions into cmdX hook fields;
// it must stay in package cmd (thin dispatch glue).
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdaudit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdexec"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfind"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprofiles"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreconcile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdworkdir"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	HistoryRunnerFn = func(args []string) *apperror.AppError {
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
	cmdclone.RunCFRPPriorVersionPrivatizeFn = RunCFRPPriorVersionPrivatize
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

	wireWave2Hooks()
}
