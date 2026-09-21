package cmdagy

import (
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdantigravity"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompt"
)

var agyPromptReadCmd = &cobra.Command{
	Use:     "read [conv-id]",
	Aliases: []string{"show", "view", "cat"},
	Short:   "Read and display user prompts for a conversation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAgyPromptRead(args)
	},
}

var agyPromptLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List available prompt templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdprompt.RunPromptList(args)
	},
}

func initAgyPromptSubcommands() {
	cmdantigravity.PromptCmd.AddCommand(agyPromptReadCmd)
	agyPromptCmd.AddCommand(agyPromptReadCmd)
	agyPromptCmd.AddCommand(agyPromptLsCmd)
	cmdantigravity.SetPromptDispatcher(dispatchPromptFromAntigravityCmd)
}

func dispatchPromptFromAntigravityCmd(repoRoot, promptPath, title, content string) (string, bool) {
	ideProcRes := DetectRunningAntigravityIDE()
	pid := resolveActiveOrZeroPID(ideProcRes)
	res := DispatchPromptToAntigravity(repoRoot, promptPath, title, content, pid)

	return res.Message, res.IsSuccess
}
