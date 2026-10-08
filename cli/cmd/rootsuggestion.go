package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/suggestion"
)

// InterceptUnknownCommand resolves suggestions via suggestion.DefaultEngine,
// prints the Catppuccin suggestion box, attaches suggestions to AppError, and terminates.
func InterceptUnknownCommand(command string) {
	group := suggestion.DefaultEngine().ResolveCommand(command)
	if group.HasSuggestions() {
		_ = suggestion.RenderBox(os.Stderr, group)
	}
	logUnknownCommandTelemetry(command, group)
	appErr := buildUnknownCommandAppError(command, group)
	cliexit.HandleError(appErr, 1)
}

// ResolveCommandSuggestions resolves suggestions for any given command token.
func ResolveCommandSuggestions(token string) suggestion.SuggestionGroup {
	return suggestion.DefaultEngine().ResolveCommand(token)
}

// RenderErrorSuggestions renders attached suggestions from an AppError if present.
func RenderErrorSuggestions(w io.Writer, appErr *apperror.AppError) bool {
	if appErr == nil || !appErr.HasSuggestions() {
		return false
	}
	group := suggestion.SuggestionGroup{
		Title:       "Remediation Suggestions",
		Reason:      appErr.Message,
		Suggestions: appErr.Suggestions,
	}
	_ = suggestion.RenderBox(w, group)
	return true
}

func logUnknownCommandTelemetry(command string, group suggestion.SuggestionGroup) {
	cmdList := extractSuggestionCommandList(group.Suggestions)
	rawArgs := strings.Join(os.Args[1:], " ")
	store.LogFailedCommand(command, rawArgs, "root", "E1001", group.Reason, cmdList)
}

func extractSuggestionCommandList(items []suggestion.Suggestion) []string {
	out := make([]string, 0, len(items))
	for _, s := range items {
		out = append(out, s.Command)
	}
	return out
}

func buildUnknownCommandAppError(command string, group suggestion.SuggestionGroup) *apperror.AppError {
	appErr := apperror.NewWithDetails(
		"cmd.dispatch",
		"E1001",
		fmt.Sprintf("unrecognized command: %s", command),
		"cmd.root",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		map[string]any{"command": command},
	)
	return appErr.WithSuggestionGroup(group)
}
