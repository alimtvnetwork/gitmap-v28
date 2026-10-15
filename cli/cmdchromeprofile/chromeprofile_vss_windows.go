//go:build windows

package cmdchromeprofile

// Snapshot represents an optional point-in-time snapshot handle.
// On Windows we use standard shared-read file copy fallback rather than
// invoking shadow copy utilities, avoiding antivirus false-positive
// heuristics on unsigned CLI binaries.
type Snapshot struct {
	ID         string
	DevicePath string
	Volume     string
}

// CreateSnapshot returns ok=false so callers use standard shared-read copy.
func CreateSnapshot(_ string) (Snapshot, bool) {
	return Snapshot{}, false
}

// TranslatePath returns srcAbs unchanged when no snapshot is active.
func (s Snapshot) TranslatePath(srcAbs string) string {
	if s.DevicePath == "" {
		return srcAbs
	}
	return srcAbs
}

// Delete is a no-op when no snapshot is active.
func (s Snapshot) Delete() {}
