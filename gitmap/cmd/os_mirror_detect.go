// Package cmd — os_mirror_detect.go detects problematic regional mirrors in sources.list.
package cmd

import (
	"os"
	"strings"
)

// HasRegionalMirrorGlitch checks if /etc/apt/sources.list contains regional mirrors prone to glitches.
func HasRegionalMirrorGlitch(path string) bool {
	if path == "" {
		path = "/etc/apt/sources.list"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	content := string(data)
	return strings.Contains(content, "my.archive.ubuntu.com") ||
		strings.Contains(content, ".archive.ubuntu.com") && !strings.Contains(content, "http://archive.ubuntu.com")
}
