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
)

// AgyWatchPromptsCmd loops and displays running/pending prompts.
var AgyWatchPromptsCmd = &cobra.Command{
	Use:   "watch-prompts",
	Short: "Watch running and pending prompts",
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
	AgyWatchPromptsCmd.Flags().BoolVar(&isWatchPromptsSSH, "ssh", false, "Use SSH to aggregate")
	AgyWatchPromptsCmd.Flags().IntVar(&watchPromptsCompact, "compact", 0, "Compact words")
	AgyWatchPromptsCmd.Flags().IntVar(&watchPromptsCompact, "words", 0, "Words (alias for compact)")
	AgyWatchPromptsCmd.Flags().IntVarP(&watchPromptsCount, "count", "c", 0, "Number of prompts")
	AgyCmd.AddCommand(AgyWatchPromptsCmd)
}

func runAgyWatchPrompts(args []string) *apperror.AppError {
	if isWatchPromptsJSON {
		fmt.Println(`{ "status": "watching", "interval": 30 }`)
		return nil
	}

	fmt.Println("Watching prompts (press Ctrl+C to stop)...")
	for {
		queues, _ := DiscoverAllWorkspaceQueues()
		fmt.Printf("Discovered %d active queues\n", len(queues))
		
		time.Sleep(30 * time.Second)
	}
}
