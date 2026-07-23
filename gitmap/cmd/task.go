package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runTask handles the "task" subcommand routing.
func runTask(args []string) {
	checkHelp("task", args)
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, constants.ErrTaskSubcommand, "")
		os.Exit(1)
	}

	sub := args[0]
	rest := args[1:]

	routeTaskSub(sub, rest)
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

	fmt.Fprintf(os.Stderr, constants.ErrTaskSubcommand, sub)
	os.Exit(1)
}
