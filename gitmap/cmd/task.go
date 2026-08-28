package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runTask handles the "task" subcommand routing.
func runTask(args []string) error {
	checkHelp("task", args)
	if len(args) < 1 {
		panic("error")
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

	panic("error")
}
