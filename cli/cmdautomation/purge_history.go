package cmdautomation

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type purgeTotals struct {
	totalBytes  int64
	purgedCount int
}

// RunPurgeHistory traces and identifies large historical git blobs.
func RunPurgeHistory(opts PurgeHistoryOptions) PurgeHistoryResultMonad {
	start := time.Now()
	baseDir := resolveArtifactBaseDir(opts.Dir)
	items := discoverHistoryBlobs(opts, baseDir)
	purged := executeHistoryPurge(items, opts, baseDir)
	res := aggregatePurgeHistory(purged, start)
	return result.Ok(res)
}

func resolveMinSizeBytes(minMB float64) int64 {
	hasMin := minMB > 0
	if hasMin {
		return int64(minMB * 1024 * 1024)
	}
	return 1024 * 1024
}

func extractBlobPath(parts []string) string {
	if len(parts) >= 4 {
		return filepath.ToSlash(strings.TrimSpace(parts[3]))
	}
	return "unknown"
}

func createPurgeItem(hash string, path string, size int64) PurgeHistoryItem {
	return PurgeHistoryItem{
		Path:       path,
		CommitHash: hash,
		SizeBytes:  size,
		SizeMB:     float64(size) / (1024 * 1024),
		IsPurged:   false,
	}
}

func parseBlobLine(line string, minBytes int64, pathPattern string) PurgeHistoryItem {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 3 {
		return PurgeHistoryItem{}
	}
	objType := parts[1]
	if objType != "blob" {
		return PurgeHistoryItem{}
	}
	size, _ := strconv.ParseInt(parts[2], 10, 64)
	if size < minBytes {
		return PurgeHistoryItem{}
	}
	path := extractBlobPath(parts)
	if len(pathPattern) > 0 && strings.Contains(path, pathPattern) == false {
		return PurgeHistoryItem{}
	}
	return createPurgeItem(parts[0], path, size)
}

func parseBlobCheckOutput(batchOutput string, minBytes int64, pathPattern string) []PurgeHistoryItem {
	lines := strings.Split(batchOutput, "\n")
	items := make([]PurgeHistoryItem, 0, len(lines))
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if len(trimmed) == 0 {
			continue
		}
		it := parseBlobLine(trimmed, minBytes, pathPattern)
		if it.SizeBytes > 0 {
			items = append(items, it)
		}
	}
	return items
}

func fetchObjectsList(baseDir string) string {
	cmd := exec.Command("git", "rev-list", "--objects", "--all")
	cmd.Dir = baseDir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func checkObjectBatch(baseDir string, objects string) string {
	if len(objects) == 0 {
		return ""
	}
	formatArg := "--batch-check=%(objectname) %(objecttype) %(objectsize) %(rest)"
	cmd := exec.Command("git", "cat-file", formatArg)
	cmd.Dir = baseDir
	cmd.Stdin = strings.NewReader(objects)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func discoverHistoryBlobs(opts PurgeHistoryOptions, baseDir string) []PurgeHistoryItem {
	minBytes := resolveMinSizeBytes(opts.MinSizeMB)
	objects := fetchObjectsList(baseDir)
	batch := checkObjectBatch(baseDir, objects)
	return parseBlobCheckOutput(batch, minBytes, opts.PathPattern)
}

func createSafetyBackup(baseDir string) string {
	tag := time.Now().Format("20060102-150405")
	branch := fmt.Sprintf("backup/history-purge-%s", tag)
	cmd := exec.Command("git", "branch", branch)
	cmd.Dir = baseDir
	_ = cmd.Run()
	return branch
}

func executeHistoryPurge(items []PurgeHistoryItem, opts PurgeHistoryOptions, baseDir string) []PurgeHistoryItem {
	if opts.IsDryRun || len(items) == 0 {
		return items
	}
	_ = createSafetyBackup(baseDir)
	return items
}

func computePurgeTotals(items []PurgeHistoryItem) purgeTotals {
	var tot purgeTotals
	for _, it := range items {
		tot.totalBytes += it.SizeBytes
		if it.IsPurged {
			tot.purgedCount++
		}
	}
	return tot
}

func aggregatePurgeHistory(items []PurgeHistoryItem, start time.Time) PurgeHistoryResult {
	tot := computePurgeTotals(items)
	return PurgeHistoryResult{
		TotalFound:  len(items),
		TotalPurged: tot.purgedCount,
		TotalBytes:  tot.totalBytes,
		TotalMB:     float64(tot.totalBytes) / (1024 * 1024),
		Items:       items,
		Duration:    time.Since(start),
		IsSuccess:   true,
	}
}
