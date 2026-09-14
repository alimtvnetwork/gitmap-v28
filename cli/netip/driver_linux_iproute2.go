package netip

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// LinuxIPRoute2Driver manages IP addresses directly via iproute2 for container environments.
type LinuxIPRoute2Driver struct{}

// NewLinuxIPRoute2Driver creates a new iproute2 container fallback driver.
func NewLinuxIPRoute2Driver() *LinuxIPRoute2Driver {
	return &LinuxIPRoute2Driver{}
}

// Name returns the driver identifier.
func (d *LinuxIPRoute2Driver) Name() string {
	return DriverNameLinuxIPRoute2
}

// ListInterfaces lists system network adapters.
func (d *LinuxIPRoute2Driver) ListInterfaces(ctx context.Context) InterfaceSliceResult {
	return ListSystemInterfaces()
}

// GetInterface retrieves a specific interface.
func (d *LinuxIPRoute2Driver) GetInterface(ctx context.Context, name string) InterfaceResult {
	return GetSystemInterface(name)
}

// ApplyConfig configures interface IP, link state, and default route via iproute2.
func (d *LinuxIPRoute2Driver) ApplyConfig(ctx context.Context, opts ChangeOptions) *apperror.AppError {
	cidr := FormatCIDR(opts.IP, opts.Netmask)
	if err := applyIPRoute2Address(ctx, opts.InterfaceName, cidr); err != nil {
		return err
	}

	if err := setIPRoute2LinkUp(ctx, opts.InterfaceName); err != nil {
		return err
	}

	return setIPRoute2DefaultGateway(ctx, opts.InterfaceName, opts.Gateway)
}

// RevertConfig restores previous configuration from snapshot.
func (d *LinuxIPRoute2Driver) RevertConfig(ctx context.Context, snap Snapshot) *apperror.AppError {
	opts := snapshotToOptions(snap)

	return d.ApplyConfig(ctx, opts)
}

func applyIPRoute2Address(ctx context.Context, iface, cidr string) *apperror.AppError {
	return runIPRoute2(ctx, "addr", "add", cidr, "dev", iface)
}

func setIPRoute2LinkUp(ctx context.Context, iface string) *apperror.AppError {
	return runIPRoute2(ctx, "link", "set", iface, "up")
}

func setIPRoute2DefaultGateway(ctx context.Context, iface, gateway string) *apperror.AppError {
	if gateway == "" {
		return nil
	}

	return runIPRoute2(ctx, "route", "replace", "default", "via", gateway, "dev", iface)
}

func runIPRoute2(ctx context.Context, args ...string) *apperror.AppError {
	if isNetIPTestMode() {
		return nil
	}

	cmd := exec.CommandContext(ctx, "ip", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		msg := fmt.Sprintf("ip %s failed: %s", strings.Join(args, " "), string(out))
		return apperror.WrapWithDetails(err, "runIPRoute2", "E_IPROUTE2_EXEC", msg, "driver_iproute2", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}
