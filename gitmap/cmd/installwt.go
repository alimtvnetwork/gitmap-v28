package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// runWTSettingsOnly syncs Windows Terminal settings from the bundled folder.
func runWTSettingsOnly() {
	fmt.Println("  Syncing Windows Terminal settings...")
	syncWTSettings()
}

// syncWTSettings copies Windows Terminal settings to the LocalState directory.
func syncWTSettings() {
	target := wtSettingsTarget()
	if target == "" {
		return
	}

	fmt.Printf("  -> Windows Terminal settings target: %s\n", target)

	sourcePath := resolveSettingsPath(
		filepath.Join("settings", "04 - windows-terminal"),
		filepath.Join("data", "wt-settings"),
	)

	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: Windows Terminal settings source not found at %s: %v\n", sourcePath, err)

		return
	}

	err = os.MkdirAll(target, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to create directory %s: %v\n", target, err)

		return
	}

	copied := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := strings.ToLower(entry.Name())
		if name == "readme.txt" || name == "readme.md" {
			continue
		}

		src := filepath.Join(sourcePath, entry.Name())
		dst := filepath.Join(target, entry.Name())

		data, readErr := os.ReadFile(src)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "  Error: failed to read %s: %v\n", src, readErr)

			continue
		}

		writeErr := os.WriteFile(dst, data, 0o644)
		if writeErr != nil {
			fmt.Fprintf(os.Stderr, "  Error: failed to write %s: %v\n", dst, writeErr)

			continue
		}

		copied++
	}

	if copied > 0 {
		fmt.Printf("  -> Synced %d file(s) to %s\n", copied, target)
	} else {
		fmt.Println("  -> No settings files found to sync. Add your settings.json to the settings/04 - windows-terminal/ folder.")
	}
}

// wtSettingsTarget finds the Windows Terminal LocalState directory.
func wtSettingsTarget() string {
	if runtime.GOOS != "windows" {
		fmt.Fprintf(os.Stderr, "  Error: Windows Terminal settings sync is only supported on Windows (current OS: %s)\n", runtime.GOOS)

		return ""
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		fmt.Fprint(os.Stderr, "  Error: LOCALAPPDATA environment variable not set\n")

		return ""
	}

	// Search for Microsoft.WindowsTerminal_* package folder.
	packagesDir := filepath.Join(localAppData, "Packages")

	entries, err := os.ReadDir(packagesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: could not read Packages directory: %v\n", err)

		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "Microsoft.WindowsTerminal_") {
			target := filepath.Join(packagesDir, entry.Name(), "LocalState")

			return target
		}
	}

	fmt.Fprint(os.Stderr, "  Error: Windows Terminal package not found in LocalAppData\\Packages\n")

	return ""
}
