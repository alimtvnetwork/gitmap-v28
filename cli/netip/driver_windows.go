package netip

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// WindowsDriver configures network interfaces via netsh.
type WindowsDriver struct{}

// NewWindowsDriver creates a new Windows netsh driver.
func NewWindowsDriver() *WindowsDriver {
	return &WindowsDriver{}
}

// Name returns the driver identifier.
func (d *WindowsDriver) Name() string {
	return DriverNameWindows
}

// ListInterfaces lists system network adapters.
func (d *WindowsDriver) ListInterfaces(ctx context.Context) InterfaceSliceResult {
	return ListSystemInterfaces()
}

// GetInterface retrieves a specific interface.
func (d *WindowsDriver) GetInterface(ctx context.Context, name string) InterfaceResult {
	return GetSystemInterface(name)
}

// ApplyConfig configures static or DHCP IP and DNS on Windows.
func (d *WindowsDriver) ApplyConfig(ctx context.Context, opts ChangeOptions) *apperror.AppError {
	if opts.IsDHCP {
		return applyWindowsDHCP(ctx, opts.InterfaceName)
	}

	return applyWindowsStatic(ctx, opts)
}

// RevertConfig restores previous configuration from snapshot.
func (d *WindowsDriver) RevertConfig(ctx context.Context, snap Snapshot) *apperror.AppError {
	if snap.IsDHCP {
		return applyWindowsDHCP(ctx, snap.InterfaceName)
	}

	opts := ChangeOptions{
		InterfaceName: snap.InterfaceName,
		IP:            snap.IP,
		Netmask:       snap.Netmask,
		Gateway:       snap.Gateway,
		DNS:           snap.DNS,
		IsDHCP:        snap.IsDHCP,
	}

	return applyWindowsStatic(ctx, opts)
}

func applyWindowsStatic(ctx context.Context, opts ChangeOptions) *apperror.AppError {
	err := setWindowsAddress(ctx, opts.InterfaceName, opts.IP, opts.Netmask, opts.Gateway)
	if err != nil {
		return err
	}

	return setWindowsDNS(ctx, opts.InterfaceName, opts.DNS)
}

func setWindowsAddress(ctx context.Context, iface, ip, netmask, gateway string) *apperror.AppError {
	args := []string{"interface", "ip", "set", "address", "name=" + iface, "static", ip, netmask}
	if gateway != "" {
		args = append(args, gateway)
	}

	return runNetsh(ctx, args...)
}

func setWindowsDNS(ctx context.Context, iface string, dnsList []string) *apperror.AppError {
	if len(dnsList) == 0 {
		return nil
	}

	err := runNetsh(ctx, "interface", "ip", "set", "dns", "name="+iface, "static", dnsList[0])
	if err != nil {
		return err
	}

	return addSecondaryDNSList(ctx, iface, dnsList[1:])
}

func addSecondaryDNSList(ctx context.Context, iface string, secondary []string) *apperror.AppError {
	for i, dns := range secondary {
		idxStr := fmt.Sprintf("index=%d", i+2)
		if err := runNetsh(ctx, "interface", "ip", "add", "dns", "name="+iface, dns, idxStr); err != nil {
			return err
		}
	}

	return nil
}

func applyWindowsDHCP(ctx context.Context, iface string) *apperror.AppError {
	err := runNetsh(ctx, "interface", "ip", "set", "address", "name="+iface, "dhcp")
	if err != nil {
		return err
	}

	return runNetsh(ctx, "interface", "ip", "set", "dns", "name="+iface, "dhcp")
}

func runNetsh(ctx context.Context, args ...string) *apperror.AppError {
	cmd := exec.CommandContext(ctx, "netsh", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		msg := fmt.Sprintf("netsh %s failed: %s", strings.Join(args, " "), string(out))
		return apperror.WrapWithDetails(err, "runNetsh", "E_NETSH_EXEC", msg, "driver_windows", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}
