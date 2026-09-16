package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// routeTaskSub routes to the appropriate task subcommand.
func routeTaskSub(sub string, args []string) {
	if sub == constants.CmdTaskCreate {
		runTaskCreate(args)

		return
	}

	if sub == constants.CmdTaskList {
		runTaskList()

		return
	}

	if sub == constants.CmdTaskRun {
		runTaskRun(args)

		return
	}

	if sub == constants.CmdTaskShow {
		runTaskShow(args)

		return
	}

	if sub == constants.CmdTaskDelete {
		runTaskDelete(args)

		return
	}

	err := apperror.NewWithDetails(
		"cmd.task.route",
		"E1091",
		fmt.Sprintf("unknown task subcommand '%s'", sub),
		"cmd.task",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		map[string]any{"subcommand": sub},
	)
	cliexit.HandleError(err, 1)
}
