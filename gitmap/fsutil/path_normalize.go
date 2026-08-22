package fsutil

import (
	"path/filepath"
	"strings"
)

// NormalizeSlashes converts backslashes to forward slashes.
func NormalizeSlashes(path string) string {
	return filepath.ToSlash(filepath.Clean(strings.ReplaceAll(path, "\\", "/")))
}

// TrimTrailingSlashes removes trailing slashes and backslashes.
func TrimTrailingSlashes(path string) string {
	trimmed := strings.TrimRight(path, "/\\")
	if len(trimmed) == 0 && len(path) > 0 {
		return path
	}
	return trimmed
}

// CanonicalPath returns a clean, absolute, normalized path.
func CanonicalPath(path string) (string, error) {
	trimmed := TrimTrailingSlashes(path)
	abs, err := filepath.Abs(trimmed)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

// EqualPaths checks case-insensitive equality across path formats.
func EqualPaths(pathA, pathB string) bool {
	normA := strings.ToLower(NormalizeSlashes(pathA))
	normB := strings.ToLower(NormalizeSlashes(pathB))
	return normA == normB
}

// IsSubdirectory reports whether child is inside parent directory.
func IsSubdirectory(parent, child string) bool {
	pNorm := strings.ToLower(NormalizeSlashes(parent))
	cNorm := strings.ToLower(NormalizeSlashes(child))
	if pNorm == cNorm {
		return true
	}
	if !strings.HasSuffix(pNorm, "/") {
		pNorm += "/"
	}
	return strings.HasPrefix(cNorm, pNorm)
}
