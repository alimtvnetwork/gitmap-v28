package osclean

import (
	"strings"
)

// IsAntigravityProtected returns true if path is an Antigravity IDE configuration or setting.
func IsAntigravityProtected(rawPath string) bool {
	if len(rawPath) == 0 {
		return false
	}
	normalized := normalizeCleanPath(rawPath)
	baseName := extractUniversalBaseName(normalized)

	if isProtectedBaseName(baseName) {
		return true
	}
	if !isAntigravityScope(normalized) {
		return false
	}

	return isProtectedWithinScope(normalized, baseName)
}

func normalizeCleanPath(path string) string {
	slash := strings.ReplaceAll(path, "\\", "/")

	return strings.ToLower(strings.TrimRight(slash, "/"))
}

func extractUniversalBaseName(normalized string) string {
	if idx := strings.LastIndex(normalized, "/"); idx >= 0 {
		return normalized[idx+1:]
	}

	return normalized
}

func isAntigravityScope(normalized string) bool {
	return strings.Contains(normalized, "/.gemini") ||
		strings.Contains(normalized, "/antigravity")
}

func isAntigravityRoot(normalized string) bool {
	return strings.HasSuffix(normalized, "/.gemini") ||
		strings.HasSuffix(normalized, "/.gemini/config") ||
		strings.HasSuffix(normalized, "/.gemini/antigravity") ||
		strings.HasSuffix(normalized, "/appdata/roaming/antigravity")
}

func isProtectedWithinScope(normalized string, baseName string) bool {
	if isAntigravityRoot(normalized) {
		return true
	}
	if hasProtectedDirToken(normalized + "/") {
		return true
	}

	return !hasSafeCacheToken(normalized + "/")
}
