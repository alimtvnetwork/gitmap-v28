package cmdagy

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdantigravity"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompt"
)

var (
	sendPromptProjectFlag string
	sendPromptMsgFlag     string
	sendPromptTitleFlag   string
)

var agySendPromptCmd = &cobra.Command{
	Use:     "send-prompt",
	Aliases: []string{"sendprompt", "sp"},
	Short:   "Send a prompt to target Antigravity project session",
	RunE:    RunAgySendPrompt,
}

var agyPromptSendCmd = &cobra.Command{
	Use:     "send",
	Aliases: []string{"dispatch"},
	Short:   "Send a prompt to target Antigravity project session",
	RunE:    RunAgySendPrompt,
}

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

// RunAgySendPrompt executes prompt dispatch to a target Antigravity project.
func RunAgySendPrompt(cmd *cobra.Command, args []string) error {
	proj, msg, title := resolveSendPromptArgs(args)

	return ExecuteSendPrompt(proj, msg, title)
}

func resolveSendPromptArgs(args []string) (string, string, string) {
	proj := sendPromptProjectFlag
	msg := sendPromptMsgFlag
	title := sendPromptTitleFlag

	if proj == "" && len(args) > 0 {
		proj = args[0]
	}
	if msg == "" && len(args) > 1 {
		msg = strings.Join(args[1:], " ")
	}

	return proj, msg, title
}

func initSendPromptFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&sendPromptProjectFlag, "project", "p", "", "Target project slug, alias, ID, or path")
	cmd.Flags().StringVarP(&sendPromptMsgFlag, "prompt", "m", "", "Prompt message content or file path")
	cmd.Flags().StringVarP(&sendPromptTitleFlag, "title", "t", "", "Custom prompt title")
}

func init() {
	agyPromptInjectCmd.Flags().StringVarP(&injectProjectFlag, "project", "p", "", "Target project path or alias")
	agyPromptInjectCmd.Flags().StringVarP(&injectConvFlag, "conversation", "c", "", "Target conversation ID or default")
	agyPromptInjectCmd.Flags().StringVarP(&injectTitleFlag, "title", "t", "", "Custom prompt title")
	agyPromptInjectCmd.Flags().BoolVar(&injectSkipFlag, "skip-inject", false, "Stage prompt file only without injecting to queue")

	initSendPromptFlags(agySendPromptCmd)
	initSendPromptFlags(agyPromptSendCmd)
}

func initAgyPromptSubcommands() {
	cmdantigravity.PromptCmd.AddCommand(agyPromptReadCmd)
	cmdantigravity.PromptCmd.AddCommand(agyPromptInjectCmd)
	cmdantigravity.PromptCmd.AddCommand(agySendPromptCmd)
	AgyCmd.AddCommand(agySendPromptCmd)
	agyPromptCmd.AddCommand(agyPromptReadCmd)
	agyPromptCmd.AddCommand(agyPromptLsCmd)
	agyPromptCmd.AddCommand(agyPromptInjectCmd)
	agyPromptCmd.AddCommand(agyPromptSendCmd)
	agyPromptCmd.AddCommand(AgyInjectPromptsCmd)
	agyPromptCmd.AddCommand(AgyEnhancePromptsCmd)
	agyPromptCmd.AddCommand(cmdprompt.PromptAddCmd)
	agyPromptCmd.RunE = RunAgySendPrompt
	cmdantigravity.SetPromptDispatcher(dispatchPromptFromAntigravityCmd)
}

func dispatchPromptFromAntigravityCmd(repoRoot, promptPath, title, content string) (string, bool) {
	ideProcRes := DetectRunningAntigravityIDE()
	pid := resolveActiveOrZeroPID(ideProcRes)
	res := DispatchPromptToAntigravity(repoRoot, promptPath, title, content, pid)

	return res.Message, res.IsSuccess
}
