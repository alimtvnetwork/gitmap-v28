package cmd

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhelp"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
)

// runHelpDispatch handles the `help` subcommand including topic
// help, `help example <cmd>`, --groups, --compact, and the default
// usage screen.
func runHelpDispatch() error {
	isHelpExample := len(os.Args) >= 3 && strings.EqualFold(os.Args[2], "example")
	if isHelpExample {
		return dispatchHelpExample(os.Args[3:])
	}

	isHelpTopicCmd := len(os.Args) >= 3 && strings.EqualFold(os.Args[2], "topic")
	if isHelpTopicCmd {
		return cmdhelp.RunHelpTopic(os.Args[3:])
	}

	hasTopic := len(os.Args) >= 3 && !isFlagToken(os.Args[2])
	if hasTopic {
		return dispatchHelpTopic(os.Args[2])
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

func renderSpecificHelpTopic(topic string, mode render.PrettyModeType) {
	prettyOff := mode == render.PrettyOff
	helpdoc.PrintWithMode(topic, mode)
	helpdoc.PrintExampleSection(topic, prettyOff)
	printUsageFooterShort()
}

// dispatchHelpTopic renders help for a named topic, falling back to filtered
// usage when the topic has no dedicated help text.
func dispatchHelpTopic(rawTopic string) error {
	topic := normalizeHelpTopic(rawTopic)
	if isHelpCategoryGroup(topic) {
		printUsageCategoryGroup(topic)

		return nil
	}

	_, mode := ParsePrettyFlag(os.Args[3:])
	prettyOff := mode == render.PrettyOff
	if tryRenderRichTopic(topic) {
		helpdoc.PrintExampleSection(topic, prettyOff)

		return nil
	}

	if _, err := helpdoc.ReadRaw(topic); err != nil {
		printUsageFiltered(topic)

		return nil
	}

	renderSpecificHelpTopic(topic, mode)

	return nil
}

// dispatchHelpExample handles `gitmap help example [cmd]`.
func dispatchHelpExample(args []string) error {
	cleaned, mode := ParsePrettyFlag(args)
	prettyOff := mode == render.PrettyOff

	if len(cleaned) == 0 {
		helpdoc.PrintHelpExampleIndex(prettyOff)

		return nil
	}

	helpdoc.PrintHelpExample(cleaned[0], prettyOff)

	return nil
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
	case "res":
		return "resolve"
	case "cf", "conflict":
		return "conflicts"
	case "rc", "rec", "clone-now", "cnow", "relclone":
		return "reclone"
	case "split", "split-large", "split-files":
		return "split-large-files"
	case "lfs", "lfsconvert", "git-lfs":
		return "lfs-convert"
	case "repo-flatten":
		return "flatten"
	case "cin", "commitin":
		return "commit-in"
	case "mergein":
		return "merge-in"
	case "mrg":
		return "merge-right"
	case "ml":
		return "merge-left"
	case "mb":
		return "merge-both"
	case "th", "thelp":
		return "help"
	}

	return topic
}
