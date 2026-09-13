package netip

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// LinuxNMCLIDriver manages network interfaces via NetworkManager nmcli.
type LinuxNMCLIDriver struct{}

// NewLinuxNMCLIDriver creates a new nmcli driver instance.
func NewLinuxNMCLIDriver() *LinuxNMCLIDriver {
	return &LinuxNMCLIDriver{}
}

// Name returns the driver identifier.
func (d *LinuxNMCLIDriver) Name() string {
	return DriverNameLinuxNMCLI
}

// ListInterfaces lists system network adapters.
func (d *LinuxNMCLIDriver) ListInterfaces(ctx context.Context) InterfaceSliceResult {
	return ListSystemInterfaces()
}

// GetInterface retrieves a specific interface.
func (d *LinuxNMCLIDriver) GetInterface(ctx context.Context, name string) InterfaceResult {
	return GetSystemInterface(name)
}

// ApplyConfig configures NetworkManager connection via nmcli.
func (d *LinuxNMCLIDriver) ApplyConfig(ctx context.Context, opts ChangeOptions) *apperror.AppError {
	if opts.IsDHCP {
		return applyNMCLIDHCP(ctx, opts.InterfaceName)
	}

	return applyNMCLIStatic(ctx, opts)
}

// RevertConfig restores previous configuration from snapshot.
func (d *LinuxNMCLIDriver) RevertConfig(ctx context.Context, snap Snapshot) *apperror.AppError {
	opts := snapshotToOptions(snap)

	return d.ApplyConfig(ctx, opts)
}

func applyNMCLIDHCP(ctx context.Context, iface string) *apperror.AppError {
	if err := runNMCLI(ctx, "connection", "modify", iface, "ipv4.method", "auto"); err != nil {
		return err
	}

	return upNMCLIConnection(ctx, iface)
}

func applyNMCLIStatic(ctx context.Context, opts ChangeOptions) *apperror.AppError {
	args := buildNMCLIStaticArgs(opts)
	if err := runNMCLI(ctx, args...); err != nil {
		return err
	}

	return upNMCLIConnection(ctx, opts.InterfaceName)
}

func buildNMCLIStaticArgs(opts ChangeOptions) []string {
	cidr := FormatCIDR(opts.IP, opts.Netmask)
	args := []string{"connection", "modify", opts.InterfaceName, "ipv4.addresses", cidr, "ipv4.method", "manual"}
	if opts.Gateway != "" {
		args = append(args, "ipv4.gateway", opts.Gateway)
	}
	if len(opts.DNS) > 0 {
		args = append(args, "ipv4.dns", strings.Join(opts.DNS, " "))
	}

	return args
}

func upNMCLIConnection(ctx context.Context, iface string) *apperror.AppError {
	return runNMCLI(ctx, "connection", "up", iface)
}

func runNMCLI(ctx context.Context, args ...string) *apperror.AppError {
	cmd := exec.CommandContext(ctx, "nmcli", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		msg := fmt.Sprintf("nmcli %s failed: %s", strings.Join(args, " "), string(out))
		return apperror.WrapWithDetails(err, "runNMCLI", "E_NMCLI_EXEC", msg, "driver_nmcli", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}
