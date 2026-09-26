package cmdagy

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var (
	isWatchPromptsJSON  bool
	isWatchPromptsAll   bool
	isWatchPromptsSSH   bool
	watchPromptsCompact int
	watchPromptsCount   int
	watchPromptsFile    string
)

// AgyWatchPromptsCmd loops and displays running/pending prompts.
var AgyWatchPromptsCmd = &cobra.Command{
	Use:   "watch-prompts",
	Short: "Watch running and pending prompts with 30-second interval",
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgyWatchPrompts(args)
		if appErr != nil {
			return appErr
		}
		return nil
	},
}

func init() {
	AgyWatchPromptsCmd.Flags().BoolVarP(&isWatchPromptsJSON, "json", "j", false, "JSON output")
	AgyWatchPromptsCmd.Flags().BoolVarP(&isWatchPromptsAll, "all", "a", false, "All projects")
	AgyWatchPromptsCmd.Flags().BoolVar(&isWatchPromptsSSH, "ssh", false, "Use SSH to aggregate across cluster nodes")
	AgyWatchPromptsCmd.Flags().StringVarP(&watchPromptsFile, "file", "f", "", "Output results to file")
	AgyWatchPromptsCmd.Flags().IntVar(&watchPromptsCompact, "compact", 0, "Compact prompt words")
	AgyWatchPromptsCmd.Flags().IntVar(&watchPromptsCompact, "words", 0, "Words (alias for compact)")
	AgyWatchPromptsCmd.Flags().IntVarP(&watchPromptsCount, "count", "c", 0, "Number of prompts to show")
	AgyCmd.AddCommand(AgyWatchPromptsCmd)
}

func runAgyWatchPrompts(args []string) *apperror.AppError {
	fmt.Println("Watching prompts at 30s interval (press Ctrl+C to stop)...")

	for {
		snapshot, err := CollectPromptsSnapshot(isWatchPromptsAll, watchPromptsCompact, watchPromptsCount, isWatchPromptsSSH)
		if err != nil {
			return apperror.WrapSimple(err, "collect prompts snapshot")
		}

		if isWatchPromptsJSON || watchPromptsFile != "" {
			return handlePromptsJSONOutput(snapshot, watchPromptsFile)
		}

		RenderPromptsTable(snapshot)
		time.Sleep(30 * time.Second)
	}
}
