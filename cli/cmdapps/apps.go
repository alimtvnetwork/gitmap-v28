package cmdapps

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ListApps scans and discovers installed applications across the system.
func ListApps(opts ListOptions) (AppListResponse, error) {
	apps := make([]InstalledApp, 0, 32)

	if runtime.GOOS == "linux" {
		linuxApps, err := scanLinuxApps(opts)
		if err == nil {
			apps = append(apps, linuxApps...)
		}
	} else if runtime.GOOS == "windows" {
		windowsApps, err := scanWindowsApps(opts)
		if err == nil {
			apps = append(apps, windowsApps...)
		}
	}

	filtered := filterApps(apps, opts)

	resp := AppListResponse{
		Success:  true,
		Count:    len(filtered),
		Platform: runtime.GOOS,
		Data:     filtered,
	}

	return resp, nil
}

// UninstallApp performs surgical application uninstallation and cache synchronization.
func UninstallApp(target string, opts UninstallOptions) (AppUninstallResponse, error) {
	startTime := time.Now()
	resp := AppUninstallResponse{
		Success:      false,
		AppID:        target,
		Purged:       opts.IsPurge,
		RemovedFiles: make([]string, 0),
		CachesReset:  make([]string, 0),
	}

	if target == "" {
		resp.Error = "target application identifier cannot be empty"
		return resp, fmt.Errorf("%s", resp.Error)
	}

	if opts.IsDryRun {
		resp.Success = true
		resp.DurationMs = time.Since(startTime).Milliseconds()
		return resp, nil
	}

	var hasExecutionErr error
	if runtime.GOOS == "linux" {
		resp, hasExecutionErr = executeLinuxUninstall(target, opts, resp)
	} else if runtime.GOOS == "windows" {
		resp, hasExecutionErr = executeWindowsUninstall(target, opts, resp)
	} else {
		resp.Error = fmt.Sprintf("unsupported platform for uninstaller: %s", runtime.GOOS)
		return resp, fmt.Errorf("%s", resp.Error)
	}

	resp.DurationMs = time.Since(startTime).Milliseconds()
	if hasExecutionErr != nil {
		resp.Success = false
		resp.Error = hasExecutionErr.Error()
		return resp, hasExecutionErr
	}

	resp.Success = true
	return resp, nil
}

func filterApps(apps []InstalledApp, opts ListOptions) []InstalledApp {
	filtered := make([]InstalledApp, 0, len(apps))
	normalizedFilter := strings.ToLower(strings.TrimSpace(opts.FilterStr))

	for _, app := range apps {
		if opts.IsSystem && app.Scope != ScopeSystem {
			continue
		}
		if opts.IsUser && app.Scope != ScopeUser {
			continue
		}

		if normalizedFilter != "" {
			hasIdMatch := strings.Contains(strings.ToLower(app.ID), normalizedFilter)
			hasNameMatch := strings.Contains(strings.ToLower(app.Name), normalizedFilter)
			hasPkgMatch := strings.Contains(strings.ToLower(app.Package), normalizedFilter)
			hasExecMatch := strings.Contains(strings.ToLower(app.Exec), normalizedFilter)
			if !hasIdMatch && !hasNameMatch && !hasPkgMatch && !hasExecMatch {
				continue
			}
		}

		filtered = append(filtered, app)
	}

	return filtered
}

func scanLinuxApps(opts ListOptions) ([]InstalledApp, error) {
	appDirs := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		"/var/lib/snapd/desktop/applications",
		"/var/lib/flatpak/exports/share/applications",
	}

	homeDir, hasHomeErr := os.UserHomeDir()
	if hasHomeErr == nil {
		appDirs = append(appDirs,
			filepath.Join(homeDir, ".local", "share", "applications"),
			filepath.Join(homeDir, ".local", "share", "flatpak", "exports", "share", "applications"),
		)
	}

	seenIDs := make(map[string]bool)
	apps := make([]InstalledApp, 0, 32)

	for _, dir := range appDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".desktop") {
				continue
			}

			fullPath := filepath.Join(dir, entry.Name())
			app, isParsed := parseDesktopFile(fullPath, opts.IsAll)
			if !isParsed {
				continue
			}

			if seenIDs[app.ID] {
				continue
			}
			seenIDs[app.ID] = true
			apps = append(apps, app)
		}
	}

	return apps, nil
}

func parseDesktopFile(path string, isIncludeAll bool) (InstalledApp, bool) {
	file, err := os.Open(path)
	if err != nil {
		return InstalledApp{}, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isInDesktopEntry := false
	props := make(map[string]string)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			isInDesktopEntry = (line == "[Desktop Entry]")
			continue
		}

		if !isInDesktopEntry {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if _, hasKey := props[key]; !hasKey {
				props[key] = val
			}
		}
	}

	appType := props["Type"]
	if appType != "" && appType != "Application" && !isIncludeAll {
		return InstalledApp{}, false
	}

	isNoDisplay := strings.EqualFold(props["NoDisplay"], "true")
	if isNoDisplay && !isIncludeAll {
		return InstalledApp{}, false
	}

	appName := props["Name"]
	if appName == "" {
		appName = strings.TrimSuffix(filepath.Base(path), ".desktop")
	}

	cleanExec := cleanExecString(props["Exec"])
	appID := strings.ToLower(strings.TrimSuffix(filepath.Base(path), ".desktop"))

	scope := ScopeSystem
	if strings.Contains(path, "/.local/") || strings.Contains(path, "/home/") {
		scope = ScopeUser
	}

	manager := ManagerCustom
	packageName := ""
	if strings.Contains(path, "snap") {
		manager = ManagerSnap
	} else if strings.Contains(path, "flatpak") {
		manager = ManagerFlatpak
	} else {
		pkg, isTracked := traceDpkgPackage(cleanExec)
		if isTracked {
			manager = ManagerApt
			packageName = pkg
		}
	}

	version := ""
	if packageName != "" {
		version = getDpkgVersion(packageName)
	}

	app := InstalledApp{
		ID:          appID,
		Name:        appName,
		Exec:        cleanExec,
		Icon:        props["Icon"],
		Package:     packageName,
		Version:     version,
		Manager:     manager,
		Scope:       scope,
		DesktopFile: path,
		IsRemovable: true,
	}

	return app, true
}

func cleanExecString(rawExec string) string {
	if rawExec == "" {
		return ""
	}
	fields := strings.Fields(rawExec)
	cleaned := make([]string, 0, len(fields))
	for _, field := range fields {
		if strings.HasPrefix(field, "%") {
			continue
		}
		cleaned = append(cleaned, field)
	}
	return strings.Join(cleaned, " ")
}

func traceDpkgPackage(execPath string) (string, bool) {
	if execPath == "" {
		return "", false
	}
	binary := strings.Fields(execPath)[0]
	if !filepath.IsAbs(binary) {
		absPath, err := exec.LookPath(binary)
		if err == nil {
			binary = absPath
		}
	}

	out, err := exec.Command("dpkg", "-S", binary).Output()
	if err != nil {
		return "", false
	}

	parts := strings.Split(string(out), ":")
	if len(parts) > 0 {
		pkg := strings.TrimSpace(parts[0])
		return pkg, true
	}

	return "", false
}

func getDpkgVersion(packageName string) string {
	out, err := exec.Command("dpkg-query", "-W", "-f=${Version}", packageName).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
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

	packageName := target
	desktopFile := ""
	if targetApp != nil {
		if targetApp.Package != "" {
			packageName = targetApp.Package
		}
		desktopFile = targetApp.DesktopFile
		resp.Package = packageName
	}

	// Purge or remove package via apt
	isDpkgPackage := isPackageInstalledApt(packageName)
	if isDpkgPackage {
		var aptCmd *exec.Cmd
		if opts.IsPurge {
			aptCmd = exec.Command("sudo", "DEBIAN_FRONTEND=noninteractive", "apt-get", "purge", "-y", packageName)
		} else {
			aptCmd = exec.Command("sudo", "DEBIAN_FRONTEND=noninteractive", "apt-get", "remove", "-y", packageName)
		}
		output, err := aptCmd.CombinedOutput()
		if err != nil {
			// Try fallback without sudo if already root
			if os.Geteuid() == 0 {
				var rootCmd *exec.Cmd
				if opts.IsPurge {
					rootCmd = exec.Command("apt-get", "purge", "-y", packageName)
				} else {
					rootCmd = exec.Command("apt-get", "remove", "-y", packageName)
				}
				_ = rootCmd.Run()
			} else {
				return resp, fmt.Errorf("apt uninstallation failed for %s: %v (%s)", packageName, err, strings.TrimSpace(string(output)))
			}
		}
	}

	// Clean desktop file if present
	possibleDesktopFiles := []string{
		desktopFile,
		filepath.Join("/usr/share/applications", target+".desktop"),
		filepath.Join("/usr/share/applications", strings.ToLower(target)+".desktop"),
		filepath.Join("/usr/share/applications", strings.Title(strings.ReplaceAll(target, "-", " "))+".desktop"),
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
		if _, err := os.Stat(dPath); err == nil {
			rmErr := os.Remove(dPath)
			if rmErr != nil {
				_ = exec.Command("sudo", "rm", "-f", dPath).Run()
			}
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

func scanWindowsApps(opts ListOptions) ([]InstalledApp, error) {
	apps := make([]InstalledApp, 0, 16)

	// Scan global npm packages
	npmApps, err := scanNpmGlobalPackages()
	if err == nil {
		apps = append(apps, npmApps...)
	}

	return apps, nil
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
		version := ""
		data, readErr := os.ReadFile(pkgJsonPath)
		if readErr == nil {
			var parsed struct {
				Version string `json:"version"`
			}
			if json.Unmarshal(data, &parsed) == nil {
				version = parsed.Version
			}
		}

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
	if appData != "" {
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
	}

	return resp, nil
}
