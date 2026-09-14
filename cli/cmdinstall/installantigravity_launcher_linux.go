//go:build linux

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AntigravityLauncherDeployConfig encapsulates parameters for canonical Antigravity launcher deployment.
type AntigravityLauncherDeployConfig struct {
	BinPath      string
	UserAppDir   string
	DesktopDir   string
	IconName     string
	IsPrivileged bool
}

// AntigravityStaleLauncherNames lists all stale and duplicate Antigravity desktop launcher file names.
var AntigravityStaleLauncherNames = []string{
	"antigravity-ide.desktop",
	"Google Antigravity.desktop",
	"Google-Antigravity.desktop",
	"Antigravity.desktop",
	"antigravity.desktop",
}

func isDirWritable(dir string) bool {
	probe := filepath.Join(dir, fmt.Sprintf(".probe_%d", time.Now().UnixNano()))
	if err := os.WriteFile(probe, []byte("1"), 0644); err != nil {
		return false
	}
	_ = os.Remove(probe)

	return true
}

func checkDirExists(dir string) bool {
	if dir == "" {
		return false
	}
	info, err := os.Stat(dir)

	return err == nil && info.IsDir()
}

func resolveUserAppsDir() string {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".local", "share", "applications")
}

func resolveUserDesktopDir() string {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, "Desktop")
}

func purgeStaleLaunchersFromDir(dir string) {
	for _, name := range AntigravityStaleLauncherNames {
		_ = os.Remove(filepath.Join(dir, name))
	}
}

func purgeSystemStaleLaunchers() {
	systemDir := "/usr/share/applications"
	if isDirWritable(systemDir) {
		purgeStaleLaunchersFromDir(systemDir)
	}
}

// PurgeDuplicateAntigravityLaunchers removes all stale Antigravity desktop launchers.
func PurgeDuplicateAntigravityLaunchers() {
	purgeStaleLaunchersFromDir(resolveUserAppsDir())
	purgeStaleLaunchersFromDir(resolveUserDesktopDir())
	purgeSystemStaleLaunchers()
}

func markDesktopFileTrusted(desktopFilePath string) {
	gioBin, err := exec.LookPath("gio")
	if err != nil {
		return
	}
	_ = exec.Command(gioBin, "set", desktopFilePath, "metadata::trusted", "true").Run()
}

func buildAntigravityDesktopContent(binPath, iconName string) string {
	resolvedIcon := iconName
	if resolvedIcon == "" {
		resolvedIcon = AntigravityCanonicalIconName
	}

	return fmt.Sprintf("[Desktop Entry]\nName=Google Antigravity\nComment=AI-First Development Platform & Agent Orchestration IDE\nGenericName=Text Editor / IDE\nExec=%s %%U\nIcon=%s\nType=Application\nStartupNotify=true\nStartupWMClass=Antigravity\nCategories=Development;IDE;Utility;\n", binPath, resolvedIcon)
}

func writeDesktopFileWithChmod(destPath, content string) error {
	if err := os.WriteFile(destPath, []byte(content), 0755); err != nil {
		return apperror.WrapSimple(err, "write.desktopFile")
	}
	if err := os.Chmod(destPath, 0755); err != nil {
		return apperror.WrapSimple(err, "chmod.desktopFile")
	}

	return nil
}

func copyFileWithChmod(srcPath, destPath string, mode os.FileMode) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return apperror.WrapSimple(err, "read.desktopFile")
	}
	if err := os.WriteFile(destPath, data, mode); err != nil {
		return apperror.WrapSimple(err, "write.desktopFileCopy")
	}
	_ = os.Chmod(destPath, mode)

	return nil
}

func deployToUserDesktop(srcFile, desktopDir string) {
	if !checkDirExists(desktopDir) {
		return
	}
	destFile := filepath.Join(desktopDir, "antigravity.desktop")
	_ = copyFileWithChmod(srcFile, destFile, 0755)
	markDesktopFileTrusted(destFile)
}

func deployToSystemApplications(srcFile string, isPrivileged bool) {
	if !isPrivileged {
		return
	}
	destFile := filepath.Join("/usr/share/applications", "antigravity.desktop")
	_ = copyFileWithChmod(srcFile, destFile, 0755)
}

func touchDesktopFile(filePath string) {
	now := time.Now()
	_ = os.Chtimes(filePath, now, now)
}

func touchAntigravityLaunchers(userFile, desktopDir string) {
	touchDesktopFile(userFile)
	if checkDirExists(desktopDir) {
		touchDesktopFile(filepath.Join(desktopDir, "antigravity.desktop"))
	}
}

func refreshDesktopLauncherDatabases(userAppDir string, isPrivileged bool) {
	updateDesktopDatabase(userAppDir)
	if isPrivileged {
		updateDesktopDatabase("/usr/share/applications")
	}
}

// DeployCanonicalAntigravityLauncher installs canonical Antigravity desktop launcher across user and system targets.
func DeployCanonicalAntigravityLauncher(cfg AntigravityLauncherDeployConfig) error {
	PurgeDuplicateAntigravityLaunchers()
	userFile := filepath.Join(cfg.UserAppDir, "antigravity.desktop")
	content := buildAntigravityDesktopContent(cfg.BinPath, cfg.IconName)
	if err := writeDesktopFileWithChmod(userFile, content); err != nil {
		return err
	}
	deployToUserDesktop(userFile, cfg.DesktopDir)
	deployToSystemApplications(userFile, cfg.IsPrivileged)
	refreshDesktopLauncherDatabases(cfg.UserAppDir, cfg.IsPrivileged)
	touchAntigravityLaunchers(userFile, cfg.DesktopDir)

	return nil
}
