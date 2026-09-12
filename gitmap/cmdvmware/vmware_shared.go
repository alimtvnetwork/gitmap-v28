package cmdvmware

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

	ensureVMwareToolsInstalled()

	return nil
}

func ensureVMwareToolsInstalled() {
	if _, err := exec.LookPath("vmhgfs-fuse"); err == nil {
		return
	}

	fmt.Println("  ⚠ vmhgfs-fuse not found. Attempting to install open-vm-tools via apt...")
	if _, err := exec.LookPath("apt-get"); err != nil {
		return
	}

	cmd := exec.Command("sudo", "apt-get", "install", "-y", "open-vm-tools", "open-vm-tools-desktop")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

// ResolveUserDesktopDir returns the user Desktop directory accounting for SUDO_USER.
func ResolveUserDesktopDir() string {
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

func unmountIfMounted(mountPoint string) {
	if !isMountActive(mountPoint) {
		return
	}

	_ = exec.Command("fusermount", "-u", mountPoint).Run()
	_ = exec.Command("sudo", "umount", "-l", mountPoint).Run()
}

func ensureVMwareServiceRunning() {
	if _, err := exec.LookPath("systemctl"); err != nil {
		return
	}

	_ = exec.Command("sudo", "systemctl", "start", "open-vm-tools").Run()
}

func executePrimaryMount(mountPoint string) ([]byte, error) {
	cmd := exec.Command("sudo", "vmhgfs-fuse", "-o", "allow_other", "-o", "auto_unmount", ".host:/", mountPoint)

	return cmd.CombinedOutput()
}

func tryFallbackMount(mountPoint string) ([]byte, error) {
	cmd := exec.Command("sudo", "mount", "-t", "fuse.vmhgfs-fuse", ".host:/", mountPoint, "-o", "allow_other")

	return cmd.CombinedOutput()
}

func isFallbackMountSuccess(mountPoint string) bool {
	_, fbErr := tryFallbackMount(mountPoint)

	return fbErr == nil
}

func getHostShares() []string {
	out, err := exec.Command("vmware-hgfsclient").CombinedOutput()
	if err != nil || len(out) == 0 {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var shares []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			shares = append(shares, trimmed)
		}
	}

	return shares
}

func buildMountDiagnosticHelp(errText string, runErr error) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("mount failed: %s (%v)\n", errText, runErr))
	sb.WriteString("\n  [Diagnostic & Remediation]\n")
	sb.WriteString("  → Host has not enabled Shared Folders or no shares are configured in VMware.\n")
	sb.WriteString("  → To fix in VMware Workstation / Player / Fusion:\n")
	sb.WriteString("    1. Open VM -> Settings -> Options -> Shared Folders\n")
	sb.WriteString("    2. Change setting to 'Always enabled'\n")
	sb.WriteString("    3. Click 'Add...' and configure at least one folder path from the host\n")
	sb.WriteString("    4. Click OK / Save, then re-run: gitmap vmware shared enable\n")

	return sb.String()
}

func formatVmwareMountError(out []byte, runErr error) string {
	errText := strings.TrimSpace(string(out))
	if errText == "" && runErr != nil {
		errText = runErr.Error()
	}

	hasConnErr := strings.Contains(errText, "-107") || strings.Contains(errText, "cannot open connection")
	hasShares := len(getHostShares()) > 0
	if hasConnErr || !hasShares {
		return buildMountDiagnosticHelp(errText, runErr)
	}

	return fmt.Sprintf("mount failed: %s (%v)", errText, runErr)
}

func newMountHostError(out []byte, err error) error {
	msg := formatVmwareMountError(out, err)

	return apperror.NewWithDetails("cmd.vmware.mountHostShare", "E4003", msg, "cmd.vmware", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
}

func mountHostShare(mountPoint string) error {
	unmountIfMounted(mountPoint)
	ensureVMwareServiceRunning()
	out, err := executePrimaryMount(mountPoint)
	if err == nil {
		return nil
	}

	if isFallbackMountSuccess(mountPoint) {
		return nil
	}

	return newMountHostError(out, err)
}

func createDesktopSymlink(mountPoint string) error {
	desktopDir := ResolveUserDesktopDir()
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

func runVmwareSharedEnable(args []string) error {
	checkHelp(constants.CmdVmware, args)
	isDryRun := hasDryRunFlag(args) || hasShortDryRunFlag(args)
	if err := checkVMwarePrerequisites(); err != nil {
		return err
	}

	fmt.Println("▶ gitmap vmware shared enable")
	if isDryRun {
		fmt.Printf("  [dry-run] Would verify mount point %s\n", defaultMountPoint)
		fmt.Printf("  [dry-run] Would mount .host:/ at %s (vmhgfs-fuse)\n", defaultMountPoint)
		fmt.Printf("  [dry-run] Would create Desktop/SharedDirectories symlink -> %s\n", defaultMountPoint)
		fmt.Printf("  [dry-run] Would register @reboot crontab persistence\n")

		return nil
	}

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

	if err := EnsureCrontabPersistence(); err != nil {
		return err
	}

	fmt.Printf("  ✓ Registered @reboot crontab persistence\n")

	return nil
}

func hasShortDryRunFlag(args []string) bool {
	for _, a := range args {
		if a == "-n" {
			return true
		}
	}

	return false
}
