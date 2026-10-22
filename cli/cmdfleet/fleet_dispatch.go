package cmdfleet

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// fleetHandlers maps a lowercased `nodes` subcommand to its handler (prompt 38 map dispatch).
var fleetHandlers = map[string]func([]string) *appfault.AppError{
	"jump":    RunNodesJump,
	"version": RunNodesVersion,
	"--help":  printFleetDispatchHelp,
	"-h":      printFleetDispatchHelp,
	"help":    printFleetDispatchHelp,
}

func printFleetDispatchHelp([]string) *appfault.AppError {
	fmt.Print(nodesJumpHelp)
	fmt.Print(nodesVersionHelp)
	return nil
}

// DispatchFleet routes `nodes <subcommand>` to the Group E fleet commands.
func DispatchFleet(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("nodes: missing subcommand (jump | version)")
	}
	if handler, ok := fleetHandlers[strings.ToLower(args[0])]; ok {
		return handler(args[1:])
	}
	return appfault.NewValidationError(fmt.Sprintf("nodes: unknown subcommand %q (jump | version)", args[0]))
}
