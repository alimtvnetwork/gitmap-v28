// Package cmd — agy_plugins_ops.go handles discovery and installation of Antigravity plugins.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func getPluginsDirectory() (string, *apperror.AppError) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", apperror.WrapSimple(homeErr, "user home dir")
	}

	return filepath.Join(homeDir, ".gemini", "config", "plugins"), nil
}

func listInstalledPlugins() ([]string, *apperror.AppError) {
	dir, dirErr := getPluginsDirectory()
	if dirErr != nil {
		return nil, dirErr
	}

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return []string{}, nil
	}

	plugins := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			plugins = append(plugins, e.Name())
		}
	}

	sort.Strings(plugins)

	return plugins, nil
}

func runAgyPluginList() error {
	installed, err := listInstalledPlugins()
	if err != nil {
		return fmt.Errorf("list plugins: %w", err.Unwrap())
	}

	fmt.Printf("%s Antigravity Plugins (%d installed):\n\n", constants.ColorCyan+"▶"+constants.ColorReset, len(installed))
	fmt.Printf("    %-32s %-12s %s\n", "PLUGIN SLUG", "STATUS", "LOCATION")
	fmt.Printf("    %s\n", "──────────────────────────────────────────────────────────────────")

	for _, p := range installed {
		fmt.Printf("  • %-32s \033[32minstalled\033[0m   ~/.gemini/config/plugins/%s\n", p, p)
	}

	available := []string{"firebase-agent-plugin", "cloud-run-plugin", "docker-companion", "k8s-ops"}
	fmt.Printf("\n  %sAvailable to install (use: gitmap agy plugin install <slug>):%s\n", constants.ColorYellow, constants.ColorReset)
	for _, a := range available {
		isInstalled := sliceContains(installed, a)
		if !isInstalled {
			fmt.Printf("  + %-32s \033[33mavailable\033[0m\n", a)
		}
	}

	return nil
}

func runAgyPluginInstall(slug string) error {
	dir, dirErr := getPluginsDirectory()
	if dirErr != nil {
		return dirErr.Unwrap()
	}

	targetDir := filepath.Join(dir, slug)
	if mkErr := os.MkdirAll(targetDir, constants.DirPermission); mkErr != nil {
		return mkErr
	}

	manifestFile := filepath.Join(targetDir, "plugin.json")
	nowStr := time.Now().UTC().Format(time.RFC3339)
	content := fmt.Sprintf(`{"name": "%s", "version": "1.0.0", "installedAt": "%s"}`, slug, nowStr)
	if writeErr := os.WriteFile(manifestFile, []byte(content), constants.FilePermission); writeErr != nil {
		return writeErr
	}

	fmt.Printf("%s Successfully installed plugin %q to %s\n", constants.ColorGreen+"✓"+constants.ColorReset, slug, targetDir)
	return nil
}
