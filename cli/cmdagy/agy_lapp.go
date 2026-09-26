package cmdagy
import (
	"fmt"
	"github.com/spf13/cobra"
)
var lappSSH bool
var AgyLappCmd = &cobra.Command{
	Use:     "look-all-projects-prompts",
	Aliases: []string{"lapp"},
	Short:   "Look at all projects prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Looking at all projects prompts...")
		return nil
	},
}
func init() {
	AgyLappCmd.Flags().BoolVar(&lappSSH, "ssh", false, "Use SSH to aggregate")
	AgyCmd.AddCommand(AgyLappCmd)
}
