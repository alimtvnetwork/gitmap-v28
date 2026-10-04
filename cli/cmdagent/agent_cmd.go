package cmdagent

import (
	"strings"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	// AgentCmd is the root Cobra command for AI agent task orchestration and telemetry.
	AgentCmd = &cobra.Command{
		Use:     "agent",
		Aliases: []string{"ai-agents", "agents"},
		Short:   "AI agent task orchestrator, multi-tier split-db management, and telemetry",
		RunE:    runDefaultAgentCmd,
	}

	// TaskCmd handles parent task lifecycle and queries.
	TaskCmd = &cobra.Command{
		Use:     "task",
		Aliases: []string{"t"},
		Short:   "Manage parent agent tasks and status rollups",
		RunE:    runDefaultTaskCmd,
	}

	// SubtaskCmd handles atomic subtask queueing, claiming, and completion.
	SubtaskCmd = &cobra.Command{
		Use:     "subtask",
		Aliases: []string{"st", "sub"},
		Short:   "Manage atomic subtasks, worker claiming, and evidence recording",
		RunE:    runDefaultSubtaskCmd,
	}
)

func runDefaultAgentCmd(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func runDefaultTaskCmd(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func runDefaultSubtaskCmd(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

// DispatchAgent routes CLI arguments to agent subcommands with prefix stripping.
func DispatchAgent(args []string) error {
	clean := stripAgentPrefix(args)
	hasClean := len(clean) > 0
	if !hasClean {
		return AgentCmd.Help()
	}
	first := clean[0]
	isHelp := first == "-h" || first == "--help" || first == "help"
	if isHelp {
		return AgentCmd.Help()
	}
	AgentCmd.SetArgs(clean)

	return AgentCmd.Execute()
}

func stripAgentPrefix(args []string) []string {
	hasArgs := len(args) > 0
	if !hasArgs {
		return args
	}
	first := strings.ToLower(args[0])
	isTrigger := isAgentTrigger(first)
	if isTrigger {
		return args[1:]
	}

	return args
}

func isAgentTrigger(token string) bool {
	return token == constants.CmdAgent || token == constants.CmdAgentAlias || token == constants.CmdAgentAlias2
}

func toError(appErr *appfault.AppError) error {
	hasErr := appErr != nil
	if hasErr {
		return appErr
	}

	return nil
}

var (
	logCmd = &cobra.Command{
		Use:                "log",
		Short:              "Record granular in-flight telemetry into agent split database",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunAgentLog(args))
		},
	}

	crashedCmd = &cobra.Command{
		Use:                "crashed",
		Aliases:            []string{"crash"},
		Short:              "Detect crashed or abandoned agents and display autopsy report",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunAgentCrashed(args))
		},
	}

	diagnoseCmd = &cobra.Command{
		Use:                "diagnose",
		Aliases:            []string{"diag", "autopsy"},
		Short:              "Perform deep forensic inspection of agent tasks and logs",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunAgentDiagnose(args))
		},
	}

	clearCmd = &cobra.Command{
		Use:                "clear",
		Short:              "Clear or archive completed agent task runs",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunAgentClear(args))
		},
	}

	resetCmd = &cobra.Command{
		Use:                "reset",
		Short:              "Reset agent databases and telemetry",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunAgentReset(args))
		},
	}

	tempClearCmd = &cobra.Command{
		Use:                "temp-clear",
		Aliases:            []string{"clean-temp", "purge-temp"},
		Short:              "Delete entire ephemeral .ai-memory/temp-agents/ directory",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunAgentTempClear(args))
		},
	}
)

func init() {
	AgentCmd.AddCommand(TaskCmd)
	AgentCmd.AddCommand(SubtaskCmd)
	AgentCmd.AddCommand(logCmd)
	AgentCmd.AddCommand(crashedCmd)
	AgentCmd.AddCommand(diagnoseCmd)
	AgentCmd.AddCommand(clearCmd)
	AgentCmd.AddCommand(resetCmd)
	AgentCmd.AddCommand(tempClearCmd)

	initTaskCommands()
	initSubtaskCommands()
}
