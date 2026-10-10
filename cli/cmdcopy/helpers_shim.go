package cmdcopy

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// gitTopLevel returns the absolute path of the current repo's top-level dir.
func gitTopLevel() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("empty top-level")
	}

	return filepath.Clean(root), nil
}

func isExistingFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isFileFlagWithArg(arg string) bool {
	lower := strings.ToLower(arg)

	return lower == "--file" || lower == "--filepath" || lower == "--out" || lower == "--output" || lower == "-o" || lower == "-f"
}

func isTerminalInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}
