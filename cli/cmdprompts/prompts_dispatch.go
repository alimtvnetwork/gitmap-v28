package cmdprompts

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// promptsHandlers maps a lowercased `prompts` subcommand to its handler (prompt 38 map dispatch).
var promptsHandlers = map[string]func([]string) *appfault.AppError{
	"drift":  RunPromptsDrift,
	"--help": printPromptsDispatchHelp,
	"-h":     printPromptsDispatchHelp,
	"help":   printPromptsDispatchHelp,
}

func printPromptsDispatchHelp([]string) *appfault.AppError {
	fmt.Print(promptsDriftHelp)
	return nil
}

// DispatchPrompts routes `prompts <subcommand>` to the Group E prompt commands.
func DispatchPrompts(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("prompts: missing subcommand (drift)")
	}
	if handler, ok := promptsHandlers[strings.ToLower(args[0])]; ok {
		return handler(args[1:])
	}
	return appfault.NewValidationError(fmt.Sprintf("prompts: unknown subcommand %q (drift)", args[0]))
}
