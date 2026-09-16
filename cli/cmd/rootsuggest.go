package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var primaryTopCommands = []string{
	"scan", "clone", "create", "clone-sync", "pull", "push", "pull-all",
	"status", "reconcile", "fix", "stash", "wip", "discard", "exec",
	"release", "pull-release", "changelog", "cd", "group", "open",
	"history", "stats", "export", "import", "profile", "bookmark",
	"dashboard", "version", "help", "diff", "amend", "sync",
	"add", "rm", "mv", "prune", "revert", "ip", "zsh", "user",
	"service", "schedule", "macro", "os", "storage", "pipeline",
	"servers-clients", "servers-client", "sc", "clients", "cluster",
	"agy", "fix-pipeline", "tasks", "task", "author", "sponsor", "credits",
}

func buildUnknownCommandMessage(command string, suggestions []string) string {
	if looksLikeURLToken(command) {
		return fmt.Sprintf(constants.ErrUnknownCommandURLHint, command)
	}

	msg := fmt.Sprintf(constants.ErrUnknownCommand, command)
	if len(suggestions) > 0 {
		msg += fmt.Sprintf("\n  Did you mean: %s?", strings.Join(suggestions, ", "))
	}

	return msg
}

func printCommandSuggestions(suggestions []string) {
	if len(suggestions) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("  Did you mean one of these?")
	for _, s := range suggestions {
		fmt.Printf("    gitmap %s\n", s)
	}

	fmt.Println()
}

func handleUnknownCommand(command string) {
	suggestions := suggestTopLevelCommands(command)
	printCommandSuggestions(suggestions)
	msg := buildUnknownCommandMessage(command, suggestions)
	printUsage()
	dispatchErr := apperror.NewWithDetails(
		"cmd.dispatch",
		"E1001",
		msg,
		"cmd.root",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		map[string]any{"command": command},
	)
	cliexit.HandleError(dispatchErr, 1)
}
