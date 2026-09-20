package cmdagy

import (
	"github.com/spf13/cobra"
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
	Use:     "ls [limit]",
	Aliases: []string{"list"},
	Short:   "List conversation prompts across workspaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAgyPromptLs(args)
	},
}

func initAgyPromptSubcommands() {
	agyPromptCmd.AddCommand(agyPromptReadCmd)
	agyPromptCmd.AddCommand(agyPromptLsCmd)
}
