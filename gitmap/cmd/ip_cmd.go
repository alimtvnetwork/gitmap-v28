package cmd

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

var SSHJoinCmd = &cobra.Command{
	Use:     "ssh-join",
	Aliases: []string{"sj"},
	Short:   "Join an SSH machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHJoin(cmd, args, cmd.Context())
	},
}

func runSSHJoin(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) > 0 {
		switch args[0] {
		case "add", "rm", "ls":
			// Handled by subcommands
		default:
			return apperror.New("runSSHJoin", "E_INTERNAL_ERROR", map[string]any{"arg": args[0]})
		}
	}
	return nil
}

var SJRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove machine-alias or ip",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJRm(cmd, args, cmd.Context())
	},
}

func runSJRm(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) != 1 {
		return apperror.New("runSJRm", "E_INTERNAL_ERROR", map[string]any{"msg": "invalid argument count"})
	}
	return nil
}

func init() {
	SSHJoinCmd.AddCommand(SJRmCmd)
}
