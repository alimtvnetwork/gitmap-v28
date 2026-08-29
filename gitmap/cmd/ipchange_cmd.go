package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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

//nolint:revive
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

//nolint:revive
func executeIPChange(ctx context.Context, newIP string, doPing bool) error {
	var swapErr error
	interfaceName := ""
	if runtime.GOOS == "windows" {
		interfaceName = "Ethernet"
	}
	swapErr = swapIP(ctx, interfaceName, "", newIP)

	if swapErr != nil {
		return apperror.Wrap(swapErr, "executeIPChange", map[string]any{"ip": newIP})
	}

	if doPing {
		if !validatePing(ctx, "8.8.8.8", 3) {
			fmt.Println("reverting")
			// swap back placeholder
			_ = swapIP(ctx, interfaceName, newIP, "192.168.1.100") // rollback
			return apperror.New("executeIPChange", "E_INTERNAL_ERROR", map[string]any{"msg": "ping failed, reverting"})
		}
	}
	return nil
}

func init() {
	// Handled by root or dispatch
}
