package cmdinstall

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var devToolAliases = map[string]string{
	"powershell": "pwsh", "pwsh": "pwsh",
	"agm": "agm", "ag-manager": "agm", "antigravity-manager": "agm",
	"vim": "vim", "vi": "vim",
	"code": "vscode", "vs-code": "vscode", "vscode": "vscode",
	"cursor": "cursor", "cur": "cursor",
}

func normalizeDevToolName(tool string) string {
	t := strings.ToLower(strings.TrimSpace(tool))
	if canonical, isFound := devToolAliases[t]; isFound {
		return canonical
	}
	return t
}

// RunToolUninstall dispatches uninstallation for supported dev tools.
func RunToolUninstall(tool string, isDryRun, isForce, isPurge bool) error {
	switch normalizeDevToolName(tool) {
	case "pwsh":
		return uninstallPwsh(isDryRun, isPurge)
	case "agm":
		return uninstallAgm(isDryRun, isPurge)
	case "vim":
		return uninstallVim(isDryRun, isPurge)
	case "vscode":
		return uninstallVSCode(isDryRun, isPurge)
	case "cursor":
		return uninstallCursor(isDryRun, isPurge)
	default:
		return fmt.Errorf("unsupported dev tool for uninstallation: %s", tool)
	}
}

func resolveAptUninstallCmd(pkg string, isPurge bool) string {
	action := "remove"
	if isPurge {
		action = "purge"
	}
	return fmt.Sprintf("sudo apt %s %s -y", action, pkg)
}

func uninstallPwsh(isDryRun, isPurge bool) error {
	cmd := resolveAptUninstallCmd("powershell", isPurge)
	plan := UninstallPlan{Tool: "pwsh", PackageRemoval: cmd, BinaryWrappers: []string{"/usr/local/bin/pwsh"}}
	if isPurge {
		plan.Configuration = resolveConfigPath(".config", "powershell")
	}
	if isDryRun {
		renderUninstallDryRun("pwsh", plan)
		return nil
	}
	_ = runSystemCmd(cmd)
	return removePathTargets(plan.BinaryWrappers)
}

func resolveAgmLocations() ([]string, []string) {
	wrappers := []string{"/usr/local/bin/agm"}
	launchers := []string{"/usr/share/applications/agm.desktop"}
	if home, err := os.UserHomeDir(); err == nil {
		wrappers = append(wrappers, filepath.Join(home, ".local", "bin", "agm"))
		launchers = append(launchers, filepath.Join(home, ".local", "share", "applications", "agm.desktop"))
	}
	return wrappers, launchers
}

func uninstallAgm(isDryRun, isPurge bool) error {
	wrappers, launchers := resolveAgmLocations()
	plan := UninstallPlan{Tool: "agm", BinaryWrappers: wrappers, DesktopLaunchers: launchers}
	if isPurge {
		plan.Configuration = resolveConfigPath(".config", "antigravity")
	}
	if isDryRun {
		renderUninstallDryRun("agm", plan)
		return nil
	}
	return removePathTargets(append(wrappers, launchers...))
}

func resolveVimConfigPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{filepath.Join(home, ".vim"), filepath.Join(home, ".vimrc")}
}

func uninstallVim(isDryRun, isPurge bool) error {
	cmd := resolveAptUninstallCmd("vim", isPurge)
	plan := UninstallPlan{Tool: "vim", PackageRemoval: cmd}
	if isPurge {
		plan.Configuration = resolveVimConfigPaths()
	}
	if isDryRun {
		renderUninstallDryRun("vim", plan)
		return nil
	}
	_ = runSystemCmd(cmd)
	return removePathTargets(plan.Configuration)
}

func resolveVSCodeLaunchers() []string {
	launchers := []string{"/usr/share/applications/code.desktop"}
	if home, err := os.UserHomeDir(); err == nil {
		launchers = append(launchers, filepath.Join(home, ".local", "share", "applications", "code.desktop"))
	}
	return launchers
}

func uninstallVSCode(isDryRun, isPurge bool) error {
	cmd := resolveAptUninstallCmd("code", isPurge)
	launchers := resolveVSCodeLaunchers()
	plan := UninstallPlan{Tool: "vscode", PackageRemoval: cmd, BinaryWrappers: []string{"/usr/local/bin/code"}, DesktopLaunchers: launchers}
	if isPurge {
		plan.Configuration = resolveConfigPath(".config", "Code")
	}
	if isDryRun {
		renderUninstallDryRun("vscode", plan)
		return nil
	}
	_ = runSystemCmd(cmd)
	return removePathTargets(append(plan.BinaryWrappers, launchers...))
}

func resolveCursorLocations() ([]string, []string, []string) {
	targets := []string{"/opt/cursor"}
	wrappers := []string{"/usr/local/bin/cursor"}
	launchers := []string{"/usr/share/applications/cursor.desktop"}
	if home, err := os.UserHomeDir(); err == nil {
		wrappers = append(wrappers, filepath.Join(home, ".local", "bin", "cursor"))
		launchers = append(launchers, filepath.Join(home, ".local", "share", "applications", "cursor.desktop"))
	}
	return targets, wrappers, launchers
}

func uninstallCursor(isDryRun, isPurge bool) error {
	targets, wrappers, launchers := resolveCursorLocations()
	plan := UninstallPlan{Tool: "cursor", PackageRemoval: "rm -rf /opt/cursor", BinaryWrappers: wrappers, DesktopLaunchers: launchers, DesktopDock: "unpin cursor.desktop from org.gnome.shell favorite-apps"}
	if isPurge {
		plan.Configuration = resolveConfigPath(".config", "Cursor")
	}
	if isDryRun {
		renderUninstallDryRun("cursor", plan)
		return nil
	}
	_ = unpinCursorGnomeDock()
	return removePathTargets(append(targets, append(wrappers, launchers...)...))
}
