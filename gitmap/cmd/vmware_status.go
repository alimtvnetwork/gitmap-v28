package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runVmwareSharedStatus() error {
	fmt.Println("▶ gitmap vmware shared status")
	active := isMountActive(defaultMountPoint)
	fmt.Printf("  Mount (%s): active=%t\n", defaultMountPoint, active)

	desktopDir := resolveUserDesktopDir()
	link := filepath.Join(desktopDir, "SharedDirectories")
	target, err := os.Readlink(link)
	hasLink := err == nil
	fmt.Printf("  Desktop Symlink: present=%t (target=%s)\n", hasLink, target)

	return nil
}

func isMountActive(mountPoint string) bool {
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return false
	}

	return strings.Contains(string(data), mountPoint)
}

func runVmwareStatus(args []string) error {
	fmt.Println("▶ gitmap vmware status")
	fmt.Printf("  Linux Guest: %t\n", isLinuxOS())
	fmt.Printf("  VMware Hypervisor: %t\n", isVMwareHypervisor())
	fmt.Printf("  Mount Point: %s (active=%t)\n", defaultMountPoint, isMountActive(defaultMountPoint))

	return nil
}
