//go:build !windows

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func sweepIconResolutions(hicolorDir, iconName string) {
	for _, size := range StandardIconResolutions {
		iconFile := filepath.Join(hicolorDir, fmt.Sprintf("%dx%d", size, size), "apps", iconName+".png")
		_ = os.Remove(iconFile)
	}

	runGtkUpdateIconCache(hicolorDir)
}

func sweepPixmaps(iconName string) {
	home, _ := os.UserHomeDir()
	_ = os.Remove(filepath.Join(home, ".local", "share", "pixmaps", iconName+".png"))
	_ = os.Remove(filepath.Join("/usr/share/pixmaps", iconName+".png"))
}

func sweepAntigravityIcons() {
	home, _ := os.UserHomeDir()
	userHicolor := filepath.Join(home, ".local", "share", "icons", "hicolor")
	sweepIconResolutions(userHicolor, AntigravityCanonicalIconName)
	sweepIconResolutions("/usr/share/icons/hicolor", AntigravityCanonicalIconName)
	sweepPixmaps(AntigravityCanonicalIconName)
}

func cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir string) {
	_ = os.Remove(filepath.Join(binDir, "antigravity"))
	_ = os.Remove(filepath.Join(binDir, "agy"))
	_ = os.Remove(filepath.Join(desktopDir, "antigravity.desktop"))
	_ = os.RemoveAll(installDir)
	sweepAntigravityIcons()
	PurgeDuplicateAntigravityLaunchers()
	updateDesktopDatabase(desktopDir)
}

func updateDesktopDatabase(desktopDir string) {
	if bin, err := exec.LookPath("update-desktop-database"); err == nil {
		_ = exec.Command(bin, desktopDir).Run()
	}
}

func resolveAppIconPath(installDir string) string {
	candidates := []string{
		filepath.Join(installDir, "resources", "app", "resources", "icon.png"),
		filepath.Join(installDir, "resources", "icon.png"),
		filepath.Join(installDir, "icon.png"),
	}

	for _, p := range candidates {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}

	return AntigravityCanonicalIconName
}
