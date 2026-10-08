package cmdllm

import "github.com/alimtvnetwork/gitmap-v28/cli/cmd/llm"

func RunLlm(args []string) error {
	if appErr := llm.Run(args); appErr != nil {
		return appErr
	}

	return nil
}
