//go:build windows

package cmd

import (
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
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
	volRoot := resolveVolumeRoot(path)
	volPtr, err := windows.UTF16PtrFromString(volRoot)
	if err != nil {
		return nil, err
	}

	return calculateDiskMetrics(volRoot, volPtr)
}

func resolveVolumeRoot(path string) string {
	vol := filepath.VolumeName(path)
	if vol == "" {
		return "C:\\"
	}

	return vol + "\\"
}

func calculateDiskMetrics(volRoot string, volPtr *uint16) (*DiskSpaceInfo, error) {
	var freeBytes, totalBytes, totalFreeBytes uint64
	err := windows.GetDiskFreeSpaceEx(volPtr, &freeBytes, &totalBytes, &totalFreeBytes)
	if err != nil {
		return nil, err
	}

	usedBytes := totalBytes - freeBytes
	var pct float64
	if totalBytes > 0 {
		pct = (float64(usedBytes) / float64(totalBytes)) * 100.0
	}

	return &DiskSpaceInfo{
		DrivePath:   volRoot,
		Filesystem:  queryFilesystemName(volPtr),
		TotalBytes:  totalBytes,
		UsedBytes:   usedBytes,
		FreeBytes:   freeBytes,
		UsedPercent: pct,
	}, nil
}

func queryFilesystemName(volPtr *uint16) string {
	var fsBuf [256]uint16
	modKernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetVolumeInfo := modKernel32.NewProc("GetVolumeInformationW")

	ret, _, _ := procGetVolumeInfo.Call(
		uintptr(unsafe.Pointer(volPtr)),
		0, 0, 0, 0, 0,
		uintptr(unsafe.Pointer(&fsBuf[0])),
		uintptr(len(fsBuf)),
	)
	if ret == 0 {
		return "NTFS"
	}

	return syscall.UTF16ToString(fsBuf[:])
}
