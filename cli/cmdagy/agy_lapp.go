package cmdagy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var lappSSH bool

var AgyLappCmd = &cobra.Command{
	Use:     "look-all-projects-prompts",
	Aliases: []string{"lapp"},
	Short:   "Look at all projects prompts across local and cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		snapshot, err := CollectPromptsSnapshot(true, 0, 0, lappSSH)
		if err != nil {
			return apperror.WrapSimple(err, "collect all project prompts")
		}

		RenderPromptsTable(snapshot)
		return nil
	},
}

func init() {
	AgyLappCmd.Flags().BoolVar(&lappSSH, "ssh", false, "Use SSH to aggregate across cluster nodes")
	AgyCmd.AddCommand(AgyLappCmd)
}
