package macro

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// DirTracker maintains current and previous working directories during interactive recording.
type DirTracker struct {
	CurrentDir string
	PrevDir    string
}

// NewDirTracker initializes a new directory tracker.
func NewDirTracker(initialDir string) *DirTracker {
	return &DirTracker{
		CurrentDir: initialDir,
		PrevDir:    initialDir,
	}
}

// ProcessCd checks if a command changes directory and updates tracked paths.
func (dt *DirTracker) ProcessCd(cmdText string) bool {
	target := parseCdTarget(cmdText)
	if target == "" {
		return false
	}
	resolved := dt.resolveTarget(target)
	info, err := os.Stat(resolved)
	if err != nil || !info.IsDir() {
		return false
	}
	dt.PrevDir = dt.CurrentDir
	dt.CurrentDir = resolved
	fmt.Printf("  ➜ 📁 Directory: %s%s%s\n", constants.ColorGreen, dt.CurrentDir, constants.ColorReset)
	return true
}

func (dt *DirTracker) resolveTarget(target string) string {
	if target == "-" {
		return dt.PrevDir
	}
	if filepath.IsAbs(target) {
		return filepath.Clean(target)
	}
	return filepath.Clean(filepath.Join(dt.CurrentDir, target))
}
