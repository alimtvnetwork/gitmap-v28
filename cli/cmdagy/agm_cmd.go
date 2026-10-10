package cmdagy

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/spf13/cobra"
)

func init() {
	cmdinstall.RunAGMVersionTagsLSFn = RunAGMVersionTagsLS
	cmdinstall.RegisterAgyInstallSubcommand(AgyCmd)
	cmdinstall.RemoteAgmUpdateFleetFn = func(target, except string) error {
		args := []string{"agm", "--target", target}
		if except != "" {
			args = append(args, "--except", except)
		}
		return RunSSHUpdateFn(args)
	}
}

func DispatchAgm(ctx context.Context, args []string, root *cobra.Command) error {
	return cmdinstall.DispatchAgm(ctx, args, root)
}
