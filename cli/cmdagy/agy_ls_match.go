// Package cmdagy — agy_ls_match.go provides directory existence and string filter matching helpers.
package cmdagy

import (
	"os"
	"strings"
)

func matchesFilter(name, path, filter string) bool {
	term := strings.ToLower(filter)
	matchesName := strings.Contains(strings.ToLower(name), term)
	matchesPath := strings.Contains(strings.ToLower(path), term)

	return matchesName || matchesPath
}

func checkDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
