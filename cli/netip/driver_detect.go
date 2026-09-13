package netip

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ResolveDriverByName returns a driver by explicit name.
func ResolveDriverByName(name string) (Driver, *apperror.AppError) {
	switch name {
	case DriverNameWindows:
		return NewWindowsDriver(), nil
	case DriverNameLinuxNetplan:
		return NewLinuxNetplanDriver(), nil
	case DriverNameLinuxNMCLI:
		return NewLinuxNMCLIDriver(), nil
	case DriverNameLinuxIPRoute2:
		return NewLinuxIPRoute2Driver(), nil
	default:
		return nil, apperror.New("ResolveDriverByName", "E_UNKNOWN_DRIVER", map[string]any{"driver": name})
	}
}

// DetectDriver auto-detects the appropriate driver for current OS environment.
func DetectDriver() (Driver, *apperror.AppError) {
	if runtime.GOOS == "windows" {
		return NewWindowsDriver(), nil
	}

	if runtime.GOOS == "linux" {
		return detectLinuxDriver(), nil
	}

	return nil, apperror.New("DetectDriver", "E_UNSUPPORTED_OS", map[string]any{"os": runtime.GOOS})
}

func detectLinuxDriver() Driver {
	if isNetplanAvailable() {
		return NewLinuxNetplanDriver()
	}

	if isNMCLIAvailable() {
		return NewLinuxNMCLIDriver()
	}

	return NewLinuxIPRoute2Driver()
}

func isNetplanAvailable() bool {
	if _, err := exec.LookPath("netplan"); err == nil {
		return true
	}

	_, statErr := os.Stat("/etc/netplan")

	return statErr == nil
}

func isNMCLIAvailable() bool {
	_, err := exec.LookPath("nmcli")

	return err == nil
}
