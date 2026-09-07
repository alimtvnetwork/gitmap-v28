package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	defaultMountPoint = "/mnt/hgfs"
	crontabRebootLine = "@reboot /usr/bin/vmhgfs-fuse -o allow_other -o auto_unmount .host:/ /mnt/hgfs"
)

func runVmwareShared(args []string) error {
	if len(args) == 0 {
		return runVmwareSharedStatus()
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case constants.SubCmdSharedEnable, "mount":
		return runVmwareSharedEnable(args[1:])
	case constants.SubCmdSharedStatus:
		return runVmwareSharedStatus()
	default:
		return unknownVmwareSubcommandError(sub)
	}
}

func isLinuxOS() bool {
	return runtime.GOOS == "linux"
}

func isVMwareHypervisor() bool {
	paths := []string{
		"/sys/class/dmi/id/sys_vendor",
		"/sys/class/dmi/id/product_name",
	}

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil && strings.Contains(strings.ToLower(string(data)), "vmware") {
			return true
		}
	}

	out, err := exec.Command("systemd-detect-virt").Output()

	return err == nil && strings.Contains(strings.ToLower(string(out)), "vmware")
}

func checkVMwarePrerequisites() error {
	if !isLinuxOS() {
		return apperror.NewWithDetails("cmd.vmware", "E4002", "vmware features are only supported on Linux guest environments", "cmd.vmware", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
	}

	if !isVMwareHypervisor() {
		fmt.Println("  ⚠ Warning: VMware hypervisor not detected; continuing per user request...")
	}

	return nil
}

func resolveUserDesktopDir() string {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" && sudoUser != "root" {
		return filepath.Join("/home", sudoUser, "Desktop")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/root", "Desktop")
	}

	return filepath.Join(home, "Desktop")
}

func ensureMountDirectory(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	cmd := exec.Command("sudo", "mkdir", "-p", dir)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "vmware.ensureMountDirectory")
	}

	return nil
}

func mountHostShare(mountPoint string) error {
	cmd := exec.Command("sudo", "vmhgfs-fuse", "-o", "allow_other", "-o", "auto_unmount", ".host:/", mountPoint)
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperror.NewWithDetails("cmd.vmware.mountHostShare", "E4003", fmt.Sprintf("mount failed: %s (%v)", string(out), err), "cmd.vmware", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}

func createDesktopSymlink(mountPoint string) error {
	desktopDir := resolveUserDesktopDir()
	if err := os.MkdirAll(desktopDir, 0755); err != nil {
		return apperror.WrapSimple(err, "vmware.createDesktopSymlink.mkdir")
	}

	linkPath := filepath.Join(desktopDir, "SharedDirectories")
	_ = os.Remove(linkPath)

	if err := os.Symlink(mountPoint, linkPath); err != nil {
		return apperror.WrapSimple(err, "vmware.createDesktopSymlink.symlink")
	}

	return nil
}

func ensureCrontabPersistence() error {
	out, _ := exec.Command("crontab", "-l").CombinedOutput()
	current := string(out)
	if strings.Contains(current, "vmhgfs-fuse") && strings.Contains(current, defaultMountPoint) {
		return nil
	}

	newLine := crontabRebootLine + "\n"
	updated := current + newLine
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(updated)
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperror.NewWithDetails("cmd.vmware.crontab", "E4004", fmt.Sprintf("failed updating crontab: %s (%v)", string(out), err), "cmd.vmware", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}

func runVmwareSharedEnable(args []string) error {
	if err := checkVMwarePrerequisites(); err != nil {
		return err
	}

	fmt.Println("▶ gitmap vmware shared enable")
	if err := ensureMountDirectory(defaultMountPoint); err != nil {
		return err
	}
	fmt.Printf("  ✓ Verified mount point %s\n", defaultMountPoint)

	if err := mountHostShare(defaultMountPoint); err != nil {
		return err
	}
	fmt.Printf("  ✓ Mounted .host:/ at %s\n", defaultMountPoint)

	if err := createDesktopSymlink(defaultMountPoint); err != nil {
		return err
	}
	fmt.Printf("  ✓ Created Desktop/SharedDirectories symlink\n")

	if err := ensureCrontabPersistence(); err != nil {
		return err
	}
	fmt.Printf("  ✓ Registered @reboot crontab persistence\n")

	return nil
}

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
