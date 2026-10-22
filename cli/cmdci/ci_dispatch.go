package cmdci

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// ciHandlers maps a lowercased `ci` subcommand to its handler (prompt 38 map dispatch).
var ciHandlers = map[string]func([]string) *appfault.AppError{
	"open":   RunCiOpen,
	"--help": printCiDispatchHelp,
	"-h":     printCiDispatchHelp,
	"help":   printCiDispatchHelp,
}

func printCiDispatchHelp([]string) *appfault.AppError {
	fmt.Print(ciOpenHelp)
	return nil
}

// DispatchCi routes `ci <subcommand>` to the Group E CI commands.
func DispatchCi(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("ci: missing subcommand (open)")
	}
	if handler, ok := ciHandlers[strings.ToLower(args[0])]; ok {
		return handler(args[1:])
	}
	return appfault.NewValidationError(fmt.Sprintf("ci: unknown subcommand %q (open)", args[0]))
}
