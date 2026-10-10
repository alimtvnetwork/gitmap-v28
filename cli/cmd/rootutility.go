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
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsummary"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzsh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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

func extractRemoteUpdateTarget(args []string) (string, []string) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if (arg == "--remote" || arg == "-r" || arg == "--node" || arg == "-n") && i+1 < len(args) {
			target := args[i+1]
			clean := append([]string{}, args[:i]...)
			clean = append(clean, args[i+2:]...)
			return target, clean
		}
		if strings.HasPrefix(arg, "--remote=") {
			clean := append([]string{}, args[:i]...)
			clean = append(clean, args[i+1:]...)
			return strings.TrimPrefix(arg, "--remote="), clean
		}
		if strings.HasPrefix(arg, "--node=") {
			clean := append([]string{}, args[:i]...)
			clean = append(clean, args[i+1:]...)
			return strings.TrimPrefix(arg, "--node="), clean
		}
	}
	return "", args
}

func runUpdateHelp() error {
	cmdName := subcommandName()
	args := argsTail()
	if cmdupdate.IsFleetUpdateCommand(cmdName, args) {
		return cmdupdate.RunFleetUpdateDispatch(cmdName, args)
	}
	if len(args) > 0 && (args[0] == "ssh" || args[0] == "remote") {
		return cmdssh.RunSSHUpdateCLI(args[1:])
	}
	if hasSSHArg(args) {
		cleaned := stripSSHFlag(args)
		return cmdssh.RunSSHUpdateCLI(cleaned)
	}
	checkHelp("update", argsTail())

	if isVersionListRequest(args) && isAgmUpdateTarget(args) {
		return cmdagy.RunAGMVersionTagsLS()
	}
	if isVersionListRequest(args) {
		return RunGitMapVersionTagsLS()
	}

	remoteTarget, cleanArgs := extractRemoteUpdateTarget(args)
	if remoteTarget != "" {
		pkg := resolveUpdatePackage(cleanArgs)
		return cmdssh.RunSSHUpdateCLI([]string{pkg, remoteTarget})
	}

	targetVer := extractVersionFromArgs(args)
	applyPinIfRequested(args, targetVer)

	if isAgmUpdateTarget(args) {
		return runUpdateAgManagerTarget(args)
	}

	targetVer = resolveEffectiveVersion("gitmap", targetVer)
	if targetVer != "" {
		cmdupdate.SetTargetVersion(targetVer)
	}

	if len(args) > 0 && !isKnownUpdateTargetOrFlag(args[0]) {
		handleUnknownUpdateTarget(args[0])
		return nil
	}

	return cmdupdate.RunUpdate()
}

func applyPinIfRequested(args []string, targetVer string) {
	if !hasPinArg(args) || targetVer == "" {
		return
	}
	appName := "gitmap"
	if isAgmUpdateTarget(args) {
		appName = "agm"
	}
	_ = cmdinstaller.PinVersion(appName, targetVer)
}

func resolveEffectiveVersion(appName, requestedVer string) string {
	if requestedVer != "" {
		return requestedVer
	}
	pinned, err := cmdinstaller.GetPinnedVersion(appName)
	if err == nil && pinned != "" {
		return pinned
	}
	return ""
}

func hasPinArg(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "--pin" || low == "-pin" {
			return true
		}
	}
	return false
}

func hasSSHArg(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "--ssh" || low == "-ssh" || low == "--remote" {
			return true
		}
	}
	return false
}

func stripSSHFlag(args []string) []string {
	var cleaned []string
	for _, a := range args {
		low := strings.ToLower(a)
		if low != "--ssh" && low != "-ssh" && low != "--remote" {
			cleaned = append(cleaned, a)
		}
	}
	return cleaned
}

func isKnownUpdateTargetOrFlag(token string) bool {
	if strings.HasPrefix(token, "-") {
		return true
	}
	if isSemverLike(token) {
		return true
	}
	low := strings.ToLower(token)
	switch low {
	case "all", "all-nodes", "allnodes", "ls", "list", "ssh", "remote", "gitmap", "agm", "agy", "ag-manager", "antigravity-manager":
		return true
	case "zip", "uaz", "update-all-zip", "updateallzip", "--include-others", "--include-other":
		return true
	default:
		return false
	}
}

func suggestUpdateTarget(token string) string {
	low := strings.ToLower(token)
	switch low {
	case "al", "allnodes", "all-node", "nodes", "hosts", "cluster":
		return "all"
	case "l", "lss", "inventory", "apps", "soft":
		return "ls"
	case "ag", "ang", "manager":
		return "agm"
	case "git", "gitm", "gm":
		return "gitmap"
	case "rem", "ss":
		return "ssh"
	case "uaz", "update-all-zip", "updateallzip":
		return "update-all-zip"
	case "zp", "zi":
		return "zip"
	}
	return ""
}

func formatUnknownUpdateTargetHeader(token string) string {
	suggestion := suggestUpdateTarget(token)
	if suggestion != "" {
		return fmt.Sprintf("gitmap update: Unknown update target '%s'.\n  Did you mean: gitmap update %s?", token, suggestion)
	}
	return fmt.Sprintf("gitmap update: Unknown update target '%s'.", token)
}

func printUpdateUsageHelp() {
	fmt.Println("Available update commands:")
	fmt.Println("  gitmap update                          - Self-update local gitmap binary")
	fmt.Println("  gitmap update all                      - Update all fleet cluster nodes in parallel (alias: gitmap ua)")
	fmt.Println("  gitmap update all --include-others     - Update all cluster nodes and discovered hosts")
	fmt.Println("  gitmap update all zip --include-others - Distribute zip packages via SCP across all nodes")
	fmt.Println("  gitmap update-all-zip --include-others - Fast zip distribution across all nodes (alias: uaz)")
	fmt.Println("  gitmap update ls                       - List installed software inventory across fleet nodes")
	fmt.Println("  gitmap update agm                      - Update Antigravity Manager")
	fmt.Println("  gitmap update ssh <target>             - Update package on specified SSH target")
	fmt.Println()
}

func handleUnknownUpdateTarget(token string) {
	header := formatUnknownUpdateTargetHeader(token)
	sugg := suggestUpdateTarget(token)
	var suggList []string
	if sugg != "" {
		suggList = []string{"gitmap update " + sugg}
	} else {
		suggList = []string{"gitmap update", "gitmap update all", "gitmap update ls"}
	}
	store.LogFailedCommand(token, strings.Join(os.Args[1:], " "), "update", "E1001", header, suggList)
	fmt.Printf("\n%s\n\n", header)
	printUpdateUsageHelp()
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
	return low == "agm" || low == "agy" || low == "ag-manager" || low == "antigravity-manager"
}

func runUpdateAgManagerTarget(args []string) error {
	remoteTarget, cleanArgs := extractRemoteUpdateTarget(args)
	if remoteTarget != "" {
		return cmdssh.RunSSHUpdateCLI([]string{"agm", remoteTarget})
	}
	if isVersionListRequest(cleanArgs) {
		return cmdagy.RunAGMVersionTagsLS()
	}
	ver := extractVersionFromArgs(cleanArgs)
	ver = resolveEffectiveVersion("agm", ver)
	opts := cmdinstall.InstallOptions{
		DryRun:  hasDryRunArg(cleanArgs),
		Version: ver,
	}
	return cmdinstall.RunUpdateAgManagerWithOpts(opts)
}

func isVersionListRequest(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(strings.TrimSpace(a))
		if low == "ls" || low == "list" || low == "versions" || low == "tags" || low == "ls-versions" {
			return true
		}
	}
	return false
}

func isSemverLike(s string) bool {
	s = strings.TrimPrefix(s, "v")
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		if len(p) == 0 {
			return false
		}
		for _, r := range p {
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	return true
}

func extractVersionFromArgs(args []string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if (a == "--version" || a == "-v") && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(a, "--version=") {
			return strings.TrimPrefix(a, "--version=")
		}
		if !strings.HasPrefix(a, "-") && isSemverLike(a) {
			return a
		}
	}
	return ""
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
		{[]string{constants.CmdPower, constants.CmdPowerAlias, constants.CmdPowerAlias2}, func() error { return cmdos.RunPower(argsTail()) }},
		{[]string{constants.CmdOS, "linutil"}, func() error { return cmdos.RunOS(argsTail()) }},
		{[]string{"autologin", "auto-login", "al"}, func() error { return cmdos.RunOS(append([]string{"autologin"}, argsTail()...)) }},
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
		{[]string{"pe", "pipeline-errors", "pipeline_errors", "ee"}, func() error {
			args := argsTail()
			if len(args) > 0 && strings.EqualFold(args[0], "all") {
				return cmdsummary.RunPipelineErrorsAll(args[1:])
			}
			return cmdpipeline.RunPipelineErrors(args)
		}},
		{[]string{"pipe-error-all", "pipe-errors-all"}, func() error { return cmdsummary.RunPipelineErrorsAll(argsTail()) }},
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
// help, --groups, --compact, and the default usage screen.
func runHelpDispatch() error {
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
// usage when the topic has no dedicated help text.
func dispatchHelpTopic(rawTopic string) {
	topic := normalizeHelpTopic(rawTopic)
	if isHelpCategoryGroup(topic) {
		printUsageCategoryGroup(topic)
		return
	}
	if tryRenderRichTopic(topic) {
		return
	}
	_, err := helpdoc.ReadRaw(topic)
	if err != nil {
		printUsageFiltered(topic)

		return
	}

	_, mode := ParsePrettyFlag(os.Args[3:])
	helpdoc.PrintWithMode(topic, mode)
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
	case "rr", "rra", "rrq", "rerun-restart", "rerun-all", "rerun-queue":
		return "rerun"
	case constants.CmdShutdownUntil, constants.CmdShutdownUntilAlias, constants.CmdShutdownUntilGreen:
		return constants.CmdShutdownUntil
	case constants.CmdFixReleaseTagsAlias, "fix-release-tag", "fix-tags", constants.CmdFixReleaseTagsAliasSolid:
		return constants.CmdFixReleaseTags
	}

	return topic
}
