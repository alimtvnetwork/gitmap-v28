package cmdssh

import (
	"github.com/spf13/cobra"
)

// SJAddCmd represents the 'add' subcommand for ssh-join.
var SJAddCmd = &cobra.Command{
	Use:     "add <user@ip|ip> [alias] [flags]",
	Aliases: []string{"join", "new", "enroll"},
	Short:   "Enroll an SSH machine by user@ip or IP address",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		return executeEnrollCLI(ctx, args)
	},
}
