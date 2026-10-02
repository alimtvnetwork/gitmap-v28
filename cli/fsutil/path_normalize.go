// Package fsutil — path_normalize.go provides cross-platform path utilities for installer portable paths.
package fsutil

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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

// IsPathCaseInsensitive reports whether the given path should be treated case-insensitively.
// It returns true on Windows (runtime.GOOS == "windows") or if the path exhibits Windows volume syntax (e.g. C:\ or c:/).
func IsPathCaseInsensitive(p string) bool {
	if runtime.GOOS == "windows" {
		return true
	}
	clean := filepath.Clean(strings.TrimSpace(p))
	vol := filepath.VolumeName(clean)
	if len(vol) >= 2 && vol[1] == ':' {
		return true
	}

	return false
}

// EqualPaths checks if two paths are identical after normalization, respecting OS case sensitivity.
// On Windows (or for Windows-style volume paths), it uses zero-allocation strings.EqualFold.
// On Unix/Linux, it enforces exact byte equality.
func EqualPaths(p1, p2 string) bool {
	n1 := NormalizeToForwardSlashes(p1)
	n2 := NormalizeToForwardSlashes(p2)
	if IsPathCaseInsensitive(p1) || IsPathCaseInsensitive(p2) {
		return strings.EqualFold(n1, n2)
	}

	return n1 == n2
}

// CanonicalPathKey produces a normalized path key respecting host OS case sensitivity.
// On Windows (or for Windows-style paths), it lowercases the path. On Unix/Linux, it preserves casing.
func CanonicalPathKey(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}
	slashed := filepath.ToSlash(filepath.Clean(trimmed))
	if IsPathCaseInsensitive(slashed) {
		return strings.ToLower(slashed)
	}

	return slashed
}

// IsSubdirectory checks if child is within parent, respecting OS case sensitivity.
func IsSubdirectory(parent, child string) bool {
	p := strings.TrimRight(NormalizeToForwardSlashes(parent), "/")
	c := NormalizeToForwardSlashes(child)

	if IsPathCaseInsensitive(parent) || IsPathCaseInsensitive(child) {
		pPrefix := p + "/"
		if len(c) > len(pPrefix) && strings.EqualFold(c[:len(pPrefix)], pPrefix) {
			return true
		}

		return false
	}

	return strings.HasPrefix(c, p+"/")
}
