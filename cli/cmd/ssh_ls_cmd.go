package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// SSHLsCmd represents the gitmap ssh ls command.
var SSHLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all joined SSH machines",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHLs(cmd, args, cmd.Context())
	},
}

//nolint:revive
func runSSHLs(cmd *cobra.Command, args []string, ctx context.Context) error {
	return nil
}
