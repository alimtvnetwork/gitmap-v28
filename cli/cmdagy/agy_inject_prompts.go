package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var (
	injectPromptsPrefix  string
	injectPromptsRerun   int
	isInjectPromptsWatch bool
	injectPromptsSSH     string
)

// AgyInjectPromptsCmd provides iterative prompt injection.
var AgyInjectPromptsCmd = &cobra.Command{
	Use:     "inject-prompts <folder-or-text>",
	Aliases: []string{"ip", "inject-prompts-ssh"},
	Short:   "Inject prompts iteratively with templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgyInjectPrompts(args)
		if appErr != nil {
			return appErr
		}
		return nil
	},
}

func init() {
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "prefix", "", "Prefix template")
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "pfx", "", "Prefix template (alias)")
	AgyInjectPromptsCmd.Flags().BoolVar(&isInjectPromptsWatch, "ui-ux", false, "Use UI/UX template")
	AgyInjectPromptsCmd.Flags().BoolVar(&isInjectPromptsWatch, "uu", false, "Use UI/UX template (alias)")
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "suffix", "", "Suffix template")
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsSSH, "nodes", "", "Nodes")
	AgyInjectPromptsCmd.Flags().IntVar(&injectPromptsRerun, "rerun", 0, "Rerun loop count")
	AgyInjectPromptsCmd.Flags().BoolVar(&isInjectPromptsWatch, "watch", false, "Watch until prompt completes")
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsSSH, "ssh-only-node", "", "SSH node targeting")
	AgyCmd.AddCommand(AgyInjectPromptsCmd)
}

func runAgyInjectPrompts(args []string) *apperror.AppError {
	if len(args) == 0 {
		return apperror.NewSimple("ip requires an argument", "EIP001")
	}
	
	fmt.Printf("Injecting prompt from: %s\n", args[0])
	fmt.Printf("Rerun count: %d\n", injectPromptsRerun)
	
	return nil
}
