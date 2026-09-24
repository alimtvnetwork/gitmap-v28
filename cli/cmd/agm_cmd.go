package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

func init() {
	cmdagy.RunAgySSHFn = cmdssh.RunSSHAgyCLI
	cmdinstall.RemoteAgmUpdateFleetFn = func(target, except string) error {
		args := []string{"agm", "--target", target}
		if except != "" {
			args = append(args, "--except", except)
		}
		return cmdssh.RunSSHUpdateCLI(args)
	}
}

func dispatchAgm(ctx context.Context, args []string, root *cobra.Command) error {
	return cmdinstall.DispatchAgm(ctx, args, root)
}
