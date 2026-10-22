package cmdhooks

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// hooksHandlers maps a lowercased `hooks` subcommand to its handler (prompt 38 map dispatch).
var hooksHandlers = map[string]func([]string) *appfault.AppError{
	"install": RunHooksInstall,
	"--help":  printHooksDispatchHelp,
	"-h":      printHooksDispatchHelp,
	"help":    printHooksDispatchHelp,
}

func printHooksDispatchHelp([]string) *appfault.AppError {
	fmt.Print(hooksInstallHelp)
	return nil
}

// DispatchHooks routes `hooks <subcommand>` to the Group E hook commands.
func DispatchHooks(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("hooks: missing subcommand (install)")
	}
	if handler, ok := hooksHandlers[strings.ToLower(args[0])]; ok {
		return handler(args[1:])
	}
	return appfault.NewValidationError(fmt.Sprintf("hooks: unknown subcommand %q (install)", args[0]))
}
