package cmd

import (
	"context"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func swapIP(ctx context.Context, interfaceName, oldIP, newIP string) error {
	cmdAdd := exec.CommandContext(ctx, "ip", "addr", "add", newIP, "dev", "eth0")
	if err := cmdAdd.Run(); err != nil {
		return apperror.Wrap(err, "swapIPLinux", map[string]any{"op": "add"})
	}
	if oldIP != "" {
		cmdDel := exec.CommandContext(ctx, "ip", "addr", "del", oldIP, "dev", "eth0")
		if err := cmdDel.Run(); err != nil {
			return apperror.Wrap(err, "swapIPLinux", map[string]any{"op": "del"})
		}
	}
	return nil
}
