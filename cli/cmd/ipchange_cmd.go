package cmd

import (
	"context"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var IPChangeCmd = &cobra.Command{
	Use:   "ip-change [new-ip]",
	Short: "Change IP address for a machine (requires root/admin privileges)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runIPChangeCmd(cmd, args, cmd.Context())
	},
}

//nolint:revive
func runIPChangeCmd(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) < 1 {
		return apperror.New("runIPChangeCmd", "E_INTERNAL_ERROR", map[string]any{"msg": "requires new-ip argument"})
	}

	newIP := args[0]

	return executeIPChange(ctx, newIP, true)
}

func validatePing(ctx context.Context, targetHost string, count int) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "ping", "-n", strconv.Itoa(count), targetHost)
	} else {
		cmd = exec.CommandContext(ctx, "ping", "-c", strconv.Itoa(count), targetHost)
	}

	err := cmd.Run()

	return err == nil
}

func executeIPChange(ctx context.Context, newIP string, doPing bool) error {
	mgr := createCMDNetIPManager()
	opts, valOpts := parseCMDChangeOptions([]string{newIP}, false)
	valOpts.IsValidationActive = doPing
	res := mgr.ChangeIP(ctx, opts, valOpts)
	if res.IsFailure() {
		return res.AppError()
	}
	if res.Value.IsReverted {
		return apperror.New("executeIPChange", "E_INTERNAL_ERROR", map[string]any{"msg": res.Value.Message})
	}

	return nil
}

func init() {
	// Handled by root or dispatch
}
