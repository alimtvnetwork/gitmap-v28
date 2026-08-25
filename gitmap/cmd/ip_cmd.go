package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

var SSHJoinCmd = &cobra.Command{
	Use:     "ssh-join",
	Aliases: []string{"sj", "ssh-joined", "ssh-joiner"},
	Short:   "Join an SSH machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHJoin(cmd, args, cmd.Context())
	},
}

func runSSHJoin(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) > 0 {
		switch args[0] {
		case "add", "rm", "ls", "history", "add-auth":
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

var IPCmd = &cobra.Command{
	Use:   "ip",
	Short: "Print local IP",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runIPCmd(cmd, args, cmd.Context())
	},
}

func runIPCmd(cmd *cobra.Command, args []string, ctx context.Context) error {
	return executeIPCmd(ctx, true, os.Stdout)
}

func executeIPCmd(ctx context.Context, skipLoopback bool, writer io.Writer) error {
	ipStr, err := GetLocalIP(ctx, skipLoopback, "")
	if err != nil {
		return apperror.New("executeIPCmd", "E_INTERNAL_ERROR", map[string]any{"err": err.Error()})
	}
	if skipLoopback && (ipStr == "127.0.0.1" || ipStr == "::1") {
		return apperror.New("executeIPCmd", "E_INTERNAL_ERROR", map[string]any{"err": "only loopback found"})
	}

	fmt.Fprintln(writer, ipStr)
	return nil
}

func init() {
	SSHJoinCmd.AddCommand(SJRmCmd)
	SSHJoinCmd.AddCommand(SJAddAuthCmd)
	SSHJoinCmd.AddCommand(SJLsCmd)
	SSHJoinCmd.AddCommand(SJHistCmd)
}
