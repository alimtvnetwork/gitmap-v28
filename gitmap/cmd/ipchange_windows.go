//go:build windows

package cmd

import (
	"context"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

var netshExecutor = exec.CommandContext

//nolint:revive
func swapIPWindows(ctx context.Context, interfaceName string, newIP string) error {
	if interfaceName == "" {
		interfaceName = "Ethernet"
	}
	cmd := netshExecutor(ctx, "netsh", "interface", "ip", "set", "address", "name="+interfaceName, "static", newIP, "255.255.255.0")
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "swapIPWindows", map[string]any{"interface": interfaceName})
	}
	return nil
}

//nolint:revive
func swapIP(ctx context.Context, interfaceName, oldIP, newIP string) error {
	return swapIPWindows(ctx, interfaceName, newIP)
}
