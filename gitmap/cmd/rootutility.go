package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/helptext"
)

// isFlagToken returns true when arg looks like a CLI flag (-x or --xx).
func isFlagToken(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

// dispatchUtility routes setup, update, doctor, and other utility commands.
func dispatchUtility(command string) bool {
	return runDispatchTable(command, utilityDispatchEntries())
}

// utilityDispatchEntries returns the routing table for utility commands.
func utilityDispatchEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdUpdate}, func() { checkHelp("update", argsTail()); runUpdate() }},
		{[]string{constants.CmdUpdateRunner}, func() { runUpdateRunner() }},
		{[]string{constants.CmdUpdateCleanup}, func() { runUpdateCleanup() }},
		{
			[]string{constants.CmdInstalledDir, constants.CmdInstalledDirAlias},
			func() { checkHelp("installed-dir", argsTail()); runInstalledDir() },
		},
		{[]string{constants.CmdRevert}, func() { runRevert(argsTail()) }},
		{[]string{constants.CmdRm, constants.CmdRmAlias, constants.CmdRmAlias2}, func() { runRm(argsTail()) }},
		{[]string{constants.CmdRevertRunner}, func() { runRevertRunner() }},
		{
			[]string{constants.CmdVersion, constants.CmdVersionAlias},
			func() { checkHelp("version", argsTail()); fmt.Printf(constants.MsgVersionFmt, constants.Version) },
		},
		{[]string{constants.CmdHelp}, runHelpDispatch},
		{[]string{constants.CmdDocs, constants.CmdDocsAlias}, func() { runDocs(argsTail()) }},
		{[]string{constants.CmdHelpDashboard, constants.CmdHelpDashboardAlias}, func() { runHelpDashboard(argsTail()) }},
		{[]string{constants.CmdLLMDocs, constants.CmdLLMDocsAlias}, func() { runLLMDocs(argsTail()) }},
		{[]string{constants.CmdSetSourceRepo}, func() { runSetSourceRepo() }},
		{[]string{constants.CmdSf}, func() { runSf(argsTail()) }},
		{[]string{constants.CmdProbe}, func() { runProbe(argsTail()) }},
		{[]string{"vscode", "vsc"}, func() { runVSCode(argsTail()) }},
		{[]string{constants.CmdFindNext, constants.CmdFindNextAlias}, func() { runFindNext(argsTail()) }},
		{[]string{constants.CmdVSCodePMPath, constants.CmdVSCodePMPathAlias}, func() { runVSCodePMPath(argsTail()) }},
		{[]string{constants.CmdVSCodeWorkspace, constants.CmdVSCodeWorkspaceAlias}, func() { runVSCodeWorkspace(argsTail()) }},
		{[]string{constants.CmdLFSCommon, constants.CmdLFSCommonAlias}, func() { runLFSCommon(argsTail()) }},
		{[]string{constants.CmdReinstall}, func() { runReinstall(argsTail()) }},
		{[]string{constants.CmdWhoAmI, constants.CmdWhoAmIAlias}, func() { checkHelp("whoami", argsTail()); runWhoAmI(argsTail()) }},
		{[]string{constants.CmdSSHBind, constants.CmdSSHBindAlias}, func() { checkHelp("ssh-bind", argsTail()); runSSHBind(argsTail()) }},
		{[]string{constants.CmdFixAuth, constants.CmdFixAuthAlias}, func() { checkHelp("fix-auth", argsTail()); runFixAuth(argsTail()) }},
	}
}

// runHelpDispatch handles the `help` subcommand including topic
// help, --groups, --compact, and the default usage screen.
func runHelpDispatch() {
	hasTopic := len(os.Args) >= 3 && isFlagToken(os.Args[2]) == false
	if hasTopic == true {
		dispatchHelpTopic(os.Args[2])
		return
	}

	if hasFlag(constants.FlagJSON) == true {
		printUsageJSON(resolveFilterQuery())
		return
	}

	q := resolveFilterQuery()
	needsFilter := len(q) > 0 || hasFlag(constants.FlagFilter) == true || hasFlag(constants.FlagFilterShort) == true
	if needsFilter == true {
		printUsageFiltered(q)

		return
	}

	printUsage()
}

// dispatchHelpTopic renders help for a named topic, falling back to filtered
// usage when the topic has no dedicated help text.
func dispatchHelpTopic(rawTopic string) {
	topic := normalizeHelpTopic(rawTopic)
	_, err := helptext.ReadRaw(topic)
	if err != nil {
		printUsageFiltered(topic)
		return
	}

	_, mode := ParsePrettyFlag(os.Args[3:])
	helptext.PrintWithMode(topic, mode)
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
	}
	return topic
}
