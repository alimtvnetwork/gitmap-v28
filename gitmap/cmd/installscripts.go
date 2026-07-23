package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// scriptsConfig holds the powershell.json structure for deploy path resolution.
type scriptsConfig struct {
	DeployPath string `json:"deployPath"`
}

// runInstallScripts clones/copies gitmap scripts to a platform-specific folder.
func runInstallScripts() {
	targetDir := resolveScriptsDir()

	fmt.Printf(constants.MsgScriptsTarget, targetDir)

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScriptsMkdir, targetDir, err)
		os.Exit(1)
	}

	// Clone the repo into a temp dir, then copy scripts out.
	repoURL := "https://" + constants.GitmapRepoPrefix + ".git"
	fmt.Printf(constants.MsgScriptsCloning, repoURL)

	tmpDir, err := os.MkdirTemp("", "gitmap-scripts-clone-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScriptsTemp, err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	cloneCmd := exec.Command("git", "clone", "--depth", "1", repoURL, tmpDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err := cloneCmd.Run(); err != nil {
		_ = os.RemoveAll(tmpDir)
		fmt.Fprintf(os.Stderr, constants.ErrScriptsClone, err)
		exitWith(1)
	}

	// Copy scripts from gitmap/scripts/ and root scripts.
	scriptSources := []struct {
		src  string
		name string
	}{
		{filepath.Join(tmpDir, "gitmap", "scripts", "install.ps1"), "install.ps1"},
		{filepath.Join(tmpDir, "gitmap", "scripts", "install.sh"), "install.sh"},
		{filepath.Join(tmpDir, "gitmap", "scripts", "uninstall.ps1"), "uninstall.ps1"},
		{filepath.Join(tmpDir, "gitmap", "scripts", "Get-LastRelease.ps1"), "Get-LastRelease.ps1"},
		{filepath.Join(tmpDir, "run.ps1"), "run.ps1"},
		{filepath.Join(tmpDir, "run.sh"), "run.sh"},
	}

	copied := 0

	for _, s := range scriptSources {
		data, err := os.ReadFile(s.src)
		if err != nil {
			fmt.Printf(constants.MsgScriptsSkip, s.name)

			continue
		}

		dest := filepath.Join(targetDir, s.name)

		if err := os.WriteFile(dest, data, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrScriptsCopy, s.name, err)

			continue
		}

		fmt.Printf(constants.MsgScriptsCopied, s.name)
		copied++
	}

	fmt.Println()
	fmt.Printf(constants.MsgScriptsDone, copied, targetDir)
}

// resolveScriptsDir returns the target directory for scripts.
// Windows: reads deployPath from powershell.json, defaults to D:\gitmap-scripts.
// Linux/macOS: ~/Desktop/gitmap-scripts.
func resolveScriptsDir() string {
	if runtime.GOOS == "windows" {
		return resolveScriptsDirWindows()
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return filepath.Join(home, "Desktop", "gitmap-scripts")
}

// resolveScriptsDirWindows reads powershell.json for the deploy drive.
// Search order:
//  1. <binaryDir>/powershell.json — written by install-quick.ps1
//  2. ./gitmap/powershell.json    — repo checkout
//  3. ./powershell.json           — repo root
//  4. Default: D:\gitmap-scripts
func resolveScriptsDirWindows() string {
	candidates := []string{}

	if exe, err := os.Executable(); err == nil {
		if resolved, evalErr := filepath.EvalSymlinks(exe); evalErr == nil {
			candidates = append(candidates, filepath.Join(filepath.Dir(resolved), "powershell.json"))
		} else {
			candidates = append(candidates, filepath.Join(filepath.Dir(exe), "powershell.json"))
		}
	}

	candidates = append(candidates,
		filepath.Join("gitmap", "powershell.json"),
		"powershell.json",
	)

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var cfg scriptsConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}

		if cfg.DeployPath != "" {
			drive := filepath.VolumeName(cfg.DeployPath)
			if drive != "" {
				return filepath.Join(drive+"\\", "gitmap-scripts")
			}
		}
	}

	// Default to D:\gitmap-scripts if no config found.
	return `D:\gitmap-scripts`
}
