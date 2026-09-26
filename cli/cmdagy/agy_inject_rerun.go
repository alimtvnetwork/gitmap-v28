package cmdagy
import (
	"fmt"
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
	Short:   "Inject prompt and rerun",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Injecting and rerunning...")
		return nil
	},
}

var AgyInjectRerunSSHCmd = &cobra.Command{
	Use:     "inject-prompts-rerun-ssh <prompt-text>",
	Aliases: []string{"iprs"},
	Short:   "Inject prompt and rerun via SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Injecting and rerunning over SSH...")
		return nil
	},
}

func init() {
	AgyInjectRerunCmd.Flags().StringVarP(&iprPromptName, "prompt-name", "p", "", "Prompt name")
	AgyInjectRerunCmd.Flags().BoolVar(&iprWatch, "watch", false, "Watch until prompt completes")
	AgyInjectRerunCmd.Flags().IntVarP(&iprN, "number", "n", 2, "Rerun count")
	AgyInjectRerunCmd.Flags().StringVar(&iprPrefix, "prefix", "", "Prefix template")
	AgyInjectRerunCmd.Flags().StringVar(&iprSuffix, "suffix", "", "Suffix template")
	AgyInjectRerunCmd.Flags().StringVar(&iprNodes, "ssh-only-node", "", "Specific SSH node")
	AgyCmd.AddCommand(AgyInjectRerunCmd)

	AgyInjectRerunSSHCmd.Flags().StringVarP(&iprPromptName, "prompt-name", "p", "", "Prompt name")
	AgyInjectRerunSSHCmd.Flags().BoolVar(&iprWatch, "watch", false, "Watch until prompt completes")
	AgyInjectRerunSSHCmd.Flags().IntVarP(&iprN, "number", "n", 2, "Rerun count")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprPrefix, "prefix", "", "Prefix template")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprSuffix, "suffix", "", "Suffix template")
	AgyInjectRerunSSHCmd.Flags().StringVar(&iprNodes, "nodes", "", "Specific SSH nodes")
	AgyCmd.AddCommand(AgyInjectRerunSSHCmd)
}
