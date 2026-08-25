package cmd

import (
	"context"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func swapIP(ctx context.Context, interfaceName, oldIP, newIP string) error {
	cmd := exec.CommandContext(ctx, "netsh", "interface", "ip", "set", "address", "name="+interfaceName, "static", newIP, "255.255.255.0")
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "swapIPWindows", map[string]any{"interface": interfaceName})
	}
	return nil
}
