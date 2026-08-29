package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runTask handles the "task" subcommand routing.
func runTask(args []string) error {
	checkHelp("task", args)
	if len(args) < 1 {
		err := apperror.NewWithDetails(
			"cmd.task.run",
			"E1090",
			"missing required task subcommand; usage: gitmap task [create|list|run|show|delete]",
			"cmd.task",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return nil
	}

	sub := args[0]
	rest := args[1:]

	routeTaskSub(sub, rest)
	return nil
}

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
