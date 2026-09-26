package cmdagy

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var wappSSH bool

var AgyWappCmd = &cobra.Command{
	Use:     "watch-all-projects-prompts",
	Aliases: []string{"wapp"},
	Short:   "Watch all projects prompts across local and cluster nodes (30s interval)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Watching all projects prompts at 30s interval (press Ctrl+C to stop)...")
		for {
			snapshot, err := CollectPromptsSnapshot(true, 0, 0, wappSSH)
			if err != nil {
				return apperror.WrapSimple(err, "collect all project prompts")
			}

			RenderPromptsTable(snapshot)
			time.Sleep(30 * time.Second)
		}
	},
}

func init() {
	AgyWappCmd.Flags().BoolVar(&wappSSH, "ssh", false, "Use SSH to aggregate across cluster nodes")
	AgyCmd.AddCommand(AgyWappCmd)
}
