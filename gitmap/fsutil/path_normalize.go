// Package fsutil — path_normalize.go provides cross-platform path utilities for installer portable paths.
package fsutil

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// NormalizeToForwardSlashes converts all backslashes to forward slashes.
func NormalizeToForwardSlashes(p string) string {
	return strings.ReplaceAll(filepath.Clean(p), "\\", "/")
}

// MakeRelativeToRoot calculates relative path from root and normalizes with forward slashes.
func MakeRelativeToRoot(base, target string) (string, error) {
	if strings.TrimSpace(base) == "" || strings.TrimSpace(target) == "" {
		return "", apperror.New("MakeRelativeToRoot", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "base and target cannot be empty",
		})
	}

	rel, err := filepath.Rel(base, target)
	if err != nil {
		appErr := apperror.Wrap(err, "MakeRelativeToRoot", map[string]any{
			"base":   base,
			"target": target,
		})
		appErr.Code = "E_INSTALLER_PATH_ERROR"
		return "", appErr
	}

	return NormalizeToForwardSlashes(rel), nil
}

// TrimTrailingSlashes removes trailing slashes from path string.
func TrimTrailingSlashes(p string) string {
	return strings.TrimRight(p, "/\\")
}

// NormalizeSlashes converts backslashes to forward slashes.
func NormalizeSlashes(p string) string {
	return NormalizeToForwardSlashes(p)
}

// EqualPaths checks if two paths are identical after normalization.
func EqualPaths(p1, p2 string) bool {
	return NormalizeToForwardSlashes(p1) == NormalizeToForwardSlashes(p2)
}

// IsSubdirectory checks if child is within parent.
func IsSubdirectory(parent, child string) bool {
	p := NormalizeToForwardSlashes(parent)
	c := NormalizeToForwardSlashes(child)
	return strings.HasPrefix(c, p+"/")
}
