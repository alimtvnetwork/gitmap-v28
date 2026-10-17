package cmdci

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// DispatchCi routes `ci <subcommand>` to the Group E CI commands.
func DispatchCi(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("ci: missing subcommand (open)")
	}
	switch strings.ToLower(args[0]) {
	case "open":
		return RunCiOpen(args[1:])
	case "--help", "-h", "help":
		fmt.Print(ciOpenHelp)
		return nil
	default:
		return appfault.NewValidationError(fmt.Sprintf("ci: unknown subcommand %q (open)", args[0]))
	}
}
