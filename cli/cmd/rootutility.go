package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzsh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
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
		{[]string{"error"}, func() error { return runErrorCmd(argsTail()) }},
		{[]string{constants.CmdUpdate}, runUpdateHelp},
		{[]string{constants.CmdUpdateRunner}, runUpdateRunner},
		{[]string{constants.CmdUpdateCleanup}, runUpdateCleanup},
		{[]string{constants.CmdInstalledDir, constants.CmdInstalledDirAlias}, runInstalledDirHelp},
		{[]string{constants.CmdRevert}, func() error { return runRevert(argsTail()) }},
		{[]string{constants.CmdRm, constants.CmdRmAlias, constants.CmdRmAlias2}, func() error { return runRm(argsTail()) }},
		{[]string{constants.CmdRevertRunner}, runRevertRunner},
		{[]string{constants.CmdVersion, constants.CmdVersionAlias, "--version", "-version", "-v"}, printVersionBlock},
		{[]string{constants.CmdHelp, "--help", "-h"}, runHelpDispatch},
	}
}

func printIdentityLong() error {
	printGitmapIdentityBlockLong()

	return nil
}

func printVersionBlock() error {
	checkHelp("version", argsTail())
	fmt.Printf(constants.MsgVersionFmt, constants.Version)

	return nil
}

func runInstalledDirHelp() error {
	checkHelp("installed-dir", argsTail())

	return runInstalledDir()
}

func extractRemoteUpdateTarget(args []string) (string, []string) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if (arg == "--remote" || arg == "-r" || arg == "--node" || arg == "-n") && i+1 < len(args) {
			target := args[i+1]
			clean := append(args[:i], args[i+2:]...)
			return target, clean
		}
		if strings.HasPrefix(arg, "--remote=") {
			return strings.TrimPrefix(arg, "--remote="), append(args[:i], args[i+1:]...)
		}
		if strings.HasPrefix(arg, "--node=") {
			return strings.TrimPrefix(arg, "--node="), append(args[:i], args[i+1:]...)
		}
	}
	return "", args
}

func runUpdateHelp() error {
	checkHelp("update", argsTail())
	remoteTarget, cleanArgs := extractRemoteUpdateTarget(argsTail())
	if remoteTarget != "" {
		pkg := resolveUpdatePackage(cleanArgs)
		return cmdssh.RunSSHUpdateCLI([]string{pkg, remoteTarget})
	}
	if isAgmUpdateTarget(argsTail()) {
		return runUpdateAgManagerTarget(argsTail())
	}

	return runUpdate()
}

func resolveUpdatePackage(args []string) string {
	if isAgmUpdateTarget(args) {
		return "agm"
	}
	return "gitmap"
}

func isAgmUpdateTarget(args []string) bool {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") && isAgmToken(arg) {
			return true
		}
	}
	return false
}

func isAgmToken(arg string) bool {
	low := strings.ToLower(arg)
	return low == "agm" || low == "ag-manager" || low == "antigravity-manager"
}

func runUpdateAgManagerTarget(args []string) error {
	remoteTarget, _ := extractRemoteUpdateTarget(args)
	if remoteTarget != "" {
		return cmdssh.RunSSHUpdateCLI([]string{"agm", remoteTarget})
	}
	opts := cmdinstall.InstallOptions{
		DryRun: hasDryRunArg(args),
	}
	return cmdinstall.RunUpdateAgManagerWithOpts(opts)
}

func hasDryRunArg(args []string) bool {
	for _, a := range args {
		if a == "--dry-run" || a == "-n" {
			return true
		}
	}
	return false
}

func utilityToolEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdDocs, constants.CmdDocsAlias}, func() error { return runDocs(argsTail()) }},
		{[]string{constants.CmdHelpDashboard, constants.CmdHelpDashboardAlias}, func() error { return runHelpDashboard(argsTail()) }},
		{[]string{constants.CmdLLMDocs, constants.CmdLLMDocsAlias}, func() error { return runLLMDocs(argsTail()) }},
		{[]string{constants.CmdSetSourceRepo}, runSetSourceRepo},
		{[]string{constants.CmdSf}, func() error { return runSf(argsTail()) }},
		{[]string{constants.CmdProbe}, func() error { return runProbe(argsTail()) }},
		{[]string{"vscode", "vsc"}, func() error { return runVSCode(argsTail()) }},
		{[]string{constants.CmdFindNext, constants.CmdFindNextAlias}, func() error { return runFindNext(argsTail()) }},
		{[]string{constants.CmdVSCodePMPath, constants.CmdVSCodePMPathAlias}, func() error { return runVSCodePMPath(argsTail()) }},
		{[]string{constants.CmdVSCodeWorkspace, constants.CmdVSCodeWorkspaceAlias}, func() error { return runVSCodeWorkspace(argsTail()) }},
		{[]string{constants.CmdLFSCommon, constants.CmdLFSCommonAlias}, func() error { return runLFSCommon(argsTail()) }},
		{[]string{constants.CmdReinstall}, func() error { return runReinstall(argsTail()) }},
	}
}

func utilitySystemEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdPower, constants.CmdPowerAlias, constants.CmdPowerAlias2}, func() error { return runPower(argsTail()) }},
		{[]string{constants.CmdOS}, func() error { return runOS(argsTail()) }},
		{[]string{"zsh"}, func() error { return runZsh(argsTail()) }},
		{[]string{constants.CmdFixLink, constants.CmdFixLinkAlias, constants.CmdFixLinkAlias2}, func() error { return runOSFixLink(argsTail()) }},
		{[]string{constants.CmdWhoAmI, constants.CmdWhoAmIAlias}, func() error { checkHelp("whoami", argsTail()); return runWhoAmI(argsTail()) }},
		{[]string{constants.CmdSSHBind, constants.CmdSSHBindAlias}, func() error { checkHelp("ssh-bind", argsTail()); return runSSHBind(argsTail()) }},
		{[]string{constants.CmdFixAuth, constants.CmdFixAuthAlias}, func() error { checkHelp("fix-auth", argsTail()); return runFixAuth(argsTail()) }},
		{[]string{"ai-clean", "aiclean", "clean-ai"}, func() error { return cmdos.RunOSAICleanCLI(argsTail()) }},
	}
}

func runZsh(args []string) error {
	checkHelp("zsh", args)

	return cmdzsh.RunZsh(args)
}

func utilityPipelineEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"pipeline-fix", "fix-pipeline", "aef", "agy-errors-fix", "fix-agy"}, func() error { return cmdagy.RunPipelineFixAgyCLI(argsTail()) }},
		{[]string{"pipeline-ai", "pl-ai", "plai", "pipeline_ai"}, func() error { return runPipelineAI(argsTail()) }},
		{[]string{"pipeline", "pipelines", "pl"}, func() error { return runPipeline(argsTail()) }},
		{[]string{"pe", "pipeline-errors", "pipeline_errors"}, func() error { return runPipelineErrors(argsTail()) }},
		{[]string{"pd", "pipeline-details", "pipeline_details"}, func() error { return runPipelineDetails(argsTail()) }},
		{[]string{"error-logs", "error-log", "errorlogs", "errorlog", "errors", "err", "errorslogs", "errors-log", "errors-logs", "last-failed-logs"}, func() error { return runPipeline(append([]string{os.Args[1]}, argsTail()...)) }},
		{[]string{"logs", "log"}, func() error { return runPipeline(append([]string{"logs"}, argsTail()...)) }},
		{[]string{"waittime", "wait-time", "eta"}, func() error { return runPipeline(append([]string{"waittime"}, argsTail()...)) }},
		{[]string{"repo"}, func() error { return runRepoCommand(argsTail()) }},
		{[]string{"ui"}, func() error { return runUI(argsTail()) }},
	}
}

func utilityDesktopEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"copy", "cp-mem", "copy-to-memory"}, func() error { return runCopyCmd(argsTail()) }},
		{[]string{"paste", "paste-mem"}, func() error { return runPasteCmd(argsTail()) }},
		{[]string{"explorer", "open-explorer", "folder", "open-folder", "browse-folder"}, func() error { return runExplorerCmd(argsTail()) }},
		{[]string{"open-url", "browse", "browse-url", "open-browser"}, func() error { return runBrowseCmd(argsTail()) }},
		{[]string{"cat", "view", "type"}, func() error { return runCatCmd(argsTail()) }},
		{[]string{"touch"}, func() error { return runTouchCmd(argsTail()) }},
		{[]string{"mkfile", "create-file"}, func() error { return runMkfileCmd(argsTail()) }},
	}
}

// runHelpDispatch handles the `help` subcommand including topic
// help, --groups, --compact, and the default usage screen.
func runHelpDispatch() error {
	hasTopic := len(os.Args) >= 3 && !isFlagToken(os.Args[2])
	if hasTopic {
		dispatchHelpTopic(os.Args[2])

		return nil
	}

	if hasFlag(constants.FlagJSON) {
		printUsageJSON(resolveFilterQuery())

		return nil
	}

	q := resolveFilterQuery()
	needsFilter := len(q) > 0 || hasFlag(constants.FlagFilter) || hasFlag(constants.FlagFilterShort)
	if needsFilter {
		printUsageFiltered(q)

		return nil
	}

	printUsage()

	return nil
}

// dispatchHelpTopic renders help for a named topic, falling back to filtered
// usage when the topic has no dedicated help text.
func dispatchHelpTopic(rawTopic string) {
	topic := normalizeHelpTopic(rawTopic)
	if tryRenderRichTopic(topic) {
		return
	}
	_, err := helptext.ReadRaw(topic)
	if err != nil {
		printUsageFiltered(topic)

		return
	}

	_, mode := ParsePrettyFlag(os.Args[3:])
	helptext.PrintWithMode(topic, mode)
	printUsageFooterShort()
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
	}

	return topic
}
