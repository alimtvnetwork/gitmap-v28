package cmdagy
import (
	"fmt"
	"github.com/spf13/cobra"
)

var AgyLookProjectsSSHCmd = &cobra.Command{
	Use:   "look-projects-ssh",
	Short: "Table of all VM projects via SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching projects over SSH...")
		return nil
	},
}

var AgyRestEnableCmd = &cobra.Command{
	Use:   "rest-enable",
	Short: "Enable REST OS service",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Enabling REST service...")
		return nil
	},
}

func init() {
	AgyCmd.AddCommand(AgyLookProjectsSSHCmd)
	AgyCmd.AddCommand(AgyRestEnableCmd)
}
