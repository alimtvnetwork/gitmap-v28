// Package cmd — open_target.go: target resolution and execution for gitmap open.
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func isEditorTarget(arg string) bool {
	low := strings.ToLower(strings.TrimSpace(arg))

	return low == "code" || low == "vscode"
}

func isWebTarget(arg string) bool {
	trimmed := strings.TrimSpace(arg)
	if hasWebProtocol(trimmed) {
		return true
	}

	return hasDomainFormat(trimmed)
}

func hasWebProtocol(target string) bool {
	low := strings.ToLower(target)

	return strings.HasPrefix(low, "http://") || strings.HasPrefix(low, "https://") || strings.HasPrefix(low, "mailto:")
}

func hasDomainFormat(target string) bool {
	low := strings.ToLower(target)
	tlds := []string{".com", ".org", ".net", ".io", ".dev", ".co", ".app", ".ai", ".tv", ".me", ".edu", ".gov"}
	for _, tld := range tlds {
		if strings.Contains(low, tld) {
			return isNonExistentLocalPath(target)
		}
	}

	return false
}

func isNonExistentLocalPath(target string) bool {
	_, err := os.Stat(target)

	return os.IsNotExist(err)
}

func normalizeWebUrl(target string) string {
	trimmed := strings.TrimSpace(target)
	if hasWebProtocol(trimmed) {
		return trimmed
	}

	return "https://" + trimmed
}

func launchEditor(args []string) error {
	targetPath := "."
	if len(args) > 1 && args[1] != "" {
		targetPath = expandTilde(args[1])
	}

	bin := findEditorBinary()
	if bin == "" {
		absPath, _ := filepath.Abs(targetPath)

		return launchNativeOpener(absPath)
	}

	cmd := exec.Command(bin, targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return apperror.WrapSimple(err, "launch editor")
	}

	return nil
}

func findEditorBinary() string {
	candidates := []string{"code"}
	if runtime.GOOS == "windows" {
		candidates = []string{"code.cmd", "code"}
	}

	for _, cand := range candidates {
		if path, err := exec.LookPath(cand); err == nil {
			return path
		}
	}

	return ""
}
