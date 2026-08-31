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
func dispatchUtility(command string) (bool, error) {
	return runDispatchTable(command, utilityDispatchEntries())
}

// utilityDispatchEntries returns the routing table for utility commands.
func utilityDispatchEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"error"}, func() error { return runErrorCmd(argsTail()) }},
		{[]string{constants.CmdUpdate}, func() error { checkHelp("update", argsTail()); return runUpdate() }},
		{[]string{constants.CmdUpdateRunner}, runUpdateRunner},
		{[]string{constants.CmdUpdateCleanup}, runUpdateCleanup},
		{
			[]string{constants.CmdInstalledDir, constants.CmdInstalledDirAlias},
			func() error { checkHelp("installed-dir", argsTail()); return runInstalledDir() },
		},
		{[]string{constants.CmdRevert}, func() error { return runRevert(argsTail()) }},
		{[]string{constants.CmdRm, constants.CmdRmAlias, constants.CmdRmAlias2}, func() error { return runRm(argsTail()) }},
		{[]string{constants.CmdRevertRunner}, runRevertRunner},
		{
			[]string{constants.CmdVersion, constants.CmdVersionAlias},
			func() error {
				checkHelp("version", argsTail())
				fmt.Printf(constants.MsgVersionFmt, constants.Version)
				return nil
			},
		},
		{[]string{constants.CmdHelp, "--help", "-h"}, runHelpDispatch},
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
		{[]string{constants.CmdWhoAmI, constants.CmdWhoAmIAlias}, func() error { checkHelp("whoami", argsTail()); return runWhoAmI(argsTail()) }},
		{[]string{constants.CmdSSHBind, constants.CmdSSHBindAlias}, func() error { checkHelp("ssh-bind", argsTail()); return runSSHBind(argsTail()) }},
		{[]string{constants.CmdFixAuth, constants.CmdFixAuthAlias}, func() error { checkHelp("fix-auth", argsTail()); return runFixAuth(argsTail()) }},
		{[]string{"pipeline-ai", "pl-ai", "plai", "pipeline_ai"}, func() error { return runPipelineAI(argsTail()) }},
		{[]string{"pipeline", "pl"}, func() error { return runPipeline(argsTail()) }},
		{[]string{"error-logs", "error-log", "errorlogs", "errorlog", "errors", "err"}, func() error { return runPipeline(append([]string{"error-logs"}, argsTail()...)) }},
		{[]string{"logs", "log"}, func() error { return runPipeline(append([]string{"logs"}, argsTail()...)) }},
		{[]string{"waittime", "wait-time", "eta"}, func() error { return runPipeline(append([]string{"waittime"}, argsTail()...)) }},
		{[]string{"ui"}, func() error { return runUI(argsTail()) }},
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
