package cmdfleet

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// DispatchFleet routes `nodes <subcommand>` to the Group E fleet commands.
func DispatchFleet(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("nodes: missing subcommand (jump | version)")
	}
	switch strings.ToLower(args[0]) {
	case "jump":
		return RunNodesJump(args[1:])
	case "version":
		return RunNodesVersion(args[1:])
	case "--help", "-h", "help":
		fmt.Print(nodesJumpHelp)
		fmt.Print(nodesVersionHelp)
		return nil
	default:
		return appfault.NewValidationError(fmt.Sprintf("nodes: unknown subcommand %q (jump | version)", args[0]))
	}
}
