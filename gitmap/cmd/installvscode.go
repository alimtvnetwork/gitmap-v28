package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// runVSCodeSettingsOnly syncs VS Code settings without installing the binary.
func runVSCodeSettingsOnly() {
	fmt.Println("  Syncing VS Code settings (settings-only mode)...")
	syncVSCodeSettings()
}

// syncVSCodeSettings copies VS Code settings from the bundled settings folder.
func syncVSCodeSettings() {
	target := vsCodeSettingsTarget()
	if target == "" {
		return
	}

	fmt.Printf("  -> VS Code settings target: %s\n", target)

	err := os.MkdirAll(target, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to create directory %s: %v\n", target, err)

		return
	}

	// Resolve settings source path.
	sourcePath := resolveSettingsPath(
		filepath.Join("settings", "02 - vscode"),
		filepath.Join("data", "vscode-settings"),
	)

	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: VS Code settings source not found at %s: %v\n", sourcePath, err)

		return
	}

	copied := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.EqualFold(name, "readme.txt") || strings.EqualFold(name, "readme.md") {
			continue
		}

		src := filepath.Join(sourcePath, name)
		dst := filepath.Join(target, name)

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
		fmt.Printf("  -> Synced %d settings file(s) to %s\n", copied, target)
	} else {
		fmt.Println("  -> No settings files found to sync.")
	}

	// Also sync extensions list if present.
	syncVSCodeExtensions(sourcePath)
}

// vsCodeSettingsTarget returns the VS Code user settings directory.
func vsCodeSettingsTarget() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			fmt.Fprint(os.Stderr, "  Error: APPDATA environment variable not set\n")

			return ""
		}

		return filepath.Join(appData, "Code", "User")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error: could not resolve home directory: %v\n", err)

			return ""
		}

		return filepath.Join(home, "Library", "Application Support", "Code", "User")
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error: could not resolve home directory: %v\n", err)

			return ""
		}

		return filepath.Join(home, ".config", "Code", "User")
	}
}

// syncVSCodeExtensions installs extensions from an extensions.txt file.
func syncVSCodeExtensions(settingsDir string) {
	extFile := filepath.Join(settingsDir, "extensions.txt")

	data, err := os.ReadFile(extFile)
	if err != nil {
		return
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 {
		return
	}

	// Check if code CLI is available.
	codePath, err := exec.LookPath("code")
	if err != nil {
		fmt.Println("  -> extensions.txt found but 'code' CLI not in PATH, skipping extension sync.")

		return
	}

	fmt.Printf("  -> Installing %d VS Code extension(s)...\n", len(lines))

	installed := 0

	for _, ext := range lines {
		ext = strings.TrimSpace(ext)
		if ext == "" || strings.HasPrefix(ext, "#") {
			continue
		}

		cmd := exec.Command(codePath, "--install-extension", ext, "--force")
		out, runErr := cmd.CombinedOutput()

		if runErr != nil {
			fmt.Fprintf(os.Stderr, "  ! Failed to install extension %s: %v\n", ext, runErr)

			continue
		}

		_ = out
		installed++
	}

	fmt.Printf("  -> Installed %d/%d extension(s).\n", installed, len(lines))
}
