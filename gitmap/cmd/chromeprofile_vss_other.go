//go:build !windows

package cmd

// Snapshot is the non-Windows stub for the VSS shadow-copy helper
// (#8). The Windows build provides a real implementation; on other
// platforms callers always receive ok=false and fall back to the
// regular skip-list copy path.
type Snapshot struct {
	ID         string
	DevicePath string
	Volume     string
}

// CreateSnapshot is a no-op on non-Windows builds.
func CreateSnapshot(_ string) (Snapshot, bool) { return Snapshot{}, false }

// TranslatePath returns srcAbs unchanged on non-Windows builds.
func (s Snapshot) TranslatePath(srcAbs string) string { return srcAbs }

// Delete is a no-op on non-Windows builds.
func (s Snapshot) Delete() {}
