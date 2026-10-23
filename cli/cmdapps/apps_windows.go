package cmdapps

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

func scanWindowsApps(opts ListOptions) ([]InstalledApp, error) {
	apps := make([]InstalledApp, 0, 16)

	// Scan global npm packages
	npmApps, err := scanNpmGlobalPackages()
	if err == nil {
		apps = append(apps, npmApps...)
	}

	return apps, nil
}

func readPackageVersion(pkgJsonPath string) string {
	data, readErr := os.ReadFile(pkgJsonPath)
	if readErr != nil {
		return ""
	}
	var parsed struct {
		Version string `json:"version"`
	}
	if json.Unmarshal(data, &parsed) == nil {
		return parsed.Version
	}

	return ""
}

func scanNpmGlobalPackages() ([]InstalledApp, error) {
	apps := make([]InstalledApp, 0, 4)

	appData := os.Getenv("APPDATA")
	if appData == "" {
		return apps, nil
	}

	npmModules := filepath.Join(appData, "npm", "node_modules")
	entries, err := os.ReadDir(npmModules)
	if err != nil {
		return apps, nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkgName := entry.Name()
		pkgJsonPath := filepath.Join(npmModules, pkgName, "package.json")
		version := readPackageVersion(pkgJsonPath)

		app := InstalledApp{
			ID:          pkgName,
			Name:        pkgName,
			Exec:        filepath.Join(appData, "npm", pkgName+".cmd"),
			Package:     pkgName,
			Version:     version,
			Manager:     ManagerNpm,
			Scope:       ScopeUser,
			IsRemovable: true,
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func executeWindowsUninstall(target string, opts UninstallOptions, resp AppUninstallResponse) (AppUninstallResponse, error) {
	// Try npm global uninstall first
	npmCmd := exec.Command("npm", "uninstall", "-g", target)
	_ = npmCmd.Run()

	appData := os.Getenv("APPDATA")
	if appData == "" {
		return resp, nil
	}

	npmDir := filepath.Join(appData, "npm")
	shims := []string{
		filepath.Join(npmDir, target),
		filepath.Join(npmDir, target+".cmd"),
		filepath.Join(npmDir, target+".ps1"),
	}
	for _, shim := range shims {
		if _, statErr := os.Stat(shim); statErr == nil {
			_ = os.Remove(shim)
			resp.RemovedFiles = append(resp.RemovedFiles, shim)
		}
	}
	nodeModulePath := filepath.Join(npmDir, "node_modules", target)
	if _, statErr := os.Stat(nodeModulePath); statErr == nil {
		_ = os.RemoveAll(nodeModulePath)
		resp.RemovedFiles = append(resp.RemovedFiles, nodeModulePath)
	}

	return resp, nil
}
