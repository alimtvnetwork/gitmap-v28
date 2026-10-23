package cmdapps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolveTargetPackageInfo(targetApp *InstalledApp, defaultPkg string) (string, string) {
	if targetApp == nil {
		return defaultPkg, ""
	}
	pkg := defaultPkg
	if targetApp.Package != "" {
		pkg = targetApp.Package
	}

	return pkg, targetApp.DesktopFile
}

func runAptUninstall(packageName string, isPurge bool) error {
	action := "remove"
	if isPurge {
		action = "purge"
	}
	aptCmd := exec.Command("sudo", "DEBIAN_FRONTEND=noninteractive", "apt-get", action, "-y", packageName)
	output, err := aptCmd.CombinedOutput()
	if err == nil {
		return nil
	}
	if os.Geteuid() != 0 {
		return apperror.NewExecutionError(fmt.Sprintf("apt uninstallation failed for %s: %v (%s)", packageName, err, strings.TrimSpace(string(output))))
	}
	rootCmd := exec.Command("apt-get", action, "-y", packageName)

	return rootCmd.Run()
}

func tryAptUninstallIfInstalled(packageName string, isPurge bool) error {
	if !isPackageInstalledApt(packageName) {
		return nil
	}

	return runAptUninstall(packageName, isPurge)
}

func removeDesktopFile(dPath string) bool {
	if _, err := os.Stat(dPath); err != nil {
		return false
	}
	rmErr := os.Remove(dPath)
	if rmErr != nil {
		_ = exec.Command("sudo", "rm", "-f", dPath).Run()
	}

	return true
}

func toTitleCase(s string) string {
	if s == "" {
		return ""
	}
	words := strings.Fields(s)
	for i, w := range words {
		r := []rune(w)
		r[0] = unicode.ToUpper(r[0])
		words[i] = string(r)
	}

	return strings.Join(words, " ")
}

func executeLinuxUninstall(target string, opts UninstallOptions, resp AppUninstallResponse) (AppUninstallResponse, error) {
	allApps, _ := scanLinuxApps(ListOptions{IsAll: true})
	var targetApp *InstalledApp
	for i := range allApps {
		if strings.EqualFold(allApps[i].ID, target) ||
			strings.EqualFold(allApps[i].Package, target) ||
			strings.EqualFold(allApps[i].Name, target) ||
			strings.EqualFold(allApps[i].DesktopFile, target) {
			targetApp = &allApps[i]
			break
		}
	}

	packageName, desktopFile := resolveTargetPackageInfo(targetApp, target)
	if targetApp != nil {
		resp.Package = packageName
	}

	if err := tryAptUninstallIfInstalled(packageName, opts.IsPurge); err != nil {
		return resp, err
	}

	// Clean desktop file if present
	possibleDesktopFiles := []string{
		desktopFile,
		filepath.Join("/usr/share/applications", target+".desktop"),
		filepath.Join("/usr/share/applications", strings.ToLower(target)+".desktop"),
		filepath.Join("/usr/share/applications", toTitleCase(strings.ReplaceAll(target, "-", " "))+".desktop"),
	}

	homeDir, hasHomeErr := os.UserHomeDir()
	if hasHomeErr == nil {
		possibleDesktopFiles = append(possibleDesktopFiles,
			filepath.Join(homeDir, ".local", "share", "applications", target+".desktop"),
			filepath.Join(homeDir, ".local", "share", "applications", strings.ToLower(target)+".desktop"),
		)
	}

	for _, dPath := range possibleDesktopFiles {
		if dPath == "" {
			continue
		}
		if removeDesktopFile(dPath) {
			resp.RemovedFiles = append(resp.RemovedFiles, dPath)
		}
	}

	// Refresh GNOME desktop database and icon caches
	_ = exec.Command("sudo", "update-desktop-database", "/usr/share/applications").Run()
	_ = exec.Command("update-desktop-database", filepath.Join(homeDir, ".local", "share", "applications")).Run()
	_ = exec.Command("sudo", "gtk-update-icon-cache", "-f", "-t", "/usr/share/icons/hicolor").Run()
	resp.CachesReset = append(resp.CachesReset, "update-desktop-database", "gtk-update-icon-cache")

	return resp, nil
}

func isPackageInstalledApt(packageName string) bool {
	if packageName == "" {
		return false
	}
	out, err := exec.Command("dpkg", "-s", packageName).Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "Status: install ok installed")
}
