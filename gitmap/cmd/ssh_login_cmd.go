package cmd

import (
	"context"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
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
	return executeSSHLogin(ctx, target, false)
}

func executeSSHLogin(ctx context.Context, target string, force bool) error {
	sshTarget, err := ParseSSHTarget(target, "root", 22)
	if err != nil {
		return err
	}

	if !strings.Contains(target, "@") && net.ParseIP(target) == nil {
		db, err := store.OpenDefault()
		if err == nil {
			host, err := store.GetHostByAlias(ctx, target, db.Conn())
			if err == nil {
				sshTarget.Username = host.Username
				sshTarget.IP = host.IP
			}
			db.Close()
		}
	}

	return SpawnSSH(ctx, *sshTarget, nil)
}
