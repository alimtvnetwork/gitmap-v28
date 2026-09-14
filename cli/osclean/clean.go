package osclean

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// CleanTempDirectories sweeps OS ephemeral directories.
func CleanTempDirectories(opts CleanOptions) CleanResult {
	dirs := resolveTempDirectories()
	var totalStats CleanStats

	for _, dir := range dirs {
		stats := cleanSingleDirectory(dir, opts)
		totalStats.RemovedFilesCount += stats.RemovedFilesCount
		totalStats.RemovedDirsCount += stats.RemovedDirsCount
		totalStats.FreedBytes += stats.FreedBytes
	}

	return result.Ok(totalStats)
}

func resolveTempDirectories() []string {
	if runtime.GOOS == constants.OSWindows {
		return resolveWindowsTempDirs()
	}

	return resolveUnixTempDirs()
}

func resolveWindowsTempDirs() []string {
	seen := make(map[string]bool)
	var dirs []string
	appendIfPresent(&dirs, seen, os.TempDir())
	appendIfPresent(&dirs, seen, os.Getenv("TEMP"))
	appendIfPresent(&dirs, seen, filepath.Join(os.Getenv("LOCALAPPDATA"), "Temp"))

	return dirs
}

func resolveUnixTempDirs() []string {
	seen := make(map[string]bool)
	var dirs []string
	appendIfPresent(&dirs, seen, "/tmp")
	appendIfPresent(&dirs, seen, "/var/tmp")
	if home := os.Getenv("HOME"); home != "" {
		appendIfPresent(&dirs, seen, filepath.Join(home, ".cache", "gitmap"))
	}

	return dirs
}

func appendIfPresent(dirs *[]string, seen map[string]bool, path string) {
	if path == "" || seen[path] {
		return
	}
	seen[path] = true
	*dirs = append(*dirs, path)
}

func cleanSingleDirectory(dir string, opts CleanOptions) CleanStats {
	var stats CleanStats
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return stats
	}

	for _, entry := range entries {
		processEntry(dir, entry, opts, &stats)
	}

	return stats
}

func processEntry(dir string, entry os.DirEntry, opts CleanOptions, stats *CleanStats) {
	fullPath := filepath.Join(dir, entry.Name())
	info, statErr := entry.Info()
	if statErr != nil {
		return
	}

	if opts.IsDryRun {
		recordDryRunStats(info, stats)
		return
	}

	removeEntry(fullPath, info, stats)
}

func recordDryRunStats(info os.FileInfo, stats *CleanStats) {
	if info.IsDir() {
		stats.RemovedDirsCount++
		return
	}
	stats.RemovedFilesCount++
	stats.FreedBytes += info.Size()
}

func removeEntry(fullPath string, info os.FileInfo, stats *CleanStats) {
	if info.IsDir() {
		removeDirectoryEntry(fullPath, stats)
		return
	}
	removeFileEntry(fullPath, info, stats)
}

func removeDirectoryEntry(fullPath string, stats *CleanStats) {
	if err := os.RemoveAll(fullPath); err == nil {
		stats.RemovedDirsCount++
	}
}

func removeFileEntry(fullPath string, info os.FileInfo, stats *CleanStats) {
	if err := os.Remove(fullPath); err == nil {
		stats.RemovedFilesCount++
		stats.FreedBytes += info.Size()
	}
}
