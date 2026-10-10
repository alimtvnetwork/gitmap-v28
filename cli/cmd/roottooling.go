package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdauthor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdjoin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservercmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdregoldens"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdauditlegacy"
	
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdindex"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreplace"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfind"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdllm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdadd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgitrm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsequence"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfoldertree"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddownload"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvhost"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstartup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdapps"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdselfinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdimport"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdenv"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtask"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtemprelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautomation"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdopen"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchangelog"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsize"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdorphans"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstale"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdseowrite"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlist"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgomod"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcompletion"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbackup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlatestbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddesktopsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmergeai"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdasset"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcg"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchromeprofile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddoctor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdports"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrun"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservice"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsetup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtoken"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvmware"
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
		{[]string{constants.CmdReleaseNotes}, func() error { return cmdchangelog.RunReleaseNotes(argsTail()) }},
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

func toolingUtilEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"schedule", "sc", "crontab", "cron"}, func() error { return cmdschedule.RunSchedule(argsTail()) }},
		{[]string{constants.CmdDownloaderConfig, constants.CmdDownloaderConfigAlias}, func() error { return cmddownload.RunDownloaderConfig(argsTail()) }},
		{[]string{constants.CmdUnzipCompact, constants.CmdUnzipCompactAlias}, func() error { return cmdzip.RunUnzipCompact(argsTail()) }},
		{[]string{constants.CmdFolder, "tree"}, func() error { return cmdfoldertree.RunFolder(argsTail()) }},
		{[]string{
			constants.CmdLowercase,
			constants.CmdLowerCaseFix,
			constants.CmdLowerCaseFixAlias,
			constants.CmdLowerCaseFixShort,
			constants.CmdLowercaseReadme,
			constants.CmdReadmeLower,
			"lower",
			"lower-case-readme",
			"readme-lowercase",
			"lcr",
			"lc-fix",
		}, func() error {
			checkHelp("lowercase", argsTail())
			return runLowerCaseFixCLI(argsTail())
		}},
		{[]string{constants.CmdFixSeqFiles, constants.CmdFixSeqFilesAlias}, func() error { return runFixSeqFiles(argsTail()) }},
		{[]string{constants.CmdSequence, constants.CmdSequenceAlias}, func() error { return cmdsequence.RunSequence(argsTail()) }},
		{[]string{constants.CmdGitRm}, func() error { return cmdgitrm.RunGitRm(argsTail()) }},
		{[]string{constants.CmdCommitPush, constants.CmdCommitPushAlias}, func() error { return cmdcommit.RunCommitPush(argsTail()) }},
		{[]string{constants.CmdPullCommitPush, constants.CmdPullCommitPushAlias}, func() error { return cmdcommit.RunPullCommitPush(argsTail()) }},
		{[]string{constants.CmdCommitPushBug, constants.CmdCommitPushBugAlias}, func() error { return cmdcommit.RunCommitPushBug(argsTail()) }},
		{[]string{constants.CmdCommitPushFeature, constants.CmdCommitPushFeatureAlias}, func() error { return cmdcommit.RunCommitPushFeature(argsTail()) }},
		{[]string{constants.CmdCommitPushChore}, func() error { return cmdcommit.RunCommitPushChore(argsTail()) }},
		{[]string{"cpc"}, func() error { return runCpcDisambiguated(argsTail()) }},
		{[]string{constants.CmdCommitPushRelease, constants.CmdCommitPushReleaseAlias}, func() error { return cmdcommit.RunCommitPushRelease(argsTail()) }},
		{[]string{constants.CmdRmGit, constants.CmdRmGitAlias}, func() error { return cmdcommit.RunRmGit(argsTail()) }},
		{[]string{constants.CmdGitReset, constants.CmdGitResetAlias}, func() error { return cmdcommit.RunGitReset(argsTail()) }},
		{[]string{constants.CmdIgnore}, func() error { return cmdignore.RunIgnoreCLI(argsTail()) }},
		{[]string{constants.CmdIgnoreRm}, func() error { return cmdignore.RunIgnoreRm(argsTail()) }},
		{[]string{constants.CmdAdd}, func() error { return cmdadd.RunAdd(argsTail()) }},
		{[]string{constants.CmdLlm}, func() error { return cmdllm.RunLlm(argsTail()) }},
		{[]string{"locate", "find-tool", "vcvars"}, func() error { return RunLocateTopLevel(argsTail()) }},
		{[]string{constants.CmdFind, "f"}, func() error { return cmdfind.RunFind(argsTail()) }},
		{[]string{constants.CmdFindFiles, "ff"}, func() error { return cmdfind.RunFindFiles(argsTail()) }},
		{[]string{constants.CmdFindFilesAny, "ffa"}, func() error { return cmdfind.RunFindFilesAny(argsTail()) }},
		{[]string{constants.CmdFindFilesStartsWith, "ffs"}, func() error { return cmdfind.RunFindFilesStartsWith(argsTail()) }},
		{[]string{constants.CmdFindFilesEndsWith, "ffe"}, func() error { return cmdfind.RunFindFilesEndsWith(argsTail()) }},
		{[]string{constants.CmdListFiles, "lf"}, func() error { return cmdfind.RunListFiles(argsTail()) }},
		{[]string{constants.CmdFindRegex}, func() error { return cmdfind.RunFindRegex(argsTail()) }},
		{[]string{constants.CmdFindRead}, func() error { return cmdfind.RunFindRead(argsTail()) }},
		{[]string{constants.CmdFindReadJson}, func() error { return cmdfind.RunFindReadJson(argsTail()) }},
		{[]string{constants.CmdFindRegexRead}, func() error { return cmdfind.RunFindRegexRead(argsTail()) }},
		{[]string{constants.CmdFindRegexReadJson}, func() error { return cmdfind.RunFindRegexReadJson(argsTail()) }},
		{[]string{constants.CmdFindHelp}, func() error { cmdfind.RenderFindHelp(); return nil }},
		{[]string{constants.CmdSearchHelp}, func() error { cmdsearch.RenderSearchHelp(); return nil }},
		{[]string{constants.CmdRegexHelp}, func() error { cmdsearch.RenderSearchHelp(); return nil }},
		{[]string{constants.CmdSearch}, func() error { return cmdsearch.RunSearch(argsTail()) }},
		{[]string{constants.CmdReplace}, func() error { return cmdreplace.RunReplace(argsTail()) }},
		{[]string{constants.CmdReplaceRegex}, func() error { return cmdsearch.RunReplaceRegex(argsTail()) }},
		{[]string{constants.CmdRepoSearch, constants.CmdRepoSearchAlias}, func() error { return cmdsearch.RunRepoSearch(argsTail()) }},
		{[]string{"_index"}, func() error { return cmdindex.RunIndex(argsTail()) }},
		{[]string{constants.CmdRepoRegex, constants.CmdRepoRegexAlias}, func() error { return cmdsearch.RunRepoRegex(argsTail()) }},
		{[]string{constants.CmdRepoSearchJson, constants.CmdRepoSearchJsonAlias}, func() error { return cmdsearch.RunRepoSearchJson(argsTail()) }},
		{[]string{constants.CmdRepoSearchRegexJson}, func() error { return cmdsearch.RunRepoSearchRegexJson(argsTail()) }},
		{[]string{constants.CmdSearchReplaceAll}, func() error { return cmdsearch.RunSearchReplaceAll(argsTail()) }},
		{[]string{"asset", "download-prnt", "prnt", "prnt-download"}, func() error { return cmdasset.Run(argsTail()) }},
		{[]string{constants.CmdZip}, func() error { return cmdzip.RunZip(argsTail()) }},
		{[]string{"mkdir"}, func() error { return RunMkdir(argsTail()) }},
		{[]string{"cat"}, func() error { return cmdmacro.RunCat(argsTail()) }},
		{[]string{constants.CmdAppend}, func() error { return RunAppend(argsTail()) }},
		{[]string{constants.CmdWrite}, func() error { return RunWrite(argsTail()) }},
		{[]string{constants.CmdHead}, func() error { return RunHead(argsTail()) }},
		{[]string{constants.CmdTail}, func() error { return RunTail(argsTail()) }},
		{[]string{constants.CmdFileSearch}, func() error { return runFileSearch(argsTail()) }},
		{[]string{"reset-and-rescan"}, func() error { return runReset([]string{"--confirm", "--rescan"}) }},
		{[]string{constants.CmdReplace, constants.CmdReplaceAlias}, func() error { return cmdreplace.RunReplace(argsTail()) }},
		{[]string{constants.CmdRegoldens, constants.CmdRegoldensAlias}, func() error { return cmdregoldens.RunRegoldens(argsTail()) }},
		{[]string{constants.CmdAuditLegacy, constants.CmdAuditLegacyAlias, constants.CmdAuditLegacyAlias2}, func() error { return cmdauditlegacy.RunAuditLegacy(argsTail()) }},
		{[]string{constants.CmdFixRepo, constants.CmdFixRepoAlias}, func() error { return cmdfixrepo.RunFixRepo(argsTail()) }},
		{[]string{constants.CmdFixGit, constants.CmdFixGitAlias, "--fix-git", "fixgit"}, func() error { return cmdfixgit.RunFixGit(argsTail()) }},
		{[]string{constants.CmdUndo, constants.CmdUndoAlias}, func() error { return cmdundo.RunUndo(argsTail()) }},
		{[]string{constants.CmdHistoryPurge, constants.CmdHistoryPurgeAlias}, func() error { return cmdhistory.RunHistoryPurge(argsTail()) }},
		{[]string{constants.CmdHistoryPin, constants.CmdHistoryPinAlias}, func() error { return cmdhistory.RunHistoryPin(argsTail()) }},
		{[]string{"author"}, func() error { return cmdauthor.RunAuthor(argsTail()) }},
		{[]string{"sponsor"}, func() error { return cmdauthor.RunSponsor(argsTail()) }},
		{[]string{"credits"}, func() error { return cmdauthor.RunCredits(argsTail()) }},
	}
}

func toolingChromeEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdChromeProfileCopy}, func() error { return cmdchromeprofile.RunProfileCopy(argsTail()) }},
		{[]string{constants.CmdChromeProfileExport, constants.CmdChromeProfileExportAlias}, func() error { return cmdchromeprofile.RunProfileExport(argsTail()) }},
		{[]string{constants.CmdChromeProfileImport, constants.CmdChromeProfileImportAlias}, func() error { return cmdchromeprofile.RunProfileImport(argsTail()) }},
		{[]string{constants.CmdChromeProfileList, constants.CmdChromeProfileListAlias, constants.CmdChromeProfileListAlias2}, func() error { return cmdchromeprofile.RunProfileList(argsTail()) }},
		{[]string{constants.CmdChromeProfileDelete, constants.CmdChromeProfileDeleteAlias}, func() error { return cmdchromeprofile.RunProfileDelete(argsTail()) }},
		{[]string{constants.CmdChromeProfileMerge, constants.CmdChromeProfileMergeAlias}, func() error { return cmdchromeprofile.RunProfileMerge(argsTail()) }},
		{[]string{"chrome-profile-copy-all", "cpc-all", "copy-all"}, func() error { return cmdchromeprofile.RunCopyAll(argsTail()) }},
		{[]string{"chrome-profile-export-all", "cpe-all", "export-all"}, func() error { return cmdchromeprofile.RunExportAll(argsTail()) }},
		{[]string{"chrome-profile-import-all", "cpi-all", "import-all"}, func() error { return cmdchromeprofile.RunImportAll(argsTail()) }},
		{[]string{constants.CmdChrome, constants.CmdChromeAlias, constants.CmdChromeAlias2}, func() error { return cmdchrome.RunChrome(argsTail()) }},
	}
}

func runCpcDisambiguated(args []string) error {
	if cmdcommit.IsCommitPushHelpArg(args) {
		return cmdcommit.RunCommitPushChore(args)
	}

	if cmdpull.IsGitRepoCWD() {
		return cmdcommit.RunCommitPushChore(args)
	}

	return cmdchromeprofile.RunProfileCopy(args)
}

func toolingNetworkEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdServe, constants.CmdServeAlias}, func() error { return cmdservercmd.RunServe(argsTail()) }},
		{[]string{constants.CmdJoin, constants.CmdJoinAlias}, func() error {
			if err := cmdjoin.RunJoin(argsTail()); err != nil {
				return err
			}
			return nil
		}},
		{[]string{"ip"}, func() error { return runIP(argsTail()) }},
		{[]string{"ports", "port", "open-ports", "listening-ports"}, func() error { return cmdports.Run(argsTail()) }},
	}
}

func handleWhichSubcommand(args []string) error {
	if len(args) == 0 {
		return cmdconfig.RunWhichFormatCLI(args)
	}
	first := strings.ToLower(args[0])
	if first == "config" || first == "configs" {
		return cmdconfig.RunWhatConfigsCLI(args[1:])
	}
	if first == "format" || first == "json" {
		return cmdconfig.RunWhichFormatCLI(args[1:])
	}
	return cmdconfig.RunWhichFormatCLI(args)
}

func handleFormatSubcommand(args []string) error {
	if len(args) == 0 {
		return cmdconfig.RunWhichFormatCLI(args)
	}
	first := strings.ToLower(args[0])
	if first == "merge" {
		return cmdmerge.RunMergeJSONCLI(args[1:])
	}
	if first == "import" || first == "import-all" || first == "importall" {
		return cmdimport.RunImportAllJSONCLI(args[1:])
	}
	if first == "which" || first == "inspect" || first == "check" {
		return cmdconfig.RunWhichFormatCLI(args[1:])
	}
	return cmdconfig.RunWhichFormatCLI(args)
}

func handleJSONSubcommand(args []string) error {
	if len(args) == 0 {
		return cmdconfig.RunWhichFormatCLI(args)
	}
	first := strings.ToLower(args[0])
	if first == "merge" {
		return cmdmerge.RunMergeJSONCLI(args[1:])
	}
	if first == "import" || first == "import-all" || first == "importall" {
		return cmdimport.RunImportAllJSONCLI(args[1:])
	}
	return cmdconfig.RunWhichFormatCLI(args)
}

func toolingSystemEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"history clean", "history-clean"}, func() error {
			return cmdhistory.RunHistoryPurgeCLI(argsTail())
		}},
		{[]string{"history-undo", "hu", "history restore", "history-restore", "undo-history"}, func() error {
			return cmdhistory.RunHistoryUndoCLI(argsTail())
		}},
		{[]string{"history"}, func() error {
			return handleHistoryParentDispatcher(argsTail())
		}},
	}
}

func handleHistoryParentDispatcher(args []string) error {
	if len(args) == 0 {
		return cmdhistory.RunHistory(args)
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "purge", "clean", "p":
		return cmdhistory.RunHistoryPurgeCLI(args[1:])
	case "undo", "restore", "u":
		return cmdhistory.RunHistoryUndoCLI(args[1:])
	default:
		return cmdhistory.RunHistory(args)
	}
}
