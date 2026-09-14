package macro

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func isSudoActive() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	sudoUser := os.Getenv("SUDO_USER")

	return sudoUser != "" && sudoUser != "root"
}

func restoreSudoOwnership(targetPath string) {
	if !isSudoActive() {
		return
	}
	sudoUser := os.Getenv("SUDO_USER")
	chownPath, lookErr := exec.LookPath("chown")
	if lookErr != nil {
		return
	}
	cmd := exec.Command(chownPath, "-R", sudoUser, targetPath)
	_ = cmd.Run()
}

func attemptPermissionHealing(targetPath string) {
	if runtime.GOOS == "windows" {
		return
	}
	_ = os.Chmod(targetPath, 0777)
	parent := filepath.Dir(targetPath)
	if parent != "" && parent != targetPath {
		_ = os.Chmod(parent, 0777)
	}
}
