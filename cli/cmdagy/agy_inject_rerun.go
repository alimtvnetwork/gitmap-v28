package cmdagy

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var iprPromptName string
var iprWatch bool
var iprPrefix string
var iprSuffix string
var iprNodes string
var iprN int

var AgyInjectRerunCmd = &cobra.Command{
	Use:     "inject-prompts-rerun <prompt-text>",
	Aliases: []string{"ipr"},
	Short:   "Inject prompt and rerun for N iterations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: gitmap agy inject-prompts-rerun <prompt-text> [-n 2] [--watch]")
		}
		cwd, _ := os.Getwd()
		prefixText, suffixText := resolveAppliedTemplates(iprPrefix, iprSuffix, false)
		content := decoratePromptBody(args[0], prefixText, suffixText)
		title := iprPromptName
		if title == "" {
			title = "Injected Rerun Prompt"
		}

		executePromptWithRerun(cwd, content, title, iprN, iprWatch)
		return nil
	},
}

var AgyInjectRerunSSHCmd = &cobra.Command{
	Use:     "inject-prompts-rerun-ssh <prompt-text>",
	Aliases: []string{"iprs"},
	Short:   "Inject prompt and rerun via SSH cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: gitmap agy inject-prompts-rerun-ssh <prompt-text> [--nodes ...] [-n 2]")
		}
		fmt.Printf("Injecting and rerunning across SSH nodes (%s)...\n", iprNodes)
		cwd, _ := os.Getwd()
		prefixText, suffixText := resolveAppliedTemplates(iprPrefix, iprSuffix, false)
		content := decoratePromptBody(args[0], prefixText, suffixText)
		title := iprPromptName
		if title == "" {
			title = "SSH Injected Rerun Prompt"
		}

		executePromptWithRerun(cwd, content, title, iprN, iprWatch)
		return nil
	},
}

func init() {
	AgyInjectRerunCmd.Flags().StringVarP(&iprPromptName, "prompt-name", "p", "", "Prompt title name")
	AgyInjectRerunCmd.Flags().BoolVar(&iprWatch, "watch", false, "Watch until prompt completes")
	AgyInjectRerunCmd.Flags().IntVarP(&iprN, "number", "n", 2, "Rerun count")
	AgyInjectRerunCmd.Flags().StringVar(&iprPrefix, "prefix", "", "Prefix template")
	AgyInjectRerunCmd.Flags().StringVar(&iprSuffix, "suffix", "", "Suffix template")
	AgyInjectRerunCmd.Flags().StringVar(&iprNodes, "ssh-only-node", "", "Specific SSH node")
	AgyCmd.AddCommand(AgyInjectRerunCmd)

	AgyInjectRerunSSHCmd.Flags().StringVarP(&iprPromptName, "prompt-name", "p", "", "Prompt title name")
	AgyInjectRerunSSHCmd.Flags().BoolVar(&iprWatch, "watch", false, "Watch until prompt completes")
	AgyInjectRerunSSHCmd.Flags().IntVarP(&iprN, "number", "n", 2, "Rerun count")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprPrefix, "prefix", "", "Prefix template")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprSuffix, "suffix", "", "Suffix template")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprNodes, "nodes", "", "Specific SSH nodes")
	AgyCmd.AddCommand(AgyInjectRerunSSHCmd)
}
