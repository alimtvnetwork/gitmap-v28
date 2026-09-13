//go:build !windows

package cmdinstall

import (
	"os"
	"os/exec"
	"path/filepath"
)

func cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir string) {
	_ = os.Remove(filepath.Join(binDir, "antigravity"))
	_ = os.Remove(filepath.Join(binDir, "agy"))
	_ = os.Remove(filepath.Join(desktopDir, "antigravity.desktop"))
	_ = os.RemoveAll(installDir)
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
	return "antigravity"
}
