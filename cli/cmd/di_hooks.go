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
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
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
)

func init() {
	cmdpipeline.PipelineAgyFixRunner = cmdagy.RunPipelineFixAgyCLI
	cmdagy.SSHConnectionsFetcher = cmdssh.FetchAllSSHConnections
	cmdagy.SSHNodeDialer = cmdssh.DialSSHConnectionWithFallback
	cmdssh.JoinRunner = cmdssh.RunSSHJoinCLI
	cmdssh.ProfileRunner = runProfile
	cmdmacro.MacroSyncRunner = cmdssh.RunSSHMacroSyncCLI

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
	cmdchrome.RunFindDuplicatesFn = runFindDuplicates
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

	cmdupdate.RunPostUpdateMigrateFn = runPostUpdateMigrate
	cmdupdate.RequireOnlineFn = requireOnline
	cmdupdate.ResolveDeployedAndConfigPathsFn = resolveDeployedAndConfigPaths
	cmdupdate.PrintGitmapIdentityBlockLongFn = printGitmapIdentityBlockLong

	cmdpull.LoadAllRecordsDBFn = loadAllRecordsDB
	cmdpull.LoadRecordsByGroupFn = loadRecordsByGroup
	cmdpull.CreatePendingTaskFn = createPendingTask
	cmdpull.CompletePendingTaskFn = completePendingTask
	cmdpull.FailPendingTaskFn = failPendingTask
	cmdpull.RequireOnlineFn = requireOnline
	cmdpull.CheckHelpFn = checkHelp
	cmdpull.ApplyTransportFlagFn = ApplyTransportFlag
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
		return apperror.WrapSimple(runHistory(args), "history")
	}
	cmdpull.HasAliasFn = HasAlias
	cmdpull.GetAliasSlugFn = GetAliasSlug
	cmdpull.GetAliasPathFn = GetAliasPath
	cmdpull.RunStatusFn = runStatus
	cmdpull.RunSpecialRepoProbeOnPullFn = RunSpecialRepoProbeForPull
	cmdpull.PrintRemediationSummaryNoPromptFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		PrintRemediationSummaryNoPrompt(cmdItems)
	}
	cmdpull.PrintRemediationSummaryAutoFixFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		PrintRemediationSummaryAutoFix(cmdItems)
	}
	cmdpull.PrintRemediationSummaryFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		PrintRemediationSummary(cmdItems)
	}

	cmdclone.CreatePendingTaskFn = createPendingTask
	cmdclone.CompletePendingTaskFn = completePendingTask
	cmdclone.FailPendingTaskFn = failPendingTask
	cmdclone.RequireOnlineFn = requireOnline
	cmdclone.CheckHelpFn = checkHelp
	cmdclone.WriteShellHandoffFn = WriteShellHandoff
	cmdclone.EscapeCwdIfInsideFn = escapeCwdIfInside
	cmdclone.FinalizeErrorReportFn = cmdscan.FinalizeErrorReport
	cmdclone.RunCodingGuidelinesInstallFn = func(dir string) error {
		return RunCodingGuidelinesInstall(CodingGuidelinesOpts{WorkingDir: dir})
	}
	cmdclone.CommitCodingGuidelinesFn = func(dir string, isSkipCommit, isSkipPush bool) error {
		return CommitCodingGuidelines(CGCommitOpts{WorkingDir: dir, IsSkipCommit: isSkipCommit, IsSkipPush: isSkipPush})
	}
	cmdclone.RunCFRPPriorVersionPrivatizeFn = runCFRPPriorVersionPrivatize
	cmdclone.RunGitHubDesktopOptimizeFn = runGitHubDesktopOptimize
	cmdclone.ResolveEndpointStringFn = resolveEndpointString
	cmdclone.ResolveReleaseAliasPathFn = resolveReleaseAliasPath
	cmdclone.RunStatusFn = runStatus

	cmdscan.CreatePendingTaskFn = createPendingTask
	cmdscan.CompletePendingTaskFn = completePendingTask
	cmdscan.FailPendingTaskFn = failPendingTask
	cmdscan.SyncRecordsToVSCodePMFn = cmdvscode.SyncRecordsToVSCodePM
	cmdscan.SyncRecordsToIDEsFn = cmdide.SyncScanRecordsToIDEs
	cmdscan.RunPruneStaleDBFn = runPruneStaleDB
	cmdscan.CheckHelpFn = checkHelp
	cmdscan.CheckSpecialReposOnScanFn = func(workBaseDir string, isQuiet bool) {
		_, _ = CheckSpecialReposOnScan(workBaseDir, isQuiet)
	}

	cmdos.RunPowerNeverSleepFn = runPowerNeverSleep
	cmdos.RunPowerSetFn = runPowerSet
	cmdos.RunPowerResetFn = runPowerReset
	cmdos.CheckHelpFn = checkHelp

	cmdschedule.CheckHelpFn = checkHelp
	cmdmacro.RunScheduleFn = cmdschedule.RunSchedule

	cmdworkdir.CheckHelpFn = checkHelp
}
