package cmdagy

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var iptPromptName string
var iptWatch bool
var iptRerun int
var iptPrefix string
var iptSuffix string

var AgyInjectTxtCmd = &cobra.Command{
	Use:     "inject-prompts-txt <prompt-text>",
	Aliases: []string{"ipt"},
	Short:   "Inject raw prompt text into current project conversation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return apperror.NewValidationError("usage: gitmap agy inject-prompts-txt <prompt-text> [--watch] [--rerun N]")
		}
		cwd, _ := os.Getwd()
		prefixText, suffixText := resolveAppliedTemplates(iptPrefix, iptSuffix, false)
		content := decoratePromptBody(args[0], prefixText, suffixText)
		title := iptPromptName
		if title == "" {
			title = "Injected Raw Text Prompt"
		}

		executePromptWithRerun(cwd, content, title, iptRerun, iptWatch)
		return nil
	},
}

func init() {
	AgyInjectTxtCmd.Flags().StringVarP(&iptPromptName, "prompt-name", "p", "", "Prompt title name")
	AgyInjectTxtCmd.Flags().BoolVar(&iptWatch, "watch", false, "Watch until prompt completes")
	AgyInjectTxtCmd.Flags().IntVar(&iptRerun, "rerun", 0, "Rerun loop count")
	AgyInjectTxtCmd.Flags().StringVar(&iptPrefix, "prefix", "", "Prefix template")
	AgyInjectTxtCmd.Flags().StringVar(&iptSuffix, "suffix", "", "Suffix template")
	AgyCmd.AddCommand(AgyInjectTxtCmd)
}
