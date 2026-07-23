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
	if len(os.Args) >= 3 && !isFlagToken(os.Args[2]) {
		topic := os.Args[2]
		switch topic {
		case constants.CmdRmAlias, constants.CmdRmAlias2:
			topic = constants.CmdRm
		case constants.CmdStaleAlias:
			topic = constants.CmdStale
		case constants.CmdRecentAlias:
			topic = constants.CmdRecent
		case constants.CmdPRAlias:
			topic = constants.CmdPR
		}
		if _, err := helptext.ReadRaw(topic); err == nil {
			_, mode := ParsePrettyFlag(os.Args[3:])
			helptext.PrintWithMode(topic, mode)

			return
		}
		// Unknown topic — treat as a filter query so users can type
		// `gitmap help ssh` and get the same hits as `gitmap help -f ssh`.
		printUsageFiltered(topic)

		return
	}
	if hasFlag(constants.FlagJSON) {
		printUsageJSON(resolveFilterQuery())

		return
	}
	if q := resolveFilterQuery(); len(q) > 0 || hasFlag(constants.FlagFilter) || hasFlag(constants.FlagFilterShort) {
		printUsageFiltered(q)

		return
	}
	if hasFlag(constants.FlagGroups) {
		printHelpGroups()

		return
	}
	if hasFlag(constants.FlagCompact) {
		printUsageCompact()

		return
	}
	printUsage()
}
