package cmdpe

import (
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// peHandlers maps a lowercased `pe` subcommand to its handler (prompt 38 map dispatch).
var peHandlers = map[string]func([]string) *appfault.AppError{
	"triage": RunPeTriage,
	"rerun":  RunPeRerun,
	"--help": printPeDispatchHelp,
	"-h":     printPeDispatchHelp,
	"help":   printPeDispatchHelp,
}

func printPeDispatchHelp([]string) *appfault.AppError {
	fmt.Print(peTriageHelp)
	fmt.Print(peRerunHelp)
	return nil
}

// DispatchPe routes `pe <subcommand>` to the Group E pipeline commands.
func DispatchPe(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("pe: missing subcommand (triage | rerun)")
	}
	if handler, ok := peHandlers[strings.ToLower(args[0])]; ok {
		return handler(args[1:])
	}
	return appfault.NewValidationError(fmt.Sprintf("pe: unknown subcommand %q (triage | rerun)", args[0]))
}
