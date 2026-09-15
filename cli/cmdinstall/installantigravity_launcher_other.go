//go:build !linux

package cmdinstall

// PurgeDuplicateAntigravityLaunchers is a no-op on non-Linux platforms.
func PurgeDuplicateAntigravityLaunchers() {}
