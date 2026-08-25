package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// SSHAliasCmd represents the gitmap ssh as command.
var SSHAliasCmd = &cobra.Command{
	Use:   "as [ip] [alias name]",
	Short: "Create an SSH alias for an IP using the 'as' keyword",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHAlias(cmd, args, cmd.Context())
	},
}

func runSSHAlias(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) < 2 {
		return apperror.New("runSSHAlias", "E_INTERNAL_ERROR", nil)
	}

	ip := args[0]
	aliasName := args[1]
	_ = ip
	_ = aliasName

	return nil
}
