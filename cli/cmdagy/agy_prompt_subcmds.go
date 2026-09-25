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

var agyPromptInjectCmd = &cobra.Command{
	Use:     "inject <slug-or-file> [target-project]",
	Aliases: []string{"inj"},
	Short:   "Inject a prompt into an Antigravity project and conversation",
	RunE:    RunAgyPromptInject,
}

func init() {
	agyPromptInjectCmd.Flags().StringVarP(&injectProjectFlag, "project", "p", "", "Target project path or alias")
	agyPromptInjectCmd.Flags().StringVarP(&injectConvFlag, "conversation", "c", "", "Target conversation ID or default")
	agyPromptInjectCmd.Flags().StringVarP(&injectTitleFlag, "title", "t", "", "Custom prompt title")
	agyPromptInjectCmd.Flags().BoolVar(&injectSkipFlag, "skip-inject", false, "Stage prompt file only without injecting to queue")
}

func initAgyPromptSubcommands() {
	cmdantigravity.PromptCmd.AddCommand(agyPromptReadCmd)
	cmdantigravity.PromptCmd.AddCommand(agyPromptInjectCmd)
	agyPromptCmd.AddCommand(agyPromptReadCmd)
	agyPromptCmd.AddCommand(agyPromptLsCmd)
	agyPromptCmd.AddCommand(agyPromptInjectCmd)
	agyPromptCmd.AddCommand(cmdprompt.PromptAddCmd)
	cmdantigravity.SetPromptDispatcher(dispatchPromptFromAntigravityCmd)
}

func dispatchPromptFromAntigravityCmd(repoRoot, promptPath, title, content string) (string, bool) {
	ideProcRes := DetectRunningAntigravityIDE()
	pid := resolveActiveOrZeroPID(ideProcRes)
	res := DispatchPromptToAntigravity(repoRoot, promptPath, title, content, pid)

	return res.Message, res.IsSuccess
}
