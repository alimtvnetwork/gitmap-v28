// Package cmdzsh provides default plugin clean up for .zshrc.
package cmdzsh

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RemoveDefaultGitPlugin removes default plugins=(git) entry from ~/.zshrc.
func RemoveDefaultGitPlugin(targetHome string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}

	home := ResolveTargetHome(targetHome)
	zshrcPath := filepath.Join(home, ".zshrc")
	if !HasFile(zshrcPath) {
		return nil
	}

	return cleanGitPluginLine(zshrcPath)
}

func cleanGitPluginLine(zshrcPath string) *apperror.AppError {
	bytes, err := os.ReadFile(zshrcPath)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.cleanGitPluginLine.read")
	}

	lines := strings.Split(string(bytes), constants.NewLineUnix)
	filtered, isModified := filterDefaultPluginLines(lines)
	if !isModified {
		return nil
	}

	output := strings.Join(filtered, constants.NewLineUnix)

	return writeZshrcBytes(zshrcPath, []byte(output))
}

func filterDefaultPluginLines(lines []string) ([]string, bool) {
	var out []string
	isModified := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "plugins=(git)" {
			isModified = true
			continue
		}

		out = append(out, line)
	}

	return out, isModified
}
