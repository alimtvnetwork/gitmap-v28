package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdapps"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautomation"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbackup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcg"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcompletion"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddesktopsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddoctor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdenv"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgomod"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdimport"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlatestbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlist"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmergeai"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdopen"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdorphans"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrun"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdselfinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdseowrite"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservice"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsetup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsize"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstale"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstartup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtask"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtemprelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtoken"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvhost"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvmware"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdworkdir"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzip"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// dispatchTooling routes dev tooling and maintenance commands.
func dispatchTooling(command string) (bool, error) {
	return runDispatchTable(command, toolingDispatchEntries())
}

// toolingDispatchEntries returns the routing table for tooling commands.
func toolingDispatchEntries() []dispatchEntry {
	entries := make([]dispatchEntry, 0, 65)
	entries = append(entries, toolingWorkspaceEntries()...)
	entries = append(entries, toolingDevEntries()...)
	entries = append(entries, toolingAuditEntries()...)
	entries = append(entries, toolingOpsEntries()...)
	entries = append(entries, toolingInstallEntries()...)
	entries = append(entries, toolingUtilEntries()...)
	entries = append(entries, toolingChromeEntries()...)
	entries = append(entries, toolingNetworkEntries()...)

	return append(entries, toolingSystemEntries()...)
}

func toolingWorkspaceEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdDesktopSync, constants.CmdDesktopSyncAlias}, func() error { checkHelp("desktop-sync", argsTail()); return cmddesktopsync.RunDesktopSync(argsTail()) }},
		{[]string{constants.CmdGitHubDesktop, constants.CmdGitHubDesktopAlias, "github", "desktop", "gh-desktop"}, func() error { return cmdvscode.RunGitHubDesktop(argsTail()) }},
		{[]string{constants.CmdRescan, constants.CmdRescanAlias}, func() error { checkHelp("rescan", argsTail()); return cmdscan.RunRescan() }},
		{[]string{constants.CmdRescanSubtree, constants.CmdRescanSubtreeAlias}, func() error { return cmdscan.RunRescanSubtree(argsTail()) }},
		{[]string{constants.CmdSetup}, func() error { return cmdsetup.RunSetup(argsTail()) }},
		{[]string{constants.CmdDoctor}, func() error { checkHelp("doctor", argsTail()); return cmddoctor.RunDoctorCmd(argsTail()) }},
		{[]string{constants.CmdLatestBranch, constants.CmdLatestBranchAlias}, func() error { return cmdlatestbranch.RunLatestBranch(argsTail()) }},
		{[]string{constants.CmdBranch, constants.CmdBranchAlias}, func() error { return cmdbranch.RunBranch(argsTail()) }},
		{[]string{constants.CmdPendingCommits, constants.CmdPendingCommitsAlias}, func() error { return cmdpending.RunPendingCommits(argsTail()) }},
		{[]string{constants.CmdSends}, func() error { return RunSends(argsTail()) }},
	}
}

func toolingDevEntries() []dispatchEntry {
	entries := make([]dispatchEntry, 0, 20)
	entries = append(entries, toolingDevGeneralEntries()...)
	entries = append(entries, toolingDevSSHEntries()...)
	return append(entries, toolingDevDeployEntries()...)
}

func toolingDevGeneralEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdListVersions, constants.CmdListVersionsAlias}, func() error { return cmdlist.RunListVersions(argsTail()) }},
		{[]string{constants.CmdListReleases, constants.CmdListReleasesAlias, constants.CmdReleases}, func() error { return cmdlist.RunListReleases(argsTail()) }},
		{[]string{constants.CmdSEOWrite, constants.CmdSEOWriteAlias}, func() error { return cmdseowrite.RunSEOWrite(argsTail()) }},
		{[]string{constants.CmdGoMod, constants.CmdGoModAlias}, func() error { return cmdgomod.RunGoMod(argsTail()) }},
		{[]string{constants.CmdCompletion, constants.CmdCompletionAlias}, func() error { return cmdcompletion.RunCompletion(argsTail()) }},
		{[]string{constants.CmdZipGroup, constants.CmdZipGroupShort}, func() error { return cmdzip.RunZipGroup(argsTail()) }},
		{[]string{constants.CmdAlias, constants.CmdAliasShort}, func() error { return runAlias(argsTail()) }},
		{[]string{constants.CmdNewCommands, constants.CmdNewCommandsAlias}, func() error { return cmdpending.RunNewCommands(argsTail()) }},
		{[]string{constants.CmdBackup}, func() error { return cmdbackup.RunBackup(argsTail()) }},
	}
}

func toolingDevSSHEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"ssh-clone", "ssh-c"}, func() error { return cmdssh.RunSSHCloneCLI(argsTail()) }},
		{[]string{"token", "git-token", "access-token"}, func() error { return cmdtoken.Run(argsTail()) }},
		{[]string{"login", "signin"}, func() error { return cmdlogin.Run(argsTail()) }},
		{[]string{"logout", "signout"}, func() error { return cmdlogin.RunLogout(argsTail()) }},
		{[]string{constants.CmdSSH, "ssh-key", "ssh-keys", "auth-key", "auth-key-add", "ssh-key-add"}, func() error { return cmdssh.RunSSH(argsTail()) }},
		{[]string{"deploy-bin", "deploy-binary", "push-bin", "sync-bin"}, func() error { return cmdssh.RunSSHDeployBinCLI(argsTail()) }},
		{[]string{"pull-inventory", "fetch-inventory", "sync-inventory"}, func() error { return cmdssh.RunSSHPullInventoryCLI(argsTail()) }},
		{[]string{"scp", "ssh-cp", "ssh-copy"}, func() error { return cmdssh.RunSSHCPCLI(argsTail()) }},
	}
}

func toolingDevDeployEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"deploy"}, func() error { return cmdssh.RunSSHDeployCLI("deploy", argsTail()) }},
		{[]string{"deploy-ide", "ide-deploy"}, func() error { return cmdagy.RunAgyDeployCLI(append([]string{"--all"}, argsTail()...)) }},
		{[]string{"agy-deploy", "deploy-agy"}, func() error { return cmdagy.RunAgyDeployCLI(argsTail()) }},
		{[]string{"deploy-right"}, func() error { return cmdssh.RunSSHDeployCLI("deploy-right", argsTail()) }},
		{[]string{"deploy-left"}, func() error { return cmdssh.RunSSHDeployCLI("deploy-left", argsTail()) }},
		{[]string{"deploy-config", "deploy-config-ssh", "deploy-ssh-config"}, func() error { return cmdssh.RunSSHDeployConfigSSHCLI(argsTail()) }},
		{[]string{"deploy-keys", "deploy-keys-all", "deploy-all-keys", "deploy-all-key", "deploy-key", "deploykeys", "deploy-key-all"}, func() error { return cmdssh.RunSSHDeployKeysCLI(argsTail()) }},
		{[]string{"export-ssh", "ssh-export", "export-ssh-nodes", "nodes-export-json"}, func() error { return cmdssh.RunSSHNodesExportJSON(argsTail()) }},
		{[]string{"import-ssh", "ssh-import", "import-ssh-nodes", "nodes-import-json"}, func() error { return cmdssh.RunSSHNodesImportJSON(argsTail()) }},
	}
}

func toolingAuditEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdStale, constants.CmdStaleAlias}, func() error { return cmdstale.RunStale(argsTail()) }},
		{[]string{constants.CmdOrphans}, func() error { return cmdorphans.RunOrphans(argsTail()) }},
		{[]string{constants.CmdDedupe}, func() error { return cmdscan.RunDedupe(argsTail()) }},
		{[]string{constants.CmdSize}, func() error { return cmdsize.RunSize(argsTail()) }},
		{[]string{constants.CmdReleaseNotes}, func() error { return cmdrelease.RunReleaseNotes(argsTail()) }},
		{[]string{constants.CmdReleaseDry}, func() error { return cmdrelease.RunReleaseDry(argsTail()) }},
		{[]string{constants.CmdTagRename}, func() error { return cmdrelease.RunTagRename(argsTail()) }},
		{[]string{constants.CmdRecent, constants.CmdRecentAlias}, func() error { return runRecent(argsTail()) }},
		{[]string{constants.CmdTodo}, func() error { return runTodo(argsTail()) }},
		{[]string{"run-errors", "run-err"}, func() error { return cmdrun.RunErrorsCmd(argsTail()) }},
		{[]string{"run-history"}, func() error { return cmdrun.RunHistoryCmd(argsTail()) }},
	}
}

func toolingOpsEntries() []dispatchEntry {
	entries := []dispatchEntry{
		{[]string{constants.CmdOpen, constants.CmdOpenAlias}, func() error { return cmdopen.RunOpen(argsTail()) }},
		{[]string{constants.CmdPR, constants.CmdPRAlias}, func() error { return runPR(argsTail()) }},
		{[]string{constants.CmdBlameStats}, func() error { return runBlameStats(argsTail()) }},
		{[]string{constants.CmdSnapshot}, func() error { return RunSnapshot(argsTail()) }},
		{[]string{constants.CmdRollback}, func() error { return RunRollback(argsTail()) }},
		{[]string{constants.CmdGuard}, func() error { return cmdautomation.RunGuard(argsTail()) }},
		{[]string{constants.CmdPrune, constants.CmdPruneAlias}, func() error { return cmdprune.RunPrune(argsTail()) }},
		{[]string{constants.CmdTempRelease, constants.CmdTempReleaseShort}, func() error { return cmdtemprelease.RunTempRelease(argsTail()) }},
		{[]string{constants.CmdTask, constants.CmdTasks, constants.CmdTaskAlias}, func() error { return cmdtask.RunTasks(argsTail()) }},
		{[]string{constants.CmdEnv, constants.CmdEnvAlias}, func() error { return cmdenv.RunEnv(argsTail()) }},
		{[]string{constants.CmdService, constants.CmdServiceAlias, "services"}, func() error { return cmdservice.Run(argsTail()) }},
		{[]string{"run", "run-macro", "exec-macro"}, func() error { return cmdmacro.RunMacroRootRun(argsTail()) }},
		{[]string{"run-until"}, func() error { return cmdmacro.RunMacroRootRunUntil(argsTail()) }},
		{[]string{
			"clean-dev", "cleandev", "dev-cleanup", "devcleanup", "dev-clean",
			"devtools-cache", "dev-tools-cache", "dev-tool-cache", "devtool-cache",
			"clear-devtools", "clear-dev-tools", "clear-dev-tools-cache", "clear-devtools-cache",
			"clear-devtool", "clear-dev-tool", "clean-devtools", "clean-dev-tools", "clean-dev-tools-cache",
			"devtools-cache-clear", "dev-tools-cache-clear",
		}, func() error { return RunCleanDevTopLevel(argsTail()) }},
		{[]string{
			"terminal", "clean-terminal", "clear-terminal", "terminal-clean", "terminal-clear",
			"clean-term", "clear-term",
		}, func() error { return RunTerminalTopLevel(argsTail()) }},
		{[]string{"devtool", "devtools", "dev-tool", "dev-tools", "dt"}, func() error { return RunDevToolTopLevel(argsTail()) }},
		{[]string{"clean"}, func() error { return RunCleanTopLevel(argsTail()) }},
		{[]string{"clear", "cls"}, func() error { return RunClearTopLevel(argsTail()) }},
		{[]string{"winutil"}, func() error { return RunWinUtilTopLevel(argsTail()) }},
	}
	return append(entries, toolingSpecialRepoEntries()...)
}

func toolingSpecialRepoEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdRepoSecretsAlias, constants.CmdRepoSecrets}, func() error { return rsCmd(argsTail()) }},
		{[]string{constants.CmdRepoCacheAlias, constants.CmdRepoCache, constants.CmdRepoStorageAlias}, func() error { return rcCmd(argsTail()) }},
		{[]string{constants.CmdSummary}, func() error { return cmdsummary.RunSummary(argsTail()) }},
		{[]string{constants.CmdFullSummary, constants.CmdFullStatus, constants.CmdFs}, func() error { return cmdsummary.RunFullSummary(argsTail()) }},
		{[]string{constants.CmdFullSummaryPE, constants.CmdFullStatusPE, constants.CmdFsPE, constants.CmdFspe}, func() error { return cmdsummary.RunFullSummary(append([]string{"+pe"}, argsTail()...)) }},
		{[]string{constants.CmdMergeAI, constants.CmdMergeAIAlias}, func() error { return cmdmergeai.RunMergeAI(argsTail()) }},
	}
}

func toolingInstallEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"muse", "meta-muse", "metamuse"}, func() error { return RunMuseCLI(argsTail()) }},
		{[]string{"installer"}, func() error { return cmdinstaller.RunInstallerCLI(argsTail()) }},
		{[]string{"which-format", "whichformat", "format-which", "format-inspect", "format-check", "which-json"}, func() error { return cmdconfig.RunWhichFormatCLI(argsTail()) }},
		{[]string{"which"}, func() error { return handleWhichSubcommand(argsTail()) }},
		{[]string{"format"}, func() error { return handleFormatSubcommand(argsTail()) }},
		{[]string{"json"}, func() error { return handleJSONSubcommand(argsTail()) }},
		{[]string{"import-all-json", "import-all", "importall", "import-json-all", "importjsonall"}, func() error { return cmdimport.RunImportAllJSONCLI(argsTail()) }},
		{[]string{"what-configs", "what-config", "whatconfigs", "whatconfig", "wc"}, func() error { return cmdconfig.RunWhatConfigsCLI(argsTail()) }},
		{[]string{"merge-json", "mergejson", "json-merge"}, func() error { return cmdmerge.RunMergeJSONCLI(argsTail()) }},
		{[]string{"pin", "version-pin"}, func() error { return RunPinCLI(argsTail()) }},
		{[]string{"unpin", "version-unpin"}, func() error { return RunUnpinCLI(argsTail()) }},
		{[]string{"cg", "coding-guide", "coding-guidelines", "ct"}, func() error { return cmdcg.RunCG(argsTail()) }},
		{[]string{"install-version-json", "init-version"}, func() error { return cmdcg.RunCG(append([]string{"install-version-json"}, argsTail()...)) }},
		{[]string{"install-prompts", "install-prompt"}, func() error { return cmdcg.RunCG(append([]string{"install-prompts"}, argsTail()...)) }},
		{[]string{"prompts-status"}, func() error { return cmdcg.RunCG(append([]string{"prompts-status"}, argsTail()...)) }},
		{[]string{"prompts-version"}, func() error { return cmdcg.RunCG(append([]string{"prompts-version"}, argsTail()...)) }},
		{[]string{"workdir", "work-dir", "wd"}, func() error { return cmdworkdir.RunWorkDir(argsTail()) }},
		{[]string{"os", "os-update"}, func() error { return cmdos.RunOSCLI(argsTail()) }},
		{[]string{"change-password", "passwd", "chpasswd"}, func() error { return cmdos.RunChangePasswordCLI(argsTail()) }},
		{[]string{"machine", "machines", "machine-name", "hostname"}, func() error { return cmdos.RunMachineCLI(argsTail()) }},
		{[]string{"os-info", "osinfo", "sysinfo", "system-info", "which-os", "whichos", "os-which"}, func() error { return cmdos.RunOSInfoCLI(argsTail()) }},
		{[]string{"bash", "git-bash"}, func() error { return RunBash(argsTail()) }},
		{[]string{"shell", "sh"}, func() error { return RunShell(argsTail()) }},
		{[]string{"powershell", "pwsh", "ps"}, func() error { return RunPowerShell(argsTail()) }},
		{[]string{"sj", "ssh-joiner", "ssh-join", "ssh-joined"}, func() error { return runSJ(argsTail()) }},
		{[]string{"sjc", "ssh-join-common", "ssh-join-c", "join-common"}, func() error { return runSJC(argsTail()) }},
		{[]string{"se", "ssh-exe", "ssh-exec", "ssh-execute", "mm", "multi-machines", "multi-machine", "multimachine", "multimachines"}, func() error { return cmdssh.RunSSHExec(argsTail()) }},
		{[]string{"remote"}, func() error { return cmdssh.RunRemote(argsTail()) }},
		{[]string{constants.CmdCargo, constants.CmdCargoAlias}, func() error { return cmdcargo.RunCargo(argsTail()) }},
		{[]string{constants.CmdInstall, constants.CmdInstallAlias}, func() error { return cmdinstall.RunInstall(argsTail()) }},
		{[]string{constants.CmdUninstall, constants.CmdUninstallAlias}, func() error { return runUninstall(argsTail()) }},
		{[]string{"apps", "app"}, func() error { return cmdapps.RunAppsDispatch(argsTail()) }},
		{[]string{"uninstall-agy"}, func() error { return runUninstall(append([]string{"agy"}, argsTail()...)) }},
		{[]string{"uninstall-agy-all"}, func() error { return runUninstall(append([]string{"agy-all"}, argsTail()...)) }},
		{[]string{"uninstall-agm"}, func() error { return runUninstall(append([]string{"agm"}, argsTail()...)) }},
		{[]string{"uninstall-agm-all"}, func() error { return runUninstall(append([]string{"agm-all"}, argsTail()...)) }},
		{[]string{"uninstall-copilot"}, func() error { return runUninstall(append([]string{"copilot"}, argsTail()...)) }},
		{[]string{"uninstall-copilot-all"}, func() error { return runUninstall(append([]string{"copilot-all"}, argsTail()...)) }},
		{[]string{"uninstall-edge"}, func() error { return runUninstall(append([]string{"edge"}, argsTail()...)) }},
		{[]string{"uninstall-edge-all"}, func() error { return runUninstall(append([]string{"edge-all"}, argsTail()...)) }},
		{[]string{constants.CmdExportConfig, constants.CmdExportConfigAlias}, func() error { return cmdconfig.RunExportConfig(argsTail()) }},
		{[]string{constants.CmdImportConfig, constants.CmdImportConfigAlias, constants.CmdImportConfigTypo}, func() error { return cmdconfig.RunImportConfig(argsTail()) }},
		{[]string{constants.CmdStartupAdd, constants.CmdStartupAddAlias}, func() error { return cmdstartup.RunStartupAdd(argsTail()) }},
		{[]string{constants.CmdStartupList, constants.CmdStartupListAlias}, func() error { return cmdstartup.RunStartupList(argsTail()) }},
		{[]string{constants.CmdStartupRemove, constants.CmdStartupRemoveAlias}, func() error { return cmdstartup.RunStartupRemove(argsTail()) }},
		{[]string{constants.CmdSelfInstall}, func() error { return cmdselfinstall.RunSelfInstall(argsTail()) }},
		{[]string{constants.CmdSelfUninstall}, func() error { return cmdselfinstall.RunSelfUninstall(argsTail()) }},
		{[]string{constants.CmdSelfUninstallRunner}, cmdselfinstall.RunSelfUninstallRunner},
		{[]string{constants.CmdPending}, cmdpending.RunPending},
		{[]string{constants.CmdDoPending, constants.CmdDoPendingAlias}, func() error { return cmdpending.RunDoPending(argsTail()) }},
		{[]string{constants.CmdVmware, constants.CmdVmwareAlias}, func() error { return cmdvmware.Run(argsTail()) }},
		{[]string{"perms", "permissions"}, func() error { return cmdsetup.RunSetupPerms(argsTail()) }},
		{[]string{constants.CmdVHost}, func() error { return cmdvhost.RunVHost(argsTail()) }},
		{[]string{constants.CmdNginx, constants.CmdNginxAlias}, func() error { return cmdvhost.RunNginx(argsTail()) }},
	}
}
