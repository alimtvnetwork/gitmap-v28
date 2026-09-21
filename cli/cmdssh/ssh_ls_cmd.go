package cmdssh

import (
	"context"
	"os"

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

// SSHNodesCmd represents the gitmap ssh nodes command.
var SSHNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "List all registered SSH nodes/machines",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHLs(cmd, args, cmd.Context())
	},
}

// RunSSHNodesCLI executes the ssh nodes/ls command.
func RunSSHNodesCLI(ctx context.Context, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		return dispatchNodesArgs(ctx, args)
	}

	return printSJList(ctx, os.Stdout, 0)
}

func dispatchNodesArgs(ctx context.Context, args []string) error {
	sub := args[0]
	isClear := sub == "clear"
	if isClear {
		return runSSHClearNodes(args[1:])
	}

	isRm := sub == "rm" || sub == "remove" || sub == "delete"
	if isRm {
		return runSJRm(nil, args[1:], ctx)
	}

	isReset := sub == "reset"
	if isReset {
		return RunSSHResetCLI(args[1:])
	}

	return printSJList(ctx, os.Stdout, 0)
}

//nolint:revive
func runSSHLs(cmd *cobra.Command, args []string, ctx context.Context) error {
	return RunSSHNodesCLI(ctx, args)
}
