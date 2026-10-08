package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// primaryTopCommands contains the fast core primary tier. The full repository
// command space is dynamically pooled from completion.AllCommands() in
// collectTopCommandCandidates() for exhaustive typo suggestion matching.
var primaryTopCommands = []string{
	"scan", "clone", "clone-only-missing", "com", "create", "clone-sync", "pull", "push", "pull-all",
	"pull-all-efficient", "pae", "pull-ae", "pull-all-efficient-table", "paet",
	"status", "reconcile", "fix", "stash", "wip", "discard", "exec",
	"release", "pull-release", "changelog", "cd", "group", "open",
	"history", "stats", "export", "import", "profile", "bookmark",
	"dashboard", "version", "help", "diff", "amend", "sync",
	"commit", "cpf", "cpb", "cpr", "search", "find", "apps", "install", "uninstall", "login", "setup",
	"add", "rm", "mv", "prune", "revert", "ip", "zsh", "user",
	"service", "schedule", "macro", "os", "storage", "pipeline", "pe", "pd", "ee",
	"servers-clients", "servers-client", "sc", "clients", "cluster",
	"agy", "fix-pipeline", "tasks", "task", "author", "sponsor", "credits",
	"update", "ua", "ssh", "ssh-join", "ssh-exec", "ssh-nodes", "install-exec", "deploy",
	"clean", "clear", "terminal", "clear-terminal", "clean-terminal", "devtool",
	"deploy-all-keys", "deploy-keys-all", "deploy-keys", "which-format", "import-all-json",
	"what-configs", "wc", "merge-json", "failed-commands", "fc", "unknown-commands", "errors",
	"gitignore", "agm", "ta", "nodes", "ping", "ports", "port",
}

func buildUnknownCommandMessage(command string, suggestions []string) string {
	if looksLikeURLToken(command) {
		return fmt.Sprintf(constants.ErrUnknownCommandURLHint, command)
	}

	msg := fmt.Sprintf(constants.ErrUnknownCommand, command)
	if len(suggestions) > 0 {
		msg += fmt.Sprintf("\n  It is not there, but here is a suggestion you can try: %s", strings.Join(suggestions, ", "))
	}

	return msg
}

func handleUnknownCommand(command string) {
	InterceptUnknownCommand(command)
}
