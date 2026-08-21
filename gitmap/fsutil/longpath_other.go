//go:build !windows

package fsutil

// EnsureLongPath is a pass-through on non-Windows platforms.
func EnsureLongPath(path string) string {
	return path
}

// StripLongPathPrefix is a pass-through on non-Windows platforms.
func StripLongPathPrefix(path string) string {
	return path
}
