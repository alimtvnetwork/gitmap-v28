// Package cmdagy — agy_clean_cache.go provides core cache cleaning logic.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/osclean"
)

// CalculateDirStats recursively walks a directory to compute total size and file count.
func CalculateDirStats(dirPath string) (int64, int, error) {
	var totalSize int64
	var fileCount int

	walkErr := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	return totalSize, fileCount, walkErr
}

// CleanDirectoryContents removes all items inside dirPath while preserving the parent directory.
func CleanDirectoryContents(dirPath string) (int64, int, []string) {
	var freedBytes int64
	var deletedFiles int
	var warnings []string

	if osclean.IsAntigravityProtected(dirPath) {
		warnings = append(warnings, fmt.Sprintf("skipped protected directory %s", dirPath))
		return freedBytes, deletedFiles, warnings
	}

	entries, readErr := os.ReadDir(dirPath)
	if readErr != nil {
		warnings = append(warnings, fmt.Sprintf("read dir %s: %v", dirPath, readErr))
		return freedBytes, deletedFiles, warnings
	}

	for _, entry := range entries {
		subPath := filepath.Join(dirPath, entry.Name())
		if osclean.IsAntigravityProtected(subPath) {
			warnings = append(warnings, fmt.Sprintf("skipped protected path %s", subPath))
			continue
		}
		subSize, subFiles, _ := CalculateDirStats(subPath)
		if removeErr := os.RemoveAll(subPath); removeErr != nil {
			warnings = append(warnings, fmt.Sprintf("failed removing %s: %v", subPath, removeErr))
			continue
		}
		freedBytes += subSize
		deletedFiles += countOrOne(subFiles)
	}

	return freedBytes, deletedFiles, warnings
}

func countOrOne(count int) int {
	if count > 0 {
		return count
	}
	return 1
}
