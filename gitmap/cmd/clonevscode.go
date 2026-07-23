package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// isVSCodeAvailable checks if the VS Code CLI is on PATH.
func isVSCodeAvailable() bool {
	if _, err := exec.LookPath(constants.VSCodeBin); err == nil {
		return true
	}

	return resolveVSCodeExecutable() != ""
}

// openInVSCode opens the given folder in VS Code.
// Tries multiple strategies to bypass "Another instance running as administrator":
// 1. code --reuse-window (standard)
// 2. code --new-window (bypasses some admin conflicts)
// 3. Launch Code.exe with an isolated user-data dir in a detached process.
func openInVSCode(absPath string) {
	if !isVSCodeAvailable() {
		fmt.Fprintf(os.Stdout, constants.MsgVSCodeNotFound)

		return
	}

	fmt.Printf(constants.MsgVSCodeOpening, absPath)

	if tryVSCodeCLI(absPath) {
		fmt.Println(constants.MsgVSCodeOpened)

		return
	}

	if tryVSCodeDetached(absPath) {
		fmt.Println(constants.MsgVSCodeOpened)

		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrVSCodeAdminLock)
}

// tryVSCodeCLI attempts both supported VS Code CLI launch modes.
func tryVSCodeCLI(absPath string) bool {
	if tryVSCodeReuse(absPath) {
		return true
	}

	return tryVSCodeNewWindow(absPath)
}

// tryVSCodeReuse attempts to open with --reuse-window.
func tryVSCodeReuse(absPath string) bool {
	return runVSCodeCommand(constants.VSCodeBin,
		constants.VSCodeFlagReuseWindow,
		absPath) == nil
}

// tryVSCodeNewWindow attempts to open with --new-window.
func tryVSCodeNewWindow(absPath string) bool {
	return runVSCodeCommand(constants.VSCodeBin,
		constants.VSCodeFlagNewWindow,
		absPath) == nil
}

// tryVSCodeDetached launches Code.exe via cmd /C start using a separate
// user-data directory so it does not try to hand off to an elevated instance.
func tryVSCodeDetached(absPath string) bool {
	exePath := resolveVSCodeExecutable()
	if exePath == "" {
		return false
	}

	userDataDir := filepath.Join(os.TempDir(), constants.VSCodeUserDataDirName)
	if err := os.MkdirAll(userDataDir, constants.DirPermission); err != nil {
		return false
	}

	cmd := exec.Command(
		constants.CmdWindowsShell,
		constants.CmdArgSlashC,
		constants.CmdArgStart,
		constants.CmdArgEmpty,
		exePath,
		constants.VSCodeFlagNewWindow,
		constants.VSCodeFlagUserDataDir,
		userDataDir,
		absPath,
	)

	return cmd.Run() == nil
}

// runVSCodeCommand executes a VS Code CLI command and waits for its exit code.
func runVSCodeCommand(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return err
	}

	return fmt.Errorf("%w: %s", err, trimmed)
}

// resolveVSCodeExecutable finds the desktop executable used for detached launch.
func resolveVSCodeExecutable() string {
	for _, candidate := range vscodeExecutableCandidates() {
		if isExistingFile(candidate) {
			return candidate
		}
	}

	return ""
}

// vscodeExecutableCandidates returns likely Code.exe locations.
func vscodeExecutableCandidates() []string {
	cliPath, _ := exec.LookPath(constants.VSCodeBin)

	return []string{
		lookPathOrEmpty(constants.VSCodeExeBin),
		candidateVSCodeExecutable(cliPath),
		filepath.Join(os.Getenv(constants.VSCodeEnvLocalAppData),
			constants.VSCodeProgramsDirName,
			constants.VSCodeInstallDirName,
			constants.VSCodeExeBin),
		filepath.Join(os.Getenv(constants.VSCodeEnvProgramFiles),
			constants.VSCodeInstallDirName,
			constants.VSCodeExeBin),
		filepath.Join(os.Getenv(constants.VSCodeEnvProgramFilesX86),
			constants.VSCodeInstallDirName,
			constants.VSCodeExeBin),
	}
}

// lookPathOrEmpty returns the PATH match or an empty string.
func lookPathOrEmpty(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}

	return path
}

// candidateVSCodeExecutable resolves Code.exe from the CLI's bin directory.
func candidateVSCodeExecutable(cliPath string) string {
	if cliPath == "" {
		return ""
	}

	return filepath.Clean(filepath.Join(filepath.Dir(cliPath), "..", constants.VSCodeExeBin))
}

// isExistingFile reports whether path exists and is a file.
func isExistingFile(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
