package completion

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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

// Install writes the completion script and adds a source line to the profile.
func Install(shell string) error {
	script, err := Generate(shell)
	if err != nil {
		return err
	}

	scriptPath, profilePath := resolvePaths(shell)

	return writeAndSource(script, scriptPath, profilePath, shell)
}

// resolvePaths returns the script file path and profile path for the shell.
func resolvePaths(shell string) (string, string) {
	switch shell {
	case constants.ShellPowerShell:
		return resolvePowerShellPaths()
	case constants.ShellBash:
		return resolveBashPaths()
	default:
		return resolveZshPaths()
	}
}

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

// resolveBashPaths returns paths for Bash completion.
func resolveBashPaths() (string, string) {
	home, _ := os.UserHomeDir()
	scriptPath := filepath.Join(home, ".local", "share", constants.CompDirName, constants.CompFileBash)
	profile := filepath.Join(home, ".bashrc")

	return scriptPath, profile
}

// resolveZshPaths returns paths for Zsh completion.
func resolveZshPaths() (string, string) {
	home, _ := os.UserHomeDir()
	scriptPath := filepath.Join(home, ".local", "share", constants.CompDirName, constants.CompFileZsh)
	profile := filepath.Join(home, ".zshrc")

	return scriptPath, profile
}

// writeAndSource writes the script file and adds a source line to the profile.
func writeAndSource(script, scriptPath, profilePath, shell string) error {
	if err := writeScriptFile(scriptPath, script); err != nil {
		return err
	}

	profilePaths := []string{profilePath}
	if shell == constants.ShellPowerShell {
		profilePaths = resolvePowerShellProfilePaths()
	}

	return addSourceLines(scriptPath, profilePaths, shell)
}

// addSourceLines appends the source command to every resolved profile.
func addSourceLines(scriptPath string, profilePaths []string, shell string) error {
	for _, profilePath := range profilePaths {
		if err := addSourceLine(scriptPath, profilePath, shell); err != nil {
			return err
		}
	}

	return nil
}

// writeScriptFile creates directories and writes the completion script.
func writeScriptFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

// addSourceLine appends the source command to the profile if absent.
func addSourceLine(scriptPath, profilePath, shell string) error {
	sourceLine := buildSourceLine(scriptPath, shell)
	if err := os.MkdirAll(filepath.Dir(profilePath), 0o755); err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}

	existing, err := os.ReadFile(profilePath)
	if err == nil && strings.Contains(string(existing), sourceLine) {
		fmt.Fprintf(os.Stderr, constants.MsgCompAlreadyDone, shell)

		return nil
	}

	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\n# gitmap shell completion\n%s\n", sourceLine)
	if err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCompProfileWrite, profilePath)

	return nil
}

// buildSourceLine returns the shell-appropriate source command.
func buildSourceLine(scriptPath, shell string) string {
	if shell == constants.ShellPowerShell {
		return fmt.Sprintf(". '%s'", scriptPath)
	}

	return fmt.Sprintf("source '%s'", scriptPath)
}
