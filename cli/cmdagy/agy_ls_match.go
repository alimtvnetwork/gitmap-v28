// Package cmdagy — agy_ls_match.go provides directory existence and string filter matching helpers.
package cmdagy

import (
	"os"
	"strings"
)

func matchesFilter(name, path, filter string) bool {
	if strings.EqualFold(name, filter) || strings.EqualFold(path, filter) {
		return true
	}

	term := strings.ToLower(filter)
	if strings.Contains(strings.ToLower(name), term) {
		return true
	}

	return strings.Contains(strings.ToLower(path), term)
}

func checkDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
