package cmdagy
import (
	"fmt"
	"github.com/spf13/cobra"
)
var iptPromptName string
var iptWatch bool
var iptRerun int
var iptPrefix string
var iptSuffix string

var AgyInjectTxtCmd = &cobra.Command{
	Use:     "inject-prompts-txt <prompt-text>",
	Aliases: []string{"ipt"},
	Short:   "Inject raw prompt text",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Injecting text prompt...")
		return nil
	},
}
func init() {
	AgyInjectTxtCmd.Flags().StringVarP(&iptPromptName, "prompt-name", "p", "", "Prompt name")
	AgyInjectTxtCmd.Flags().BoolVar(&iptWatch, "watch", false, "Watch until prompt completes")
	AgyInjectTxtCmd.Flags().IntVar(&iptRerun, "rerun", 1, "Rerun count")
	AgyInjectTxtCmd.Flags().StringVar(&iptPrefix, "prefix", "", "Prefix template")
	AgyInjectTxtCmd.Flags().StringVar(&iptSuffix, "suffix", "", "Suffix template")
	AgyCmd.AddCommand(AgyInjectTxtCmd)
}
