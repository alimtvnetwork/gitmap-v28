package cmdpe

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// DispatchPe routes `pe <subcommand>` to the Group E pipeline commands.
func DispatchPe(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("pe: missing subcommand (triage | rerun)")
	}
	switch strings.ToLower(args[0]) {
	case "triage":
		return RunPeTriage(args[1:])
	case "rerun":
		return RunPeRerun(args[1:])
	case "--help", "-h", "help":
		fmt.Print(peTriageHelp)
		fmt.Print(peRerunHelp)
		return nil
	default:
		return appfault.NewValidationError(fmt.Sprintf("pe: unknown subcommand %q (triage | rerun)", args[0]))
	}
}
