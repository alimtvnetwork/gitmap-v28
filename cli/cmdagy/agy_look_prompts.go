package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var (
	isLookPromptsJSON  bool
	isLookPromptsAll   bool
	isLookPromptsSSH   bool
	lookPromptsCompact int
	lookPromptsCount   int
	lookPromptsFile    string
)

var AgyLookPromptsCmd = &cobra.Command{
	Use:   "look-prompts",
	Short: "Snapshot of running and pending prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgyLookPrompts(args)
		if appErr != nil {
			return appErr
		}
		return nil
	},
}

func init() {
	AgyLookPromptsCmd.Flags().BoolVarP(&isLookPromptsJSON, "json", "j", false, "JSON output")
	AgyLookPromptsCmd.Flags().BoolVarP(&isLookPromptsAll, "all", "a", false, "All projects")
	AgyLookPromptsCmd.Flags().BoolVar(&isLookPromptsSSH, "ssh", false, "Use SSH to aggregate")
	AgyLookPromptsCmd.Flags().StringVarP(&lookPromptsFile, "file", "f", "", "Output to file")
	AgyLookPromptsCmd.Flags().IntVar(&lookPromptsCompact, "compact", 0, "Compact words")
	AgyLookPromptsCmd.Flags().IntVar(&lookPromptsCompact, "words", 0, "Words (alias for compact)")
	AgyLookPromptsCmd.Flags().IntVarP(&lookPromptsCount, "count", "c", 0, "Number of prompts")
	AgyCmd.AddCommand(AgyLookPromptsCmd)
}

func runAgyLookPrompts(args []string) *apperror.AppError {
	queues, _ := DiscoverAllWorkspaceQueues()
	fmt.Printf("Discovered %d active queues\n", len(queues))
	return nil
}
