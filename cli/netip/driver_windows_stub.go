//go:build !windows

package netip

// NewWindowsDriver returns nil on non-Windows platforms.
func NewWindowsDriver() Driver {
	return nil
}
