package cmdcopy

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy/constants"
)

// gitTopLevel returns the absolute path of the current repo's top-level dir.
func gitTopLevel() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, constants.GitArgShowTopLevel)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", apperror.NewValidationError(constants.ErrMsgEmptyTopLevel)
	}

	return filepath.Clean(root), nil
}

func isExistingFile(path string) bool {
	info, err := os.Stat(path)

	return err == nil && !info.IsDir()
}

func isFileFlagWithArg(arg string) bool {
	lower := strings.ToLower(arg)
	for _, flag := range constants.AllCopyFileFlags() {
		if lower == flag {
			return true
		}
	}

	return false
}

func isTerminalInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}
