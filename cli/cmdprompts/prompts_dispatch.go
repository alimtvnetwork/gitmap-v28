package cmdprompts

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// DispatchPrompts routes `prompts <subcommand>` to the Group E prompt commands.
func DispatchPrompts(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("prompts: missing subcommand (drift)")
	}
	switch strings.ToLower(args[0]) {
	case "drift":
		return RunPromptsDrift(args[1:])
	case "--help", "-h", "help":
		fmt.Print(promptsDriftHelp)
		return nil
	default:
		return appfault.NewValidationError(fmt.Sprintf("prompts: unknown subcommand %q (drift)", args[0]))
	}
}
