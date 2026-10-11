package completion

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	powerShellBinLegacy       = "powershell"
	powerShellBinCore         = "pwsh"
	powerShellDirCore         = "PowerShell"
	powerShellDirLegacy       = "WindowsPowerShell"
	powerShellProfileAllHosts = "profile.ps1"
	powerShellProfileCurrent  = "Microsoft.PowerShell_profile.ps1"
	powerShellProfileProbeCmd = "$PROFILE.CurrentUserAllHosts; $PROFILE.CurrentUserCurrentHost"
)

// resolvePowerShellPaths returns paths for PowerShell completion.
func resolvePowerShellPaths() (string, string) {
	appData := os.Getenv("APPDATA")
	if len(appData) == 0 {
		home, _ := os.UserHomeDir()
		appData = filepath.Join(home, ".config")
	}

	scriptPath := filepath.Join(appData, constants.CompDirName, constants.CompFilePS)
	profile := defaultPSProfile()
	paths := resolvePowerShellProfilePaths()
	if len(paths) > 0 {
		profile = paths[0]
	}

	return scriptPath, profile
}

// defaultPSProfile returns the default PowerShell profile path.
func defaultPSProfile() string {
	home, _ := os.UserHomeDir()
	paths := defaultPowerShellProfilePaths(home, runtime.GOOS)
	if len(paths) > 0 {
		return paths[0]
	}

	return filepath.Join(home, ".config", "powershell", powerShellProfileAllHosts)
}

// resolvePowerShellProfilePaths returns all profile targets that should receive gitmap shell integration.
func resolvePowerShellProfilePaths() []string {
	paths := make([]string, 0, 7)
	if profile := strings.TrimSpace(os.Getenv("PROFILE")); len(profile) > 0 {
		paths = append(paths, profile)
	}

	paths = append(paths, probePowerShellProfilePaths()...)
	home, _ := os.UserHomeDir()
	paths = append(paths, defaultPowerShellProfilePaths(home, runtime.GOOS)...)

	return uniqueProfilePaths(paths)
}

// probePowerShellProfilePaths asks available PowerShell engines for their current-user profile paths.
func probePowerShellProfilePaths() []string {
	bins := []string{powerShellBinLegacy, powerShellBinCore}
	paths := make([]string, 0, 4)
	for _, bin := range bins {
		out, err := exec.Command(bin, "-NoProfile", "-Command", powerShellProfileProbeCmd).Output()
		if err != nil {
			continue
		}

		paths = append(paths, parsePowerShellProfileOutput(string(out))...)
	}

	return paths
}

// parsePowerShellProfileOutput converts PowerShell probe output into profile paths.
func parsePowerShellProfileOutput(output string) []string {
	lines := strings.Split(output, "\n")
	paths := make([]string, 0, len(lines))
	for _, line := range lines {
		path := strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if len(path) > 0 {
			paths = append(paths, path)
		}
	}

	return paths
}

// defaultPowerShellProfilePaths returns the standard current-user profile files for the OS.
func defaultPowerShellProfilePaths(home, goos string) []string {
	if len(home) == 0 {
		return nil
	}

	if goos == "windows" {
		docs := filepath.Join(home, "Documents")

		return []string{
			filepath.Join(docs, powerShellDirCore, powerShellProfileAllHosts),
			filepath.Join(docs, powerShellDirCore, powerShellProfileCurrent),
			filepath.Join(docs, powerShellDirLegacy, powerShellProfileAllHosts),
			filepath.Join(docs, powerShellDirLegacy, powerShellProfileCurrent),
		}
	}

	base := filepath.Join(home, ".config", "powershell")

	return []string{
		filepath.Join(base, powerShellProfileAllHosts),
		filepath.Join(base, powerShellProfileCurrent),
	}
}

// uniqueProfilePaths removes empty and duplicate profile entries while preserving order.
func uniqueProfilePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	unique := make([]string, 0, len(paths))
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if len(path) == 0 {
			continue
		}

		if _, ok := seen[path]; ok {
			continue
		}

		seen[path] = struct{}{}
		unique = append(unique, path)
	}

	return unique
}
