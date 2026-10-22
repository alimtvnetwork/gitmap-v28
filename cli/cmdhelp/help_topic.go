// Package cmdhelp — help_topic.go: `gitmap help topic <name>` / `gitmap ht`.
//
// Topic index (proposal P-07): a learning path per workflow — the relevant
// commands in the order an operator runs them. Read-only.
package cmdhelp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// helpTopicStep is one ordered step inside a topic: the command plus a
// one-line note on when it runs.
type helpTopicStep struct {
	command string
	note    string
}

// helpTopic is a learning-path index: the commands for one workflow, in order.
type helpTopic struct {
	name        string
	description string
	steps       []helpTopicStep
}

// helpTopics is the topic index (proposal P-07).
var helpTopics = []helpTopic{
	{
		name:        "releases",
		description: "Ship a version: check state, dry-run, notes, release, watch CI.",
		steps: []helpTopicStep{
			{command: "gitmap summary", note: "repo state before you ship"},
			{command: "gitmap release-dry", note: "preview the release with no side effects"},
			{command: "gitmap release-notes", note: "changelog between tags"},
			{command: "gitmap release vX.Y.Z", note: "cut the release"},
			{command: "gitmap pe", note: "watch the release pipeline"},
		},
	},
	{
		name:        "pipeline",
		description: "Triage and recover a red CI pipeline.",
		steps: []helpTopicStep{
			{command: "gitmap pe", note: "list pipeline errors"},
			{command: "gitmap pe triage", note: "ranked failure bundle, worst first"},
			{command: "gitmap pe rerun <run-id>", note: "re-run failed jobs without the browser"},
			{command: "gitmap pipeline-ai status", note: "wait on pipeline state"},
		},
	},
	{
		name:        "fleet",
		description: "Reach and verify every machine in the fleet.",
		steps: []helpTopicStep{
			{command: "gitmap nodes ping", note: "which nodes are reachable"},
			{command: "gitmap nodes version", note: "version skew across the fleet"},
			{command: "gitmap nodes jump <node>", note: "one-shot shell on a node"},
			{command: "gitmap ssh <node>", note: "raw SSH session"},
		},
	},
	{
		name:        "migration",
		description: "Move a repo set between machines.",
		steps: []helpTopicStep{
			{command: "gitmap scan export", note: "dump the repo index to portable JSON"},
			{command: "gitmap repo import <file> --into <dir>", note: "clone and register on the new machine"},
			{command: "gitmap scan", note: "re-index the new workspace"},
		},
	},
	{
		name:        "scanner",
		description: "Keep the repo index truthful.",
		steps: []helpTopicStep{
			{command: "gitmap scan", note: "discover repos on disk"},
			{command: "gitmap scan verify", note: "index vs filesystem consistency"},
			{command: "gitmap scan diff", note: "what changed since the last scan"},
			{command: "gitmap orphans", note: "on-disk repos missing from the index"},
		},
	},
	{
		name:        "themes",
		description: "Terminal palettes: list, preview, activate.",
		steps: []helpTopicStep{
			{command: "gitmap theme list", note: "available palettes and the active one"},
			{command: "gitmap theme preview <name>", note: "try a palette without switching"},
			{command: "gitmap --theme <name> <command>", note: "run one command in a palette"},
		},
	},
	{
		name:        "conflicts",
		description: "Inspect and resolve merge conflicts using IDE 3-way editor or CLI.",
		steps: []helpTopicStep{
			{command: "gitmap conflicts", note: "inspect conflicted files in repo"},
			{command: "gitmap resolve --ide", note: "launch VS Code / Cursor 3-Way Merge Editor"},
			{command: "gitmap resolve --theirs", note: "accept remote incoming changes"},
			{command: "gitmap resolve --ours", note: "accept local current changes"},
			{command: "gitmap resolve --abort", note: "safely abort in-progress merge or rebase"},
			{command: "gitmap cpf \"resolve conflicts\"", note: "stage, commit, and push once resolved"},
		},
	},
	{
		name:        "storage",
		description: "Manage large files, Git LFS migration, and repo history flattening.",
		steps: []helpTopicStep{
			{command: "gitmap split large files --dry-run", note: "preview large files to chunk into ZIP parts"},
			{command: "gitmap lfs convert --dry-run", note: "preview binary files to migrate to Git LFS"},
			{command: "gitmap repo flatten --dry-run", note: "preview squashing bloated repo history"},
			{command: "gitmap repo undo-flatten <repo>", note: "restore flattened repo from backup archive"},
		},
	},
}

// RunHelpTopic executes `gitmap help topic <name>` / `gitmap ht <name>`.
func RunHelpTopic(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(helpTopicHelpText)

		return nil
	}

	if len(args) == 0 {
		printHelpTopicIndex()

		return nil
	}

	name := strings.ToLower(strings.TrimSpace(args[0]))
	for _, topic := range helpTopics {
		if topic.name == name {
			printHelpTopic(topic)

			return nil
		}
	}

	return appfault.NewValidationError(
		"unknown topic " + strconv.Quote(name) + " — valid topics: " + helpTopicNames())
}

// helpTopicNames returns the comma-separated topic names for error messages.
func helpTopicNames() string {
	names := make([]string, 0, len(helpTopics))
	for _, topic := range helpTopics {
		names = append(names, topic.name)
	}

	return strings.Join(names, ", ")
}

// printHelpTopicIndex lists every topic with its one-line description.
func printHelpTopicIndex() {
	fmt.Println("TOPICS:")
	for _, topic := range helpTopics {
		fmt.Printf("  %-10s %s\n", topic.name, topic.description)
	}
	fmt.Println("run: gitmap help topic <name>   (alias: gitmap ht <name>)")
}

// printHelpTopic renders one topic's ordered command list.
func printHelpTopic(topic helpTopic) {
	fmt.Printf("TOPIC %s: %s\n", topic.name, topic.description)
	for i, step := range topic.steps {
		fmt.Printf("  %d) %s — %s\n", i+1, step.command, step.note)
	}
}
