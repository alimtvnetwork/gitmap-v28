// Package cmdzsh provides parsers for .zshrc theme and plugin inspection.
package cmdzsh

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func inspectThemeInZshrc(zshrcPath string) string {
	lines := readZshrcLines(zshrcPath)

	return findThemePrefix(lines)
}

func readZshrcLines(path string) []string {
	if !HasFile(path) {
		return nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	return strings.Split(string(bytes), constants.NewLineUnix)
}

func findThemePrefix(lines []string) string {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ZSH_THEME=") {
			return extractQuotedValue(trimmed)
		}
	}

	return ""
}

func extractQuotedValue(line string) string {
	idx := strings.Index(line, "=")
	if idx < 0 {
		return ""
	}

	val := strings.TrimSpace(line[idx+1:])

	return strings.Trim(val, `"'`)
}

func inspectPluginsInZshrc(zshrcPath string) []string {
	lines := readZshrcLines(zshrcPath)

	return findPluginsLine(lines)
}

func findPluginsLine(lines []string) []string {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "plugins=(") {
			return extractPluginNames(trimmed)
		}
	}

	return nil
}

func extractPluginNames(line string) []string {
	idxStart := strings.Index(line, "(")
	idxEnd := strings.LastIndex(line, ")")
	if idxStart < 0 || idxEnd <= idxStart {
		return nil
	}

	raw := strings.TrimSpace(line[idxStart+1 : idxEnd])
	if raw == "" {
		return nil
	}

	return strings.Fields(raw)
}
