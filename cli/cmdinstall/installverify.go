package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var guiTools = map[string]bool{
	constants.ToolNpp:           true,
	constants.ToolNppInstall:    true,
	constants.ToolGitHubDesktop: true,
	constants.ToolDbeaver:       true,
	constants.ToolOBS:           true,
	constants.ToolStickyNotes:   true,
	constants.ToolQBittorrent:   true,
	constants.ToolUTorrent:      true,
}

// isGUITool returns true for tools that open a GUI window on --version.
func isGUITool(tool string) bool {
	if isGUI := guiTools[tool]; isGUI {
		return true
	}

	return false
}

func verifyExeDirect(tool, exePath string) bool {
	if exePath == "" {
		return false
	}

	if _, err := os.Stat(exePath); err == nil {
		fmt.Printf(constants.MsgInstallSuccess, tool)
		fmt.Printf(constants.MsgInstallExeFound, exePath)
		_ = runPostInstall(tool)

		return true
	}

	return false
}

func reportVerifiedTool(tool, version string) {
	fmt.Printf(constants.MsgInstallSuccess, tool)
	fmt.Printf("  → Detected version: %s\n", version)
	verifyExePath(tool)
	_ = runPostInstall(tool)
}

func verifyByVersion(tool string) {
	if isGUITool(tool) {
		reportVerificationFailure(tool, toolBinaryName(tool))

		return
	}

	bin := toolBinaryName(tool)
	version := getInstalledVersion(bin)
	if version == "" {
		reportVerificationFailure(tool, bin)

		return
	}

	reportVerifiedTool(tool, version)
}

// verifyInstallation confirms a tool is accessible after install.
func verifyInstallation(tool string) {
	fmt.Printf(constants.MsgInstallVerifying, tool)

	if isVerified := verifyExeDirect(tool, expectedExePath(tool)); isVerified {
		return
	}

	verifyByVersion(tool)
}

// verifyExePath checks the expected exe path exists after install.
func verifyExePath(tool string) {
	exePath := expectedExePath(tool)
	if exePath == "" {
		return
	}

	fmt.Printf(constants.MsgInstallExeVerify, tool, exePath)
	if _, err := os.Stat(exePath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrInstallExeNotFound, exePath)

		return
	}

	fmt.Printf(constants.MsgInstallExeFound, exePath)
}

var windowsExeMap = map[string]string{
	constants.ToolNpp:         `C:\Program Files\Notepad++\notepad++.exe`,
	constants.ToolVSCode:      `C:\Program Files\Microsoft VS Code\Code.exe`,
	constants.ToolDbeaver:     `C:\Program Files\DBeaver\dbeaver.exe`,
	constants.ToolOBS:         `C:\Program Files\obs-studio\bin\64bit\obs64.exe`,
	constants.ToolNginx:       `C:\tools\nginx\nginx.exe`,
	constants.ToolQBittorrent: `C:\Program Files\qBittorrent\qbittorrent.exe`,
	constants.ToolUTorrent:    `C:\Program Files (x86)\uTorrent\uTorrent.exe`,
}

// expectedExePath returns the expected binary path for a tool.
func expectedExePath(tool string) string {
	if runtime.GOOS != "windows" {
		return ""
	}

	if path, isFound := windowsExeMap[tool]; isFound {
		return path
	}

	return ""
}

func detectDirectExe(tool string) string {
	exePath := expectedExePath(tool)
	if exePath == "" {
		return ""
	}

	if _, err := os.Stat(exePath); err == nil {
		return "installed (at " + exePath + ")"
	}

	return ""
}

// detectInstalledVersion checks if a tool is already installed.
func detectInstalledVersion(tool string) string {
	if directPath := detectDirectExe(tool); directPath != "" {
		return directPath
	}

	if isGUITool(tool) {
		return ""
	}

	return getInstalledVersion(toolBinaryName(tool))
}

func versionFlag(binary string) string {
	if binary == "nginx" {
		return "-v"
	}

	return "--version"
}

func execVersion(path, flag string) string {
	out, err := exec.Command(path, flag).CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// getInstalledVersion runs --version or -v and returns the output.
func getInstalledVersion(binary string) string {
	path := resolveToolBinaryPath(binary)
	if path == "" {
		return ""
	}

	return execVersion(path, versionFlag(binary))
}

var toolBinaryMap = map[string]string{
	constants.ToolVSCode:        "code",
	constants.ToolNodeJS:        "node",
	constants.ToolYarn:          "yarn",
	constants.ToolBun:           "bun",
	constants.ToolPnpm:          "pnpm",
	constants.ToolPython:        "python3",
	constants.ToolGo:            "go",
	constants.ToolGit:           "git",
	constants.ToolGitLFS:        "git-lfs",
	constants.ToolGHCLI:         "gh",
	constants.ToolGitHubDesktop: "github-desktop",
	constants.ToolCPP:           "g++",
	constants.ToolPHP:           "php",
	constants.ToolPowerShell:    "pwsh",
	constants.ToolNpp:           "notepad++",
	constants.ToolNppInstall:    "notepad++",
	constants.ToolNginx:         "nginx",
	constants.ToolWordPress:     "wp",
	constants.ToolLaravel:       "laravel",
	constants.ToolVMware:        "vmhgfs-fuse",
	constants.ToolAntigravity:   "antigravity",
	constants.ToolAgy:           "agy",
	constants.ToolAgManager:     "ag-manager",
	constants.ToolQBittorrent:   "qbittorrent",
	constants.ToolUTorrent:      "utorrent",
	constants.ToolZsh:           "zsh",
}

// toolBinaryName maps tool names to their binary/executable names.
func toolBinaryName(tool string) string {
	if binary, isFound := toolBinaryMap[tool]; isFound {
		return binary
	}

	return tool
}

// runPostInstall executes post-install actions for specific tools.
func runPostInstall(tool string) error {
	if tool == constants.ToolGitLFS {
		return runPostInstallGitLFS()
	}

	if tool == constants.ToolGit {
		return runPostInstallGit()
	}

	return nil
}

// runPostInstallGitLFS runs git lfs install.
func runPostInstallGitLFS() error {
	cmd := exec.Command("git", "lfs", "install")

	return cmd.Run()
}

// runPostInstallGit configures git longpaths.
func runPostInstallGit() error {
	cmd := exec.Command("git", "config", "--global", "core.longpaths", "true")

	return cmd.Run()
}
