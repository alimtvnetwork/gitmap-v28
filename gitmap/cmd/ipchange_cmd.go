package cmd

import (
	"github.com/spf13/cobra"
)

var IPChangeCmd = &cobra.Command{
	Use:   "ip-change",
	Short: "Change IP address for a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	// Usually this is added to RootCmd or another parent command.
	// For now we scaffold the command itself.
}
