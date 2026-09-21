package cmd

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai"
)

func stripSubcommandArg(args []string) []string {
	if len(args) == 0 {
		return args
	}

	switch args[0] {
	case "run-smart", "run-incremental", "smart", "test", "smart-test":
		return args[1:]
	default:
		return args
	}
}

// runSmartTestCLI executes the smart incremental Go test runner via python streaming.
func runSmartTestCLI(args []string) error {
	ctx := context.Background()
	scriptPath, err := cmdai.ResolveScriptPath("06-cicd-local-runner.py")
	if err != nil {
		return err
	}

	cleanArgs := stripSubcommandArg(args)
	cmdArgs := append([]string{"--smart"}, cleanArgs...)
	appErr := cmdai.ExecuteScriptStreaming(ctx, scriptPath, cmdArgs)
	if appErr != nil {
		return appErr
	}

	return nil
}
