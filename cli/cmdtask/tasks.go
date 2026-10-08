// Package cmd — tasks.go manages pending and completed task execution queues.
package cmdtask

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagent"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runTasks routes the "tasks" and "task" command suite.
func RunTasks(args []string) error {
	checkHelp("task", args)
	if len(args) == 0 {
		return runTasksList()
	}

	sub := strings.ToLower(args[0])
	rest := args[1:]

	if isAgentTaskSubcommand(sub) {
		return cmdagent.DispatchTaskSQLite(args)
	}

	if isLegacyTaskSubcommand(sub) {
		routeTaskSub(sub, rest)

		return nil
	}

	return dispatchTaskSubcommand(sub, rest)
}

func isAgentTaskSubcommand(sub string) bool {
	switch sub {
	case "init", "add", "add-subtasks", "claim", "complete", "fail", "log-action", "status", "schema", "diagnose":
		return true
	default:
		return false
	}
}

func isLegacyTaskSubcommand(sub string) bool {
	return sub == constants.CmdTaskCreate || sub == constants.CmdTaskRun ||
		sub == constants.CmdTaskShow || sub == constants.CmdTaskDelete
}

func dispatchTaskSubcommand(sub string, args []string) error {
	switch sub {
	case constants.CmdTaskList, "ls":
		return runTasksList()
	case constants.CmdTaskHistory, "hist", "hi":
		return runTasksHistory(args)
	case constants.CmdTaskUndo:
		return runTasksUndo(args)
	case constants.CmdTaskRedo:
		return runTasksRedo(args)
	case constants.CmdTaskClear:
		runPendingClear(args)

		return nil
	default:
		return runTasksList()
	}
}
