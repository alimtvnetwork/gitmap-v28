// Package cmdzsh provides ZSH theme management and .zshrc mutation logic.
package cmdzsh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ApplyTheme modifies or adds ZSH_THEME in target .zshrc.
func ApplyTheme(targetHome, themeName string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}

	home := ResolveTargetHome(targetHome)
	zshrcPath := filepath.Join(home, ".zshrc")
	if !HasFile(zshrcPath) {
		entry := fmt.Sprintf(`ZSH_THEME="%s"%s`, themeName, constants.NewLineUnix)

		return writeZshrcBytes(zshrcPath, []byte(entry))
	}

	return updateThemeInZshrc(zshrcPath, themeName)
}

func updateThemeInZshrc(zshrcPath, themeName string) *apperror.AppError {
	bytes, err := os.ReadFile(zshrcPath)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.updateThemeInZshrc.read")
	}

	lines := strings.Split(string(bytes), constants.NewLineUnix)
	updated := renderThemeLines(lines, themeName)
	output := strings.Join(updated, constants.NewLineUnix)

	return writeZshrcBytes(zshrcPath, []byte(output))
}

func renderThemeLines(lines []string, themeName string) []string {
	entry := fmt.Sprintf(`ZSH_THEME="%s"`, themeName)
	out, hasSet := replaceThemeLine(lines, entry)
	if !hasSet {
		out = append(out, entry)
	}

	return out
}

func replaceThemeLine(lines []string, entry string) ([]string, bool) {
	var out []string
	hasSet := false
	for _, l := range lines {
		if strings.HasPrefix(l, "ZSH_THEME=") {
			out = append(out, entry)
			hasSet = true
			continue
		}

		out = append(out, l)
	}

	return out, hasSet
}
