package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautofix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbrowse"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy"
	cmdcopyconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddocs"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfindnext"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixauth"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixreleasetags"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdflatten"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlfs"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlfscommon"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdllm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprobe"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsplit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdui"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdurl"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddownload"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzsh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// isFlagToken returns true when arg looks like a CLI flag (-x or --xx).
func isFlagToken(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

// dispatchUtility routes setup, update, doctor, and other utility commands.
func dispatchUtility(command string) (bool, error) {
	return runDispatchTable(command, utilityDispatchEntries())
}

// utilityDispatchEntries returns the routing table for utility commands.
func utilityDispatchEntries() []dispatchEntry {
	entries := make([]dispatchEntry, 0, 45)
	entries = append(entries, utilityCoreEntries()...)
	entries = append(entries, utilityToolEntries()...)
	entries = append(entries, utilitySystemEntries()...)
	entries = append(entries, utilityPipelineEntries()...)
	entries = append(entries, utilityDesktopEntries()...)

	return entries
}

func utilityCoreEntries() []dispatchEntry {
	entries := make([]dispatchEntry, 0, 15)
	entries = append(entries, utilityCoreBaseEntries()...)
	entries = append(entries, utilityUrlEntries()...)

	return entries
}

func utilityUrlEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"url"}, func() error { return cmdurl.RunUrlCommand(argsTail()) }},
		{[]string{"read-url", "readurl"}, func() error { return cmdurl.RunUrlCommand(append([]string{"read"}, argsTail()...)) }},
		{[]string{"get", "get-url"}, func() error { return cmdurl.RunUrlCommand(append([]string{"read"}, argsTail()...)) }},
	}
}

func utilityCoreBaseEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"binary", "info"}, printIdentityLong},
		{[]string{"error"}, func() error { return cmderrors.RunErrorCmd(argsTail()) }},
		{[]string{constants.CmdUpdate, "ua", "update-all", "updateall", "uaz", "update-all-zip", "updateallzip"}, runUpdateHelp},
		{[]string{constants.CmdUpdateRunner}, cmdupdate.RunUpdateRunner},
		{[]string{constants.CmdUpdateCleanup}, cmdupdate.RunUpdateCleanup},
		{[]string{constants.CmdInstalledDir, constants.CmdInstalledDirAlias}, runInstalledDirHelp},
		{[]string{constants.CmdRevert}, func() error { return cmdrevert.RunRevert(argsTail()) }},
		{[]string{constants.CmdRm, constants.CmdRmAlias, constants.CmdRmAlias2}, func() error { return cmdrm.RunRm(argsTail()) }},
		{[]string{constants.CmdRevertRunner}, cmdrevert.RunRevertRunner},
		{[]string{constants.CmdVersion, constants.CmdVersionAlias, "--version", "-version", "-v", "versions"}, printVersionBlock},
		{[]string{constants.CmdHelp, "--help", "-h"}, runHelpDispatch},
	}
}

func printIdentityLong() error {
	printGitmapIdentityBlockLong()

	return nil
}

func printVersionBlock() error {
	args := argsTail()
	if isVersionListRequest(args) {
		return RunGitMapVersionTagsLS()
	}
	checkHelp("version", args)
	fmt.Printf(constants.MsgVersionFmt, constants.Version)

	return nil
}

func runInstalledDirHelp() error {
	checkHelp("installed-dir", argsTail())

	return cmdinstall.RunInstalledDir()
}

const (
	CmdResolve         = "resolve"
	CmdConflicts       = "conflicts"
	CmdSplitLargeFiles = "split-large-files"
	CmdLFSConvert      = "lfs-convert"
	CmdFlatten         = "flatten"
	CmdDb              = "db"
)

func utilityToolEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{CmdDb, "database"}, func() error { return cmddb.RunDb(argsTail()) }},
		{[]string{CmdResolve, "res"}, func() error { return cmdconflict.RunResolveCmd(argsTail()) }},
		{[]string{CmdConflicts, "conflict"}, func() error { return cmdconflict.RunConflictsCmd(argsTail()) }},
		{[]string{CmdSplitLargeFiles, "split-large", "splitlargefiles"}, func() error { return cmdsplit.RunSplitLargeFilesCmd(argsTail()) }},
		{[]string{CmdLFSConvert, "lfsconvert"}, func() error { return cmdlfs.RunLFSConvertCmd(argsTail()) }},
		{[]string{CmdFlatten, "repo-flatten", "repoflatten"}, func() error { return cmdflatten.RunFlattenCmd(argsTail()) }},
		{[]string{"undo-flatten", "repo-undo-flatten"}, func() error { return cmdundo.RunUndoFlattenCLI(argsTail()) }},
		{[]string{constants.CmdDocs, constants.CmdDocsAlias}, func() error { return cmddocs.RunDocs(argsTail()) }},
		{[]string{constants.CmdHelpDashboard, constants.CmdHelpDashboardAlias}, func() error { return runHelpDashboard(argsTail()) }},
		{[]string{constants.CmdLLMDocs, constants.CmdLLMDocsAlias, "ld"}, func() error { return cmdllm.RunLLMDocs(argsTail()) }},
		{[]string{"llm-train", "train", "llmtrain"}, func() error { return cmdllm.RunLlm(append([]string{"train"}, argsTail()...)) }},
		{[]string{"train-muse", "llm-train-muse"}, func() error { return cmdllm.RunLlm(append([]string{"train", "muse"}, argsTail()...)) }},
		{[]string{"ai-server", "ai-ping", "aum-server"}, func() error { return cmdai.RunAIMemoryServerCmd(argsTail()) }},
		{[]string{constants.CmdSetSourceRepo}, RunSetSourceRepo},
		{[]string{constants.CmdSf}, func() error { return RunSf(argsTail()) }},
		{[]string{constants.CmdProbe}, func() error { return cmdprobe.RunProbe(argsTail()) }},
		{[]string{"vscode", "vsc"}, func() error { return cmdvscode.RunVSCode(argsTail()) }},
		{[]string{"ide", "ides"}, func() error { return cmdide.RunIDE(argsTail()) }},
		{[]string{constants.CmdFindNext, constants.CmdFindNextAlias}, func() error { return cmdfindnext.RunFindNext(argsTail()) }},
		{[]string{constants.CmdVSCodePMPath, constants.CmdVSCodePMPathAlias}, func() error { return cmdvscode.RunVSCodePMPath(argsTail()) }},
		{[]string{constants.CmdVSCodeWorkspace, constants.CmdVSCodeWorkspaceAlias}, func() error { return cmdvscode.RunVSCodeWorkspace(argsTail()) }},
		{[]string{constants.CmdLFSCommon, constants.CmdLFSCommonAlias}, func() error { return cmdlfscommon.RunLFSCommon(argsTail()) }},
		{[]string{constants.CmdReinstall}, func() error { return runReinstall(argsTail()) }},
		{[]string{"peat", "pea"}, func() error { return runPeatCmd(argsTail()) }},
		{[]string{"install-exec", "in-exec", "setup-exec"}, func() error { return cmdssh.RunSSHInstallExecCLI(argsTail()) }},
		{[]string{"py", "python"}, func() error { return cmdpy.RunPy(argsTail()) }},
		{[]string{constants.CmdPlans, constants.CmdPlansAlias}, func() error { return cmdplans.RunPlans(argsTail()) }},
		{[]string{"download", "dl"}, func() error { return cmddownload.RunDownloadCLI(argsTail()) }},
	}
}

func runPeatCmd(args []string) error {
	if len(args) > 0 && args[0] == "deploy" {
		return cmdmacro.ExecuteMacroDeploySSH(args[1:])
	}
	return cmdmacro.RunMacroCmd(args)
}

func utilitySystemEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"cpu", "processor", "proc"}, func() error { return cmdos.RunCPUCLI(argsTail()) }},
		{[]string{"ram", "mem", "memory"}, func() error { return cmdos.RunRAMCLI(argsTail()) }},
		{[]string{"disk"}, func() error { return cmdos.RunDiskCLI(argsTail()) }},
		{[]string{constants.CmdPower, constants.CmdPowerAlias, constants.CmdPowerAlias2}, func() error { return cmdos.RunPower(argsTail()) }},
		{[]string{constants.CmdOS, "linutil"}, func() error { return cmdos.RunOSCLI(argsTail()) }},
		{[]string{"autologin", "auto-login", "al"}, func() error { return cmdos.RunOSCLI(append([]string{"autologin"}, argsTail()...)) }},
		{[]string{"zsh"}, func() error { return runZsh(argsTail()) }},
		{[]string{constants.CmdFixLink, constants.CmdFixLinkAlias, constants.CmdFixLinkAlias2}, func() error { return cmdos.RunOSFixLink(argsTail()) }},
		{[]string{constants.CmdWhoAmI, constants.CmdWhoAmIAlias}, func() error { checkHelp("whoami", argsTail()); return RunWhoAmI(argsTail()) }},
		{[]string{constants.CmdSSHBind, constants.CmdSSHBindAlias}, func() error { checkHelp("ssh-bind", argsTail()); return cmdssh.RunSSHBind(argsTail()) }},
		{[]string{constants.CmdFixAuth, constants.CmdFixAuthAlias}, func() error { checkHelp("fix-auth", argsTail()); return cmdfixauth.RunFixAuth(argsTail()) }},
		{[]string{constants.CmdFixCredential, constants.CmdFixCredentialAlias}, func() error { checkHelp("fix-credential", argsTail()); return cmdfixauth.RunFixCredential(argsTail()) }},
		{[]string{constants.CmdFixReleaseTags, constants.CmdFixReleaseTagsAlias, "fix-release-tag", "fix-tags", constants.CmdFixReleaseTagsAliasSolid}, func() error { return cmdfixreleasetags.RunFixReleaseTags(argsTail()) }},
		{[]string{"ai-clean", "aiclean", "clean-ai"}, func() error { return cmdos.RunOSAICleanCLI(argsTail()) }},
		{[]string{constants.CmdShutdownUntil, constants.CmdShutdownUntilAlias, constants.CmdShutdownUntilGreen}, func() error { return cmdagy.RunSUGCLI(argsTail()) }},
		{[]string{constants.CmdWatchPromptsRunning, constants.CmdWatchPromptsRunningAlias, "watch-running-prompts"}, func() error { return cmdagy.RunWPRCLI(argsTail()) }},
	}
}

func runZsh(args []string) error {
	checkHelp("zsh", args)

	return cmdzsh.RunZsh(args)
}

func utilityPipelineEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"pull-fix", "pf"}, func() error { return cmdpull.RunBatchPullFix(argsTail()) }},
		{[]string{"pipeline-fix", "fix-pipeline", "aef", "agy-errors-fix", "fix-agy"}, func() error { return cmdagy.RunPipelineFixAgyCLI(argsTail()) }},
		{[]string{"pipeline-ai", "pl-ai", "plai", "pipeline_ai"}, func() error { return cmdpipeline.RunPipelineAI(argsTail()) }},
		{[]string{"pipeline", "pipelines", "pl"}, func() error { return cmdpipeline.RunPipeline(argsTail()) }},
		{[]string{"pe", "pipeline-errors", "pipeline_errors", "ee"}, func() error { return cmdpipeline.RunPipelineErrors(argsTail()) }},
		{[]string{"pipe-error-all", "pipe-errors-all"}, func() error { return cmdpipeline.RunPipelineErrors(append([]string{"all"}, argsTail()...)) }},
		{[]string{"te"}, func() error { return cmdpipeline.RunPipelineErrors(argsTail()) }},
		{[]string{"pd", "pipeline-details", "pipeline_details"}, func() error { return cmdpipeline.RunPipelineDetails(argsTail()) }},
		{[]string{"e", "errors", "internal-errors", "errs"}, func() error { return cmderrors.RunErrorsCLI(argsTail()) }},
		{[]string{"failed-commands", "failed-command", "fc", "unknown-commands", "unknown-command", "failed-to-detect", "failed-to-detect-commands", "failed-commands-count", "fcc"}, func() error {
			return cmderrors.RunFailedCommandsWithCmd(os.Args[1], argsTail())
		}},
		{[]string{"error-logs", "error-log", "errorlogs", "errorlog", "errorslogs", "errors-log", "errors-logs", "last-failed-logs"}, func() error { return cmdpipeline.RunPipeline(append([]string{os.Args[1]}, argsTail()...)) }},
		{[]string{"logs", "log"}, func() error { return cmdpipeline.RunPipeline(append([]string{"logs"}, argsTail()...)) }},
		{[]string{"waittime", "wait-time", "eta"}, func() error { return cmdpipeline.RunPipeline(append([]string{"waittime"}, argsTail()...)) }},
		{[]string{"repo"}, func() error { return cmdrepo.RunRepoCommand(argsTail()) }},
		{[]string{"ui"}, func() error { return cmdui.RunUICmd(argsTail()) }},
	}
}

func utilityDesktopEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{cmdcopyconstants.CmdCopy, cmdcopyconstants.CmdCopyAliasMem, cmdcopyconstants.CmdCopyAliasToMem}, func() error { return cmdcopy.RunCopyCmd(argsTail()) }},
		{[]string{cmdcopyconstants.CmdPaste, cmdcopyconstants.CmdPasteAliasMem}, func() error { return cmdcopy.RunPasteCmd(argsTail()) }},
		{[]string{"explorer", "open-explorer", "folder", "open-folder", "browse-folder"}, func() error { return RunExplorerCmd(argsTail()) }},
		{[]string{"open-url", "browse", "browse-url", "open-browser"}, func() error { return cmdbrowse.RunBrowseCmd(argsTail()) }},
		{[]string{"cat", "view", "type"}, func() error { return cmdmacro.RunCatCmd(argsTail()) }},
		{[]string{"touch"}, func() error { return cmdmacro.RunTouchCmd(argsTail()) }},
		{[]string{"mkfile", "create-file"}, func() error { return cmdmacro.RunMkfileCmd(argsTail()) }},
		{[]string{"fix"}, func() error { return cmdautofix.RunFixCmd(argsTail()) }},
	}
}
