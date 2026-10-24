package cmdhooks

import (
	"fmt"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks/enums"
)

// hooksHandlers maps a subcommand action enum to its handler (prompt 38 map dispatch).
var hooksHandlers = map[enums.HookActionType]func([]string) *appfault.AppError{
	enums.HookActionInstall:   RunHooksInstall,
	enums.HookActionHelpLong:  printHooksDispatchHelp,
	enums.HookActionHelpShort: printHooksDispatchHelp,
	enums.HookActionHelp:      printHooksDispatchHelp,
}

func printHooksDispatchHelp([]string) *appfault.AppError {
	fmt.Print(constants.HelpTextHooksInstall)

	return nil
}

// DispatchHooks routes `hooks <subcommand>` to the Group E hook commands.
func DispatchHooks(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError(constants.ErrMissingSubcommand)
	}

	action, parseErr := enums.ParseHookActionType(args[0])
	if parseErr != nil {
		return appfault.NewValidationError(fmt.Sprintf(constants.ErrUnknownSubcommandFmt, args[0]))
	}

	handler, hasHandler := hooksHandlers[action]
	if !hasHandler {
		return appfault.NewValidationError(fmt.Sprintf(constants.ErrUnknownSubcommandFmt, args[0]))
	}

	return handler(args[1:])
}
