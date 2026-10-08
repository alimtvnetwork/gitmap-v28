package cmdos

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runPower handles the `gitmap power` command family.
func RunPower(args []string) error {
	checkHelp("power", args)

	if len(args) == 0 {
		return runPowerStatus()
	}

	subcmd := args[0]
	subArgs := args[1:]

	return dispatchPowerSubcommand(subcmd, subArgs)
}

func dispatchPowerSubcommand(subcmd string, subArgs []string) error {
	switch subcmd {
	case constants.SubCmdPowerStatus, "st":
		return runPowerStatus()
	case constants.SubCmdPowerNever, "never", "ns":
		return RunPowerNeverSleep()
	case constants.SubCmdPowerSet:
		return RunPowerSet(subArgs)
	case constants.SubCmdPowerReset, "restore":
		return RunPowerReset()
	case constants.SubCmdPowerHistory, "hist", "log", "logs":
		return runPowerHistory(subArgs)
	default:
		return apperror.NewSimple(fmt.Sprintf("unknown power subcommand %q (see 'gitmap power --help')", subcmd), "E_INVALID_SUBCMD")
	}
}
