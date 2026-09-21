//go:build !windows

package cmdagy

// FocusAntigravityWindow is a non-Windows no-op.
func FocusAntigravityWindow(targetPID int) bool {
	return false
}
