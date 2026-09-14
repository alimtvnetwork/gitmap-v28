package netip

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DefaultNetplanConfigFile defines the standard Netplan drop-in configuration path.
const DefaultNetplanConfigFile = "/etc/netplan/99-gitmap-config.yaml"

// LinuxNetplanDriver manages network configuration via Netplan YAML.
type LinuxNetplanDriver struct {
	configPath string
}

// NewLinuxNetplanDriver creates a new Netplan driver targeting default config path.
func NewLinuxNetplanDriver() *LinuxNetplanDriver {
	return &LinuxNetplanDriver{configPath: DefaultNetplanConfigFile}
}

// Name returns the driver identifier.
func (d *LinuxNetplanDriver) Name() string {
	return DriverNameLinuxNetplan
}

// ListInterfaces lists system network adapters.
func (d *LinuxNetplanDriver) ListInterfaces(ctx context.Context) InterfaceSliceResult {
	return ListSystemInterfaces()
}

// GetInterface retrieves a specific interface.
func (d *LinuxNetplanDriver) GetInterface(ctx context.Context, name string) InterfaceResult {
	return GetSystemInterface(name)
}

// ApplyConfig writes Netplan YAML and applies changes.
func (d *LinuxNetplanDriver) ApplyConfig(ctx context.Context, opts ChangeOptions) *apperror.AppError {
	yamlContent := generateNetplanYAML(opts)
	if err := writeNetplanConfig(d.configPath, yamlContent); err != nil {
		return err
	}

	return runNetplan(ctx, "apply")
}

// TryConfig tests Netplan configuration with auto-rollback timeout.
func (d *LinuxNetplanDriver) TryConfig(ctx context.Context, timeoutSec int) *apperror.AppError {
	timeoutStr := fmt.Sprintf("--timeout=%d", timeoutSec)

	return runNetplan(ctx, "try", timeoutStr)
}

// RevertConfig restores previous configuration from snapshot.
func (d *LinuxNetplanDriver) RevertConfig(ctx context.Context, snap Snapshot) *apperror.AppError {
	opts := snapshotToOptions(snap)

	return d.ApplyConfig(ctx, opts)
}

func generateNetplanYAML(opts ChangeOptions) string {
	if opts.IsDHCP {
		return renderDHCP(opts.InterfaceName)
	}

	return renderStaticEthernet(opts)
}

func renderDHCP(iface string) string {
	return fmt.Sprintf("network:\n  version: 2\n  renderer: networkd\n  ethernets:\n    %s:\n      dhcp4: true\n", iface)
}

func renderStaticEthernet(opts ChangeOptions) string {
	cidr := FormatCIDR(opts.IP, opts.Netmask)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("network:\n  version: 2\n  renderer: networkd\n  ethernets:\n    %s:\n      dhcp4: false\n      addresses:\n        - %s\n", opts.InterfaceName, cidr))
	if opts.Gateway != "" {
		sb.WriteString(renderRoutes(opts.Gateway))
	}
	if len(opts.DNS) > 0 {
		sb.WriteString(renderDNS(opts.DNS))
	}

	return sb.String()
}

func renderRoutes(gateway string) string {
	return fmt.Sprintf("      routes:\n        - to: default\n          via: %s\n", gateway)
}

func renderDNS(dnsList []string) string {
	joined := strings.Join(dnsList, ", ")

	return fmt.Sprintf("      nameservers:\n        addresses: [%s]\n", joined)
}

func writeNetplanConfig(path, yamlStr string) *apperror.AppError {
	if isNetIPTestMode() {
		return nil
	}

	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0o755)

	if err := os.WriteFile(path, []byte(yamlStr), 0o600); err != nil {
		return apperror.WrapSimple(err, "writeNetplanConfig")
	}

	return nil
}

func runNetplan(ctx context.Context, args ...string) *apperror.AppError {
	if isNetIPTestMode() {
		return nil
	}

	cmd := exec.CommandContext(ctx, "netplan", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		msg := fmt.Sprintf("netplan %s failed: %s", strings.Join(args, " "), string(out))
		return apperror.WrapWithDetails(err, "runNetplan", "E_NETPLAN_EXEC", msg, "driver_netplan", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}
