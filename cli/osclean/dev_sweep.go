package osclean

import (
	"os"
	"path/filepath"
	"runtime"
)

// SweepTarget cleans or simulates cleaning on a directory target.
func SweepTarget(dir string, isDryRun bool, hasStripReadOnly bool) CategoryCleanStats {
	var stats CategoryCleanStats
	if !isExistingDir(dir) {
		return stats
	}

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		stats.Errors = append(stats.Errors, readErr.Error())
		return stats
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		sweepEntry(fullPath, entry, isDryRun, hasStripReadOnly, &stats)
	}

	return stats
}

func isExistingDir(path string) bool {
	if len(path) == 0 {
		return false
	}
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

func sweepEntry(p string, entry os.DirEntry, isDryRun, hasStrip bool, stats *CategoryCleanStats) {
	if isDryRun {
		recordDryRunEntry(p, entry, stats)
		return
	}

	if hasStrip && runtime.GOOS == "windows" {
		stripReadOnlyAttr(p)
	}

	recordLiveRemoval(p, entry, stats)
}

func recordDryRunEntry(p string, entry os.DirEntry, stats *CategoryCleanStats) {
	if entry.IsDir() {
		dirFiles, dirBytes := measureDirRecursive(p)
		stats.DirsRemoved++
		stats.ItemsRemoved += dirFiles
		stats.BytesFreed += dirBytes
		return
	}

	stats.ItemsRemoved++
	if info, err := entry.Info(); err == nil {
		stats.BytesFreed += info.Size()
	}
}

func recordLiveRemoval(p string, entry os.DirEntry, stats *CategoryCleanStats) {
	dirFiles, dirBytes := measureDirRecursive(p)
	if removeErr := os.RemoveAll(p); removeErr != nil {
		stats.Errors = append(stats.Errors, removeErr.Error())
		return
	}

	if entry.IsDir() {
		stats.DirsRemoved++
		stats.ItemsRemoved += dirFiles
	} else {
		stats.ItemsRemoved++
	}
	stats.BytesFreed += dirBytes
}

func stripReadOnlyAttr(root string) {
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err == nil && info != nil && isReadOnlyFile(info) {
			_ = os.Chmod(p, 0o666)
		}
		return nil
	})
}

func isReadOnlyFile(info os.FileInfo) bool {
	return (info.Mode() & 0o200) == 0
}

func measureDirRecursive(root string) (int, int64) {
	var count int
	var totalBytes int64

	_ = filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			count++
			totalBytes += info.Size()
		}
		return nil
	})

	return count, totalBytes
}
