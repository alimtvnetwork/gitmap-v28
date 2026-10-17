package cmdhooks

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// DispatchHooks routes `hooks <subcommand>` to the Group E hook commands.
func DispatchHooks(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("hooks: missing subcommand (install)")
	}
	switch strings.ToLower(args[0]) {
	case "install":
		return RunHooksInstall(args[1:])
	case "--help", "-h", "help":
		fmt.Print(hooksInstallHelp)
		return nil
	default:
		return appfault.NewValidationError(fmt.Sprintf("hooks: unknown subcommand %q (install)", args[0]))
	}
}
