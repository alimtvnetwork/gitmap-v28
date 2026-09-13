package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
)

func dispatchAgm(ctx context.Context, args []string, root *cobra.Command) error {
	return cmdinstall.DispatchAgm(ctx, args, root)
}
