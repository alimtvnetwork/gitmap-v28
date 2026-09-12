package cmd

import "github.com/alimtvnetwork/gitmap-v28/cli/cmd/llm"

func runLlm(args []string) error {
	llm.Run(args)

	return nil
}
