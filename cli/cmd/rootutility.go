package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautofix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbrowse"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddocs"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfindnext"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixauth"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixreleasetags"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhelp"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlfscommon"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdllm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprobe"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdui"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddownload"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzsh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
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

func utilityToolEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdDocs, constants.CmdDocsAlias}, func() error { return cmddocs.RunDocs(argsTail()) }},
		{[]string{constants.CmdHelpDashboard, constants.CmdHelpDashboardAlias}, func() error { return runHelpDashboard(argsTail()) }},
		{[]string{constants.CmdLLMDocs, constants.CmdLLMDocsAlias, "ld"}, func() error { return cmdllm.RunLLMDocs(argsTail()) }},
		{[]string{"llm-train", "train", "llmtrain"}, func() error { return cmdllm.RunLlm(append([]string{"train"}, argsTail()...)) }},
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
		{[]string{"copy", "cp-mem", "copy-to-memory"}, func() error { return cmdcopy.RunCopyCmd(argsTail()) }},
		{[]string{"paste", "paste-mem"}, func() error { return cmdcopy.RunPasteCmd(argsTail()) }},
		{[]string{"explorer", "open-explorer", "folder", "open-folder", "browse-folder"}, func() error { return RunExplorerCmd(argsTail()) }},
		{[]string{"open-url", "browse", "browse-url", "open-browser"}, func() error { return cmdbrowse.RunBrowseCmd(argsTail()) }},
		{[]string{"cat", "view", "type"}, func() error { return cmdmacro.RunCatCmd(argsTail()) }},
		{[]string{"touch"}, func() error { return cmdmacro.RunTouchCmd(argsTail()) }},
		{[]string{"mkfile", "create-file"}, func() error { return cmdmacro.RunMkfileCmd(argsTail()) }},
		{[]string{"fix"}, func() error { return cmdautofix.RunFixCmd(argsTail()) }},
	}
}

// runHelpDispatch handles the `help` subcommand including topic
// help, `help example <cmd>`, --groups, --compact, and the default
// usage screen.
func runHelpDispatch() error {
	isHelpExample := len(os.Args) >= 3 && strings.EqualFold(os.Args[2], "example")
	if isHelpExample {
		dispatchHelpExample(os.Args[3:])

		return nil
	}

	// Program 275: `gitmap help topic <name>` routes to the topic index
	// before the generic topic lookup below.
	if len(os.Args) >= 3 && strings.EqualFold(os.Args[2], "topic") {
		if appErr := cmdhelp.RunHelpTopic(os.Args[3:]); appErr != nil {
			return appErr
		}

		return nil
	}

	hasTopic := len(os.Args) >= 3 && !isFlagToken(os.Args[2])
	if hasTopic {
		dispatchHelpTopic(os.Args[2])

		return nil
	}

	if cmdupdate.HasFlag(constants.FlagJSON) {
		printUsageJSON(resolveFilterQuery())

		return nil
	}

	q := resolveFilterQuery()
	needsFilter := len(q) > 0 || cmdupdate.HasFlag(constants.FlagFilter) || cmdupdate.HasFlag(constants.FlagFilterShort)
	if needsFilter {
		printUsageFiltered(q)

		return nil
	}

	printUsage()

	return nil
}

// dispatchHelpTopic renders help for a named topic, falling back to filtered
// usage when the topic has no dedicated help text. The registry's Examples
// section is appended after the topic body (single source of truth).
func dispatchHelpTopic(rawTopic string) {
	topic := normalizeHelpTopic(rawTopic)
	if isHelpCategoryGroup(topic) {
		printUsageCategoryGroup(topic)
		return
	}
	_, mode := ParsePrettyFlag(os.Args[3:])
	prettyOff := mode == render.PrettyOff
	if tryRenderRichTopic(topic) {
		helpdoc.PrintExampleSection(topic, prettyOff)
		return
	}
	_, err := helpdoc.ReadRaw(topic)
	if err != nil {
		printUsageFiltered(topic)

		return
	}

	helpdoc.PrintWithMode(topic, mode)
	helpdoc.PrintExampleSection(topic, prettyOff)
	printUsageFooterShort()
}

// dispatchHelpExample handles `gitmap help example [cmd]`: the canonical
// what-it-does + invocation + expected output for one command, or the
// command index when no command is given.
func dispatchHelpExample(args []string) {
	cleaned, mode := ParsePrettyFlag(args)
	prettyOff := mode == render.PrettyOff

	if len(cleaned) == 0 {
		helpdoc.PrintHelpExampleIndex(prettyOff)

		return
	}

	helpdoc.PrintHelpExample(cleaned[0], prettyOff)
}

func normalizeHelpTopic(topic string) string {
	switch topic {
	case constants.CmdRmAlias, constants.CmdRmAlias2:
		return constants.CmdRm
	case constants.CmdStaleAlias:
		return constants.CmdStale
	case constants.CmdRecentAlias:
		return constants.CmdRecent
	case constants.CmdPRAlias:
		return constants.CmdPR
	case constants.CmdClusterAlias:
		return constants.CmdCluster
	case constants.CmdPowerAlias, constants.CmdPowerAlias2:
		return constants.CmdPower
	case "rr", "rra", "rrq", "rerun-restart", "rerun-all", "rerun-queue":
		return "rerun"
	case constants.CmdShutdownUntil, constants.CmdShutdownUntilAlias, constants.CmdShutdownUntilGreen:
		return constants.CmdShutdownUntil
	case constants.CmdFixReleaseTagsAlias, "fix-release-tag", "fix-tags", constants.CmdFixReleaseTagsAliasSolid:
		return constants.CmdFixReleaseTags
	}

	return topic
}
