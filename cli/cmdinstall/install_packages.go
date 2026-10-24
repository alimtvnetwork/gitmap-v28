package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// UninstallPlan encapsulates package and directory targets for uninstallation.
type UninstallPlan struct {
	Tool             string
	PackageRemoval   string
	BinaryWrappers   []string
	DesktopLaunchers []string
	SystemIcons      []string
	DesktopDock      string
	Configuration    []string
}

func renderUninstallDryRun(tool string, plan UninstallPlan) {
	fmt.Printf("● [DryRun] Simulated Uninstallation for '%s' (OS: %s):\n", tool, runtime.GOOS)
	printPlanPackage(plan.PackageRemoval)
	renderDryRunPaths("- Binary Wrappers", plan.BinaryWrappers)
	renderDryRunPaths("- Desktop Launchers", plan.DesktopLaunchers)
	renderDryRunPaths("- System Icons", plan.SystemIcons)
	if plan.DesktopDock != "" {
		fmt.Printf("  - Desktop Dock: %s\n", plan.DesktopDock)
	}
	renderDryRunPaths("- Configuration", plan.Configuration)
	fmt.Println("✔ Dry-run completed. No packages or files were modified.")
}

func printPlanPackage(pkgCmd string) {
	if pkgCmd != "" {
		fmt.Printf("  - Package Removal: %s\n", pkgCmd)
		return
	}
	fmt.Println("  - Package Removal: N/A")
}

func renderDryRunPaths(label string, paths []string) {
	if len(paths) == 0 {
		return
	}
	fmt.Printf("  %s:\n", label)
	for _, p := range paths {
		fmt.Printf("    • %s\n", p)
	}
}

func isInsideCwd(norm string) bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}
	cwdNorm := strings.ToLower(filepath.ToSlash(filepath.Clean(cwd)))
	return norm == cwdNorm || strings.HasPrefix(norm, cwdNorm+"/")
}

func isProtectedWorkDirectory(targetPath string) bool {
	norm := strings.ToLower(filepath.ToSlash(filepath.Clean(targetPath)))
	if norm == "" || norm == "." || norm == ".." {
		return true
	}
	if strings.Contains(norm, "/git-work") || strings.Contains(norm, "/work") {
		return true
	}
	return isInsideCwd(norm)
}

func removeSingleTarget(p string) error {
	if isProtectedWorkDirectory(p) {
		return apperror.NewValidationError(fmt.Sprintf("safety guard: refused to modify protected directory %s", p))
	}
	_ = os.RemoveAll(p)
	return nil
}

func removePathTargets(paths []string) error {
	for _, p := range paths {
		if err := removeSingleTarget(p); err != nil {
			return err
		}
	}
	return nil
}

func runSystemCmd(cmdStr string) error {
	if constants.IsWindowsPlatform() {
		return nil
	}
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, "command execution failed").WithContext("output", strings.TrimSpace(string(out)))
	}
	return nil
}

func unpinCursorGnomeDock() error {
	if !constants.IsLinuxPlatform() {
		return nil
	}
	out, err := exec.Command("gsettings", "get", "org.gnome.shell", "favorite-apps").Output()
	if err != nil {
		return nil
	}
	cur := string(out)
	if !strings.Contains(cur, "'cursor.desktop'") {
		return nil
	}
	updated := strings.Replace(cur, "'cursor.desktop', ", "", 1)
	updated = strings.Replace(updated, ", 'cursor.desktop'", "", 1)
	updated = strings.Replace(updated, "'cursor.desktop'", "", 1)
	return exec.Command("gsettings", "set", "org.gnome.shell", "favorite-apps", strings.TrimSpace(updated)).Run()
}

func resolveConfigPath(subdirs ...string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{filepath.Join(append([]string{home}, subdirs...)...)}
}
