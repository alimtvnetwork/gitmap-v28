//go:build windows

package cmdvmware

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"golang.org/x/sys/windows/registry"
)

func init() {
	executeVmwareInstallPlatformFn = executeVmwareInstallWindows
	executeVmwareSharedEnablePlatformFn = executeVmwareSharedEnableWindows
	runVmwareSharedStatusPlatformFn = runVmwareSharedStatusWindows
	runVmwareStatusPlatformFn = runVmwareStatusWindows
}

// isVMwareGuestWindows checks registry and SCM for VMware Tools presence.
func isVMwareGuestWindows() bool {
	if hasVMwareRegistryKey() {
		return true
	}

	return hasVMwareService()
}

func hasVMwareRegistryKey() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\VMware, Inc.\VMware Tools`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}

	_ = k.Close()

	return true
}

func hasVMwareService() bool {
	out, err := exec.Command("sc", "query", "VMTools").Output()
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToUpper(string(out)), "VMTOOLS")
}

// getVMToolsServiceStatusWindows queries VMTools service state.
func getVMToolsServiceStatusWindows() (bool, string) {
	out, err := exec.Command("sc", "query", "VMTools").Output()
	if err != nil {
		return false, "NOT_FOUND"
	}

	text := string(out)
	if strings.Contains(text, "RUNNING") {
		return true, "RUNNING"
	}

	if strings.Contains(text, "STOPPED") {
		return false, "STOPPED"
	}

	return false, "UNKNOWN"
}

// findVmrunPath locates vmrun.exe in PATH and standard program files directories.
func findVmrunPath() string {
	if p, err := exec.LookPath("vmrun.exe"); err == nil {
		return p
	}

	if p, err := exec.LookPath("vmrun"); err == nil {
		return p
	}

	return searchStandardVmrunPaths()
}

func searchStandardVmrunPaths() string {
	for _, p := range buildVmrunCandidatePaths() {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func buildVmrunCandidatePaths() []string {
	var paths []string
	progFiles := os.Getenv("ProgramFiles")
	progFilesX86 := os.Getenv("ProgramFiles(x86)")
	for _, base := range []string{progFiles, progFilesX86} {
		if base != "" {
			paths = append(paths, buildCandidatePathsForBase(base)...)
		}
	}

	return paths
}

func buildCandidatePathsForBase(base string) []string {
	return []string{
		filepath.Join(base, "VMware", "VMware Workstation", "vmrun.exe"),
		filepath.Join(base, "VMware", "VMware Player", "vmrun.exe"),
		filepath.Join(base, "VMware", "VMware VIX", "vmrun.exe"),
	}
}

// isWindowsSharedFolderAvailable checks access to \\vmware-host\Shared Folders.
func isWindowsSharedFolderAvailable() bool {
	_, err := os.Stat(`\\vmware-host\Shared Folders`)

	return err == nil
}

func executeVmwareInstallWindows() error {
	fmt.Println("▶ gitmap vmware install (Windows)")
	isGuest := isVMwareGuestWindows()
	hasService, serviceStatus := getVMToolsServiceStatusWindows()
	vmrunPath := findVmrunPath()

	fmt.Printf("  VMware Guest: %t\n", isGuest)
	fmt.Printf("  VMware Tools Service: %s (running=%t)\n", serviceStatus, hasService)
	printVmrunPath(vmrunPath)

	return evaluateWindowsToolsStatus(isGuest, hasService)
}

func printVmrunPath(vmrunPath string) {
	if vmrunPath != "" {
		fmt.Printf("  vmrun utility: %s\n", vmrunPath)

		return
	}

	fmt.Println("  vmrun utility: not found")
}

func evaluateWindowsToolsStatus(isGuest, hasService bool) error {
	if isGuest && hasService {
		fmt.Println("  ✓ VMware Tools is installed and active on Windows guest")

		return nil
	}

	if isGuest {
		fmt.Println("  ⚠ VMware Tools detected but service is not currently running")

		return nil
	}

	fmt.Println("  ⚠ VMware Tools not detected on this system")
	fmt.Println("  To install: VMware menu -> VM -> Install VMware Tools")

	return nil
}

func executeVmwareSharedEnableWindows(isDryRun bool) error {
	fmt.Println("▶ gitmap vmware shared enable (Windows)")
	uncPath := `\\vmware-host\Shared Folders`
	if isDryRun {
		return simulateWindowsSharedEnable(uncPath)
	}

	isAvailable := isWindowsSharedFolderAvailable()
	reportUNCStatus(uncPath, isAvailable)

	return createWindowsDesktopShortcut(uncPath, "VMware Shared Folders")
}

func simulateWindowsSharedEnable(uncPath string) error {
	fmt.Printf("  [dry-run] Would verify UNC share access: %s\n", uncPath)
	fmt.Println("  [dry-run] Would create Desktop shortcut to Shared Folders")

	return nil
}

func reportUNCStatus(uncPath string, isAvailable bool) {
	if isAvailable {
		fmt.Printf("  ✓ UNC path %s is accessible\n", uncPath)

		return
	}

	fmt.Printf("  ⚠ UNC path %s not currently accessible\n", uncPath)
	fmt.Println("  Please verify Shared Folders in VMware Settings -> Options -> Shared Folders")
}

func createWindowsDesktopShortcut(target, shortcutName string) error {
	desktopDir := resolveWindowsDesktopDir()
	if err := os.MkdirAll(desktopDir, 0755); err != nil {
		return apperror.WrapSimple(err, "vmware.createDesktopShortcut.mkdir")
	}

	lnkPath := filepath.Join(desktopDir, shortcutName+".lnk")
	script := fmt.Sprintf("$w=New-Object -ComObject WScript.Shell;$s=$w.CreateShortcut('%s');$s.TargetPath='%s';$s.Save()",
		lnkPath, target)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	if err := cmd.Run(); err != nil {
		return createFallbackUrlShortcut(desktopDir, shortcutName, target)
	}

	fmt.Println("  ✓ Created Desktop shortcut to VMware Shared Folders")

	return nil
}

func createFallbackUrlShortcut(desktopDir, shortcutName, target string) error {
	urlPath := filepath.Join(desktopDir, shortcutName+".url")
	content := fmt.Sprintf("[InternetShortcut]\r\nURL=file:%s\r\n", target)
	if err := os.WriteFile(urlPath, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "vmware.createFallbackUrlShortcut")
	}

	fmt.Println("  ✓ Created Desktop network shortcut to VMware Shared Folders")

	return nil
}

func resolveWindowsDesktopDir() string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "Desktop")
	}

	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		return filepath.Join(userProfile, "Desktop")
	}

	return "."
}

func runVmwareSharedStatusWindows() error {
	fmt.Println("▶ gitmap vmware shared status (Windows)")
	uncPath := `\\vmware-host\Shared Folders`
	isAvailable := isWindowsSharedFolderAvailable()
	fmt.Printf("  UNC Shared Folders (%s): accessible=%t\n", uncPath, isAvailable)

	desktopDir := resolveWindowsDesktopDir()
	lnkPath := filepath.Join(desktopDir, "VMware Shared Folders.lnk")
	_, errLnk := os.Stat(lnkPath)
	fmt.Printf("  Desktop Shortcut: present=%t\n", errLnk == nil)

	return nil
}

func runVmwareStatusWindows() error {
	fmt.Println("▶ gitmap vmware status")
	fmt.Println("  OS Platform: Windows")
	isGuest := isVMwareGuestWindows()
	fmt.Printf("  VMware Guest: %t\n", isGuest)
	hasService, serviceStatus := getVMToolsServiceStatusWindows()
	fmt.Printf("  VMware Tools Service: %s (running=%t)\n", serviceStatus, hasService)
	vmrunPath := findVmrunPath()
	fmt.Printf("  vmrun Installed: %t (path=%s)\n", vmrunPath != "", vmrunPath)
	fmt.Printf("  Shared Folders Accessible: %t\n", isWindowsSharedFolderAvailable())

	return nil
}
