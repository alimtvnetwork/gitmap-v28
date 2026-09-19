package cmd

import "github.com/alimtvnetwork/gitmap-v28/cli/cmd/llm"

func runLlm(args []string) error {
	if appErr := llm.Run(args); appErr != nil {
		return appErr
	}

	return nil
}
