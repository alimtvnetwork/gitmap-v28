package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var primaryTopCommands = []string{
	"scan", "clone", "clone-only-missing", "com", "create", "clone-sync", "pull", "push", "pull-all",
	"pull-all-efficient", "pae", "pull-ae", "pull-all-efficient-table", "paet",
	"status", "reconcile", "fix", "stash", "wip", "discard", "exec",
	"release", "pull-release", "changelog", "cd", "group", "open",
	"history", "stats", "export", "import", "profile", "bookmark",
	"dashboard", "version", "help", "diff", "amend", "sync",
	"add", "rm", "mv", "prune", "revert", "ip", "zsh", "user",
	"service", "schedule", "macro", "os", "storage", "pipeline", "pe", "pd", "ee",
	"servers-clients", "servers-client", "sc", "clients", "cluster",
	"agy", "fix-pipeline", "tasks", "task", "author", "sponsor", "credits",
	"update", "ua", "ssh", "ssh-join", "ssh-exec", "ssh-nodes", "install-exec", "deploy",
	"clean", "clear", "terminal", "clear-terminal", "clean-terminal", "devtool",
	"deploy-all-keys", "deploy-keys-all", "deploy-keys", "which-format", "import-all-json",
	"what-configs", "wc", "merge-json", "failed-commands", "fc", "unknown-commands", "errors",
	"gitignore", "agm", "ta", "nodes", "ping",
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
	fmt.Println("  💡 Did you mean one of these?")
	for _, s := range suggestions {
		fmt.Printf("    gitmap %s\n", s)
	}

	fmt.Println()
}

func handleUnknownCommand(command string) {
	suggestions := suggestTopLevelCommands(command)
	printCommandSuggestions(suggestions)
	msg := buildUnknownCommandMessage(command, suggestions)
	store.LogFailedCommand(command, strings.Join(os.Args[1:], " "), "root", "E1001", msg, suggestions)
	fmt.Println("  Run 'gitmap help' or 'gitmap <command> --help' to view available commands.")
	fmt.Println("  Run 'gitmap failed-commands' (or 'gitmap fc') to inspect failed command logs & suggestions.")
	fmt.Println()
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
