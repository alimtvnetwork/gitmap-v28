package cmdagm

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// DispatchAgm routes `agm queue <subcommand>` to the Group E AGM commands.
func DispatchAgm(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("agm: missing subcommand (queue audit)")
	}
	if strings.ToLower(args[0]) == "queue" && len(args) > 1 && strings.ToLower(args[1]) == "audit" {
		return RunAgmQueueAudit(args[2:])
	}
	if args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		fmt.Print(agmQueueAuditHelp)
		return nil
	}

	return appfault.NewValidationError(fmt.Sprintf("agm: unknown subcommand %q (queue audit)", strings.Join(args, " ")))
}
