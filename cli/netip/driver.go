package netip

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func isNetIPTestMode() bool {
	if flag.Lookup("test.v") != nil {
		return true
	}
	if os.Getenv("GITMAP_TEST") != "" {
		return true
	}
	return strings.HasSuffix(os.Args[0], ".test") || strings.HasSuffix(os.Args[0], ".test.exe")
}

const (
	// DriverNameWindows identifies the Windows netsh backend.
	DriverNameWindows = "windows"

	// DriverNameLinuxNetplan identifies the Netplan v2 YAML backend.
	DriverNameLinuxNetplan = "linux_netplan"

	// DriverNameLinuxNMCLI identifies the NetworkManager backend.
	DriverNameLinuxNMCLI = "linux_nmcli"

	// DriverNameLinuxIPRoute2 identifies the ephemeral iproute2 container fallback backend.
	DriverNameLinuxIPRoute2 = "linux_iproute2"
)

// Driver abstracts platform-specific network configuration engines.
type Driver interface {
	Name() string
	ApplyConfig(ctx context.Context, opts ChangeOptions) *apperror.AppError
	RevertConfig(ctx context.Context, snap Snapshot) *apperror.AppError
	ListInterfaces(ctx context.Context) InterfaceSliceResult
	GetInterface(ctx context.Context, name string) InterfaceResult
}

func snapshotToOptions(snap Snapshot) ChangeOptions {
	return ChangeOptions{
		InterfaceName: snap.InterfaceName,
		IP:            snap.IP,
		Netmask:       snap.Netmask,
		Gateway:       snap.Gateway,
		DNS:           snap.DNS,
		IsDHCP:        snap.IsDHCP,
	}
}
