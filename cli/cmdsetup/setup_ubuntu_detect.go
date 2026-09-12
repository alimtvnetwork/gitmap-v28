package cmdsetup

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	lookPathFunc    = exec.LookPath
	statPathFunc    = os.Stat
	execCommandFunc = exec.Command
	getenvFunc      = os.Getenv
	userHomeDirFunc = os.UserHomeDir
	stdinStatFunc   = os.Stdin.Stat
)

func isFilePresent(path string) bool {
	_, statErr := statPathFunc(path)

	return statErr == nil
}

func findZshBinary() (string, bool) {
	path, lookErr := lookPathFunc("zsh")
	if lookErr == nil {
		return path, true
	}

	for _, candidate := range []string{"/bin/zsh", "/usr/bin/zsh", "/usr/local/bin/zsh"} {
		if isFilePresent(candidate) {
			return candidate, true
		}
	}

	return "", false
}

func getZshVersion(zshPath string) string {
	cmd := execCommandFunc(zshPath, "--version")
	outBytes, runErr := cmd.Output()
	hasOutput := runErr == nil && len(outBytes) > 0
	if hasOutput {
		return parseVersionFromOutput(string(outBytes))
	}

	return ""
}

func isDirPresent(path string) bool {
	stat, statErr := statPathFunc(path)
	if statErr == nil && stat.IsDir() {
		return true
	}

	return false
}

func isOhMyZshInstalled() bool {
	zshEnv := getenvFunc("ZSH")
	if len(zshEnv) > 0 && isDirPresent(zshEnv) {
		return true
	}

	homeDir, homeErr := userHomeDirFunc()
	if homeErr == nil && isDirPresent(filepath.Join(homeDir, ".oh-my-zsh")) {
		return true
	}

	return false
}

func isStdinTerminal() bool {
	stat, statErr := stdinStatFunc()
	if statErr == nil {
		isCharDevice := (stat.Mode() & os.ModeCharDevice) != 0

		return isCharDevice
	}

	return false
}

func isSkipZshEnv() bool {
	val := getenvFunc(constants.EnvGitmapSkipZsh)
	isSkip := val == "1"

	return isSkip
}
