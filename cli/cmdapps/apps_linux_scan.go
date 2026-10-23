package cmdapps

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func parseDesktopProps(scanner *bufio.Scanner) map[string]string {
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
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if _, hasKey := props[key]; !hasKey {
			props[key] = val
		}
	}

	return props
}

func resolvePackageManager(path, cleanExec string) (PackageManager, string) {
	if strings.Contains(path, "snap") {
		return ManagerSnap, ""
	}
	if strings.Contains(path, "flatpak") {
		return ManagerFlatpak, ""
	}
	pkg, isTracked := traceDpkgPackage(cleanExec)
	if isTracked {
		return ManagerApt, pkg
	}

	return ManagerCustom, ""
}

func parseDesktopFile(path string, isIncludeAll bool) (InstalledApp, bool) {
	file, err := os.Open(path)
	if err != nil {
		return InstalledApp{}, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	props := parseDesktopProps(scanner)

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

	manager, packageName := resolvePackageManager(path, cleanExec)

	version := ""
	if packageName != "" {
		version = getDpkgVersion(packageName)
	}

	app := InstalledApp{
		ID:          appID,
		Name:        appName,
		Exec:        cleanExec,
		DesktopFile: path,
		Package:     packageName,
		Version:     version,
		Manager:     manager,
		Scope:       scope,
		Icon:        props["Icon"],
		IsRemovable: (manager != ManagerCustom && manager != ""),
	}

	return app, true
}

func cleanExecString(execStr string) string {
	fields := strings.Fields(execStr)
	cleaned := make([]string, 0, len(fields))
	for _, f := range fields {
		if strings.HasPrefix(f, "%") {
			continue
		}
		cleaned = append(cleaned, f)
	}
	return strings.Join(cleaned, " ")
}

func resolveBinaryPath(execPath string) string {
	binary := strings.Fields(execPath)[0]
	if filepath.IsAbs(binary) {
		return binary
	}
	absPath, err := exec.LookPath(binary)
	if err == nil {
		return absPath
	}

	return binary
}

func traceDpkgPackage(execPath string) (string, bool) {
	if execPath == "" {
		return "", false
	}
	binary := resolveBinaryPath(execPath)
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
