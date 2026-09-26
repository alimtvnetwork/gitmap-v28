package cmdagy
import (
	"fmt"
	"time"
	"github.com/spf13/cobra"
)
var wappSSH bool
var AgyWappCmd = &cobra.Command{
	Use:     "watch-all-projects-prompts",
	Aliases: []string{"wapp"},
	Short:   "Watch all projects prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Watching all projects prompts...")
		for {
			time.Sleep(30 * time.Second)
		}
	},
}
func init() {
	AgyWappCmd.Flags().BoolVar(&wappSSH, "ssh", false, "Use SSH to aggregate")
	AgyCmd.AddCommand(AgyWappCmd)
}
