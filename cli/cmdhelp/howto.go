// Package cmdhelp — howto.go: `gitmap howto <task...>` keyword task→recipe matching.
//
// Pure keyword matching, no AI calls. An unknown task prints the topic index
// so the operator can continue with `help topic <name>`.
package cmdhelp

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// howtoRecipe maps a task keyword cluster to an ordered command recipe.
type howtoRecipe struct {
	keywords []string
	title    string
	topic    string
	steps    []string
}

// howtoRecipes is the task→recipe table (proposal P-08).
var howtoRecipes = []howtoRecipe{
	{
		keywords: []string{"ship", "release", "version", "publish", "deploy"},
		title:    "ship a new version",
		topic:    "releases",
		steps: []string{
			"gitmap summary",
			"gitmap release-dry",
			"gitmap release-notes",
			"gitmap release vX.Y.Z",
			"gitmap pe",
		},
	},
	{
		keywords: []string{"pipeline", "red", "ci", "fail", "flaky", "rerun", "broken build"},
		title:    "fix a red pipeline",
		topic:    "pipeline",
		steps: []string{
			"gitmap pe",
			"gitmap pe triage",
			"gitmap pe rerun <run-id>",
		},
	},
	{
		keywords: []string{"theme", "color", "color", "palette"},
		title:    "change the terminal theme",
		topic:    "themes",
		steps: []string{
			"gitmap theme list",
			"gitmap theme preview <name>",
			"gitmap --theme <name> <command>",
		},
	},
	{
		keywords: []string{"install", "setup", "verify install"},
		title:    "verify the gitmap installation",
		topic:    "",
		steps: []string{
			"gitmap install verify",
			"gitmap install urls <tag>",
		},
	},
	{
		keywords: []string{"scan", "index", "orphan", "discover"},
		title:    "keep the repo index truthful",
		topic:    "scanner",
		steps: []string{
			"gitmap scan",
			"gitmap scan verify",
			"gitmap scan diff",
		},
	},
	{
		keywords: []string{"fleet", "node", "ssh", "remote", "machine"},
		title:    "work across fleet machines",
		topic:    "fleet",
		steps: []string{
			"gitmap nodes ping",
			"gitmap nodes version",
			"gitmap nodes jump <node>",
		},
	},
	{
		keywords: []string{"migrate", "migration", "move", "export", "import", "transfer"},
		title:    "migrate a repo set to a new machine",
		topic:    "migration",
		steps: []string{
			"gitmap scan export",
			"gitmap repo import <file> --into <dir>",
			"gitmap scan",
		},
	},
}

// RunHowto executes `gitmap howto <task...>` — pure keyword matching.
func RunHowto(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(howtoHelpText)
		cliexit.Exit(0)
	}

	task := strings.ToLower(strings.TrimSpace(strings.Join(args, " ")))
	if task == "" {
		return appfault.NewValidationError("missing task description\n" + howtoUsageLine +
			"example: gitmap howto \"ship a new version\"")
	}

	for _, recipe := range howtoRecipes {
		if hasHowtoKeyword(task, recipe.keywords) {
			printHowtoRecipe(recipe)

			return nil
		}
	}

	printHowtoUnknown(task)

	return nil
}

// hasHowtoKeyword reports whether the task text contains any recipe keyword.
func hasHowtoKeyword(task string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(task, keyword) {
			return true
		}
	}

	return false
}

// printHowtoRecipe renders the matched recipe as a numbered command list.
func printHowtoRecipe(recipe howtoRecipe) {
	fmt.Printf("RECIPE — %s:\n", recipe.title)
	for i, step := range recipe.steps {
		fmt.Printf("  %d) %s\n", i+1, step)
	}
	if recipe.topic != "" {
		fmt.Printf("see: gitmap help topic %s\n", recipe.topic)
	}
}

// printHowtoUnknown reports no match and points at the topic index.
func printHowtoUnknown(task string) {
	fmt.Printf("no recipe matches %q.\n\n", task)
	printHelpTopicIndex()
}
