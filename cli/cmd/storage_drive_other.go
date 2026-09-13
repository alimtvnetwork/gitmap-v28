//go:build !windows

package cmd

import (
	"golang.org/x/sys/unix"
)

// DiskSpaceInfo encapsulates disk drive capacity metrics.
type DiskSpaceInfo struct {
	DrivePath   string
	Filesystem  string
	TotalBytes  uint64
	UsedBytes   uint64
	FreeBytes   uint64
	UsedPercent float64
}

func getDiskSpaceMetrics(path string) (*DiskSpaceInfo, error) {
	target := resolveUnixPath(path)
	var stat unix.Statfs_t
	if err := unix.Statfs(target, &stat); err != nil {
		return nil, err
	}

	return calculateUnixMetrics(target, stat), nil
}

func resolveUnixPath(path string) string {
	if path == "" {
		return "/"
	}

	return path
}

func calculateUnixMetrics(target string, stat unix.Statfs_t) *DiskSpaceInfo {
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free
	var pct float64
	if total > 0 {
		pct = (float64(used) / float64(total)) * 100.0
	}

	return &DiskSpaceInfo{
		DrivePath:   target,
		Filesystem:  "POSIX",
		TotalBytes:  total,
		UsedBytes:   used,
		FreeBytes:   free,
		UsedPercent: pct,
	}
}
