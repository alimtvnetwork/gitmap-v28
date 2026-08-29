//go:build windows

package fsutil

import (
	"path/filepath"
	"strings"
)

const winLongPathPrefix = `\\?\`
const winLongPathUNCPrefix = `\\?\UNC\`
//nolint:unused
const winMaxPathThreshold = 240

// EnsureLongPath prefixes path with \\?\ if needed on Windows.
func EnsureLongPath(path string) string {
	if len(path) == 0 {
		return path
	}
	if strings.HasPrefix(path, winLongPathPrefix) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	if strings.HasPrefix(abs, `\\`) {
		return winLongPathUNCPrefix + strings.TrimPrefix(abs, `\\`)
	}
	return winLongPathPrefix + abs
}

// StripLongPathPrefix removes \\?\ or \\?\UNC\ for DB persistence.
func StripLongPathPrefix(path string) string {
	if strings.HasPrefix(path, winLongPathUNCPrefix) {
		return `\\` + strings.TrimPrefix(path, winLongPathUNCPrefix)
	}
	if strings.HasPrefix(path, winLongPathPrefix) {
		return strings.TrimPrefix(path, winLongPathPrefix)
	}
	return path
}
