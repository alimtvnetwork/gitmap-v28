package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// SSHLoginCmd represents the gitmap ssh login command.
var SSHLoginCmd = &cobra.Command{
	Use:   "login [target]",
	Short: "Login via SSH to the specified target",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHLogin(cmd, args, cmd.Context())
	},
}

func runSSHLogin(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) < 1 {
		return apperror.New("runSSHLogin", "E_INTERNAL_ERROR", nil)
	}

	target := args[0]
	_ = target

	return nil
}
